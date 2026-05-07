package listoperator

import (
	"iter"

	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	listoperatornodesentities "github.com/odysseythink/gofy/backend/entities/nodes/list_operator"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
)

type ListOperatorNode struct {
	*base.BaseNode[*listoperatornodesentities.ListOperatorNodeData]
}

func (n *ListOperatorNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_LIST_OPERATOR
}

func (n *ListOperatorNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}

func (n *ListOperatorNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *listoperatornodesentities.ListOperatorNodeData) map[string][]string {

	return map[string][]string{}
}
func New() *ListOperatorNode {
	return &ListOperatorNode{
		BaseNode: &base.BaseNode[*listoperatornodesentities.ListOperatorNodeData]{},
	}
}
