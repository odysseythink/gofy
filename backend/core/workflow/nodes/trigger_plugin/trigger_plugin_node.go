package trigger_plugin

import (
	"iter"

	"github.com/odysseythink/gofy/backend/constants"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	triggerpluginentities "github.com/odysseythink/gofy/backend/entities/nodes/trigger_plugin"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/models"
)

type TriggerPluginNode struct {
	*base.BaseNode[*triggerpluginentities.TriggerEventNodeData]
}

func (n *TriggerPluginNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_TRIGGER_PLUGIN
}

func (n *TriggerPluginNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	nodeData := n.NodeData

	// Build metadata for trigger info
	metadata := map[workflowenumtypes.NodeRunMetadataKey]any{
		workflowenumtypes.NodeRunMetadataKey_TRIGGER_INFO: map[string]any{
			"provider_id":              nodeData.ProviderID,
			"event_name":               nodeData.EventName,
			"plugin_unique_identifier": nodeData.PluginUniqueIdentifier,
		},
	}

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
		Status:   models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:   nodeInputs,
		Outputs:  outputs,
		Metadata: metadata,
	}, nil
}

func (n *TriggerPluginNode) ExtractVariableSelectorToVariableMapping(graphConfig map[string]any, nodeID string, nodeData *triggerpluginentities.TriggerEventNodeData) map[string][]string {
	return map[string][]string{}
}

func New() *TriggerPluginNode {
	return &TriggerPluginNode{
		BaseNode: &base.BaseNode[*triggerpluginentities.TriggerEventNodeData]{},
	}
}
