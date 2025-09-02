package advancedchat

import (
	"iter"
	"strings"

	uuid "github.com/satori/go.uuid"
	achatconfigmgr "mlib.com/gofy/server/core/app/config_manageres/advanced_chat"
	acresponseconverter "mlib.com/gofy/server/core/app/generator_response_convertes/advanced_chat"
	actaskpipeline "mlib.com/gofy/server/core/app/generator_task_pipelines/advanced_chat"
	msggenerator "mlib.com/gofy/server/core/app/generatores/message_based"
	msgqueuemgr "mlib.com/gofy/server/core/app/queue_manager/message"
	acrunner "mlib.com/gofy/server/core/app/runner/advanced_chat"
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	conversationmgr "mlib.com/gofy/server/core/manageres/conversation"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	dbengine "mlib.com/gofy/server/db_engine"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appresponseentities "mlib.com/gofy/server/entities/app/response"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

type AdvancedChatAppGenerator[T1 interface {
	*models.Account | *models.EndUser
}] struct {
	*msggenerator.MessageBasedAppGenerator[*appgeneratorentities.AdvancedChatAppGenerateEntity]
	_dialogue_count int
	config_manager  *achatconfigmgr.AdvancedChatAppConfigManager
}

func New[T1 interface {
	*models.Account | *models.EndUser
}]() *AdvancedChatAppGenerator[T1] {
	return &AdvancedChatAppGenerator[T1]{
		MessageBasedAppGenerator: msggenerator.New[*appgeneratorentities.AdvancedChatAppGenerateEntity](),
		config_manager:           achatconfigmgr.New(),
	}
}

func (generator *AdvancedChatAppGenerator[T1]) _generate_worker(
	application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity,
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

		// chatbot app
		runner := acrunner.New(
			application_generate_entity,
			queue_manager,
			conversation,
			message,
			generator._dialogue_count,
		)

		runner.Run()
	}()
}

func (generator *AdvancedChatAppGenerator[T1]) _handle_advanced_chat_response(
	application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity,
	wf *models.Workflow,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
	conversation *models.Conversation,
	message *models.Message,
	user T1,
	stream bool,
) (*appresponseentities.ChatbotAppBlockingResponse, iter.Seq[*appresponseentities.ChatbotAppStreamResponse]) {
	/*
		Handle response.
		:param application_generate_entity: application generate entity
		:param workflow: workflow
		:param queue_manager: queue manager
		:param conversation: conversation
		:param message: message
		:param user: account or end user
		:param stream: is stream
		:return:
	*/
	// init generate task pipeline
	generate_task_pipeline := actaskpipeline.New(
		application_generate_entity,
		wf,
		queue_manager,
		conversation,
		message,
		user,
		stream,
		generator._dialogue_count,
	)

	return generate_task_pipeline.Process()
}

