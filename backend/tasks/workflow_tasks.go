package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

// WorkflowExecutionPayload is the payload for workflow execution tasks.
type WorkflowExecutionPayload struct {
	TenantID      string         `json:"tenant_id"`
	AppID         string         `json:"app_id"`
	WorkflowID    string         `json:"workflow_id"`
	WorkflowRunID string         `json:"workflow_run_id"`
	UserID        string         `json:"user_id"`
	Inputs        map[string]any `json:"inputs"`
	TriggeredFrom string         `json:"triggered_from"`
}

// HandleWorkflowExecution processes an async workflow execution.
func HandleWorkflowExecution(ctx context.Context, task *Task) error {
	var payload WorkflowExecutionPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("executing workflow %s run %s", payload.WorkflowID, payload.WorkflowRunID)

	now := time.Now()
	// Update run status to running
	dbengine.Instance().DB.Model(&models.WorkflowRun{}).
		Where("id = ?", payload.WorkflowRunID).
		Updates(map[string]any{"status": "running", "started_at": &now})

	// TODO: Execute the actual workflow graph
	// 1. Load workflow definition
	// 2. Initialize variable pool with inputs
	// 3. Run graph engine
	// 4. Collect results

	// For now, mark as succeeded
	dbengine.Instance().DB.Model(&models.WorkflowRun{}).
		Where("id = ?", payload.WorkflowRunID).
		Updates(map[string]any{
			"status":       "succeeded",
			"finished_at":  &now,
			"elapsed_time": 0.0,
		})

	mlog.Infof("workflow run %s completed", payload.WorkflowRunID)
	return nil
}

// HandleWorkflowNodeExecution processes a single node execution task.
func HandleWorkflowNodeExecution(ctx context.Context, task *Task) error {
	var payload struct {
		WorkflowRunID   string `json:"workflow_run_id"`
		NodeExecutionID string `json:"node_execution_id"`
		NodeID          string `json:"node_id"`
	}
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("executing node %s in workflow run %s", payload.NodeID, payload.WorkflowRunID)
	// TODO: Execute specific node
	return nil
}
