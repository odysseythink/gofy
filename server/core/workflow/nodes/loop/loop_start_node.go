package loop

import (
	"iter"

	"mlib.com/gofy/server/core/workflow/nodes/base"
	loopnodesentities "mlib.com/gofy/server/entities/nodes/loop"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type LoopStartNode struct {
	*base.BaseNode[*loopnodesentities.LoopStartNodeData]
}

func (n *LoopStartNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_LOOP_START
}

func (n *LoopStartNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *LoopStartNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *loopnodesentities.LoopStartNodeData) map[string][]string {
	return map[string][]string{}
}

func NewLoopStartNode() *LoopStartNode {
	return &LoopStartNode{
		BaseNode: &base.BaseNode[*loopnodesentities.LoopStartNodeData]{},
	}
}
