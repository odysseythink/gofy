package workflow

import (
	"iter"

	fileupload "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/file_upload"
	wfappcfgmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/workflow"
	wfappgeneratorreponseconvertor "github.com/odysseythink/gofy/backend/core/app/generator_response_convertes/workflow"
	wfappgeneratortaskpipeline "github.com/odysseythink/gofy/backend/core/app/generator_task_pipelines/workflow"
	baseappgenerator "github.com/odysseythink/gofy/backend/core/app/generatores/base"
	wfappqueuemanager "github.com/odysseythink/gofy/backend/core/app/queue_manager/workflow"
	workflowapprunner "github.com/odysseythink/gofy/backend/core/app/runner/workflow"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	modelruntimeexceptions "github.com/odysseythink/gofy/backend/core/exceptions/model_runtime"
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	appgeneratorentities "github.com/odysseythink/gofy/backend/entities/app/generator"
	appqueueentities "github.com/odysseythink/gofy/backend/entities/app/queue"
	appresponseentities "github.com/odysseythink/gofy/backend/entities/app/response"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	filefactory "github.com/odysseythink/gofy/backend/factories/file_factory"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

type WorkflowAppGenerator[T1 interface {
	*models.Account | *models.EndUser
}] struct {
	*baseappgenerator.BaseAppGenerator
}

func New[T1 interface {
	*models.Account | *models.EndUser
}]() *WorkflowAppGenerator[T1] {
	return &WorkflowAppGenerator[T1]{
		BaseAppGenerator: &baseappgenerator.BaseAppGenerator{},
	}
}

func (g *WorkflowAppGenerator[T1]) _generate_worker(
	// flask_app: Flask,
	application_generate_entity *appgeneratorentities.WorkflowAppGenerateEntity,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.WorkflowQueueMessage],
	// context: contextvars.Context,
	// workflow_thread_pool_id: Optional[str] = None,
) {
	/*
		Generate worker in a new thread.
		:param flask_app: Flask app
		:param application_generate_entity: application generate entity
		:param queue_manager: queue manager
		:param workflow_thread_pool_id: workflow thread pool id
		:return:
	*/
	defer func() {
		if r := recover(); r != nil {
			if real_exp, ok := r.(*exceptions.GenerateTaskStoppedError); ok {
				mlog.Error(real_exp)
				return
			} else if _, ok := r.(*modelruntimeexceptions.InvokeAuthorizationError); ok {
				queue_manager.PublishError(queue_manager, modelruntimeexceptions.NewInvokeAuthorizationError("Incorrect API key provided"), appenumtypes.PublishFrom_APPLICATION_MANAGER)
			} else if real_exp, ok := r.(*exceptions.ValidationError); ok {
				mlog.Errorf("Validation Error when generating")
				queue_manager.PublishError(queue_manager, real_exp, appenumtypes.PublishFrom_APPLICATION_MANAGER)
			} else if real_exp, ok := r.(*exceptions.ValueError); ok {
				mlog.Errorf("Error when generating")
				queue_manager.PublishError(queue_manager, real_exp, appenumtypes.PublishFrom_APPLICATION_MANAGER)
			} else if real_exp, ok := r.(error); ok {
				mlog.Errorf("Unknown Error when generating:%v", real_exp)
				queue_manager.PublishError(queue_manager, real_exp, appenumtypes.PublishFrom_APPLICATION_MANAGER)
			}
		}
	}()
	runner := workflowapprunner.NewWorkflowAppRunner(
		application_generate_entity,
		queue_manager,
		// workflow_thread_pool_id=workflow_thread_pool_id,
	)

	runner.Run()
}

func (g *WorkflowAppGenerator[T1]) _generate(
	app_model *models.App,
	wf *models.Workflow,
	user T1, /* Union[Account, EndUser]*/
	application_generate_entity *appgeneratorentities.WorkflowAppGenerateEntity,
	invoke_from appenumtypes.InvokeFrom,
	streaming bool, /* = True*/
) (map[string]any, iter.Seq[string]) {
	// init queue manager
	queue_manager := wfappqueuemanager.New(
		application_generate_entity.TaskID,
		application_generate_entity.UserID,
		application_generate_entity.InvokeFrom,
		string(app_model.Mode),
	)

	// new thread
	go func() {
		g._generate_worker(application_generate_entity, queue_manager)
	}()

	// return response or stream generator
	blocking_response, stream_response := g._handle_response(
		application_generate_entity,
		wf,
		queue_manager,
		user,
		streaming,
	)

	response_converter := &wfappgeneratorreponseconvertor.WorkflowAppGenerateResponseConvert{}
	if blocking_response != nil {
		return response_converter.ConvertBlocking(response_converter, blocking_response, invoke_from), nil
	}
	if stream_response != nil {
		return nil, response_converter.ConvertStream(response_converter, stream_response, invoke_from)
	}
	panic(exceptions.NewValueError("missing response"))
}

