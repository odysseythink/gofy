package advancedchat

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"time"

	msgcyclemgr "github.com/odysseythink/gofy/backend/core/app/cycle_manage/message"
	wfcyclemgr "github.com/odysseythink/gofy/backend/core/app/cycle_manage/workflow"
	"github.com/odysseythink/gofy/backend/core/app/generator_task_pipelines/base"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/file"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appgeneratorentities "github.com/odysseythink/gofy/backend/entities/app/generator"
	appqueueentities "github.com/odysseythink/gofy/backend/entities/app/queue"
	appresponseentities "github.com/odysseythink/gofy/backend/entities/app/response"
	appresponserentities "github.com/odysseythink/gofy/backend/entities/app/responser"
	apptaskentities "github.com/odysseythink/gofy/backend/entities/app/task"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	wfenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/events"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type AdvancedChatAppGenerateTaskPipeline struct {
	_base_task_pipeline          *base.BaseGeneratorTaskPipeline[*appgeneratorentities.AdvancedChatAppGenerateEntity, *appqueueentities.MessageQueueMessage]
	_user_id                     string
	_created_by_role             models.CreatedByRole
	_workflow_cycle_manager      *wfcyclemgr.WorkflowCycleManage[*appgeneratorentities.AdvancedChatAppGenerateEntity]
	_task_state                  *apptaskentities.WorkflowTaskState
	_message_cycle_manager       *msgcyclemgr.MessageCycleManage[*appgeneratorentities.AdvancedChatAppGenerateEntity, *apptaskentities.WorkflowTaskState]
	_application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity
	_workflow_id                 string
	_workflow_features_dict      map[string]any
	_conversation_id             string
	_conversation_mode           models.AppMode
	_message_id                  string
	_message_created_at          int64
	_recorded_files              []map[string]any
	_workflow_run_id             string
}

func New[T *models.Account | *models.EndUser](
	application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity,
	wf *models.Workflow,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
	conversation *models.Conversation,
	message *models.Message,
	user T,
	stream bool,
	dialogue_count int,
) *AdvancedChatAppGenerateTaskPipeline {
	pl := &AdvancedChatAppGenerateTaskPipeline{
		_base_task_pipeline: base.New(application_generate_entity, queue_manager, stream),
	}
	user_session_id := ""
	if real_user, ok := any(user).(*models.EndUser); ok {
		pl._user_id = real_user.ID
		user_session_id = real_user.SessionID
		pl._created_by_role = models.CreatedByRole_END_USER
	} else if real_user, ok := any(user).(*models.Account); ok {
		pl._user_id = real_user.ID
		user_session_id = real_user.ID
		pl._created_by_role = models.CreatedByRole_ACCOUNT
	}

	pl._workflow_cycle_manager = wfcyclemgr.New(
		application_generate_entity,
		map[wfenumtypes.SystemVariableKey]any{
			wfenumtypes.SystemVariableKey_QUERY:           message.Query,
			wfenumtypes.SystemVariableKey_FILES:           application_generate_entity.Files,
			wfenumtypes.SystemVariableKey_CONVERSATION_ID: conversation.ID,
			wfenumtypes.SystemVariableKey_USER_ID:         user_session_id,
			wfenumtypes.SystemVariableKey_DIALOGUE_COUNT:  dialogue_count,
			wfenumtypes.SystemVariableKey_APP_ID:          application_generate_entity.AppConfig.AppID,
			wfenumtypes.SystemVariableKey_WORKFLOW_ID:     wf.ID,
			wfenumtypes.SystemVariableKey_WORKFLOW_RUN_ID: application_generate_entity.WorkflowRunID,
		},
	)

	pl._task_state = apptaskentities.NewWorkflowTaskState()
	pl._message_cycle_manager = msgcyclemgr.New(application_generate_entity, pl._task_state)

	pl._application_generate_entity = application_generate_entity
	pl._workflow_id = wf.ID
	pl._workflow_features_dict = wf.FeaturesDict()
	pl._conversation_id = conversation.ID
	pl._conversation_mode = conversation.Mode
	pl._message_id = message.ID
	pl._message_created_at = message.CreatedAt.Unix()
	return pl
}
func (pl *AdvancedChatAppGenerateTaskPipeline) _get_message() *models.Message {
	message := new(models.Message)
	err := dbengine.Instance().DB.Model(&models.Message{}).Where("id =?", pl._message_id).First(message).Error
	if err != nil {
		mlog.Errorf("get message failed:%v", err)
		message = nil
	}
	if message == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Message not found: %s", pl._message_id)))
	}
	return message
}

