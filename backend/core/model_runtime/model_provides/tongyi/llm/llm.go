package llm

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/odysseythink/mlog"
	"github.com/pkoukk/tiktoken-go"
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/exceptions"
	dashscopeexception "mlib.com/gofy/server/core/exceptions/dashscope"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	commontypes "mlib.com/gofy/server/types/common"
	dashscopetypes "mlib.com/gofy/server/types/dashscope"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Input struct {
	Messages []Message `json:"messages"`
}

type Parameters struct {
	ResultFormat string `json:"result_format"`
}

type RequestBody struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type TongyiLargeLanguageModel struct {
	*base.LargeLanguageModel
	tokenizers map[string]*tiktoken.Tiktoken
}

func (m *TongyiLargeLanguageModel) InvokeStream(
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		responses := m.generate_call(model, credentials, prompt_messages, model_parameters, tools, stop, true)

		full_text := ""
		tool_calls := []*dashscopetypes.ResponseToolCall{}
		index := 0
		for response := range responses {
			mlog.Debugf("------response=%#v", response)
			if response.StatusCode != 200 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)))
			}
			if response.Output == nil || len(response.Output.Choices) == 0 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)))
			}
			resp_finish_reason := response.Output.Choices[0].FinishReason
			if resp_finish_reason != "" && resp_finish_reason != "null" {
				resp_content := ""
				assistant_prompt_message := modelruntimeentities.NewAssistantPromptMessage("", "", nil)
				tool_calls = response.Output.Choices[0].Message.ToolCalls
				if response.Output.Choices[0].Message.Content != nil {
					// special for qwen-vl
					if val, ok := response.Output.Choices[0].Message.Content.([]any); ok {
						if _, ok := val[0].(map[string]any); ok {
							if _, ok := val[0].(map[string]any)["text"]; ok {
								if _, ok := val[0].(map[string]any)["text"].(string); ok {
									resp_content = val[0].(map[string]any)["text"].(string)
								}
							}
						}
					} else if val, ok := response.Output.Choices[0].Message.Content.(string); ok {
						resp_content = val
					}

					// transform assistant message to prompt message
					assistant_prompt_message.Content = strings.Replace(resp_content, full_text, "", 1)
					full_text = resp_content
				}
				if len(tool_calls) > 0 {
					message_tool_calls := []*modelruntimeentities.ToolCall{}
					for _, tool_call_obj := range tool_calls {
						message_tool_call := &modelruntimeentities.ToolCall{
							ID:   tool_call_obj.Function.Name,
							Type: "function",
							Function: modelruntimeentities.ToolCallFunction{
								Name: tool_call_obj.Function.Name, Arguments: tool_call_obj.Function.Arguments,
							},
						}
						message_tool_calls = append(message_tool_calls, message_tool_call)
					}
					assistant_prompt_message.ToolCalls = message_tool_calls
				}
				// transform usage
				usage := response.Usage
				llmusage := m.CalcResponseUsage(m, model, credentials, usage.InputTokens, usage.OutputTokens)
				if !yield(&modelruntimeentities.LLMResultChunk{
					Model:          model,
					PromptMessages: prompt_messages,
					Delta: &modelruntimeentities.LLMResultChunkDelta{
						Index: index, Message: assistant_prompt_message, FinishReason: resp_finish_reason, Usage: llmusage,
					},
				}) {
					return
				}
			} else {
				resp_content := "" // response.Output.Choices[0].Message.content
				if response.Output.Choices[0].Message.Content == nil {
					tool_calls = response.Output.Choices[0].Message.ToolCalls
					continue
				}
				if val, ok := response.Output.Choices[0].Message.Content.([]any); ok {
					if _, ok := val[0].(map[string]any); ok {
						if _, ok := val[0].(map[string]any)["text"]; ok {
							if _, ok := val[0].(map[string]any)["text"].(string); ok {
								resp_content = val[0].(map[string]any)["text"].(string)
							}
						}
					}
				} else if val, ok := response.Output.Choices[0].Message.Content.(string); ok {
					resp_content = val
				}
				// special for qwen-vl
				// if isinstance(resp_content, list){
				// 	resp_content = resp_content[0]["text"]
				// }
				// transform assistant message to prompt message
				assistant_prompt_message := modelruntimeentities.NewAssistantPromptMessage(
					strings.Replace(resp_content, full_text, "", 1), "", nil,
				)
				full_text = resp_content
				if !yield(&modelruntimeentities.LLMResultChunk{
					Model:          model,
					PromptMessages: prompt_messages,
					Delta: &modelruntimeentities.LLMResultChunkDelta{
						Index: index, Message: assistant_prompt_message,
					},
				}) {
					return
				}
			}
			index++
		}
	}
}
func (m *TongyiLargeLanguageModel) Invoke(
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
) *modelruntimeentities.LLMResult {
	/*
	   Invoke large language model
	   :param model: model name
	   :param credentials: model credentials
	   :param prompt_messages: prompt messages
	   :param model_parameters: model parameters
	   :param tools: tools for tool calling
	   :param stop: stop words
	   :param stream: is stream response
	   :param user: unique user id
	   :return: full response or stream response chunk generator result
	*/
	// invoke model without code wrapper
	response := m.generate_call(model, credentials, prompt_messages, model_parameters, tools, stop, false)
	return m.handleGenerateResponse(model, credentials, response, prompt_messages)
}
func (m *TongyiLargeLanguageModel) GetNumTokens(
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	tools []*modelruntimeentities.PromptMessageTool,
) int {
	/*
	   Get number of tokens for given prompt messages
	   :param model: model name
	   :param credentials: model credentials
	   :param prompt_messages: prompt messages
	   :param tools: tools for tool calling
	   :return
	*/
	// Check if the model was added via GetCustomizableModelSchema
	if m.GetCustomizableModelSchema(model, credentials) != nil {
		// For custom models, tokens are not calculated.
		return 0
	}
	if slices.Contains([]string{"qwen-turbo-chat", "qwen-plus-chat"}, model) {
		model = strings.ReplaceAll(model, "-chat", "")
	}
	if model == "farui-plus" {
		model = "qwen-farui-plus"
	}
	// if you don't want download dictionary at runtime, you can use offline loader

	var tokenizer *tiktoken.Tiktoken
	if _, ok := m.tokenizers[model]; ok {
		tokenizer = m.tokenizers[model]
	} else {
		// tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
		tke, err := tiktoken.GetEncoding("cl100k_base")
		if err != nil {
			mlog.Errorf("getEncoding: %v", err)
			return 0
		}
		tokenizer = tke
		m.tokenizers[model] = tokenizer
	}
	// convert string to token ids
	str := m._convert_messages_to_prompt(prompt_messages)
	tokens := tokenizer.Encode(str, nil, nil)
	return len(tokens)
}

