package variableassigner

import (
	"encoding/json"
	"iter"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/constants"
	varassignernodesexceptions "mlib.com/gofy/server/core/exceptions/nodes/variable_assigner"
	"mlib.com/gofy/server/core/variables"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	variableassignernodesentities "mlib.com/gofy/server/entities/nodes/variable_assigner"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	vaenumtypes "mlib.com/gofy/server/enum_types/nodes/variable_assigner"
	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
)

var (
	EMPTY_VALUE_MAPPING = map[variableenumtypes.VariableType]any{
		variableenumtypes.Variable_STRING:       "",
		variableenumtypes.Variable_NUMBER:       0,
		variableenumtypes.Variable_OBJECT:       map[string]any{},
		variableenumtypes.Variable_ARRAY_ANY:    []any{},
		variableenumtypes.Variable_ARRAY_STRING: []string{},
		variableenumtypes.Variable_ARRAY_NUMBER: []any{},
		variableenumtypes.Variable_ARRAY_OBJECT: []map[string]any{},
	}
)

type VariableAssignerNode struct {
	*base.BaseNode[*variableassignernodesentities.VariableAssignerNodeData]
}

func (n *VariableAssignerNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_VARIABLE_ASSIGNER
}

func (n *VariableAssignerNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	inputs := n.NodeData.ModelDump()
	process_data := map[string]any{}
	// NOTE: This node has no outputs
	updated_variable_selectors := [][]string{}
	var ret *workflowentities.NodeRunResult
	func() {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if exp, ok := r.(*varassignernodesexceptions.VariableOperatorNodeError); ok {
					ret = &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Inputs:      inputs,
						ProcessData: process_data,
						Error:       exp.Error(),
					}
				} else if exp, ok := r.(error); ok {
					ret = &workflowentities.NodeRunResult{
						Status:      models.WorkflowNodeExecutionStatus_FAILED,
						Inputs:      inputs,
						ProcessData: process_data,
						Error:       exp.Error(),
					}
				} else {
					panic(r)
				}
			}
		}()
		for _, item := range n.NodeData.Items {
			variable := n.GetGraphRuntimeState().VariablePool.Get(item.VariableSelector)

			// ==================== Validation Part

			// Check if variable exists
			if variable == nil {
				panic(varassignernodesexceptions.NewVariableNotFoundError(item.VariableSelector, ""))
			}
			// Check if operation is supported
			if !is_operation_supported(variable.ValueType(), item.Operation) {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(item.Operation, string(variable.ValueType()), ""))
			}
			// Check if variable input is supported
			if item.InputType == vaenumtypes.Input_VARIABLE && !is_variable_input_supported(item.Operation) {
				panic(varassignernodesexceptions.NewInputTypeNotSupportedError(item.Operation, vaenumtypes.Input_VARIABLE, ""))
			}
			// Check if constant input is supported
			if item.InputType == vaenumtypes.Input_CONSTANT && !is_constant_input_supported(variable.ValueType(), item.Operation) {
				panic(varassignernodesexceptions.NewInputTypeNotSupportedError(item.Operation, vaenumtypes.Input_CONSTANT, ""))
			}
			// Get value from variable pool
			if item.InputType == vaenumtypes.Input_VARIABLE &&
				item.Operation != vaenumtypes.Operation_CLEAR &&
				item.Value != nil {
				var selector []string
				if _, ok := item.Value.([]string); ok {
					selector = item.Value.([]string)
				} else if _, ok := item.Value.([]any); ok {
					for _, v := range item.Value.([]any) {
						if _, ok := v.(string); !ok {
							mlog.Errorf("item.Value=%#v must be type of []string", item.Value)
							panic(varassignernodesexceptions.NewVariableNotFoundError(nil, ""))
						} else {
							if selector == nil {
								selector = make([]string, 0)
							}
							selector = append(selector, v.(string))
						}
					}
				} else {
					mlog.Errorf("item.Value=%#v must be type of []string", item.Value)
					panic(varassignernodesexceptions.NewVariableNotFoundError(nil, ""))
				}
				value := n.GetGraphRuntimeState().VariablePool.Get(selector)
				if value == nil {
					panic(varassignernodesexceptions.NewVariableNotFoundError(selector, ""))
				}
				// Skip if value is NoneSegment
				if value.ValueType() == variableenumtypes.Variable_NONE {
					continue
				}
				item.Value = value.GetValue()
			}
			// If set string / bytes / bytearray to object, try convert string to object.
			_, ok := item.Value.(string)
			if item.Operation == vaenumtypes.Operation_SET &&
				variable.ValueType() == variableenumtypes.Variable_OBJECT &&
				ok {
				tmpmap := map[string]any{}
				err := json.Unmarshal([]byte(item.Value.(string)), &tmpmap)
				if err != nil {
					tmplist := []any{}
					err = json.Unmarshal([]byte(item.Value.(string)), &tmplist)
					if err != nil {
						mlog.Errorf("can't json unmarshal value=%s to map or list:%v", item.Value, err)
						panic(varassignernodesexceptions.NewInvalidInputValueError(item.Value, ""))
					} else {
						item.Value = tmplist
					}
				} else {
					item.Value = tmpmap
				}
			}
			// Check if input value is valid
			if !is_input_value_valid(variable.ValueType(), item.Operation, item.Value) {
				panic(varassignernodesexceptions.NewInvalidInputValueError(item.Value, ""))
			}
			// ==================== Execution Part

			updated_value := n._handle_item(
				variable,
				item.Operation,
				item.Value,
			)
			variable.SetValue(updated_value)
			// variable = variable.model_copy(update={"value": updated_value})
			n.GetGraphRuntimeState().VariablePool.Add(variable.GetSelector(), variable)
			updated_variable_selectors = append(updated_variable_selectors, variable.GetSelector())
		}
	}()
	if ret != nil {
		return ret, nil
	}

	// The `updated_variable_selectors` is a list contains list[str] which not hashable,
	// remove the duplicated items first.

	// Update variables
	for _, selector := range updated_variable_selectors {
		variable := n.GetGraphRuntimeState().VariablePool.Get(selector)
		if variable == nil {
			panic(varassignernodesexceptions.NewVariableNotFoundError(selector, ""))
		}
		process_data[variable.GetName()] = variable.GetValue()

		if variable.GetSelector()[0] == constants.CONVERSATION_VARIABLE_NODE_ID {
			conversation_id := n.GetGraphRuntimeState().VariablePool.Get([]string{"sys", "conversation_id"})
			if conversation_id == nil {
				panic(varassignernodesexceptions.NewConversationIDNotFoundError(""))
			}
			if _, ok := conversation_id.GetValue().(string); !ok {
				mlog.Errorf("sys.conversation_id=%#v must be string", conversation_id.GetValue())
				panic(varassignernodesexceptions.NewConversationIDNotFoundError("sys.conversation_id must be string"))
			}
			update_conversation_variable(conversation_id.GetValue().(string), variable)
		}
	}

	return &workflowentities.NodeRunResult{
		Status:      models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:      inputs,
		ProcessData: process_data,
	}, nil
}

