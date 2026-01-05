package variableaggregator

import (
	"iter"
	"strings"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	variableaggregatornodesentities "mlib.com/gofy/server/entities/nodes/variable_aggregator"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

type VariableAggregatorNode struct {
	*base.BaseNode[*variableaggregatornodesentities.VariableAggregatorNodeData]
}

func (n *VariableAggregatorNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_VARIABLE_AGGREGATOR
}

func (n *VariableAggregatorNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	// Get variables
	outputs := map[string]any{}
	inputs := map[string]any{}

	if n.NodeData.AdvancedSettings == nil || !n.NodeData.AdvancedSettings.GroupEnabled {
		for _, selector := range n.NodeData.Variables {
			variable := n.GetGraphRuntimeState().VariablePool.Get(selector)
			if variable != nil {
				outputs = map[string]any{"output": variable.ToObject()}

				inputs = map[string]any{strings.Join(selector[1:], "."): variable.ToObject()}
				break
			}
		}
	} else {
		for _, group := range n.NodeData.AdvancedSettings.Groups {
			for _, selector := range group.Variables {
				variable := n.GetGraphRuntimeState().VariablePool.Get(selector)

				if variable != nil {
					outputs[group.GroupName] = map[string]any{"output": variable.ToObject()}
					inputs[strings.Join(selector[1:], ".")] = variable.ToObject()
					break
				}
			}
		}
	}

	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Outputs: outputs,
		Inputs:  inputs,
	}, nil
}
func (n *VariableAggregatorNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	typed_node_data, err := mapstruct.MapToStruct1[*variableaggregatornodesentities.VariableAggregatorNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)

	return map[string][]string{}
}
func New() *VariableAggregatorNode {
	return &VariableAggregatorNode{
		BaseNode: &base.BaseNode[*variableaggregatornodesentities.VariableAggregatorNodeData]{},
	}
}
func init() {
	nodesconstants.Regist(New())
}
