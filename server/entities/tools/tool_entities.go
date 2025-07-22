package tools

import (
	"encoding/json"
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/mlog"
)

type ToolInvokeMessage struct {
	Type toolsenumtypes.MessageType `json:"id"`
	/*
	   plain text, image url or link url
	*/
	Message any `json:"message"`
	// TODO: Use a BaseModel for meta
	Meta   map[string]any `json:"meta"`
	SaveAs string         `json:"save_as"`
}

func NewToolInvokeMessage() *ToolInvokeMessage {
	return &ToolInvokeMessage{
		Type: toolsenumtypes.Message_TEXT,
	}
}

type ToolInvokeMessageBinary struct {
	MimeType string/*/*= Field(..., description="The mimetype of the binary")*/ `json:"mimetype"`
	Url      string/*/*= Field(..., description="The url of the binary")*/ `json:"url"`
	SaveAs   string         `json:"save_as"`
	FileVar  map[string]any `json:"file_var"`
}

type ToolParameterOption struct {
	Value string/*/*= Field(..., description="The value of the option")*/ `json:"value"`
	Label commontypes.I18nObject/*/*= Field(..., description="The label of the option")*/ `json:"label"`
}

func NewToolParameterOption(data any) *ToolParameterOption {
	tool_parameter_option := new(ToolParameterOption)
	if data == nil {
		return tool_parameter_option
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_parameter_option)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolParameterOption failed, because of provaded invalid initial data"))
		}
		return tool_parameter_option
	} else if _, ok := data.(*ToolParameterOption); ok {
		tool_parameter_option.Value = data.(*ToolParameterOption).Value
		tool_parameter_option.Label = commontypes.I18nObject{
			ZhHans: data.(*ToolParameterOption).Label.ZhHans,
			EnUS:   data.(*ToolParameterOption).Label.EnUS,
			PtBR:   data.(*ToolParameterOption).Label.PtBR,
			JaJP:   data.(*ToolParameterOption).Label.JaJP,
		}
		return tool_parameter_option
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		return tool_parameter_option
	}
}

type ToolParameter struct {

	// deprecated, should not use.
	// SYSTEM_FILES = "systme-files"

	Name             string/*= Field(..., description="The name of the parameter")*/ `json:"name"`
	Label            commontypes.I18nObject/*= Field(..., description="The label presented to the user")*/ `json:"label"`
	HumanDescription *commontypes.I18nObject/*= Field(None, description="The description presented to the user")*/ `json:"human_description"`
	Placeholder      *commontypes.I18nObject/*= Field(None, description="The placeholder presented to the user")*/ `json:"placeholder"`
	Type             toolsenumtypes.ToolParameterType/*= Field(..., description="The type of the parameter")*/ `json:"type"`
	Form             toolsenumtypes.ToolParameterForm/*= Field(..., description="The form of the parameter, schema/form/llm")*/ `json:"form"`
	LLMDescription   string                 `json:"llm_description"`
	Required         bool                   `json:"required"`
	Default          any                    `json:"default"`
	Min              float64                `json:"min"`
	Max              float64                `json:"max"`
	Options          []*ToolParameterOption `json:"options"`
}

