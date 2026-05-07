package event

import (
	"time"

	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	"github.com/odysseythink/gofy/backend/models"
)

type NodeEventType string

const (
	/*
	   QueueEvent enum
	*/

	NodeEvent_RUN_COMPLETED          NodeEventType = "run_completed"
	NodeEvent_RUN_STREAM_CHUNK       NodeEventType = "run_stream_Chunk"
	NodeEvent_RUN_RETRIEVER_RESOURCE NodeEventType = "run_retriever_resource"
	NodeEvent_MODEL_INVOKE_COMPLETED NodeEventType = "model_invoke_completed"
)

type RunCompletedEvent struct {
	RunResult *workflowentities.NodeRunResult `json:"run_result"`
}

func (e *RunCompletedEvent) Event() NodeEventType {
	return NodeEvent_RUN_COMPLETED
}

type RunStreamChunkEvent struct {
	ChunkContent         string   `json:"chunk_content"`
	FromVariableSelector []string `json:"from_variable_selector"`
}

func (e *RunStreamChunkEvent) Event() NodeEventType {
	return NodeEvent_RUN_STREAM_CHUNK
}

type RunRetrieverResourceEvent struct {
	RetrieverResources []map[string]any `json:"retriever_resources"`
	Context            string           `json:"context"`
}

func (e *RunRetrieverResourceEvent) Event() NodeEventType {
	return NodeEvent_RUN_RETRIEVER_RESOURCE
}

type ModelInvokeCompletedEvent struct {
	Text         string                         `json:"text"`
	Usage        *modelruntimeentities.LLMUsage `json:"usage"`
	FinishReason string                         `json:"finish_reason,omitempty"`
}

func (e *ModelInvokeCompletedEvent) Event() NodeEventType {
	return NodeEvent_MODEL_INVOKE_COMPLETED
}

type RunRetryEvent struct {
	Error      string    `json:"error"`
	RetryIndex int       `json:"retry_index"`
	StartAt    time.Time `json:"start_at"`
}

type SingleStepRetryEvent struct {
	*workflowentities.NodeRunResult
	ElapsedTime float64 `json:"elapsed_time"`
}

func NewSingleStepRetryEvent() *SingleStepRetryEvent {
	return &SingleStepRetryEvent{
		NodeRunResult: &workflowentities.NodeRunResult{
			Status: models.WorkflowNodeExecutionStatus_RETRY,
		},
	}
}

type NodeEventer interface {
	// *RunCompletedEvent | *RunStreamChunkEvent | *RunRetrieverResourceEvent | *ModelInvokeCompletedEvent
	Event() NodeEventType
}
