package base

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
)

const (
	CODE_PLACEHOLDER   = "{{code}}"
	INPUTS_PLACEHOLDER = "{{inputs}}"
	RESULT_TAG         = "<<RESULT>>"
)

type TemplateTransformer interface {
	GetRunnerScript() string
	SerializeInputs(inputs map[string]any) string
	AssembleRunnerScript(cls TemplateTransformer, code string, inputs map[string]any) string
	GetPreloadScript() string
	TransformCaller(cls TemplateTransformer, code string, inputs map[string]any) (string, string)
	ExtractResultStrFromResponse(response string) string
	TransformResponse(cls TemplateTransformer, response string) map[string]any
}

type BaseTemplateTransformer struct {
}

func (former *BaseTemplateTransformer) SerializeInputs(inputs map[string]any) string {
	inputs_json_str, _ := json.Marshal(inputs)
	inputBase64Encoded := base64.StdEncoding.EncodeToString(inputs_json_str)
	return inputBase64Encoded
}
func (former *BaseTemplateTransformer) AssembleRunnerScript(cls TemplateTransformer, code string, inputs map[string]any) string {
	// assemble runner script
	script := cls.GetRunnerScript()
	script = strings.ReplaceAll(script, CODE_PLACEHOLDER, code)
	inputs_str := cls.SerializeInputs(inputs)
	script = strings.ReplaceAll(script, INPUTS_PLACEHOLDER, inputs_str)
	return script

}
func (former *BaseTemplateTransformer) GetPreloadScript() string {
	/*
		Get preload script
	*/
	return ""

}
func (former *BaseTemplateTransformer) TransformCaller(cls TemplateTransformer, code string, inputs map[string]any) (string, string) {
	/*
	   Transform code to python runner
	   :param code: code
	   :param inputs: inputs
	   :return: runner, preload
	*/
	runner_script := cls.AssembleRunnerScript(cls, code, inputs)
	preload_script := cls.GetPreloadScript()
	return runner_script, preload_script
}
func (former *BaseTemplateTransformer) ExtractResultStrFromResponse(response string) string {
	// 构建正则表达式模式
	pattern := fmt.Sprintf("%s(.*)%s", regexp.QuoteMeta(RESULT_TAG), regexp.QuoteMeta(RESULT_TAG))

	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		mlog.Errorf("failed to compile regexp patten:%v", err)
		panic(exceptions.NewValueError("failed to compile regex: " + err.Error()))
	}

	// 执行正则表达式匹配
	match := re.FindStringSubmatch(response)
	mlog.Debugf("--------match=%#v", match)
	if len(match) < 2 {
		mlog.Errorf("failed to parse result")
		panic(exceptions.NewValueError("failed to parse result"))
	}

	// 返回匹配到的字符串
	return match[1]
}
func (former *BaseTemplateTransformer) TransformResponse(cls TemplateTransformer, response string) map[string]any {
	/*
	   Transform response to dict
	   :param response: response
	   :return:
	*/
	// 使用strconv.Unquote解析转义序列
	response = response[:len(response)-1]
	decodedStr, err := strconv.Unquote(strings.Replace(strconv.Quote(response), `\\u`, `\u`, -1))
	if err != nil {
		mlog.Errorf("strconv.Unquote failed:%v", err)
		panic(exceptions.NewValueError("failed to parse unicode response"))
	}
	mlog.Debugf("------response=%s", decodedStr)
	tmpstr := cls.ExtractResultStrFromResponse(decodedStr)

	mlog.Debugf("------extract result str=%s", tmpstr)
	var result map[string]any
	err = json.Unmarshal([]byte(tmpstr), &result)
	if err != nil {
		mlog.Errorf("failed to parse response:%v", err)
		panic(exceptions.NewValueError("failed to parse response"))
	}

	return result
}
