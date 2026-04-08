package documentextractor

import (
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
)

// CodeNodeData represents answer node data
type DocumentExtractorNodeData struct {
	*basenodesentities.BaseNodeData
	VariableSelector []string `json:"variable_selector"`
}
