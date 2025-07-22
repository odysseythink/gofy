package workflow

import (
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
)

type NodeRunResult struct {
	Status           models.WorkflowNodeExecutionStatus           `json:"status"`
	Inputs           map[string]any                               `json:"inputs"`
	ProcessData      map[string]any                               `json:"process_data"`                 // process data
	Outputs          map[string]any                               `json:"outputs"`                      // node outputs
	Metadata         map[workflowenumtypes.NodeRunMetadataKey]any `json:"metadata"`                     //# node metadata
	LLmUsage         *modelruntimeentities.LLMUsage               `json:"llm_usage"`                    // llm usage
	EdgeSourceHandle string                                       `json:"edge_source_handle,omitempty"` // source handle id of node with multiple branches

	Error     string `json:"error"`      // error message if status is failed
	ErrorType string `json:"error_type"` // error type if status is failed

	// single step node run retry
	RetryIndex int `json:"retry_index"`
}

func NewNodeRunResult() *NodeRunResult {
	return &NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_RUNNING,
	}
}
