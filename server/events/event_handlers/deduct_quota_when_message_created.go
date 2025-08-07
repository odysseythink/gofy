package eventhandlers

import (
	"gorm.io/gorm"
	dbengine "mlib.com/gofy/server/db_engine"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	"mlib.com/gofy/server/models"
)

// @message_was_created.connect
func DeductQuotaWhenMessageCreatedHandle(message *models.Message, application_generate_entity any) {
	var model_config *appconfigentities.ModelConfigWithCredentialsEntity
	var tenant_id string
	switch real_application_generate_entity := application_generate_entity.(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		model_config = real_application_generate_entity.ModelConf
		tenant_id = real_application_generate_entity.AppConfig.TenantID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		model_config = real_application_generate_entity.ModelConf
		tenant_id = real_application_generate_entity.AppConfig.TenantID
	default:
		return
	}
	provider_model_bundle := model_config.ProviderModelBundle
	provider_configuration := provider_model_bundle.Configuration

	if provider_configuration.UsingProviderType != providerenumtypes.Provider_SYSTEM {
		return
	}
	system_configuration := provider_configuration.SystemConfiguration

	var quota_unit providerenumtypes.QuotaUnitType
	for _, quota_configuration := range system_configuration.QuotaConfigurations {
		if quota_configuration.QuotaType == system_configuration.CurrentQuotaType {
			quota_unit = quota_configuration.QuotaUnit
			if quota_configuration.QuotaLimit == -1 {
				return
			}
			break
		}
	}
	var used_quota int
	if string(quota_unit) != "" {
		if quota_unit == providerenumtypes.QuotaUnit_TOKENS {
			used_quota = message.MessageTokens + message.AnswerTokens
		} else if quota_unit == providerenumtypes.QuotaUnit_CREDITS {
			used_quota = 1
		} else {
			used_quota = 1
		}
	}
	if used_quota != 0 && string(system_configuration.CurrentQuotaType) != "" {
		dbengine.Instance().DB.Model(&models.Provider{}).UpdateColumn("quota_used", gorm.Expr("quota_used + ?", used_quota)).Where("tenant_id = ? and provider_name = ? and provider_type = ? and quota_type = ? and quota_limit > quota_used", tenant_id, model_config.Provider, providerenumtypes.Provider_SYSTEM, system_configuration.CurrentQuotaType)
	}
}