func (m *TongyiLargeLanguageModel) ValidateCredentials(model string, credentials map[string]any) {
	/*
	   Validate model credentials
	   :param model: model name
	   :param credentials: model credentials
	   :return
	*/

	m.generate(
		model,
		credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{
			"temperature": 0.5,
		},
		nil,
		nil,
		false,
		"",
	)
}

func (m *TongyiLargeLanguageModel) ProviderName() string {
	return "tongyi"
}

func (m *TongyiLargeLanguageModel) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (m *TongyiLargeLanguageModel) generate_call(
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*= True*/
) iter.Seq[*dashscopetypes.Response] {
	return func(yield func(*dashscopetypes.Response) bool) {
		if slices.Contains([]string{"qwen-turbo-chat", "qwen-plus-chat"}, model) {
			model = strings.ReplaceAll(model, "-chat", "")
		}
		// 创建 HTTP 客户端
		client := &http.Client{}

		// 构建请求体
		requestBody := map[string]any{
			// 此处以qwen-plus为例，可按需更换模型名称。模型列表：https://help.aliyun.com/zh/model-studio/getting-started/models
			"model": model,
			"input": map[string]any{
				// "messages": []Message{
				// 	{
				// 		Role:    "system",
				// 		Content: "You are a helpful assistant.",
				// 	},
				// 	{
				// 		Role:    "user",
				// 		Content: "你是谁？",
				// 	},
				// },
			},
			"parameters": map[string]any{},
		}
		url := ""
		extra_model_kwargs := map[string]any{}
		if len(tools) > 0 {
			extra_model_kwargs["tools"] = m.convertTools(tools)
		}
		if len(stop) > 0 {
			requestBody["stop"] = stop
		}

		parameters := map[string]any{}
		maps.Copy(parameters, model_parameters)
		// maps.Copy(params, credentials_kwargs)
		maps.Copy(parameters, extra_model_kwargs)
		if _, ok := parameters["result_format"]; !ok {
			parameters["result_format"] = "message"
		}
		requestBody["parameters"] = parameters
		model_schema := m.GetModelSchema(m, model, credentials)
		if slices.Contains(model_schema.Features, modelruntimeenumtypes.ModelFeature_VISION) {
			messages, _ := m._convert_prompt_messages_to_tongyi_messages(prompt_messages, true)
			requestBody["input"] = map[string]any{"messages": messages}
			url = "https://dashscope.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation"
			// response = MultiModalConversation.call(**params, stream=stream)
		} else {
			// nothing different between chat model and completion model in tongyi
			messages, _ := m._convert_prompt_messages_to_tongyi_messages(prompt_messages, false)
			requestBody["input"] = map[string]any{"messages": messages}
			// response = Generation.call(**params, result_format="message", stream=stream)
			url = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			mlog.Error("json marshal failed:", err)
			panic(dashscopeexception.NewInvalidParameter(fmt.Sprintf("json unmarshal request body(%#v) failed:%v", requestBody, err)))
		}

		// 创建 POST 请求
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			mlog.Error("http.NewRequest failed:", err)
			panic(dashscopeexception.NewRequestFailure("", err.Error(), "", ""))
		}

		// 设置请求头
		// 若没有配置环境变量，请用百炼API Key将下行替换为：apiKey := "sk-xxx"
		apiKey := ""
		if credentials == nil {
			apiKey = ""
		} else {
			if _, ok := credentials["dashscope_api_key"]; ok {
				if _, ok := credentials["dashscope_api_key"].(string); ok {
					apiKey = credentials["dashscope_api_key"].(string)
				}
			}
			// apiKey = credentials["dashscope_api_key"]
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		if stream {
			req.Header.Set("X-DashScope-SSE", "enable")
		}
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			mlog.Errorf("dump request failed:%v", err)
		} else {
			fmt.Printf("REQUEST:\n%s\n", string(reqDump))
		}
		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			mlog.Error("http.Request do failed:", err)
			panic(dashscopeexception.NewRequestFailure("", err.Error(), "", ""))
		}
		defer resp.Body.Close()
		rspDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			mlog.Errorf("dump response failed:%v", err)
		} else {
			fmt.Printf("RESPONSE:\n%s\n", string(rspDump))
		}
		if stream {
			sc := bufio.NewScanner(resp.Body)
			for {
				m := new(dashscopetypes.Response)
				http_status := 0
				if !sc.Scan() {
					mlog.Error("scan failed:", sc.Err())
					if sc.Err() != nil {
						panic(dashscopeexception.NewUnsupportedApiProtocol(sc.Err().Error()))
					}
					return
				}
				if strings.HasPrefix(sc.Text(), "id:") {
					if !sc.Scan() {
						mlog.Error("scan failed:", sc.Err())
						panic(dashscopeexception.NewUnsupportedApiProtocol(sc.Err().Error()))
					}
					if !strings.HasPrefix(sc.Text(), "event:") {
						mlog.Error("scan failed:the id line followed with event line")
						panic(dashscopeexception.NewUnsupportedApiProtocol("the id line followed with event line"))
					}
					if !sc.Scan() {
						mlog.Error("scan failed:", sc.Err())
						panic(dashscopeexception.NewUnsupportedApiProtocol(sc.Err().Error()))
					}
					if !strings.HasPrefix(sc.Text(), ":HTTP_STATUS/") {
						mlog.Error("scan failed:the event line followed with HTTP_STATUS line")
						panic(dashscopeexception.NewUnsupportedApiProtocol("the event line followed with HTTP_STATUS line"))
					}
					http_status, err = strconv.Atoi(strings.TrimPrefix(sc.Text(), ":HTTP_STATUS/"))
					if err != nil {
						mlog.Errorf("parse http status(%s) failed:%v", sc.Text(), err)
						panic(dashscopeexception.NewUnsupportedApiProtocol(fmt.Sprintf("parse http status(%s) failed:%v", sc.Text(), err)))
					}
					if !sc.Scan() {
						mlog.Error("scan failed:", sc.Err())
						panic(dashscopeexception.NewUnsupportedApiProtocol(sc.Err().Error()))
					}
					if !strings.HasPrefix(sc.Text(), "data:") {
						mlog.Error("scan failed:the HTTP_STATUS line followed with data line")
						panic(dashscopeexception.NewUnsupportedApiProtocol(sc.Err().Error()))
					}
					err = json.Unmarshal([]byte(strings.TrimPrefix(sc.Text(), "data:")), m)
					if err != nil {
						mlog.Errorf("json.Unmarshal(%s) failed:%v", sc.Text(), err)
						panic(dashscopeexception.NewUnsupportedApiProtocol(fmt.Sprintf("json.Unmarshal(%s) failed:%v", sc.Text(), err)))
					}
					m.StatusCode = http_status
					if !yield(m) {
						return
					}
				}
			}
		} else {
			dec := json.NewDecoder(resp.Body)
			for {
				m := new(dashscopetypes.Response)
				if err := dec.Decode(&m); err != nil {
					mlog.Errorf("json.Unmarshal failed:%v", err)
					panic(dashscopeexception.NewUnsupportedApiProtocol(fmt.Sprintf("json.Unmarshal failed:%v", err)))
				}
				if !yield(m) {
					return
				}
			}
		}
	}
}

