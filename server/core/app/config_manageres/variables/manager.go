package variables

import (
	"maps"
	"slices"

	agententities "mlib.com/gofy/server/entities/agent"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	promptentities "mlib.com/gofy/server/entities/prompt"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	"mlib.com/gofy/server/utils/mapstruct"
)

type BasicVariablesConfigManager struct{

}

func (mgr *BasicVariablesConfigManager) convert(config map[string]any) ([]*appconfigentities.VariableEntity, []*appconfigentities.ExternalDataVariableEntity){
	external_data_variables := []*appconfigentities.ExternalDataVariableEntity{}
	variable_entities := []*appconfigentities.VariableEntity{}

	// old external_data_tools
	external_data_tools := mapstruct.Get(config,"external_data_tools", []map[string]any{})
	for _, external_data_tool := range external_data_tools{
		if _, ok := external_data_tool["enabled"]; !ok {
			continue
		}
		if _, ok := external_data_tool["enabled"].(bool); !ok || !external_data_tool["enabled"].(bool){
			continue
		}

		external_data_variables=append(external_data_variables,&appconfigentities.ExternalDataVariableEntity{
				Variable: mapstruct.Get(external_data_tool,"variable",""),
				Type:mapstruct.Get(external_data_tool,"type",""),
				Config:mapstruct.Get(external_data_tool,"config",""),
			})
	}
	// variables and external_data_tools
	for _, variables := range mapstruct.Get(config,"user_input_form", []map[string]any{}){
		variable_keys := slices.Sort(maps.Keys(variables))
		if len(variable_keys) == 0{
			continue
		}
		variable_type := variable_keys[0]
		if variable_type == string(appconfigenumtypes.VariableEntity_EXTERNAL_DATA_TOOL){
			variable = variables[variable_type]
			if "config" not in variable{
				continue
			}
			external_data_variables.append(
				ExternalDataVariableEntity(
					variable=variable["variable"], type=variable["type"], config=variable["config"]
				)
			)
		} else if variable_type in {
			VariableEntityType.TEXT_INPUT,
			VariableEntityType.PARAGRAPH,
			VariableEntityType.NUMBER,
			VariableEntityType.SELECT,
		}{
			variable = variables[variable_type]
			variable_entities.append(
				VariableEntity(
					type=variable_type,
					variable=variable.get("variable"),
					description=variable.get("description") or "",
					label=variable.get("label"),
					required=variable.get("required", False),
					max_length=variable.get("max_length"),
					options=variable.get("options") or [],
				)
			)
		}
	}
	return variable_entities, external_data_variables

}
func (mgr *BasicVariablesConfigManager) validate_and_set_defaults(tenant_id  string, config map[string]any) (dict, list[str]]{
	related_config_keys = []
	config, current_related_config_keys = cls.validate_variables_and_set_defaults(config)
	related_config_keys.extend(current_related_config_keys)

	config, current_related_config_keys = cls.validate_external_data_tools_and_set_defaults(tenant_id, config)
	related_config_keys.extend(current_related_config_keys)

	return config, related_config_keys

}
func (mgr *BasicVariablesConfigManager) validate_variables_and_set_defaults(config map[string]any) (dict, list[str]]{
	if not config.get("user_input_form"){
		config["user_input_form"] = []
	}
	if not isinstance(config["user_input_form"], list){
		raise ValueError("user_input_form must be a list of objects")
	}
	variables = []
	for item in config["user_input_form"]{
		key = list(item.keys())[0]
		if key not in {"text-input", "select", "paragraph", "number", "external_data_tool"}{
			raise ValueError("Keys in user_input_form list can only be 'text-input', 'paragraph'  or 'select'")
		}
		form_item = item[key]
		if "label" not in form_item{
			raise ValueError("label is required in user_input_form")
		}
		if not isinstance(form_item["label"], str){
			raise ValueError("label in user_input_form must be of string type")
		}
		if "variable" not in form_item{
			raise ValueError("variable is required in user_input_form")
		}
		if not isinstance(form_item["variable"], str){
			raise ValueError("variable in user_input_form must be of string type")
		}
		pattern = re.compile(r"^(?!\d)[\u4e00-\u9fa5A-Za-z0-9_\U0001F300-\U0001F64F\U0001F680-\U0001F6FF]{1,100}$")
		if pattern.match(form_item["variable"]) is None{
			raise ValueError("variable in user_input_form must be a string, and cannot start with a number")
		}
		variables.append(form_item["variable"])

		if "required" not in form_item or not form_item["required"]{
			form_item["required"] = False
		}
		if not isinstance(form_item["required"], bool){
			raise ValueError("required in user_input_form must be of boolean type")
		}
		if key == "select"{
			if "options" not in form_item or not form_item["options"]{
				form_item["options"] = []
			}
			if not isinstance(form_item["options"], list){
				raise ValueError("options in user_input_form must be a list of strings")
			}
			if "default" in form_item and form_item["default"] and form_item["default"] not in form_item["options"]{
				raise ValueError("default value in user_input_form must be in the options list")
			}
		}
	}
	return config, ["user_input_form"]

}
func (mgr *BasicVariablesConfigManager) validate_external_data_tools_and_set_defaults(tenant_id  string, config map[string]any) (dict, list[str]]{
	if not config.get("external_data_tools"){
		config["external_data_tools"] = []
	}
	if not isinstance(config["external_data_tools"], list){
		raise ValueError("external_data_tools must be of list type")
	}
	for tool in config["external_data_tools"]{
		if "enabled" not in tool or not tool["enabled"]{
			tool["enabled"] = False
		}
		if not tool["enabled"]{
			continue
		}
		if "type" not in tool or not tool["type"]{
			raise ValueError("external_data_tools[].type is required")
		}
		typ = tool["type"]
		config = tool["config"]

		ExternalDataToolFactory.validate_config(name=typ, tenant_id=tenant_id, config=config)
	}
	return config, ["external_data_tools"]
}