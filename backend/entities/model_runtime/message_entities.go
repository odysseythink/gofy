package modelruntime

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
)

type PromptMessageRole string

/*
Enum class for prompt message.
*/
const (
	PromptMessageRole_SYSTEM    PromptMessageRole = "system"
	PromptMessageRole_USER      PromptMessageRole = "user"
	PromptMessageRole_ASSISTANT PromptMessageRole = "assistant"
	PromptMessageRole_TOOL      PromptMessageRole = "tool"
)

type PromptMessageContentType string

/*
Enum class for prompt message content type.
*/
const (
	PromptMessageContent_TEXT     PromptMessageContentType = "text"
	PromptMessageContent_IMAGE    PromptMessageContentType = "image"
	PromptMessageContent_AUDIO    PromptMessageContentType = "audio"
	PromptMessageContent_VIDEO    PromptMessageContentType = "video"
	PromptMessageContent_DOCUMENT PromptMessageContentType = "document"
)

// PromptMessageTool represents a tool for prompt messages.
type PromptMessageTool struct {
	Name        string
	Description string
	Parameters  map[string]any
}

// PromptMessageFunction represents a function for prompt messages.
type PromptMessageFunction struct {
	Type     string
	Function PromptMessageTool
}

func NewPromptMessageFunction() *PromptMessageFunction {
	return &PromptMessageFunction{
		Type: "function",
	}
}

// PromptMessageContenter represents the content of a prompt message.
type PromptMessageContenter interface {
	Type() PromptMessageContentType
	ToDict() map[string]any
	Data() string
	SetData(string)
}

func NewPromptMessageContenter(args any) PromptMessageContenter {
	if real_args, ok := args.(string); ok && real_args != "" {
		var tmp map[string]any
		if err := json.Unmarshal([]byte(real_args), &tmp); err != nil {
			mlog.Errorf("json unmarshal args=%s to dict failed:%v", real_args, err)
			panic(exceptions.NewValueError(fmt.Sprintf("json unmarshal args=%s to dict failed:%v", real_args, err)))
		}
		return NewPromptMessageContenter(tmp)
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		if _, ok := real_args["type"]; !ok {
			mlog.Errorf("args=%#v must have type field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have type field", real_args)))
		}
		if _, ok := real_args["type"].(string); !ok {
			mlog.Errorf("args=%#v must have string type field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string type field", real_args)))
		}
		switch PromptMessageContentType(real_args["type"].(string)) {
		case PromptMessageContent_TEXT:
			if _, ok := real_args["data"]; !ok {
				mlog.Errorf("args=%#v must have data field", real_args)
				panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have data field", real_args)))
			}
			if _, ok := real_args["data"].(string); !ok {
				mlog.Errorf("args=%#v must have string type field", real_args)
				panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string data field", real_args)))
			}
			return NewTextPromptMessageContent(real_args["data"].(string))
		case PromptMessageContent_IMAGE:
			return NewImagePromptMessageContent(real_args)
		case PromptMessageContent_AUDIO:
			return NewAudioPromptMessageContent(real_args)
		case PromptMessageContent_VIDEO:
			return NewVideoPromptMessageContent(real_args)
		case PromptMessageContent_DOCUMENT:
			return NewDocumentPromptMessageContent(real_args)
		default:
			mlog.Errorf("args[type]=%s not supported", real_args["type"].(string))
			panic(exceptions.NewValueError(fmt.Sprintf("args[type]=%s not supported", real_args["type"].(string))))
		}
	} else {
		mlog.Errorf("unsupported args=%#v", args)
		panic(exceptions.NewValueError(fmt.Sprintf("unsupported args=%#v", args)))
	}
}

// TextPromptMessageContent represents text content for prompt messages.
type TextPromptMessageContent struct {
	data string
}

func (pmc *TextPromptMessageContent) Type() PromptMessageContentType {
	return PromptMessageContent_TEXT
}
func (pmc *TextPromptMessageContent) ToDict() map[string]any {
	return map[string]any{"type": pmc.Type(), "data": pmc.data}
}
func (pmc *TextPromptMessageContent) Data() string {
	return pmc.data
}
func (pmc *TextPromptMessageContent) SetData(val string) {
	pmc.data = val
}
func NewTextPromptMessageContent(data string) *TextPromptMessageContent {
	return &TextPromptMessageContent{
		data: data,
	}
}

