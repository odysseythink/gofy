package tool

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	toolnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/tool"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
)

type ToolNode struct {
	*base.BaseNode[*toolnodesentities.ToolNodeData]
}

func (n *ToolNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_TOOL
}

func (n *ToolNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *ToolNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *toolnodesentities.ToolNodeData) map[string][]string {

	return map[string][]string{}
}
func New() *ToolNode {
	return &ToolNode{
		BaseNode: &base.BaseNode[*toolnodesentities.ToolNodeData]{},
	}
}
