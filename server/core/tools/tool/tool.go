package tool

import (
	"encoding/json"
	"maps"
	"strings"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	toolfilemanager "mlib.com/gofy/server/core/tools/tool_file_manager"
	toolsentities "mlib.com/gofy/server/entities/tools"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/mlog"
)

type ToolInvokeResponseType interface {
	*toolsentities.ToolInvokeMessage | []*toolsentities.ToolInvokeMessage
}

type Tool struct {
	Identity            *toolsentities.ToolIdentity    `json:"identity"`
	Parameters          []*toolsentities.ToolParameter `json:"parameters"`
	Description         *toolsentities.ToolDescription `json:"description"`
	IsTeamAuthorization bool                           `json:"is_team_authorization"`

	// pydantic configs
	ModelConfig map[string]any                         `json:"model_config"`
	Runtime     *ToolRuntime                           `json:"runtime"`
	Variables   *toolsentities.ToolRuntimeVariablePool `json:"variables"`
}

type ToolRuntime struct {
	TenantID          string                        `json:"tenant_id"`
	ToolID            string                        `json:"tool_id"`
	InvokeFrom        appenumtypes.InvokeFrom       `json:"invoke_from"`
	ToolInvokeFrom    toolsenumtypes.ToolInvokeFrom `json:"tool_invoke_from"`
	Credentials       map[string]any                `json:"credentials"`
	RuntimeParameters map[string]any                `json:"runtime_parameters"`
}

func NewToolRuntime(data any) *ToolRuntime {
	tool_runtime := new(ToolRuntime)
	if data == nil {
		return tool_runtime
	}
	if _, ok := data.(map[string]any); ok {
		bindata, _ := json.Marshal(data.(map[string]any))
		err := json.Unmarshal(bindata, tool_runtime)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
			panic(exceptions.NewValueError("new ToolRuntime failed, because of provaded invalid initial data"))
		}
		return tool_runtime
	} else if _, ok := data.(*ToolRuntime); ok {
		tool_runtime.TenantID = data.(*ToolRuntime).TenantID
		tool_runtime.ToolID = data.(*ToolRuntime).ToolID
		tool_runtime.InvokeFrom = data.(*ToolRuntime).InvokeFrom
		tool_runtime.ToolInvokeFrom = data.(*ToolRuntime).ToolInvokeFrom
		tool_runtime.Credentials = maps.Clone(data.(*ToolRuntime).Credentials)
		tool_runtime.RuntimeParameters = maps.Clone(data.(*ToolRuntime).RuntimeParameters)
		return tool_runtime
	} else {
		mlog.Warningf("provaded initial data=%#v is unsupported", data)
		return tool_runtime
	}
}

func (t *Tool) ForkToolRuntime(runtime map[string]any) *Tool {
	/*
	   fork a new tool with meta data

	   :param meta: the meta data of a tool call processing, tenant_id is required
	   :return: the new tool
	*/
	new_tool := &Tool{
		Identity:            toolsentities.NewToolIdentity(t.Identity),
		Description:         toolsentities.NewToolDescription(t.Description),
		IsTeamAuthorization: t.IsTeamAuthorization,
		Runtime:             NewToolRuntime(t.Runtime),
	}
	for _, v := range t.Parameters {
		if new_tool.Parameters == nil {
			new_tool.Parameters = make([]*toolsentities.ToolParameter, 0)
		}
		new_tool.Parameters = append(new_tool.Parameters, v.Copy())
	}

	return new_tool
}

