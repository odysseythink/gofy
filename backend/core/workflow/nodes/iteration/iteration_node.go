package iteration

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	iterationnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/iteration"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
)

type IterationNode struct {
	*base.BaseNode[*iterationnodesentities.IterationNodeData]
}

func (n *IterationNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_ITERATION
}

func (n *IterationNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *IterationNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *iterationnodesentities.IterationNodeData) map[string][]string {
	return map[string][]string{}
}
func New() *IterationNode {
	return &IterationNode{
		BaseNode: &base.BaseNode[*iterationnodesentities.IterationNodeData]{},
	}
}
