package agentchat

import (
	"slices"

	uuid "github.com/satori/go.uuid"
	"github.com/odysseythink/gofy/backend/core/app/runner/base"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	modelmanager "github.com/odysseythink/gofy/backend/core/manageres/model_manager"
	modelmgr "github.com/odysseythink/gofy/backend/core/manageres/model_manager"
	"github.com/odysseythink/gofy/backend/core/memory"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appgeneratorentities "github.com/odysseythink/gofy/backend/entities/app/generator"
	appqueueentities "github.com/odysseythink/gofy/backend/entities/app/queue"
	agententities "github.com/odysseythink/gofy/backend/entities/agent"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type AgentChatAppRunner struct {
	*base.AppRunner[*appqueueentities.MessageQueueMessage]
}
func (r *AgentChatAppRunner[T]) DirectOutput(
        queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
        app_generate_entity *appgeneratorentities.AgentChatAppGenerateEntity,
        prompt_messages []modelruntimeentities.PromptMessager,
        text string,
        stream bool,
        usage *modelruntimeentities.LLMUsage,
    ) {
        /*
        Direct output
        :param queue_manager: application queue manager
        :param app_generate_entity: app generate entity
        :param prompt_messages: prompt messages
        :param text: text
        :param stream: stream
        :param usage: usage
        :return:
        */
        if stream{
            index := 0
            for _, token := range text{
                chunk := &modelruntimeentities.LLMResultChunk{
                    Model:app_generate_entity.ModelConf.Model,
                    PromptMessages: prompt_messages,
                    Delta: &modelruntimeentities.LLMResultChunkDelta{Index:index, Message:modelruntimeentities.NewAssistantPromptMessage(string(token), "", nil)},
                }

                queue_manager.Publish(&appqueueentities.QueueLLMChunkEvent{Chunk:chunk}, appenumtypes.PublishFrom_APPLICATION_MANAGER)
                index += 1
                time.sleep(0.01)
				}
}
	if usage == nil {
		usage = modelruntimeentities.NewLLMUsage()
	}
        queue_manager.Publish(
            &appqueueentities.QueueMessageEndEvent{
                LLMResult: &modelruntimeentities.LLMResult{
ID               :uuid.NewV4().String(),
Model            :app_generate_entity.ModelConf.Model,
PromptMessages   :prompt_messages,
Message          :modelruntimeentities.NewAssistantPromptMessage(text, "", nil),
Usage            :usage,
                },
            },
            appenumtypes.PublishFrom_APPLICATION_MANAGER,
        )
}
func (runner *AgentChatAppRunner) Run() (
	application_generate_entity *appgeneratorentities.AgentChatAppGenerateEntity,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
	conversation *models.Conversation,
	message *models.Message,
) {
	/*
	   Run assistant application
	   :param application_generate_entity: application generate entity
	   :param queue_manager: application queue manager
	   :param conversation: conversation
	   :param message: message
	   :return
	*/
	app_config := application_generate_entity.AppConfig

	app_record := new(models.App)
	err := dbengine.Instance().DB.Where(&models.App{}).Where("id = ?", app_config.AppID).First(app_record).Error
	if err != nil {
		mlog.Errorf("get app failed:%v", err)
		app_record = nil
	}
	if app_record == nil {
		panic(exceptions.NewValueError("App not found"))
	}

	inputs := map[string]string{}
	for k, v := range application_generate_entity.Inputs {
		if _, ok := v.(string); ok {
			inputs[k] = v.(string)
		}
	}
	query := application_generate_entity.Query
	files := application_generate_entity.Files

	var mem *memory.TokenBufferMemory
	if application_generate_entity.ConversationID != "" {
		// get memory of conversation (read-only)
		model_instance := modelmgr.NewModelInstance(
			application_generate_entity.ModelConf.ProviderModelBundle,
			application_generate_entity.ModelConf.Model,
		)
		mem = memory.NewTokenBufferMemory(conversation, model_instance)
	}
	// organize all inputs and template to prompt messages
	// Include: prompt template, inputs, query(optional), files(optional)
	//          memory(optional)
	prompt_messages, _ := runner.OrganizePromptMessages(
		app_record,
		application_generate_entity.ModelConf,
		&app_config.PromptTemplate,
		inputs,
		files,
		query,
		"",
		mem,
	)

	if query != ""{
	    // annotation reply
	    annotation_reply := runner.QueryAppAnnotationsToReply(
	        app_record,
	        message,
	        query,
	        application_generate_entity.UserID,
	        application_generate_entity.InvokeFrom,
	    )

	    if annotation_reply != nil {
	        queue_manager.Publish(
	            &appqueueentities.QueueAnnotationReplyEvent{MessageAnnotationID:annotation_reply.ID},
	            appenumtypes.PublishFrom_APPLICATION_MANAGER,
	        )

	        runner.DirectOutput(
	            queue_manager,
	            application_generate_entity,
	            prompt_messages,
	            annotation_reply.Content,
	            application_generate_entity.Stream,
	        )
	        return
		}
	}
	// fill in variable inputs from external data tools if exists
	external_data_tools := app_config.ExternalDataVariables
	if len(external_data_tools) > 0{
	    inputs = runner.FillInInputsFromExternalDataTools(
	        app_record.TenantID,
	        app_record.ID,
	        external_data_tools,
	        inputs,
	        query,
	    )
	}
	// reorganize all inputs and template to prompt messages
	// Include: prompt template, inputs, query(optional), files(optional)
	//          memory(optional), external data, dataset context(optional)
	// prompt_messages, _ = runner.OrganizePromptMessages(
	//     app_record,
	//     application_generate_entity.ModelConf,
	//     &app_config.PromptTemplate,
	//     inputs,
	//     files,
	//     query,
	//     mem,
	// )

	// check hosting moderation
	// hosting_moderation_result = runner.check_hosting_moderation(
	//     application_generate_entity=application_generate_entity,
	//     queue_manager=queue_manager,
	//     prompt_messages=prompt_messages,
	// )

	// if hosting_moderation_result{
	//     return

	agent_entity := app_config.Agent
	if agent_entity == nil {
				mlog.Errorf("agent_entity can't be nil")
		panic(exceptions.NewValueError("agent_entity can't be nil"))
	}

	// init model instance
	model_instance := modelmanager.NewModelInstance(
	    application_generate_entity.ModelConf.ProviderModelBundle,
	    application_generate_entity.ModelConf.Model,
	)
	prompt_messages, _ = runner.OrganizePromptMessages(
	    app_record,
	    application_generate_entity.ModelConf,
	    &app_config.PromptTemplate,
	    inputs,
	    files,
	    query,
	    mem,
	)

	// change function call strategy based on LLM model
	// llm_model = cast(LargeLanguageModel, model_instance.model_type_instance)
	model_schema := model_instance.ModelTypeInstance.GetModelSchema(model_instance.ModelTypeInstance, model_instance.Model, model_instance.Credentials)
	if  model_schema == nil{
		mlog.Errorf("Model schema not found")
	    panic(exceptions.NewValueError("Model schema not found"))
	}
	intersection_list := []modelruntimeenumtypes.ModelFeature{}
	for _, v := range model_schema.Features {
		if slices.Contains([]modelruntimeenumtypes.ModelFeature{modelruntimeenumtypes.ModelFeature_MULTI_TOOL_CALL, modelruntimeenumtypes.ModelFeature_TOOL_CALL}, v) {
			intersection_list = append(intersection_list, v)
		}
	}
	slices.Sort(intersection_list)
	intersection_list = slices.Compact(intersection_list)

	if len(intersection_list) > 0{
	    agent_entity.Strategy = agententities.AgentEntityStrategy_FUNCTION_CALLING
	}

	conversation_result := new(models.Conversation)
	err := dbengine.Instance().DB.Model(&models.Conversation{}).Where("id =?", conversation.ID).First(conversation_result).Error
	if err != nil {
				mlog.Errorf("get Conversation failed:%v", err)
	    conversation_result = nil
	}
	if conversation_result == nil {
	    panic(exceptions.NewValueError("Conversation not found"))
	}
		message_result := new(models.Message)
	err = dbengine.Instance().DB.Model(&models.Message{}).Where("id =?", message.ID).First(message_result).Error
	if err != nil {
				mlog.Errorf("get Message failed:%v", err)
	    message_result = nil
	}
	if message_result  == nil {
	    panic(exceptions.NewValueError("Message not found"))
	}


	runner_cls: type[FunctionCallAgentRunner] | type[CotChatAgentRunner] | type[CotCompletionAgentRunner]
	// start agent runner
	if agent_entity.strategy == AgentEntity.Strategy.CHAIN_OF_THOUGHT{
	    // check LLM mode
	    if model_schema.model_properties.get(ModelPropertyKey.MODE) == LLMMode.CHAT.value{
	        runner_cls = CotChatAgentRunner
	    elif model_schema.model_properties.get(ModelPropertyKey.MODE) == LLMMode.COMPLETION.value{
	        runner_cls = CotCompletionAgentRunner
	    else{
	        raise ValueError(f"Invalid LLM mode: {model_schema.model_properties.get(ModelPropertyKey.MODE)}")
	elif agent_entity.strategy == AgentEntity.Strategy.FUNCTION_CALLING{
	    runner_cls = FunctionCallAgentRunner
	else{
	    raise ValueError(f"Invalid agent strategy: {agent_entity.strategy}")

	runner = runner_cls(
	    tenant_id=app_config.tenant_id,
	    application_generate_entity=application_generate_entity,
	    conversation=conversation_result,
	    app_config=app_config,
	    model_config=application_generate_entity.model_conf,
	    config=agent_entity,
	    queue_manager=queue_manager,
	    message=message_result,
	    user_id=application_generate_entity.user_id,
	    memory=memory,
	    prompt_messages=prompt_message,
	    model_instance=model_instance,
	)

	invoke_result = runner.run(
	    message=message,
	    query=query,
	    inputs=inputs,
	)

	// handle invoke result
	runner._handle_invoke_result(
	    invoke_result=invoke_result,
	    queue_manager=queue_manager,
	    stream=application_generate_entity.stream,
	    agent=True,
	)
}