func (m *TongyiLargeLanguageModel) generate(
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*= True*/
	user string,
) any /*-> Union[LLMResult, Generator]*/ {
	/*
	   Invoke large language model
	   :param model: model name
	   :param credentials: credentials
	   :param prompt_messages: prompt messages
	   :param tools: tools for tool calling
	   :param model_parameters: model parameters
	   :param stop: stop words
	   :param stream: is stream response
	   :param user: unique user id
	   :return: full response or stream response chunk generator result
	*/
	// transform credentials to kwargs for model instance
	response := m.generate_call(model, credentials, prompt_messages, model_parameters, tools, stop, stream)
	if stream {
		return m.handleGenerateStreamResponse(model, credentials, response, prompt_messages)
	}
	return m.handleGenerateResponse(model, credentials, response, prompt_messages)
}

func (m *TongyiLargeLanguageModel) handleGenerateResponse(
	model string, credentials map[string]any, responses iter.Seq[*dashscopetypes.Response], prompt_messages []modelruntimeentities.PromptMessager,
) *modelruntimeentities.LLMResult {
	/*
	   Handle llm response
	   :param model: model name
	   :param credentials: credentials
	   :param response: response
	   :param prompt_messages: prompt messages
	   :return: llm response
	*/
	for response := range responses {
		mlog.Debugf("------response=%#v", response)
		if response.StatusCode != 200 && response.StatusCode != 0 {
			mlog.Errorf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)
			panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)))
		}
		if response.Output == nil {
			mlog.Errorf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)
			panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)))
		}
		resp_content := ""
		if len(response.Output.Choices) > 0 && response.Output.Choices[0].Message.Content != nil {
			// special for qwen-vl
			if val, ok := response.Output.Choices[0].Message.Content.([]any); ok {
				if _, ok := val[0].(map[string]any); ok {
					if _, ok := val[0].(map[string]any)["text"]; ok {
						if _, ok := val[0].(map[string]any)["text"].(string); ok {
							resp_content = val[0].(map[string]any)["text"].(string)
						}
					}
				}
			} else if val, ok := response.Output.Choices[0].Message.Content.(string); ok {
				resp_content = val
			}
		} else if response.Output.Text != "" {
			resp_content = response.Output.Text
		}

		assistant_prompt_message := modelruntimeentities.NewAssistantPromptMessage(resp_content, "", nil)
		// transform usage
		usage := m.CalcResponseUsage(m, model, credentials, response.Usage.InputTokens, response.Usage.OutputTokens)
		// transform response
		result := &modelruntimeentities.LLMResult{
			Model:          model,
			Message:        assistant_prompt_message,
			PromptMessages: prompt_messages,
			Usage:          usage,
		}
		return result
	}
	panic(dashscopeexception.NewServiceUnavailableError("no response"))
}

