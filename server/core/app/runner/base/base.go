package base

import (
	"fmt"
	"iter"

	annotationreply "mlib.com/gofy/server/core/app/features/annotation_reply"
	"mlib.com/gofy/server/core/exceptions"
	extdatatool "mlib.com/gofy/server/core/external_data_tool"
	"mlib.com/gofy/server/core/file"
	"mlib.com/gofy/server/core/memory"
	"mlib.com/gofy/server/core/prompt"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
)

type AppRunner[T interface {
	*appqueueentities.MessageQueueMessage | *appqueueentities.WorkflowQueueMessage
}] struct {
}

func (r *AppRunner[T]) HandleInvokeResult(
	invoke_result any, /*: Union[LLMResult, Generator[Any, None, None]]*/
	queue_manager appqueueentities.AppQueueManager[T],
	stream bool,
	agent bool, /* = false*/
) {
	/*
		Handle invoke result
		:param invoke_result: invoke result
		:param queue_manager: application queue manager
		:param stream: stream
		:param agent: agent
		:return:
	*/
	if v, ok := invoke_result.(*modelruntimeentities.LLMResult); ok && !stream {
		r.handleInvokeResultDirect(v, queue_manager, agent)
	} else if v, ok := invoke_result.(iter.Seq[*modelruntimeentities.LLMResultChunk]); ok && stream {
		r.handleInvokeResultStream(v, queue_manager, agent)
	} else {
		panic(exceptions.NewNotImplementedError(fmt.Sprintf("unsupported invoke result type: %#v", invoke_result)))
	}
}
func (r *AppRunner[T]) handleInvokeResultDirect(
	invoke_result *modelruntimeentities.LLMResult, queue_manager appqueueentities.AppQueueManager[T], agent bool,
) {
	/*
		Handle invoke result direct
		:param invoke_result: invoke result
		:param queue_manager: application queue manager
		:param agent: agent
		:return:
	*/
	queue_manager.Publish(
		&appqueueentities.QueueMessageEndEvent{
			LLMResult: invoke_result,
		},
		appenumtypes.PublishFrom_APPLICATION_MANAGER,
	)
}
func (r *AppRunner[T]) handleInvokeResultStream(
	invoke_result iter.Seq[*modelruntimeentities.LLMResultChunk], queue_manager appqueueentities.AppQueueManager[T], agent bool,
) {
	/*
		Handle invoke result
		:param invoke_result: invoke result
		:param queue_manager: application queue manager
		:param agent: agent
		:return:
	*/
	model := ""
	prompt_messages := []modelruntimeentities.PromptMessager{}
	text := ""
	var usage *modelruntimeentities.LLMUsage
	for result := range invoke_result {
		if !agent {
			queue_manager.Publish(&appqueueentities.QueueLLMChunkEvent{Chunk: result}, appenumtypes.PublishFrom_APPLICATION_MANAGER)
		} else {
			queue_manager.Publish(&appqueueentities.QueueAgentMessageEvent{Chunk: result}, appenumtypes.PublishFrom_APPLICATION_MANAGER)
		}
		text += result.Delta.Message.Content

		if model == "" {
			model = result.Model
		}
		if len(prompt_messages) == 0 {
			prompt_messages = result.PromptMessages
		}
		if result.Delta.Usage != nil {
			usage = result.Delta.Usage
		}
	}
	if usage == nil {
		usage = &modelruntimeentities.LLMUsage{}
	}
	llm_result := &modelruntimeentities.LLMResult{
		Model:          model,
		PromptMessages: prompt_messages,
		// Message:        modelruntimeentities.NewAssistantPromptMessage(text, "", nil),
		Usage: usage,
	}
	llm_result.Message = modelruntimeentities.NewAssistantPromptMessage(text, "", nil)

	queue_manager.Publish(&appqueueentities.QueueMessageEndEvent{LLMResult: llm_result}, appenumtypes.PublishFrom_APPLICATION_MANAGER)
}
func (r *AppRunner[T]) OrganizePromptMessages(
	app_record *models.App,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	prompt_template_entity *appconfigentities.PromptTemplateEntity,
	inputs map[string]string,
	files []*file.File,
	query string,
	context string,
	mem *memory.TokenBufferMemory,
) ([]modelruntimeentities.PromptMessager, []string) {
	/*
	   Organize prompt messages
	   :param context:
	   :param app_record: app record
	   :param model_config: model config entity
	   :param prompt_template_entity: prompt template entity
	   :param inputs: inputs
	   :param files: files
	   :param query: query
	   :param memory: memory
	   :param image_detail_config: the image quality config
	   :return:
	*/
	//  get prompt without memory and context
	// prompt_transform: Union[SimplePromptTransform, AdvancedPromptTransform]
	prompt_transform := &prompt.SimplePromptTransform{}
	return prompt_transform.GetPrompt(
		app_record.Mode,
		prompt_template_entity,
		inputs,
		query,
		files,
		context,
		mem,
		model_config,
	)

}

func (r *AppRunner[T]) QueryAppAnnotationsToReply(
	app_record *models.App, message *models.Message, query string, user_id string, invoke_from appenumtypes.InvokeFrom,
) *models.MessageAnnotation {
	/*
	   Query app annotations to reply
	   :param app_record: app record
	   :param message: message
	   :param query: query
	   :param user_id: user id
	   :param invoke_from: invoke from
	   :return:
	*/
	annotation_reply_feature := &annotationreply.AnnotationReplyFeature{}
	return annotation_reply_feature.Query(app_record, message, query, user_id, invoke_from)
}
func (r *AppRunner[T]) FillInInputsFromExternalDataTools(
	tenant_id string,
	app_id string,
	external_data_tools []*appconfigentities.ExternalDataVariableEntity,
	inputs map[string]any,
	query string,
) map[string]any {
	/*
	   Fill in variable inputs from external data tools if exists.

	   :param tenant_id: workspace id
	   :param app_id: app id
	   :param external_data_tools: external data tools configs
	   :param inputs: the inputs
	   :param query: the query
	   :return: the filled inputs
	*/
	external_data_fetch_feature := &extdatatool.ExternalDataFetch{}
	return external_data_fetch_feature.Fetch(
		tenant_id, app_id, external_data_tools, inputs, query,
	)
}
