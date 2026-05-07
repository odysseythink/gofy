package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
)

type AgentService struct{}

func (s *AgentService) GetAgentLogs(appID, conversationID, messageID string) []*models.MessageAgentThought {
	var thoughts []*models.MessageAgentThought
	query := dbengine.Instance().DB.Where("message_id = ?", messageID)
	if conversationID != "" {
		query = query.Where("conversation_id = ?", conversationID)
	}
	query.Order("position ASC").Find(&thoughts)
	return thoughts
}

func (s *AgentService) ListAgentProviders(tenantID string) []map[string]any {
	// TODO: Query plugin system for agent strategy providers
	return []map[string]any{}
}
