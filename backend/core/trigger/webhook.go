package trigger

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	uuid "github.com/satori/go.uuid"
)

// WebhookProvider handles webhook-based triggers.
type WebhookProvider struct {
	baseURL string
}

func NewWebhookProvider(baseURL string) *WebhookProvider {
	return &WebhookProvider{baseURL: baseURL}
}

func (wp *WebhookProvider) Type() TriggerType {
	return TriggerTypeWebhook
}

func (wp *WebhookProvider) Setup(config map[string]any) error {
	return nil
}

func (wp *WebhookProvider) Teardown() error {
	return nil
}

func (wp *WebhookProvider) HandleEvent(event *TriggerEvent) error {
	return nil
}

// CreateWebhookTrigger creates a webhook trigger for a workflow node.
func (wp *WebhookProvider) CreateWebhookTrigger(appID, nodeID, tenantID, createdBy string) (*models.WorkflowWebhookTrigger, error) {
	webhookID, err := generateWebhookID()
	if err != nil {
		return nil, err
	}

	trigger := &models.WorkflowWebhookTrigger{
		ID:        uuid.NewV4().String(),
		AppID:     appID,
		NodeID:    nodeID,
		TenantID:  tenantID,
		WebhookID: webhookID,
		CreatedBy: createdBy,
	}

	if err := dbengine.Instance().DB.Create(trigger).Error; err != nil {
		return nil, fmt.Errorf("failed to create webhook trigger: %w", err)
	}

	return trigger, nil
}

// GetWebhookTrigger finds a webhook trigger by ID.
func (wp *WebhookProvider) GetWebhookTrigger(webhookID string) *models.WorkflowWebhookTrigger {
	var trigger models.WorkflowWebhookTrigger
	if err := dbengine.Instance().DB.Where("webhook_id = ?", webhookID).First(&trigger).Error; err != nil {
		return nil
	}
	return &trigger
}

// DeleteWebhookTrigger removes a webhook trigger.
func (wp *WebhookProvider) DeleteWebhookTrigger(appID, nodeID string) error {
	return dbengine.Instance().DB.Where("app_id = ? AND node_id = ?", appID, nodeID).Delete(&models.WorkflowWebhookTrigger{}).Error
}

// WebhookURL returns the webhook URL.
func (wp *WebhookProvider) WebhookURL(webhookID string) string {
	return fmt.Sprintf("%s/webhooks/%s", wp.baseURL, webhookID)
}

func generateWebhookID() (string, error) {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
