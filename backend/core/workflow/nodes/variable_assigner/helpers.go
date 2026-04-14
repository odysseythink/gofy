package variableassigner

import (
	"encoding/json"
	"slices"

	"github.com/odysseythink/mlog"
	varassignernodesexceptions "mlib.com/gofy/server/core/exceptions/nodes/variable_assigner"
	"mlib.com/gofy/server/core/variables"
	dbengine "mlib.com/gofy/server/db_engine"
	vaenumtypes "mlib.com/gofy/server/enum_types/nodes/variable_assigner"
	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
	"mlib.com/gofy/server/models"
)

func update_conversation_variable(conversation_id string, variable variables.Variabler) {
	conversation_variable := new(models.ConversationVariable)
	err := dbengine.Instance().DB.Model(&models.ConversationVariable{}).Where("id = ? and conversation_id =?", variable.GetID(), conversation_id).First(conversation_variable).Error
	if err != nil {
		mlog.Errorf("get ConversationVariable failed:%v", err)
		panic(varassignernodesexceptions.NewVariableOperatorNodeError("conversation variable not found in the database"))
	}
	bindata, _ := json.Marshal(variable.ToDict(variable))
	conversation_variable.Data = string(bindata)
	dbengine.Instance().DB.Updates(&models.ConversationVariable{ID: conversation_variable.ID, Data: string(bindata)})
}
func is_operation_supported(variable_type variableenumtypes.VariableType, operation vaenumtypes.OperationType) bool {
	switch operation {
	case vaenumtypes.Operation_OVER_WRITE:
		return true
	case vaenumtypes.Operation_CLEAR:
		return true
	case vaenumtypes.Operation_SET:
		return slices.Contains([]variableenumtypes.VariableType{variableenumtypes.Variable_OBJECT, variableenumtypes.Variable_STRING, variableenumtypes.Variable_NUMBER}, variable_type)
	case vaenumtypes.Operation_ADD:
		// Only number variable can be added, subtracted, multiplied or divided
		return variable_type == variableenumtypes.Variable_NUMBER
	case vaenumtypes.Operation_SUBTRACT:
		// Only number variable can be added, subtracted, multiplied or divided
		return variable_type == variableenumtypes.Variable_NUMBER
	case vaenumtypes.Operation_MULTIPLY:
		// Only number variable can be added, subtracted, multiplied or divided
		return variable_type == variableenumtypes.Variable_NUMBER
	case vaenumtypes.Operation_DIVIDE:
		// Only number variable can be added, subtracted, multiplied or divided
		return variable_type == variableenumtypes.Variable_NUMBER
	case vaenumtypes.Operation_EXTEND:
		// Only array variable can be appended or extended
		return slices.Contains([]variableenumtypes.VariableType{
			variableenumtypes.Variable_ARRAY_ANY,
			variableenumtypes.Variable_ARRAY_OBJECT,
			variableenumtypes.Variable_ARRAY_STRING,
			variableenumtypes.Variable_ARRAY_NUMBER,
			variableenumtypes.Variable_ARRAY_FILE,
		}, variable_type)
	case vaenumtypes.Operation_APPEND:
		// Only array variable can be appended or extended
		return slices.Contains([]variableenumtypes.VariableType{
			variableenumtypes.Variable_ARRAY_ANY,
			variableenumtypes.Variable_ARRAY_OBJECT,
			variableenumtypes.Variable_ARRAY_STRING,
			variableenumtypes.Variable_ARRAY_NUMBER,
			variableenumtypes.Variable_ARRAY_FILE,
		}, variable_type)
	}
	return false
}
func is_variable_input_supported(operation vaenumtypes.OperationType) bool {
	if slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_SET, vaenumtypes.Operation_ADD, vaenumtypes.Operation_SUBTRACT, vaenumtypes.Operation_MULTIPLY, vaenumtypes.Operation_DIVIDE}, operation) {
		return false
	}
	return true

}
func is_constant_input_supported(variable_type variableenumtypes.VariableType, operation vaenumtypes.OperationType) bool {
	switch variable_type {
	case variableenumtypes.Variable_STRING:
		return slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_OVER_WRITE, vaenumtypes.Operation_SET}, operation)
	case variableenumtypes.Variable_OBJECT:
		return slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_OVER_WRITE, vaenumtypes.Operation_SET}, operation)
	case variableenumtypes.Variable_NUMBER:
		return slices.Contains([]vaenumtypes.OperationType{
			vaenumtypes.Operation_OVER_WRITE,
			vaenumtypes.Operation_SET,
			vaenumtypes.Operation_ADD,
			vaenumtypes.Operation_SUBTRACT,
			vaenumtypes.Operation_MULTIPLY,
			vaenumtypes.Operation_DIVIDE,
		}, operation)
	default:
		return false
	}
}
func is_input_value_valid(variable_type variableenumtypes.VariableType, operation vaenumtypes.OperationType, value any) bool {
	if operation == vaenumtypes.Operation_CLEAR {
		return true
	}
	switch variable_type {
	case variableenumtypes.Variable_STRING:
		_, ok := value.(string)
		return ok

	case variableenumtypes.Variable_NUMBER:
		switch real_value := value.(type) {
		case int:
			if operation == vaenumtypes.Operation_DIVIDE && real_value == 0 {
				return false
			}
		case float32:
			if operation == vaenumtypes.Operation_DIVIDE && real_value == 0. {
				return false
			}
		case float64:
			if operation == vaenumtypes.Operation_DIVIDE && real_value == 0. {
				return false
			}
		default:
			return false
		}

		return true

	case variableenumtypes.Variable_OBJECT:
		_, ok := value.(map[string]any)
		return ok

	// Array & Append
	case variableenumtypes.Variable_ARRAY_ANY:
		if operation == vaenumtypes.Operation_APPEND {
			switch value.(type) {
			case int:
				return true
			case float32:
				return true
			case float64:
				return true
			case string:
				return true
			case map[string]any:
				return true
			default:
				return false
			}
		} else if slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_EXTEND, vaenumtypes.Operation_OVER_WRITE}, operation) {
			if real_value, ok := value.([]any); ok {
				for _, item := range real_value {
					switch item.(type) {
					case int:
					case float32:
					case float64:
					case string:
					case map[string]any:
					default:
						return false
					}
				}
				return true
			} else {
				return false
			}
		}
	case variableenumtypes.Variable_ARRAY_STRING:
		if operation == vaenumtypes.Operation_APPEND {
			_, ok := value.(string)
			return ok
		} else if slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_EXTEND, vaenumtypes.Operation_OVER_WRITE}, operation) {
			if real_value, ok := value.([]any); ok {
				for _, item := range real_value {
					switch item.(type) {
					case string:
					default:
						return false
					}
				}
				return true
			} else if _, ok := value.([]string); ok {
				return true
			} else {
				return false
			}
		}
	case variableenumtypes.Variable_ARRAY_NUMBER:
		if operation == vaenumtypes.Operation_APPEND {
			switch value.(type) {
			case int:
				return true
			case float32:
				return true
			case float64:
				return true
			default:
				return false
			}
		} else if slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_EXTEND, vaenumtypes.Operation_OVER_WRITE}, operation) {
			if real_value, ok := value.([]any); ok {
				for _, item := range real_value {
					switch item.(type) {
					case int:
					case float32:
					case float64:
					default:
						return false
					}
				}
				return true
			} else if _, ok := value.([]int); ok {
				return true
			} else if _, ok := value.([]float32); ok {
				return true
			} else if _, ok := value.([]float64); ok {
				return true
			} else {
				return false
			}
		}
	case variableenumtypes.Variable_ARRAY_OBJECT:
		if operation == vaenumtypes.Operation_APPEND {
			_, ok := value.(map[string]any)
			return ok
		} else if slices.Contains([]vaenumtypes.OperationType{vaenumtypes.Operation_EXTEND, vaenumtypes.Operation_OVER_WRITE}, operation) {
			if real_value, ok := value.([]any); ok {
				for _, item := range real_value {
					switch item.(type) {
					case map[string]any:
					default:
						return false
					}
				}
				return true
			} else if _, ok := value.([]map[string]any); ok {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
