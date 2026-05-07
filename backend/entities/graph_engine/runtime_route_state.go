package graphengine

import (
	"fmt"
	"time"

	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	graphengineenumtypes "github.com/odysseythink/gofy/backend/enum_types/graph_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// RouteNodeState represents the state of a route node
type RouteNodeState struct {
	ID            string                                        `json:"id"`
	NodeID        string                                        `json:"node_id"`
	NodeRunResult *workflowentities.NodeRunResult               `json:"node_run_result,omitempty"`
	Status        graphengineenumtypes.RouteNodeStateStatusType `json:"status"`
	StartAt       time.Time                                     `json:"start_at"`
	PausedAt      *time.Time                                    `json:"paused_at,omitempty"`
	FinishedAt    *time.Time                                    `json:"finished_at,omitempty"`
	FailedReason  string                                        `json:"failed_reason,omitempty"`
	PausedBy      string                                        `json:"paused_by,omitempty"`
	Index         int                                           `json:"index"`
}

func NewRouteNodeState() *RouteNodeState {
	return &RouteNodeState{
		ID:    uuid.NewV4().String(),
		Index: 1,
	}
}

// SetFinished sets the finished state of the node
func (rns *RouteNodeState) SetFinished(runResult *workflowentities.NodeRunResult) {
	if rns.Status == graphengineenumtypes.RouteNodeStateStatus_SUCCESS || rns.Status == graphengineenumtypes.RouteNodeStateStatus_FAILED || rns.Status == graphengineenumtypes.RouteNodeStateStatus_EXCEPTION {
		mlog.Errorf("route state %s already finished", rns.ID)
		panic(fmt.Errorf("route state %s already finished", rns.ID))
	}

	switch runResult.Status {
	case models.WorkflowNodeExecutionStatus_SUCCEEDED:
		rns.Status = graphengineenumtypes.RouteNodeStateStatus_SUCCESS
	case models.WorkflowNodeExecutionStatus_FAILED:
		rns.Status = graphengineenumtypes.RouteNodeStateStatus_FAILED
		rns.FailedReason = runResult.Error
	case models.WorkflowNodeExecutionStatus_EXCEPTION:
		rns.Status = graphengineenumtypes.RouteNodeStateStatus_EXCEPTION
		rns.FailedReason = runResult.Error
	default:
		mlog.Errorf("invalid route status %v", runResult.Status)
		panic(fmt.Errorf("invalid route status %v", runResult.Status))
	}

	rns.NodeRunResult = runResult
	now := time.Now().UTC()
	rns.FinishedAt = &now
}

// RuntimeRouteState represents the runtime state of routes
type RuntimeRouteState struct {
	Routes           map[string][]string        `json:"routes"`
	NodeStateMapping map[string]*RouteNodeState `json:"node_state_mapping"`
}

// CreateNodeState creates a new node state
func (rrs *RuntimeRouteState) CreateNodeState(nodeID string) *RouteNodeState {
	if rrs.NodeStateMapping == nil {
		rrs.NodeStateMapping = make(map[string]*RouteNodeState)
	}
	state := &RouteNodeState{
		ID:      uuid.NewV4().String(),
		NodeID:  nodeID,
		StartAt: time.Now(),
		Status:  graphengineenumtypes.RouteNodeStateStatus_RUNNING,
		Index:   1,
	}
	if rrs.NodeStateMapping == nil {
		rrs.NodeStateMapping = make(map[string]*RouteNodeState)
	}
	rrs.NodeStateMapping[state.ID] = state
	return state
}

// AddRoute adds a route to the graph state
func (rrs *RuntimeRouteState) AddRoute(sourceNodeStateID string, targetNodeStateID string) {
	if rrs.Routes == nil {
		rrs.Routes = make(map[string][]string)
	}
	if _, exists := rrs.Routes[sourceNodeStateID]; !exists {
		rrs.Routes[sourceNodeStateID] = []string{}
	}
	rrs.Routes[sourceNodeStateID] = append(rrs.Routes[sourceNodeStateID], targetNodeStateID)
}

// GetRoutesWithNodeStateBySourceNodeStateID gets routes with node state by source node state ID
func (rrs *RuntimeRouteState) GetRoutesWithNodeStateBySourceNodeStateID(sourceNodeStateID string) []*RouteNodeState {
	var result []*RouteNodeState
	for _, targetStateID := range rrs.Routes[sourceNodeStateID] {
		if state, exists := rrs.NodeStateMapping[targetStateID]; exists {
			result = append(result, state)
		}
	}
	return result
}