func (pl *AdvancedChatAppGenerateTaskPipeline) _message_end_to_stream_response() *appresponseentities.MessageEndStreamResponse {
	/*
	   Message end to stream response.
	   :return:
	*/
	extras := map[string]any{}
	if len(pl._task_state.Metadata) > 0 {
		extras["metadata"] = maps.Clone(pl._task_state.Metadata)

		delete(extras["metadata"].(map[string]any), "annotation_reply")
	} else {
		extras["metadata"] = map[string]any{}
	}
	return appresponseentities.NewMessageEndStreamResponse(
		pl._message_id,
		pl._application_generate_entity.TaskID,
		extras["metadata"].(map[string]any),
	)
}

func (pl *AdvancedChatAppGenerateTaskPipeline) _save_message(graph_runtime_state *graphengineentities.GraphRuntimeState) {
	message := pl._get_message()
	mlog.Debugf("------message=%#v", message)
	mlog.Debugf("------pl._task_state=%#v", pl._task_state)
	message.Answer = pl._task_state.Answer
	message.ProviderResponseLatency = time.Since(pl._base_task_pipeline.StartAt).Seconds()
	message.SetMessageMetadata(pl._task_state.Metadata)

	if graph_runtime_state != nil && graph_runtime_state.LLMUsage != nil {
		usage := graph_runtime_state.LLMUsage
		message.MessageTokens = usage.PromptTokens
		message.MessageUnitPrice = usage.PromptUnitPrice
		message.MessagePriceUnit = usage.PromptPriceUnit
		message.AnswerTokens = usage.CompletionTokens
		message.AnswerUnitPrice = usage.CompletionUnitPrice
		message.AnswerPriceUnit = usage.CompletionPriceUnit
		message.TotalPrice = usage.TotalPrice
		message.Currency = usage.Currency
		pl._task_state.Metadata["usage"] = usage
	} else {
		pl._task_state.Metadata["usage"] = modelruntimeentities.NewLLMUsage()
	}
	events.Instance.MessageWasCreatedSig.Emit(message, pl._application_generate_entity)
}