func (t *Tool) LoadVariables(variables *toolsentities.ToolRuntimeVariablePool) {
	/*
	   load variables from database

	   :param conversation_id: the conversation id
	*/
	t.Variables = variables
}
func (t *Tool) SetImageVariable(variable_name string, image_key string) {
	/*
	   set an image variable
	*/
	if t.Variables == nil {
		return
	}
	if t.Identity == nil {
		return
	}
	t.Variables.SetFile(t.Identity.Name, variable_name, image_key)
}
func (t *Tool) SetTextVariable(variable_name string, text string) {
	/*
	   set a text variable
	*/
	if t.Variables == nil {
		return
	}
	if t.Identity == nil {
		return
	}
	t.Variables.SetText(t.Identity.Name, variable_name, text)
}
func (t *Tool) GetVariable(name string) toolsentities.ToolRuntimeVariabler {
	/*
	   get a variable

	   :param name: the name of the variable
	   :return: the variable
	*/
	if t.Variables == nil {
		return nil
	}

	for _, variable := range t.Variables.Pool {
		if variable.GetName() == name {
			return variable
		}
	}

	return nil
}
func (t *Tool) GetDefaultImageVariable() *toolsentities.ToolRuntimeImageVariable {
	/*
	   get the default image variable

	   :return: the image variable
	*/
	if t.Variables == nil {
		return nil
	}

	v := t.GetVariable(string(toolsenumtypes.ToolRuntimeVariable_IMAGE))
	if _, ok := any(v).(*toolsentities.ToolRuntimeImageVariable); ok {
		return any(v).(*toolsentities.ToolRuntimeImageVariable)
	} else {
		return nil
	}
}
func (t *Tool) GetVariableFile(name string) []byte {
	/*
	   get a variable file

	   :param name: the name of the variable
	   :return: the variable file
	*/
	variable := t.GetVariable(name)
	if variable == nil {
		return nil
	}
	if _, ok := any(variable).(*toolsentities.ToolRuntimeImageVariable); !ok {
		return nil
	}
	real_variable := any(variable).(*toolsentities.ToolRuntimeImageVariable)

	message_file_id := real_variable.Value
	// get file binary
	file_binary, _ := (&toolfilemanager.ToolFileManager{}).GetFileBinaryByMessageFileID(message_file_id)
	if file_binary == nil {
		return nil
	}
	return file_binary
}
func (t *Tool) list_variables() []toolsentities.ToolRuntimeVariabler {
	/*
	   list all variables

	   :return: the variables
	*/
	if t.Variables == nil {
		return nil
	}

	return t.Variables.Pool
}
func (t *Tool) ListDefaultImageVariables() []toolsentities.ToolRuntimeVariabler {
	/*
	   list all image variables

	   :return: the image variables
	*/
	if t.Variables == nil {
		return nil
	}

	result := []toolsentities.ToolRuntimeVariabler{}

	for _, variable := range t.Variables.Pool {
		if strings.HasPrefix(variable.GetName(), string(toolsenumtypes.ToolRuntimeVariable_IMAGE)) {
			result = append(result, variable)
		}
	}
	return result
}
func (t *Tool) _transform_tool_parameters_type(tool_parameters map[string]any) map[string]any {
	/*
	   Transform tool parameters type
	*/
	// Temp fix for the issue that the tool parameters will be converted to empty while validating the credentials
	result := maps.Clone(tool_parameters)
	for _, parameter := range t.Parameters {
		if _, ok := tool_parameters[parameter.Name]; ok {
			result[parameter.Name] = parameter.Type.CastValue(tool_parameters[parameter.Name])
		}
	}
	return result
}
func (t *Tool) Invoke(user_id string, tool_parameters map[string]any, tooler Tooler) []*toolsentities.ToolInvokeMessage {
	// update tool_parameters
	// TODO: Fix type error.
	if t.Runtime == nil {
		return nil
	}
	if len(t.Runtime.RuntimeParameters) > 0 {
		// Convert Mapping to dict before updating
		maps.Copy(tool_parameters, t.Runtime.RuntimeParameters)
	}
	// try parse tool parameters into the correct type
	tool_parameters = t._transform_tool_parameters_type(tool_parameters)
	rsp := []*toolsentities.ToolInvokeMessage{}
	result, results := tooler.RealInvoke(user_id, tool_parameters)

	if results == nil {
		rsp = append(rsp, result)
	}

	return rsp
}

func (t *Tool) GetRuntimeParameters() []*toolsentities.ToolParameter {
	/*
	   get the runtime parameters

	   interface for developer to dynamic change the parameters of a tool depends on the variables pool

	   :return: the runtime parameters
	*/
	return t.Parameters
}
func (t *Tool) GetAllRuntimeParameters() []*toolsentities.ToolParameter {
	/*
	   get all runtime parameters

	   :return: all runtime parameters
	*/
	parameters := t.Parameters
	user_parameters := t.GetRuntimeParameters()

	// override parameters
	for _, parameter := range user_parameters {
		// check if parameter in tool parameters
		found := false
		var matched_tool_parameter *toolsentities.ToolParameter
		for _, tool_parameter := range parameters {
			if tool_parameter.Name == parameter.Name {
				found = true
				matched_tool_parameter = tool_parameter
				break
			}
		}
		if found {
			// override parameter
			matched_tool_parameter.Type = parameter.Type
			matched_tool_parameter.Form = parameter.Form
			matched_tool_parameter.Required = parameter.Required
			matched_tool_parameter.Default = parameter.Default
			matched_tool_parameter.Options = parameter.Options
			matched_tool_parameter.LLMDescription = parameter.LLMDescription
		} else {
			// add new parameter
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}
func (t *Tool) CreateImageMessage(image string, save_as string) *toolsentities.ToolInvokeMessage {
	/*
	   create an image message

	   :param image: the url of the image
	   :return: the image message
	*/
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_IMAGE,
		Message: image,
		SaveAs:  save_as,
	}
}
func (t *Tool) CreateFileMessage(file *file.File) *toolsentities.ToolInvokeMessage {
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_FILE,
		Message: "",
		Meta:    map[string]any{"file": file},
		SaveAs:  "",
	}
}
func (t *Tool) CreateLinkMessage(link string, save_as string) *toolsentities.ToolInvokeMessage {
	/*
	   create a link message

	   :param link: the url of the link
	   :return: the link message
	*/
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_LINK,
		Message: link,
		SaveAs:  save_as,
	}
}
func (t *Tool) CreateTextMessage(text string, save_as string) *toolsentities.ToolInvokeMessage {
	/*
	   create a text message

	   :param text: the text
	   :return: the text message
	*/
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_TEXT,
		Message: text,
		SaveAs:  save_as,
	}
}
func (t *Tool) CreateBlobMessage(blob []byte, meta map[string]any, save_as string) *toolsentities.ToolInvokeMessage {
	/*
	   create a blob message

	   :param blob: the blob
	   :return: the blob message
	*/
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_BLOB,
		Message: blob,
		Meta:    meta,
		SaveAs:  save_as,
	}
}
func (t *Tool) CreateJsonMessage(object map[string]any) *toolsentities.ToolInvokeMessage {
	/*
	   create a json message
	*/
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_JSON,
		Message: object,
	}
}