// MultiModalPromptMessageContent represents multi-modal content for prompt messages.
type MultiModalPromptMessageContent interface {
	PromptMessageContenter
	Format() string
	Base64Data() string
	URL() string
	MimeType() string
}

type BaseMultiModalPromptMessageContent struct {
	format      string
	base64_data string
	url         string
	mime_type   string
}

func (c *BaseMultiModalPromptMessageContent) Format() string {
	return c.format
}
func (c *BaseMultiModalPromptMessageContent) Base64Data() string {
	return c.base64_data
}
func (c *BaseMultiModalPromptMessageContent) URL() string {
	return c.url
}
func (c *BaseMultiModalPromptMessageContent) MimeType() string {
	return c.mime_type
}

func (c *BaseMultiModalPromptMessageContent) Data() string {
	if c.url != "" {
		return c.url
	}
	return fmt.Sprintf("data:%s;base64,%s", c.mime_type, c.base64_data)
}

func (c *BaseMultiModalPromptMessageContent) SetData(val string) {
	c.url = val
}

func NewBaseMultiModalPromptMessageContent(param map[string]any) *BaseMultiModalPromptMessageContent {
	c := &BaseMultiModalPromptMessageContent{}
	if param != nil {
		if _, ok := param["format"]; ok {
			if _, ok := param["format"].(string); ok {
				c.format = param["format"].(string)
			} else {
				mlog.Errorf("args=%#v must have string format field", param)
				panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string format field", param)))
			}
		}
		if _, ok := param["base64_data"]; ok {
			if _, ok := param["base64_data"].(string); ok {
				c.base64_data = param["base64_data"].(string)
			} else {
				mlog.Errorf("args=%#v must have string base64_data field", param)
				panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string base64_data field", param)))
			}
		}
		if _, ok := param["url"]; ok {
			if _, ok := param["url"].(string); ok {
				c.url = param["url"].(string)
			} else {
				mlog.Errorf("args=%#v must have string url field", param)
				panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string url field", param)))
			}
		}
		if _, ok := param["mime_type"]; ok {
			if _, ok := param["mime_type"].(string); ok {
				c.mime_type = param["mime_type"].(string)
			} else {
				mlog.Errorf("args=%#v must have string mime_type field", param)
				panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string mime_type field", param)))
			}
		}
	}
	return c
}

type VideoPromptMessageContent struct {
	*BaseMultiModalPromptMessageContent
}

func (pmc *VideoPromptMessageContent) ToDict() map[string]any {
	return map[string]any{"type": pmc.Type(), "format": pmc.format, "base64_data": pmc.base64_data, "url": pmc.url, "mime_type": pmc.mime_type}
}
func NewVideoPromptMessageContent(param map[string]any) *VideoPromptMessageContent {
	return &VideoPromptMessageContent{
		BaseMultiModalPromptMessageContent: NewBaseMultiModalPromptMessageContent(param),
	}
}
func (pmc *VideoPromptMessageContent) Type() PromptMessageContentType {
	return PromptMessageContent_VIDEO
}

type AudioPromptMessageContent struct {
	*BaseMultiModalPromptMessageContent
}

func (pmc *AudioPromptMessageContent) ToDict() map[string]any {
	return map[string]any{"type": pmc.Type(), "format": pmc.format, "base64_data": pmc.base64_data, "url": pmc.url, "mime_type": pmc.mime_type}
}
func (pmc *AudioPromptMessageContent) Type() PromptMessageContentType {
	return PromptMessageContent_AUDIO
}
func NewAudioPromptMessageContent(param map[string]any) *AudioPromptMessageContent {
	return &AudioPromptMessageContent{
		BaseMultiModalPromptMessageContent: NewBaseMultiModalPromptMessageContent(param),
	}
}

type ImagePromptMessageContentDETAIL string

const (
	ImagePromptMessageContentDETAIL_LOW  = "low"
	ImagePromptMessageContentDETAIL_HIGH = "high"
)

