package cyclemanage

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	appresponseentities "mlib.com/gofy/server/entities/app/response"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
)

type applicationGenerateEntityType interface {
	*appgeneratorentities.AdvancedChatAppGenerateEntity | *appgeneratorentities.WorkflowAppGenerateEntity
}
type WorkflowCycleManage[T applicationGenerateEntityType] struct {
	workflowRun               *models.WorkflowRun
	workflowNodeExecutions    map[string]*models.WorkflowNodeExecution
	applicationGenerateEntity T
	workflowSystemVariables   map[workflowenumtypes.SystemVariableKey]any
}

func New[T *appgeneratorentities.AdvancedChatAppGenerateEntity | *appgeneratorentities.WorkflowAppGenerateEntity](
	applicationGenerateEntity T, //  Union[AdvancedChatAppGenerateEntity, WorkflowAppGenerateEntity]
	workflowSystemVariables map[workflowenumtypes.SystemVariableKey]any,
) *WorkflowCycleManage[T] {
	return &WorkflowCycleManage[T]{
		workflowRun:               nil,
		workflowNodeExecutions:    make(map[string]*models.WorkflowNodeExecution),
		applicationGenerateEntity: applicationGenerateEntity,
		workflowSystemVariables:   workflowSystemVariables,
	}
}

func (mgr *WorkflowCycleManage[T]) GetWorkflowRun(workflow_run_id string) *models.WorkflowRun {
	if mgr.workflowRun != nil && mgr.workflowRun.ID == workflow_run_id {
		return mgr.workflowRun
	}
	workflow_run := new(models.WorkflowRun)
	err := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", workflow_run_id).First(workflow_run).Error
	if err != nil {
		mlog.Errorf("find WorkflowRun failed:%v", err)
		panic(exceptions.NewWorkflowRunNotFoundError(workflow_run_id))
	}

	mgr.workflowRun = workflow_run

	return workflow_run
}
func (mgr *WorkflowCycleManage[T]) GetWorkflowNodeExecution(node_execution_id string) *models.WorkflowNodeExecution {
	if _, ok := mgr.workflowNodeExecutions[node_execution_id]; !ok {
		panic(exceptions.NewValueError(fmt.Sprintf("Workflow node execution not found: %s", node_execution_id)))
	}
	cached_workflow_node_execution := mgr.workflowNodeExecutions[node_execution_id]
	return cached_workflow_node_execution
}

func (mgr *WorkflowCycleManage[T]) HandleWorkflowRunStart(
	workflow_id string,
	user_id string,
	created_by_role models.CreatedByRole,
) *models.WorkflowRun {
	wf := new(models.Workflow)
	err := dbengine.Instance().DB.Model(&models.Workflow{}).Where("id = ?", workflow_id).First(wf).Error
	if err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Workflow(%s) not found", workflow_id)))
	}
	var max_sequence int
	err = dbengine.Instance().DB.Model(&models.WorkflowRun{}).Select("sequence_number").Where("tenant_id = ? and app_id = ?", wf.TenantID, wf.AppID).Order("sequence_number DESC").First(&max_sequence).Error
	if err != nil {
		mlog.Warningf("no WorkflowRun exist")
	}
	new_sequence_number := max_sequence + 1

	var inputs map[string]any
	var triggered_from models.WorkflowRunTriggeredFrom
	if realentity, ok := any(mgr.applicationGenerateEntity).(*appgeneratorentities.AdvancedChatAppGenerateEntity); ok {
		inputs = realentity.Inputs
		if realentity.InvokeFrom == appenumtypes.InvokeFrom_DEBUGGER {
			triggered_from = models.WorkflowRunTriggeredFrom_DEBUGGING
		} else {
			triggered_from = models.WorkflowRunTriggeredFrom_APP_RUN
		}
	} else if realentity, ok := any(mgr.applicationGenerateEntity).(*appgeneratorentities.WorkflowAppGenerateEntity); ok {
		inputs = realentity.Inputs
		if realentity.InvokeFrom == appenumtypes.InvokeFrom_DEBUGGER {
			triggered_from = models.WorkflowRunTriggeredFrom_DEBUGGING
		} else {
			triggered_from = models.WorkflowRunTriggeredFrom_APP_RUN
		}
	}

	for key, value := range mgr.workflowSystemVariables {
		if string(key) == "conversation" {
			continue
		}
		inputs[fmt.Sprintf("sys.%v", key)] = value
	}

	// handle special values
	inputs = (&workflow.WorkflowEntry{}).HandleSpecialValues(inputs)

	// init workflow run
	// TODO: This workflow_run_id should always not be None, maybe we can use a more elegant way to handle this
	workflow_run_id := ""
	if _, ok := mgr.workflowSystemVariables[workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID]; ok {
		if _, ok := mgr.workflowSystemVariables[workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID].(string); ok {
			workflow_run_id = mgr.workflowSystemVariables[workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID].(string)
		}
	}
	if workflow_run_id == "" {
		workflow_run_id = uuid.NewV4().String()
	}
	now := time.Now()
	workflow_run := &models.WorkflowRun{}
	workflow_run.ID = workflow_run_id
	workflow_run.TenantID = wf.TenantID
	workflow_run.AppID = wf.AppID
	workflow_run.SequenceNumber = new_sequence_number
	workflow_run.WorkflowID = wf.ID
	workflow_run.Type = string(wf.Type)
	workflow_run.TriggeredFrom = string(triggered_from)
	workflow_run.Version = wf.Version
	workflow_run.Graph = wf.Graph
	bindata, _ := json.Marshal(inputs)
	workflow_run.Inputs = string(bindata)
	workflow_run.Status = models.WorkflowRunStatus_RUNNING
	workflow_run.CreatedByRole = created_by_role
	workflow_run.CreatedBy = user_id
	workflow_run.CreatedAt = &now

	dbengine.Instance().DB.Create(workflow_run)

	return workflow_run
}

