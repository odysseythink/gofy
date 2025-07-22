package base

import (
	"fmt"
	"iter"

	"mlib.com/gofy/server/core/exceptions"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
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