type ImagePromptMessageContent struct {
	*BaseMultiModalPromptMessageContent
	Detail ImagePromptMessageContentDETAIL
}

func (pmc *ImagePromptMessageContent) ToDict() map[string]any {
	return map[string]any{"type": pmc.Type(), "detail": pmc.Detail, "format": pmc.format, "base64_data": pmc.base64_data, "url": pmc.url, "mime_type": pmc.mime_type}
}
func (pmc *ImagePromptMessageContent) Type() PromptMessageContentType {
	return PromptMessageContent_IMAGE
}
func NewImagePromptMessageContent(param map[string]any) *ImagePromptMessageContent {
	c := &ImagePromptMessageContent{
		BaseMultiModalPromptMessageContent: NewBaseMultiModalPromptMessageContent(param),
		Detail:                             ImagePromptMessageContentDETAIL_LOW,
	}
	if param != nil {
		if _, ok := param["detail"].(string); ok {
			c.Detail = ImagePromptMessageContentDETAIL(param["detail"].(string))
		} else {
			mlog.Errorf("args=%#v must have string detail field", param)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string detail field", param)))
		}
	}
	return c
}

type DocumentPromptMessageContent struct {
	*BaseMultiModalPromptMessageContent
}

func (pmc *DocumentPromptMessageContent) ToDict() map[string]any {
	return map[string]any{"type": pmc.Type(), "format": pmc.format, "base64_data": pmc.base64_data, "url": pmc.url, "mime_type": pmc.mime_type}
}
func (pmc *DocumentPromptMessageContent) Type() PromptMessageContentType {
	return PromptMessageContent_DOCUMENT
}
func NewDocumentPromptMessageContent(param map[string]any) *DocumentPromptMessageContent {
	return &DocumentPromptMessageContent{
		BaseMultiModalPromptMessageContent: NewBaseMultiModalPromptMessageContent(param),
	}
}

// PromptMessager is an abstract base class for prompt messages.
type PromptMessager interface {
	IsEmpty() bool
	Role() PromptMessageRole
	GetContent() any
	SetContent(any) error
	Name() string
	ToDict() map[string]any
}

