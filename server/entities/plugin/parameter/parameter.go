package parameter

import (
	"fmt"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
	commonparameter "mlib.com/gofy/server/entities/parameter"
	commontypes "mlib.com/gofy/server/types/common"
)


type PluginParameterOption struct{
    Value string `json:"value"` //description="The value of the option")
    Label commontypes.I18nObject `json:"label"` //description="The label of the option")
    Icon string `json:"icon"` //description="The icon of the option, can be a url or a base64 encoded image"
}
func(option *PluginParameterOption)TransformIDToStr(value any) string{
    switch real_val := value.(type){
    case string:
        return real_val
    default:
        fmt.Sprintf("%v", real_val)
    }
}

type PluginParameterType string
const (
    PluginParameter_STRING = PluginParameterType(commonparameter.CommonParameter_STRING)
    PluginParameter_NUMBER = PluginParameterType(commonparameter.CommonParameter_NUMBER)
    PluginParameter_BOOLEAN = PluginParameterType(commonparameter.CommonParameter_BOOLEAN)
    PluginParameter_SELECT = PluginParameterType(commonparameter.CommonParameter_SELECT)
    PluginParameter_SECRET_INPUT = PluginParameterType(commonparameter.CommonParameter_SECRET_INPUT)
    PluginParameter_FILE = PluginParameterType(commonparameter.CommonParameter_FILE)
    PluginParameter_FILES = PluginParameterType(commonparameter.CommonParameter_FILES)
    PluginParameter_APP_SELECTOR = PluginParameterType(commonparameter.CommonParameter_APP_SELECTOR)
    PluginParameter_MODEL_SELECTOR = PluginParameterType(commonparameter.CommonParameter_MODEL_SELECTOR)
    PluginParameter_TOOLS_SELECTOR = PluginParameterType(commonparameter.CommonParameter_TOOLS_SELECTOR)
    PluginParameter_ANY = PluginParameterType(commonparameter.CommonParameter_ANY)
    PluginParameter_DYNAMIC_SELECT = PluginParameterType(commonparameter.CommonParameter_DYNAMIC_SELECT)

    // deprecated, should not use.
    PluginParameter_SYSTEM_FILES = PluginParameterType(commonparameter.CommonParameter_SYSTEM_FILES)

    // MCP object and array type parameters
    PluginParameter_ARRAY = PluginParameterType(commonparameter.CommonParameter_ARRAY)
    PluginParameter_OBJECT = PluginParameterType(commonparameter.CommonParameter_OBJECT)
)

func (typ PluginParameterType)AsNormalType()string{
    if slices.Contains([]PluginParameterType{PluginParameter_SECRET_INPUT,PluginParameter_SELECT,},typ) {
        return "string"
        }
    return string(typ)
}

