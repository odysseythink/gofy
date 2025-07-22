package graphengine

import (
	"time"

	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
)

// GraphEngineEvent is a base event for the graph engine.
type GraphEngineEvent interface {
	EventName() string
}

// ##########################################
//  Graph Events
// ##########################################

// BaseGraphEvent is a base event for graph-related events.
type BaseGraphEvent struct {
}

func (e *BaseGraphEvent) EventName() string {
	return "base_graph_event"
}

// GraphRunStartedEvent is triggered when a graph run starts.
type GraphRunStartedEvent struct {
	*BaseGraphEvent
}

func (e *GraphRunStartedEvent) EventName() string {
	return "graph_run_started_event"
}

// GraphRunSucceededEvent is triggered when a graph run succeeds.
type GraphRunSucceededEvent struct {
	*BaseGraphEvent
	Outputs map[string]any `json:"outputs"`
}

func (e *GraphRunSucceededEvent) EventName() string {
	return "graph_run_succeeded_event"
}

// GraphRunFailedEvent is triggered when a graph run fails.
type GraphRunFailedEvent struct {
	*BaseGraphEvent
	Error           string `json:"error"`
	ExceptionsCount int    `json:"exceptions_count"`
}

func (e *GraphRunFailedEvent) EventName() string {
	return "graph_run_failed_event"
}

// GraphRunPartialSucceededEvent is triggered when a graph run partially succeeds.
type GraphRunPartialSucceededEvent struct {
	*BaseGraphEvent
	ExceptionsCount int            `json:"exceptions_count"`
	Outputs         map[string]any `json:"outputs"`
}

func (e *GraphRunPartialSucceededEvent) EventName() string {
	return "graph_run_partial_succeeded_event"
}

// ##########################################
//  Node Events
// ##########################################

// BaseNodeEvent is a base event for node-related events.
type BaseNodeEvent struct {
	// *GraphEngineEvent
	ID                        string                          `json:"id"`
	NodeID                    string                          `json:"node_id"`
	NodeType                  nodesenumtypes.NodeType         `json:"node_type"`
	NodeData                  *basenodesentities.BaseNodeData `json:"node_data"`
	RouteNodeState            *RouteNodeState                 `json:"route_node_state"` // Assuming RouteNodeState is an integer type
	ParallelID                string                          `json:"parallel_id"`
	ParallelStartNodeID       string                          `json:"parallel_start_node_id"`
	ParentParallelID          string                          `json:"parent_parallel_id"`
	ParentParallelStartNodeID string                          `json:"parent_parallel_start_node_id"`
	InIterationID             string                          `json:"in_iteration_id"`
}

func (e *BaseNodeEvent) EventName() string {
	return "base_node_event"
}

// NodeRunStartedEvent is triggered when a node run starts.
type NodeRunStartedEvent struct {
	*BaseNodeEvent
	PredecessorNodeID string `json:"predecessor_node_id"`
	ParallelModeRunID string `json:"parallel_mode_run_id"`
}

func (e *NodeRunStartedEvent) EventName() string {
	return "node_run_started_event"
}

// NodeRunStreamChunkEvent is triggered when a node run produces a stream chunk.
type NodeRunStreamChunkEvent struct {
	*BaseNodeEvent
	ChunkContent         string   `json:"chunk_content"`
	FromVariableSelector []string `json:"from_variable_selector"`
}

func (e *NodeRunStreamChunkEvent) EventName() string {
	return "node_run_stream_chunk_event"
}

// NodeRunRetrieverResourceEvent is triggered when a node run retrieves resources.
type NodeRunRetrieverResourceEvent struct {
	*BaseNodeEvent
	RetrieverResources []map[string]any `json:"retriever_resources"`
	Context            string           `json:"context"`
}

func (e *NodeRunRetrieverResourceEvent) EventName() string {
	return "node_run_retriever_resource_event"
}

// NodeRunSucceededEvent is triggered when a node run succeeds.
type NodeRunSucceededEvent struct {
	*BaseNodeEvent
}

func (e *NodeRunSucceededEvent) EventName() string {
	return "node_run_succeeded_event"
}

// NodeRunFailedEvent is triggered when a node run fails.
type NodeRunFailedEvent struct {
	*BaseNodeEvent
	Error string `json:"error"`
}

func (e *NodeRunFailedEvent) EventName() string {
	return "node_run_failed_event"
}

// NodeRunExceptionEvent is triggered when a node run encounters an exception.
type NodeRunExceptionEvent struct {
	*BaseNodeEvent
	Error string `json:"error"`
}

func (e *NodeRunExceptionEvent) EventName() string {
	return "node_run_exception_event"
}

// NodeInIterationFailedEvent is triggered when a node in an iteration fails.
type NodeInIterationFailedEvent struct {
	*BaseNodeEvent
	Error string `json:"error"`
}

func (e *NodeInIterationFailedEvent) EventName() string {
	return "node_in_iteration_failed_event"
}

