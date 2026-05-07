package workflow

import (
	"iter"
	"time"

	workflowcyclemanage "github.com/odysseythink/gofy/backend/core/app/cycle_manage/workflow"
	basegeneratortaskpipeline "github.com/odysseythink/gofy/backend/core/app/generator_task_pipelines/base"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appgeneratorentities "github.com/odysseythink/gofy/backend/entities/app/generator"
	appqueueentities "github.com/odysseythink/gofy/backend/entities/app/queue"
	appresponseentities "github.com/odysseythink/gofy/backend/entities/app/response"
	appresponserentities "github.com/odysseythink/gofy/backend/entities/app/responser"
	apptaskentities "github.com/odysseythink/gofy/backend/entities/app/task"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// WorkflowAppGenerateTaskPipeline 类的 Go 实现
type WorkflowAppGenerateTaskPipeline struct {
	baseTaskPipeline          *basegeneratortaskpipeline.BaseGeneratorTaskPipeline[*appgeneratorentities.WorkflowAppGenerateEntity, *appqueueentities.WorkflowQueueMessage]
	userID                    string
	createdByRole             models.CreatedByRole
	workflowCycleManager      *workflowcyclemanage.WorkflowCycleManage[*appgeneratorentities.WorkflowAppGenerateEntity]
	applicationGenerateEntity *appgeneratorentities.WorkflowAppGenerateEntity
	workflowID                string
	workflowFeaturesDict      map[string]interface{}
	taskState                 *apptaskentities.WorkflowTaskState
	workflowRunID             string
}

func New[T *models.Account | *models.EndUser](
	application_generate_entity *appgeneratorentities.WorkflowAppGenerateEntity,
	wf *models.Workflow,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.WorkflowQueueMessage],
	user T, /*: Union[Account, EndUser]*/
	stream bool,
) *WorkflowAppGenerateTaskPipeline {
	pl := &WorkflowAppGenerateTaskPipeline{
		baseTaskPipeline: basegeneratortaskpipeline.New(application_generate_entity, queue_manager, stream),
	}
	user_session_id := ""
	switch realuser := any(user).(type) {
	case *models.EndUser:
		pl.userID = realuser.ID
		user_session_id = realuser.SessionID
		pl.createdByRole = models.CreatedByRole_END_USER
	case *models.Account:
		pl.userID = realuser.ID
		user_session_id = realuser.ID
		pl.createdByRole = models.CreatedByRole_ACCOUNT
	}
	pl.workflowCycleManager = workflowcyclemanage.New(
		application_generate_entity,
		map[workflowenumtypes.SystemVariableKey]any{
			workflowenumtypes.SystemVariableKey_FILES:           application_generate_entity.Files,
			workflowenumtypes.SystemVariableKey_USER_ID:         user_session_id,
			workflowenumtypes.SystemVariableKey_APP_ID:          application_generate_entity.AppConfig.AppID,
			workflowenumtypes.SystemVariableKey_WORKFLOW_ID:     wf.ID,
			workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID: application_generate_entity.WorkflowRunID,
		},
	)

	pl.applicationGenerateEntity = application_generate_entity
	pl.workflowRunID = wf.ID
	pl.workflowFeaturesDict = wf.FeaturesDict()
	pl.taskState = &apptaskentities.WorkflowTaskState{}
	pl.workflowRunID = ""
	pl.workflowID = wf.ID
	return pl
}

func (pl *WorkflowAppGenerateTaskPipeline) Process() (*appresponseentities.WorkflowAppBlockingResponse, iter.Seq[*appresponseentities.WorkflowAppStreamResponse]) {
	/*
	   Process generate task pipeline.
	   :return:
	*/
	generator := pl.wrapperProcessStreamResponse()
	if pl.baseTaskPipeline.Stream {
		return nil, pl._to_stream_response(generator)
	} else {
		rsp := pl._to_blocking_response(generator)
		return rsp, nil
	}
}
func (pl *WorkflowAppGenerateTaskPipeline) _to_blocking_response(generator iter.Seq[appresponserentities.StreamResponser]) *appresponseentities.WorkflowAppBlockingResponse {
	/*
		To blocking response.
		:return:
	*/
	for stream_response := range generator {
		switch realstreamresponse := any(stream_response).(type) {
		case *appresponseentities.ErrorStreamResponse:
			panic(realstreamresponse.Err)
		case *appresponseentities.WorkflowFinishStreamResponse:
			return appresponseentities.NewWorkflowAppBlockingResponse(realstreamresponse, pl.applicationGenerateEntity.TaskID)
		default:
			continue
		}
	}
	panic(exceptions.NewValueError("queue listening stopped unexpectedly."))
}