func NewPromptMessager(args any) PromptMessager {
	if real_args, ok := args.(string); ok && real_args != "" {
		var tmp map[string]any
		if err := json.Unmarshal([]byte(real_args), &tmp); err != nil {
			mlog.Errorf("json unmarshal args=%s to dict failed:%v", real_args, err)
			panic(exceptions.NewValueError(fmt.Sprintf("json unmarshal args=%s to dict failed:%v", real_args, err)))
		}
		return NewPromptMessager(tmp)
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		if _, ok := real_args["role"]; !ok {
			mlog.Errorf("args=%#v must have role field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have role field", real_args)))
		}
		if _, ok := real_args["role"].(string); !ok {
			mlog.Errorf("args=%#v must have string role field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string role field", real_args)))
		}
		if _, ok := real_args["name"]; !ok {
			mlog.Errorf("args=%#v must have name field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have name field", real_args)))
		}
		if _, ok := real_args["name"].(string); !ok {
			mlog.Errorf("args=%#v must have string name field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string name field", real_args)))
		}
		name := real_args["name"].(string)
		if _, ok := real_args["content"]; !ok {
			mlog.Errorf("args=%#v must have content field", real_args)
			panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have content field", real_args)))
		}
		switch real_content := real_args["content"].(type) {
		case string:
			switch PromptMessageRole(real_args["role"].(string)) {
			case PromptMessageRole_SYSTEM:
				return NewSystemPromptMessage(real_content, name)
			case PromptMessageRole_USER:
				return NewUserPromptMessage(real_content, name)
			case PromptMessageRole_ASSISTANT:
				if _, ok := real_args["tool_calls"]; !ok {
					mlog.Errorf("args=%#v must have tool_calls field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have tool_calls field", real_args)))
				}
				tool_calls_dict := []map[string]any{}
				if _, ok := real_args["tool_calls"].([]any); !ok {
					if _, ok := real_args["tool_calls"].([]map[string]any); !ok {
						mlog.Errorf("args=%#v must have []map[string]any tool_calls field", real_args)
						panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have []map[string]any tool_calls field", real_args)))
					} else {
						tool_calls_dict = real_args["tool_calls"].([]map[string]any)
					}
				} else {
					for _, v := range real_args["tool_calls"].([]any) {
						if _, ok := v.(map[string]any); !ok {
							mlog.Errorf("args=%#v must have []map[string]any tool_calls field", real_args)
							panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have []map[string]any tool_calls field", real_args)))
						} else {
							tool_calls_dict = append(tool_calls_dict, v.(map[string]any))
						}
					}
				}
				tool_calls := []*ToolCall{}
				for _, v := range tool_calls_dict {
					tool_calls = append(tool_calls, NewToolCall(v))
				}
				return NewAssistantPromptMessage(real_content, name, tool_calls)
			case PromptMessageRole_TOOL:
				if _, ok := real_args["tool_call_id"]; !ok {
					mlog.Errorf("args=%#v must have tool_call_id field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have tool_call_id field", real_args)))
				}
				if _, ok := real_args["tool_call_id"].(string); !ok {
					mlog.Errorf("args=%#v must have string tool_call_id field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string tool_call_id field", real_args)))
				}
				tool_call_id := real_args["tool_call_id"].(string)
				return NewToolPromptMessage(real_content, name, tool_call_id)
			default:
				mlog.Errorf("unsupported role=%#v", real_args["role"].(string))
				panic(exceptions.NewValueError(fmt.Sprintf("unsupported role=%#v", real_args["role"].(string))))
			}
		case []any:
			pmcs := []PromptMessageContenter{}
			for _, v := range real_content {
				pmcs = append(pmcs, NewPromptMessageContenter(v))
			}
			switch PromptMessageRole(real_args["role"].(string)) {
			case PromptMessageRole_SYSTEM:
				return NewSystemPromptMessage(pmcs, name)
			case PromptMessageRole_USER:
				return NewUserPromptMessage(pmcs, name)
			case PromptMessageRole_ASSISTANT:
				mlog.Error("assistant prompt message don't supported []PromptMessageContenter content")
				panic(exceptions.NewValueError("assistant prompt message don't supported []PromptMessageContenter content"))
			case PromptMessageRole_TOOL:
				if _, ok := real_args["tool_call_id"]; !ok {
					mlog.Errorf("args=%#v must have tool_call_id field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have tool_call_id field", real_args)))
				}
				if _, ok := real_args["tool_call_id"].(string); !ok {
					mlog.Errorf("args=%#v must have string tool_call_id field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string tool_call_id field", real_args)))
				}
				tool_call_id := real_args["tool_call_id"].(string)
				return NewToolPromptMessage(pmcs, name, tool_call_id)
			default:
				mlog.Errorf("unsupported role=%#v", real_args["role"].(string))
				panic(exceptions.NewValueError(fmt.Sprintf("unsupported role=%#v", real_args["role"].(string))))
			}
		case []map[string]any:
			pmcs := []PromptMessageContenter{}
			for _, v := range real_content {
				pmcs = append(pmcs, NewPromptMessageContenter(v))
			}
			switch PromptMessageRole(real_args["role"].(string)) {
			case PromptMessageRole_SYSTEM:
				return NewSystemPromptMessage(pmcs, name)
			case PromptMessageRole_USER:
				return NewUserPromptMessage(pmcs, name)
			case PromptMessageRole_ASSISTANT:
				mlog.Error("assistant prompt message don't supported []PromptMessageContenter content")
				panic(exceptions.NewValueError("assistant prompt message don't supported []PromptMessageContenter content"))
			case PromptMessageRole_TOOL:
				if _, ok := real_args["tool_call_id"]; !ok {
					mlog.Errorf("args=%#v must have tool_call_id field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have tool_call_id field", real_args)))
				}
				if _, ok := real_args["tool_call_id"].(string); !ok {
					mlog.Errorf("args=%#v must have string tool_call_id field", real_args)
					panic(exceptions.NewValueError(fmt.Sprintf("args=%#v must have string tool_call_id field", real_args)))
				}
				tool_call_id := real_args["tool_call_id"].(string)
				return NewToolPromptMessage(pmcs, name, tool_call_id)
			default:
				mlog.Errorf("unsupported role=%#v", real_args["role"].(string))
				panic(exceptions.NewValueError(fmt.Sprintf("unsupported role=%#v", real_args["role"].(string))))
			}
		default:
			mlog.Errorf("unsupported content=%#v", real_content)
			panic(exceptions.NewValueError(fmt.Sprintf("unsupported content=%#v", real_content)))
		}

	} else {
		mlog.Errorf("unsupported args=%#v", args)
		panic(exceptions.NewValueError(fmt.Sprintf("unsupported args=%#v", args)))
	}
}

