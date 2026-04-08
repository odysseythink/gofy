package graphengine

import (
	"sync"
	"time"
)

// ExecutionStatus represents the current state of a workflow execution.
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusSucceeded ExecutionStatus = "succeeded"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusPaused    ExecutionStatus = "paused"
	ExecutionStatusAborted   ExecutionStatus = "aborted"
	ExecutionStatusPartial   ExecutionStatus = "partial_succeeded"
	ExecutionStatusSkipped   ExecutionStatus = "skipped"
)

// NodeExecutionState tracks the execution state of a single node.
type NodeExecutionState struct {
	NodeID     string
	NodeType   string
	Status     ExecutionStatus
	StartTime  time.Time
	EndTime    time.Time
	Inputs     map[string]any
	Outputs    map[string]any
	Error      string
	RetryCount int
	Metadata   map[string]any
}

// GraphStateManager manages the execution state of a workflow graph.
type GraphStateManager struct {
	mu              sync.RWMutex
	status          ExecutionStatus
	startTime       time.Time
	endTime         time.Time
	nodeStates      map[string]*NodeExecutionState
	executedNodeIDs []string
	pendingNodeIDs  []string
	currentNodeID   string
	totalTokens     int
	totalSteps      int
	exceptions      int
	error           string
}

// NewGraphStateManager creates a new state manager.
func NewGraphStateManager() *GraphStateManager {
	return &GraphStateManager{
		status:     ExecutionStatusPending,
		nodeStates: make(map[string]*NodeExecutionState),
	}
}

// Status returns the current execution status.
func (sm *GraphStateManager) Status() ExecutionStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.status
}

// SetStatus updates the execution status.
func (sm *GraphStateManager) SetStatus(status ExecutionStatus) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.status = status
	if status == ExecutionStatusRunning && sm.startTime.IsZero() {
		sm.startTime = time.Now()
	}
	if status == ExecutionStatusSucceeded || status == ExecutionStatusFailed || status == ExecutionStatusAborted {
		sm.endTime = time.Now()
	}
}

// StartNode marks a node as started.
func (sm *GraphStateManager) StartNode(nodeID, nodeType string, inputs map[string]any) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.currentNodeID = nodeID
	sm.nodeStates[nodeID] = &NodeExecutionState{
		NodeID:    nodeID,
		NodeType:  nodeType,
		Status:    ExecutionStatusRunning,
		StartTime: time.Now(),
		Inputs:    inputs,
	}
}

// CompleteNode marks a node as completed.
func (sm *GraphStateManager) CompleteNode(nodeID string, outputs map[string]any, metadata map[string]any) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if state, ok := sm.nodeStates[nodeID]; ok {
		state.Status = ExecutionStatusSucceeded
		state.EndTime = time.Now()
		state.Outputs = outputs
		state.Metadata = metadata
	}
	sm.executedNodeIDs = append(sm.executedNodeIDs, nodeID)
	sm.totalSteps++
	// Remove from pending
	for i, id := range sm.pendingNodeIDs {
		if id == nodeID {
			sm.pendingNodeIDs = append(sm.pendingNodeIDs[:i], sm.pendingNodeIDs[i+1:]...)
			break
		}
	}
}

// FailNode marks a node as failed.
func (sm *GraphStateManager) FailNode(nodeID string, err string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if state, ok := sm.nodeStates[nodeID]; ok {
		state.Status = ExecutionStatusFailed
		state.EndTime = time.Now()
		state.Error = err
	}
	sm.exceptions++
}

// SkipNode marks a node as skipped.
func (sm *GraphStateManager) SkipNode(nodeID, reason string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.nodeStates[nodeID] = &NodeExecutionState{
		NodeID: nodeID,
		Status: ExecutionStatusSkipped,
		Error:  reason,
	}
}

// RetryNode increments retry count for a node.
func (sm *GraphStateManager) RetryNode(nodeID string) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if state, ok := sm.nodeStates[nodeID]; ok {
		state.RetryCount++
		return state.RetryCount
	}
	return 0
}

// AddPendingNode adds a node to the pending list.
func (sm *GraphStateManager) AddPendingNode(nodeID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.pendingNodeIDs = append(sm.pendingNodeIDs, nodeID)
}

// AddTokens accumulates token usage.
func (sm *GraphStateManager) AddTokens(tokens int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.totalTokens += tokens
}

// GetNodeState returns the state of a specific node.
func (sm *GraphStateManager) GetNodeState(nodeID string) *NodeExecutionState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.nodeStates[nodeID]
}

// IsNodeExecuted checks if a node has been executed.
func (sm *GraphStateManager) IsNodeExecuted(nodeID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _, id := range sm.executedNodeIDs {
		if id == nodeID {
			return true
		}
	}
	return false
}

// CurrentNodeID returns the currently executing node.
func (sm *GraphStateManager) CurrentNodeID() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.currentNodeID
}

// ExecutedNodeIDs returns all executed node IDs in order.
func (sm *GraphStateManager) ExecutedNodeIDs() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make([]string, len(sm.executedNodeIDs))
	copy(result, sm.executedNodeIDs)
	return result
}

// PendingNodeIDs returns all pending node IDs.
func (sm *GraphStateManager) PendingNodeIDs() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make([]string, len(sm.pendingNodeIDs))
	copy(result, sm.pendingNodeIDs)
	return result
}

// ElapsedTime returns the total elapsed time in seconds.
func (sm *GraphStateManager) ElapsedTime() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.elapsedTime()
}

// elapsedTime is the non-locking internal implementation.
func (sm *GraphStateManager) elapsedTime() float64 {
	end := sm.endTime
	if end.IsZero() {
		end = time.Now()
	}
	if sm.startTime.IsZero() {
		return 0
	}
	return end.Sub(sm.startTime).Seconds()
}

// TotalTokens returns accumulated token count.
func (sm *GraphStateManager) TotalTokens() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.totalTokens
}

// TotalSteps returns the number of executed steps.
func (sm *GraphStateManager) TotalSteps() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.totalSteps
}

// Exceptions returns the exception count.
func (sm *GraphStateManager) Exceptions() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.exceptions
}

// Summary returns a summary of the execution.
func (sm *GraphStateManager) Summary() map[string]any {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return map[string]any{
		"status":         string(sm.status),
		"total_steps":    sm.totalSteps,
		"total_tokens":   sm.totalTokens,
		"exceptions":     sm.exceptions,
		"elapsed_time":   sm.elapsedTime(),
		"executed_nodes": len(sm.executedNodeIDs),
		"pending_nodes":  len(sm.pendingNodeIDs),
	}
}
