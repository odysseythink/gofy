package trigger_schedule

import (
	"iter"

	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	triggerscheduleentities "mlib.com/gofy/server/entities/nodes/trigger_schedule"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type TriggerScheduleNode struct {
	*base.BaseNode[*triggerscheduleentities.TriggerScheduleNodeData]
}

func (n *TriggerScheduleNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_TRIGGER_SCHEDULE
}

func (n *TriggerScheduleNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	nodeInputs := n.GetGraphRuntimeState().VariablePool.UserInputs
	systemInputs := n.GetGraphRuntimeState().VariablePool.SystemVariables

	// Set system variables as node outputs
	for v := range systemInputs {
		nodeInputs[constants.SYSTEM_VARIABLE_NODE_ID+"."+string(v)] = systemInputs[v]
	}

	outputs := make(map[string]any)
	for k, v := range nodeInputs {
		outputs[k] = v
	}

	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:  nodeInputs,
		Outputs: outputs,
	}, nil
}

func (n *TriggerScheduleNode) ExtractVariableSelectorToVariableMapping(graphConfig map[string]any, nodeID string, nodeData *triggerscheduleentities.TriggerScheduleNodeData) map[string][]string {
	return map[string][]string{}
}

func New() *TriggerScheduleNode {
	return &TriggerScheduleNode{
		BaseNode: &base.BaseNode[*triggerscheduleentities.TriggerScheduleNodeData]{},
	}
}
