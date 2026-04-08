package start

import (
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
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