func (m *TongyiLargeLanguageModel) handleGenerateStreamResponse(
	model string,
	credentials map[string]any,
	responses iter.Seq[*dashscopetypes.Response],
	prompt_messages []modelruntimeentities.PromptMessager,
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	/*
	   Handle llm stream response
	   :param model: model name
	   :param credentials: credentials
	   :param responses: response
	   :param prompt_messages: prompt messages
	   :return: llm response chunk generator result
	*/
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		full_text := ""
		tool_calls := []*dashscopetypes.ResponseToolCall{}
		index := 0
		for response := range responses {
			mlog.Debugf("------response=%#v", response)
			if response.StatusCode != 200 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)))
			}
			if response.Output == nil || len(response.Output.Choices) == 0 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.StatusCode, response.Message)))
			}
			resp_finish_reason := response.Output.Choices[0].FinishReason
			if resp_finish_reason != "" && resp_finish_reason != "null" {
				resp_content := ""
				assistant_prompt_message := modelruntimeentities.NewAssistantPromptMessage("", "", nil)
				tool_calls = response.Output.Choices[0].Message.ToolCalls
				if response.Output.Choices[0].Message.Content != nil {
					// special for qwen-vl
					if val, ok := response.Output.Choices[0].Message.Content.([]any); ok {
						if _, ok := val[0].(map[string]any); ok {
							if _, ok := val[0].(map[string]any)["text"]; ok {
								if _, ok := val[0].(map[string]any)["text"].(string); ok {
									resp_content = val[0].(map[string]any)["text"].(string)
								}
							}
						}
					} else if val, ok := response.Output.Choices[0].Message.Content.(string); ok {
						resp_content = val
					}

					// transform assistant message to prompt message
					assistant_prompt_message.Content = strings.Replace(resp_content, full_text, "", 1)
					full_text = resp_content
				}
				if len(tool_calls) > 0 {
					message_tool_calls := []*modelruntimeentities.ToolCall{}
					for _, tool_call_obj := range tool_calls {
						message_tool_call := &modelruntimeentities.ToolCall{
							ID:   tool_call_obj.Function.Name,
							Type: "function",
							Function: modelruntimeentities.ToolCallFunction{
								Name: tool_call_obj.Function.Name, Arguments: tool_call_obj.Function.Arguments,
							},
						}
						message_tool_calls = append(message_tool_calls, message_tool_call)
					}
					assistant_prompt_message.ToolCalls = message_tool_calls
				}
				// transform usage
				usage := response.Usage
				llmusage := m.CalcResponseUsage(m, model, credentials, usage.InputTokens, usage.OutputTokens)
				if !yield(&modelruntimeentities.LLMResultChunk{
					Model:          model,
					PromptMessages: prompt_messages,
					Delta: &modelruntimeentities.LLMResultChunkDelta{
						Index: index, Message: assistant_prompt_message, FinishReason: resp_finish_reason, Usage: llmusage,
					},
				}) {
					return
				}
			} else {
				resp_content := "" // response.Output.Choices[0].Message.content
				if response.Output.Choices[0].Message.Content == nil {
					tool_calls = response.Output.Choices[0].Message.ToolCalls
					continue
				}
				if val, ok := response.Output.Choices[0].Message.Content.([]any); ok {
					if _, ok := val[0].(map[string]any); ok {
						if _, ok := val[0].(map[string]any)["text"]; ok {
							if _, ok := val[0].(map[string]any)["text"].(string); ok {
								resp_content = val[0].(map[string]any)["text"].(string)
							}
						}
					}
				} else if val, ok := response.Output.Choices[0].Message.Content.(string); ok {
					resp_content = val
				}
				// special for qwen-vl
				// if isinstance(resp_content, list){
				// 	resp_content = resp_content[0]["text"]
				// }
				// transform assistant message to prompt message
				assistant_prompt_message := modelruntimeentities.NewAssistantPromptMessage(
					strings.Replace(resp_content, full_text, "", 1), "", nil,
				)
				full_text = resp_content
				if !yield(&modelruntimeentities.LLMResultChunk{
					Model:          model,
					PromptMessages: prompt_messages,
					Delta: &modelruntimeentities.LLMResultChunkDelta{
						Index: index, Message: assistant_prompt_message,
					},
				}) {
					return
				}
			}
			index++
		}
	}
}

