package httprequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/spf13/viper"
	"mlib.com/gofy/server/core/exceptions"
	httprequestnodesentities "mlib.com/gofy/server/entities/nodes/http_request"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	"mlib.com/mlog"
)

var (
	BODY_TYPE_TO_CONTENT_TYPE = map[string]string{
		"json":                  "application/json",
		"x-www-form-urlencoded": "application/x-www-form-urlencoded",
		"form-data":             "multipart/form-data",
		"raw-text":              "text/plain",
	}
)

type Executor struct {
	Method  string
	URL     string
	Params  [][2]string
	Content any //string // | bytes | None
	Data    map[string]string
	// files: Mapping[str, tuple[str | None, bytes, str]] | None
	Json       any
	Headers    map[string]string
	Auth       *httprequestnodesentities.HttpRequestNodeAuthorization
	Timeout    *httprequestnodesentities.HttpRequestNodeTimeout
	MaxRetries int

	Boundary string

	variable_pool *workflowentities.VariablePool
	node_data     *httprequestnodesentities.HttpRequestNodeData
}

func NewExecutor(
	node_data *httprequestnodesentities.HttpRequestNodeData,
	timeout *httprequestnodesentities.HttpRequestNodeTimeout,
	variable_pool *workflowentities.VariablePool,
	max_retries int, /* dify_config.SSRF_DEFAULT_MAX_RETRIES*/
) *Executor {
	// If authorization API key is present, convert the API key using the variable pool
	if node_data.Authorization != nil && node_data.Authorization.Type == "api-key" {
		if node_data.Authorization.Config == nil {
			panic(exceptions.NewAuthorizationConfigError("authorization config is required"))
		}
		node_data.Authorization.Config.APIKey = variable_pool.ConvertTemplate(node_data.Authorization.Config.APIKey).Text()
	}
	if !slices.Contains(httprequestnodesentities.METHODS, node_data.Method) {
		panic(exceptions.NewInvalidHttpMethodError("method=" + node_data.Method))
	}
	e := &Executor{
		URL:        node_data.URL,
		Method:     node_data.Method,
		Auth:       node_data.Authorization,
		Timeout:    timeout,
		Params:     make([][2]string, 0),
		Headers:    make(map[string]string),
		MaxRetries: max_retries,

		// init template
		variable_pool: variable_pool,
		node_data:     node_data,
	}
	e.initialize()

	return e
}
func (e *Executor) initialize() {
	e.init_url()
	e.init_params()
	e.init_headers()
	e.init_body()
}
func (e *Executor) init_url() {
	e.URL = e.variable_pool.ConvertTemplate(e.node_data.URL).Text()

	// check if url is a valid URL
	if e.URL == "" {
		panic(exceptions.NewInvalidURLError("url is required"))
	}
	if !strings.HasPrefix(e.URL, "http://") && !strings.HasPrefix(e.URL, "https://") {
		panic(exceptions.NewInvalidURLError("url should start with http:// or https://"))
	}
}
func (e *Executor) init_params() {
	/*
	   Almost same as _init_headers(), difference:
	   1. response a list tuple to support same key, like 'aa=1&aa=2'
	   2. param value may have '\n', we need to splitlines then extract the variable value.
	*/
	result := [][2]string{}

	for line := range strings.SplitSeq(e.node_data.Params, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		mlog.Debug("---line=", line)
		tmplist := strings.SplitN(line, ":", 2)
		if len(tmplist) != 2 {
			mlog.Error("invalid line=", line)
			continue
		}
		mlog.Debug("---tmplist=", tmplist)
		// key, *value = line.split(":", 1)
		if strings.TrimSpace(tmplist[0]) == "" {
			continue
		}
		value_str := strings.TrimSpace(tmplist[1])
		result = append(result, [2]string{e.variable_pool.ConvertTemplate(strings.TrimSpace(tmplist[0])).Text(), e.variable_pool.ConvertTemplate(value_str).Text()})
	}
	e.Params = result
}
func (e *Executor) init_headers() {
	/*
	   Convert the header string of frontend to a dictionary.

	   Each line in the header string represents a key-value pair.
	   Keys and values are separated by ':'.
	   Empty values are allowed.

	   Examples:
	       'aa:bb\n cc:dd'  -> {'aa': 'bb', 'cc': 'dd'}
	       'aa:\n cc:dd\n'  -> {'aa': '', 'cc': 'dd'}
	       'aa\n cc : dd'   -> {'aa': '', 'cc': 'dd'}

	*/
	e.Headers = map[string]string{}
	headers := e.variable_pool.ConvertTemplate(e.node_data.Headers).Text()
	for line := range strings.SplitSeq(headers, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		mlog.Debug("---line=", line)
		tmplist := strings.SplitN(line, ":", 2)
		if len(tmplist) != 2 {
			mlog.Error("invalid line=", line)
			continue
		}

		// key, *value = line.split(":", 1)
		if strings.TrimSpace(tmplist[0]) == "" {
			continue
		}
		value_str := strings.TrimSpace(tmplist[1])
		e.Headers[strings.TrimSpace(tmplist[0])] = value_str
	}
}
func (e *Executor) init_body() {
	body := e.node_data.Body
	if body != nil {
		data := body.Data
		switch body.Type {
		case "none":
			e.Content = ""
		case "raw-text":
			if len(data) != 1 {
				panic(exceptions.NewRequestBodyError("raw-text body type should have exactly one item"))
			}
			e.Content = e.variable_pool.ConvertTemplate(data[0].Value).Text()
		case "json":
			if len(data) != 1 {
				panic(exceptions.NewRequestBodyError("json body type should have exactly one item"))
			}
			json_string := e.variable_pool.ConvertTemplate(data[0].Value).Text()
			var tmp []any
			err := json.Unmarshal([]byte(json_string), &tmp)
			if err != nil {
				var tmpmap map[string]any
				err = json.Unmarshal([]byte(json_string), &tmpmap)
				if err != nil {
					mlog.Errorf("failed to parse json str(%s) :%v", json_string, err)
					panic(exceptions.NewRequestBodyError("Failed to parse JSON: " + json_string))
				} else {
					e.Json = tmpmap
				}
			} else {
				e.Json = tmp
			}
			// e.Json = e._parse_object_contains_variables(json_object)
		// case "binary":
		//     if len(data) != 1{
		//         return exceptions.NewRequestBodyError("binary body type should have exactly one item")
		// 	}
		//     file_selector = data[0].File
		//     file_variable = e.variable_pool.GetFile(file_selector)
		//     if file_variable == nil {
		//         return exceptions.NewFileFetchError(f"cannot fetch file with selector {file_selector}")
		// 	}
		//     file = file_variable.value
		//     e.Content = file_manager.download(file)
		case "x-www-form-urlencoded":
			form_data := map[string]string{}
			for _, item := range data {
				if item.Key != "" && item.Value != "" {
					form_data[e.variable_pool.ConvertTemplate(item.Key).Text()] = e.variable_pool.ConvertTemplate(item.Value).Text()
				}
			}

			e.Data = form_data
		case "form-data":
			form_data := map[string]string{}
			for _, item := range data {
				if item.Type == "text" && item.Key != "" && item.Value != "" {
					form_data[e.variable_pool.ConvertTemplate(item.Key).Text()] = e.variable_pool.ConvertTemplate(item.Value).Text()
				}
			}

			e.Data = form_data
		}
	}
}