func (typ PluginParameterType) CastParameterValue(value Any) any{
    	defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if real_exp, ok := r.(*exceptions.ValueError); ok {
					panic(real_exp)
				}  else {
                    panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value {%#v} is not in correct type of {%s}.",value, typ.AsNormalType())))
				}
			}

        switch typ{
            case PluginParameter_STRING , PluginParameter_SECRET_INPUT , PluginParameter_SELECT:
                if value ==nil{
                    return ""
                } else {
                    if _, ok := value.(string); ok {
                        return value.(string)
                    } else {
                        real_val, err := cast.ToStringE(value)
                        if err != nil {
                            panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to string failed:%v", value, err)))
                        }
                        return real_val
                    }
                }
            case PluginParameter_BOOLEAN:
                if value == nil{
                    return false
                } else if _, ok := value.(string); ok {
                    // Allowed YAML boolean value strings: https://yaml.org/type/bool.html
                    // and also '0' for False and '1' for True
                    switch strings.ToLower( value){
                        case "true","yes" , "y" , "1":
                            return true
                        case "false" , "no" , "n" , "0":
                            return false
                        default:
                            real_val, err := cast.ToBoolE(value)
                            if err != nil {
                                panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to bool failed:%v", value, err)))
                            }
                            return real_val
                        }
                } else if _, ok := value.(bool); ok {
                    return value.(bool)
                } else {
                            real_val, err := cast.ToBoolE(value)
                            if err != nil {
                                panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to bool failed:%v", value, err)))
                            }
                            return real_val
                }
            case PluginParameter_NUMBER:
                switch real_value := value.(type){
                case int, float64, float32:
                    return real_value
                case string:
                    if strings.Contains(real_value, "."){
                        
                    }
                }
                if isinstance(value, int | float){
                    return value
                } else if isinstance(value, str) and value{
                    if "." in value{
                        return float(value)
                    } else {
                        return int(value)
            case PluginParameter_SYSTEM_FILES | PluginParameter_FILES:
                if not isinstance(value, list){
                    return [value]
                return value
            case PluginParameter_FILE:
                if isinstance(value, list){
                    if len(value) != 1{
                        raise ValueError("This parameter only accepts one file but got multiple files while invoking.")
                    } else {
                        return value[0]
                return value
            case PluginParameter_MODEL_SELECTOR | PluginParameter_APP_SELECTOR:
                if not isinstance(value, dict){
                    raise ValueError("The selector must be a dictionary.")
                return value
            case PluginParameter_TOOLS_SELECTOR:
                if value and not isinstance(value, list){
                    raise ValueError("The tools selector must be a list.")
                return value
            case PluginParameter_ANY:
                if value and not isinstance(value, str | dict | list | NumberType){
                    raise ValueError("The var selector must be a string, dictionary, list or number.")
                return value
            case PluginParameter_ARRAY:
                if not isinstance(value, list){
                    # Try to parse JSON string for arrays
                    if isinstance(value, str){
                        try:
                            import json

                            parsed_value = json.loads(value)
                            if isinstance(parsed_value, list){
                                return parsed_value
                        except (json.JSONDecodeError, ValueError):
                            pass
                    return [value]
                return value
            case PluginParameter_OBJECT:
                if not isinstance(value, dict){
                    # Try to parse JSON string for objects
                    if isinstance(value, str){
                        try:
                            import json

                            parsed_value = json.loads(value)
                            if isinstance(parsed_value, dict){
                                return parsed_value
                        except (json.JSONDecodeError, ValueError):
                            pass
                    return {}
                return value
            case _:
                return str(value)
                }
		}()
            }

type MCPServerParameterType string
const (
    MCPServerParameter_ARRAY MCPServerParameterType= "array"
    MCPServerParameter_OBJECT MCPServerParameterType= "object"
)

type PluginParameterAutoGenerateType string
const (
    PluginParameterAutoGenerate_PROMPT_INSTRUCTION PluginParameterAutoGenerateType = "prompt_instruction"
)
type PluginParameterAutoGenerate struct{
    Type PluginParameterAutoGenerateType `json:"type"`
}
type PluginParameterTemplate struct{
    Enabled bool `json:"enabled"` //description="Whether the parameter is jinja enabled")
}
type PluginParameter[T1 float64|int|string, T2 float64|int] struct{
    Name string `json:"name"` //description="The name of the parameter")
    Label commontypes.I18nObject `json:"label"` //description="The label presented to the user")
    Placeholder *commontypes.I18nObject `json:"placeholder"` //description="The placeholder presented to the user")
    Scope string `json:"scope"`
    AutoGenerate *PluginParameterAutoGenerate `json:"auto_generate"`
    Template *PluginParameterTemplate `json:"template"`
    Required bool `json:"required"`
    Default T1 `json:"default"`
    Min T2 `json:"min"`
    Max T2 `json:"max"`
    Precision int `json:"precision"`
    Options []*PluginParameterOption `json:"options"`
}




def init_frontend_parameter(rule: PluginParameter, type: enum.StrEnum, value: Any):
    """
    init frontend parameter by rule
    """
    parameter_value = value
    if not parameter_value and parameter_value != 0:
        # get default value
        parameter_value = rule.default
        if not parameter_value and rule.required:
            raise ValueError(f"tool parameter {rule.name} not found in tool config")

    if type == PluginParameter_SELECT:
        # check if tool_parameter_config in options
        options = [x.value for x in rule.options]
        if parameter_value is not None and parameter_value not in options:
            raise ValueError(f"tool parameter {rule.name} value {parameter_value} not in options {options}")

    return cast_parameter_value(type, parameter_value)
