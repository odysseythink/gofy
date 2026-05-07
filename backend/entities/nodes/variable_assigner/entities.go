package variableassigner

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
	vaenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes/variable_assigner"
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
