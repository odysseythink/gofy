package questionclassifier

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
	"github.com/odysseythink/gofy/backend/entities/nodes/llm"
	"github.com/odysseythink/gofy/backend/entities/prompt"
)

type ClassConfig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// QuestionClassifierNodeData represents answer node data
type QuestionClassifierNodeData struct {
	*basenodesentities.BaseNodeData
	QueryVariableSelector []string             `json:"query_variable_selector"`
	Model                 llm.ModelConfig      `json:"model"`
	Classes               []*ClassConfig       `json:"classes"`
	Instruction           string               `json:"instruction"`
	Memory                *prompt.MemoryConfig `json:"memory"`
	Vision                *llm.VisionConfig    `json:"vision"`
}

func New() *QuestionClassifierNodeData {
	return &QuestionClassifierNodeData{
		Vision: llm.NewVisionConfig(),
	}
}
