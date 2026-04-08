package events

import "time"

// EventType defines workflow event types.
type EventType string

const (
	EventGraphStarted       EventType = "graph_started"
	EventGraphSucceeded     EventType = "graph_succeeded"
	EventGraphFailed        EventType = "graph_failed"
	EventNodeStarted        EventType = "node_started"
	EventNodeSucceeded      EventType = "node_succeeded"
	EventNodeFailed         EventType = "node_failed"
	EventNodeRetrying       EventType = "node_retrying"
	EventIterationStarted   EventType = "iteration_started"
	EventIterationNext      EventType = "iteration_next"
	EventIterationCompleted EventType = "iteration_completed"
	EventLoopStarted        EventType = "loop_started"
	EventLoopNext           EventType = "loop_next"
	EventLoopCompleted      EventType = "loop_completed"
	EventPauseRequested     EventType = "pause_requested"
	EventResumed            EventType = "resumed"
	EventTextChunk          EventType = "text_chunk"
)

// GraphEvent is the base event emitted by the workflow engine.
type GraphEvent struct {
	Type      EventType      `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data,omitempty"`
}

// NodeEvent is emitted during node execution.
type NodeEvent struct {
	GraphEvent
	NodeID      string         `json:"node_id"`
	NodeType    string         `json:"node_type"`
	NodeTitle   string         `json:"node_title,omitempty"`
	Inputs      map[string]any `json:"inputs,omitempty"`
	Outputs     map[string]any `json:"outputs,omitempty"`
	Error       string         `json:"error,omitempty"`
	ElapsedTime float64        `json:"elapsed_time,omitempty"`
	RetryCount  int            `json:"retry_count,omitempty"`
}

// IterationEvent is emitted during iteration/loop execution.
type IterationEvent struct {
	GraphEvent
	NodeID string `json:"node_id"`
	Index  int    `json:"index"`
	Total  int    `json:"total,omitempty"`
}

// TextChunkEvent is emitted for streaming text output.
type TextChunkEvent struct {
	GraphEvent
	NodeID string `json:"node_id"`
	Text   string `json:"text"`
}

// PauseEvent is emitted when workflow is paused.
type PauseEvent struct {
	GraphEvent
	NodeID     string `json:"node_id"`
	ReasonType string `json:"reason_type"`
	FormID     string `json:"form_id,omitempty"`
	Message    string `json:"message,omitempty"`
}

// NewGraphEvent creates a new graph-level event.
func NewGraphEvent(eventType EventType, data map[string]any) *GraphEvent {
	return &GraphEvent{Type: eventType, Timestamp: time.Now(), Data: data}
}

// NewNodeEvent creates a new node-level event.
func NewNodeEvent(eventType EventType, nodeID, nodeType string) *NodeEvent {
	return &NodeEvent{
		GraphEvent: GraphEvent{Type: eventType, Timestamp: time.Now()},
		NodeID:     nodeID,
		NodeType:   nodeType,
	}
}
