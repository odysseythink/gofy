package templatetransform

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
)

// TemplateTransformNodeData represents answer node data
type TemplateTransformNodeData struct {
	*basenodesentities.BaseNodeData
	Variables []*workflowentities.VariableSelector
	Template  string
}