func (n *VariableAssignerNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *variableassignernodesentities.VariableAssignerNodeData) map[string][]string {

	return map[string][]string{}
}

func (n *VariableAssignerNode) _handle_item(
	variable variables.Variabler,
	operation vaenumtypes.OperationType,
	value any,
) any {
	switch operation {
	case vaenumtypes.Operation_OVER_WRITE:
		return value
	case vaenumtypes.Operation_CLEAR:
		return EMPTY_VALUE_MAPPING[variable.ValueType()]
	case vaenumtypes.Operation_APPEND:
		if real_var_value, ok := variable.GetValue().([]any); ok {
			switch real_val := value.(type) {
			case int:
				return append(real_var_value, real_val)
			case float32:
				return append(real_var_value, real_val)
			case float64:
				return append(real_var_value, real_val)
			case string:
				return append(real_var_value, real_val)
			case map[string]any:
				return append(real_var_value, real_val)
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		} else if real_var_value, ok := variable.GetValue().([]string); ok {
			if _, ok := value.(string); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return append(real_var_value, value.(string))
			}
		} else if real_var_value, ok := variable.GetValue().([]int); ok {
			if _, ok := value.(int); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return append(real_var_value, value.(int))
			}
		} else if real_var_value, ok := variable.GetValue().([]float32); ok {
			if _, ok := value.(float32); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return append(real_var_value, value.(float32))
			}
		} else if real_var_value, ok := variable.GetValue().([]float64); ok {
			if _, ok := value.(float64); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return append(real_var_value, value.(float64))
			}
		} else if real_var_value, ok := variable.GetValue().([]map[string]any); ok {
			if _, ok := value.(map[string]any); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return append(real_var_value, value.(map[string]any))
			}
		}
		panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
	case vaenumtypes.Operation_EXTEND:
		if real_var_value, ok := variable.GetValue().([]any); ok {
			switch real_val := value.(type) {
			case []any:
				return append(real_var_value, real_val...)
			case []int:
				for _, item := range real_val {
					real_var_value = append(real_var_value, item)
				}
				return real_var_value
			case []float32:
				for _, item := range real_val {
					real_var_value = append(real_var_value, item)
				}
				return real_var_value
			case []float64:
				for _, item := range real_val {
					real_var_value = append(real_var_value, item)
				}
				return real_var_value
			case []string:
				for _, item := range real_val {
					real_var_value = append(real_var_value, item)
				}
				return real_var_value
			case []map[string]any:
				for _, item := range real_val {
					real_var_value = append(real_var_value, item)
				}
				return real_var_value
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		} else if real_var_value, ok := variable.GetValue().([]string); ok {
			switch real_val := value.(type) {
			case []any:
				for _, item := range real_val {
					if _, ok := item.(string); !ok {
						panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
					} else {
						real_var_value = append(real_var_value, item.(string))
					}
				}
				return real_var_value
			case []string:
				return append(real_var_value, real_val...)
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		} else if real_var_value, ok := variable.GetValue().([]int); ok {
			switch real_val := value.(type) {
			case []any:
				for _, item := range real_val {
					if _, ok := item.(int); !ok {
						panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
					} else {
						real_var_value = append(real_var_value, item.(int))
					}
				}
				return real_var_value
			case []int:
				return append(real_var_value, real_val...)
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		} else if real_var_value, ok := variable.GetValue().([]float32); ok {
			switch real_val := value.(type) {
			case []any:
				for _, item := range real_val {
					if _, ok := item.(float32); !ok {
						panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
					} else {
						real_var_value = append(real_var_value, item.(float32))
					}
				}
				return real_var_value
			case []float32:
				return append(real_var_value, real_val...)
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		} else if real_var_value, ok := variable.GetValue().([]float64); ok {
			switch real_val := value.(type) {
			case []any:
				for _, item := range real_val {
					if _, ok := item.(float64); !ok {
						panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
					} else {
						real_var_value = append(real_var_value, item.(float64))
					}
				}
				return real_var_value
			case []float64:
				return append(real_var_value, real_val...)
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		} else if real_var_value, ok := variable.GetValue().([]map[string]any); ok {
			switch real_val := value.(type) {
			case []any:
				for _, item := range real_val {
					if _, ok := item.(map[string]any); !ok {
						panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
					} else {
						real_var_value = append(real_var_value, item.(map[string]any))
					}
				}
				return real_var_value
			case []map[string]any:
				return append(real_var_value, real_val...)
			}
			panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
		}
		panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
	case vaenumtypes.Operation_SET:
		return value
	case vaenumtypes.Operation_ADD:
		if real_var_value, ok := variable.GetValue().(int); ok {
			if _, ok := value.(int); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value + value.(int)
			}
		} else if real_var_value, ok := variable.GetValue().(float32); ok {
			if _, ok := value.(float32); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value + value.(float32)
			}
		} else if real_var_value, ok := variable.GetValue().(float64); ok {
			if _, ok := value.(float64); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value + value.(float64)
			}
		}
		panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
	case vaenumtypes.Operation_SUBTRACT:
		if real_var_value, ok := variable.GetValue().(int); ok {
			if _, ok := value.(int); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value - value.(int)
			}
		} else if real_var_value, ok := variable.GetValue().(float32); ok {
			if _, ok := value.(float32); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value - value.(float32)
			}
		} else if real_var_value, ok := variable.GetValue().(float64); ok {
			if _, ok := value.(float64); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value - value.(float64)
			}
		}
		panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
	case vaenumtypes.Operation_MULTIPLY:
		if real_var_value, ok := variable.GetValue().(int); ok {
			if _, ok := value.(int); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value * value.(int)
			}
		} else if real_var_value, ok := variable.GetValue().(float32); ok {
			if _, ok := value.(float32); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value * value.(float32)
			}
		} else if real_var_value, ok := variable.GetValue().(float64); ok {
			if _, ok := value.(float64); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value * value.(float64)
			}
		}
		panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
	case vaenumtypes.Operation_DIVIDE:
		if real_var_value, ok := variable.GetValue().(int); ok {
			if _, ok := value.(int); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value / value.(int)
			}
		} else if real_var_value, ok := variable.GetValue().(float32); ok {
			if _, ok := value.(float32); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value / value.(float32)
			}
		} else if real_var_value, ok := variable.GetValue().(float64); ok {
			if _, ok := value.(float64); !ok {
				panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
			} else {
				return real_var_value / value.(float64)
			}
		}
		panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))

	}
	panic(varassignernodesexceptions.NewOperationNotSupportedError(operation, string(variable.ValueType()), ""))
}
func New() *VariableAssignerNode {
	return &VariableAssignerNode{
		BaseNode: &base.BaseNode[*variableassignernodesentities.VariableAssignerNodeData]{},
	}
}
