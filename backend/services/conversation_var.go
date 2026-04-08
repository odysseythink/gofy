package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type ConversationVariableService struct{}

func (s *ConversationVariableService) GetConversationVariables(appID, conversationID string) []*models.ConversationVariable {
	var vars []*models.ConversationVariable
	dbengine.Instance().DB.Where("app_id = ? AND conversation_id = ?", appID, conversationID).Find(&vars)
	return vars
}

func (s *ConversationVariableService) UpdateConversationVariable(appID, conversationID string, data string) error {
	return dbengine.Instance().DB.Model(&models.ConversationVariable{}).
		Where("app_id = ? AND conversation_id = ?", appID, conversationID).
		Update("data", data).Error
}
