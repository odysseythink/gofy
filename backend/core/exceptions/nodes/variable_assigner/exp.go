package variableassigner

import (
	"fmt"

	vaenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes/variable_assigner"
)

type VariableOperatorNodeError struct {
	desc string
}

func NewVariableOperatorNodeError(desc string) *VariableOperatorNodeError {
	return &VariableOperatorNodeError{
		desc: desc,
	}
}
func (e *VariableOperatorNodeError) Error() string {
	return fmt.Sprintf("variable operator node error:%s", e.desc)
}

func NewOperationNotSupportedError(operation vaenumtypes.OperationType, variable_type, desc string) *VariableOperatorNodeError {
	desc = fmt.Sprintf("Operation {%s} is not supported for type {%s}:", operation, variable_type) + desc
	return NewVariableOperatorNodeError(desc)
}

func NewInputTypeNotSupportedError(operation vaenumtypes.OperationType, input_type vaenumtypes.InputType, desc string) *VariableOperatorNodeError {
	desc = fmt.Sprintf("Input type {%s} is not supported for operation {%s}:", input_type, operation) + desc
	return NewVariableOperatorNodeError(desc)
}

func NewVariableNotFoundError(variable_selector []string, desc string) *VariableOperatorNodeError {
	desc = fmt.Sprintf("Variable {%v} not found:", variable_selector) + desc
	return NewVariableOperatorNodeError(desc)
}

func NewInvalidInputValueError(value any, desc string) *VariableOperatorNodeError {
	desc = fmt.Sprintf("Invalid input value {%#v}:", value) + desc
	return NewVariableOperatorNodeError(desc)
}

func NewConversationIDNotFoundError(desc string) *VariableOperatorNodeError {
	desc = "conversation_id not found:" + desc
	return NewVariableOperatorNodeError(desc)
}
