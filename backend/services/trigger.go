package services

import (
	"fmt"

	"mlib.com/gofy/server/core/trigger"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"

	uuid "github.com/satori/go.uuid"
)

type TriggerService struct{}

// CreateAppTrigger creates a trigger for an app.
func (s *TriggerService) CreateAppTrigger(tenantID, appID, nodeID, triggerType, title, providerName string) (*models.AppTrigger, error) {
	mgr := trigger.GetTriggerManager()
	if mgr == nil {
		return nil, fmt.Errorf("trigger manager not initialized")
	}
	return mgr.CreateAppTrigger(tenantID, appID, nodeID, trigger.TriggerType(triggerType), title, providerName)
}

// DeleteAppTrigger removes an app trigger.
func (s *TriggerService) DeleteAppTrigger(triggerID string) error {
	mgr := trigger.GetTriggerManager()
	if mgr == nil {
		return fmt.Errorf("trigger manager not initialized")
	}
	return mgr.DeleteAppTrigger(triggerID)
}

// GetAppTriggers lists triggers for an app.
func (s *TriggerService) GetAppTriggers(tenantID, appID string) []*models.AppTrigger {
	mgr := trigger.GetTriggerManager()
	if mgr == nil {
		return nil
	}
	return mgr.GetAppTriggers(tenantID, appID)
}

// EnableAppTrigger enables or disables a trigger.
func (s *TriggerService) EnableAppTrigger(triggerID string, enabled bool) error {
	mgr := trigger.GetTriggerManager()
	if mgr == nil {
		return fmt.Errorf("trigger manager not initialized")
	}
	return mgr.EnableAppTrigger(triggerID, enabled)
}

// HandleWebhook processes an incoming webhook event.
func (s *TriggerService) HandleWebhook(webhookID string, data map[string]any) error {
	// Find the webhook trigger
	var wt models.WorkflowWebhookTrigger
	if err := dbengine.Instance().DB.Where("webhook_id = ?", webhookID).First(&wt).Error; err != nil {
		return fmt.Errorf("webhook not found: %s", webhookID)
	}

	// Create trigger log
	log := &models.WorkflowTriggerLog{
		ID:              uuid.NewV4().String(),
		TenantID:        wt.TenantID,
		AppID:           wt.AppID,
		TriggerType:     "webhook",
		TriggerData:     fmt.Sprintf("%v", data),
		TriggerMetadata: fmt.Sprintf(`{"webhook_id":"%s"}`, webhookID),
		Inputs:          "{}",
		Status:          "pending",
		QueueName:       "default",
		CreatedByRole:   "system",
		CreatedBy:       "system",
	}
	dbengine.Instance().DB.Create(log)

	// Fire event
	mgr := trigger.GetTriggerManager()
	if mgr != nil {
		mgr.FireEvent(&trigger.TriggerEvent{
			TenantID:    wt.TenantID,
			AppID:       wt.AppID,
			TriggerType: "webhook",
			NodeID:      wt.NodeID,
			Data:        data,
		})
	}

	return nil
}

// GetTriggerLogs returns trigger execution logs.
func (s *TriggerService) GetTriggerLogs(tenantID, appID string, page, pageSize int) ([]*models.WorkflowTriggerLog, int64) {
	var logs []*models.WorkflowTriggerLog
	var total int64
	query := dbengine.Instance().DB.Where("tenant_id = ? AND app_id = ?", tenantID, appID)
	query.Model(&models.WorkflowTriggerLog{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)
	return logs, total
}

// CreateSchedulePlan creates a cron schedule for a workflow trigger.
func (s *TriggerService) CreateSchedulePlan(appID, nodeID, tenantID, cronExpr, timezone string) (*models.WorkflowSchedulePlan, error) {
	sp := trigger.NewScheduleProvider()
	return sp.CreateSchedulePlan(appID, nodeID, tenantID, cronExpr, timezone)
}

// GetSchedulePlans returns schedule plans for an app.
func (s *TriggerService) GetSchedulePlans(appID, tenantID string) []*models.WorkflowSchedulePlan {
	sp := trigger.NewScheduleProvider()
	return sp.GetSchedulePlans(appID, tenantID)
}
