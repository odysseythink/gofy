package llm

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	uuid "github.com/satori/go.uuid"
	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	dashscopeexception "mlib.com/gofy/server/core/exceptions/dashscope"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntime "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/gofy/server/utils/validate"
	"mlib.com/mlog"
)

type ResponseMetadata struct {
	RequestId string `json:"RequestId"`
	Action    string `json:"Action"`
	Version   string `json:"Version"`
	Service   string `json:"Service"`
	Region    string `json:"Region"`
	Error     struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error"`
}
type MessageType int

const (
	Message_INVALID                   MessageType = 0   //默认非法值
	Message_SIMPLE_AGENT_ORIGIN       MessageType = 10  //simple agent 的原始输出
	Message_SIMPLE_AGENT_ORIGIN_SLICE MessageType = 11  //simple agent 的原始输出分片，后缀 DELTA 代表分片输出
	Message_RAG_SUM                   MessageType = 100 //rag 总结回复
	Message_RAG_SUM_SLICE             MessageType = 101 //rag 总结回复的分片，后缀 DELTA 代表分片输出
	Message_RAG_RECOMMEND_QUSTION     MessageType = 102 //rag 推荐问
	Message_RAG_RESULT                MessageType = 103 //rag 召回结果"
)

type MixtureContentType int

const (
	MixtureContent_AUDIO      MixtureContentType = 0 // 音频
	MixtureContent_TEXT       MixtureContentType = 1 // 文本
	MixtureContent_VIDEO      MixtureContentType = 2 // 视频
	MixtureContent_IMAGE_TEXT MixtureContentType = 3 // 图文混合消息/图片
)

type Image struct {
	URL    string `json:"URL"`
	FileID string `json:"FileID"`
}
type Video struct {
	URL    string `json:"URL"`
	FileID string `json:"FileID"`
}
type MixtureContent struct {
	Type  MixtureContentType `json:"Type"`
	Text  string             `json:"Text"`
	Image Image              `json:"Image"`
	Video Video              `json:"Video"`
}
type OutputMessage struct {
	ConversationID  string            `json:"ConversationID"`
	RoundID         string            `json:"RoundID"`
	MessageID       string            `json:"MessageID"`
	Text            string            `json:"Text"`
	Type            MessageType       `json:"Type"`
	FinishReason    int               `json:"FinishReason"`
	OutputType      MessageType       `json:"OutputType"`
	MixtureContents []*MixtureContent `json:"MixtureContents"`
}
type LLMModelUsage struct {
	CompletionTokens int64 `json:"CompletionTokens"`
	PromptTokens     int64 `json:"PromptTokens"`
	TotalTokens      int64 `json:"TotalTokens"`
}
type AgentResult struct {
	AgentType string           `json:"AgentType"`
	AgentID   string           `json:"AgentID"`
	Outputs   []*OutputMessage `json:"Outputs"`
	Usage     *LLMModelUsage   `json:"usage"`
}
type AssistantSyncAnalysisResult struct {
	AgentResults []*AgentResult `json:"AgentResults"`
	FailAgentIds []string       `json:"FailAgentIds"`
}
type Response struct {
	ResponseMetadata ResponseMetadata             `json:"ResponseMetadata"`
	Result           *AssistantSyncAnalysisResult `json:"Result"`
}

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

type LadderLargeLanguageModel struct {
	*base.LargeLanguageModel
	tokenizers map[string]*tiktoken.Tiktoken
}

const (
	LADDER_URL = "https://speech.bytedance.com/ih/robot/llm/assistant/sync"
)

func GenStringSignature(body []byte, accountId, userId, secretKey string) string {
	return signature(getSignature(body, accountId, userId, secretKey))
}

func signature(kSigning []byte) string {
	return hex.EncodeToString(kSigning)
}

