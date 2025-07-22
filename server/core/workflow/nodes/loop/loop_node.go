package loop

import (
	"iter"

	"mlib.com/gofy/server/core/workflow/nodes/base"
	loopnodesentities "mlib.com/gofy/server/entities/nodes/loop"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type LoopNode struct {
	*base.BaseNode[*loopnodesentities.LoopNodeData]
}

func (n *LoopNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_LOOP
}

func (n *LoopNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *LoopNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *loopnodesentities.LoopNodeData) map[string][]string {

	return map[string][]string{}
}
func New() *LoopNode {
	return &LoopNode{
		BaseNode: &base.BaseNode[*loopnodesentities.LoopNodeData]{},
	}
}