func (g *WorkflowAppGenerator[T1]) SingleIterationGenerate(
	app_model *models.App,
	wf *models.Workflow,
	node_id string,
	user T1,
	args map[string]any,
	streaming bool, /* = True*/
) (map[string]any, iter.Seq[string]) {
	/*
		Generate App response.

		:param app_model: App
		:param workflow: Workflow
		:param user: account or end user
		:param args: request args
		:param invoke_from: invoke from source
		:param stream: is stream
	*/
	if node_id == "" {
		panic(exceptions.NewValueError("node_id is required"))
	}
	if _, ok := args["inputs"]; !ok {
		panic(exceptions.NewValueError("inputs is required"))
	}
	if _, ok := args["inputs"].(map[string]any); !ok {
		panic(exceptions.NewValueError("inputs must be dict"))
	}
	user_id := ""
	switch realuser := any(user).(type) {
	case *models.Account:
		user_id = realuser.ID
	case *models.EndUser:
		user_id = realuser.ID
	}

	// convert to app config
	app_config := (&wfappcfgmgr.WorkflowAppConfigManager{}).GetAppConfig(app_model, wf)

	// init application generate entity
	application_generate_entity := &appgeneratorentities.WorkflowAppGenerateEntity{
		AppGenerateEntity: &appgeneratorentities.AppGenerateEntity[*appconfigentities.WorkflowUIBasedAppConfig]{
			TaskID: uuid.NewV4().String(),

			Inputs:     make(map[string]any),
			UserID:     user_id,
			Stream:     streaming,
			InvokeFrom: appenumtypes.InvokeFrom_DEBUGGER,
			Extras:     map[string]any{"auto_generate_conversation_name": false},
			AppConfig:  app_config,
		},

		SingleIterationRun: &appgeneratorentities.SingleIterationRunEntity{
			NodeID: node_id,
			Inputs: args["inputs"].(map[string]any),
		},
		WorkflowRunID: uuid.NewV4().String(),
	}
	// contexts.tenant_id.set(application_generate_entity.app_config.tenant_id)

	return g._generate(
		app_model,
		wf,
		user,
		application_generate_entity,
		appenumtypes.InvokeFrom_DEBUGGER,
		streaming,
	)
}

func (g *WorkflowAppGenerator[T1]) _handle_response(
	application_generate_entity *appgeneratorentities.WorkflowAppGenerateEntity,
	wf *models.Workflow,
	queue_manager appqueueentities.AppQueueManager[*appqueueentities.WorkflowQueueMessage],
	user T1, /*: Union[Account, EndUser]*/
	stream bool, /* = False*/
) (*appresponseentities.WorkflowAppBlockingResponse, iter.Seq[*appresponseentities.WorkflowAppStreamResponse]) {
	/*
		Handle response.
		:param application_generate_entity: application generate entity
		:param workflow: workflow
		:param queue_manager: queue manager
		:param user: account or end user
		:param stream: is stream
		:return:
	*/
	// init generate task pipeline
	defer func() {
		if r := recover(); r != nil {
			if exp, ok := r.(*exceptions.ValueError); ok {
				mlog.Errorf("Fails to process generate task pipeline, task_id: %s", application_generate_entity.TaskID)
				panic(exp)
			} else {
				panic(r)
			}
		}
	}()
	var generate_task_pipeline *wfappgeneratortaskpipeline.WorkflowAppGenerateTaskPipeline
	switch realuser := any(user).(type) {
	case *models.Account:
		generate_task_pipeline = wfappgeneratortaskpipeline.New(
			application_generate_entity,
			wf,
			queue_manager,
			realuser,
			stream,
		)
	case *models.EndUser:
		generate_task_pipeline = wfappgeneratortaskpipeline.New(
			application_generate_entity,
			wf,
			queue_manager,
			realuser,
			stream,
		)
	}

	rsp, generator := generate_task_pipeline.Process()

	return rsp, generator
}

func (g *WorkflowAppGenerator[T1]) Generate(
	app_model *models.App,
	wf *models.Workflow,
	user T1, /* Account | EndUser*/
	args map[string]any,
	invoke_from appenumtypes.InvokeFrom,
	streaming bool, /* = True*/
	call_depth int,
) (map[string]any, iter.Seq[string]) {
	user_id := ""
	if realuser, ok := any(user).(*models.Account); ok {
		user_id = realuser.ID
	} else if realuser, ok := any(user).(*models.EndUser); ok {
		user_id = realuser.ID
	}
	files := []map[string]any{}
	if _, ok := args["files"]; ok {
		if _, ok := args["files"].([]map[string]any); ok {
			files = args["files"].([]map[string]any)
		}
	}

	// parse files
	file_extra_config := (&fileupload.FileUploadConfigManager{}).Convert(wf.FeaturesDict(), false)
	system_files := filefactory.BuildFromMappings(
		files,
		app_model.TenantID,
		file_extra_config,
	)

	// convert to app config
	app_config := (&wfappcfgmgr.WorkflowAppConfigManager{}).GetAppConfig(
		app_model,
		wf,
	)

	var inputs map[string]any
	if _, ok := args["inputs"]; ok {
		if _, ok := args["inputs"].(map[string]any); ok {
			inputs = args["inputs"].(map[string]any)
		}
	}
	inputs = g.PrepareUserInputs(inputs, app_config.Variables, app_model.TenantID)

	workflow_run_id := uuid.NewV4().String()
	// init application generate entity
	application_generate_entity := &appgeneratorentities.WorkflowAppGenerateEntity{
		AppGenerateEntity: &appgeneratorentities.AppGenerateEntity[*appconfigentities.WorkflowUIBasedAppConfig]{
			TaskID:           uuid.NewV4().String(),
			FileUploadConfig: file_extra_config,
			Inputs:           inputs,
			Files:            system_files,
			UserID:           user_id,
			Stream:           streaming,
			InvokeFrom:       invoke_from,
			CallDepth:        call_depth,
			AppConfig:        app_config,
		},

		WorkflowRunID: workflow_run_id,
	}
	// contexts.tenant_id.set(application_generate_entity.app_config.tenant_id)

	return g._generate(
		app_model,
		wf,
		user,
		application_generate_entity,
		invoke_from,
		streaming,
	)
}
