package start

import (
	"iter"

	"mlib.com/gofy/server/constants"
	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	startnodesentities "mlib.com/gofy/server/entities/nodes/start"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

type StartNode struct {
	*base.BaseNode[*startnodesentities.StartNodeData]
}

func (n *StartNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_START
}

func (n *StartNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	node_inputs := n.GetGraphRuntimeState().VariablePool.UserInputs
	system_inputs := n.GetGraphRuntimeState().VariablePool.SystemVariables

	// TODO: System variables should be directly accessible, no need for special handling
	// Set system variables as node outputs.
	for v := range system_inputs {
		node_inputs[constants.SYSTEM_VARIABLE_NODE_ID+"."+string(v)] = system_inputs[v]
	}
	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:  node_inputs,
		Outputs: node_inputs,
	}, nil
}
func (n *StartNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	typed_node_data, err := mapstruct.MapToStruct1[*startnodesentities.StartNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)

	return map[string][]string{}
}

func New() *StartNode {
	return &StartNode{
		BaseNode: &base.BaseNode[*startnodesentities.StartNodeData]{},
	}
}

func init() {
	nodesconstants.Regist(New())
}
