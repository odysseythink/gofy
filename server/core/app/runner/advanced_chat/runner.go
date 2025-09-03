package advancedchat

import (
	"slices"

	wfbasedrunner "mlib.com/gofy/server/core/app/runner/workflow_based"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/variables"
	"mlib.com/gofy/server/core/workflow"
	wfcallbacks "mlib.com/gofy/server/core/workflow/callbacks"
	"mlib.com/gofy/server/core/workflow/graph"
	dbengine "mlib.com/gofy/server/db_engine"
	appgeneratorentities "mlib.com/gofy/server/entities/app/generator"
	appqueueentities "mlib.com/gofy/server/entities/app/queue"
	wfentities "mlib.com/gofy/server/entities/workflow"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	wfenumtypes "mlib.com/gofy/server/enum_types/workflow"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

type AdvancedChatAppRunner struct {
	*wfbasedrunner.WorkflowBasedAppRunner[*appqueueentities.MessageQueueMessage]
	application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity
	conversation                *models.Conversation
	message                     *models.Message
	_dialogue_count             int
}

func New(
	application_generate_entity *appgeneratorentities.AdvancedChatAppGenerateEntity,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.MessageQueueMessage],
	conversation *models.Conversation,
	message *models.Message,
	dialogue_count int,
) *AdvancedChatAppRunner {
	return &AdvancedChatAppRunner{
		WorkflowBasedAppRunner:      wfbasedrunner.New(queue_manager),
		application_generate_entity: application_generate_entity,
		conversation:                conversation,
		message:                     message,
		_dialogue_count:             dialogue_count,
	}
}

func (runner *AdvancedChatAppRunner) _complete_with_stream_output(text string, stopped_by appenumtypes.QueueStopEvent_StopBy) {
	/*
		Direct output
	*/
	runner.PublishEvent(&appqueueentities.QueueTextChunkEvent{Text: text})

	runner.PublishEvent(&appqueueentities.QueueStopEvent{StoppedBy: stopped_by})
}

func (runner *AdvancedChatAppRunner) Run() {
	app_config := runner.application_generate_entity.AppConfig

	app_record := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", app_config.AppID).First(app_record).Error
	if err != nil {
		mlog.Errorf("get App failed:%v", err)
		app_record = nil
	}

	if app_record == nil {
		panic(exceptions.NewValueError("App not found"))
	}
	wf := runner.GetWorkflow(app_record, app_config.WorkflowID)
	if wf == nil {
		panic(exceptions.NewValueError("Workflow not initialized"))
	}
	user_id := ""
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_WEB_APP, appenumtypes.InvokeFrom_SERVICE_API}, runner.application_generate_entity.InvokeFrom) {
		end_user := new(models.EndUser)
		err := dbengine.Instance().DB.Model(&models.EndUser{}).Where("id = ?", runner.application_generate_entity.UserID).First(end_user).Error
		if err != nil {
			mlog.Errorf("get EndUser failed:%v", err)
			end_user = nil
		} else {
			user_id = end_user.SessionID
		}
	} else {
		user_id = runner.application_generate_entity.UserID
	}
	workflow_callbacks := []wfcallbacks.WorkflowCallback{}
	workflow_callbacks = append(workflow_callbacks, &wfcallbacks.WorkflowLoggingCallback{})
	var graph *graph.Graph
	var variable_pool *wfentities.VariablePool
	if runner.application_generate_entity.SingleIterationRun != nil {
		// if only single iteration run is requested
		graph, variable_pool = runner.GetGraphAndVariablePoolOfSingleIteration(
			wf,
			runner.application_generate_entity.SingleIterationRun.NodeID,
			runner.application_generate_entity.SingleIterationRun.Inputs,
		)
	} else {
		inputs := runner.application_generate_entity.Inputs
		query := runner.application_generate_entity.Query
		files := runner.application_generate_entity.Files

		// Init conversation variables
		var db_conversation_variables []*models.ConversationVariable
		err := dbengine.Instance().DB.Model(&models.ConversationVariable{}).Where("app_id = ? and conversation_id = ?", runner.conversation.AppID, runner.conversation.ID).Find(&db_conversation_variables).Error
		if err != nil {
			mlog.Errorf("get ConversationVariable failed:%v", err)
			db_conversation_variables = nil
		}

		if len(db_conversation_variables) == 0 {
			// Create conversation variables if they don't exist.
			for _, variable := range wf.GetConversationVariables() {
				tmp := (&models.ConversationVariable{}).FromVariable(runner.conversation.AppID, runner.conversation.ID, variable)
				if tmp != nil {
					if db_conversation_variables == nil {
						db_conversation_variables = make([]*models.ConversationVariable, 0)
					}
					db_conversation_variables = append(db_conversation_variables, tmp)
					dbengine.Instance().DB.Create(tmp)
				}
			}
		}
		// Convert database entities to variables.
		conversation_variables := []variables.Variabler{}
		for _, item := range db_conversation_variables {
			conversation_variables = append(conversation_variables, item.ToVariable())
		}

		// Create a variable pool.
		system_inputs := map[wfenumtypes.SystemVariableKey]any{
			wfenumtypes.SystemVariableKey_QUERY:           query,
			wfenumtypes.SystemVariableKey_FILES:           files,
			wfenumtypes.SystemVariableKey_CONVERSATION_ID: runner.conversation.ID,
			wfenumtypes.SystemVariableKey_USER_ID:         user_id,
			wfenumtypes.SystemVariableKey_DIALOGUE_COUNT:  runner._dialogue_count,
			wfenumtypes.SystemVariableKey_APP_ID:          app_config.AppID,
			wfenumtypes.SystemVariableKey_WORKFLOW_ID:     app_config.WorkflowID,
			wfenumtypes.SystemVariableKey_WORKFLOW_RUN_ID: runner.application_generate_entity.WorkflowRunID,
		}

		// init variable pool
		variable_pool = wfentities.NewVariablePool(
			system_inputs,
			inputs,
			wf.GetEnvironmentVariables(),
			conversation_variables,
		)

		// init graph
		graph = runner.InitGraph(wf.GraphDict())
	}
	user_from := models.UserFrom_END_USER
	if slices.Contains([]appenumtypes.InvokeFrom{appenumtypes.InvokeFrom_EXPLORE, appenumtypes.InvokeFrom_DEBUGGER}, runner.application_generate_entity.InvokeFrom) {
		user_from = models.UserFrom_ACCOUNT
	}
	// RUN WORKFLOW
	workflow_entry := workflow.NewWorkflowEntry(
		wf.TenantID,
		wf.AppID,
		wf.ID,
		wf.Type,

		wf.GraphDict(),
		graph,
		runner.application_generate_entity.UserID,
		user_from,
		runner.application_generate_entity.InvokeFrom,
		runner.application_generate_entity.CallDepth,
		variable_pool,
	)

	generator := workflow_entry.Run(
		workflow_callbacks,
	)
	for event := range generator {
		if any(event) == nil {
			return
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
					panic(r)
				}
			}()
			runner.HandleEvent(workflow_entry, event)
		}()

	}
	mlog.Debugf("------ workflow entry run return")
}
