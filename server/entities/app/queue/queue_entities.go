package queue

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	dbengine "mlib.com/gofy/server/db_engine"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type AppQueueEventer interface {
	Event() appenumtypes.QueueEvent
}

// QueueLLMChunkEvent represents a queue event for LLM chunk.
type QueueLLMChunkEvent struct {
	// AppQueueEvent
	Chunk *modelruntimeentities.LLMResultChunk `json:"chunk"`
}

func (aqe *QueueLLMChunkEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_LLM_CHUNK
}

// QueueIterationStartEvent represents a queue event for iteration start.
type QueueIterationStartEvent struct {
	NodeExecutionID           string                          `json:"node_execution_id"`
	NodeID                    string                          `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType         `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData `json:"node_data"`
	ParallelID                string                          `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                          `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                          `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                          `json:"parent_parallel_start_node_id,omitempty"`
	StartAt                   time.Time                       `json:"start_at"`
	NodeRunIndex              int                             `json:"node_run_index"`
	Inputs                    map[string]any                  `json:"inputs,omitempty"`
	PredecessorNodeID         string                          `json:"predecessor_node_id,omitempty"`
	Metadata                  map[string]any                  `json:"metadata,omitempty"`
}

func (aqe *QueueIterationStartEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_ITERATION_START
}

// QueueIterationNextEvent represents a queue event for iteration next.
type QueueIterationNextEvent struct {
	Index                     int                             `json:"index"`
	NodeExecutionID           string                          `json:"node_execution_id"`
	NodeID                    string                          `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType         `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData `json:"node_data"`
	ParallelID                string                          `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                          `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                          `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                          `json:"parent_parallel_start_node_id,omitempty"`
	ParallelModeRunID         string                          `json:"parallel_mode_run_id,omitempty"`
	NodeRunIndex              int                             `json:"node_run_index"`
	Output                    any                             `json:"output,omitempty"`
	Duration                  float64                         `json:"duration,omitempty"`
}

func (aqe *QueueIterationNextEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_ITERATION_NEXT
}