func (pl *WorkflowAppGenerateTaskPipeline) _to_stream_response(generator iter.Seq[appresponserentities.StreamResponser]) iter.Seq[*appresponseentities.WorkflowAppStreamResponse] {
	/*
		To stream response.
		:return:
	*/
	return func(yield func(*appresponseentities.WorkflowAppStreamResponse) bool) {
		workflow_run_id := ""
		for stream_response := range generator {
			switch realstreamresponse := any(stream_response).(type) {
			case *appresponseentities.WorkflowStartStreamResponse:
				workflow_run_id = realstreamresponse.WorkflowRunID
			}
			if !yield(appresponseentities.NewWorkflowAppStreamResponse(workflow_run_id, stream_response)) {
				return
			}
		}
	}
}

func (pl *WorkflowAppGenerateTaskPipeline) textChunkToStreamResponse(
	text string, from_variable_selector []string,
) *appresponseentities.TextChunkStreamResponse {
	/*
		Handle completed event.
		:param text: text
		:return:
	*/
	return appresponseentities.NewTextChunkStreamResponse(pl.applicationGenerateEntity.TaskID, text, from_variable_selector)
}
func (pl *WorkflowAppGenerateTaskPipeline) saveWorkflowAppLog(workflow_run *models.WorkflowRun) {
	/*
	   Save workflow app log.
	   :return:
	*/
	var created_from models.WorkflowAppLogCreatedFrom
	if pl.applicationGenerateEntity.InvokeFrom == appenumtypes.InvokeFrom_SERVICE_API {
		created_from = models.WorkflowAppLogCreatedFrom_SERVICE_API
	} else if pl.applicationGenerateEntity.InvokeFrom == appenumtypes.InvokeFrom_EXPLORE {
		created_from = models.WorkflowAppLogCreatedFrom_INSTALLED_APP
	} else if pl.applicationGenerateEntity.InvokeFrom == appenumtypes.InvokeFrom_WEB_APP {
		created_from = models.WorkflowAppLogCreatedFrom_WEB_APP
	} else {
		// not save log for debugging
		return
	}
	now := time.Now()
	workflow_app_log := &models.WorkflowAppLog{
		ID:            uuid.NewV4().String(),
		TenantID:      workflow_run.TenantID,
		AppID:         workflow_run.AppID,
		WorkflowID:    workflow_run.WorkflowID,
		WorkflowRunID: workflow_run.ID,
		CreatedFrom:   string(created_from),
		CreatedByRole: pl.createdByRole,
		CreatedBy:     pl.userID,
		CreatedAt:     &now,
	}
	dbengine.Instance().DB.Create(workflow_app_log)
}
func (pl *WorkflowAppGenerateTaskPipeline) _process_stream_response() iter.Seq[appresponserentities.StreamResponser] {
	/*
		Process stream response.
		:return:
	*/

	return func(yield func(appresponserentities.StreamResponser) bool) {

		var graph_runtime_state *graphengineentities.GraphRuntimeState

		for queue_message := range pl.baseTaskPipeline.QueueManager.Listen(pl.baseTaskPipeline.QueueManager) {
			mlog.Debugf("------receive message=%#v", queue_message)
			event := queue_message.Eventer
			mlog.Debugf("------receive event=%#v", event)
			if _, ok := any(event).(*appqueueentities.QueuePingEvent); ok {
				if !yield(pl.baseTaskPipeline.PingStreamResponse()) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueErrorEvent); ok {
				mlog.Debugf("QueueErrorEvent ")
				exp := pl.baseTaskPipeline.HandleError(realev, "")
				yield(pl.baseTaskPipeline.ErrorToStreamResponse(exp))
				return
			} else if realev, ok := any(event).(*appqueueentities.QueueWorkflowStartedEvent); ok {
				// override graph runtime state
				graph_runtime_state = realev.GraphRuntimeState
				// init workflow run
				workflow_run := pl.workflowCycleManager.HandleWorkflowRunStart(
					pl.workflowID,
					pl.userID,
					pl.createdByRole,
				)
				pl.workflowRunID = workflow_run.ID
				start_resp := pl.workflowCycleManager.WorkflowStartToStreamResponse(pl.applicationGenerateEntity.TaskID, workflow_run)
				if !yield(start_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueNodeRetryEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)

				workflow_node_execution := pl.workflowCycleManager.HandleWorkflowNodeExecutionRetried(workflow_run, realev)
				response := pl.workflowCycleManager.WorkflowNodeRetryToStreamResponse(
					realev,
					pl.applicationGenerateEntity.TaskID,
					workflow_node_execution,
				)

				if !yield(response) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueNodeStartedEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				workflow_node_execution := pl.workflowCycleManager.HandleNodeExecutionStart(workflow_run, realev)
				node_start_response := pl.workflowCycleManager.WorkflowNodeStartToStreamResponse(
					realev,
					pl.applicationGenerateEntity.TaskID,
					workflow_node_execution,
				)
				if !yield(node_start_response) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueNodeSucceededEvent); ok {
				workflow_node_execution := pl.workflowCycleManager.HandleWorkflowNodeExecutionSuccess(realev)
				node_success_response := pl.workflowCycleManager.WorkflowNodeFinishToStreamResponse(
					realev,
					pl.applicationGenerateEntity.TaskID,
					workflow_node_execution,
				)

				if !yield(node_success_response) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueNodeFailedEvent); ok {
				workflow_node_execution := pl.workflowCycleManager.HandleWorkflowNodeExecutionFailed(realev)
				node_failed_response := pl.workflowCycleManager.WorkflowNodeFinishToStreamResponse(
					realev,
					pl.applicationGenerateEntity.TaskID,
					workflow_node_execution,
				)

				if !yield(node_failed_response) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueNodeInIterationFailedEvent); ok {
				workflow_node_execution := pl.workflowCycleManager.HandleWorkflowNodeExecutionFailed(realev)
				node_failed_response := pl.workflowCycleManager.WorkflowNodeFinishToStreamResponse(
					realev,
					pl.applicationGenerateEntity.TaskID,
					workflow_node_execution,
				)

				if !yield(node_failed_response) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueNodeExceptionEvent); ok {
				workflow_node_execution := pl.workflowCycleManager.HandleWorkflowNodeExecutionFailed(realev)
				node_failed_response := pl.workflowCycleManager.WorkflowNodeFinishToStreamResponse(
					realev,
					pl.applicationGenerateEntity.TaskID,
					workflow_node_execution,
				)

				if !yield(node_failed_response) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueParallelBranchRunStartedEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				parallel_start_resp := pl.workflowCycleManager.WorkflowParallelBranchStartToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
					realev,
				)
				if !yield(parallel_start_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueParallelBranchRunSucceededEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				parallel_finish_resp := pl.workflowCycleManager.WorkflowParallelBranchFinishedToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
					realev,
				)
				if !yield(parallel_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueParallelBranchRunFailedEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				parallel_finish_resp := pl.workflowCycleManager.WorkflowParallelBranchFinishedToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
					realev,
				)
				if !yield(parallel_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueIterationStartEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				iter_start_resp := pl.workflowCycleManager.WorkflowIterationStartToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
					realev,
				)
				if !yield(iter_start_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueIterationNextEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				iter_next_resp := pl.workflowCycleManager.WorkflowIterationNextToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
					realev,
				)
				if !yield(iter_next_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueIterationCompletedEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}

				workflow_run := pl.workflowCycleManager.GetWorkflowRun(pl.workflowRunID)
				iter_finish_resp := pl.workflowCycleManager.WorkflowIterationCompletedToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
					realev,
				)
				if !yield(iter_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueWorkflowSucceededEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				workflow_run := pl.workflowCycleManager.HandleWorkflowRunSuccess(
					pl.workflowRunID,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					realev.Outputs,
					"",
				)
				mlog.Debugf("------workflow_run=%#v", workflow_run)
				// save workflow app log
				pl.saveWorkflowAppLog(workflow_run)
				workflow_finish_resp := pl.workflowCycleManager.WorkflowFinishToStreamResponse(
					pl.applicationGenerateEntity.TaskID,
					workflow_run,
				)
				dbengine.Instance().DB.Save(workflow_run)
				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueWorkflowPartialSuccessEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				workflow_run := pl.workflowCycleManager.HandleWorkflowRunPartialSuccess(
					pl.workflowRunID,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					realev.Outputs,
					realev.ExceptionsCount,
					"",
				)
				// save workflow app log
				pl.saveWorkflowAppLog(workflow_run)
				workflow_finish_resp := pl.workflowCycleManager.WorkflowFinishToStreamResponse(pl.applicationGenerateEntity.TaskID, workflow_run)
				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueWorkflowFailedEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				workflow_run := pl.workflowCycleManager.HandleWorkflowRunFailed(
					pl.workflowRunID,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					models.WorkflowRunStatus_FAILED,
					realev.Error,
					"",
					realev.ExceptionsCount,
				)
				// save workflow app log
				pl.saveWorkflowAppLog(workflow_run)
				workflow_finish_resp := pl.workflowCycleManager.WorkflowFinishToStreamResponse(pl.applicationGenerateEntity.TaskID, workflow_run)
				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield WorkflowFinishStreamResponse failed")
					return
				}

			} else if realev, ok := any(event).(*appqueueentities.QueueStopEvent); ok {
				if pl.workflowRunID == "" {
					mlog.Errorf("workflow run not initialized.")
					panic(exceptions.NewValueError("workflow run not initialized."))
				}
				if graph_runtime_state == nil {
					mlog.Errorf("graph runtime state not initialized.")
					panic(exceptions.NewValueError("graph runtime state not initialized."))
				}

				workflow_run := pl.workflowCycleManager.HandleWorkflowRunFailed(
					pl.workflowRunID,
					graph_runtime_state.StartAt,
					graph_runtime_state.TotalTokens,
					graph_runtime_state.NodeRunSteps,
					models.WorkflowRunStatus_STOPPED,
					realev.GetStopReason(),
					"",
					0,
				)

				// save workflow app log
				pl.saveWorkflowAppLog(workflow_run)
				workflow_finish_resp := pl.workflowCycleManager.WorkflowFinishToStreamResponse(pl.applicationGenerateEntity.TaskID, workflow_run)
				if !yield(workflow_finish_resp) {
					mlog.Errorf("yield failed")
					return
				}
			} else if realev, ok := any(event).(*appqueueentities.QueueTextChunkEvent); ok {
				delta_text := realev.Text
				if delta_text == "" {
					continue
				}

				pl.taskState.Answer += delta_text
				if !yield(pl.textChunkToStreamResponse(delta_text, realev.FromVariableSelector)) {
					mlog.Errorf("yield TextChunkStreamResponse failed")
					return
				}
			} else {
				continue
			}
		}
	}
}
func (pl *WorkflowAppGenerateTaskPipeline) wrapperProcessStreamResponse() iter.Seq[appresponserentities.StreamResponser] {
	return func(yield func(appresponserentities.StreamResponser) bool) {
		// tts_publisher = None
		// task_id := pl.applicationGenerateEntity.task_id
		// tenant_id := pl.applicationGenerateEntity.app_config.tenant_id
		// features_dict := pl.workflowFeaturesDict

		for response := range pl._process_stream_response() {
			if !yield(response) {
				return
			}
		}
	}
}