func (m *TongyiLargeLanguageModel) toCredentialKwargs(credentials map[string]any) map[string]any {
	/*
	   Transform credentials to kwargs for model instance
	   :param credentials{
	   :return
	*/
	credentials_kwargs := map[string]any{}
	if credentials == nil {
		credentials_kwargs["api_key"] = ""
	} else {
		credentials_kwargs["api_key"] = credentials["dashscope_api_key"]
	}

	return credentials_kwargs
}

func (m *TongyiLargeLanguageModel) _convert_one_message_to_text(message modelruntimeentities.PromptMessager) (string, error) {
	/*
	   Convert a single message to a string.
	   :param message: modelruntimeentities.PromptMessage to convert.
	   :return: String representation of the message.
	*/
	human_prompt := "\n\nHuman:"
	ai_prompt := "\n\nAssistant:"
	message_text := ""

	if message.Role() == modelruntimeentities.PromptMessageRole_USER {
		if real_content, ok := message.GetContent().(string); ok {
			message_text = fmt.Sprintf("%s %s", human_prompt, real_content)
		} else if real_content, ok := message.GetContent().([]modelruntimeentities.PromptMessageContenter); ok {
			message_text = ""
			for _, sub_content := range real_content {
				if sub_message, ok := any(sub_content).(*modelruntimeentities.TextPromptMessageContent); ok {
					message_text = fmt.Sprintf("%s %s", human_prompt, sub_message.Data)
					break
				}
			}
		} else {
			return "", exceptions.NewValueError(fmt.Sprintf("Got unknown type %#v", message))
		}
	} else if message.Role() == modelruntimeentities.PromptMessageRole_ASSISTANT {
		if real_content, ok := message.GetContent().(string); ok {
			message_text = fmt.Sprintf("%s %s", ai_prompt, real_content)
		} else {
			return "", exceptions.NewValueError(fmt.Sprintf("Got unknown type %#v", message))
		}
	} else if message.Role() == modelruntimeentities.PromptMessageRole_SYSTEM {
		if real_content, ok := message.GetContent().(string); ok {
			message_text = real_content
		} else {
			return "", exceptions.NewValueError(fmt.Sprintf("Got unknown type %#v", message))
		}
	} else if message.Role() == modelruntimeentities.PromptMessageRole_TOOL {
		if real_content, ok := message.GetContent().(string); ok {
			message_text = real_content
		} else {
			return "", exceptions.NewValueError(fmt.Sprintf("Got unknown type %#v", message))
		}
	} else {
		return "", exceptions.NewValueError(fmt.Sprintf("Got unknown type %#v", message))
	}

	return message_text, nil
}