// QueueIterationCompletedEvent represents a queue event for iteration completed.
type QueueIterationCompletedEvent struct {
	NodeExecutionID           string                          `json:"node_execution_id"`
	NodeID                    string                          `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType         `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData `json:"node_data"`
	ParallelID                string                          `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                          `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                          `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                          `json:"parent_parallel_start_node_id,omitempty"`
	StartAt                   time.Time                       `json:"start_at"`
	NodeRunIndex              int                             `json:"node_run_index"`
	Inputs                    map[string]any                  `json:"inputs,omitempty"`
	Outputs                   map[string]any                  `json:"outputs,omitempty"`
	Metadata                  map[string]any                  `json:"metadata,omitempty"`
	Steps                     int                             `json:"steps"`
	Error                     string                          `json:"error,omitempty"`
}

func (aqe *QueueIterationCompletedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_ITERATION_COMPLETED
}

// QueueTextChunkEvent represents a queue event for text chunk.
type QueueTextChunkEvent struct {
	Text                 string   `json:"text"`
	FromVariableSelector []string `json:"from_variable_selector,omitempty"`
	InIterationID        string   `json:"in_iteration_id,omitempty"`
}

func (aqe *QueueTextChunkEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_TEXT_CHUNK
}

// QueueAgentMessageEvent represents a queue event for agent message.
type QueueAgentMessageEvent struct {
	Chunk *modelruntimeentities.LLMResultChunk `json:"chunk"`
}

func (aqe *QueueAgentMessageEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_AGENT_MESSAGE
}

// QueueMessageReplaceEvent represents a queue event for message replace.
type QueueMessageReplaceEvent struct {
	Text string `json:"text"`
}

func (aqe *QueueMessageReplaceEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_MESSAGE_REPLACE
}

// QueueRetrieverResourcesEvent represents a queue event for retriever resources.
type QueueRetrieverResourcesEvent struct {
	RetrieverResources []map[string]any `json:"retriever_resources"`
	InIterationID      string           `json:"in_iteration_id,omitempty"`
}

func (aqe *QueueRetrieverResourcesEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_RETRIEVER_RESOURCES
}

// QueueAnnotationReplyEvent represents a queue event for annotation reply.
type QueueAnnotationReplyEvent struct {
	MessageAnnotationID string `json:"message_annotation_id"`
}

func (aqe *QueueAnnotationReplyEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_ANNOTATION_REPLY
}

// QueueMessageEndEvent represents a queue event for message end.
type QueueMessageEndEvent struct {
	LLMResult *modelruntimeentities.LLMResult `json:"llm_result,omitempty"`
}

func (aqe *QueueMessageEndEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_MESSAGE_END
}

// QueueAdvancedChatMessageEndEvent represents a queue event for advanced chat message end.
type QueueAdvancedChatMessageEndEvent struct {
}

func (aqe *QueueAdvancedChatMessageEndEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_ADVANCED_CHAT_MESSAGE_END
}

// QueueWorkflowStartedEvent represents a queue event for workflow started.
type QueueWorkflowStartedEvent struct {
	GraphRuntimeState *graphengineentities.GraphRuntimeState `json:"graph_runtime_state"`
}

func (aqe *QueueWorkflowStartedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_WORKFLOW_STARTED
}

// QueueWorkflowSucceededEvent represents a queue event for workflow succeeded.
type QueueWorkflowSucceededEvent struct {
	Outputs map[string]any `json:"outputs,omitempty"`
}

func (aqe *QueueWorkflowSucceededEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_WORKFLOW_SUCCEEDED
}

// QueueWorkflowFailedEvent represents a queue event for workflow failed.
type QueueWorkflowFailedEvent struct {
	Error           string `json:"error"`
	ExceptionsCount int    `json:"exceptions_count"`
}

func (aqe *QueueWorkflowFailedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_WORKFLOW_FAILED
}

// QueueWorkflowPartialSuccessEvent represents a queue event for workflow partial success.
type QueueWorkflowPartialSuccessEvent struct {
	ExceptionsCount int            `json:"exceptions_count"`
	Outputs         map[string]any `json:"outputs,omitempty"`
}

func (aqe *QueueWorkflowPartialSuccessEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_WORKFLOW_PARTIAL_SUCCEEDED
}

type QueueNodeEventer interface {
	AppQueueEventer
	Handle(*models.WorkflowRun, *models.WorkflowNodeExecution) *models.WorkflowNodeExecution
}

// QueueNodeStartedEvent represents a queue event for node started.
type QueueNodeStartedEvent struct {
	NodeExecutionID           string                  `json:"node_execution_id"`
	NodeID                    string                  `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType `json:"node_type"`
	NodeData                  any/* *basenodesentities.BaseNodeData*/ `json:"node_data"`
	NodeRunIndex              int       `json:"node_run_index"`
	PredecessorNodeID         string    `json:"predecessor_node_id,omitempty"`
	ParallelID                string    `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string    `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string    `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string    `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string    `json:"in_iteration_id,omitempty"`
	StartAt                   time.Time `json:"start_at"`
	ParallelModeRunID         string    `json:"parallel_mode_run_id,omitempty"`
}

func (aqe *QueueNodeStartedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_NODE_STARTED
}

func (aqe *QueueNodeStartedEvent) Handle(workflow_run *models.WorkflowRun, workflow_node_execution *models.WorkflowNodeExecution) *models.WorkflowNodeExecution {
	now := time.Now()
	if workflow_node_execution != nil {
		workflow_node_execution = nil
	}
	workflow_node_execution = &models.WorkflowNodeExecution{
		ID:                uuid.NewV4().String(),
		TenantID:          workflow_run.TenantID,
		AppID:             workflow_run.AppID,
		WorkflowID:        workflow_run.WorkflowID,
		TriggeredFrom:     string(models.WorkflowNodeExecutionTriggeredFrom_WORKFLOW_RUN),
		WorkflowRunID:     workflow_run.ID,
		Index:             aqe.NodeRunIndex,
		PredecessorNodeID: aqe.PredecessorNodeID,
		NodeExecutionID:   aqe.NodeExecutionID,
		NodeID:            aqe.NodeID,
		NodeType:          aqe.NodeType,
		Title:             "", //aqe.NodeData.Title,
		Status:            string(models.WorkflowNodeExecutionStatus_RUNNING),
		CreatedByRole:     workflow_run.CreatedByRole,
		CreatedBy:         workflow_run.CreatedBy,
		CreatedAt:         &now,
	}
	bindata, _ := json.Marshal(map[workflowenumtypes.NodeRunMetadataKey]string{
		workflowenumtypes.NodeRunMetadataKey_PARALLEL_MODE_RUN_ID: aqe.ParallelModeRunID,
		workflowenumtypes.NodeRunMetadataKey_ITERATION_ID:         aqe.InIterationID,
	})
	workflow_node_execution.ExecutionMetadata = string(bindata)
	dbengine.Instance().DB.Create(workflow_node_execution)

	return workflow_node_execution
}

