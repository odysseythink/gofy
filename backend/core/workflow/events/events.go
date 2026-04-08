package events

import (
	"time"
)

// ========== Event Types ==========

type EventType string

const (
	// Graph lifecycle events
	EventGraphRunStarted          EventType = "graph_run_started"
	EventGraphRunSucceeded        EventType = "graph_run_succeeded"
	EventGraphRunFailed           EventType = "graph_run_failed"
	EventGraphRunAborted          EventType = "graph_run_aborted"
	EventGraphRunPaused           EventType = "graph_run_paused"
	EventGraphRunResumed          EventType = "graph_run_resumed"
	EventGraphRunPartialSucceeded EventType = "graph_run_partial_succeeded"

	// Node execution events
	EventNodeExecutionStarted   EventType = "node_execution_started"
	EventNodeExecutionSucceeded EventType = "node_execution_succeeded"
	EventNodeExecutionFailed    EventType = "node_execution_failed"
	EventNodeExecutionRetrying  EventType = "node_execution_retrying"
	EventNodeExecutionSkipped   EventType = "node_execution_skipped"

	// Agent events
	EventAgentLoopStarted  EventType = "agent_loop_started"
	EventAgentLoopPending  EventType = "agent_loop_pending"
	EventAgentLoopFinished EventType = "agent_loop_finished"
	EventAgentMessage      EventType = "agent_message"

	// Human input events
	EventHumanInputRequired  EventType = "human_input_required"
	EventHumanInputSubmitted EventType = "human_input_submitted"
	EventHumanInputTimeout   EventType = "human_input_timeout"

	// Iteration events
	EventIterationStarted   EventType = "iteration_started"
	EventIterationNext      EventType = "iteration_next"
	EventIterationCompleted EventType = "iteration_completed"

	// Loop events
	EventLoopStarted   EventType = "loop_started"
	EventLoopNext      EventType = "loop_next"
	EventLoopCompleted EventType = "loop_completed"

	// Streaming events
	EventTextChunk       EventType = "text_chunk"
	EventMessageReplace  EventType = "message_replace"
	EventStreamCompleted EventType = "stream_completed"

	// Parallel execution events
	EventParallelBranchStarted   EventType = "parallel_branch_started"
	EventParallelBranchCompleted EventType = "parallel_branch_completed"
)

// ========== Event Interface ==========

// GraphEvent is the interface all events implement.
type GraphEvent interface {
	GetEventType() EventType
	GetEventTimestamp() time.Time
}

// ========== Base Event ==========

type BaseGraphEvent struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func newBaseEvent(eventType EventType) BaseGraphEvent {
	return BaseGraphEvent{Type: eventType, Timestamp: time.Now()}
}

func (e BaseGraphEvent) GetEventType() EventType    { return e.Type }
func (e BaseGraphEvent) GetEventTimestamp() time.Time { return e.Timestamp }

// ========== Graph Events ==========

type GraphRunStartedEvent struct {
	BaseGraphEvent
	GraphID       string         `json:"graph_id"`
	WorkflowID    string         `json:"workflow_id"`
	WorkflowRunID string         `json:"workflow_run_id"`
	Inputs        map[string]any `json:"inputs,omitempty"`
}

func NewGraphRunStartedEvent(graphID, workflowID, workflowRunID string, inputs map[string]any) *GraphRunStartedEvent {
	return &GraphRunStartedEvent{
		BaseGraphEvent: newBaseEvent(EventGraphRunStarted),
		GraphID: graphID, WorkflowID: workflowID, WorkflowRunID: workflowRunID, Inputs: inputs,
	}
}

type GraphRunSucceededEvent struct {
	BaseGraphEvent
	GraphID       string         `json:"graph_id"`
	WorkflowRunID string         `json:"workflow_run_id"`
	Outputs       map[string]any `json:"outputs,omitempty"`
	TotalTokens   int            `json:"total_tokens"`
	TotalSteps    int            `json:"total_steps"`
	ElapsedTime   float64        `json:"elapsed_time"`
}

type GraphRunFailedEvent struct {
	BaseGraphEvent
	GraphID       string  `json:"graph_id"`
	WorkflowRunID string  `json:"workflow_run_id"`
	Error         string  `json:"error"`
	ElapsedTime   float64 `json:"elapsed_time"`
}

type GraphRunAbortedEvent struct {
	BaseGraphEvent
	GraphID       string `json:"graph_id"`
	WorkflowRunID string `json:"workflow_run_id"`
	Reason        string `json:"reason"`
}

type GraphRunPausedEvent struct {
	BaseGraphEvent
	GraphID       string `json:"graph_id"`
	WorkflowRunID string `json:"workflow_run_id"`
	NodeID        string `json:"node_id"`
	ReasonType    string `json:"reason_type"`
	FormID        string `json:"form_id,omitempty"`
}

type GraphRunResumedEvent struct {
	BaseGraphEvent
	GraphID       string `json:"graph_id"`
	WorkflowRunID string `json:"workflow_run_id"`
}

type GraphRunPartialSucceededEvent struct {
	BaseGraphEvent
	GraphID       string         `json:"graph_id"`
	WorkflowRunID string         `json:"workflow_run_id"`
	Outputs       map[string]any `json:"outputs,omitempty"`
	Exceptions    int            `json:"exceptions"`
}

// ========== Node Events ==========

