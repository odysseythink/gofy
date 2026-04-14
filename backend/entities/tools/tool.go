package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	pluginparameter "mlib.com/gofy/server/entities/plugin/parameter"
	providerentities "mlib.com/gofy/server/entities/provider"
	ragentities "mlib.com/gofy/server/entities/rag"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	commontypes "mlib.com/gofy/server/types/common"
)

type Messager interface {
	Type() toolsenumtypes.MessageType
	ToDict() map[string]any
}

type TextMessage struct {
	Text string `json:"text"`
}

func (msg *TextMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_TEXT
}
func (msg *TextMessage) ToDict() map[string]any {
	return map[string]any{
		"type": msg.Type(),
		"text": msg.Text,
	}
}

func (msg TextMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type JsonMessage struct {
	JsonObject map[string]any `json:"json_object"`
}

func (msg *JsonMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_JSON
}
func (msg *JsonMessage) ToDict() map[string]any {
	return map[string]any{
		"type":        msg.Type(),
		"json_object": msg.JsonObject,
	}
}
func (msg JsonMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type BlobMessage struct {
	Blob []byte `json:"blob"`
}

func (msg *BlobMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_BLOB
}
func (msg *BlobMessage) ToDict() map[string]any {
	return map[string]any{
		"type": msg.Type(),
		"blob": base64.RawStdEncoding.EncodeToString(msg.Blob),
	}
}
func (msg BlobMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type BlobChunkMessage struct {
	ID          string `json:"id"`           //description="The id of the blob")
	Sequence    int    `json:"sequence"`     //description="The sequence of the chunk")
	TotalLength int    `json:"total_length"` //description="The total length of the blob")
	Blob        []byte `json:"blob"`         //description="The blob data of the chunk")
	End         bool   `json:"end"`          //description="Whether the chunk is the last chunk")
}

func (msg *BlobChunkMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_BLOB_CHUNK
}
func (msg *BlobChunkMessage) ToDict() map[string]any {
	return map[string]any{
		"type":         msg.Type(),
		"id":           msg.ID,
		"sequence":     msg.Sequence,
		"total_length": msg.TotalLength,
		"blob":         base64.RawStdEncoding.EncodeToString(msg.Blob),
		"end":          msg.End,
	}
}
func (msg BlobChunkMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type FileMessage struct {
}

func (msg *FileMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_FILE
}
func (msg *FileMessage) ToDict() map[string]any {
	return map[string]any{
		"type": msg.Type(),
	}
}
func (msg FileMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type VariableMessage struct {
	VariableName  string `json:"variable_name"`  //description="The name of the variable")
	VariableValue any    `json:"variable_value"` //description="The value of the variable")
	Stream        bool   `json:"stream"`         //description="Whether the variable is streamed")
}

func NewVariableMessage(args any) *VariableMessage {
	newfrommap := func(dict map[string]any) *VariableMessage {
		bindata, _ := json.Marshal(dict)
		msg := new(VariableMessage)
		err := json.Unmarshal(bindata, msg)
		if err != nil {
			mlog.Error("dict {%v} Unmarshal VariableMessage failed:%v", dict, err)
			panic(exceptions.NewValueError(fmt.Sprintf("dict {%v} Unmarshal VariableMessage failed:%v", dict, err)))
		}

		switch msg.VariableValue.(type) {
		case map[string]any, []any, string, int, float64, bool:
		default:
			mlog.Error("Only basic types and lists are allowed.")
			panic(exceptions.NewValueError("Only basic types and lists are allowed."))
		}

		// if stream is true, the value must be a string
		if msg.Stream {
			if _, ok := msg.VariableValue.(string); !ok {
				mlog.Errorf("When 'stream' is True, 'variable_value' must be a string.")
				panic(exceptions.NewValueError("When 'stream' is True, 'variable_value' must be a string."))
			}
		}

		if slices.Contains([]string{"json", "text", "files"}, msg.VariableName) {
			mlog.Errorf("When 'stream' is True, 'variable_value' must be a string.")
			panic(exceptions.NewValueError(fmt.Sprintf("The variable name '{%s}' is reserved.", msg.VariableName)))
		}
		return msg
	}
	switch real_args := args.(type) {
	case []byte:
		tmpdict := map[string]any{}
		if err := json.Unmarshal([]byte(real_args), &tmpdict); err != nil {
			mlog.Errorf("args {%s} must be dict %v", string(real_args), err)
			panic(exceptions.NewValueError(fmt.Sprintf("args {%s} must be dict %v", string(real_args), err)))
		}
		return newfrommap(tmpdict)
	case string:
		tmpdict := map[string]any{}
		if err := json.Unmarshal([]byte(real_args), &tmpdict); err != nil {
			mlog.Errorf("args {%s} must be dict %v", string(real_args), err)
			panic(exceptions.NewValueError(fmt.Sprintf("args {%s} must be dict %v", string(real_args), err)))
		}
		return newfrommap(tmpdict)
	case map[string]any:
		return newfrommap(real_args)
	default:
		mlog.Errorf("args {%#v} is not support ", args)
		panic(exceptions.NewValueError(fmt.Sprintf("args {%#v} is not support ", args)))
	}
}
func (msg *VariableMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_VARIABLE
}
func (msg *VariableMessage) ToDict() map[string]any {
	return map[string]any{
		"type":           msg.Type(),
		"variable_name":  msg.VariableName,
		"variable_value": msg.VariableValue,
		"stream":         msg.Stream,
	}
}
func (msg VariableMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type LogMessage struct {
	ID       string                       `json:"id"`
	Label    string                       `json:"label"`     //description="The label of the log")
	ParentID *string                      `json:"parent_id"` //description="Leave empty for root log")
	Error    *string                      `json:"error"`     //description="The error message")
	Status   toolsenumtypes.LogStatusType `json:"status"`    //description="The status of the log")
	Data     map[string]any               `json:"data"`      //description="Detailed log data")
	Metadata map[string]any               `json:"metadata"`  //description="The metadata of the log")
}

func (msg *LogMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_LOG
}
func (msg *LogMessage) ToDict() map[string]any {
	ret := map[string]any{
		"type":      msg.Type(),
		"id":        msg.ID,
		"label":     msg.Label,
		"status":    msg.Status,
		"data":      msg.Data,
		"metadata":  msg.Metadata,
		"parent_id": "",
		"error":     "",
	}
	if msg.ParentID != nil {
		ret["parent_id"] = *msg.ParentID
	}
	if msg.Error != nil {
		ret["error"] = *msg.Error
	}
	return ret
}
func (msg LogMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type RetrieverResourceMessage struct {
	RetrieverResources []*ragentities.RetrievalSourceMetadata `json:"retriever_resources"` //description="retriever resources")
	Context            string                                 `json:"context"`             //description="context")
}

func (msg *RetrieverResourceMessage) Type() toolsenumtypes.MessageType {
	return toolsenumtypes.Message_RETRIEVER_RESOURCES
}
func (msg *RetrieverResourceMessage) ToDict() map[string]any {
	return map[string]any{
		"type":                msg.Type(),
		"context":             msg.Context,
		"retriever_resources": msg.RetrieverResources,
	}
}
func (msg RetrieverResourceMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type ToolInvokeMessage struct {
	Type    toolsenumtypes.MessageType `json:"type"`
	Message Messager                   `json:"message"`
	Meta    map[string]any             `json:"meta"`
}

func (msg *ToolInvokeMessage) ToDict() map[string]any {
	return map[string]any{
		"meta":    msg.Meta,
		"message": msg.Message.ToDict(),
	}
}
func (msg ToolInvokeMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((&msg).ToDict())
}

type ToolInvokeMessageBinary struct {
	MimeType string         `json:"mimetype"` //description="The mimetype of the binary"
	Url      string         `json:"url"`      //description="The url of the binary"
	FileVar  map[string]any `json:"file_var"`
}

type ToolParameter struct {
	*pluginparameter.PluginParameter

	Type             toolsenumtypes.ToolParameterType     `json:"type"`              //description="The type of the parameter")
	HumanDescription *commontypes.I18nObject              `json:"human_description"` //description="The description presented to the user")
	Form             toolsenumtypes.ToolParameterFormType `json:"form"`              //description="The form of the parameter, schema/form/llm"
	LLMDescription   string                               `json:"llm_description"`
	// MCP object and array type parameters use this field to store the schema
	InputSchema map[string]any `json:"input_schema"`
}

func (tp *ToolParameter) Copy() *ToolParameter {
	tool_parameter := new(ToolParameter)

	tool_parameter.Name = tp.Name
	tool_parameter.Label = commontypes.I18nObject{
		ZhHans: tp.Label.ZhHans,
		EnUS:   tp.Label.EnUS,
		PtBR:   tp.Label.PtBR,
		JaJP:   tp.Label.JaJP,
	}
	if tp.HumanDescription != nil {
		tool_parameter.HumanDescription = &commontypes.I18nObject{
			ZhHans: tp.HumanDescription.ZhHans,
			EnUS:   tp.HumanDescription.EnUS,
			PtBR:   tp.HumanDescription.PtBR,
			JaJP:   tp.HumanDescription.JaJP,
		}
	}
	if tp.Placeholder != nil {
		tool_parameter.Placeholder = &commontypes.I18nObject{
			ZhHans: tp.Placeholder.ZhHans,
			EnUS:   tp.Placeholder.EnUS,
			PtBR:   tp.Placeholder.PtBR,
			JaJP:   tp.Placeholder.JaJP,
		}
	}
	tool_parameter.Type = tp.Type
	tool_parameter.Form = tp.Form
	tool_parameter.LLMDescription = tp.LLMDescription
	tool_parameter.Required = tp.Required
	tool_parameter.Default = tp.Default
	tool_parameter.Min = tp.Min
	tool_parameter.Max = tp.Max
	for _, v := range tp.Options {
		if tool_parameter.Options == nil {
			tool_parameter.Options = make([]*pluginparameter.PluginParameterOption, 0)
		}
		tool_parameter.Options = append(tool_parameter.Options, v.Copy())
	}
	return tool_parameter
}
func (tp *ToolParameter) InitFrontendParameter(value any) any {
	return pluginparameter.InitFrontendParameter(tp.PluginParameter, string(tp.Type), value)
}
func GetSimpleInstance(
	name string,
	llm_description string,
	param_type toolsenumtypes.ToolParameterType,
	required bool,
	options []string,
) *ToolParameter {
	// """
	// get a simple tool parameter

	// :param name: the name of the parameter
	// :param llm_description: the description presented to the LLM
	// :param type: the type of the parameter
	// :param required: if the parameter is required
	// :param options: the options of the parameter
	// """
	// convert options to ToolParameterOption
	// FIXME fix the type error
	var option_objs []*pluginparameter.PluginParameterOption
	for _, v := range options {
		if option_objs == nil {
			option_objs = make([]*pluginparameter.PluginParameterOption, len(options))
		}
		option_objs = append(option_objs, &pluginparameter.PluginParameterOption{
			Value: v,
			Label: commontypes.I18nObject{
				EnUS:   v,
				ZhHans: v,
			},
		})
	}
	return &ToolParameter{
		PluginParameter: &pluginparameter.PluginParameter{
			Name: name,
			Label: commontypes.I18nObject{
				EnUS:   "",
				ZhHans: "",
			},
			Placeholder: nil,
			Required:    required,
			Options:     option_objs, // type: ignore
		},

		HumanDescription: &commontypes.I18nObject{
			EnUS:   "",
			ZhHans: "",
		},
		Type:           param_type,
		Form:           toolsenumtypes.ToolParameterForm_LLM,
		LLMDescription: llm_description,
	}
}

type ToolProviderIdentity struct {
	Author      string                         `json:"author"`      //description="The author of the tool"
	Name        string                         `json:"name"`        //description="The name of the tool"
	Description commontypes.I18nObject         `json:"description"` //description="The description of the tool"
	Icon        string                         `json:"icon"`        //description="The icon of the tool"
	IconDark    string                         `json:"icon_dark"`   //description="The dark icon of the tool"
	Label       commontypes.I18nObject         `json:"label"`       //description="The label of the tool"
	Tags        []toolsenumtypes.ToolLabelType `json:"tags"`        //description="The tags of the tool"
}

type ToolIdentity struct {
	Author   string                 `json:"author"`   //description="The author of the tool"
	Name     string                 `json:"name"`     //description="The name of the tool"
	Label    commontypes.I18nObject `json:"label"`    //description="The label of the tool"
	Provider string                 `json:"provider"` //description="The provider of the tool"
	Icon     string                 `json:"icon"`
}

func NewToolIdentity(data any) *ToolIdentity {
	tool_identity := new(ToolIdentity)
	if data == nil {
		return tool_identity
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_identity)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolIdentity failed, because of provaded invalid initial data"))
		}
		return tool_identity
	} else if _, ok := data.(*ToolIdentity); ok {
		tool_identity.Author = data.(*ToolIdentity).Author
		tool_identity.Name = data.(*ToolIdentity).Name
		tool_identity.Label = commontypes.I18nObject{
			ZhHans: data.(*ToolIdentity).Label.ZhHans,
			EnUS:   data.(*ToolIdentity).Label.EnUS,
			PtBR:   data.(*ToolIdentity).Label.PtBR,
			JaJP:   data.(*ToolIdentity).Label.JaJP,
		}
		tool_identity.Provider = data.(*ToolIdentity).Provider
		tool_identity.Icon = data.(*ToolIdentity).Icon
		return tool_identity
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		return tool_identity
	}
}

type ToolDescription struct {
	Human commontypes.I18nObject `json:"human"` //description="The description presented to the user"
	LLM   string                 `json:"llm"`   //description="The description presented to the LLM"
}

func NewToolDescription(data any) *ToolDescription {
	tool_description := new(ToolDescription)
	if data == nil {
		return tool_description
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_description)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolDescription failed, because of provaded invalid initial data"))
		}
		return tool_description
	} else if _, ok := data.(*ToolDescription); ok {
		tool_description.LLM = data.(*ToolDescription).LLM
		tool_description.Human = commontypes.I18nObject{
			ZhHans: data.(*ToolDescription).Human.ZhHans,
			EnUS:   data.(*ToolDescription).Human.EnUS,
			PtBR:   data.(*ToolDescription).Human.PtBR,
			JaJP:   data.(*ToolDescription).Human.JaJP,
		}
		return tool_description
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		return tool_description
	}
}

type ToolEntity struct {
	Identity             ToolIdentity     `json:"identity"`
	Parameters           []*ToolParameter `json:"parameters"`
	Description          *ToolDescription `json:"description"`
	OutputSchema         map[string]any   `json:"output_schema"`
	HasRuntimeParameters bool             `json:"has_runtime_parameters"` //description="Whether the tool has runtime parameters")

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}

func NewToolEntity(data any) *ToolEntity {
	tool_entity := new(ToolEntity)
	if data == nil {
		return tool_entity
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_entity)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolEntity failed, because of provided invalid initial data"))
		}
		return tool_entity
	} else if _, ok := data.(*ToolEntity); ok {
		tool_entity.Identity = ToolIdentity{
			Author:   data.(*ToolEntity).Identity.Author,
			Name:     data.(*ToolEntity).Identity.Name,
			Label:    data.(*ToolEntity).Identity.Label,
			Provider: data.(*ToolEntity).Identity.Provider,
			Icon:     data.(*ToolEntity).Identity.Icon,
		}
		tool_entity.Parameters = []*ToolParameter{}
		for _, v := range data.(*ToolEntity).Parameters {
			tool_entity.Parameters = append(tool_entity.Parameters, v.Copy())
		}
		if data.(*ToolEntity).Description != nil {
			tool_entity.Description = &ToolDescription{
				Human: data.(*ToolEntity).Description.Human,
				LLM:   data.(*ToolEntity).Description.LLM,
			}
		}
		maps.Copy(tool_entity.OutputSchema, data.(*ToolEntity).OutputSchema)
		tool_entity.HasRuntimeParameters = data.(*ToolEntity).HasRuntimeParameters
		maps.Copy(tool_entity.ModelConfig, data.(*ToolEntity).ModelConfig)
		return tool_entity
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		return tool_entity
	}
}

type OAuthSchema struct {
	ClientSchema      []providerentities.ProviderConfig `json:"client_schema"`      //description="The schema of the OAuth client")
	CredentialsSchema []providerentities.ProviderConfig `json:"credentials_schema"` //description="The schema of the OAuth credentials"
}

type ToolProviderEntity struct {
	Identity          ToolProviderIdentity               `json:"identity"`
	PluginID          string                             `json:"plugin_id"`
	CredentialsSchema []*providerentities.ProviderConfig `json:"credentials_schema"`
	OauthSchema       *OAuthSchema                       `json:"oauth_schema"`
}

type ToolProviderEntityWithPlugin struct {
	*ToolProviderEntity
	Tools []*ToolEntity `json:"tools"`
}

type WorkflowToolParameterConfiguration struct {
	Name        string                               `json:"name"`        //description="The name of the parameter"
	Description string                               `json:"description"` //description="The description of the parameter"
	Form        toolsenumtypes.ToolParameterFormType `json:"form"`        //description="The form of the parameter"
}

type ToolInvokeMeta struct {
	TimeCost   float64        `json:"time_cost"` //description="The time cost of the tool invoke"
	Error      string         `json:"error"`
	ToolConfig map[string]any `json:"tool_config"`
}

func EmptyToolInvokeMeta() *ToolInvokeMeta {
	return &ToolInvokeMeta{
		TimeCost:   0.0,
		Error:      "",
		ToolConfig: map[string]any{},
	}
}
func ErrorToolInvokeMeta(errmsg string) *ToolInvokeMeta {
	return &ToolInvokeMeta{
		TimeCost:   0.0,
		Error:      errmsg,
		ToolConfig: map[string]any{},
	}
}

func (meta *ToolInvokeMeta) ToDict() map[string]any {
	return map[string]any{
		"time_cost":   meta.TimeCost,
		"error":       meta.Error,
		"tool_config": meta.ToolConfig,
	}
}

type ToolLabel struct {
	Name  string                 `json:"name"`  //description="The name of the tool"
	Label commontypes.I18nObject `json:"label"` //description="The label of the tool"
	Icon  string                 `json:"icon"`  //description="The icon of the tool"
}

type ToolSelector[T int | float64 | string] struct {
	ProviderID        string         `json:"provider_id"`        //description="The id of the provider")
	CredentialID      *string        `json:"credential_id"`      //description="The id of the credential")
	ToolName          string         `json:"tool_name"`          //description="The name of the tool")
	ToolDescription   string         `json:"tool_description"`   //description="The description of the tool")
	ToolConfiguration map[string]any `json:"tool_configuration"` //description="Configuration, type form")
	ToolParameters    map[string]struct {
		Name        string                                   `json:"name"`        //description="The name of the parameter"
		Type        toolsenumtypes.ToolParameterType         `json:"type"`        //description="The type of the parameter"
		Required    bool                                     `json:"required"`    //description="Whether the parameter is required"
		Description string                                   `json:"description"` //description="The description of the parameter"
		Default     T                                        `json:"default"`
		Options     []*pluginparameter.PluginParameterOption `json:"options"`
	} `json:"tool_parameters"` //description="Parameters, type llm")
}

func (ts *ToolSelector[T]) GofyModelIdentity() string {
	return TOOL_SELECTOR_MODEL_IDENTITY
}
func (ts *ToolSelector[T]) ToPluginParameter() map[string]any {
	bindata, _ := json.Marshal(ts)
	res := map[string]any{}
	err := json.Unmarshal(bindata, &res)
	if err != nil {
		mlog.Error("unmarshal ToolSelector to map failed:", err)
		return nil
	}
	res["gofy_model_identity"] = ts.GofyModelIdentity()
	return res
}

// type ToolCredentialsOption struct {
// 	Value string                 `json:"value"` //description="The value of the option"
// 	Label commontypes.I18nObject `json:"label"` //description="The label of the option"
// }

// type ToolProviderCredentials struct {
// 	Name        string                         `json:"name"` //description="The name of the credentials"
// 	Type        toolsenumtypes.CredentialsType `json:"type"` //description="The type of the credentials"
// 	Required    bool                           `json:"required"`
// 	Default     any                            `json:"default"`
// 	Options     []*ToolCredentialsOption       `json:"options"`
// 	Label       *commontypes.I18nObject        `json:"label"`
// 	Help        *commontypes.I18nObject        `json:"help"`
// 	URL         string                         `json:"url"`
// 	Placeholder *commontypes.I18nObject        `json:"placeholder"`
// }

// func (tpc *ToolProviderCredentials) ToDict() map[string]any {
// 	return map[string]any{
// 		"name":        tpc.Name,
// 		"type":        tpc.Type,
// 		"required":    tpc.Required,
// 		"default":     tpc.Default,
// 		"options":     tpc.Options,
// 		"help":        tpc.Help,
// 		"label":       tpc.Label,
// 		"url":         tpc.URL,
// 		"placeholder": tpc.Placeholder,
// 	}
// }

// type ToolRuntimeVariabler interface {
// 	Type() toolsenumtypes.ToolRuntimeVariableType
// 	GetName() string
// 	GetPosition() int
// 	GetToolName() string
// }
// type ToolRuntimeVariable struct {
// 	Name     string `json:"name"`      //description="The name of the variable"
// 	Position int    `json:"position"`  //description="The position of the variable"
// 	ToolName string `json:"tool_name"` //description="The name of the tool"
// }

// func (t *ToolRuntimeVariable) GetName() string {
// 	return t.Name
// }
// func (t *ToolRuntimeVariable) GetPosition() int {
// 	return t.Position
// }
// func (t *ToolRuntimeVariable) GetToolName() string {
// 	return t.ToolName
// }

// type ToolRuntimeTextVariable struct {
// 	*ToolRuntimeVariable
// 	Value string `json:"value"` //description="The value of the variable"
// }

// func (variable *ToolRuntimeTextVariable) Type() toolsenumtypes.ToolRuntimeVariableType {
// 	return toolsenumtypes.ToolRuntimeVariable_TEXT
// }
// func (variable ToolRuntimeTextVariable) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		*ToolRuntimeVariable
// 		Value string `json:"value"`
// 		Type  string `json:"type"`
// 	}{
// 		ToolRuntimeVariable: variable.ToolRuntimeVariable,
// 		Value:               variable.Value,
// 		Type:                string((&variable).Type()),
// 	})
// }

// type ToolRuntimeImageVariable struct {
// 	*ToolRuntimeVariable
// 	Value string `json:"value"` //description="The path of the image"
// }

// func (variable *ToolRuntimeImageVariable) Type() toolsenumtypes.ToolRuntimeVariableType {
// 	return toolsenumtypes.ToolRuntimeVariable_IMAGE
// }

// func (variable ToolRuntimeImageVariable) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		*ToolRuntimeVariable
// 		Value string `json:"value"`
// 		Type  string `json:"type"`
// 	}{
// 		ToolRuntimeVariable: variable.ToolRuntimeVariable,
// 		Value:               variable.Value,
// 		Type:                string((&variable).Type()),
// 	})
// }

// type ToolRuntimeVariablePool struct {
// 	ConversationID string `json:"conversation_id"` //description="The conversation id"
// 	UserID         string `json:"user_id"`         //description="The user id"
// 	TenantID       string `json:"tenant_id"`       //description="The tenant id of assistant"

// 	Pool []ToolRuntimeVariabler `json:"pool"` //description="The pool of variables"
// }

// func NewToolRuntimeVariablePool(conversation_id, user_id, tenant_id string, pool []ToolRuntimeVariabler) *ToolRuntimeVariablePool {
// 	return &ToolRuntimeVariablePool{
// 		ConversationID: conversation_id,
// 		UserID:         user_id,
// 		TenantID:       tenant_id,
// 		Pool:           pool,
// 	}

// }

// func (vp *ToolRuntimeVariablePool) SetText(tool_name string, name string, value string) {
// 	for idx, variable := range vp.Pool {
// 		if real_variable, ok := any(variable).(*ToolRuntimeTextVariable); ok {
// 			if real_variable.Name == name {
// 				real_variable := any(variable).(*ToolRuntimeTextVariable)
// 				real_variable.Value = value
// 				vp.Pool[idx] = real_variable
// 				return
// 			}
// 		}
// 	}
// 	variable := &ToolRuntimeTextVariable{
// 		ToolRuntimeVariable: &ToolRuntimeVariable{
// 			Name:     name,
// 			Position: len(vp.Pool),
// 			ToolName: tool_name,
// 		},
// 		Value: value,
// 	}
// 	if vp.Pool == nil {
// 		vp.Pool = make([]ToolRuntimeVariabler, 0)
// 	}
// 	vp.Pool = append(vp.Pool, variable)
// }
// func (vp *ToolRuntimeVariablePool) SetFile(tool_name string, value string, name string) {
// 	// check how many image variables are there
// 	image_variable_count := 0
// 	for _, variable := range vp.Pool {
// 		if variable.Type() == toolsenumtypes.ToolRuntimeVariable_IMAGE {
// 			image_variable_count += 1
// 		}
// 	}
// 	if name == "" {
// 		name = fmt.Sprintf("file_%d", image_variable_count)
// 	}
// 	for idx, variable := range vp.Pool {
// 		if real_variable, ok := any(variable).(*ToolRuntimeImageVariable); ok {
// 			if real_variable.Name == name {
// 				real_variable := any(variable).(*ToolRuntimeImageVariable)
// 				real_variable.Value = value
// 				vp.Pool[idx] = real_variable
// 				return
// 			}
// 		}

// 	}

// 	variable := &ToolRuntimeImageVariable{
// 		ToolRuntimeVariable: &ToolRuntimeVariable{
// 			Name:     name,
// 			Position: len(vp.Pool),
// 			ToolName: tool_name,
// 		},
// 		Value: value,
// 	}

// 	if vp.Pool == nil {
// 		vp.Pool = make([]ToolRuntimeVariabler, 0)
// 	}
// 	vp.Pool = append(vp.Pool, variable)
// }

// type ModelToolConfiguration struct {
// 	// """
// 	// Model tool configuration
// 	// """

// 	Type       string                                      `json:"type"`       //description="The type of the model tool"
// 	Model      string                                      `json:"model"`      //description="The model"
// 	Label      commontypes.I18nObject                      `json:"label"`      //description="The label of the model tool"
// 	Properties map[toolsenumtypes.ModelToolPropertyKey]any `json:"properties"` //description="The properties of the model tool"
// }

// type ModelToolProviderConfiguration struct {
// 	// """
// 	// Model tool provider configuration
// 	// """

// 	Provider string                    `json:"provider"` //description="The provider of the model tool"
// 	Models   []*ModelToolConfiguration `json:"models"`   //description="The models of the model tool"
// 	Label    commontypes.I18nObject    `json:"label"`    //description="The label of the model tool"
// }