// QueueNodeSucceededEvent represents a queue event for node succeeded.
type QueueNodeSucceededEvent struct {
	NodeExecutionID           string                                       `json:"node_execution_id"`
	NodeID                    string                                       `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType                      `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData              `json:"node_data"`
	ParallelID                string                                       `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                                       `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                                       `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                                       `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string                                       `json:"in_iteration_id,omitempty"`
	StartAt                   time.Time                                    `json:"start_at"`
	Inputs                    map[string]any                               `json:"inputs,omitempty"`
	ProcessData               map[string]any                               `json:"process_data,omitempty"`
	Outputs                   map[string]any                               `json:"outputs,omitempty"`
	ExecutionMetadata         map[workflowenumtypes.NodeRunMetadataKey]any `json:"execution_metadata,omitempty"`
	Error                     string                                       `json:"error,omitempty"`
	IterationDurationMap      map[string]float64                           `json:"iteration_duration_map,omitempty"`
}

func (aqe *QueueNodeSucceededEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_NODE_SUCCEEDED
}

func (aqe *QueueNodeSucceededEvent) Handle(workflow_run *models.WorkflowRun, workflow_node_execution *models.WorkflowNodeExecution) *models.WorkflowNodeExecution {
	if workflow_node_execution == nil {
		panic(exceptions.NewValueError("missing workflow_node_execution "))
	}
	inputs := file.HandleSpecialValues(aqe.Inputs)
	process_data := file.HandleSpecialValues(aqe.ProcessData)
	outputs := file.HandleSpecialValues(aqe.Outputs)
	execution_metadata, _ := json.Marshal(aqe.ExecutionMetadata)
	finished_at := time.Now()
	elapsed_time := finished_at.Sub(aqe.StartAt).Seconds()

	workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_SUCCEEDED)
	bindata, _ := json.Marshal(inputs)
	workflow_node_execution.Inputs = string(bindata)
	bindata, _ = json.Marshal(process_data)
	workflow_node_execution.ProcessData = string(bindata)
	bindata, _ = json.Marshal(outputs)
	workflow_node_execution.Outputs = string(bindata)
	workflow_node_execution.ExecutionMetadata = string(execution_metadata)
	workflow_node_execution.FinishedAt = &finished_at
	workflow_node_execution.ElapsedTime = elapsed_time

	dbengine.Instance().DB.Save(workflow_node_execution)
	return workflow_node_execution
}

// QueueNodeRetryEvent represents a queue event for node retry.
type QueueNodeRetryEvent struct {
	*QueueNodeStartedEvent
	Inputs            map[string]any                               `json:"inputs,omitempty"`
	ProcessData       map[string]any                               `json:"process_data,omitempty"`
	Outputs           map[string]any                               `json:"outputs,omitempty"`
	ExecutionMetadata map[workflowenumtypes.NodeRunMetadataKey]any `json:"execution_metadata,omitempty"`
	Error             string                                       `json:"error"`
	RetryIndex        int                                          `json:"retry_index"`
}

func (aqe *QueueNodeRetryEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_RETRY
}

