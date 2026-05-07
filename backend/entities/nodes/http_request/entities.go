package httprequest

import (
	"fmt"
	"io"
	"net/http"

	"github.com/odysseythink/confy"
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
)

var (
	METHODS = []string{
		"get",
		"post",
		"put",
		"patch",
		"delete",
		"head",
		"options",
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"HEAD",
		"OPTIONS",
	}
)

type HttpRequestNodeAuthorizationConfig struct {
	Type   string `json:"type"` // "basic", "bearer", "custom"
	APIKey string `json:"api_key"`
	Header string `json:"header"`
}

type HttpRequestNodeAuthorization struct {
	Type   string                              `json:"type"` //["no-auth", "api-key"]
	Config *HttpRequestNodeAuthorizationConfig `json:"config"`
}

// func (hra *HttpRequestNodeAuthorization) CheckConfig() error {
// 	if hra.Type == "no-auth" {
// 		if hra.Config != nil {
// 			return errors.New("config should be None for type no-auth")
// 		}
// 		return nil
// 	} else {
// 		if hra.Config == nil {
// 			return errors.New("config should be a dict")
// 		}
// 		return nil
// 	}
// }

type BodyData struct {
	Key   string   `json:"key"`
	Type  string   `json:"type"` // : Literal["file", "text"]
	Value string   `json:"value"`
	File  []string `json:"file"`
}

type HttpRequestNodeBody struct {
	Type string      `json:"type"` //: Literal["none", "form-data", "x-www-form-urlencoded", "raw-text", "json", "binary"]
	Data []*BodyData `json:"data"`
}

// func (hrnb *HttpRequestNodeBody) CheckData() []BodyData {
// 	if hrnb.Data == nil {
// 		return []BodyData{}
// 	}
// 	if len(hrnb.Data) == 1 && hrnb.Data[0].Value != "" {
// 		return []BodyData{
// 			{
// 				Key:   "",
// 				Type:  "text",
// 				Value: hrnb.Data[0].Value,
// 			},
// 		}
// 	}
// 	return hrnb.Data
// }

type HttpRequestNodeTimeout struct {
	Connect int `json:"connect"`
	Read    int `json:"read"`
	Write   int `json:"write"`
}

func NewHttpRequestNodeTimeout(connect, read, write int) *HttpRequestNodeTimeout {
	if connect <= 0 {
		connect = confy.Get[int]("http_node_config.max_connect_timeout")
	}
	if read <= 0 {
		read = confy.Get[int]("http_node_config.max_read_timeout")
	}
	if write <= 0 {
		write = confy.Get[int]("http_node_config.max_write_timeout")
	}
	return &HttpRequestNodeTimeout{
		Connect: connect,
		Read:    read,
		Write:   write,
	}
}

// CodeNodeData represents answer node data
type HttpRequestNodeData struct {
	*basenodesentities.BaseNodeData
	Method        string                        `json:"method"`
	URL           string                        `json:"url"`
	Authorization *HttpRequestNodeAuthorization `json:"authorization"`
	Headers       string                        `json:"headers"`
	Params        string                        `json:"params"`
	Body          *HttpRequestNodeBody          `json:"body"`
	Timeout       *HttpRequestNodeTimeout       `json:"timeout"`
}

func New() *HttpRequestNodeData {
	return &HttpRequestNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
	}
}

type Response struct {
	*http.Response
}

func NewResponse(response *http.Response) *Response {
	return &Response{
		Response: response,
	}
}

func (r *Response) IsFile() bool {
	// contentType := r.headers.Get("Content-Type")
	// contentDisposition := r.response.Header.Get("Content-Disposition")

	// // 检查 Content-Disposition 头部
	// if contentDisposition != "" {
	// 	disp_type, params, err := mime.ParseMediaType(contentDisposition)
	// 	if err != nil {
	// 		mlog.Errorf("ParseMediaType failed:%v", err)
	// 		// return false
	// 	} else {
	// 		if _, ok := params["filename"]; ok && params["filename"] != "" && disp_type == "attachment" {
	// 			return true
	// 		}
	// 	}
	// }

	// // 检查 Content-Type
	// if strings.HasPrefix(contentType, "application/") {
	// 	textTypes := []string{"json", "xml", "javascript", "x-www-form-urlencoded", "yaml", "graphql"}
	// 	for _, textType := range textTypes {
	// 		if strings.Contains(contentType, textType) {
	// 			return false
	// 		}
	// 	}
	// }

	// // 使用 MIME 类型分析
	// ext := mime.TypeByExtension(strings.TrimPrefix(mime.ExtensionByType(contentType), "."))
	// if ext != "" {
	// 	mainType := mime.TypeByExtension(ext)
	// 	if mainType != "" {
	// 		return strings.HasPrefix(mainType, "application/") || strings.HasPrefix(mainType, "image/") ||
	// 			strings.HasPrefix(mainType, "audio/") || strings.HasPrefix(mainType, "video/")
	// 	}
	// }

	// // 检查是否是媒体类型
	// mediaTypes := []string{"image/", "audio/", "video/"}
	// for _, mediaType := range mediaTypes {
	// 	if strings.HasPrefix(contentType, mediaType) {
	// 		return true
	// 	}
	// }

	return false
}

func (r *Response) Text() string {
	bodyBytes, _ := io.ReadAll(r.Body)
	return string(bodyBytes)
}

func (r *Response) Content() []byte {
	bodyBytes, _ := io.ReadAll(r.Body)
	return bodyBytes
}

func (r *Response) Size() int64 {
	return r.ContentLength
}

func (r *Response) ReadableSize() string {
	size := r.Size()
	if size < 1024 {
		return fmt.Sprintf("%d bytes", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
}
