package eventhandlers

import (
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appgeneratorentities "github.com/odysseythink/gofy/backend/entities/app/generator"
	"github.com/odysseythink/gofy/backend/models"
)

// @message_was_created.connect
func UpdateProviderLastUsedAtWhenMessageCreatedHandle(message *models.Message, application_generate_entity any) {
	var tenant_id string
	var provider string
	switch real_application_generate_entity := application_generate_entity.(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		provider = real_application_generate_entity.ModelConf.Provider
		tenant_id = real_application_generate_entity.AppConfig.TenantID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		provider = real_application_generate_entity.ModelConf.Provider
		tenant_id = real_application_generate_entity.AppConfig.TenantID
	default:
		return
	}

	dbengine.Instance().DB.Model(&models.Provider{}).UpdateColumn("last_used", time.Now()).Where("tenant_id = ? and provider_name = ? and provider_type = ?", tenant_id, provider)
}
