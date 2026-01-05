package listoperator

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	listoperatornodesentities "mlib.com/gofy/server/entities/nodes/list_operator"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
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
func init() {
	nodesconstants.Regist(New())
}