// BasePromptMessage provides common functionality for prompt messages.
type BasePromptMessage[T string | []PromptMessageContenter] struct {
	Content T // [T string | []PromptMessageContenter]
	name    string
}

func NewBasePromptMessage[T string | []PromptMessageContenter](content T, name string) *BasePromptMessage[T] {
	return &BasePromptMessage[T]{
		Content: content,
		name:    name,
	}
}

func (p *BasePromptMessage[T]) Name() string {
	return p.name
}

func (p *BasePromptMessage[T]) GetContent() any {
	return p.Content
}

func (p *BasePromptMessage[T]) SetContent(val any) error {
	if _, ok := val.(T); ok {
		p.Content = val.(T)
		return nil
	} else {
		mlog.Errorf("val(%#v) must be type=%s", val, reflect.TypeOf(p.Content).Name())
		return exceptions.NewValueError(fmt.Sprintf("val(%#v) must be type=%s", val, reflect.TypeOf(p.Content).Name()))
	}

}

// IsEmpty checks if the prompt message is empty.
func (p *BasePromptMessage[T]) IsEmpty() bool {
	if real_content, ok := any(p.Content).(string); ok {
		return real_content == ""
	} else if real_content, ok := any(p.Content).([]PromptMessageContenter); ok {
		return len(real_content) == 0
	}
	return false
}

// type PromptMessage(ABC, BaseModel) struct {
// /*
//     Model class for prompt message.
//     */

//     role: PromptMessageRole
//     content: Optional[str | Sequence[PromptMessageContenter]] = None
//     name: Optional[str] = None

//     def is_empty(self) -> bool:
//         """
//         Check if prompt message is empty.

//         :return: True if prompt message is empty, False otherwise
//         """
//         return not self.content
// }

// UserPromptMessage represents a user prompt message.
type UserPromptMessage[T string | []PromptMessageContenter] struct {
	*BasePromptMessage[T]
}

// NewUserPromptMessage creates a new UserPromptMessage instance.
func NewUserPromptMessage[T string | []PromptMessageContenter](content T, name string) *UserPromptMessage[T] {
	return &UserPromptMessage[T]{
		BasePromptMessage: NewBasePromptMessage(content, name),
	}
}
func (msg *UserPromptMessage[T]) ToDict() map[string]any {
	rsp := map[string]any{"role": msg.Role(), "name": msg.name}
	switch real_content := any(msg.Content).(type) {
	case string:
		rsp["content"] = real_content
	case []PromptMessageContenter:
		dict_list := []map[string]any{}
		for _, v := range real_content {
			dict_list = append(dict_list, v.ToDict())
		}
		rsp["content"] = dict_list
	}
	return rsp
}
func (msg *UserPromptMessage[T]) Role() PromptMessageRole {
	return PromptMessageRole_USER
}

// AssistantPromptMessage represents an assistant prompt message.
type AssistantPromptMessage struct {
	*BasePromptMessage[string]
	ToolCalls []*ToolCall
}

// ToolCall represents a tool call for assistant prompt messages.
type ToolCall struct {
	ID       string
	Type     string
	Function ToolCallFunction
}

func NewToolCall(args any) *ToolCall {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(ToolCall)
		if err := json.Unmarshal([]byte(real_args), rsp); err != nil {
			mlog.Errorf("json unmarshal args=%s to ToolCall failed:%v", real_args, err)
			panic(exceptions.NewValueError(fmt.Sprintf("json unmarshal args=%s to ToolCall failed:%v", real_args, err)))
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		bindata, _ := json.Marshal(real_args)
		return NewToolCall(string(bindata))
	} else {
		mlog.Errorf("unsupported args=%#v", args)
		panic(exceptions.NewValueError(fmt.Sprintf("unsupported args=%#v", args)))
	}
}
func (tc *ToolCall) ModelDump() map[string]any {
	bindata, _ := json.Marshal(tc)
	rc := map[string]any{}
	json.Unmarshal(bindata, &rc)
	return rc
}