func (aqe *QueueNodeRetryEvent) Handle(workflow_run *models.WorkflowRun, workflow_node_execution *models.WorkflowNodeExecution) *models.WorkflowNodeExecution {
	if workflow_run == nil {
		mlog.Error("missing workflow_run ")
		panic(exceptions.NewValueError("missing workflow_run "))
	}
	created_at := aqe.StartAt
	finished_at := time.Now()
	elapsed_time := finished_at.Sub(created_at).Seconds()
	inputs := file.HandleSpecialValues(aqe.Inputs)
	outputs := file.HandleSpecialValues(aqe.Outputs)
	origin_metadata := map[workflowenumtypes.NodeRunMetadataKey]any{
		workflowenumtypes.NodeRunMetadataKey_ITERATION_ID:         aqe.InIterationID,
		workflowenumtypes.NodeRunMetadataKey_PARALLEL_MODE_RUN_ID: aqe.ParallelModeRunID,
	}
	for k, v := range aqe.ExecutionMetadata {
		origin_metadata[k] = v
	}

	execution_metadata, _ := json.Marshal(origin_metadata)

	workflow_node_execution = &models.WorkflowNodeExecution{}
	workflow_node_execution.ID = uuid.NewV4().String()
	workflow_node_execution.TenantID = workflow_run.TenantID
	workflow_node_execution.AppID = workflow_run.AppID
	workflow_node_execution.WorkflowID = workflow_run.WorkflowID
	workflow_node_execution.TriggeredFrom = string(models.WorkflowNodeExecutionTriggeredFrom_WORKFLOW_RUN)
	workflow_node_execution.WorkflowRunID = workflow_run.ID
	workflow_node_execution.PredecessorNodeID = aqe.PredecessorNodeID
	workflow_node_execution.NodeExecutionID = aqe.NodeExecutionID
	workflow_node_execution.NodeID = aqe.NodeID
	workflow_node_execution.NodeType = aqe.NodeType
	workflow_node_execution.Title = "" //aqe.NodeData.Title
	workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_RETRY)
	workflow_node_execution.CreatedByRole = workflow_run.CreatedByRole
	workflow_node_execution.CreatedBy = workflow_run.CreatedBy
	workflow_node_execution.CreatedAt = &created_at
	workflow_node_execution.FinishedAt = &finished_at
	workflow_node_execution.ElapsedTime = elapsed_time
	workflow_node_execution.Error = aqe.Error
	bindata, _ := json.Marshal(inputs)
	workflow_node_execution.Inputs = string(bindata)
	bindata, _ = json.Marshal(outputs)
	workflow_node_execution.Outputs = string(bindata)
	workflow_node_execution.ExecutionMetadata = string(execution_metadata)
	workflow_node_execution.Index = aqe.NodeRunIndex

	dbengine.Instance().DB.Create(workflow_node_execution)

	return workflow_node_execution
}

// QueueNodeInIterationFailedEvent represents a queue event for node failed in iteration.
type QueueNodeInIterationFailedEvent struct {
	NodeExecutionID           string                                       `json:"node_execution_id"`
	NodeID                    string                                       `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType                      `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData              `json:"node_data"`
	ParallelID                string                                       `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                                       `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                                       `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                                       `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string                                       `json:"in_iteration_id,omitempty"`
	StartAt                   time.Time                                    `json:"start_at"`
	Inputs                    map[string]any                               `json:"inputs,omitempty"`
	ProcessData               map[string]any                               `json:"process_data,omitempty"`
	Outputs                   map[string]any                               `json:"outputs,omitempty"`
	ExecutionMetadata         map[workflowenumtypes.NodeRunMetadataKey]any `json:"execution_metadata,omitempty"`
	Error                     string                                       `json:"error"`
}

func (aqe *QueueNodeInIterationFailedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_NODE_FAILED
}

func (aqe *QueueNodeInIterationFailedEvent) Handle(workflow_run *models.WorkflowRun, workflow_node_execution *models.WorkflowNodeExecution) *models.WorkflowNodeExecution {
	if workflow_node_execution == nil {
		mlog.Error("missing workflow_node_execution ")
		panic(exceptions.NewValueError("missing workflow_node_execution "))
	}
	inputs := file.HandleSpecialValues(aqe.Inputs)
	process_data := file.HandleSpecialValues(aqe.ProcessData)
	outputs := file.HandleSpecialValues(aqe.Outputs)
	execution_metadata, _ := json.Marshal(aqe.ExecutionMetadata)
	finished_at := time.Now()
	elapsed_time := finished_at.Sub(aqe.StartAt).Seconds()

	workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_FAILED)
	bindata, _ := json.Marshal(inputs)
	workflow_node_execution.Inputs = string(bindata)
	bindata, _ = json.Marshal(process_data)
	workflow_node_execution.ProcessData = string(bindata)
	bindata, _ = json.Marshal(outputs)
	workflow_node_execution.Outputs = string(bindata)
	workflow_node_execution.FinishedAt = &finished_at
	workflow_node_execution.ElapsedTime = elapsed_time
	workflow_node_execution.ExecutionMetadata = string(execution_metadata)
	dbengine.Instance().DB.Save(workflow_node_execution)
	return workflow_node_execution
}

