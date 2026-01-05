package iteration

import (
	"iter"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	iterationnodesentities "mlib.com/gofy/server/entities/nodes/iteration"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
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
func (n *IterationNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	typed_node_data, err := mapstruct.MapToStruct1[*iterationnodesentities.IterationNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)
	return map[string][]string{}
}
func New() *IterationNode {
	return &IterationNode{
		BaseNode: &base.BaseNode[*iterationnodesentities.IterationNodeData]{},
	}
}
func init() {
	nodesconstants.Regist(New())
}