func (pl *AdvancedChatAppGenerateTaskPipeline) _process_stream_response() iter.Seq[appresponserentities.StreamResponser] {
	/*
	   Process stream response.
	   :return:
	*/
	return func(yield func(appresponserentities.StreamResponser) bool) {
		// init fake graph runtime state
		var graph_runtime_state *graphengineentities.GraphRuntimeState

		for queue_message := range pl._base_task_pipeline.QueueManager.Listen(pl._base_task_pipeline.QueueManager) {
			event := queue_message.Eventer

			if _, ok := any(event).(*appqueueentities.QueuePingEvent); ok {
				if !yield(pl._base_task_pipeline.PingStreamResponse()) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueErrorEvent); ok {
				err := pl._base_task_pipeline.HandleError(real_event, pl._message_id)
				if !yield(pl._base_task_pipeline.ErrorToStreamResponse(err)) {
					mlog.Errorf("yield failed")
					return
				}
				break
			} else if real_event, ok := any(event).(*appqueueentities.QueueWorkflowStartedEvent); ok {
				// override graph runtime state
				graph_runtime_state = real_event.GraphRuntimeState

				// init workflow run
				workflow_run := pl._workflow_cycle_manager.HandleWorkflowRunStart(
					pl._workflow_id,
					pl._user_id,
					pl._created_by_role,
				)

				pl._workflow_run_id = workflow_run.ID
				message := pl._get_message()
				if message == nil {
					mlog.Errorf("Message(%s) not found", pl._message_id)
					panic(exceptions.NewValueError(fmt.Sprintf("Message not found: %s", pl._message_id)))
				}

				message.WorkflowRunID = workflow_run.ID
				workflow_start_resp := pl._workflow_cycle_manager.WorkflowStartToStreamResponse(
					pl._application_generate_entity.TaskID, workflow_run,
				)

				if !yield(workflow_start_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueNodeRetryEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)

				workflow_node_execution := pl._workflow_cycle_manager.HandleWorkflowNodeExecutionRetried(
					workflow_run, real_event,
				)

				node_retry_resp := pl._workflow_cycle_manager.WorkflowNodeRetryToStreamResponse(
					real_event,
					pl._application_generate_entity.TaskID,
					workflow_node_execution,
				)

				if node_retry_resp != nil {
					if !yield(node_retry_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueNodeStartedEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)

				workflow_node_execution := pl._workflow_cycle_manager.HandleNodeExecutionStart(
					workflow_run, real_event,
				)

				node_start_resp := pl._workflow_cycle_manager.WorkflowNodeStartToStreamResponse(
					real_event,
					pl._application_generate_entity.TaskID,
					workflow_node_execution,
				)

				if node_start_resp != nil {
					if !yield(node_start_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueNodeSucceededEvent); ok {
				// Record files if it's an answer node or end node
				if slices.Contains([]nodesenumtypes.NodeType{nodesenumtypes.Node_ANSWER, nodesenumtypes.Node_END}, real_event.NodeType) {
					pl._recorded_files = append(pl._recorded_files, file.FetchFilesFromNodeOutputs(real_event.Outputs)...)
				}
				workflow_node_execution := pl._workflow_cycle_manager.HandleWorkflowNodeExecutionSuccess(real_event)
				node_finish_resp := pl._workflow_cycle_manager.WorkflowNodeFinishToStreamResponse(
					real_event,
					pl._application_generate_entity.TaskID,
					workflow_node_execution,
				)

				if node_finish_resp != nil {
					if !yield(node_finish_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueNodeFailedEvent); ok {
				workflow_node_execution := pl._workflow_cycle_manager.HandleWorkflowNodeExecutionFailed(real_event)
				node_finish_resp := pl._workflow_cycle_manager.WorkflowNodeFinishToStreamResponse(
					real_event,
					pl._application_generate_entity.TaskID,
					workflow_node_execution,
				)

				if node_finish_resp != nil {
					if !yield(node_finish_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueNodeInIterationFailedEvent); ok {
				workflow_node_execution := pl._workflow_cycle_manager.HandleWorkflowNodeExecutionFailed(real_event)

				node_finish_resp := pl._workflow_cycle_manager.WorkflowNodeFinishToStreamResponse(
					real_event,
					pl._application_generate_entity.TaskID,
					workflow_node_execution,
				)

				if node_finish_resp != nil {
					if !yield(node_finish_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueNodeExceptionEvent); ok {
				workflow_node_execution := pl._workflow_cycle_manager.HandleWorkflowNodeExecutionFailed(real_event)

				node_finish_resp := pl._workflow_cycle_manager.WorkflowNodeFinishToStreamResponse(
					real_event,
					pl._application_generate_entity.TaskID,
					workflow_node_execution,
				)

				if node_finish_resp != nil {
					if !yield(node_finish_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueParallelBranchRunStartedEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)

				parallel_start_resp := pl._workflow_cycle_manager.WorkflowParallelBranchStartToStreamResponse(
					pl._application_generate_entity.TaskID,
					workflow_run,
					real_event,
				)

				if !yield(parallel_start_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueParallelBranchRunSucceededEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)

				parallel_finish_resp := pl._workflow_cycle_manager.WorkflowParallelBranchFinishedToStreamResponse(
					pl._application_generate_entity.TaskID,
					workflow_run,
					real_event,
				)

				if !yield(parallel_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueParallelBranchRunFailedEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)

				parallel_finish_resp := pl._workflow_cycle_manager.WorkflowParallelBranchFinishedToStreamResponse(
					pl._application_generate_entity.TaskID,
					workflow_run,
					real_event,
				)

				if !yield(parallel_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueIterationStartEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)
				iter_start_resp := pl._workflow_cycle_manager.WorkflowIterationStartToStreamResponse(
					pl._application_generate_entity.TaskID,
					workflow_run,
					real_event,
				)

				if !yield(iter_start_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueIterationNextEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)
				iter_next_resp := pl._workflow_cycle_manager.WorkflowIterationNextToStreamResponse(
					pl._application_generate_entity.TaskID,
					workflow_run,
					real_event,
				)

				if !yield(iter_next_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueIterationCompletedEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.GetWorkflowRun(pl._workflow_run_id)

				iter_finish_resp := pl._workflow_cycle_manager.WorkflowIterationCompletedToStreamResponse(
					pl._application_generate_entity.TaskID,
					workflow_run,
					real_event,
				)

				if !yield(iter_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueWorkflowSucceededEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}
				workflow_run := pl._workflow_cycle_manager.HandleWorkflowRunSuccess(
					pl._workflow_run_id,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					real_event.Outputs,
					pl._conversation_id,
				)

				workflow_finish_resp := pl._workflow_cycle_manager.WorkflowFinishToStreamResponse(
					pl._application_generate_entity.TaskID, workflow_run,
				)

				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
				pl._base_task_pipeline.QueueManager.Publish(
					&appqueueentities.QueueAdvancedChatMessageEndEvent{}, appenumtypes.PublishFrom_TASK_PIPELINE,
				)
			} else if real_event, ok := any(event).(*appqueueentities.QueueWorkflowPartialSuccessEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.HandleWorkflowRunPartialSuccess(
					pl._workflow_run_id,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					real_event.Outputs,
					real_event.ExceptionsCount,
					"",
				)
				workflow_finish_resp := pl._workflow_cycle_manager.WorkflowFinishToStreamResponse(
					pl._application_generate_entity.TaskID, workflow_run,
				)

				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
				pl._base_task_pipeline.QueueManager.Publish(
					&appqueueentities.QueueAdvancedChatMessageEndEvent{}, appenumtypes.PublishFrom_TASK_PIPELINE,
				)
			} else if real_event, ok := any(event).(*appqueueentities.QueueWorkflowFailedEvent); ok {
				if pl._workflow_run_id == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				workflow_run := pl._workflow_cycle_manager.HandleWorkflowRunFailed(
					pl._workflow_run_id,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					models.WorkflowRunStatus_FAILED,
					real_event.Error,
					pl._conversation_id,
					real_event.ExceptionsCount,
				)
				workflow_finish_resp := pl._workflow_cycle_manager.WorkflowFinishToStreamResponse(
					pl._application_generate_entity.TaskID, workflow_run,
				)

				err_event := &appqueueentities.QueueErrorEvent{
					Err: exceptions.NewValueError(fmt.Sprintf("Run failed: %s", workflow_run.Error)),
				}
				err := pl._base_task_pipeline.HandleError(
					err_event, pl._message_id,
				)

				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
				if !yield(pl._base_task_pipeline.ErrorToStreamResponse(err)) {
					mlog.Errorf("yield failed")
					return
				}
				break
			} else if real_event, ok := any(event).(*appqueueentities.QueueStopEvent); ok {
				if pl._workflow_run_id != "" && graph_runtime_state != nil {

					workflow_run := pl._workflow_cycle_manager.HandleWorkflowRunFailed(
						pl._workflow_run_id,
						graph_runtime_state.StartAt,
						graph_runtime_state.TotalTokens,
						graph_runtime_state.NodeRunSteps,
						models.WorkflowRunStatus_STOPPED,
						real_event.GetStopReason(),
						pl._conversation_id,
						0)
					workflow_finish_resp := pl._workflow_cycle_manager.WorkflowFinishToStreamResponse(
						pl._application_generate_entity.TaskID,
						workflow_run,
					)
					// Save message
					pl._save_message(graph_runtime_state)

					if !yield(workflow_finish_resp) {
						mlog.Errorf("yield failed")
						return
					}
				}
				if !yield(pl._message_end_to_stream_response()) {
					mlog.Errorf("yield failed")
					return
				}
				break
			} else if real_event, ok := any(event).(*appqueueentities.QueueRetrieverResourcesEvent); ok {
				pl._message_cycle_manager.HandleRetrieverResources(real_event)

				message := pl._get_message()
				message.SetMessageMetadata(pl._task_state.Metadata)
				dbengine.Instance().DB.Save(message)
			} else if real_event, ok := any(event).(*appqueueentities.QueueAnnotationReplyEvent); ok {
				pl._message_cycle_manager.HandleAnnotationReply(real_event)

				message := pl._get_message()
				message.SetMessageMetadata(pl._task_state.Metadata)
				dbengine.Instance().DB.Save(message)
			} else if real_event, ok := any(event).(*appqueueentities.QueueTextChunkEvent); ok {
				delta_text := real_event.Text
				if delta_text == "" {
					continue
				}

				pl._task_state.Answer += delta_text
				if !yield(pl._message_cycle_manager.MessageToStreamResponse(
					delta_text, pl._message_id, real_event.FromVariableSelector,
				)) {
					mlog.Errorf("yield failed")
					return
				}
			} else if real_event, ok := any(event).(*appqueueentities.QueueMessageReplaceEvent); ok {
				// published by moderation
				if !yield(pl._message_cycle_manager.MessageReplaceToStreamResponse(real_event.Text)) {
					mlog.Errorf("yield failed")
					return
				}
			} else if _, ok := any(event).(*appqueueentities.QueueAdvancedChatMessageEndEvent); ok {
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				pl._save_message(graph_runtime_state)

				if !yield(pl._message_end_to_stream_response()) {
					mlog.Errorf("yield failed")
					return
				}
			} else {
				continue
			}
		}
	}
}

func (pl *AdvancedChatAppGenerateTaskPipeline) Process() (*appresponseentities.ChatbotAppBlockingResponse, iter.Seq[*appresponseentities.ChatbotAppStreamResponse]) {
	/*
	   Process generate task pipeline.
	   :return:
	*/
	// start generate conversation name thread
	pl._message_cycle_manager.GenerateConversationName(
		pl._conversation_id, pl._application_generate_entity.Query,
	)

	generator := pl._wrapper_process_stream_response()

	if pl._base_task_pipeline.Stream {
		return nil, pl._to_stream_response(generator)
	} else {
		return pl._to_blocking_response(generator), nil
	}
}
func (pl *AdvancedChatAppGenerateTaskPipeline) _to_blocking_response(
	generator iter.Seq[appresponserentities.StreamResponser],
) *appresponseentities.ChatbotAppBlockingResponse {
	/*
	   Process blocking response.
	   :return:
	*/
	for stream_response := range generator {
		if real_stream_response, ok := any(stream_response).(*appresponseentities.ErrorStreamResponse); ok {
			panic(real_stream_response.Err)
		} else if real_stream_response, ok := any(stream_response).(*appresponseentities.MessageEndStreamResponse); ok {

			return appresponseentities.NewChatbotAppBlockingResponse(
				real_stream_response.TaskID(),
				pl._message_id,
				pl._conversation_id,
				pl._task_state.Answer,
				pl._conversation_mode,
				pl._message_created_at,
				real_stream_response.Metadata,
			)
		} else {
			continue
		}
	}
	panic(exceptions.NewValueError("queue listening stopped unexpectedly."))
}

func (pl *AdvancedChatAppGenerateTaskPipeline) _to_stream_response(
	generator iter.Seq[appresponserentities.StreamResponser],
) iter.Seq[*appresponseentities.ChatbotAppStreamResponse] {
	/*
	   To stream response.
	   :return:
	*/
	return func(yield func(*appresponseentities.ChatbotAppStreamResponse) bool) {
		for stream_response := range generator {
			if !yield(&appresponseentities.ChatbotAppStreamResponse{
				AppStreamResponse: &appresponseentities.AppStreamResponse{
					StreamResponser: stream_response,
				},
				ConversationID: pl._conversation_id,
				MessageID:      pl._message_id,
				CreatedAt:      pl._message_created_at,
			}) {
				mlog.Errorf("yield failed")
				return
			}
		}
	}
}

func (pl *AdvancedChatAppGenerateTaskPipeline) _wrapper_process_stream_response() iter.Seq[appresponserentities.StreamResponser] {
	return func(yield func(appresponserentities.StreamResponser) bool) {
		for response := range pl._process_stream_response() {
			yield(response)
		}
	}
}
