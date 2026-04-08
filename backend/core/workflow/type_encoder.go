package workflow

import (
	"encoding/json"
	"time"
)

// WorkflowTypeEncoder handles custom JSON encoding for workflow types.
type WorkflowTypeEncoder struct{}

// EncodeNodeResult encodes a node execution result for storage.
func (e *WorkflowTypeEncoder) EncodeNodeResult(
	nodeID, nodeType string,
	inputs, outputs, processData, metadata map[string]any,
	status string,
	elapsed float64,
	errMsg string,
) map[string]any {
	result := map[string]any{
		"node_id":      nodeID,
		"node_type":    nodeType,
		"status":       status,
		"elapsed_time": elapsed,
		"created_at":   time.Now().Format(time.RFC3339),
	}
	if inputs != nil {
		inputsJSON, _ := json.Marshal(inputs)
		result["inputs"] = string(inputsJSON)
	}
	if outputs != nil {
		outputsJSON, _ := json.Marshal(outputs)
		result["outputs"] = string(outputsJSON)
	}
	if processData != nil {
		pdJSON, _ := json.Marshal(processData)
		result["process_data"] = string(pdJSON)
	}
	if metadata != nil {
		mdJSON, _ := json.Marshal(metadata)
		result["execution_metadata"] = string(mdJSON)
	}
	if errMsg != "" {
		result["error"] = errMsg
	}
	return result
}

// EncodeWorkflowRunResult encodes a complete workflow run for storage.
func (e *WorkflowTypeEncoder) EncodeWorkflowRunResult(
	workflowRunID string,
	status string,
	inputs, outputs map[string]any,
	totalSteps, totalTokens int,
	elapsedTime float64,
	errMsg string,
) map[string]any {
	result := map[string]any{
		"id":           workflowRunID,
		"status":       status,
		"total_steps":  totalSteps,
		"total_tokens": totalTokens,
		"elapsed_time": elapsedTime,
		"finished_at":  time.Now().Format(time.RFC3339),
	}
	if inputs != nil {
		inputsJSON, _ := json.Marshal(inputs)
		result["inputs"] = string(inputsJSON)
	}
	if outputs != nil {
		outputsJSON, _ := json.Marshal(outputs)
		result["outputs"] = string(outputsJSON)
	}
	if errMsg != "" {
		result["error"] = errMsg
	}
	return result
}

// DecodeJSONField safely decodes a JSON string field to a map.
func DecodeJSONField(jsonStr string) map[string]any {
	if jsonStr == "" {
		return nil
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil
	}
	return result
}