func (m *TongyiLargeLanguageModel) _convert_messages_to_prompt(messages []modelruntimeentities.PromptMessager) string {
	/*
	   Format a list of messages into a full prompt for the Anthropic model
	   :param messages: List of PromptMessage to combine.
	   :return: Combined string with necessary human_prompt and ai_prompt tags.
	*/
	strlist := []string{}
	for _, message := range messages {
		strmsg, err := m._convert_one_message_to_text(message)
		if err != nil {
			mlog.Errorf("convert One  Message To Text failed:%v", err)
			return ""
		}
		strlist = append(strlist, strmsg)
	}
	text := strings.Join(strlist, "")
	// trim off the trailing ' ' that might come from the "Assistant: "
	return strings.TrimRight(text, " ")
}

func (m *TongyiLargeLanguageModel) _convert_prompt_messages_to_tongyi_messages(
	prompt_messages []modelruntimeentities.PromptMessager, rich_content bool,
) ([]map[string]any, error) {
	/*
	   Convert prompt messages to tongyi messages
	   :param prompt_messages: prompt messages
	   :return: tongyi messages
	*/
	tongyi_messages := []map[string]any{}

	for _, prompt_message := range prompt_messages {
		if realpromptmsg, ok := any(prompt_message).(*modelruntimeentities.SystemPromptMessage[string]); ok {
			var content any = realpromptmsg.Content
			if rich_content {
				content = []map[string]string{{"text": realpromptmsg.Content}}
			}
			tongyi_messages = append(tongyi_messages, map[string]any{
				"role":    "system",
				"content": content,
			})
		} else if realpromptmsg, ok := any(prompt_message).(*modelruntimeentities.UserPromptMessage[string]); ok {
			var content any = realpromptmsg.Content
			if rich_content {
				content = []map[string]string{{"text": realpromptmsg.Content}}
			}
			tongyi_messages = append(tongyi_messages, map[string]any{
				"role":    "user",
				"content": content,
			})
		} else if realpromptmsg, ok := any(prompt_message).(*modelruntimeentities.UserPromptMessage[[]modelruntimeentities.PromptMessageContenter]); ok {
			sub_messages := []map[string]string{}
			for _, message_content := range realpromptmsg.Content {
				if message_content.Type() == modelruntimeentities.PromptMessageContent_TEXT {
					real_message_content := any(message_content).(*modelruntimeentities.TextPromptMessageContent)
					sub_message_dict := map[string]string{"text": real_message_content.Data()}
					sub_messages = append(sub_messages, sub_message_dict)
				} else if message_content.Type() == modelruntimeentities.PromptMessageContent_IMAGE {
					real_message_content := any(message_content).(*modelruntimeentities.ImagePromptMessageContent)
					image_url := real_message_content.Data()
					if strings.HasPrefix(real_message_content.Data(), "data:") {
						// convert image base64 data to file in /tmp
						image_url = m.saveBase64ImageToFile(real_message_content.Data())
					}
					sub_message_dict := map[string]string{"image": image_url}
					sub_messages = append(sub_messages, sub_message_dict)
				} else if message_content.Type() == modelruntimeentities.PromptMessageContent_VIDEO {
					real_message_content := any(message_content).(*modelruntimeentities.VideoPromptMessageContent)
					video_url := real_message_content.URL()
					if video_url == "" {
						return nil, modelruntimeexceptions.NewInvokeError("not support base64, please set MULTIMODAL_SEND_FORMAT to url")
					}
					sub_message_dict := map[string]string{"video": video_url}
					sub_messages = append(sub_messages, sub_message_dict)
				}
			}
			// resort sub_messages to ensure text is always at last
			sort.Slice(sub_messages, func(i int, j int) bool {
				return sub_messages[i]["text"] < sub_messages[j]["text"]
			})

			tongyi_messages = append(tongyi_messages, map[string]any{"role": "user", "content": sub_messages})
		} else if realpromptmsg, ok := any(prompt_message).(*modelruntimeentities.AssistantPromptMessage); ok {
			content := realpromptmsg.Content

			if content == "" {
				content = " "
			}
			message := map[string]any{"role": "assistant"}
			if rich_content {
				message["content"] = []map[string]string{
					{"text": realpromptmsg.Content},
				}
			} else {
				message["content"] = content
			}

			if len(realpromptmsg.ToolCalls) > 0 {
				tool_calls := []map[string]any{}
				for _, tool_call := range realpromptmsg.ToolCalls {
					tool_calls = append(tool_calls, tool_call.ModelDump())
				}
				message["tool_calls"] = tool_calls
			}
			tongyi_messages = append(tongyi_messages, message)
		} else if realpromptmsg, ok := any(prompt_message).(*modelruntimeentities.ToolPromptMessage[string]); ok {
			tongyi_messages = append(tongyi_messages, map[string]any{"role": "tool", "content": realpromptmsg.Content, "name": realpromptmsg.ToolCallID})
		} else {
			return nil, exceptions.NewValueError(fmt.Sprintf("Got unknown type %#v", prompt_message))
		}
	}
	return tongyi_messages, nil
}

