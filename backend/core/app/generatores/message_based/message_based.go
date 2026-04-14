package messagebased

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	easyuigeneratortaskpipeline "mlib.com/gofy/server/core/app/generator_task_pipelines/easy_ui"
	"mlib.com/gofy/server/core/app/generatores/base"
	"mlib.com/gofy/server/core/exceptions"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	"mlib.com/gofy/server/models"
)

type MessageBasedAppGenerator[T1 interface {
	*appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity | *appgeneratorentities.AdvancedChatAppGenerateEntity
}] struct {
	*base.BaseAppGenerator
}

func New[T1 interface {
	*appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity | *appgeneratorentities.AdvancedChatAppGenerateEntity
}]() *MessageBasedAppGenerator[T1] {
	return &MessageBasedAppGenerator[T1]{
		BaseAppGenerator: base.New(),
	}
}

func (generator *MessageBasedAppGenerator[T1]) HandleResponse(
	application_generate_entity T1,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
	conversation *models.Conversation,
	message *models.Message,
	stream bool,
) any { /*-> Union[
		ChatbotAppBlockingResponse,
		CompletionAppBlockingResponse,
		Generator[Union[ChatbotAppStreamResponse, CompletionAppStreamResponse], None, None],
	]*/
	/*
		Handle response.
		:param application_generate_entity: application generate entity
		:param queue_manager: queue manager
		:param conversation: conversation
		:param message: message
		:param user: user
		:param stream: is stream
		:return:
	*/
	// init generate task pipeline

	switch real_app_generate_entity := any(application_generate_entity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		generate_task_pipeline := easyuigeneratortaskpipeline.New(
			real_app_generate_entity,
			queue_manager,
			conversation,
			message,
			stream,
		)

		return generate_task_pipeline.Process()
	case *appgeneratorentities.CompletionAppGenerateEntity:
		generate_task_pipeline := easyuigeneratortaskpipeline.New(
			real_app_generate_entity,
			queue_manager,
			conversation,
			message,
			stream,
		)

		return generate_task_pipeline.Process()
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		generate_task_pipeline := easyuigeneratortaskpipeline.New(
			real_app_generate_entity,
			queue_manager,
			conversation,
			message,
			stream,
		)

		return generate_task_pipeline.Process()
	}
	mlog.Error("advanced_chat generate entity not supported in this function")
	panic(exceptions.NewNotImplementedError("advanced_chat generate entity not supported in this function"))
}
func (generator *MessageBasedAppGenerator[T1]) _get_conversation_introduction(application_generate_entity T1) string {
	/*
		Get conversation introduction
		:param application_generate_entity: application generate entity
		:return: conversation introduction
	*/
	var inputs map[string]any
	introduction := ""
	switch real_entity := any(application_generate_entity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		inputs = real_entity.Inputs
		introduction = real_entity.AppConfig.AdditionalFeatures.OpeningStatement
	case *appgeneratorentities.CompletionAppGenerateEntity:
		introduction = real_entity.AppConfig.AdditionalFeatures.OpeningStatement
		inputs = real_entity.Inputs
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		introduction = real_entity.AppConfig.AdditionalFeatures.OpeningStatement
		inputs = real_entity.Inputs
	case *appgeneratorentities.AdvancedChatAppGenerateEntity:
		introduction = real_entity.AppConfig.AdditionalFeatures.OpeningStatement
		inputs = real_entity.Inputs
	}

	if introduction != "" {
		// try:
		prompt_template := promptutils.NewPromptTemplateParser(introduction, false)
		prompt_inputs := map[string]string{}
		for _, k := range prompt_template.VariableKeys {
			if _, ok := inputs[k]; ok {
				prompt_inputs[k] = inputs[k].(string)
			}
		}

		introduction = prompt_template.Format(prompt_inputs, false)
		// except KeyError:
		// 	pass
	}
	return introduction
}

