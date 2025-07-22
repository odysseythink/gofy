package variableassigner

import (
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	vaenumtypes "mlib.com/gofy/server/enum_types/nodes/variable_assigner"
)

type VariableOperationItem struct {
	VariableSelector      []string `json:"variable_selector"`
	vaenumtypes.InputType `json:"input_type"`
	Operation             vaenumtypes.OperationType `json:"operation"`
	Value                 any                       `json:"value"`
}

// VariableAssignerNodeData represents answer node data
type VariableAssignerNodeData struct {
	*basenodesentities.BaseNodeData
	Version string                  `json:"version"`
	Items   []VariableOperationItem `json:"items"`
}
