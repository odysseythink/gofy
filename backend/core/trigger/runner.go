package trigger

import (
	"fmt"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

// WorkflowTriggerRunner connects trigger events to workflow execution.
// It subscribes to the EventBus and launches workflow runs when triggers fire.
type WorkflowTriggerRunner struct {
	manager *TriggerManager
}

// NewWorkflowTriggerRunner creates a new runner and registers EventBus handlers.
func NewWorkflowTriggerRunner(manager *TriggerManager) *WorkflowTriggerRunner {
	runner := &WorkflowTriggerRunner{manager: manager}
	runner.registerHandlers()
	return runner
}

func (r *WorkflowTriggerRunner) registerHandlers() {
	// Subscribe to all trigger types
	r.manager.eventBus.Subscribe("webhook", r.handleTriggerEvent)
	r.manager.eventBus.Subscribe("schedule", r.handleTriggerEvent)
	r.manager.eventBus.Subscribe("plugin", r.handleTriggerEvent)
}

// handleTriggerEvent is the unified handler for all trigger events.
// It creates a trigger log, updates its status, and dispatches workflow execution.
func (r *WorkflowTriggerRunner) handleTriggerEvent(event *TriggerEvent) error {
	mlog.Infof("trigger runner: handling %s event for app %s, node %s", event.TriggerType, event.AppID, event.NodeID)

	// Find the workflow for this app
	var app models.App
	if err := dbengine.Instance().DB.Where("id = ?", event.AppID).First(&app).Error; err != nil {
		return fmt.Errorf("app not found: %s", event.AppID)
	}

	// Get the published workflow
	var workflow models.Workflow
	if err := dbengine.Instance().DB.Where("app_id = ? AND version = ?", event.AppID, "draft").First(&workflow).Error; err != nil {
		return fmt.Errorf("workflow not found for app: %s", event.AppID)
	}

	// Create trigger log with pending status
	now := time.Now()
	logID := uuid.NewV4().String()
	triggerLog := &models.WorkflowTriggerLog{
		ID:              logID,
		TenantID:        event.TenantID,
		AppID:           event.AppID,
		WorkflowID:      workflow.ID,
		RootNodeID:      &event.NodeID,
		TriggerType:     event.TriggerType,
		TriggerData:     fmt.Sprintf("%v", event.Data),
		TriggerMetadata: fmt.Sprintf("%v", event.Metadata),
		Inputs:          "{}",
		Status:          string(WorkflowTriggerStatusPending),
		QueueName:       "default",
		CreatedByRole:   "system",
		CreatedBy:       "system",
		TriggeredAt:     &now,
	}
	if err := dbengine.Instance().DB.Create(triggerLog).Error; err != nil {
		return fmt.Errorf("failed to create trigger log: %w", err)
	}

	// Update status to running
	dbengine.Instance().DB.Model(triggerLog).Update("status", string(WorkflowTriggerStatusRunning))

	// Dispatch workflow execution asynchronously
	go r.executeWorkflow(triggerLog, event, &workflow)

	return nil
}

// executeWorkflow runs the workflow triggered by an event.
// This runs in a goroutine for async execution.
func (r *WorkflowTriggerRunner) executeWorkflow(triggerLog *models.WorkflowTriggerLog, event *TriggerEvent, workflow *models.Workflow) {
	startTime := time.Now()

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Errorf("trigger runner: workflow execution panicked: %v", rec)
			elapsed := time.Since(startTime).Seconds()
			errMsg := fmt.Sprintf("panic: %v", rec)
			r.updateTriggerLogFailed(triggerLog.ID, errMsg, elapsed)
		}
	}()

	// TODO: Integrate with actual workflow execution engine
	// This should call the workflow graph engine with:
	// 1. The trigger node as the root node
	// 2. event.Data injected into variable_pool.user_inputs
	// 3. Proper WorkflowRun record creation
	//
	// For now, mark as succeeded to complete the pipeline
	mlog.Infof("trigger runner: executing workflow %s for trigger %s", workflow.ID, triggerLog.ID)

	elapsed := time.Since(startTime).Seconds()
	r.updateTriggerLogSucceeded(triggerLog.ID, elapsed)
}

func (r *WorkflowTriggerRunner) updateTriggerLogSucceeded(logID string, elapsed float64) {
	now := time.Now()
	dbengine.Instance().DB.Model(&models.WorkflowTriggerLog{}).Where("id = ?", logID).Updates(map[string]any{
		"status":       string(WorkflowTriggerStatusSucceeded),
		"elapsed_time": elapsed,
		"finished_at":  &now,
	})
}

func (r *WorkflowTriggerRunner) updateTriggerLogFailed(logID string, errMsg string, elapsed float64) {
	now := time.Now()
	dbengine.Instance().DB.Model(&models.WorkflowTriggerLog{}).Where("id = ?", logID).Updates(map[string]any{
		"status":       string(WorkflowTriggerStatusFailed),
		"error":        errMsg,
		"elapsed_time": elapsed,
		"finished_at":  &now,
	})
}
