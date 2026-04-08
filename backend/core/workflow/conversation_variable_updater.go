package workflow

import (
	"encoding/json"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

// ConversationVariableUpdater updates conversation variables after workflow execution.
type ConversationVariableUpdater struct{}

// UpdateFromExecution updates conversation variables based on workflow execution results.
func (u *ConversationVariableUpdater) UpdateFromExecution(
	appID, conversationID string,
	variableUpdates map[string]any,
) error {
	if conversationID == "" || len(variableUpdates) == 0 {
		return nil
	}

	var existing models.ConversationVariable
	err := dbengine.Instance().DB.Where("app_id = ? AND conversation_id = ?", appID, conversationID).First(&existing).Error

	updatedData, _ := json.Marshal(variableUpdates)

	if err != nil {
		// Create new
		cv := &models.ConversationVariable{
			AppID:          appID,
			ConversationID: conversationID,
			Data:           string(updatedData),
		}
		return dbengine.Instance().DB.Create(cv).Error
	}

	// Merge with existing
	var existingData map[string]any
	json.Unmarshal([]byte(existing.Data), &existingData)
	if existingData == nil {
		existingData = make(map[string]any)
	}
	for k, v := range variableUpdates {
		existingData[k] = v
	}
	mergedData, _ := json.Marshal(existingData)
	return dbengine.Instance().DB.Model(&existing).Update("data", string(mergedData)).Error
}

// GetConversationVariables retrieves conversation variables for a conversation.
func (u *ConversationVariableUpdater) GetConversationVariables(appID, conversationID string) map[string]any {
	var cv models.ConversationVariable
	if err := dbengine.Instance().DB.Where("app_id = ? AND conversation_id = ?", appID, conversationID).First(&cv).Error; err != nil {
		return nil
	}
	var data map[string]any
	json.Unmarshal([]byte(cv.Data), &data)
	return data
}

// ClearConversationVariables removes conversation variables.
func (u *ConversationVariableUpdater) ClearConversationVariables(appID, conversationID string) error {
	return dbengine.Instance().DB.Where("app_id = ? AND conversation_id = ?", appID, conversationID).Delete(&models.ConversationVariable{}).Error
}

// UpdateVariableAssignment handles variable_assigner node outputs.
func (u *ConversationVariableUpdater) UpdateVariableAssignment(
	appID, conversationID string,
	assignments []map[string]any,
) error {
	updates := make(map[string]any)
	for _, a := range assignments {
		name, _ := a["variable"].(string)
		value := a["value"]
		if name != "" {
			updates[name] = value
		}
	}
	if len(updates) == 0 {
		return nil
	}
	return u.UpdateFromExecution(appID, conversationID, updates)
}
