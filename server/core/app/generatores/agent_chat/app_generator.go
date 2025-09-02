package agentchat

import (
	achatconfigmgr "mlib.com/gofy/server/core/app/config_manageres/advanced_chat"
	msggenerator "mlib.com/gofy/server/core/app/generatores/message_based"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	"mlib.com/gofy/server/models"
)

type AgentChatAppGenerator[T1 interface {
	*models.Account | *models.EndUser
}] struct {
	*msggenerator.MessageBasedAppGenerator[*appgeneratorentities.AgentChatAppGenerateEntity]
	config_manager *achatconfigmgr.AdvancedChatAppConfigManager
}
