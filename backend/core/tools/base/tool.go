package base

import (
	"iter"
	"maps"

	"github.com/odysseythink/gofy/backend/core/file"
	toolsentities "github.com/odysseythink/gofy/backend/entities/tools"
	toolsenumtypes "github.com/odysseythink/gofy/backend/enum_types/tools"
)

type ToolInvokeResponseType interface {
	*toolsentities.ToolInvokeMessage | []*toolsentities.ToolInvokeMessage
}

type Toolor interface {
	ToolProviderType() toolsenumtypes.ToolProviderType
	Invoke(
		user_id string,
		tool_parameters map[string]any,
		conversation_id string,
		app_id string,
		message_id string,
	) (*toolsentities.ToolInvokeMessage, []*toolsentities.ToolInvokeMessage, iter.Seq[*toolsentities.ToolInvokeMessage])
	ForkToolRuntime(runtime *ToolRuntime) *Toolor
}

type Tool struct {
	Entity  *toolsentities.ToolEntity `json:"entity"`
	Runtime *ToolRuntime              `json:"runtime"`
}

func NewTool(entity *toolsentities.ToolEntity, runtime *ToolRuntime) *Tool {
	return &Tool{
		Entity:  entity,
		Runtime: runtime,
	}
}

func (t *Tool) ForkToolRuntime(runtime *ToolRuntime) *Tool {
	/*
	   fork a new tool with meta data

	   :param meta: the meta data of a tool call processing, tenant_id is required
	   :return: the new tool
	*/
	new_tool := &Tool{
		Entity:  toolsentities.NewToolEntity(t.Entity),
		Runtime: NewToolRuntime(t.Runtime),
	}

	return new_tool
}

func (t *Tool) Invoke1(
	instance Toolor,
	user_id string,
	tool_parameters map[string]any,
	conversation_id string,
	app_id string,
	message_id string,
) iter.Seq[*toolsentities.ToolInvokeMessage] {
	return func(yield func(*toolsentities.ToolInvokeMessage) bool) {
		if t.Runtime != nil && t.Runtime.RuntimeParameters != nil {
			maps.Copy(tool_parameters, t.Runtime.RuntimeParameters)
		}
		// try parse tool parameters into the correct type
		tool_parameters = t._transform_tool_parameters_type(tool_parameters)

		result1, result2, result3 := instance.Invoke(
			user_id,
			tool_parameters,
			conversation_id,
			app_id,
			message_id,
		)

		if result1 != nil {
			yield(result1)
		} else if result2 != nil {
			for _, v := range result2 {
				yield(v)
			}
		} else {
			for v := range result3 {
				yield(v)
			}
		}
	}
}

func (t *Tool) _transform_tool_parameters_type(tool_parameters map[string]any) map[string]any {
	/*
	   Transform tool parameters type
	*/
	// Temp fix for the issue that the tool parameters will be converted to empty while validating the credentials
	result := maps.Clone(tool_parameters)
	for _, parameter := range t.Entity.Parameters {
		if _, ok := tool_parameters[parameter.Name]; ok {
			result[parameter.Name] = parameter.Type.CastValue(tool_parameters[parameter.Name])
		}
	}
	return result
}

func (t *Tool) GetRuntimeParameters() []*toolsentities.ToolParameter {
	/*
	   get the runtime parameters

	   interface for developer to dynamic change the parameters of a tool depends on the variables pool

	   :return: the runtime parameters
	*/
	return t.Entity.Parameters
}
func (t *Tool) GetMergedRuntimeParameters(
	conversation_id string,
	app_id string,
	message_id string,
) []*toolsentities.ToolParameter {
	/*
	   get all runtime parameters

	   :return: all runtime parameters
	*/
	parameters := t.Entity.Parameters
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
		Message: &toolsentities.TextMessage{Text: image},
	}
}
func (t *Tool) CreateFileMessage(file *file.File) *toolsentities.ToolInvokeMessage {
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_FILE,
		Message: &toolsentities.FileMessage{},
		Meta:    map[string]any{"file": file},
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
		Message: &toolsentities.TextMessage{Text: link},
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
		Message: &toolsentities.TextMessage{Text: text},
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
		Message: &toolsentities.BlobMessage{Blob: blob},
		Meta:    meta,
	}
}
func (t *Tool) CreateJsonMessage(object map[string]any) *toolsentities.ToolInvokeMessage {
	/*
	   create a json message
	*/
	return &toolsentities.ToolInvokeMessage{
		Type:    toolsenumtypes.Message_JSON,
		Message: &toolsentities.JsonMessage{JsonObject: object},
	}
}