func NewToolParameter(data any) *ToolParameter {
	tool_parameter := new(ToolParameter)
	if data == nil {
		return tool_parameter
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_parameter)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolParameter failed, because of provaded invalid initial data"))
		}
		return tool_parameter
	} else if _, ok := data.(*ToolParameter); ok {
		tool_parameter.Name = data.(*ToolParameter).Name
		tool_parameter.Label = commontypes.I18nObject{
			ZhHans: data.(*ToolParameter).Label.ZhHans,
			EnUS:   data.(*ToolParameter).Label.EnUS,
			PtBR:   data.(*ToolParameter).Label.PtBR,
			JaJP:   data.(*ToolParameter).Label.JaJP,
		}
		if data.(*ToolParameter).HumanDescription != nil {
			tool_parameter.HumanDescription = &commontypes.I18nObject{
				ZhHans: data.(*ToolParameter).HumanDescription.ZhHans,
				EnUS:   data.(*ToolParameter).HumanDescription.EnUS,
				PtBR:   data.(*ToolParameter).HumanDescription.PtBR,
				JaJP:   data.(*ToolParameter).HumanDescription.JaJP,
			}
		}
		if data.(*ToolParameter).Placeholder != nil {
			tool_parameter.Placeholder = &commontypes.I18nObject{
				ZhHans: data.(*ToolParameter).Placeholder.ZhHans,
				EnUS:   data.(*ToolParameter).Placeholder.EnUS,
				PtBR:   data.(*ToolParameter).Placeholder.PtBR,
				JaJP:   data.(*ToolParameter).Placeholder.JaJP,
			}
		}
		tool_parameter.Type = data.(*ToolParameter).Type
		tool_parameter.Form = data.(*ToolParameter).Form
		tool_parameter.LLMDescription = data.(*ToolParameter).LLMDescription
		tool_parameter.Required = data.(*ToolParameter).Required
		tool_parameter.Default = data.(*ToolParameter).Default
		tool_parameter.Min = data.(*ToolParameter).Min
		tool_parameter.Max = data.(*ToolParameter).Max
		for _, v := range data.(*ToolParameter).Options {
			if tool_parameter.Options == nil {
				tool_parameter.Options = make([]*ToolParameterOption, 0)
			}
			tool_parameter.Options = append(tool_parameter.Options, NewToolParameterOption(v))
		}
		return tool_parameter
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		return tool_parameter
	}
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
	var param_options []*ToolParameterOption
	for _, v := range options {
		if param_options == nil {
			param_options = make([]*ToolParameterOption, len(options))
		}
		param_options = append(param_options, &ToolParameterOption{
			Value: v,
			Label: commontypes.I18nObject{
				EnUS:   v,
				ZhHans: v,
			},
		})
	}
	return &ToolParameter{
		Name: name,
		Label: commontypes.I18nObject{
			EnUS:   "",
			ZhHans: "",
		},
		HumanDescription: &commontypes.I18nObject{
			EnUS:   "",
			ZhHans: "",
		},
		Placeholder:    nil,
		Type:           param_type,
		Form:           toolsenumtypes.ToolParameterForm_LLM,
		LLMDescription: llm_description,
		Required:       required,
		Options:        param_options, // type: ignore
	}
}

type ToolProviderIdentity struct {
	Author      string/*= Field(..., description="The author of the tool")*/ `json:"author"`
	Name        string/*= Field(..., description="The name of the tool")*/ `json:"name"`
	Description commontypes.I18nObject/*= Field(..., description="The description of the tool")*/ `json:"description"`
	Icon        string/*= Field(..., description="The icon of the tool")*/ `json:"icon"`
	Label       commontypes.I18nObject/*= Field(..., description="The label of the tool")*/ `json:"label"`
	Tags        []toolsenumtypes.ToolLabelType/*= Field(default=[],description="The tags of the tool",)*/ `json:"tags"`
}