func (mgr *WorkflowCycleManage[T]) HandleWorkflowRunSuccess(
	workflow_run_id string,
	start_at time.Time,
	total_tokens int,
	total_steps int,
	outputs map[string]any,
	conversation_id string,
) *models.WorkflowRun {
	/*
		Workflow run success
		:param workflow_run: workflow run
		:param start_at: start time
		:param total_tokens: total tokens
		:param total_steps: total steps
		:param outputs: outputs
		:param conversation_id: conversation id
		:return
	*/
	workflow_run := mgr.GetWorkflowRun(workflow_run_id)

	outputs = (&workflow.WorkflowEntry{}).HandleSpecialValues(outputs)

	workflow_run.Status = models.WorkflowRunStatus_SUCCEEDED
	bindata, _ := json.Marshal(outputs)
	workflow_run.Outputs = string(bindata)
	workflow_run.ElapsedTime = time.Since(start_at).Seconds()
	workflow_run.TotalTokens = total_tokens
	workflow_run.TotalSteps = total_steps
	now := time.Now()
	workflow_run.FinishedAt = &now

	return workflow_run
}
func (mgr *WorkflowCycleManage[T]) HandleWorkflowRunPartialSuccess(
	workflow_run_id string,
	start_at time.Time,
	total_tokens int,
	total_steps int,
	outputs map[string]any,
	exceptions_count int,
	conversation_id string,
) *models.WorkflowRun {
	workflow_run := mgr.GetWorkflowRun(workflow_run_id)

	if outputs != nil {
		outputs = (&workflow.WorkflowEntry{}).HandleSpecialValues(outputs)
	} else {
		outputs = (&workflow.WorkflowEntry{}).HandleSpecialValues(nil)
	}

	now := time.Now()
	workflow_run.Status = models.WorkflowRunStatus_PARTIAL_SUCCESSED
	bindata, _ := json.Marshal(outputs)
	workflow_run.Outputs = string(bindata)
	workflow_run.ElapsedTime = time.Since(start_at).Seconds()
	workflow_run.TotalTokens = total_tokens
	workflow_run.TotalSteps = total_steps
	workflow_run.FinishedAt = &now
	workflow_run.ExceptionsCount = exceptions_count

	return workflow_run
}
func (mgr *WorkflowCycleManage[T]) HandleWorkflowRunFailed(
	workflow_run_id string,
	start_at time.Time,
	total_tokens int,
	total_steps int,
	status models.WorkflowRunStatus,
	errmsg string,
	conversation_id string,
	exceptions_count int,
) *models.WorkflowRun {
	/*
		Workflow run failed
		:param workflow_run: workflow run
		:param start_at: start time
		:param total_tokens: total tokens
		:param total_steps: total steps
		:param status: status
		:param error: error message
		:return
	*/
	workflow_run := mgr.GetWorkflowRun(workflow_run_id)
	workflow_run.Status = string(status)
	workflow_run.Error = errmsg
	workflow_run.ElapsedTime = time.Since(start_at).Seconds()
	workflow_run.TotalTokens = total_tokens
	workflow_run.TotalSteps = total_steps
	now := time.Now()
	workflow_run.FinishedAt = &now
	workflow_run.ExceptionsCount = exceptions_count

	ids := []string{}
	dbengine.Instance().DB.Model(&models.WorkflowNodeExecution{}).Select("node_execution_id").Where("tenant_id = ? and app_id = ? and workflow_id =? and triggered_from = ? and workflow_run_id = ? and status = ?", workflow_run.TenantID, workflow_run.AppID, workflow_run.WorkflowID, models.WorkflowNodeExecutionTriggeredFrom_WORKFLOW_RUN, workflow_run.ID, models.WorkflowNodeExecutionStatus_RUNNING).Find(&ids)

	// Use mgr.GetWorkflowNodeExecution here to make sure the cache is updated
	running_workflow_node_executions := []*models.WorkflowNodeExecution{}
	for _, id := range ids {
		if id != "" {
			tmp := mgr.GetWorkflowNodeExecution(id)
			running_workflow_node_executions = append(running_workflow_node_executions, tmp)
		}
	}

	for _, workflow_node_execution := range running_workflow_node_executions {
		now := time.Now()
		workflow_node_execution.Status = string(models.WorkflowNodeExecutionStatus_FAILED)
		workflow_node_execution.Error = errmsg
		workflow_node_execution.FinishedAt = &now
		workflow_node_execution.ElapsedTime = time.Since(*workflow_node_execution.CreatedAt).Seconds()
	}
	return workflow_run
}
func (mgr *WorkflowCycleManage[T]) HandleNodeExecutionStart(
	workflow_run *models.WorkflowRun, event *appqueueentities.QueueNodeStartedEvent,
) *models.WorkflowNodeExecution {
	workflow_node_execution := event.Handle(workflow_run, nil)

	mgr.workflowNodeExecutions[event.NodeExecutionID] = workflow_node_execution
	return workflow_node_execution
}
func (mgr *WorkflowCycleManage[T]) HandleWorkflowNodeExecutionSuccess(event *appqueueentities.QueueNodeSucceededEvent) *models.WorkflowNodeExecution {
	workflow_node_execution := mgr.GetWorkflowNodeExecution(event.NodeExecutionID)
	return event.Handle(nil, workflow_node_execution)
}
func (mgr *WorkflowCycleManage[T]) HandleWorkflowNodeExecutionFailed(event appqueueentities.QueueNodeEventer) *models.WorkflowNodeExecution {
	/*
		Workflow node execution failed
		:param event: queue node failed event
		:return
	*/
	//: QueueNodeFailedEvent | QueueNodeInIterationFailedEvent | QueueNodeExceptionEvent
	switch ev := event.(type) {
	case *appqueueentities.QueueNodeFailedEvent:
		workflow_node_execution := mgr.GetWorkflowNodeExecution(ev.NodeExecutionID)
		return ev.Handle(nil, workflow_node_execution)
	case *appqueueentities.QueueNodeInIterationFailedEvent:
		workflow_node_execution := mgr.GetWorkflowNodeExecution(ev.NodeExecutionID)
		return ev.Handle(nil, workflow_node_execution)
	case *appqueueentities.QueueNodeExceptionEvent:
		workflow_node_execution := mgr.GetWorkflowNodeExecution(ev.NodeExecutionID)
		return ev.Handle(nil, workflow_node_execution)
	default:
		panic(exceptions.NewValueError("event must be QueueNodeFailedEvent, QueueNodeInIterationFailedEvent or QueueNodeExceptionEvent"))
	}
}
func (mgr *WorkflowCycleManage[T]) HandleWorkflowNodeExecutionRetried(workflow_run *models.WorkflowRun, event *appqueueentities.QueueNodeRetryEvent) *models.WorkflowNodeExecution {
	/*
		Workflow node execution failed
		:param event: queue node failed event
		:return
	*/
	workflow_node_execution := event.Handle(workflow_run, nil)

	mgr.workflowNodeExecutions[event.NodeExecutionID] = workflow_node_execution
	return workflow_node_execution
}