func (generator *AdvancedChatAppGenerator[T1]) Generate(
	app_model *models.App,
	wf *models.Workflow,
	user T1,
	args map[string]any,
	invoke_from appenumtypes.InvokeFrom,
	streaming bool, /* = true*/
) (map[string]any, iter.Seq[string]) {
	/*
		Generate App response.

		:param app_model *models.App
		:param wf *models.Workflow
		:param user: account or end user
		:param args: request args
		:param invoke_from: invoke from source
		:param stream: is stream
	*/
	if _, ok := args["query"]; !ok {
		panic(exceptions.NewValueError("query is required"))
	}
	if _, ok := args["query"].(string); !ok {
		panic(exceptions.NewValueError("query must be a string"))
	}
	if _, ok := args["inputs"]; !ok {
		panic(exceptions.NewValueError("inputs is required"))
	}
	if _, ok := args["inputs"].(map[string]any); !ok {
		panic(exceptions.NewValueError("inputs must be dict"))
	}

	query := args["query"].(string)
	query = strings.ReplaceAll(query, "\x00", "")
	inputs := args["inputs"].(map[string]any)
	auto_generate_name := false
	if _, ok := args["auto_generate_name"]; ok {
		if _, ok := args["auto_generate_name"].(bool); ok {
			auto_generate_name = args["auto_generate_name"].(bool)
		}
	}
	extras := map[string]any{"auto_generate_conversation_name": auto_generate_name}

	// get conversation
	var conversation *models.Conversation
	conversation_id := ""
	if _, ok := args["conversation_id"]; ok {
		if _, ok := args["conversation_id"].(string); ok {
			conversation_id = args["conversation_id"].(string)
		}
	}
	if conversation_id != "" {
		conversation = conversationmgr.GetConversationByUser(
			app_model, conversation_id, user,
		)
	}
	// // parse files
	// files = args["files"] if args.get("files") else []
	// file_extra_config = FileUploadConfigManager.convert(wf.features_dict, is_vision=false)
	// if file_extra_config:
	// 	file_objs = file_factory.build_from_mappings(
	// 		mappings=files,
	// 		tenant_id=app_model.tenant_id,
	// 		config=file_extra_config,
	// 	)
	// else:
	// 	file_objs = []

	// convert to app config
	app_config := generator.config_manager.GetAppConfig(app_model, wf)

	if invoke_from == appenumtypes.InvokeFrom_DEBUGGER {
		// always enable retriever resource in debugger mode
		app_config.AdditionalFeatures.ShowRetrieveSource = true
	}
	parent_message_id := ""
	if invoke_from != appenumtypes.InvokeFrom_SERVICE_API {
		if _, ok := args["parent_message_id"]; ok {
			if _, ok := args["parent_message_id"].(string); ok {
				parent_message_id = args["parent_message_id"].(string)
			}
		}
	}
	user_id := ""
	switch real_user := any(user).(type) {
	case *models.Account:
		user_id = real_user.ID
	case *models.EndUser:
		user_id = real_user.ID
	}
	workflow_run_id := uuid.NewV4().String()
	// init application generate entity
	application_generate_entity := &appgeneratorentities.AdvancedChatAppGenerateEntity{
		ConversationAppGenerateEntity: &appgeneratorentities.ConversationAppGenerateEntity[*appconfigentities.AdvancedChatAppConfig]{
			AppGenerateEntity: &appgeneratorentities.AppGenerateEntity[*appconfigentities.AdvancedChatAppConfig]{
				TaskID:     uuid.NewV4().String(),
				AppConfig:  app_config,
				Inputs:     map[string]any{},
				UserID:     user_id,
				Stream:     streaming,
				InvokeFrom: invoke_from,
				Extras:     extras,
			},
			ParentMessageID: parent_message_id,
		},
		WorkflowRunID: workflow_run_id,
		Query:         query,
	}
	if conversation != nil {
		application_generate_entity.ConversationID = conversation.ID
		application_generate_entity.Inputs = conversation.Inputs()
	} else {
		application_generate_entity.Inputs = generator.PrepareUserInputs(inputs, app_config.Variables, app_config.TenantID)
	}

	return generator._generate(
		wf,
		user,
		invoke_from,
		application_generate_entity,
		conversation,
		streaming,
	)
}
func (generator *AdvancedChatAppGenerator[T1]) single_iteration_generate(
	app_model *models.App, wf *models.Workflow, node_id string, user T1, args map[string]any, streaming bool, /* = true*/
) (map[string]any, iter.Seq[string]) {
	/*
		Generate App response.

		:param app_model *models.App
		:param wf *models.Workflow
		:param user: account or end user
		:param args: request args
		:param invoke_from: invoke from source
		:param stream: is stream
	*/
	if node_id == "" {
		panic(exceptions.NewValueError("node_id is required"))
	}
	if _, ok := args["inputs"]; !ok {
		panic(exceptions.NewValueError("inputs is required"))
	}
	if _, ok := args["inputs"].(map[string]any); !ok {
		panic(exceptions.NewValueError("inputs must be dict"))
	}
	user_id := ""
	switch real_user := any(user).(type) {
	case *models.Account:
		user_id = real_user.ID
	case *models.EndUser:
		user_id = real_user.ID
	}

	// convert to app config
	app_config := generator.config_manager.GetAppConfig(app_model, wf)

	// init application generate entity
	application_generate_entity := &appgeneratorentities.AdvancedChatAppGenerateEntity{
		ConversationAppGenerateEntity: &appgeneratorentities.ConversationAppGenerateEntity[*appconfigentities.AdvancedChatAppConfig]{
			AppGenerateEntity: &appgeneratorentities.AppGenerateEntity[*appconfigentities.AdvancedChatAppConfig]{
				TaskID:     uuid.NewV4().String(),
				AppConfig:  app_config,
				Inputs:     map[string]any{},
				UserID:     user_id,
				Stream:     streaming,
				InvokeFrom: appenumtypes.InvokeFrom_DEBUGGER,
				Extras:     map[string]any{"auto_generate_conversation_name": false},
			},
		},
		Query: "",
		SingleIterationRun: &appgeneratorentities.SingleIterationRunEntity{
			NodeID: node_id,
			Inputs: args["inputs"].(map[string]any),
		},
	}
	return generator._generate(
		wf,
		user,
		appenumtypes.InvokeFrom_DEBUGGER,
		application_generate_entity,
		nil,
		streaming,
	)
}
func (generator *AdvancedChatAppGenerator[T1]) _generate(
	wf *models.Workflow,
	user T1,
	invoke_from appenumtypes.InvokeFrom,
	application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity,
	conversation *models.Conversation,
	stream bool, /* = true*/
) (map[string]any, iter.Seq[string]) {
	/*
		Generate App response.

		:param wf *models.Workflow
		:param user: account or end user
		:param invoke_from: invoke from source
		:param application_generate_entity: application generate entity
		:param conversation: conversation
		:param stream: is stream
	*/
	is_first_conversation := false
	if conversation == nil {
		is_first_conversation = true
	}
	// init generate records
	conversation, message := generator.InitGenerateRecords(application_generate_entity, conversation)
	mlog.Debugf("------message=%#v", message)

	if is_first_conversation {
		// update conversation features
		conversation.OverrideModelConfigsStr = wf.FeaturesStr
		dbengine.Instance().DB.Save(conversation)
	}
	// get conversation dialogue count
	generator._dialogue_count = promptutils.GetThreadMessagesLength(conversation.ID)

	// init queue manager
	queue_manager := msgqueuemgr.New(
		application_generate_entity.TaskID,
		application_generate_entity.UserID,
		application_generate_entity.InvokeFrom,
		conversation.ID,
		string(conversation.Mode),
		message.ID,
	)

	// new thread
	go func() {
		generator._generate_worker(application_generate_entity, queue_manager, conversation.ID, message.ID)
		mlog.Debugf("_generate worker exit")
	}()

	// return response or stream generator
	blocking_response, stream_response := generator._handle_advanced_chat_response(
		application_generate_entity,
		wf,
		queue_manager,
		conversation,
		message,
		user,
		stream,
	)
	converter := (&acresponseconverter.AdvancedChatAppGeneratorResponseConvert{})
	if blocking_response != nil {
		return converter.ConvertBlocking(converter, blocking_response, invoke_from), nil
	}
	if stream_response != nil {
		return nil, converter.ConvertStream(converter, stream_response, invoke_from)
	}
	panic(exceptions.NewValueError("_handle_advanced_chat_response contains no response"))
}