type ToolDescription struct {
	Human commontypes.I18nObject/*= Field(..., description="The description presented to the user")*/ `json:"human"`
	LLM   string/*= Field(..., description="The description presented to the LLM")*/ `json:"llm"`
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

type ToolIdentity struct {
	Author   string/*= Field(..., description="The author of the tool")*/ `json:"author"`
	Name     string/*= Field(..., description="The name of the tool")*/ `json:"name"`
	Label    commontypes.I18nObject/*= Field(..., description="The label of the tool")*/ `json:"label"`
	Provider string/*= Field(..., description="The provider of the tool")*/ `json:"provider"`
	Icon     string `json:"icon"`
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

type ToolCredentialsOption struct {
	Value string/*= Field(..., description="The value of the option")*/ `json:"value"`
	Label commontypes.I18nObject/*= Field(..., description="The label of the option")*/ `json:"label"`
}

type ToolProviderCredentials struct {
	Name        string/*= Field(..., description="The name of the credentials")*/ `json:"name"`
	Type        toolsenumtypes.CredentialsType/*= Field(..., description="The type of the credentials")*/ `json:"type"`
	Required    bool                     `json:"required"`
	Default     any                      `json:"default"`
	Options     []*ToolCredentialsOption `json:"options"`
	Label       *commontypes.I18nObject  `json:"label"`
	Help        *commontypes.I18nObject  `json:"help"`
	URL         string                   `json:"url"`
	Placeholder *commontypes.I18nObject  `json:"placeholder"`
}

func (tpc *ToolProviderCredentials) ToDict() map[string]any {
	return map[string]any{
		"name":        tpc.Name,
		"type":        tpc.Type,
		"required":    tpc.Required,
		"default":     tpc.Default,
		"options":     tpc.Options,
		"help":        tpc.Help,
		"label":       tpc.Label,
		"url":         tpc.URL,
		"placeholder": tpc.Placeholder,
	}
}

type ToolRuntimeVariabler interface {
	Type() toolsenumtypes.ToolRuntimeVariableType
	GetName() string
	GetPosition() int
	GetToolName() string
}
type ToolRuntimeVariable struct {
	// Type     toolsenumtypes.ToolRuntimeVariableType/*= Field(..., description="The type of the variable")*/ `json:"type"`
	Name     string/*= Field(..., description="The name of the variable")*/ `json:"name"`
	Position int/*= Field(..., description="The position of the variable")*/ `json:"position"`
	ToolName string/*= Field(..., description="The name of the tool")*/ `json:"tool_name"`
}

func (t *ToolRuntimeVariable) GetName() string {
	return t.Name
}
func (t *ToolRuntimeVariable) GetPosition() int {
	return t.Position
}
func (t *ToolRuntimeVariable) GetToolName() string {
	return t.ToolName
}

type ToolRuntimeTextVariable struct {
	*ToolRuntimeVariable
	Value string/*= Field(..., description="The value of the variable")*/ `json:"value"`
}

func (variable *ToolRuntimeTextVariable) Type() toolsenumtypes.ToolRuntimeVariableType {
	return toolsenumtypes.ToolRuntimeVariable_TEXT
}
func (variable ToolRuntimeTextVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*ToolRuntimeVariable
		Value string `json:"value"`
		Type  string `json:"type"`
	}{
		ToolRuntimeVariable: variable.ToolRuntimeVariable,
		Value:               variable.Value,
		Type:                string((&variable).Type()),
	})
}

type ToolRuntimeImageVariable struct {
	*ToolRuntimeVariable
	Value string/*= Field(..., description="The path of the image")*/ `json:"value"`
}

func (variable *ToolRuntimeImageVariable) Type() toolsenumtypes.ToolRuntimeVariableType {
	return toolsenumtypes.ToolRuntimeVariable_IMAGE
}

func (variable ToolRuntimeImageVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		*ToolRuntimeVariable
		Value string `json:"value"`
		Type  string `json:"type"`
	}{
		ToolRuntimeVariable: variable.ToolRuntimeVariable,
		Value:               variable.Value,
		Type:                string((&variable).Type()),
	})
}

type ToolRuntimeVariablePool struct {
	ConversationID string/*= Field(..., description="The conversation id")*/ `json:"conversation_id"`
	UserID         string/*= Field(..., description="The user id")*/ `json:"user_id"`
	TenantID       string/*= Field(..., description="The tenant id of assistant")*/ `json:"tenant_id"`

	Pool []ToolRuntimeVariabler/*= Field(..., description="The pool of variables")*/ `json:"pool"`
}

func NewToolRuntimeVariablePool(conversation_id, user_id, tenant_id string, pool []ToolRuntimeVariabler) *ToolRuntimeVariablePool {
	return &ToolRuntimeVariablePool{
		ConversationID: conversation_id,
		UserID:         user_id,
		TenantID:       tenant_id,
		Pool:           pool,
	}

}