func generate_random_string(n int) string {
	/*
	   Generate a random string of lowercase ASCII letters.

	   Args:
	       n (int): The length of the random string to generate.

	   Returns:
	       str: A random string of lowercase ASCII letters with length n.

	   Example:
	       >>> _generate_random_string(5)
	       'abcde'
	*/

	tmp := make([]byte, n)
	for i := range n {
		tmp[i] = byte(97 + rand.Intn(122-97))
	}
	return string(tmp)
}

func (e *Executor) assembling_headers() map[string]string {
	headers := maps.Clone(e.Headers)
	if e.Auth.Type == "api-key" {
		if e.Auth.Config == nil {
			panic(exceptions.NewAuthorizationConfigError("e.authorization config is required"))
		}
		if e.Auth.Config.APIKey == "" {
			panic(exceptions.NewAuthorizationConfigError("api_key is required"))
		}
		if e.Auth.Config.Header == "" {
			e.Auth.Config.Header = "Authorization"
		}
		if e.Auth.Config.Type == "bearer" {
			headers[e.Auth.Config.Header] = "Bearer " + e.Auth.Config.APIKey
		} else if e.Auth.Config.Type == "basic" {
			headers[e.Auth.Config.Header] = "Basic " + e.Auth.Config.APIKey
		} else if e.Auth.Config.Type == "custom" {
			headers[e.Auth.Config.Header] = e.Auth.Config.APIKey
		}
	}
	return headers
}

