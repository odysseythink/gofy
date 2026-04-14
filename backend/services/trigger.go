package services

import (
	"encoding/json"
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

// GetWebhookTrigger returns the webhook-trigger row for (app, node, tenant).
// userID is accepted but not filtered on — webhook rows are tenant-scoped.
func (s *TriggerService) GetWebhookTrigger(appID, nodeID, tenantID, userID string) *models.WorkflowWebhookTrigger {
	var wt models.WorkflowWebhookTrigger
	err := dbengine.Instance().DB.
		Where("app_id = ? AND node_id = ? AND tenant_id = ?", appID, nodeID, tenantID).
		First(&wt).Error
	if err != nil {
		return nil
	}
	return &wt
}

// ListSubscriptions returns plugin-trigger subscriptions for (tenant, provider).
func (s *TriggerService) ListSubscriptions(tenantID, provider string) []*models.TriggerSubscription {
	var subs []*models.TriggerSubscription
	dbengine.Instance().DB.
		Where("tenant_id = ? AND provider_id = ?", tenantID, provider).
		Order("created_at DESC").
		Find(&subs)
	return subs
}

// CreateSubscription creates a plugin trigger subscription for the tenant and
// registers a fresh endpoint_id that external systems can POST events to.
func (s *TriggerService) CreateSubscription(
	tenantID, userID, provider, name string,
	parameters, properties, credentials map[string]any,
	credentialType string,
) (*models.TriggerSubscription, error) {
	paramsJSON, err := json.Marshal(orEmpty(parameters))
	if err != nil {
		return nil, fmt.Errorf("marshal parameters: %w", err)
	}
	propsJSON, err := json.Marshal(orEmpty(properties))
	if err != nil {
		return nil, fmt.Errorf("marshal properties: %w", err)
	}
	credsJSON, err := json.Marshal(orEmpty(credentials))
	if err != nil {
		return nil, fmt.Errorf("marshal credentials: %w", err)
	}

	sub := &models.TriggerSubscription{
		ID:             uuid.NewV4().String(),
		Name:           name,
		TenantID:       tenantID,
		UserID:         userID,
		ProviderID:     provider,
		EndpointID:     uuid.NewV4().String(),
		Parameters:     string(paramsJSON),
		Properties:     string(propsJSON),
		Credentials:    string(credsJSON),
		CredentialType: credentialType,
	}
	if err := dbengine.Instance().DB.Create(sub).Error; err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	return sub, nil
}

// DeleteSubscription removes a trigger subscription by ID.
func (s *TriggerService) DeleteSubscription(subscriptionID string) error {
	res := dbengine.Instance().DB.
		Where("id = ?", subscriptionID).
		Delete(&models.TriggerSubscription{})
	if res.Error != nil {
		return fmt.Errorf("delete subscription: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("subscription %s not found", subscriptionID)
	}
	return nil
}

// HandlePluginEvent dispatches an inbound event delivered to the given endpoint.
// Looks up the subscription by endpoint_id, records a trigger log, and fans out
// to the trigger manager so matching workflows can fire.
func (s *TriggerService) HandlePluginEvent(endpointID string, data map[string]any) error {
	var sub models.TriggerSubscription
	if err := dbengine.Instance().DB.Where("endpoint_id = ?", endpointID).First(&sub).Error; err != nil {
		return fmt.Errorf("subscription for endpoint %s not found", endpointID)
	}

	dataJSON, _ := json.Marshal(data)
	metaJSON, _ := json.Marshal(map[string]any{
		"endpoint_id":     endpointID,
		"subscription_id": sub.ID,
		"provider_id":     sub.ProviderID,
	})
	log := &models.WorkflowTriggerLog{
		ID:              uuid.NewV4().String(),
		TenantID:        sub.TenantID,
		TriggerType:     "plugin",
		TriggerData:     string(dataJSON),
		TriggerMetadata: string(metaJSON),
		Inputs:          "{}",
		Status:          "pending",
		QueueName:       "default",
		CreatedByRole:   "system",
		CreatedBy:       "system",
	}
	dbengine.Instance().DB.Create(log)

	if mgr := trigger.GetTriggerManager(); mgr != nil {
		mgr.FireEvent(&trigger.TriggerEvent{
			TenantID:    sub.TenantID,
			TriggerType: "plugin",
			Data: map[string]any{
				"subscription_id": sub.ID,
				"provider_id":     sub.ProviderID,
				"endpoint_id":     endpointID,
				"payload":         data,
			},
		})
	}
	return nil
}

func orEmpty(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}
