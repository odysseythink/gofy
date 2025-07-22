package templatetransform

import (
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	workflowentities "mlib.com/gofy/server/entities/workflow"
)

// TemplateTransformNodeData represents answer node data
type TemplateTransformNodeData struct {
	*basenodesentities.BaseNodeData
	Variables []*workflowentities.VariableSelector
	Template  string
}
