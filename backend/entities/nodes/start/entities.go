package start

import (
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
)

// AnswerNodeData represents answer node data
type StartNodeData struct {
	*basenodesentities.BaseNodeData
	Variables []*appconfigentities.VariableEntity `json:"outputs"`
}

func New() *StartNodeData {
	return &StartNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
		Variables:    make([]*appconfigentities.VariableEntity, 0),
	}
}
