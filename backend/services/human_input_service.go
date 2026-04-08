package services

import (
	"fmt"
	"time"

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
