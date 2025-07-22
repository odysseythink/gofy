package end

import (
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	workflowentities "mlib.com/gofy/server/entities/workflow"
)

// AnswerNodeData represents answer node data
type EndNodeData struct {
	*basenodesentities.BaseNodeData
	Outputs []*workflowentities.VariableSelector `json:"outputs"`
}

// EndStreamParam entity
type EndStreamParam struct {
	EndDependencies                  map[string][]string   `json:"end_dependencies"`
	EndStreamVariableSelectorMapping map[string][][]string `json:"end_stream_variable_selector_mapping"`
}

func New() *EndNodeData {
	return &EndNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
		Outputs:      make([]*workflowentities.VariableSelector, 0),
	}
}
