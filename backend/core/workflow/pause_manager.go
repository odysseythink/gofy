package workflow

import (
	"encoding/json"
	"fmt"
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/storage"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// PauseReasonType defines why a workflow was paused.
type PauseReasonType string

const (
	PauseReasonHumanInput PauseReasonType = "human_input"
	PauseReasonTimeout    PauseReasonType = "timeout"
	PauseReasonManual     PauseReasonType = "manual"
)

// PauseRequest represents a request to pause a workflow.
type PauseRequest struct {
	WorkflowID    string
	WorkflowRunID string
	NodeID        string
	ReasonType    PauseReasonType
	FormID        string
	Message       string
}

// GraphState represents serializable workflow execution state.
type GraphState struct {
	VariablePoolData map[string]any `json:"variable_pool"`
	CurrentNodeID    string         `json:"current_node_id"`
	ExecutedNodes    []string       `json:"executed_nodes"`
	PendingNodes     []string       `json:"pending_nodes"`
	Timestamp        time.Time      `json:"timestamp"`
}

// PauseManager handles workflow pause and resume operations.
type PauseManager struct{}

// NewPauseManager creates a new pause manager.
func NewPauseManager() *PauseManager {
	return &PauseManager{}
}

// PauseWorkflow pauses a running workflow and saves its state.
func (pm *PauseManager) PauseWorkflow(req PauseRequest, state *GraphState) error {
	// Serialize state
	stateBytes, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to serialize graph state: %w", err)
	}

	// Store state in storage backend
	stateKey := fmt.Sprintf("workflow_pause/%s/%s.json", req.WorkflowRunID, time.Now().Format("20060102150405"))
	if err := storage.Save(stateKey, stateBytes); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	// Create pause record
	pause := &models.WorkflowPause{
		ID:             uuid.NewV4().String(),
		WorkflowID:     req.WorkflowID,
		WorkflowRunID:  req.WorkflowRunID,
		StateObjectKey: stateKey,
	}
	if err := dbengine.Instance().DB.Create(pause).Error; err != nil {
		return fmt.Errorf("failed to create pause record: %w", err)
	}

	// Create pause reason
	reason := &models.WorkflowPauseReason{
		ID:      uuid.NewV4().String(),
		PauseID: pause.ID,
		Type:    string(req.ReasonType),
		FormID:  req.FormID,
		Message: req.Message,
		NodeID:  req.NodeID,
	}
	if err := dbengine.Instance().DB.Create(reason).Error; err != nil {
		return fmt.Errorf("failed to create pause reason: %w", err)
	}

	// Update workflow run status
	dbengine.Instance().DB.Model(&models.WorkflowRun{}).
		Where("id = ?", req.WorkflowRunID).
		Update("status", "paused")

	mlog.Infof("workflow run %s paused at node %s (reason: %s)", req.WorkflowRunID, req.NodeID, req.ReasonType)
	return nil
}

// ResumeWorkflow resumes a paused workflow.
func (pm *PauseManager) ResumeWorkflow(workflowRunID string) (*GraphState, error) {
	// Find pause record
	var pause models.WorkflowPause
	if err := dbengine.Instance().DB.Where("workflow_run_id = ? AND resumed_at IS NULL", workflowRunID).First(&pause).Error; err != nil {
		return nil, fmt.Errorf("no active pause found for workflow run %s", workflowRunID)
	}

	// Load state from storage
	stateBytes, err := storage.LoadOnce(pause.StateObjectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	var state GraphState
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		return nil, fmt.Errorf("failed to deserialize state: %w", err)
	}

	// Mark as resumed
	now := time.Now()
	dbengine.Instance().DB.Model(&pause).Update("resumed_at", &now)

	// Update workflow run status
	dbengine.Instance().DB.Model(&models.WorkflowRun{}).
		Where("id = ?", workflowRunID).
		Update("status", "running")

	mlog.Infof("workflow run %s resumed", workflowRunID)
	return &state, nil
}

// GetPauseReason returns the reason a workflow was paused.
func (pm *PauseManager) GetPauseReason(workflowRunID string) (*models.WorkflowPauseReason, error) {
	var pause models.WorkflowPause
	if err := dbengine.Instance().DB.Where("workflow_run_id = ? AND resumed_at IS NULL", workflowRunID).First(&pause).Error; err != nil {
		return nil, fmt.Errorf("no active pause found")
	}

	var reason models.WorkflowPauseReason
	if err := dbengine.Instance().DB.Where("pause_id = ?", pause.ID).First(&reason).Error; err != nil {
		return nil, fmt.Errorf("pause reason not found")
	}

	return &reason, nil
}

// IsPaused checks if a workflow run is currently paused.
func (pm *PauseManager) IsPaused(workflowRunID string) bool {
	var count int64
	dbengine.Instance().DB.Model(&models.WorkflowPause{}).
		Where("workflow_run_id = ? AND resumed_at IS NULL", workflowRunID).
		Count(&count)
	return count > 0
}