// NodeRunRetryEvent is triggered when a node run is retried.
type NodeRunRetryEvent struct {
	*NodeRunStartedEvent
	Error      string    `json:"error"`
	RetryIndex int       `json:"retry_index"`
	StartAt    time.Time `json:"start_at"`
}

func (e *NodeRunRetryEvent) EventName() string {
	return "node_run_retry_event"
}

// ##########################################
//  Parallel Branch Events
// ##########################################

// BaseParallelBranchEvent is a base event for parallel branch-related events.
type BaseParallelBranchEvent struct {
	ParallelID                string `json:"parallel_id"`
	ParallelStartNodeID       string `json:"parallel_start_node_id"`
	ParentParallelID          string `json:"parent_Parallel_id"`
	ParentParallelStartNodeID string `json:"parent_Parallel_start_node_id"`
	InIterationID             string `json:"in_iteration_id"`
}

func (e *BaseParallelBranchEvent) EventName() string {
	return "base_parallel_branch_event"
}

// ParallelBranchRunStartedEvent is triggered when a parallel branch run starts.
type ParallelBranchRunStartedEvent struct {
	*BaseParallelBranchEvent
}

func (e *ParallelBranchRunStartedEvent) EventName() string {
	return "parallel_branch_run_started_event"
}

// ParallelBranchRunSucceededEvent is triggered when a parallel branch run succeeds.
type ParallelBranchRunSucceededEvent struct {
	*BaseParallelBranchEvent
}

func (e *ParallelBranchRunSucceededEvent) EventName() string {
	return "parallel_branch_run_succeeded_event"
}

// ParallelBranchRunFailedEvent is triggered when a parallel branch run fails.
type ParallelBranchRunFailedEvent struct {
	*BaseParallelBranchEvent
	Error string `json:"error"`
}

func (e *ParallelBranchRunFailedEvent) EventName() string {
	return "parallel_branch_run_failed_event"
}

// ##########################################
//  Iteration Events
// ##########################################

// BaseIterationEvent is a base event for iteration-related events.
type BaseIterationEvent struct {
	IterationID               string                          `json:"iteration_id"`
	IterationNodeID           string                          `json:"iteration_node_id"`
	IterationNodeType         nodesenumtypes.NodeType         `json:"iteration_node_type"`
	IterationNodeData         *basenodesentities.BaseNodeData `json:"iteration_node_data"`
	ParallelID                string                          `json:"parallel_id"`
	ParallelStartNodeID       string                          `json:"parallel_start_node_id"`
	ParentParallelID          string                          `json:"parent_Parallel_id"`
	ParentParallelStartNodeID string                          `json:"parent_Parallel_start_node_id"`
	ParallelModeRunID         string                          `json:"parallel_mode_run_id"`
}

func (e *BaseIterationEvent) EventName() string {
	return "base_iteration_event"
}

// IterationRunStartedEvent is triggered when an iteration run starts.
type IterationRunStartedEvent struct {
	*BaseIterationEvent
	StartAt           time.Time      `json:"start_at"`
	Inputs            map[string]any `json:"inputs"`
	Metadata          map[string]any `json:"metadata"`
	PredecessorNodeID string         `json:"predecessor_node_id"`
}

func (e *IterationRunStartedEvent) EventName() string {
	return "iteration_run_started_event"
}

// IterationRunNextEvent is triggered when an iteration proceeds to the next step.
type IterationRunNextEvent struct {
	*BaseIterationEvent
	Index              int     `json:"index"`
	PreIterationOutput any     `json:"pre_iteration_output"`
	Duration           float64 `json:"duration"`
}

func (e *IterationRunNextEvent) EventName() string {
	return "iteration_run_next_event"
}

// IterationRunSucceededEvent is triggered when an iteration run succeeds.
type IterationRunSucceededEvent struct {
	*BaseIterationEvent
	StartAt              time.Time          `json:"start_at"`
	Inputs               map[string]any     `json:"inputs"`
	Outputs              map[string]any     `json:"outputs"`
	Metadata             map[string]any     `json:"metadata"`
	Steps                int                `json:"steps"`
	IterationDurationMap map[string]float64 `json:"iteration_duration_map"`
}

func (e *IterationRunSucceededEvent) EventName() string {
	return "iteration_run_succeeded_event"
}

// IterationRunFailedEvent is triggered when an iteration run fails.
type IterationRunFailedEvent struct {
	*BaseIterationEvent
	StartAt  time.Time      `json:"start_at"`
	Inputs   map[string]any `json:"inputs"`
	Outputs  map[string]any `json:"outputs"`
	Metadata map[string]any `json:"metadata"`
	Steps    int            `json:"steps"`
	Error    string         `json:"error"`
}

func (e *IterationRunFailedEvent) EventName() string {
	return "iteration_run_failed_event"
}

// InNodeEvent represents events that can occur within a node.
type InNodeEvent interface {
	BaseNodeEvent | BaseParallelBranchEvent | BaseIterationEvent
}