func (m *TongyiLargeLanguageModel) saveBase64ImageToFile(base64_image string) string {
	/*
	   Save base64 image to file
	   'data:{upload_file.mime_type};base64,{encoded_string}'
	   :param base64_image: base64 image data
	   :return: image file path
	*/
	// get mime type and encoded string
	tmplist := strings.Split(base64_image, ",")
	mime_type := strings.Split(strings.Split(tmplist[0], ";")[0], ":")[1]
	encoded_string := tmplist[1]
	// mime_type, encoded_string = base64_image.split(",")[0].split(";")[0].split(":")[1], base64_image.split(",")[1]
	// save image to file

	file_path := filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", uuid.NewV4().String(), strings.Split(mime_type, "/")[1]))
	bindata, err := base64.StdEncoding.DecodeString(encoded_string)
	if err != nil {
		mlog.Errorf("base64 decode failed:%v", err)
		return ""
	}
	err = os.WriteFile(file_path, bindata, 0644)
	if err != nil {
		mlog.Errorf("write file failed:%v", err)
		return ""
	}
	return "file://" + file_path
}

func (m *TongyiLargeLanguageModel) convertTools(tools []*modelruntimeentities.PromptMessageTool) []map[string]any {
	/*
	   Convert tools
	*/
	tool_definitions := []map[string]any{}
	for _, tool := range tools {
		properties := map[string]any{}
		if _, ok := tool.Parameters["properties"]; ok {
			if _, ok := tool.Parameters["properties"].(map[string]any); ok {
				properties = tool.Parameters["properties"].(map[string]any)
			}
		}
		var required_properties any
		if _, ok := tool.Parameters["required"]; ok {
			required_properties = tool.Parameters["required"]
		}

		properties_definitions := map[string]any{}
		for p_key, p_val := range properties {
			if _, ok := p_val.(map[string]any); !ok {
				continue
			}
			p_val_dict := p_val.(map[string]any)
			desc := ""
			if _, ok := p_val.(map[string]any)["description"]; ok {
				if _, ok := p_val_dict["description"].(string); ok {
					desc = p_val_dict["description"].(string)
				}
			}

			if _, ok := p_val_dict["enum"]; ok {
				if _, ok := p_val_dict["enum"].([]string); ok {
					desc += fmt.Sprintf("; Only accepts one of the following predefined options: [%s]", strings.Join(p_val_dict["enum"].([]string), ", "))
				} else if _, ok := p_val_dict["enum"].([]any); ok {
					var new_string_enum_list []string
					for _, v := range p_val_dict["enum"].([]any) {
						if _, ok := v.(string); ok {
							if new_string_enum_list == nil {
								new_string_enum_list = make([]string, 0)
							}
							new_string_enum_list = append(new_string_enum_list, v.(string))
						} else {
							new_string_enum_list = nil
							break
						}
					}
					if new_string_enum_list != nil {
						desc += fmt.Sprintf("; Only accepts one of the following predefined options: [%s]", strings.Join(new_string_enum_list, ", "))
					}
				}
			}
			properties_definitions[p_key] = map[string]any{
				"description": desc,
				"type":        p_val_dict["type"],
			}
		}
		tool_definition := map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        tool.Name,
				"description": tool.Description,
				"parameters":  properties_definitions,
				"required":    required_properties,
			},
		}
		tool_definitions = append(tool_definitions, tool_definition)
	}
	return tool_definitions

}

