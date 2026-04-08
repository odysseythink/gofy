package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type HumanInputService struct{}

func (s *HumanInputService) GetFormByToken(formToken string) (*models.HumanInputFormRecipient, *models.HumanInputForm, error) {
	var recipient models.HumanInputFormRecipient
	if err := dbengine.Instance().DB.Where("access_token = ?", formToken).First(&recipient).Error; err != nil {
		return nil, nil, fmt.Errorf("form not found")
	}
	var form models.HumanInputForm
	if err := dbengine.Instance().DB.Where("id = ?", recipient.FormID).First(&form).Error; err != nil {
		return nil, nil, fmt.Errorf("form not found")
	}
	return &recipient, &form, nil
}

func (s *HumanInputService) SubmitForm(formID string, recipientID string, selectedActionID string, formData string) error {
	var form models.HumanInputForm
	if err := dbengine.Instance().DB.Where("id = ?", formID).First(&form).Error; err != nil {
		return fmt.Errorf("form not found")
	}
	if form.Status != "waiting" {
		return fmt.Errorf("form is not in waiting status")
	}
	if time.Now().After(form.ExpirationTime) {
		return fmt.Errorf("form has expired")
	}
	now := time.Now()
	return dbengine.Instance().DB.Model(&form).Updates(map[string]any{
		"status":                    "submitted",
		"selected_action_id":        selectedActionID,
		"submitted_data":            formData,
		"submitted_at":              &now,
		"completed_by_recipient_id": recipientID,
	}).Error
}

func (s *HumanInputService) GetFormsByWorkflowRun(workflowRunID string) []*models.HumanInputForm {
	var forms []*models.HumanInputForm
	dbengine.Instance().DB.Where("workflow_run_id = ?", workflowRunID).Order("created_at DESC").Find(&forms)
	return forms
}

func (s *HumanInputService) IsFormExpired(form *models.HumanInputForm) bool {
	return time.Now().After(form.ExpirationTime)
}

func (s *HumanInputService) TimeoutExpiredForms() int64 {
	result := dbengine.Instance().DB.Model(&models.HumanInputForm{}).
		Where("status = ? AND expiration_time < ?", "waiting", time.Now()).
		Update("status", "timeout")
	return result.RowsAffected
}

// CreateForm creates a new human input form for a workflow node.
func (s *HumanInputService) CreateForm(
	tenantID, appID, workflowRunID, nodeID string,
	formDefinition string,
	renderedContent string,
	expirationMinutes int,
) (*models.HumanInputForm, error) {
	expiration := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	form := &models.HumanInputForm{
		ID:              uuid.NewV4().String(),
		TenantID:        tenantID,
		AppID:           appID,
		WorkflowRunID:   &workflowRunID,
		FormKind:        "runtime",
		NodeID:          nodeID,
		FormDefinition:  formDefinition,
		RenderedContent: renderedContent,
		Status:          "waiting",
		ExpirationTime:  expiration,
	}
	if err := dbengine.Instance().DB.Create(form).Error; err != nil {
		return nil, err
	}
	return form, nil
}

// CreateDelivery creates a delivery record for a form.
func (s *HumanInputService) CreateDelivery(
	formID string,
	deliveryMethodType string,
	channelPayload string,
) (*models.HumanInputDelivery, error) {
	delivery := &models.HumanInputDelivery{
		ID:                 uuid.NewV4().String(),
		FormID:             formID,
		DeliveryMethodType: deliveryMethodType,
		ChannelPayload:     channelPayload,
	}
	if err := dbengine.Instance().DB.Create(delivery).Error; err != nil {
		return nil, err
	}
	return delivery, nil
}

// CreateRecipient creates a recipient for a form delivery.
func (s *HumanInputService) CreateRecipient(
	formID, deliveryID, recipientType, recipientPayload string,
) (*models.HumanInputFormRecipient, error) {
	token := generateAccessToken()
	recipient := &models.HumanInputFormRecipient{
		ID:               uuid.NewV4().String(),
		FormID:           formID,
		DeliveryID:       deliveryID,
		RecipientType:    recipientType,
		RecipientPayload: recipientPayload,
		AccessToken:      &token,
	}
	if err := dbengine.Instance().DB.Create(recipient).Error; err != nil {
		return nil, err
	}
	return recipient, nil
}

// GetFormByID retrieves a form by ID.
func (s *HumanInputService) GetFormByID(formID string) *models.HumanInputForm {
	var form models.HumanInputForm
	if err := dbengine.Instance().DB.Where("id = ?", formID).First(&form).Error; err != nil {
		return nil
	}
	return &form
}

// GetFormByWorkflowRunAndNode retrieves a form by workflow run and node.
func (s *HumanInputService) GetFormByWorkflowRunAndNode(workflowRunID, nodeID string) *models.HumanInputForm {
	var form models.HumanInputForm
	if err := dbengine.Instance().DB.Where("workflow_run_id = ? AND node_id = ?", workflowRunID, nodeID).First(&form).Error; err != nil {
		return nil
	}
	return &form
}

// GetFormRecipients returns all recipients for a form.
func (s *HumanInputService) GetFormRecipients(formID string) []*models.HumanInputFormRecipient {
	var recipients []*models.HumanInputFormRecipient
	dbengine.Instance().DB.Where("form_id = ?", formID).Find(&recipients)
	return recipients
}

// GetFormDeliveries returns all deliveries for a form.
func (s *HumanInputService) GetFormDeliveries(formID string) []*models.HumanInputDelivery {
	var deliveries []*models.HumanInputDelivery
	dbengine.Instance().DB.Where("form_id = ?", formID).Find(&deliveries)
	return deliveries
}

func generateAccessToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