// ToolCallFunction represents a function for tool calls.
type ToolCallFunction struct {
	Name      string
	Arguments string
}

// NewAssistantPromptMessage creates a new AssistantPromptMessage instance.
func NewAssistantPromptMessage(content string, name string, toolCalls []*ToolCall) *AssistantPromptMessage {
	return &AssistantPromptMessage{
		BasePromptMessage: NewBasePromptMessage(content, name),
		ToolCalls:         toolCalls,
	}
}
func (msg *AssistantPromptMessage) ToDict() map[string]any {
	rsp := map[string]any{"role": msg.Role(), "name": msg.name, "content": msg.Content}
	dict_list := []map[string]any{}
	for _, v := range msg.ToolCalls {
		dict_list = append(dict_list, v.ModelDump())
	}
	rsp["tool_calls"] = dict_list
	return rsp
}
func (msg *AssistantPromptMessage) Role() PromptMessageRole {
	return PromptMessageRole_ASSISTANT
}

// IsEmpty checks if the assistant prompt message is empty.
func (a *AssistantPromptMessage) IsEmpty() bool {
	if !a.BasePromptMessage.IsEmpty() || len(a.ToolCalls) > 0 {
		return false
	}
	return true
}

// SystemPromptMessage represents a system prompt message.
type SystemPromptMessage[T string | []PromptMessageContenter] struct {
	*BasePromptMessage[T]
}

// NewSystemPromptMessage creates a new SystemPromptMessage instance.
func NewSystemPromptMessage[T string | []PromptMessageContenter](content T, name string) *SystemPromptMessage[T] {
	return &SystemPromptMessage[T]{
		BasePromptMessage: NewBasePromptMessage(content, name),
	}
}
func (msg *SystemPromptMessage[T]) ToDict() map[string]any {
	rsp := map[string]any{"role": msg.Role(), "name": msg.name}
	switch real_content := any(msg.Content).(type) {
	case string:
		rsp["content"] = real_content
	case []PromptMessageContenter:
		dict_list := []map[string]any{}
		for _, v := range real_content {
			dict_list = append(dict_list, v.ToDict())
		}
		rsp["content"] = dict_list
	}
	return rsp
}
func (msg *SystemPromptMessage[T]) Role() PromptMessageRole {
	return PromptMessageRole_SYSTEM
}

// ToolPromptMessage represents a tool prompt message.
type ToolPromptMessage[T string | []PromptMessageContenter] struct {
	*BasePromptMessage[T]
	ToolCallID string
}

// NewToolPromptMessage creates a new ToolPromptMessage instance.
func NewToolPromptMessage[T string | []PromptMessageContenter](content T, name string, toolCallID string) *ToolPromptMessage[T] {
	return &ToolPromptMessage[T]{
		BasePromptMessage: NewBasePromptMessage(content, name),
		ToolCallID:        toolCallID,
	}
}
func (msg *ToolPromptMessage[T]) ToDict() map[string]any {
	rsp := map[string]any{"role": msg.Role(), "name": msg.name, "tool_call_id": msg.ToolCallID}
	switch real_content := any(msg.Content).(type) {
	case string:
		rsp["content"] = real_content
	case []PromptMessageContenter:
		dict_list := []map[string]any{}
		for _, v := range real_content {
			dict_list = append(dict_list, v.ToDict())
		}
		rsp["content"] = dict_list
	}
	return rsp
}
func (msg *ToolPromptMessage[T]) Role() PromptMessageRole {
	return PromptMessageRole_TOOL
}

// IsEmpty checks if the tool prompt message is empty.
func (t *ToolPromptMessage[T]) IsEmpty() bool {
	if !t.BasePromptMessage.IsEmpty() || t.ToolCallID != "" {
		return false
	}
	return true
}
