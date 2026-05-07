package documentextractor

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
)

// CodeNodeData represents answer node data
type DocumentExtractorNodeData struct {
	*basenodesentities.BaseNodeData
	VariableSelector []string `json:"variable_selector"`
}
