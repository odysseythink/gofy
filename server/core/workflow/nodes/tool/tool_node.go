package tool

import (
	"iter"

	"mlib.com/gofy/server/core/workflow/nodes/base"
	toolnodesentities "mlib.com/gofy/server/entities/nodes/tool"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
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
