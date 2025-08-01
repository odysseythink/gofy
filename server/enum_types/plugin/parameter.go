package plugin

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"mlib.com/gofy/server/core/exceptions"
	parameterenumtypes "mlib.com/gofy/server/enum_types/parameter"
	"mlib.com/gofy/server/utils"
	"mlib.com/mconfig/cast"
	"mlib.com/mlog"
)

type PluginParameterType string

const (
	PluginParameter_STRING         = PluginParameterType(parameterenumtypes.CommonParameter_STRING)
	PluginParameter_NUMBER         = PluginParameterType(parameterenumtypes.CommonParameter_NUMBER)
	PluginParameter_BOOLEAN        = PluginParameterType(parameterenumtypes.CommonParameter_BOOLEAN)
	PluginParameter_SELECT         = PluginParameterType(parameterenumtypes.CommonParameter_SELECT)
	PluginParameter_SECRET_INPUT   = PluginParameterType(parameterenumtypes.CommonParameter_SECRET_INPUT)
	PluginParameter_FILE           = PluginParameterType(parameterenumtypes.CommonParameter_FILE)
	PluginParameter_FILES          = PluginParameterType(parameterenumtypes.CommonParameter_FILES)
	PluginParameter_APP_SELECTOR   = PluginParameterType(parameterenumtypes.CommonParameter_APP_SELECTOR)
	PluginParameter_MODEL_SELECTOR = PluginParameterType(parameterenumtypes.CommonParameter_MODEL_SELECTOR)
	PluginParameter_TOOLS_SELECTOR = PluginParameterType(parameterenumtypes.CommonParameter_TOOLS_SELECTOR)
	PluginParameter_ANY            = PluginParameterType(parameterenumtypes.CommonParameter_ANY)
	PluginParameter_DYNAMIC_SELECT = PluginParameterType(parameterenumtypes.CommonParameter_DYNAMIC_SELECT)

	// deprecated, should not use.
	PluginParameter_SYSTEM_FILES = PluginParameterType(parameterenumtypes.CommonParameter_SYSTEM_FILES)

	// MCP object and array type parameters
	PluginParameter_ARRAY  = PluginParameterType(parameterenumtypes.CommonParameter_ARRAY)
	PluginParameter_OBJECT = PluginParameterType(parameterenumtypes.CommonParameter_OBJECT)
)

func AsNormalType(typ string) string {
	if slices.Contains([]string{string(PluginParameter_SECRET_INPUT), string(PluginParameter_SELECT)}, typ) {
		return "string"
	}
	return string(typ)
}

func CastParameterValue(typ string, value any) any {
	defer func() {
		if r := recover(); r != nil {
			mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
			if real_exp, ok := r.(*exceptions.ValueError); ok {
				panic(real_exp)
			} else {
				panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value {%#v} is not in correct type of {%s}.", value, AsNormalType(typ))))
			}
		}
	}()
	switch PluginParameterType(typ) {
	case PluginParameter_STRING, PluginParameter_SECRET_INPUT, PluginParameter_SELECT:
		if value == nil {
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
		if value == nil {
			return false
		} else if _, ok := value.(string); ok {
			// Allowed YAML boolean value strings: https://yaml.org/type/bool.html
			// and also '0' for False and '1' for True
			switch strings.ToLower(value.(string)) {
			case "true", "yes", "y", "1":
				return true
			case "false", "no", "n", "0":
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
		switch real_value := value.(type) {
		case int, float64, float32:
			return real_value
		case string:
			if strings.Contains(real_value, ".") {
				tmp, err := strconv.ParseFloat(real_value, 64)
				if err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to float failed:%v", value, err)))
				}
				return tmp
			} else {
				tmp, err := strconv.ParseInt(real_value, 10, 64)
				if err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to float failed:%v", value, err)))
				}
				return tmp
			}
		}
	case PluginParameter_SYSTEM_FILES, PluginParameter_FILES:
		if _, ok := value.([]any); !ok {
			return []any{value}
		}
		return value
	case PluginParameter_FILE:
		if _, ok := value.([]any); ok {
			if len(value.([]any)) != 1 {
				panic(exceptions.NewValueError("This parameter only accepts one file but got multiple files while invoking."))
			} else {
				return value.([]any)[0]
			}
		}
		return value
	case PluginParameter_MODEL_SELECTOR, PluginParameter_APP_SELECTOR:
		if _, ok := value.(map[string]any); !ok {
			panic(exceptions.NewValueError("The selector must be a dictionary."))
		}
		return value
	case PluginParameter_TOOLS_SELECTOR:
		if value != nil {
			if _, ok := value.([]any); !ok {
				panic(exceptions.NewValueError("The tools selector must be a list."))
			}
		}
		return value
	case PluginParameter_ANY:
		if value != nil {
			switch value.(type) {
			case string, map[string]any, []any, int, float32, float64:
			default:
				panic(exceptions.NewValueError("The var selector must be a string, dictionary, list or number."))
			}
		}
		return value
	case PluginParameter_ARRAY:
		if _, ok := value.([]any); !ok {
			// Try to parse JSON string for arrays
			if _, ok := value.(string); ok {
				tmplist := []any{}
				if err := json.Unmarshal([]byte(value.(string)), &tmplist); err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to list failed:%v", value, err)))
				}
				return tmplist
			}
			return []any{}
		}
		return value
	case PluginParameter_OBJECT:
		if _, ok := value.(map[string]any); !ok {
			// Try to parse JSON string for objects
			if _, ok := value.(string); ok {
				tmplist := map[string]any{}
				if err := json.Unmarshal([]byte(value.(string)), &tmplist); err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to list failed:%v", value, err)))
				}
				return tmplist
			}
			return map[string]any{}
		}
		return value
	default:
		tmp, err := cast.ToStringE(value)
		if err != nil {
			panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} to string failed:%v", value, err)))
		}
		return tmp
	}
	panic(exceptions.NewValueError(fmt.Sprintf("value {%#v} is invalid", value)))
}

type MCPServerParameterType string

const (
	MCPServerParameter_ARRAY  MCPServerParameterType = "array"
	MCPServerParameter_OBJECT MCPServerParameterType = "object"
)

type PluginParameterAutoGenerateType string

const (
	PluginParameterAutoGenerate_PROMPT_INSTRUCTION PluginParameterAutoGenerateType = "prompt_instruction"
)
