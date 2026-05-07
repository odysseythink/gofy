package loop

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	loopnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/loop"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
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