// func (e *Executor) validate_and_parse_response(response *http.Response) *httprequestnodesentities.Response{
//         executor_response := &httprequestnodesentities.NewResponse(response)

//         threshold_size = (
//             dify_config.HTTP_REQUEST_NODE_MAX_BINARY_SIZE
//             if executor_response.is_file
//             else dify_config.HTTP_REQUEST_NODE_MAX_TEXT_SIZE
//         )
//         if executor_response.size > threshold_size:
//             raise ResponseSizeError(
//                 f"{'File' if executor_response.is_file else 'Text'} size is too large,"
//                 f" max size is {threshold_size / 1024 / 1024:.2f} MB,"
//                 f" but current size is {executor_response.readable_size}."
//             )

//	        return executor_response
//		}
func (e *Executor) to_log() string {
	url_parts, err := url.Parse(e.URL)
	if err != nil {
		panic(exceptions.NewInvalidURLError("url=" + e.URL))
	}
	path := url_parts.Path
	if path == "" {
		path = "/"
	}

	// Add query parameters
	if len(e.Params) > 0 {
		queryString := url.Values{}
		for _, v := range e.Params {
			queryString.Set(v[0], v[1])
		}
		path += "?" + queryString.Encode()
	} else if url_parts.Query() != nil {
		path += "?" + url_parts.Query().Encode()
	}
	raw := fmt.Sprintf("%s %s HTTP/1.1\r\n", strings.ToUpper(e.Method), path)
	raw += fmt.Sprintf("Host: %s\r\n", url_parts.Host)

	headers := e.assembling_headers()
	body := e.node_data.Body
	boundary := "----WebKitFormBoundary" + generate_random_string(16)
	if body != nil {
		if _, ok := BODY_TYPE_TO_CONTENT_TYPE[body.Type]; ok && !slices.Contains(slices.Sorted(maps.Keys(e.Headers)), "content-type") {
			headers["Content-Type"] = BODY_TYPE_TO_CONTENT_TYPE[body.Type]
		}
		if body.Type == "form-data" {
			headers["Content-Type"] = "multipart/form-data; boundary=" + boundary
		}
	}
	for k, v := range headers {
		if e.Auth != nil && e.Auth.Type == "api-key" {
			authorization_header := "Authorization"
			if e.Auth.Config != nil && e.Auth.Config.Header != "" {
				authorization_header = e.Auth.Config.Header
			}

			if strings.EqualFold(k, authorization_header) {
				raw += fmt.Sprintf("%s: %s\r\n", k, strings.Repeat("*", len(v)))
				continue
			}
		}
		raw += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body_string := ""
	if e.node_data.Body != nil {
		if e.Content != nil {
			if _, ok := e.Content.(string); ok {
				body_string = e.Content.(string)
			} else if _, ok := e.Content.([]byte); ok {
				body_string = string(e.Content.([]byte))
			}
		} else if len(e.Data) > 0 && e.node_data.Body.Type == "x-www-form-urlencoded" {
			queryString := url.Values{}
			for k, v := range e.Data {
				queryString.Set(k, v)
			}
			body_string = queryString.Encode()
		} else if len(e.Data) > 0 && e.node_data.Body.Type == "form-data" {
			for key, value := range e.Data {
				body_string += fmt.Sprintf("--%s\r\n", boundary)
				body_string += fmt.Sprintf(`Content-Disposition: form-data; name="%s"\r\n\r\n`, key)
				body_string += fmt.Sprintf("%s\r\n", value)
			}
			body_string += fmt.Sprintf("--%s--\r\n", boundary)
		} else if e.Json != nil {
			bindata, _ := json.Marshal(e.Json)
			body_string = string(bindata)
		} else if e.node_data.Body.Type == "raw-text" {
			if len(e.node_data.Body.Data) != 1 {
				panic(exceptions.NewRequestBodyError("raw-text body type should have exactly one item"))
			}
			body_string = e.node_data.Body.Data[0].Value
		}
	}
	if body_string != "" {
		raw += fmt.Sprintf("Content-Length: %d\r\n", len(body_string))
	}
	raw += "\r\n" // Empty line between headers and body
	raw += body_string

	return raw
}

func (e *Executor) do_http_request(headers map[string]string) *http.Response {
	if !slices.Contains(httprequestnodesentities.METHODS, e.Method) {
		panic(exceptions.NewInvalidHttpMethodError("Invalid http method " + e.Method))
	}
	// 构建请求
	var req *http.Request
	var err error
	var body io.Reader

	switch {
	case e.Json != nil:
		jsonData, err := json.Marshal(e.Json)
		if err != nil {
			panic(exceptions.NewHttpRequestNodeError(err.Error()))
		}
		body = bytes.NewBuffer(jsonData)
		req, err = http.NewRequest(strings.ToUpper(e.Method), e.URL, body)
		if err != nil {
			panic(exceptions.NewHttpRequestNodeError(err.Error()))
		}
		headers["Content-Type"] = "application/json"
	case e.Data != nil:
		formData := url.Values{}
		for k, v := range e.Data {
			formData.Set(k, v)
		}
		body = strings.NewReader(formData.Encode())
		req, err = http.NewRequest(strings.ToUpper(e.Method), e.URL, body)
		if err != nil {
			panic(exceptions.NewHttpRequestNodeError(err.Error()))
		}
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	case e.Content != nil:
		if _, ok := e.Content.(string); ok {
			body = bytes.NewBuffer([]byte(e.Content.(string)))
		} else if _, ok := e.Content.([]byte); ok {
			body = bytes.NewBuffer(e.Content.([]byte))
		}

		req, err = http.NewRequest(strings.ToUpper(e.Method), e.URL, body)
		if err != nil {
			panic(exceptions.NewHttpRequestNodeError(err.Error()))
		}
	default:
		req, err = http.NewRequest(strings.ToUpper(e.Method), e.URL, nil)
		if err != nil {
			panic(exceptions.NewHttpRequestNodeError(err.Error()))
		}
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 设置查询参数
	if e.Params != nil {
		query := req.URL.Query()
		for _, v := range e.Params {
			query.Set(v[0], v[1])
		}
		req.URL.RawQuery = query.Encode()
	}

	// 设置超时
	client := &http.Client{
		Timeout: time.Duration(e.Timeout.Connect+e.Timeout.Read+e.Timeout.Write) * time.Second,
	}
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		mlog.Errorf("dump request failed:%v", err)
	} else {
		fmt.Printf("REQUEST:\n%s", string(reqDump))
	}
	// 发送请求并处理响应
	resp, err := client.Do(req)
	if err != nil {
		panic(exceptions.NewHttpRequestNodeError(err.Error()))
	}
	rspDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		mlog.Errorf("dump response failed:%v", err)
	} else {
		fmt.Printf("RESPONSE:\n%s", string(rspDump))
	}
	return resp
}

func (e *Executor) validate_and_parse_response(response *http.Response) *httprequestnodesentities.Response {
	executor_response := httprequestnodesentities.NewResponse(response)

	threshold_size := viper.GetInt64WithDefault("http_node_config.max_text_size", 1048576)
	if executor_response.Size() > threshold_size {
		panic(exceptions.NewResponseSizeError(fmt.Sprintf("'Text' size is too large, max size is %.2f MB, but current size is %s.", float64(threshold_size)/float64(1024)/float64(1024), executor_response.ReadableSize())))
	}
	return executor_response
}
func (e *Executor) invoke() *httprequestnodesentities.Response {
	// assemble headers
	headers := e.assembling_headers()
	// do http request
	response := e.do_http_request(headers)
	// validate response
	return e.validate_and_parse_response(response)
}
