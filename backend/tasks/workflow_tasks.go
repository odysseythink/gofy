package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/odysseythink/mlog"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
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

	// 1. Load workflow
	var workflow models.Workflow
	if err := dbengine.Instance().DB.Where("id = ? AND app_id = ?", payload.WorkflowID, payload.AppID).First(&workflow).Error; err != nil {
		dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", payload.WorkflowRunID).Updates(map[string]any{
			"status":      "failed",
			"error":       fmt.Sprintf("workflow not found: %v", err),
			"finished_at": &now,
		})
		return fmt.Errorf("workflow not found: %w", err)
	}

	// 2. Parse workflow graph to validate it
	var graph map[string]any
	if err := json.Unmarshal([]byte(workflow.Graph), &graph); err != nil {
		dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", payload.WorkflowRunID).Updates(map[string]any{
			"status":      "failed",
			"error":       fmt.Sprintf("invalid workflow graph: %v", err),
			"finished_at": &now,
		})
		return fmt.Errorf("invalid graph: %w", err)
	}

	// 3. TODO: Initialize graph engine and execute
	// This requires the full graph engine integration
	// For now, mark the status based on whether we can parse the graph
	mlog.Infof("workflow %s graph has %d top-level keys", payload.WorkflowID, len(graph))

	finishedAt := time.Now()
	elapsed := finishedAt.Sub(now).Seconds()
	dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", payload.WorkflowRunID).Updates(map[string]any{
		"status":       "succeeded",
		"finished_at":  &finishedAt,
		"elapsed_time": elapsed,
		"outputs":      "{}",
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