func (m *TongyiLargeLanguageModel) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	/*
	   Architecture for defining customizable models
	   :param model: model name
	   :param credentials: model credentials
	   :return: AIModelEntity or None
	*/
	function_calling_type := ""
	if _, ok := credentials["function_calling_type"]; ok {
		if _, ok := credentials["function_calling_type"].(string); ok {
			function_calling_type = credentials["function_calling_type"].(string)
		}
	}
	context_size := 8000
	if _, ok := credentials["context_size"]; ok {
		if _, ok := credentials["context_size"].(int); ok {
			context_size = credentials["context_size"].(int)
		}
	}
	max_tokens := 1024
	if _, ok := credentials["max_tokens"]; ok {
		if _, ok := credentials["max_tokens"].(int); ok {
			max_tokens = credentials["max_tokens"].(int)
		}
	}
	features := []modelruntimeenumtypes.ModelFeature{}
	if function_calling_type == "tool_call" {
		features = []modelruntimeenumtypes.ModelFeature{modelruntimeenumtypes.ModelFeature_TOOL_CALL, modelruntimeenumtypes.ModelFeature_MULTI_TOOL_CALL, modelruntimeenumtypes.ModelFeature_STREAM_TOOL_CALL}
	}
	return &modelruntimeentities.AIModelEntity{
		Model:     model,
		Label:     commontypes.I18nObject{EnUS: model, ZhHans: model},
		ModelType: modelruntimeenumtypes.Model_LLM,
		Features:  features,
		FetchFrom: modelruntimeenumtypes.FetchFrom_CUSTOMIZABLE_MODEL,
		ModelProperties: map[modelruntimeenumtypes.ModelPropertyKey]any{
			modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE: context_size,
			modelruntimeenumtypes.ModelPropertyKey_MODE:         modelruntimeentities.LLMMode_CHAT,
		},

		ParameterRules: []*modelruntimeentities.ParameterRule{
			{
				Name:        "temperature",
				UseTemplate: "temperature",
				Label:       commontypes.I18nObject{EnUS: "Temperature", ZhHans: "温度"},
				Type:        modelruntimeenumtypes.ParameterType_FLOAT,
			},
			{
				Name:        "max_tokens",
				UseTemplate: "max_tokens",
				Default:     512,
				Min:         1,
				Max:         float64(max_tokens),
				Label:       commontypes.I18nObject{EnUS: "Max Tokens", ZhHans: "最大标记"},
				Type:        modelruntimeenumtypes.ParameterType_INT,
			},
			{
				Name:        "top_p",
				UseTemplate: "top_p",
				Label:       commontypes.I18nObject{EnUS: "Top P", ZhHans: "Top P"},
				Type:        modelruntimeenumtypes.ParameterType_FLOAT,
			},
			{
				Name:        "top_k",
				UseTemplate: "top_k",
				Label:       commontypes.I18nObject{EnUS: "Top K", ZhHans: "Top K"},
				Type:        modelruntimeenumtypes.ParameterType_FLOAT,
			},
			{
				Name:        "frequency_penalty",
				UseTemplate: "frequency_penalty",
				Label:       commontypes.I18nObject{EnUS: "Frequency Penalty", ZhHans: "重复惩罚"},
				Type:        modelruntimeenumtypes.ParameterType_FLOAT,
			},
		},
	}
}
