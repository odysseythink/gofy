package services

import (
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	uuid "github.com/satori/go.uuid"
)

type SavedMessageService struct{}

func (s *SavedMessageService) GetSavedMessages(appID, userID string, page, limit int) ([]*models.SavedMessage, int64) {
	var messages []*models.SavedMessage
	var total int64
	query := dbengine.Instance().DB.Where("app_id = ? AND created_by = ?", appID, userID)
	query.Model(&models.SavedMessage{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&messages)
	return messages, total
}

func (s *SavedMessageService) SaveMessage(appID, userID, messageID string) (*models.SavedMessage, error) {
	// Check if already saved
	var existing models.SavedMessage
	if err := dbengine.Instance().DB.Where("app_id = ? AND message_id = ? AND created_by = ?", appID, messageID, userID).First(&existing).Error; err == nil {
		return &existing, nil
	}
	saved := &models.SavedMessage{
		ID:        uuid.NewV4().String(),
		AppID:     appID,
		MessageID: messageID,
		CreatedBy: userID,
	}
	if err := dbengine.Instance().DB.Create(saved).Error; err != nil {
		return nil, err
	}
	return saved, nil
}

func (s *SavedMessageService) DeleteSavedMessage(appID, userID, messageID string) error {
	result := dbengine.Instance().DB.Where("app_id = ? AND message_id = ? AND created_by = ?", appID, messageID, userID).Delete(&models.SavedMessage{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("saved message not found")
	}
	return nil
}
