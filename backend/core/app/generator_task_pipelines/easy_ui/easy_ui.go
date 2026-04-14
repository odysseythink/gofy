package easyui

import (
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	messagecyclemgr "mlib.com/gofy/server/core/app/cycle_manage/message"
	"mlib.com/gofy/server/core/app/generator_task_pipelines/base"
	"mlib.com/gofy/server/core/exceptions"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	dbengine "mlib.com/gofy/server/db_engine"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appresponseentities "mlib.com/gofy/server/entities/app/response"
	appresponserentities "mlib.com/gofy/server/entities/app/responser"
	apptaskentities "mlib.com/gofy/server/entities/app/task"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/events"
	"mlib.com/gofy/server/models"
)

type EasyUIGenerateTaskPipeline[T1 interface {
	*appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity
}] struct {
	*base.BaseGeneratorTaskPipeline[T1, *appqueueentities.MessageQueueMessage]
	*messagecyclemgr.MessageCycleManage[T1, *apptaskentities.EasyUITaskState]
	ModelConfig      *appconfigentities.ModelConfigWithCredentialsEntity
	AppConfig        *appconfigentities.EasyUIBasedAppConfig
	ConversationID   string
	ConversationMode models.AppMode
	MessageID        string
	MessageCreatedAt int64
}