// QueueNodeExceptionEvent represents a queue event for node exception.
type QueueNodeExceptionEvent struct {
	NodeExecutionID           string                                       `json:"node_execution_id"`
	NodeID                    string                                       `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType                      `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData              `json:"node_data"`
	ParallelID                string                                       `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                                       `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                                       `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                                       `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string                                       `json:"in_iteration_id,omitempty"`
	StartAt                   time.Time                                    `json:"start_at"`
	Inputs                    map[string]any                               `json:"inputs,omitempty"`
	ProcessData               map[string]any                               `json:"process_data,omitempty"`
	Outputs                   map[string]any                               `json:"outputs,omitempty"`
	ExecutionMetadata         map[workflowenumtypes.NodeRunMetadataKey]any `json:"execution_metadata,omitempty"`
	Error                     string                                       `json:"error"`
}

func (aqe *QueueNodeExceptionEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_NODE_EXCEPTION
}

func (aqe *QueueNodeExceptionEvent) Handle(workflow_run *models.WorkflowRun, workflow_node_execution *models.WorkflowNodeExecution) *models.WorkflowNodeExecution {
	if workflow_node_execution == nil {
		mlog.Error("missing workflow_node_execution ")
		panic(exceptions.NewValueError("missing workflow_node_execution "))
	}
	inputs := file.HandleSpecialValues(aqe.Inputs)
	process_data := file.HandleSpecialValues(aqe.ProcessData)
	outputs := file.HandleSpecialValues(aqe.Outputs)
	execution_metadata, _ := json.Marshal(aqe.ExecutionMetadata)
	finished_at := time.Now()
	elapsed_time := finished_at.Sub(aqe.StartAt).Seconds()

	workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_EXCEPTION)
	bindata, _ := json.Marshal(inputs)
	workflow_node_execution.Inputs = string(bindata)
	bindata, _ = json.Marshal(process_data)
	workflow_node_execution.ProcessData = string(bindata)
	bindata, _ = json.Marshal(outputs)
	workflow_node_execution.Outputs = string(bindata)
	workflow_node_execution.FinishedAt = &finished_at
	workflow_node_execution.ElapsedTime = elapsed_time
	workflow_node_execution.ExecutionMetadata = string(execution_metadata)
	dbengine.Instance().DB.Save(workflow_node_execution)
	return workflow_node_execution
}

// QueueNodeFailedEvent represents a queue event for node failed.
type QueueNodeFailedEvent struct {
	NodeExecutionID           string                                       `json:"node_execution_id"`
	NodeID                    string                                       `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType                      `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData              `json:"node_data"`
	ParallelID                string                                       `json:"parallel_id,omitempty"`
	ParallelStartNodeID       string                                       `json:"parallel_start_node_id,omitempty"`
	ParentParallelID          string                                       `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string                                       `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string                                       `json:"in_iteration_id,omitempty"`
	StartAt                   time.Time                                    `json:"start_at"`
	Inputs                    map[string]any                               `json:"inputs,omitempty"`
	ProcessData               map[string]any                               `json:"process_data,omitempty"`
	Outputs                   map[string]any                               `json:"outputs,omitempty"`
	ExecutionMetadata         map[workflowenumtypes.NodeRunMetadataKey]any `json:"execution_metadata,omitempty"`
	Error                     string                                       `json:"error"`
}

func (aqe *QueueNodeFailedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_NODE_FAILED
}

func (aqe *QueueNodeFailedEvent) Handle(workflow_run *models.WorkflowRun, workflow_node_execution *models.WorkflowNodeExecution) *models.WorkflowNodeExecution {
	if workflow_node_execution == nil {
		mlog.Error("missing workflow_node_execution ")
		panic(exceptions.NewValueError("missing workflow_node_execution "))
	}
	inputs := file.HandleSpecialValues(aqe.Inputs)
	process_data := file.HandleSpecialValues(aqe.ProcessData)
	outputs := file.HandleSpecialValues(aqe.Outputs)
	execution_metadata, _ := json.Marshal(aqe.ExecutionMetadata)
	finished_at := time.Now()
	elapsed_time := finished_at.Sub(aqe.StartAt).Seconds()

	workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_FAILED)
	bindata, _ := json.Marshal(inputs)
	workflow_node_execution.Inputs = string(bindata)
	bindata, _ = json.Marshal(process_data)
	workflow_node_execution.ProcessData = string(bindata)
	bindata, _ = json.Marshal(outputs)
	workflow_node_execution.Outputs = string(bindata)
	workflow_node_execution.FinishedAt = &finished_at
	workflow_node_execution.ElapsedTime = elapsed_time
	workflow_node_execution.ExecutionMetadata = string(execution_metadata)
	dbengine.Instance().DB.Save(workflow_node_execution)
	return workflow_node_execution
}

