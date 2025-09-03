package variables

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

type BasicVariablesConfigManager struct {
}

func (mgr *BasicVariablesConfigManager) Convert(config map[string]any) ([]*appconfigentities.VariableEntity, []*appconfigentities.ExternalDataVariableEntity) {
	external_data_variables := []*appconfigentities.ExternalDataVariableEntity{}
	variable_entities := []*appconfigentities.VariableEntity{}

	// old external_data_tools
	external_data_tools := mapstruct.Get(config, "external_data_tools", []map[string]any{})
	for _, external_data_tool := range external_data_tools {
		if _, ok := external_data_tool["enabled"]; !ok {
			continue
		}
		if _, ok := external_data_tool["enabled"].(bool); !ok || !external_data_tool["enabled"].(bool) {
			continue
		}

		external_data_variables = append(external_data_variables, &appconfigentities.ExternalDataVariableEntity{
			Variable: mapstruct.Get(external_data_tool, "variable", ""),
			Type:     mapstruct.Get(external_data_tool, "type", ""),
			Config:   mapstruct.Get(external_data_tool, "config", map[string]any{}),
		})
	}
	// variables and external_data_tools
	for _, variables := range mapstruct.Get(config, "user_input_form", []map[string]any{}) {
		variable_keys := slices.Sorted(maps.Keys(variables))
		if len(variable_keys) == 0 {
			continue
		}
		variable_type := variable_keys[0]
		if variable_type == string(appconfigenumtypes.VariableEntity_EXTERNAL_DATA_TOOL) {
			if _, ok := variables[variable_type]; !ok {
				continue
			}
			if _, ok := variables[variable_type].(map[string]any); !ok {
				continue
			}
			variable := variables[variable_type].(map[string]any)
			if _, ok := variable["config"]; !ok {
				continue
			}

			external_data_variables = append(external_data_variables, &appconfigentities.ExternalDataVariableEntity{
				Variable: mapstruct.Get(variable, "variable", ""), Type: mapstruct.Get(variable, "type", ""), Config: mapstruct.Get(variable, "config", map[string]any{}),
			})
		} else if slices.Contains([]appconfigenumtypes.VariableEntityType{
			appconfigenumtypes.VariableEntity_TEXT_INPUT,
			appconfigenumtypes.VariableEntity_PARAGRAPH,
			appconfigenumtypes.VariableEntity_NUMBER,
			appconfigenumtypes.VariableEntity_SELECT,
		}, appconfigenumtypes.VariableEntityType(variable_type)) {
			if _, ok := variables[variable_type]; !ok {
				continue
			}
			if _, ok := variables[variable_type].(map[string]any); !ok {
				continue
			}
			variable := variables[variable_type].(map[string]any)
			variable_entities = append(variable_entities, &appconfigentities.VariableEntity{
				Type:        appconfigenumtypes.VariableEntityType(variable_type),
				Variable:    mapstruct.Get(variable, "variable", ""),
				Description: mapstruct.Get(variable, "description", ""),
				Label:       mapstruct.Get(variable, "label", ""),
				Required:    mapstruct.Get(variable, "required", false),
				MaxLength:   mapstruct.Get(variable, "max_length", 0),
				Options:     mapstruct.Get(variable, "options", []string{}),
			})
		}
	}
	return variable_entities, external_data_variables

}
func (mgr *BasicVariablesConfigManager) ValidateAndSetDefaults(tenant_id string, config map[string]any) (map[string]any, []string) {
	related_config_keys := []string{}
	config, current_related_config_keys := mgr.ValidateVariablesAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	config, current_related_config_keys = mgr.ValidateExternalDataToolsAndSetDefaults(tenant_id, config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	return config, related_config_keys

}
func (mgr *BasicVariablesConfigManager) ValidateVariablesAndSetDefaults(config map[string]any) (map[string]any, []string) {
	if _, ok := config["user_input_form"]; !ok {
		config["user_input_form"] = []map[string]any{}
	}
	if _, ok := config["user_input_form"].([]any); !ok {
		if _, ok := config["user_input_form"].([]map[string]any); !ok {
			panic(exceptions.NewValueError("user_input_form must be a list of objects"))
		}
	} else {
		user_input_form := []map[string]any{}
		for _, v := range config["user_input_form"].([]any) {
			if _, ok := v.(map[string]any); !ok {
				mlog.Errorf("user_input_form=%#v must be a list of objects", config["user_input_form"])
				panic(exceptions.NewValueError("user_input_form must be a list of objects"))
			} else {
				user_input_form = append(user_input_form, v.(map[string]any))
			}
		}
		config["user_input_form"] = user_input_form
	}
	re := regexp.MustCompile(`^(?![0-9])[\p{Han}A-Za-z0-9_\U0001F300-\U0001F64F\U0001F680-\U0001F6FF]{1,100}$`)
	user_input_form := mapstruct.Get(config, "user_input_form", []map[string]any{})
	variables := []string{}
	for _, item := range user_input_form {
		for key := range item {
			if !slices.Contains([]string{"text-input", "select", "paragraph", "number", "external_data_tool"}, key) {
				panic(exceptions.NewValueError("Keys in user_input_form list can only be 'text-input', 'paragraph'  or 'select'"))
			}
			if _, ok := item[key].(map[string]any); !ok {
				panic(exceptions.NewValueError(fmt.Sprintf("user_input_form[%s]=%#v must be object", key, item[key])))
			}
			form_item := item[key].(map[string]any)
			if _, ok := form_item["label"]; !ok {
				panic(exceptions.NewValueError("label is required in user_input_form"))
			}
			if _, ok := form_item["label"].(string); !ok {
				panic(exceptions.NewValueError("label in user_input_form must be of string type"))
			}
			if _, ok := form_item["variable"]; !ok {
				panic(exceptions.NewValueError("variable is required in user_input_form"))
			}
			if _, ok := form_item["variable"].(string); !ok {
				panic(exceptions.NewValueError("variable in user_input_form must be of string type"))
			}

			if !re.MatchString(form_item["variable"].(string)) {
				panic(errors.New("variable in user_input_form must be a string, and cannot start with a number"))
			}
			variables = append(variables, form_item["variable"].(string))
			if _, ok := form_item["required"]; !ok || form_item["required"] == nil {
				form_item["required"] = false
			}
			if _, ok := form_item["required"].(bool); !ok {
				panic(exceptions.NewValueError("required in user_input_form must be of boolean type"))
			}
			if key == "select" {
				if _, ok := form_item["options"]; !ok || form_item["options"] == nil {
					form_item["options"] = []any{}
				}
				if _, ok := form_item["options"].([]any); !ok {
					panic(exceptions.NewValueError("options in user_input_form must be a list of strings"))
				}
				if _, ok := form_item["default"]; ok && form_item["options"] != nil && !slices.Contains(form_item["options"].([]any), form_item["options"]) {
					panic(exceptions.NewValueError("default value in user_input_form must be in the options list"))
				}
			}
		}
	}
	return config, []string{"user_input_form"}

}
func (mgr *BasicVariablesConfigManager) ValidateExternalDataToolsAndSetDefaults(tenant_id string, config map[string]any) (map[string]any, []string) {
	if _, ok := config["external_data_tools"]; !ok {
		config["external_data_tools"] = []map[string]any{}
	}
	if _, ok := config["external_data_tools"].([]any); !ok {
		if _, ok := config["external_data_tools"].([]map[string]any); !ok {
			panic(exceptions.NewValueError("external_data_tools must be a list of objects"))
		}
	} else {
		external_data_tools := []map[string]any{}
		for _, v := range config["external_data_tools"].([]any) {
			if _, ok := v.(map[string]any); !ok {
				mlog.Errorf("external_data_tools=%#v must be a list of objects", config["external_data_tools"])
				panic(exceptions.NewValueError("external_data_tools must be a list of objects"))
			} else {
				external_data_tools = append(external_data_tools, v.(map[string]any))
			}
		}
		config["external_data_tools"] = external_data_tools
	}
	external_data_tools := mapstruct.Get(config, "external_data_tools", []map[string]any{})
	for _, tool := range external_data_tools {
		if _, ok := tool["enabled"]; !ok || tool["enabled"] == nil {
			tool["enabled"] = false
		}
		if _, ok := tool["enabled"].(bool); !ok {
			panic(exceptions.NewValueError("enabled must be bool"))
		}

		if !tool["enabled"].(bool) {
			continue
		}
		if _, ok := tool["type"]; !ok || tool["type"] == nil {
			panic(exceptions.NewValueError("external_data_tools[].type is required"))
		}
		if _, ok := tool["type"].(string); !ok {
			panic(exceptions.NewValueError("external_data_tools[].type must be string"))
		}
		// typ = tool["type"]
		// config = tool["config"]

		// ExternalDataToolFactory.validate_config(name=typ, tenant_id=tenant_id, config=config)
	}
	return config, []string{"external_data_tools"}
}