/*------------------------------------------------
             to stream responses
--------------------------------------------------*/

func (mgr *WorkflowCycleManage[T]) WorkflowStartToStreamResponse(task_id string, workflow_run *models.WorkflowRun) *appresponseentities.WorkflowStartStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	return appresponseentities.NewWorkflowStartStreamResponse(task_id, workflow_run)
}
func (mgr *WorkflowCycleManage[T]) WorkflowFinishToStreamResponse(task_id string, workflow_run *models.WorkflowRun) *appresponseentities.WorkflowFinishStreamResponse {
	return appresponseentities.NewWorkflowFinishStreamResponse(task_id, workflow_run)
}
func (mgr *WorkflowCycleManage[T]) WorkflowNodeStartToStreamResponse(event *appqueueentities.QueueNodeStartedEvent, task_id string, workflow_node_execution *models.WorkflowNodeExecution) *appresponseentities.NodeStartStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this

	if slices.Contains([]nodesenumtypes.NodeType{nodesenumtypes.Node_ITERATION, nodesenumtypes.Node_LOOP}, workflow_node_execution.NodeType) {
		return nil
	}
	if workflow_node_execution.WorkflowRunID == "" {
		return nil
	}
	return appresponseentities.NewNodeStartStreamResponse(event, task_id, workflow_node_execution)
}
func (mgr *WorkflowCycleManage[T]) WorkflowNodeFinishToStreamResponse(
	event any,
	task_id string,
	workflow_node_execution *models.WorkflowNodeExecution,
) *appresponseentities.NodeFinishStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this

	if slices.Contains([]nodesenumtypes.NodeType{nodesenumtypes.Node_ITERATION, nodesenumtypes.Node_LOOP}, workflow_node_execution.NodeType) {
		return nil
	}
	if workflow_node_execution.WorkflowRunID == "" {
		return nil
	}
	if workflow_node_execution.FinishedAt == nil {
		return nil
	}
	switch ev := event.(type) {
	case *appqueueentities.QueueNodeSucceededEvent:
		return appresponseentities.NewNodeFinishStreamResponse(ev, task_id, workflow_node_execution)
	case *appqueueentities.QueueNodeFailedEvent:
		return appresponseentities.NewNodeFinishStreamResponse(ev, task_id, workflow_node_execution)
	case *appqueueentities.QueueNodeInIterationFailedEvent:
		return appresponseentities.NewNodeFinishStreamResponse(ev, task_id, workflow_node_execution)
	case *appqueueentities.QueueNodeExceptionEvent:
		return appresponseentities.NewNodeFinishStreamResponse(ev, task_id, workflow_node_execution)
	default:
		panic(exceptions.NewValueError(fmt.Sprintf("unsurported event(%#v)", event)))
	}
}
func (mgr *WorkflowCycleManage[T]) WorkflowNodeRetryToStreamResponse(
	event *appqueueentities.QueueNodeRetryEvent,
	task_id string,
	workflow_node_execution *models.WorkflowNodeExecution,
) *appresponseentities.NodeRetryStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	if slices.Contains([]nodesenumtypes.NodeType{nodesenumtypes.Node_ITERATION, nodesenumtypes.Node_LOOP}, workflow_node_execution.NodeType) {
		return nil
	}
	if workflow_node_execution.WorkflowRunID == "" {
		return nil
	}
	if workflow_node_execution.FinishedAt == nil {
		return nil
	}
	return appresponseentities.NewNodeRetryStreamResponse(event, task_id, workflow_node_execution)
}
func (mgr *WorkflowCycleManage[T]) WorkflowParallelBranchStartToStreamResponse(
	task_id string, workflow_run *models.WorkflowRun, event *appqueueentities.QueueParallelBranchRunStartedEvent,
) *appresponseentities.ParallelBranchStartStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	return appresponseentities.NewParallelBranchStartStreamResponse(event, task_id, workflow_run)
}
func (mgr *WorkflowCycleManage[T]) WorkflowParallelBranchFinishedToStreamResponse(
	task_id string,
	workflow_run *models.WorkflowRun,
	event any,
) *appresponseentities.ParallelBranchFinishedStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	switch ev := event.(type) {
	case *appqueueentities.QueueParallelBranchRunSucceededEvent:
		return appresponseentities.NewParallelBranchFinishedStreamResponse(ev, task_id, workflow_run)
	case *appqueueentities.QueueParallelBranchRunFailedEvent:
		return appresponseentities.NewParallelBranchFinishedStreamResponse(ev, task_id, workflow_run)
	default:
		panic(exceptions.NewValueError(fmt.Sprintf("unsurported event(%#v)", event)))
	}
}
func (mgr *WorkflowCycleManage[T]) WorkflowIterationStartToStreamResponse(
	task_id string, workflow_run *models.WorkflowRun, event *appqueueentities.QueueIterationStartEvent,
) *appresponseentities.IterationNodeStartStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	return appresponseentities.NewIterationNodeStartStreamResponse(event, task_id, workflow_run)
}
func (mgr *WorkflowCycleManage[T]) WorkflowIterationNextToStreamResponse(
	task_id string, workflow_run *models.WorkflowRun, event *appqueueentities.QueueIterationNextEvent,
) *appresponseentities.IterationNodeNextStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	return appresponseentities.NewIterationNodeNextStreamResponse(event, task_id, workflow_run)
}
func (mgr *WorkflowCycleManage[T]) WorkflowIterationCompletedToStreamResponse(
	task_id string, workflow_run *models.WorkflowRun, event *appqueueentities.QueueIterationCompletedEvent,
) *appresponseentities.IterationNodeCompletedStreamResponse {
	// receive session to make sure the workflow_run won't be expired, need a more elegant way to handle this
	return appresponseentities.NewIterationNodeCompletedStreamResponse(event, task_id, workflow_run)
}
