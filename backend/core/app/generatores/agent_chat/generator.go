package agentchat

import (
	achatconfigmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/agent_chat"
	msggenerator "github.com/odysseythink/gofy/backend/core/app/generatores/message_based"
	appgeneratorentities "github.com/odysseythink/gofy/backend/entities/app/generator"
	"github.com/odysseythink/gofy/backend/models"
)

type AgentChatAppGenerator[T1 interface {
	*models.Account | *models.EndUser
}] struct {
	*msggenerator.MessageBasedAppGenerator[*appgeneratorentities.AgentChatAppGenerateEntity]
	config_manager *achatconfigmgr.AgentChatAppConfigManager
}

func New[T1 interface {
	*models.Account | *models.EndUser
}]() *AgentChatAppGenerator[T1] {
	return &AgentChatAppGenerator[T1]{
		MessageBasedAppGenerator: msggenerator.New[*appgeneratorentities.AgentChatAppGenerateEntity](),
		config_manager:           achatconfigmgr.New(),
	}
}

func (generator *AdvancedChatAppGenerator[T1]) _generate_worker(
        application_generate_entity *appgeneratorentities.AgentChatAppGenerateEntity,
        queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
        conversation_id string,
        message_id string,
    ) {
        /*
        Generate worker in a new thread.
        :param flask_app: Flask app
        :param application_generate_entity: application generate entity
        :param queue_manager: queue manager
        :param conversation_id: conversation ID
        :param message_id: message ID
        :return:
        */

	func() {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if real_exp, ok := r.(*exceptions.GenerateTaskStoppedError); ok {
					mlog.Error(real_exp)
					return
				} else if _, ok := r.(*modelruntimeexceptions.InvokeAuthorizationError); ok {
					queue_manager.PublishError(queue_manager, modelruntimeexceptions.NewInvokeAuthorizationError("Incorrect API key provided"), appenumtypes.PublishFrom_APPLICATION_MANAGER)
				} else if real_exp, ok := r.(*exceptions.ValidationError); ok {
					mlog.Errorf("Validation Error when generating")
					queue_manager.PublishError(queue_manager, real_exp, appenumtypes.PublishFrom_APPLICATION_MANAGER)
				} else if real_exp, ok := r.(*exceptions.ValueError); ok {
					mlog.Errorf("Error when generating")
					queue_manager.PublishError(queue_manager, real_exp, appenumtypes.PublishFrom_APPLICATION_MANAGER)
				} else if real_exp, ok := r.(error); ok {
					mlog.Errorf("Unknown Error when generating:%v", real_exp)
					queue_manager.PublishError(queue_manager, real_exp, appenumtypes.PublishFrom_APPLICATION_MANAGER)
				} else {
					panic(r)
				}
			}
		}()
                // get conversation and message
                conversation := generator.GetConversation(conversation_id)
                message := generator.GetMessage(message_id)
	}()



                // chatbot app
                runner = AgentChatAppRunner()
                runner.run(
                    application_generate_entity=application_generate_entity,
                    queue_manager=queue_manager,
                    conversation=conversation,
                    message=message,
                )

				}