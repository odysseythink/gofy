package humaninput

import (
	"iter"

	"mlib.com/gofy/server/core/workflow/nodes/base"
	humaninputnodesentities "mlib.com/gofy/server/entities/nodes/human_input"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type HumanInputNode struct {
	*base.BaseNode[*humaninputnodesentities.HumanInputNodeData]
}

func (n *HumanInputNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_HUMAN_INPUT
}

func (n *HumanInputNode) Run() (run_result *workflowentities.NodeRunResult, run_stream_result iter.Seq[any]) {
	nodeData := n.BaseNode.NodeData

	// Build form outputs
	outputs := map[string]any{
		"form_content":    nodeData.FormContent,
		"inputs":          nodeData.Inputs,
		"user_actions":    nodeData.UserActions,
		"delivery_methods": nodeData.DeliveryMethods,
		"timeout":         nodeData.Timeout,
		"timeout_unit":    nodeData.TimeoutUnit,
	}

	run_result = &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Outputs: outputs,
		Metadata: map[string]any{
			"requires_pause": true,
		},
	}

	return run_result, nil
}

func (n *HumanInputNode) GetDefaultConfig(filters map[string]any) map[string]any {
	return map[string]any{
		"type": "human-input",
		"config": map[string]any{
			"form_content": "",
			"inputs":       []any{},
			"user_actions": []map[string]any{
				{"id": "approve", "title": "Approve", "style": "primary"},
				{"id": "reject", "title": "Reject", "style": "danger"},
			},
			"delivery_methods": []map[string]any{
				{"type": "webapp"},
			},
			"timeout":      7,
			"timeout_unit": "days",
		},
	}
}

func (n *HumanInputNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *humaninputnodesentities.HumanInputNodeData) map[string][]string {
	return map[string][]string{}
}

func New() *HumanInputNode {
	return &HumanInputNode{
		BaseNode: &base.BaseNode[*humaninputnodesentities.HumanInputNodeData]{},
	}
}
