package trigger

import (
	"encoding/json"
	"fmt"
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// PluginTriggerProvider handles plugin-based trigger events.
// Plugin triggers work through subscriptions: external services send events
// to a registered endpoint, which are then dispatched to subscribed workflows.
type PluginTriggerProvider struct{}

func NewPluginTriggerProvider() *PluginTriggerProvider {
	return &PluginTriggerProvider{}
}

func (pp *PluginTriggerProvider) Type() TriggerType { return TriggerTypePlugin }

func (pp *PluginTriggerProvider) Setup(config map[string]any) error { return nil }

func (pp *PluginTriggerProvider) Teardown() error { return nil }

func (pp *PluginTriggerProvider) HandleEvent(event *TriggerEvent) error { return nil }

// CreatePluginTrigger creates a plugin trigger mapping for a workflow node.
func (pp *PluginTriggerProvider) CreatePluginTrigger(appID, nodeID, tenantID, providerID, eventName, subscriptionID string) (*models.WorkflowPluginTrigger, error) {
	trigger := &models.WorkflowPluginTrigger{
		ID:             uuid.NewV4().String(),
		AppID:          appID,
		NodeID:         nodeID,
		TenantID:       tenantID,
		ProviderID:     providerID,
		EventName:      eventName,
		SubscriptionID: subscriptionID,
	}

	if err := dbengine.Instance().DB.Create(trigger).Error; err != nil {
		return nil, fmt.Errorf("failed to create plugin trigger: %w", err)
	}
	return trigger, nil
}

// GetPluginTriggersBySubscription returns plugin triggers for a subscription.
func (pp *PluginTriggerProvider) GetPluginTriggersBySubscription(subscriptionID string) []*models.WorkflowPluginTrigger {
	var triggers []*models.WorkflowPluginTrigger
	dbengine.Instance().DB.Where("subscription_id = ?", subscriptionID).Find(&triggers)
	return triggers
}

// GetPluginTriggersByEvent returns plugin triggers for a specific event.
func (pp *PluginTriggerProvider) GetPluginTriggersByEvent(providerID, eventName string) []*models.WorkflowPluginTrigger {
	var triggers []*models.WorkflowPluginTrigger
	dbengine.Instance().DB.Where("provider_id = ? AND event_name = ?", providerID, eventName).Find(&triggers)
	return triggers
}

// DeletePluginTrigger removes a plugin trigger.
func (pp *PluginTriggerProvider) DeletePluginTrigger(appID, nodeID string) error {
	return dbengine.Instance().DB.Where("app_id = ? AND node_id = ?", appID, nodeID).Delete(&models.WorkflowPluginTrigger{}).Error
}

// DispatchPluginEvent dispatches an incoming plugin event to all subscribed workflows.
func (pp *PluginTriggerProvider) DispatchPluginEvent(manager *TriggerManager, endpointID string, payload map[string]any) error {
	// Find the subscription for this endpoint
	var subscription models.TriggerSubscription
	if err := dbengine.Instance().DB.Where("endpoint_id = ?", endpointID).First(&subscription).Error; err != nil {
		return fmt.Errorf("subscription not found for endpoint: %s", endpointID)
	}

	// Find all plugin triggers using this subscription
	var pluginTriggers []*models.WorkflowPluginTrigger
	dbengine.Instance().DB.Where("subscription_id = ?", subscription.ID).Find(&pluginTriggers)

	if len(pluginTriggers) == 0 {
		mlog.Infof("no workflows subscribed to endpoint %s", endpointID)
		return nil
	}

	mlog.Infof("dispatching plugin event to %d workflows for endpoint %s", len(pluginTriggers), endpointID)

	for _, pt := range pluginTriggers {
		// Check that the app trigger is enabled
		var appTrigger models.AppTrigger
		err := dbengine.Instance().DB.Where(
			"app_id = ? AND node_id = ? AND trigger_type = ? AND status = ?",
			pt.AppID, pt.NodeID, string(AppTriggerTypePlugin), string(AppTriggerStatusEnabled),
		).First(&appTrigger).Error
		if err != nil {
			continue // Trigger disabled
		}

		manager.FireEvent(&TriggerEvent{
			ID:          uuid.NewV4().String(),
			TenantID:    pt.TenantID,
			AppID:       pt.AppID,
			TriggerType: string(TriggerTypePlugin),
			NodeID:      pt.NodeID,
			Data:        payload,
			Metadata: map[string]any{
				"provider_id":     pt.ProviderID,
				"event_name":      pt.EventName,
				"subscription_id": pt.SubscriptionID,
			},
			CreatedAt: time.Now(),
		})
	}

	return nil
}

// CreateSubscription creates a new trigger subscription.
func (pp *PluginTriggerProvider) CreateSubscription(tenantID, userID, providerID, name string, parameters, properties, credentials map[string]any, credentialType string) (*models.TriggerSubscription, error) {
	endpointID := uuid.NewV4().String()

	paramsJSON, _ := json.Marshal(parameters)
	propsJSON, _ := json.Marshal(properties)
	credsJSON, _ := json.Marshal(credentials)

	sub := &models.TriggerSubscription{
		ID:             uuid.NewV4().String(),
		Name:           name,
		TenantID:       tenantID,
		UserID:         userID,
		ProviderID:     providerID,
		EndpointID:     endpointID,
		Parameters:     string(paramsJSON),
		Properties:     string(propsJSON),
		Credentials:    string(credsJSON),
		CredentialType: credentialType,
	}

	if err := dbengine.Instance().DB.Create(sub).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}
	return sub, nil
}

// DeleteSubscription removes a trigger subscription and its associated plugin triggers.
func (pp *PluginTriggerProvider) DeleteSubscription(subscriptionID string) error {
	// Delete associated plugin triggers
	dbengine.Instance().DB.Where("subscription_id = ?", subscriptionID).Delete(&models.WorkflowPluginTrigger{})
	// Delete the subscription
	return dbengine.Instance().DB.Where("id = ?", subscriptionID).Delete(&models.TriggerSubscription{}).Error
}

// GetSubscription returns a subscription by ID.
func (pp *PluginTriggerProvider) GetSubscription(subscriptionID string) *models.TriggerSubscription {
	var sub models.TriggerSubscription
	if err := dbengine.Instance().DB.Where("id = ?", subscriptionID).First(&sub).Error; err != nil {
		return nil
	}
	return &sub
}

// ListSubscriptions returns all subscriptions for a tenant and provider.
func (pp *PluginTriggerProvider) ListSubscriptions(tenantID, providerID string) []*models.TriggerSubscription {
	var subs []*models.TriggerSubscription
	dbengine.Instance().DB.Where("tenant_id = ? AND provider_id = ?", tenantID, providerID).Find(&subs)
	return subs
}
