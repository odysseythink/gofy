package codeexecutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"mlib.com/confy"
	"mlib.com/gofy/server/cluster"
	codenodesexceptions "mlib.com/gofy/server/core/exceptions/nodes/code"
	"mlib.com/gofy/server/core/helper/code_executor/template_transformer/base"
	"mlib.com/gofy/server/core/helper/code_executor/template_transformer/jinja2"
	"mlib.com/gofy/server/core/helper/code_executor/template_transformer/python3"
	codeexecutorenumtypes "mlib.com/gofy/server/enum_types/code_executor"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

var (
	code_template_transformers = map[codeexecutorenumtypes.CodeLanguage]base.TemplateTransformer{
		codeexecutorenumtypes.CodeLanguage_PYTHON3: &python3.Python3TemplateTransformer{},
		codeexecutorenumtypes.CodeLanguage_JINJA2:  &jinja2.Jinja2TemplateTransformer{},
	}
	code_language_to_running_language = map[codeexecutorenumtypes.CodeLanguage]codeexecutorenumtypes.CodeLanguage{
		codeexecutorenumtypes.CodeLanguage_PYTHON3: codeexecutorenumtypes.CodeLanguage_PYTHON3,
		codeexecutorenumtypes.CodeLanguage_JINJA2:  codeexecutorenumtypes.CodeLanguage_PYTHON3,
	}
)

type CodeExecutionResponse struct {
	Data struct {
		Stdout string `json:"stdout"`
		Error  string `json:"error"`
	} `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ExecuteCode1(language codeexecutorenumtypes.CodeLanguage, preload string, code string) string {
	if _, ok := code_language_to_running_language[language]; !ok {
		mlog.Errorf("language(%s) is not surpported", language)
		panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("language(%s) is not surpported", language)))
	}
	conn := cluster.Instance().GetRpcClientByModule("sandbox")
	if conn != nil {
		pbrsp, err := pbapi.NewSandboxClient(conn).Run(context.Background(), &pbapi.RunRequest{
			Language:      string(code_language_to_running_language[language]),
			Code:          code,
			Preload:       preload,
			EnableNetwork: true,
		})
		if err != nil {
			mlog.Errorf("remote call Run failed:%v", err)
			panic(codenodesexceptions.NewCodeExecutionError("remote call Run failed:" + err.Error()))
		} else {
			mlog.Infof("remote call Run return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				panic(pbrsp.Exp)
			} else {
				if pbrsp.Code != 0 {
					panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("Got error code: %d. Got error msg: %s", pbrsp.Code, pbrsp.Message)))
				}

				if pbrsp.Error != "" {
					panic(codenodesexceptions.NewCodeExecutionError(pbrsp.Error))
				}

				return pbrsp.Stdout
			}
		}
	} else {
		mlog.Errorf("get rpc client failed")
		panic(codenodesexceptions.NewCodeExecutionError("get rpc client failed"))
	}
}
func ExecuteCode(language codeexecutorenumtypes.CodeLanguage, preload string, code string) string {
	/*
		Execute code
		:param language: code language
		:param code: code
		:return:
	*/
	mlog.Debugf("----run code=%s", code)
	u, err := url.Parse(confy.GetWithDefault[string]("code_execution_config.endpoint", "http://127.0.0.1:8194"))
	if err != nil {
		mlog.Errorf("url parse failed:%v", err)
		panic(codenodesexceptions.NewCodeExecutionError("Invalid URL"))
	}
	u.Path = "/v1/sandbox/run"
	// 构建请求头
	headers := http.Header{
		"Content-Type": {"application/json"},
		"X-Api-Key":    {confy.GetWithDefault[string]("code_execution_config.api_key", "dify-sandbox")},
	}
	if _, ok := code_language_to_running_language[language]; !ok {
		mlog.Errorf("language(%s) is not surpported", language)
		panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("language(%s) is not surpported", language)))
	}
	// 构建请求体
	data := map[string]any{
		"language":       code_language_to_running_language[language],
		"code":           code,
		"preload":        preload,
		"enable_network": true,
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout: time.Duration(confy.GetWithDefault[float64]("code_execution_config.connect_timeout", 10.0)+confy.GetWithDefault[float64]("code_execution_config.read_timeout", 60.0)+confy.GetWithDefault[float64]("code_execution_config.write_timeout", 10.0)) * time.Second,
	}

	// 创建请求
	reqBody, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		panic(codenodesexceptions.NewCodeExecutionError("Failed to create request"))
	}
	req.Header = headers
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		mlog.Errorf("dump request failed:%v", err)
	} else {
		fmt.Printf("REQUEST:\n%s", string(reqDump))
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("Failed to execute code, which is likely a network issue, please check if the sandbox service is running. (Error: %s)", err.Error())))
	}
	defer resp.Body.Close()
	rspDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		mlog.Errorf("dump response failed:%v", err)
	} else {
		fmt.Printf("RESPONSE:\n%s", string(rspDump))
	}

	// 检查状态码
	if resp.StatusCode == http.StatusServiceUnavailable {
		panic(codenodesexceptions.NewCodeExecutionError("Code execution service is unavailable"))
	} else if resp.StatusCode != http.StatusOK {
		panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("Failed to execute code, got status code %d, please check if the sandbox service is running", resp.StatusCode)))
	}

	// 解析响应
	var response CodeExecutionResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(codenodesexceptions.NewCodeExecutionError("Failed to read response"))
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		panic(codenodesexceptions.NewCodeExecutionError("Failed to parse response"))
	}
	mlog.Debugf("------response=%#v", response)
	// 检查响应数据
	if response.Code != 0 {
		panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("Got error code: %d. Got error msg: %s", response.Code, response.Message)))
	}

	if response.Data.Error != "" {
		panic(codenodesexceptions.NewCodeExecutionError(response.Data.Error))
	}

	return response.Data.Stdout
}
func ExecuteWorkflowCodeTemplate(language codeexecutorenumtypes.CodeLanguage, code string, inputs map[string]any) map[string]any {
	/*
	   Execute code
	   :param language: code language
	   :param code: code
	   :param inputs: inputs
	   :return:
	*/
	if _, ok := code_template_transformers[language]; !ok || code_template_transformers[language] == nil {
		panic(codenodesexceptions.NewCodeExecutionError(fmt.Sprintf("Unsupported language %s", language)))
	}
	template_transformer := code_template_transformers[language]

	runner, preload := template_transformer.TransformCaller(template_transformer, code, inputs)

	response := ExecuteCode1(language, preload, runner)

	return template_transformer.TransformResponse(template_transformer, response)
}