// QueueAgentThoughtEvent represents a queue event for agent thought.
type QueueAgentThoughtEvent struct {
	AgentThoughtID string `json:"agent_thought_id"`
}

func (aqe *QueueAgentThoughtEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_AGENT_THOUGHT
}

// QueueMessageFileEvent represents a queue event for message file.
type QueueMessageFileEvent struct {
	MessageFileID string `json:"message_file_id"`
}

func (aqe *QueueMessageFileEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_MESSAGE_FILE
}

// QueueErrorEvent represents a queue event for error.
type QueueErrorEvent struct {
	Err error `json:"error,omitempty"`
}

func (aqe *QueueErrorEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_ERROR
}

// QueuePingEvent represents a queue event for ping.
type QueuePingEvent struct {
}

func (aqe *QueuePingEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_PING
}

// QueueStopEvent represents a queue event for stop.
type QueueStopEvent struct {
	StoppedBy appenumtypes.QueueStopEvent_StopBy `json:"stopped_by"`
}

func (aqe *QueueStopEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_STOP
}

var (
	reasonMapping = map[appenumtypes.QueueStopEvent_StopBy]string{
		appenumtypes.QueueStopEvent_StopBy_USER_MANUAL:       "Stopped by user.",
		appenumtypes.QueueStopEvent_StopBy_ANNOTATION_REPLY:  "Stopped by annotation reply.",
		appenumtypes.QueueStopEvent_StopBy_OUTPUT_MODERATION: "Stopped by output moderation.",
		appenumtypes.QueueStopEvent_StopBy_INPUT_MODERATION:  "Stopped by input moderation.",
	}
)

func (aqe *QueueStopEvent) GetStopReason() string {
	if reason, ok := reasonMapping[aqe.StoppedBy]; ok {
		return reason
	}
	return "Stopped by unknown reason."
}

// QueueMessage represents a queue message.
type QueueMessage struct {
	TaskID  string          `json:"task_id"`
	AppMode string          `json:"app_mode"`
	Eventer AppQueueEventer `json:"-"`
}

// MessageQueueMessage represents a message queue message.
type MessageQueueMessage struct {
	*QueueMessage
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
}

// WorkflowQueueMessage represents a workflow queue message.
type WorkflowQueueMessage struct {
	*QueueMessage
}

// QueueParallelBranchRunStartedEvent represents a queue event for parallel branch run started.
type QueueParallelBranchRunStartedEvent struct {
	ParallelID                string `json:"parallel_id"`
	ParallelStartNodeID       string `json:"parallel_start_node_id"`
	ParentParallelID          string `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string `json:"in_iteration_id,omitempty"`
}

func (aqe *QueueParallelBranchRunStartedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_PARALLEL_BRANCH_RUN_STARTED
}

// QueueParallelBranchRunSucceededEvent represents a queue event for parallel branch run succeeded.
type QueueParallelBranchRunSucceededEvent struct {
	ParallelID                string `json:"parallel_id"`
	ParallelStartNodeID       string `json:"parallel_start_node_id"`
	ParentParallelID          string `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string `json:"in_iteration_id,omitempty"`
}

func (aqe *QueueParallelBranchRunSucceededEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_PARALLEL_BRANCH_RUN_SUCCEEDED
}

// QueueParallelBranchRunFailedEvent represents a queue event for parallel branch run failed.
type QueueParallelBranchRunFailedEvent struct {
	ParallelID                string `json:"parallel_id"`
	ParallelStartNodeID       string `json:"parallel_start_node_id"`
	ParentParallelID          string `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
	InIterationID             string `json:"in_iteration_id,omitempty"`
	Error                     string `json:"error"`
}

func (aqe *QueueParallelBranchRunFailedEvent) Event() appenumtypes.QueueEvent {
	return appenumtypes.QueueEvent_PARALLEL_BRANCH_RUN_FAILED
}
