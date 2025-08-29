package services

import (
	"fmt"
	"iter"

	"github.com/spf13/viper"
	achatgenerator "mlib.com/gofy/server/core/app/generatores/advanced_chat"
	wfgenerator "mlib.com/gofy/server/core/app/generatores/workflow"
	"mlib.com/gofy/server/core/exceptions"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
)

type AppGenerateService struct {
}

func (s *AppGenerateService) getMaxActiveRequests(app_model *models.App) int {
	max_active_requests := app_model.MaxActiveRequests
	if max_active_requests <= 0 {
		max_active_requests = viper.GetInt("app_config.max_active_requests")
	}
	return max_active_requests
}

func (s *AppGenerateService) getWorkflow(app_model *models.App, invoke_from appenumtypes.InvokeFrom) *models.Workflow {
	/*
		Get workflow
		:param app_model: app model
		:param invoke_from: invoke from
		:return:
	*/
	var wf *models.Workflow
	workflow_service := &WorkflowService{}
	if invoke_from == appenumtypes.InvokeFrom_DEBUGGER {
		// fetch draft workflow by app_model
		wf = workflow_service.GetDraftWorkflow(app_model)

		if wf == nil {

		}
	} else {
		// fetch published workflow by app_model
		wf = workflow_service.GetPublishedWorkflow(app_model)

		if wf == nil {
			panic(exceptions.NewValueError("Workflow not published"))
		}
	}

	return wf
}
func (s *AppGenerateService) GenerateSingleIteration(app_model *models.App, user *models.Account, node_id string, args any, streaming bool /* = True*/) {
	if app_model.Mode == models.AppMode_ADVANCED_CHAT {
		// wf := s.getWorkflow(app_model, appenumtypes.InvokeFrom_DEBUGGER)
		// return AdvancedChatAppGenerator().single_iteration_generate(
		// 	app_model=app_model,
		// 	workflow=workflow,
		// 	node_id=node_id,
		// 	user=user,
		// 	args=args,
		// 	streaming=streaming,
		// )
	} else if app_model.Mode == models.AppMode_WORKFLOW {
		// wf := s.getWorkflow(app_model, appenumtypes.InvokeFrom_DEBUGGER)
		// return WorkflowAppGenerator().single_iteration_generate(
		// 	app_model, wf, node_id, user, args, streaming,
		// )
	} else {
		// raise ValueError(f"Invalid app mode {app_model.mode}")
	}
}

func (s *AppGenerateService) Generate(
	app_model *models.App,
	user any, /* Account | EndUser*/
	args map[string]any,
	invoke_from appenumtypes.InvokeFrom,
	streaming bool, /* = True*/
) (map[string]any, iter.Seq[string]) {
	/*
		App Content Generate
		:param app_model: app model
		:param user: user
		:param args: args
		:param invoke_from: invoke from
		:param streaming: streaming
		:return:
	*/
	// max_active_request = AppGenerateService._get_max_active_requests(app_model)
	// rate_limit = RateLimit(app_model.id, max_active_request)
	// request_id = RateLimit.gen_request_key()
	// try:
	// request_id = rate_limit.enter(request_id)
	switch app_model.Mode {
	case models.AppMode_WORKFLOW:
		wf := s.getWorkflow(app_model, invoke_from)
		switch realuser := user.(type) {
		case *models.Account:
			return (wfgenerator.New[*models.Account]()).Generate(
				app_model,
				wf,
				realuser,
				args,
				invoke_from,
				streaming,
				0,
			)
		case *models.EndUser:
			return (wfgenerator.New[*models.EndUser]()).Generate(
				app_model,
				wf,
				realuser,
				args,
				invoke_from,
				streaming,
				0,
			)
		}

	case models.AppMode_ADVANCED_CHAT:
		wf := s.getWorkflow(app_model, invoke_from)
		switch realuser := user.(type) {
		case *models.Account:
			return achatgenerator.New[*models.Account]().Generate(
				app_model,
				wf,
				realuser,
				args,
				invoke_from,
				streaming,
			)
		case *models.EndUser:
			return achatgenerator.New[*models.EndUser]().Generate(
				app_model,
				wf,
				realuser,
				args,
				invoke_from,
				streaming,
			)
		}
	case models.AppMode_AGENT_CHAT:

	}
	panic(exceptions.NewValueError(fmt.Sprintf("Invalid app mode %s", app_model.Mode)))
}
