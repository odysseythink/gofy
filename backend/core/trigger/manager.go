package trigger

import (
	"fmt"
	"sync"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

// TriggerType defines the type of trigger.
type TriggerType string

const (
	TriggerTypeWebhook  TriggerType = "webhook"
	TriggerTypeSchedule TriggerType = "schedule"
	TriggerTypePlugin   TriggerType = "plugin"
)

// TriggerManager orchestrates trigger lifecycle.
type TriggerManager struct {
	mu        sync.RWMutex
	providers map[string]TriggerProvider
	eventBus  *EventBus
}

// TriggerProvider is the interface for trigger providers.
type TriggerProvider interface {
	Type() TriggerType
	Setup(config map[string]any) error
	Teardown() error
	HandleEvent(event *TriggerEvent) error
}

var globalTriggerManager *TriggerManager

// InitTriggerManager initializes the global trigger manager.
func InitTriggerManager() {
	globalTriggerManager = &TriggerManager{
		providers: make(map[string]TriggerProvider),
		eventBus:  NewEventBus(),
	}
	globalTriggerManager.eventBus.Start()
	mlog.Info("trigger manager initialized")
}

// GetTriggerManager returns the global trigger manager.
func GetTriggerManager() *TriggerManager {
	return globalTriggerManager
}

// RegisterProvider registers a trigger provider.
func (tm *TriggerManager) RegisterProvider(name string, provider TriggerProvider) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.providers[name] = provider
}

// CreateAppTrigger creates a trigger for an app.
func (tm *TriggerManager) CreateAppTrigger(tenantID, appID, nodeID string, triggerType TriggerType, title string, providerName string) (*models.AppTrigger, error) {
	trigger := &models.AppTrigger{
		ID:           uuid.NewV4().String(),
		TenantID:     tenantID,
		AppID:        appID,
		NodeID:       &nodeID,
		TriggerType:  string(triggerType),
		Title:        title,
		ProviderName: &providerName,
		Status:       "enabled",
	}
	if err := dbengine.Instance().DB.Create(trigger).Error; err != nil {
		return nil, fmt.Errorf("failed to create app trigger: %w", err)
	}
	return trigger, nil
}

// DeleteAppTrigger removes an app trigger.
func (tm *TriggerManager) DeleteAppTrigger(triggerID string) error {
	return dbengine.Instance().DB.Where("id = ?", triggerID).Delete(&models.AppTrigger{}).Error
}

// GetAppTriggers lists triggers for an app.
func (tm *TriggerManager) GetAppTriggers(tenantID, appID string) []*models.AppTrigger {
	var triggers []*models.AppTrigger
	dbengine.Instance().DB.Where("tenant_id = ? AND app_id = ?", tenantID, appID).Find(&triggers)
	return triggers
}

// EnableAppTrigger enables or disables a trigger.
func (tm *TriggerManager) EnableAppTrigger(triggerID string, enabled bool) error {
	status := "enabled"
	if !enabled {
		status = "disabled"
	}
	return dbengine.Instance().DB.Model(&models.AppTrigger{}).Where("id = ?", triggerID).Update("status", status).Error
}

// FireEvent fires a trigger event.
func (tm *TriggerManager) FireEvent(event *TriggerEvent) {
	tm.eventBus.Publish(event)
}

// Stop shuts down the trigger manager.
func (tm *TriggerManager) Stop() {
	tm.eventBus.Stop()
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for _, p := range tm.providers {
		p.Teardown()
	}
}
