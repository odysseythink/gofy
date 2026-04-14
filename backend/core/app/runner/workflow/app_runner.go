package workflow

import (
	"fmt"
	"slices"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	wfbasedrunner "mlib.com/gofy/server/core/app/runner/workflow_based"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow"
	"mlib.com/gofy/server/core/workflow/callbacks"
	"mlib.com/gofy/server/core/workflow/graph"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	workflowenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
)

type WorkflowAppRunner struct {
	*wfbasedrunner.WorkflowBasedAppRunner[*appqueueentities.WorkflowQueueMessage]
	WorkflowAppGenerateEntity *appgeneratorentities.WorkflowAppGenerateEntity
}

func NewWorkflowAppRunner(
	application_generate_entity *appgeneratorentities.WorkflowAppGenerateEntity,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.WorkflowQueueMessage],
) *WorkflowAppRunner {
	/*
		:param application_generate_entity: application generate entity
		:param queue_manager: application queue manager
		:param workflow_thread_pool_id: workflow thread pool id
	*/
	return &WorkflowAppRunner{
		WorkflowBasedAppRunner:    wfbasedrunner.New(queue_manager),
		WorkflowAppGenerateEntity: application_generate_entity,
	}
}

func (r *WorkflowAppRunner) Run() {
	/*
	   Run application
	   :param application_generate_entity: application generate entity
	   :param queue_manager: application queue manager
	   :return:
	*/
	app_config := r.WorkflowAppGenerateEntity.AppConfig
	//  app_config = cast(WorkflowAppConfig, app_config)
	user_id := ""
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_WEB_APP, appenumtypes.InvokeFrom_SERVICE_API}, r.WorkflowAppGenerateEntity.InvokeFrom) {
		end_user := new(models.EndUser)
		err := dbengine.Instance().DB.Model(&models.EndUser{}).Where("id = ?", r.WorkflowAppGenerateEntity.UserID).First(end_user).Error
		if err != nil {
			mlog.Warningf("get end user from database failed:%v", err)
		} else {
			user_id = end_user.SessionID
		}
	} else {
		user_id = r.WorkflowAppGenerateEntity.UserID
	}
	fmt.Println("----------user_id=", user_id)
	app_record := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", app_config.AppID).First(app_record).Error
	if err != nil {
		mlog.Warningf("get app from database failed:%v", err)
		panic(exceptions.NewValueError("App not found"))
	}

	wf := r.GetWorkflow(app_record, app_config.WorkflowID)
	if wf == nil {
		panic(exceptions.NewValueError("Workflow not initialized"))
	}

	workflow_callbacks := []callbacks.WorkflowCallback{}
	if confy.Get[bool]("debug") {
		workflow_callbacks = append(workflow_callbacks, &callbacks.WorkflowLoggingCallback{})
	}
	var gh *graph.Graph
	var variable_pool *workflowentities.VariablePool
	// if only single iteration run is requested
	if r.WorkflowAppGenerateEntity.SingleIterationRun != nil {
		// if only single iteration run is requested
		gh, variable_pool = r.GetGraphAndVariablePoolOfSingleIteration(
			wf,
			r.WorkflowAppGenerateEntity.SingleIterationRun.NodeID,
			r.WorkflowAppGenerateEntity.SingleIterationRun.Inputs,
		)
	} else {
		inputs := r.WorkflowAppGenerateEntity.Inputs
		files := r.WorkflowAppGenerateEntity.Files
		// Create a variable pool.
		system_inputs := map[workflowenumtypes.SystemVariableKey]any{
			workflowenumtypes.SystemVariableKey_FILES:           files,
			workflowenumtypes.SystemVariableKey_USER_ID:         user_id,
			workflowenumtypes.SystemVariableKey_APP_ID:          app_config.AppID,
			workflowenumtypes.SystemVariableKey_WORKFLOW_ID:     app_config.WorkflowID,
			workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID: r.WorkflowAppGenerateEntity.WorkflowRunID,
		}
		mlog.Debugf("------------inputs=%#v", inputs)
		variable_pool = workflowentities.NewVariablePool(
			system_inputs,
			inputs,
			wf.GetEnvironmentVariables(),
			nil,
		)
		// init graph
		gh = r.InitGraph(wf.GraphDict())
	}

	// RUN WORKFLOW
	var user_from models.UserFrom
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_EXPLORE, appenumtypes.InvokeFrom_DEBUGGER}, r.WorkflowAppGenerateEntity.InvokeFrom) {
		user_from = models.UserFrom_ACCOUNT
	} else {
		user_from = models.UserFrom_END_USER
	}
	workflow_entry := workflow.NewWorkflowEntry(
		wf.TenantID,
		wf.AppID,
		wf.ID,
		wf.Type,
		wf.GraphDict(),
		gh,
		r.WorkflowAppGenerateEntity.UserID,
		user_from,
		r.WorkflowAppGenerateEntity.InvokeFrom,
		r.WorkflowAppGenerateEntity.CallDepth,
		variable_pool,
	)

	generator := workflow_entry.Run(workflow_callbacks)
	for event := range generator {

		r.HandleEvent(workflow_entry, event)
	}
	mlog.Debugf("------WorkflowAppRunner return")
}
