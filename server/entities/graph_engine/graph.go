package graphengine

import (
	uuid "github.com/satori/go.uuid"
)

// GraphEdge represents an edge in the graph
type GraphEdge struct {
	SourceNodeID string        `json:"source_node_id"`
	TargetNodeID string        `json:"target_node_id"`
	RunCondition *RunCondition `json:"run_condition,omitempty"`
}

// GraphParallel represents a parallel structure in the graph
type GraphParallel struct {
	ID                        string `json:"id"`
	StartFromNodeID           string `json:"start_from_node_id"`
	ParentParallelID          string `json:"parent_parallel_id,omitempty"`
	ParentParallelStartNodeID string `json:"parent_parallel_start_node_id,omitempty"`
	EndToNodeID               string `json:"end_to_node_id,omitempty"`
}

// NewGraphParallel creates a new instance of GraphParallel with default values
func NewGraphParallel(startFromNodeID string) *GraphParallel {
	return &GraphParallel{
		ID:              uuid.NewV4().String(),
		StartFromNodeID: startFromNodeID,
	}
}
