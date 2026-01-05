package listoperator

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	listoperatornodesentities "mlib.com/gofy/server/entities/nodes/list_operator"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
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
func (n *ListOperatorNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	typed_node_data, err := mapstruct.MapToStruct1[*listoperatornodesentities.ListOperatorNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)

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
