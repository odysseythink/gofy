package code

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
)

type Output struct {
	Type     string //: Literal["string", "number", "object", "array[string]", "array[number]", "array[object]"]
	Children map[string]*Output
}
type Dependency struct {
	Name    string
	Version string
}

// CodeNodeData represents answer node data
type CodeNodeData struct {
	*basenodesentities.BaseNodeData
	Variables    []*workflowentities.VariableSelector `json:"variables"`
	CodeLanguage string                               `json:"code_language"` // : Literal[CodeLanguage.PYTHON3, CodeLanguage.JAVASCRIPT]
	Code         string                               `json:"code"`
	Outputs      map[string]*Output                   `json:"outputs"`
	Dependencies []*Dependency                        `json:"dependencies"`
}

func New() *CodeNodeData {
	return &CodeNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
		Variables:    make([]*workflowentities.VariableSelector, 0),
		Outputs:      make(map[string]*Output),
		Dependencies: make([]*Dependency, 0),
	}
}