func getSignature(body []byte, accountId, userId, secretKey string) []byte {
	kAccount := hmacSHA256(body, accountId)
	kUserId := hmacSHA256([]byte(kAccount), userId)
	kSigning := hmacSHA256([]byte(kUserId), secretKey)
	return []byte(kSigning)
}
func hmacSHA256(key []byte, content string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

func LadderRunChatMessagesApp(
	sessionID,
	account_id,
	user_id,
	secret_key,
	app_id string,
	agent_dynamic_variables map[string]string,
	query string,
) (*Response, error) {
	if account_id == "" || user_id == "" || secret_key == "" || app_id == "" {
		mlog.Errorf("invalid arg")
		return nil, fmt.Errorf("invalid args")
	}
	data := map[string]any{
		"User": map[string]any{
			"AccountId": account_id,
			"UserId":    user_id,
			"AppId":     app_id,
			"BotId":     "online",
		},
		"Type":   1,
		"TextId": uuid.NewV4().String(),
		"PlainText": map[string]any{
			"Content": query,
		},
		"agentDynamicVariables": agent_dynamic_variables,
	}
	bytesData, _ := json.Marshal(data)
	sign := GenStringSignature(bytesData, account_id, user_id, secret_key)

	ladder_url := confy.GetWithDefault("ai.ladder.chat_messages_url", LADDER_URL)
	// transport := &http.Transport{
	// 	ResponseHeaderTimeout: 60 * time.Second,
	// 	IdleConnTimeout:       30 * time.Second,
	// }
	// client := &http.Client{
	// 	Transport: transport,
	// 	Timeout:   60 * time.Second,
	// }
	request, err := http.NewRequest("POST", ladder_url, bytes.NewReader(bytesData))
	if err != nil {
		mlog.Errorf("session(%s) http post url(%s) with param(%s) failed:%v", sessionID, ladder_url, string(bytesData), err)
		return nil, fmt.Errorf("http post url(%s) with param(%s) failed:%v", ladder_url, string(bytesData), err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Top-Account-Id", account_id)
	request.Header.Set("X-Top-User-Id", user_id)
	request.Header.Set("X-Top-Signature", sign)
	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		mlog.Errorf("session(%s) httputil.DumpRequestOut failed:%v", sessionID, err)
		return nil, fmt.Errorf("httputil.DumpRequestOut failed:%v", err)
	}
	mlog.Debugf("---session(%s)---request=%s", sessionID, string(dump))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		mlog.Errorf("session(%s) http.DefaultClient.Do failed:%v", sessionID, err)
		return nil, fmt.Errorf("http.DefaultClient.Do failed:%v", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		mlog.Errorf("session(%s) io.ReadAll response failed:%v", sessionID, err)
		return nil, fmt.Errorf("io.ReadAll response failed failed:%v", err)
	}
	fmt.Println(string(respBody))
	rsp := new(Response)
	err = json.Unmarshal(respBody, rsp)
	if err != nil {
		mlog.Errorf("session(%s) json.Unmarshal(%s) failed:%v\n", sessionID, string(respBody), err)
		return nil, fmt.Errorf("json.Unmarshal(%s) failed:%v", string(respBody), err)
	}
	mlog.Debugf("---session(%s)---response=%s", sessionID, string(respBody))
	if rsp.ResponseMetadata.Error.Code != "" {
		mlog.Errorf("session(%s) response return failed:%v\n", sessionID, rsp.ResponseMetadata.Error.Message)
		return nil, fmt.Errorf("response return failed:%v", rsp.ResponseMetadata.Error.Message)
	}
	if len(rsp.Result.AgentResults) < 1 {
		mlog.Errorf("session(%s) response return no AgentResults", sessionID)
		return nil, errors.New("response return no AgentResults")
	}
	if len(rsp.Result.AgentResults[0].Outputs) < 1 {
		mlog.Errorf("session(%s) response return no Outputs", sessionID)
		return nil, errors.New("response return no Outputs")
	}
	return rsp, nil
}

func (m *LadderLargeLanguageModel) InvokeStream(
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
		index := 0
		for response := range responses {
			mlog.Debugf("------response=%#v", response)
			if response.Result == nil {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			if len(response.Result.AgentResults) == 0 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			if len(response.Result.AgentResults[0].Outputs) == 0 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			if response.Result.AgentResults[0].Outputs[0].Text == "" {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			resp_content := response.Result.AgentResults[0].Outputs[0].Text

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
func (m *LadderLargeLanguageModel) Invoke(
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
func (m *LadderLargeLanguageModel) GetNumTokens(
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

func (m *LadderLargeLanguageModel) ValidateCredentials(model string, credentials map[string]any) {
	/*
	   Validate model credentials
	   :param model: model name
	   :param credentials: model credentials
	   :return
	*/
	mlog.Debugf("validate model %s with credentials %v", model, credentials)
	m.generate(
		model,
		credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewSystemPromptMessage("ping", "")},
		map[string]any{
			"temperature": 0.5,
		},
		nil,
		nil,
		false,
		"",
	)
}

func (m *LadderLargeLanguageModel) ProviderName() string {
	return "ladder"
}

func (m *LadderLargeLanguageModel) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (m *LadderLargeLanguageModel) generate_call(
	model string,
	credentials map[string]any,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_parameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	stream bool, /*= True*/
) iter.Seq[*Response] {
	return func(yield func(*Response) bool) {
		mlog.Debug("model %s", model)
		err := validate.StringMapTypeVerify(credentials, validate.Rules{
			"user_id":    {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
			"account_id": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
			"secret_key": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
			"app_id":     {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
			"agent_id":   {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("credentials %#v is invalid:%v", credentials, err)
			panic(exceptions.NewValueError("credentials is invalid: " + err.Error()))
		}
		if len(prompt_messages) == 0 {
			mlog.Errorf("empty prompt")
			panic(exceptions.NewValueError("empty prompt"))
		}
		if prompt_messages[0].Role() != modelruntimeentities.PromptMessageRole_SYSTEM {
			mlog.Errorf("first prompt message role is not system")
			panic(exceptions.NewValueError("first prompt message role is not system"))
		}
		tmp := map[string]any{
			"sys_prompt": prompt_messages[0].GetContent(),
			"model_name": model,
		}
		if len(prompt_messages) > 1 && prompt_messages[1].Role() == modelruntimeentities.PromptMessageRole_USER {
			mlog.Debugf("----user prompt message=%#v", prompt_messages[1])
			mlog.Debugf("----user prompt message content=%#v", prompt_messages[1].GetContent())
			if content, ok := prompt_messages[1].GetContent().((*modelruntime.TextPromptMessageContent)); ok {
				tmp["user_prompt"] = content.Data()
			}
		}
		bindata, _ := json.Marshal(tmp)
		agent_dynamic_variables := map[string]string{credentials["agent_id"].(string): string(bindata)}
		rsp, err := LadderRunChatMessagesApp(
			uuid.NewV4().String(),
			credentials["account_id"].(string),
			credentials["user_id"].(string),
			credentials["secret_key"].(string),
			credentials["app_id"].(string),
			agent_dynamic_variables,
			"test",
		)

		if err != nil {
			mlog.Error("ladder run app %s failed:", credentials["app_id"].(string), err)
			panic(exceptions.NewValueError(fmt.Sprintf("ladder run app %s failed:%v", credentials["app_id"].(string), err)))
		}

		if !yield(rsp) {
			return
		}

	}
}

func (m *LadderLargeLanguageModel) generate(
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

func (m *LadderLargeLanguageModel) handleGenerateResponse(
	model string, credentials map[string]any, responses iter.Seq[*Response], prompt_messages []modelruntimeentities.PromptMessager,
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
		if response.Result == nil {
			panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
		}
		if len(response.Result.AgentResults) == 0 {
			panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
		}
		if len(response.Result.AgentResults[0].Outputs) == 0 {
			panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
		}
		if response.Result.AgentResults[0].Outputs[0].Text == "" {
			panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
		}
		resp_content := response.Result.AgentResults[0].Outputs[0].Text

		assistant_prompt_message := modelruntimeentities.NewAssistantPromptMessage(resp_content, "", nil)
		// transform usage
		usage := m.CalcResponseUsage(m, model, credentials, int(response.Result.AgentResults[0].Usage.PromptTokens), int(response.Result.AgentResults[0].Usage.CompletionTokens))
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

func (m *LadderLargeLanguageModel) handleGenerateStreamResponse(
	model string,
	credentials map[string]any,
	responses iter.Seq[*Response],
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
		index := 0
		for response := range responses {
			mlog.Debugf("------response=%#v", response)
			if response.Result == nil {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			if len(response.Result.AgentResults) == 0 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			if len(response.Result.AgentResults[0].Outputs) == 0 {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}
			if response.Result.AgentResults[0].Outputs[0].Text == "" {
				panic(dashscopeexception.NewServiceUnavailableError(fmt.Sprintf("Failed to invoke model %s, status code: %v, message: %s", model, response.ResponseMetadata.Error.Code, response.ResponseMetadata.Error.Message)))
			}

			resp_content := response.Result.AgentResults[0].Outputs[0].Text // response.Output.Choices[0].Message.content

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

func (m *LadderLargeLanguageModel) toCredentialKwargs(credentials map[string]any) map[string]any {
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

func (m *LadderLargeLanguageModel) _convert_one_message_to_text(message modelruntimeentities.PromptMessager) (string, error) {
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

func (m *LadderLargeLanguageModel) _convert_messages_to_prompt(messages []modelruntimeentities.PromptMessager) string {
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

func (m *LadderLargeLanguageModel) _convert_prompt_messages_to_tongyi_messages(
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

func (m *LadderLargeLanguageModel) saveBase64ImageToFile(base64_image string) string {
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

func (m *LadderLargeLanguageModel) convertTools(tools []*modelruntimeentities.PromptMessageTool) []map[string]any {
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

func (m *LadderLargeLanguageModel) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
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