type NodeExecutionStartedEvent struct {
	BaseGraphEvent
	NodeID        string         `json:"node_id"`
	NodeType      string         `json:"node_type"`
	NodeTitle     string         `json:"node_title"`
	PredecessorID string         `json:"predecessor_id,omitempty"`
	Inputs        map[string]any `json:"inputs,omitempty"`
	ParallelID    string         `json:"parallel_id,omitempty"`
	ParallelIndex int            `json:"parallel_index,omitempty"`
}

type NodeExecutionSucceededEvent struct {
	BaseGraphEvent
	NodeID      string         `json:"node_id"`
	NodeType    string         `json:"node_type"`
	Inputs      map[string]any `json:"inputs,omitempty"`
	Outputs     map[string]any `json:"outputs,omitempty"`
	ProcessData map[string]any `json:"process_data,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	ElapsedTime float64        `json:"elapsed_time"`
	ExecutionID string         `json:"execution_id,omitempty"`
}

type NodeExecutionFailedEvent struct {
	BaseGraphEvent
	NodeID      string  `json:"node_id"`
	NodeType    string  `json:"node_type"`
	Error       string  `json:"error"`
	ElapsedTime float64 `json:"elapsed_time"`
	ExecutionID string  `json:"execution_id,omitempty"`
}

type NodeExecutionRetryingEvent struct {
	BaseGraphEvent
	NodeID     string `json:"node_id"`
	NodeType   string `json:"node_type"`
	RetryCount int    `json:"retry_count"`
	MaxRetries int    `json:"max_retries"`
	Error      string `json:"error"`
}

type NodeExecutionSkippedEvent struct {
	BaseGraphEvent
	NodeID   string `json:"node_id"`
	NodeType string `json:"node_type"`
	Reason   string `json:"reason"`
}

// ========== Agent Events ==========

type AgentLoopStartedEvent struct {
	BaseGraphEvent
	NodeID    string `json:"node_id"`
	Iteration int    `json:"iteration"`
}

type AgentLoopPendingEvent struct {
	BaseGraphEvent
	NodeID    string `json:"node_id"`
	Iteration int    `json:"iteration"`
	ToolName  string `json:"tool_name,omitempty"`
	ToolInput string `json:"tool_input,omitempty"`
}

type AgentLoopFinishedEvent struct {
	BaseGraphEvent
	NodeID    string `json:"node_id"`
	Iteration int    `json:"iteration"`
}

type AgentMessageEvent struct {
	BaseGraphEvent
	NodeID      string `json:"node_id"`
	Message     string `json:"message"`
	MessageType string `json:"type"` // thought, action, observation
}

// ========== Human Input Events ==========

type HumanInputRequiredEvent struct {
	BaseGraphEvent
	NodeID   string `json:"node_id"`
	FormID   string `json:"form_id"`
	FormKind string `json:"form_kind"`
}

type HumanInputSubmittedEvent struct {
	BaseGraphEvent
	NodeID         string         `json:"node_id"`
	FormID         string         `json:"form_id"`
	SelectedAction string         `json:"selected_action"`
	FormData       map[string]any `json:"form_data"`
}

type HumanInputTimeoutEvent struct {
	BaseGraphEvent
	NodeID string `json:"node_id"`
	FormID string `json:"form_id"`
}

// ========== Iteration Events ==========

type IterationStartedEvent struct {
	BaseGraphEvent
	NodeID string `json:"node_id"`
	Total  int    `json:"total"`
}

type IterationNextEvent struct {
	BaseGraphEvent
	NodeID string         `json:"node_id"`
	Index  int            `json:"index"`
	Output map[string]any `json:"output,omitempty"`
}

type IterationCompletedEvent struct {
	BaseGraphEvent
	NodeID  string `json:"node_id"`
	Outputs []any  `json:"outputs"`
}

// ========== Loop Events ==========

type LoopStartedEvent struct {
	BaseGraphEvent
	NodeID   string `json:"node_id"`
	MaxLoops int    `json:"max_loops"`
}

type LoopNextEvent struct {
	BaseGraphEvent
	NodeID string `json:"node_id"`
	Index  int    `json:"index"`
}

type LoopCompletedEvent struct {
	BaseGraphEvent
	NodeID string `json:"node_id"`
	Loops  int    `json:"loops"`
}

// ========== Stream Events ==========

type TextChunkEvent struct {
	BaseGraphEvent
	NodeID  string `json:"node_id"`
	Text    string `json:"text"`
	FromVar string `json:"from_variable_selector,omitempty"`
}

type MessageReplaceEvent struct {
	BaseGraphEvent
	NodeID string `json:"node_id"`
	Text   string `json:"text"`
}

type StreamCompletedEvent struct {
	BaseGraphEvent
	NodeID           string `json:"node_id"`
	EdgeSourceHandle string `json:"edge_source_handle,omitempty"`
}

// ========== Parallel Events ==========

type ParallelBranchStartedEvent struct {
	BaseGraphEvent
	ParallelID  string `json:"parallel_id"`
	BranchIndex int    `json:"branch_index"`
	StartNodeID string `json:"start_node_id"`
}

type ParallelBranchCompletedEvent struct {
	BaseGraphEvent
	ParallelID  string `json:"parallel_id"`
	BranchIndex int    `json:"branch_index"`
	EndNodeID   string `json:"end_node_id"`
	Error       string `json:"error,omitempty"`
}