func (vp *ToolRuntimeVariablePool) SetText(tool_name string, name string, value string) {
	for idx, variable := range vp.Pool {
		if real_variable, ok := any(variable).(*ToolRuntimeTextVariable); ok {
			if real_variable.Name == name {
				real_variable := any(variable).(*ToolRuntimeTextVariable)
				real_variable.Value = value
				vp.Pool[idx] = real_variable
				return
			}
		}
	}
	variable := &ToolRuntimeTextVariable{
		ToolRuntimeVariable: &ToolRuntimeVariable{
			Name:     name,
			Position: len(vp.Pool),
			ToolName: tool_name,
		},
		Value: value,
	}
	if vp.Pool == nil {
		vp.Pool = make([]ToolRuntimeVariabler, 0)
	}
	vp.Pool = append(vp.Pool, variable)
}
func (vp *ToolRuntimeVariablePool) SetFile(tool_name string, value string, name string) {
	// check how many image variables are there
	image_variable_count := 0
	for _, variable := range vp.Pool {
		if variable.Type() == toolsenumtypes.ToolRuntimeVariable_IMAGE {
			image_variable_count += 1
		}
	}
	if name == "" {
		name = fmt.Sprintf("file_%d", image_variable_count)
	}
	for idx, variable := range vp.Pool {
		if real_variable, ok := any(variable).(*ToolRuntimeImageVariable); ok {
			if real_variable.Name == name {
				real_variable := any(variable).(*ToolRuntimeImageVariable)
				real_variable.Value = value
				vp.Pool[idx] = real_variable
				return
			}
		}

	}

	variable := &ToolRuntimeImageVariable{
		ToolRuntimeVariable: &ToolRuntimeVariable{
			Name:     name,
			Position: len(vp.Pool),
			ToolName: tool_name,
		},
		Value: value,
	}

	if vp.Pool == nil {
		vp.Pool = make([]ToolRuntimeVariabler, 0)
	}
	vp.Pool = append(vp.Pool, variable)
}

type ModelToolConfiguration struct {
	// """
	// Model tool configuration
	// """

	Type       string/*= Field(..., description="The type of the model tool")*/ `json:"type"`
	Model      string/*= Field(..., description="The model")*/ `json:"model"`
	Label      commontypes.I18nObject/*= Field(..., description="The label of the model tool")*/ `json:"label"`
	Properties map[toolsenumtypes.ModelToolPropertyKey]any/*= Field(..., description="The properties of the model tool")*/ `json:"properties"`
}

type ModelToolProviderConfiguration struct {
	// """
	// Model tool provider configuration
	// """

	Provider string/*= Field(..., description="The provider of the model tool")*/ `json:"provider"`
	Models   []*ModelToolConfiguration/*= Field(..., description="The models of the model tool")*/ `json:"models"`
	Label    commontypes.I18nObject/*= Field(..., description="The label of the model tool")*/ `json:"label"`
}

type WorkflowToolParameterConfiguration struct {
	// """
	// Workflow tool configuration
	// """

	Name        string/*= Field(..., description="The name of the parameter")*/ `json:"name"`
	Description string/*= Field(..., description="The description of the parameter")*/ `json:"description"`
	Form        toolsenumtypes.ToolParameterForm/*= Field(..., description="The form of the parameter")*/ `json:"form"`
}

type ToolInvokeMeta struct {
	// """
	// Tool invoke meta
	// """

	TimeCost   float64/*= Field(..., description="The time cost of the tool invoke")*/ `json:"time_cost"`
	Error      string         `json:"error"`
	ToolConfig map[string]any `json:"tool_config"`
}

// @classmethod
// def empty(cls) -> "ToolInvokeMeta":
//     """
//     Get an empty instance of ToolInvokeMeta
//     """
//     return cls(time_cost=0.0, error=None, tool_config={})

// @classmethod
// def error_instance(cls, error string) -> "ToolInvokeMeta":
//     """
//     Get an instance of ToolInvokeMeta with error
//     """
//     return cls(time_cost=0.0, error=error, tool_config={})

// def to_dict(self) -> dict:
//     return {
//         "time_cost": self.time_cost,
//         "error": self.error,
//         "tool_config": self.tool_config,
//     }

type ToolLabel struct {
	// """
	// Tool label
	// """

	Name  string/*= Field(..., description="The name of the tool")*/ `json:"name"`
	Label commontypes.I18nObject/*= Field(..., description="The label of the tool")*/ `json:"Label"`
	Icon  string/*= Field(..., description="The icon of the tool")*/ `json:"icon"`
}
