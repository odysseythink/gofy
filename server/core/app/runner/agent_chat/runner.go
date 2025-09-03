package agentchat

import (
	"mlib.com/gofy/server/core/app/runner/base"
	"mlib.com/gofy/server/core/exceptions"
	modelmgr "mlib.com/gofy/server/core/manageres/model_manager"
	"mlib.com/gofy/server/core/memory"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type AgentChatAppRunner struct {
	*base.AppRunner[*appqueueentities.MessageQueueMessage]
}

func (runner *AgentChatAppRunner) Run() (
	application_generate_entity *appgeneratorentities.AgentChatAppGenerateEntity,
	queue_manager *appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
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
	    annotation_reply = runner.query_app_annotations_to_reply(
	        app_record=app_record,
	        message=message,
	        query=query,
	        user_id=application_generate_entity.user_id,
	        invoke_from=application_generate_entity.invoke_from,
	    )

	    if annotation_reply{
	        queue_manager.publish(
	            QueueAnnotationReplyEvent(message_annotation_id=annotation_reply.id),
	            PublishFrom.APPLICATION_MANAGER,
	        )

	        runner.direct_output(
	            queue_manager=queue_manager,
	            app_generate_entity=application_generate_entity,
	            prompt_messages=prompt_messages,
	            text=annotation_reply.content,
	            stream=application_generate_entity.stream,
	        )
	        return
		}
	}
	// fill in variable inputs from external data tools if exists
	external_data_tools = app_config.external_data_variables
	if external_data_tools{
	    inputs = runner.fill_in_inputs_from_external_data_tools(
	        tenant_id=app_record.tenant_id,
	        app_id=app_record.id,
	        external_data_tools=external_data_tools,
	        inputs=inputs,
	        query=query,
	    )

	// reorganize all inputs and template to prompt messages
	// Include: prompt template, inputs, query(optional), files(optional)
	//          memory(optional), external data, dataset context(optional)
	prompt_messages, _ = runner.organize_prompt_messages(
	    app_record=app_record,
	    model_config=application_generate_entity.model_conf,
	    prompt_template_entity=app_config.prompt_template,
	    inputs=dict(inputs),
	    files=list(files),
	    query=query or "",
	    memory=memory,
	)

	// check hosting moderation
	hosting_moderation_result = runner.check_hosting_moderation(
	    application_generate_entity=application_generate_entity,
	    queue_manager=queue_manager,
	    prompt_messages=prompt_messages,
	)

	if hosting_moderation_result{
	    return

	agent_entity = app_config.agent
	assert agent_entity is not None

	// init model instance
	model_instance = ModelInstance(
	    provider_model_bundle=application_generate_entity.model_conf.provider_model_bundle,
	    model=application_generate_entity.model_conf.model,
	)
	prompt_message, _ = runner.organize_prompt_messages(
	    app_record=app_record,
	    model_config=application_generate_entity.model_conf,
	    prompt_template_entity=app_config.prompt_template,
	    inputs=dict(inputs),
	    files=list(files),
	    query=query or "",
	    memory=memory,
	)

	// change function call strategy based on LLM model
	llm_model = cast(LargeLanguageModel, model_instance.model_type_instance)
	model_schema = llm_model.get_model_schema(model_instance.model, model_instance.credentials)
	if not model_schema{
	    raise ValueError("Model schema not found")

	if {ModelFeature.MULTI_TOOL_CALL, ModelFeature.TOOL_CALL}.intersection(model_schema.features or []){
	    agent_entity.strategy = AgentEntity.Strategy.FUNCTION_CALLING

	conversation_result = db.session.query(Conversation).where(Conversation.id == conversation.id).first()
	if conversation_result is None{
	    raise ValueError("Conversation not found")
	message_result = db.session.query(Message).where(Message.id == message.id).first()
	if message_result is None{
	    raise ValueError("Message not found")
	db.session.close()

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
