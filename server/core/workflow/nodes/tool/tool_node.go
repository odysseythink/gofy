package tool

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	toolnodesentities "mlib.com/gofy/server/entities/nodes/tool"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
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
func (n *ToolNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	typed_node_data, err := mapstruct.MapToStruct1[*toolnodesentities.ToolNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)

	return map[string][]string{}
}
func New() *ToolNode {
	return &ToolNode{
		BaseNode: &base.BaseNode[*toolnodesentities.ToolNodeData]{},
	}
}
func init() {
	nodesconstants.Regist(New())
}
