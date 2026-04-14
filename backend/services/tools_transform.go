package services

import (
	"slices"

	"github.com/odysseythink/mlog"
	pluginparameterentities "mlib.com/gofy/server/entities/plugin/parameter"
	toolsentities "mlib.com/gofy/server/entities/tools"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/types/common"
	mcptypes "mlib.com/gofy/server/types/mcp"
)

type ToolsTransformService struct {
}

func create_parameter(
	name string, description string, param_type string, required bool, input_schema map[string]any,
) *toolsentities.ToolParameter {
	/*Create a ToolParameter instance with given attributes*/
	input_schema_dict := map[string]any{}
	if input_schema != nil {
		input_schema_dict["input_schema"] = input_schema
	}
	return &toolsentities.ToolParameter{
		PluginParameter: &pluginparameterentities.PluginParameter{
			Name:     name,
			Label:    common.I18nObject{EnUS: name},
			Required: required,
		},
		Type:             toolsenumtypes.ToolParameterType(param_type),
		HumanDescription: &common.I18nObject{EnUS: description},
		Form:             toolsenumtypes.ToolParameterForm_LLM,
		LLMDescription:   description,
		InputSchema:      input_schema,
	}

}

var (
	TYPE_MAPPING  = map[string]string{"integer": "number", "float": "number"}
	COMPLEX_TYPES = []string{"array", "object"}
)

func process_properties(props map[string]map[string]any, required []string, prefix string) []*toolsentities.ToolParameter {
	/*Process properties recursively*/

	parameters := []*toolsentities.ToolParameter{}
	for name, prop := range props {
		current_description := ""
		if _, ok := prop["description"]; ok {
			if _, ok := prop["description"].(string); ok {
				current_description = prop["description"].(string)
			}
		}
		prop_type := "string"
		if _, ok := prop["type"]; ok {
			if _, ok := prop["type"].(string); ok {
				prop_type = prop["type"].(string)
			} else if _, ok := prop["type"].([]any); ok {
				if _, ok := prop["type"].([]any)[0].(string); ok {
					prop_type = prop["type"].([]any)[0].(string)
				}
			} else if _, ok := prop["type"].([]string); ok {
				prop_type = prop["type"].([]string)[0]
			}
		}

		if _, ok := TYPE_MAPPING[prop_type]; ok {
			prop_type = TYPE_MAPPING[prop_type]
		}
		var input_schema map[string]any
		if slices.Contains(COMPLEX_TYPES, prop_type) {
			input_schema = prop
		}

		parameters = append(parameters, create_parameter(name, current_description, prop_type, slices.Contains(required, name), input_schema))
	}
	return parameters
}
func (s *ToolsTransformService) ConvertMCPSchemaToParameter(schema map[string]any) []*toolsentities.ToolParameter {
	/*
	   Convert MCP JSON schema to tool parameters

	   :param schema: JSON schema dictionary
	   :return: list of ToolParameter instances
	*/
	typ := ""
	if _, ok := schema["type"]; ok {
		if _, ok := schema["type"].(string); ok {
			typ = schema["type"].(string)
		}
	}
	if _, ok := schema["properties"]; ok && typ == "object" {
		required := []string{}
		if _, ok := schema["required"]; ok {
			if _, ok := schema["required"].([]string); ok {
				required = schema["required"].([]string)
			} else if _, ok := schema["required"].([]any); ok {
				for _, val := range schema["required"].([]any) {
					if _, ok := val.(string); !ok {
						mlog.Errorf("schema[\"required\"]=%#v must be string list", schema["required"])
						return nil
					}
					required = append(required, val.(string))
				}
			}
		}
		properties := map[string]map[string]any{}
		if _, ok := schema["properties"].(map[string]map[string]any); ok {
			properties = schema["properties"].(map[string]map[string]any)
		} else if _, ok := schema["properties"].(map[string]any); ok {
			for k, v := range schema["properties"].(map[string]any) {
				if _, ok := v.(map[string]any); !ok {
					mlog.Errorf("schema[\"required\"]=%#v must be map[string]any", schema["required"])
					return nil
				}
				properties[k] = v.(map[string]any)
			}
		}
		return process_properties(properties, required, "")
	}
	return nil
}

func (s *ToolsTransformService) MCPToolToUserTool(mcp_provider *models.MCPToolProvider, tools []*mcptypes.Tool) []*toolsentities.ToolApiEntity {
	user := mcp_provider.LoadUser()
	var datas []*toolsentities.ToolApiEntity
	for _, tool := range tools {
		if datas == nil {
			datas = make([]*toolsentities.ToolApiEntity, 0)
		}
		entity := &toolsentities.ToolApiEntity{
			Name:        tool.Name,
			Label:       common.I18nObject{EnUS: tool.Name, ZhHans: tool.Name},
			Description: common.I18nObject{EnUS: tool.Description, ZhHans: tool.Description},
			Labels:      []string{},
			Parameters:  s.ConvertMCPSchemaToParameter(tool.InputSchema),
		}
		if user != nil {
			entity.Author = user.Name
		} else {
			entity.Author = "Anonymous"
		}

		datas = append(datas, entity)
	}
	return datas
}
func (s *ToolsTransformService) MCPProviderToUserProvider(db_provider *models.MCPToolProvider, for_list bool) *toolsentities.ToolProviderApiEntity {
	user := db_provider.LoadUser()
	entity := &toolsentities.ToolProviderApiEntity{
		Name:                db_provider.Name,
		Description:         common.I18nObject{EnUS: "", ZhHans: ""},
		Icon:                db_provider.Icon,
		Label:               common.I18nObject{EnUS: db_provider.Name, ZhHans: db_provider.Name},
		Type:                toolsenumtypes.ToolProvider_MCP,
		IsTeamAuthorization: db_provider.Authed,
		Tools:               s.MCPToolToUserTool(db_provider, db_provider.MCPTools()),
		ServerURL:           db_provider.MaskedServerURL(),
		ServerIdentifier:    db_provider.ServerIdentifier,
	}
	if !for_list {
		entity.ID = db_provider.ServerIdentifier
	} else {
		entity.ID = db_provider.ID
	}
	if user != nil {
		entity.Author = user.Name
	} else {
		entity.Author = "Anonymous"
	}

	if db_provider.UpdatedAt != nil {
		entity.UpdatedAt = db_provider.UpdatedAt.Unix()
	}

	return entity
}
