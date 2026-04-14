package trigger_webhook

import (
	"iter"
	"strings"

	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	triggerwebhookentities "mlib.com/gofy/server/entities/nodes/trigger_webhook"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
)

type TriggerWebhookNode struct {
	*base.BaseNode[*triggerwebhookentities.TriggerWebhookNodeData]
}

func (n *TriggerWebhookNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_TRIGGER_WEBHOOK
}

func (n *TriggerWebhookNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {
	// Get webhook data from variable pool (injected by trigger handler)
	webhookInputs := n.GetGraphRuntimeState().VariablePool.UserInputs

	// Extract configured outputs from webhook data
	outputs := n.extractConfiguredOutputs(webhookInputs)

	// Set system variables as node outputs
	systemInputs := n.GetGraphRuntimeState().VariablePool.SystemVariables
	for v := range systemInputs {
		outputs[constants.SYSTEM_VARIABLE_NODE_ID+"."+string(v)] = systemInputs[v]
	}

	return &workflowentities.NodeRunResult{
		Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Inputs:  webhookInputs,
		Outputs: outputs,
	}, nil
}

func (n *TriggerWebhookNode) extractConfiguredOutputs(webhookInputs map[string]any) map[string]any {
	outputs := make(map[string]any)
	nodeData := n.NodeData

	webhookData, _ := webhookInputs["webhook_data"].(map[string]any)
	if webhookData == nil {
		webhookData = make(map[string]any)
	}

	// Extract configured headers (case-insensitive)
	webhookHeaders, _ := webhookData["headers"].(map[string]any)
	headersLower := make(map[string]any)
	for k, v := range webhookHeaders {
		headersLower[strings.ToLower(k)] = v
	}

	for _, header := range nodeData.Headers {
		sanitized := strings.ReplaceAll(header.Name, "-", "_")
		value := getNormalized(webhookHeaders, header.Name)
		if value == nil {
			value = getNormalized(headersLower, strings.ToLower(header.Name))
		}
		outputs[sanitized] = value
	}

	// Extract configured query parameters
	queryParams, _ := webhookData["query_params"].(map[string]any)
	for _, param := range nodeData.Params {
		if queryParams != nil {
			outputs[param.Name] = queryParams[param.Name]
		}
	}

	// Extract configured body parameters
	bodyData, _ := webhookData["body"].(map[string]any)
	for _, bodyParam := range nodeData.Body {
		if nodeData.ContentType == triggerwebhookentities.ContentTypeText {
			raw, _ := bodyData["raw"].(string)
			outputs[bodyParam.Name] = raw
			continue
		}
		if nodeData.ContentType == triggerwebhookentities.ContentTypeBinary {
			outputs[bodyParam.Name] = bodyData["raw"]
			continue
		}
		if bodyParam.Type == "file" {
			files, _ := webhookData["files"].(map[string]any)
			outputs[bodyParam.Name] = files[bodyParam.Name]
		} else if bodyData != nil {
			outputs[bodyParam.Name] = bodyData[bodyParam.Name]
		}
	}

	// Include raw webhook data
	outputs["_webhook_raw"] = webhookData
	return outputs
}

func getNormalized(m map[string]any, key string) any {
	if m == nil {
		return nil
	}
	if v, ok := m[key]; ok {
		return v
	}
	// Try alternate naming
	alternate := key
	if strings.Contains(key, "-") {
		alternate = strings.ReplaceAll(key, "-", "_")
	} else {
		alternate = strings.ReplaceAll(key, "_", "-")
	}
	return m[alternate]
}

func (n *TriggerWebhookNode) ExtractVariableSelectorToVariableMapping(graphConfig map[string]any, nodeID string, nodeData *triggerwebhookentities.TriggerWebhookNodeData) map[string][]string {
	return map[string][]string{}
}

func New() *TriggerWebhookNode {
	return &TriggerWebhookNode{
		BaseNode: &base.BaseNode[*triggerwebhookentities.TriggerWebhookNodeData]{},
	}
}