func New[T1 interface {
	*appgeneratorentities.ChatAppGenerateEntity | *appgeneratorentities.CompletionAppGenerateEntity | *appgeneratorentities.AgentChatAppGenerateEntity
}](
	application_generate_entity T1,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
	conversation *models.Conversation,
	message *models.Message,
	stream bool,
) *EasyUIGenerateTaskPipeline[T1] {

	pl := &EasyUIGenerateTaskPipeline[T1]{
		BaseGeneratorTaskPipeline: base.New(application_generate_entity, queue_manager, stream),
		ConversationID:            conversation.ID,
		ConversationMode:          conversation.Mode,
		MessageID:                 message.ID,
		MessageCreatedAt:          message.CreatedAt.Unix(),
	}
	switch real_entity := any(application_generate_entity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		pl.ModelConfig = real_entity.ModelConf
		pl.AppConfig = real_entity.AppConfig
	case *appgeneratorentities.CompletionAppGenerateEntity:
		pl.ModelConfig = real_entity.ModelConf
		pl.AppConfig = real_entity.AppConfig
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		pl.ModelConfig = real_entity.ModelConf
		pl.AppConfig = real_entity.AppConfig.EasyUIBasedAppConfig
	}
	pl.StartAt = time.Now()
	task_state := &apptaskentities.EasyUITaskState{
		LLMResult: &modelruntimeentities.LLMResult{
			ID:             uuid.NewV4().String(),
			Model:          pl.ModelConfig.Model,
			PromptMessages: make([]modelruntimeentities.PromptMessager, 0),
			Message:        modelruntimeentities.NewAssistantPromptMessage("", "", nil),
			Usage:          modelruntimeentities.NewLLMUsage(),
		},
	}
	pl.MessageCycleManage = messagecyclemgr.New(application_generate_entity, task_state)
	return pl
}
func (pl *EasyUIGenerateTaskPipeline[T1]) Process() any { /* -> Union[
	    ChatbotAppBlockingResponse,
	    CompletionAppBlockingResponse,
	    Generator[Union[ChatbotAppStreamResponse, CompletionAppStreamResponse], None, None],
	]:*/
	var app_mode models.AppMode
	query := ""
	switch real_entity := any(pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		app_mode = real_entity.AppConfig.AppMode
		query = real_entity.Query
	case *appgeneratorentities.CompletionAppGenerateEntity:
		app_mode = real_entity.AppConfig.AppMode
		query = real_entity.Query
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		app_mode = real_entity.AppConfig.AppMode
		query = real_entity.Query
	}
	if app_mode != models.AppMode_COMPLETION {
		// start generate conversation name thread
		pl.GenerateConversationName(pl.ConversationID, query)
	}
	generator := pl._wrapper_process_stream_response()
	if pl.Stream {
		return pl._to_stream_response(generator)
	} else {
		return pl._to_blocking_response(generator)
	}
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _to_blocking_response(
	generator iter.Seq[appresponserentities.StreamResponser],
) any /*-> Union[ChatbotAppBlockingResponse, CompletionAppBlockingResponse]:*/ {
	/*
		Process blocking response.
		:return:
	*/
	for stream_response := range generator {
		if real_stream_response, ok := any(stream_response).(*appresponseentities.ErrorStreamResponse); ok {
			panic(real_stream_response.Err)
		} else if _, ok := any(stream_response).(*appresponseentities.MessageEndStreamResponse); ok {
			extras := map[string]any{"usage": pl.TaskState.LLMResult.Usage}
			if len(pl.TaskState.Metadata) > 0 {
				extras["metadata"] = pl.TaskState.Metadata
			}
			task_id := ""
			switch real_entity := any(pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity).(type) {
			case *appgeneratorentities.ChatAppGenerateEntity:
				task_id = real_entity.TaskID
			case *appgeneratorentities.CompletionAppGenerateEntity:
				task_id = real_entity.TaskID
			case *appgeneratorentities.AgentChatAppGenerateEntity:
				task_id = real_entity.TaskID
			}
			// response: Union[ChatbotAppBlockingResponse, CompletionAppBlockingResponse]
			if pl.ConversationMode == models.AppMode_COMPLETION {
				return appresponseentities.NewCompletionAppBlockingResponse(
					task_id,
					pl.MessageID,
					pl.TaskState.LLMResult.Message.Content,
					pl.ConversationMode,
					pl.MessageCreatedAt,
					pl.TaskState.Metadata,
				)
			} else {
				return appresponseentities.NewChatbotAppBlockingResponse(
					task_id,
					pl.MessageID,
					pl.ConversationID,
					pl.TaskState.LLMResult.Message.Content,
					pl.ConversationMode,
					pl.MessageCreatedAt,
					pl.TaskState.Metadata,
				)
			}
		} else {
			continue
		}
	}
	mlog.Errorf("queue listening stopped unexpectedly.")
	panic(exceptions.NewRuntimeError("queue listening stopped unexpectedly."))
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _to_stream_response(generator iter.Seq[appresponserentities.StreamResponser]) iter.Seq[any] {
	/*-> Generator[Union[ChatbotAppStreamResponse, CompletionAppStreamResponse], None, None]:*/
	return func(yield func(any) bool) {
		for stream_response := range generator {
			if _, ok := any(pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity).(*appgeneratorentities.CompletionAppGenerateEntity); ok {
				if !yield(&appresponseentities.CompletionAppStreamResponse{
					MessageID: pl.MessageID,
					CreatedAt: pl.MessageCreatedAt,
					AppStreamResponse: &appresponseentities.AppStreamResponse{
						StreamResponser: stream_response,
					},
				}) {
					mlog.Errorf("yield failed")
					return
				}
			} else {
				if !yield(&appresponseentities.ChatbotAppStreamResponse{
					ConversationID: pl.ConversationID,
					MessageID:      pl.MessageID,
					CreatedAt:      pl.MessageCreatedAt,
					AppStreamResponse: &appresponseentities.AppStreamResponse{
						StreamResponser: stream_response,
					},
				}) {
					mlog.Errorf("yield failed")
					return
				}
			}
		}
	}
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _wrapper_process_stream_response() iter.Seq[appresponserentities.StreamResponser] {
	return func(yield func(appresponserentities.StreamResponser) bool) {
		for response := range pl._process_stream_response() {
			if !yield(response) {
				mlog.Errorf("yield failed")
				return
			}
		}
	}
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _process_stream_response() iter.Seq[appresponserentities.StreamResponser] {
	/*
	   Process stream response.
	   :return:
	*/
	return func(yield func(appresponserentities.StreamResponser) bool) {
		for queue_message := range pl.BaseGeneratorTaskPipeline.QueueManager.Listen(pl.BaseGeneratorTaskPipeline.QueueManager) {
			mlog.Debugf("------receive message=%#v", queue_message)
			event := queue_message.Eventer
			mlog.Debugf("------receive event=%#v", event)

			if realev, ok := any(event).(*appqueueentities.QueueErrorEvent); ok {
				mlog.Debugf("QueueErrorEvent ")
				exp := pl.HandleError(realev, "")
				yield(pl.ErrorToStreamResponse(exp))
				return
			} else if realev, ok := any(event).(*appqueueentities.QueueStopEvent); ok {
				pl._handle_stop(realev)

				// Save message
				pl._save_message()
				message_end_resp := pl._message_end_to_stream_response()

				if !yield(message_end_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueMessageEndEvent); ok {
				if realev.LLMResult != nil {
					pl.TaskState.LLMResult = realev.LLMResult
				}

				// Save message
				pl._save_message()
				message_end_resp := pl._message_end_to_stream_response()

				if !yield(message_end_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueRetrieverResourcesEvent); ok {
				pl.HandleRetrieverResources(realev)
			} else if realev, ok := any(event).(*appqueueentities.QueueAnnotationReplyEvent); ok {
				annotation := pl.HandleAnnotationReply(realev)
				if annotation != nil {
					pl.TaskState.LLMResult.Message.Content = annotation.Content
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueAgentThoughtEvent); ok {
				agent_thought_response := pl._agent_thought_to_stream_response(realev)
				if agent_thought_response != nil {
					if !yield(agent_thought_response) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueLLMChunkEvent); ok {
				chunk := realev.Chunk
				delta_text := chunk.Delta.Message.Content
				if delta_text == "" {
					continue
				}
				if pl.TaskState.LLMResult.PromptMessages == nil {
					pl.TaskState.LLMResult.PromptMessages = chunk.PromptMessages
				}

				current_content := pl.TaskState.LLMResult.Message.Content
				current_content += delta_text
				pl.TaskState.LLMResult.Message.Content = current_content
				if !yield(pl.MessageToStreamResponse(delta_text, pl.MessageID, nil)) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueAgentMessageEvent); ok {
				chunk := realev.Chunk
				delta_text := chunk.Delta.Message.Content
				if delta_text == "" {
					continue
				}
				if pl.TaskState.LLMResult.PromptMessages == nil {
					pl.TaskState.LLMResult.PromptMessages = chunk.PromptMessages
				}

				current_content := pl.TaskState.LLMResult.Message.Content
				current_content += delta_text
				pl.TaskState.LLMResult.Message.Content = current_content
				if !yield(pl._agent_message_to_stream_response(delta_text, pl.MessageID)) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueMessageReplaceEvent); ok {
				if !yield(pl.MessageReplaceToStreamResponse(realev.Text)) {
					mlog.Errorf("yield failed")
					return
				}
			} else if _, ok := any(event).(*appqueueentities.QueuePingEvent); ok {
				if !yield(pl.PingStreamResponse()) {
					mlog.Errorf("yield failed")
					return
				}
			} else {
				continue
			}
		}
	}
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _save_message() error {
	/*
	   Save message.
	   :return:
	*/
	llm_result := pl.TaskState.LLMResult
	usage := llm_result.Usage

	message := new(models.Message)
	err := dbengine.Instance().DB.Model(&models.Message{}).Where("id = ?", pl.MessageID).First(message).Error
	if err != nil {
		mlog.Errorf("get Message(%s) failed:%v", pl.ConversationID, err)
		message = nil
	}
	if message == nil {
		return exceptions.NewValueError(fmt.Sprintf("message %s not found", pl.MessageID))
	}

	conversation := new(models.Conversation)
	err = dbengine.Instance().DB.Model(&models.Conversation{}).Where("id = ?", pl.ConversationID).First(conversation).Error
	if err != nil {
		mlog.Errorf("get Conversation(%s) failed:%v", pl.ConversationID, err)
		conversation = nil
	}
	if conversation == nil {
		return exceptions.NewValueError(fmt.Sprintf("Conversation %s not found", pl.ConversationID))
	}
	message_dict_list := promptutils.PromptMessagesToPromptForSaving(
		pl.ModelConfig.Mode, pl.TaskState.LLMResult.PromptMessages,
	)
	bindata, _ := json.Marshal(message_dict_list)
	message.MessageJson = string(bindata)
	message.MessageTokens = usage.PromptTokens
	message.MessageUnitPrice = usage.PromptUnitPrice
	message.MessagePriceUnit = usage.PromptPriceUnit
	message.Answer = promptutils.RemoveTemplateVariables(strings.TrimSpace(llm_result.Message.Content), false)
	message.AnswerTokens = usage.CompletionTokens
	message.AnswerUnitPrice = usage.CompletionUnitPrice
	message.AnswerPriceUnit = usage.CompletionPriceUnit
	message.ProviderResponseLatency = time.Since(pl.StartAt).Seconds()
	message.TotalPrice = usage.TotalPrice
	message.Currency = usage.Currency
	message.SetMessageMetadata(pl.TaskState.Metadata)

	return events.Instance.MessageWasCreatedSig.Emit("message_was_created", message, pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity)
}
func (pl *EasyUIGenerateTaskPipeline[T1]) _handle_stop(event *appqueueentities.QueueStopEvent) {
	/*
	   Handle stop.
	   :return:
	*/
	model_config := pl.ModelConfig
	model := model_config.Model

	model_instance := modelmanager.NewModelInstance(model_config.ProviderModelBundle, model_config.Model)

	// calculate num tokens
	prompt_tokens := 0
	if event.StoppedBy != appenumtypes.QueueStopEvent_StopBy_ANNOTATION_REPLY {
		prompt_tokens = model_instance.GetLLMNumTokens(pl.TaskState.LLMResult.PromptMessages, nil)
	}
	completion_tokens := 0
	if event.StoppedBy == appenumtypes.QueueStopEvent_StopBy_USER_MANUAL {
		completion_tokens = model_instance.GetLLMNumTokens([]modelruntimeentities.PromptMessager{pl.TaskState.LLMResult.Message}, nil)
	}
	credentials := model_config.Credentials

	// transform usage
	model_type_instance := model_config.ProviderModelBundle.ModelTypeInstance.(modelruntimeentities.LargeLanguageModeler)
	pl.TaskState.LLMResult.Usage = model_type_instance.CalcResponseUsage(model_type_instance, model, credentials, prompt_tokens, completion_tokens)
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _message_end_to_stream_response() *appresponseentities.MessageEndStreamResponse {
	/*
	   Message end to stream response.
	   :return:
	*/
	bindata, _ := json.Marshal(pl.TaskState.LLMResult.Usage)
	pl.TaskState.Metadata["usage"] = string(bindata)

	extras := map[string]any{}
	if len(pl.TaskState.Metadata) > 0 {
		extras["metadata"] = pl.TaskState.Metadata
	}
	task_id := ""
	switch real_entity := any(pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.CompletionAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		task_id = real_entity.TaskID
	}
	return appresponseentities.NewMessageEndStreamResponse(
		task_id,
		pl.MessageID,
		pl.TaskState.Metadata,
	)
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _agent_message_to_stream_response(answer string, message_id string) *appresponseentities.AgentMessageStreamResponse {
	/*
	   Agent message to stream response.
	   :param answer: answer
	   :param message_id: message id
	   :return:
	*/
	task_id := ""
	switch real_entity := any(pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.CompletionAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		task_id = real_entity.TaskID
	}
	return appresponseentities.NewAgentMessageStreamResponse(task_id, message_id, answer)
}

func (pl *EasyUIGenerateTaskPipeline[T1]) _agent_thought_to_stream_response(event *appqueueentities.QueueAgentThoughtEvent) *appresponseentities.AgentThoughtStreamResponse {
	/*
	   Agent thought to stream response.
	   :param event: agent thought event
	   :return:
	*/
	agent_thought := new(models.MessageAgentThought)
	err := dbengine.Instance().DB.Model(&models.MessageAgentThought{}).Where("id = ?", event.AgentThoughtID).First(agent_thought).Error
	if err != nil {
		mlog.Errorf("get MessageAgentThought(%s) failed:%v", event.AgentThoughtID, err)
		agent_thought = nil
	}
	task_id := ""
	switch real_entity := any(pl.BaseGeneratorTaskPipeline.ApplicationGenerateEntity).(type) {
	case *appgeneratorentities.ChatAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.CompletionAppGenerateEntity:
		task_id = real_entity.TaskID
	case *appgeneratorentities.AgentChatAppGenerateEntity:
		task_id = real_entity.TaskID
	}

	if agent_thought != nil {
		rsp := appresponseentities.NewAgentThoughtStreamResponse(agent_thought.ID, task_id)
		rsp.Position = agent_thought.Position
		rsp.Thought = agent_thought.Thought
		rsp.Observation = agent_thought.Observation
		rsp.Tool = agent_thought.Tool
		rsp.ToolLabels = agent_thought.ToolLabels()
		rsp.ToolInput = agent_thought.ToolInput
		rsp.MessageFiles = agent_thought.Files()
		return rsp
	}
	return nil
}