func (generator *MessageBasedAppGenerator[T1]) InitGenerateRecords(
	application_generate_entity T1,
	conversation *models.Conversation,
) (*models.Conversation, *models.Message) {
	/*
		Initialize generate records
		:param application_generate_entity: application generate entity
		:conversation conversation
		:return:
	*/
	var invoke_from appenumtypes.InvokeFrom
	user_id := ""
	app_model_config_id := ""
	model_provider := ""
	model_id := ""
	var override_model_configs map[string]any
	var inputs map[string]any
	query := ""
	parent_message_id := ""
	app_id := ""
	var app_mode models.AppMode
	switch real_entity := any(application_generate_entity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		app_config := real_entity.AppConfig
		invoke_from = real_entity.InvokeFrom
		user_id = real_entity.UserID
		model_provider = real_entity.ModelConf.Provider
		model_id = real_entity.ModelConf.Model
		if app_config.AppModelConfigFrom == appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS && slices.Contains([]models.AppMode{models.AppMode_AGENT_CHAT, models.AppMode_CHAT, models.AppMode_COMPLETION}, app_config.AppMode) {
			override_model_configs = app_config.AppModelConfigDict
		}
		app_model_config_id = app_config.AppModelConfigID
		query = real_entity.Query
		parent_message_id = real_entity.ParentMessageID
		app_id = app_config.AppID
		app_mode = app_config.AppMode
	case *appgeneratorentities.CompletionAppGenerateEntity:
		app_config := real_entity.AppConfig
		invoke_from = real_entity.InvokeFrom
		user_id = real_entity.UserID
		model_provider = real_entity.ModelConf.Provider
		model_id = real_entity.ModelConf.Model
		if app_config.AppModelConfigFrom == appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS && slices.Contains([]models.AppMode{models.AppMode_AGENT_CHAT, models.AppMode_CHAT, models.AppMode_COMPLETION}, app_config.AppMode) {
			override_model_configs = app_config.AppModelConfigDict
		}
		app_model_config_id = app_config.AppModelConfigID
		query = real_entity.Query
		app_id = app_config.AppID
		app_mode = app_config.AppMode
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		app_config := real_entity.AppConfig
		invoke_from = real_entity.InvokeFrom
		user_id = real_entity.UserID
		model_provider = real_entity.ModelConf.Provider
		model_id = real_entity.ModelConf.Model
		if app_config.AppModelConfigFrom == appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS && slices.Contains([]models.AppMode{models.AppMode_AGENT_CHAT, models.AppMode_CHAT, models.AppMode_COMPLETION}, app_config.AppMode) {
			override_model_configs = app_config.AppModelConfigDict
		}
		app_model_config_id = app_config.AppModelConfigID
		query = real_entity.Query
		parent_message_id = real_entity.ParentMessageID
		app_id = app_config.AppID
		app_mode = app_config.AppMode
	case *appgeneratorentities.AdvancedChatAppGenerateEntity:
		app_config := real_entity.AppConfig
		invoke_from = real_entity.InvokeFrom
		user_id = real_entity.UserID
		query = real_entity.Query
		parent_message_id = real_entity.ParentMessageID
		app_id = app_config.AppID
		app_mode = app_config.AppMode
	}
	now := time.Now()

	// get from source
	end_user_id := ""
	account_id := ""
	from_source := ""
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_WEB_APP, appenumtypes.InvokeFrom_SERVICE_API}, invoke_from) {
		from_source = "api"
		end_user_id = user_id
	} else {
		from_source = "console"
		account_id = user_id
	}

	// get conversation introduction
	introduction := generator._get_conversation_introduction(application_generate_entity)

	if conversation == nil {
		conversation = &models.Conversation{
			ID:                      uuid.NewV4().String(),
			AppID:                   app_id,
			AppModelConfigID:        app_model_config_id,
			ModelProvider:           model_provider,
			ModelID:                 model_id,
			Mode:                    app_mode,
			Name:                    "New conversation",
			Introduction:            introduction,
			SystemInstruction:       "",
			SystemInstructionTokens: 0,
			Status:                  "normal",
			FromSource:              from_source,
			FromEndUserID:           end_user_id,
			FromAccountID:           account_id,
			InvokeFrom:              string(invoke_from),
			CreatedAt:               &now,
			UpdatedAt:               &now,
		}
		conversation.SetInputs(inputs)
		if len(override_model_configs) > 0 {
			bindata, _ := json.Marshal(override_model_configs)
			conversation.OverrideModelConfigsStr = string(bindata)
		}
		dbengine.Instance().DB.Create(conversation)
	} else {
		now := time.Now()
		conversation.UpdatedAt = &now
		dbengine.Instance().DB.Save(conversation)
	}

	message := &models.Message{
		ID:              uuid.NewV4().String(),
		AppID:           app_id,
		ModelProvider:   model_provider,
		ModelID:         model_id,
		ConversationID:  conversation.ID,
		Query:           query,
		Currency:        "USD",
		FromSource:      from_source,
		FromEndUserID:   end_user_id,
		FromAccountID:   account_id,
		InvokeFrom:      string(invoke_from),
		ParentMessageID: parent_message_id,
		CreatedAt:       &now,
		UpdatedAt:       &now,
	}
	message.SetOverrideModelConfigs(override_model_configs)
	message.SetInputs(inputs)
	dbengine.Instance().DB.Create(message)

	return conversation, message
}

func (generator *MessageBasedAppGenerator[T1]) GetConversation(conversation_id string) *models.Conversation {
	/*
		Get conversation by conversation id
		:param conversation_id: conversation id
		:return: conversation
	*/
	conversation := new(models.Conversation)
	err := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ?", conversation_id).First(conversation).Error
	if err != nil {
		mlog.Errorf("get Conversation failed:%v", err)
		conversation = nil
	}

	if conversation == nil {
		panic(exceptions.NewConversationNotExistsError(""))
	}
	return conversation
}
func (generator *MessageBasedAppGenerator[T1]) GetMessage(message_id string) *models.Message {
	/*
		Get message by message id
		:param message_id: message id
		:return: message
	*/
	message := new(models.Message)
	err := dbengine.Instance().DB.Model(&models.Message{}).Where("id = ?", message_id).First(message).Error
	if err != nil {
		mlog.Errorf("get Message failed:%v", err)
		message = nil
	}

	return message
}
