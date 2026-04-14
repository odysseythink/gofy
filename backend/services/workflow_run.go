package services

import (
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
)

type WorkflowRunService struct {
}

func (s *WorkflowRunService) Get(id string) (wfr *models.WorkflowRun, err error) {
	wfr = &models.WorkflowRun{}
	err = dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", id).First(wfr).Error
	if err != nil {
		wfr = nil
	}
	return
}

func (s *WorkflowRunService) Message(wfr *models.WorkflowRun) (msg *models.Message, err error) {
	msg = &models.Message{}
	err = dbengine.Instance().DB.Model(&models.Message{}).Where("app_id = ? and workflow_run_id = ?", wfr.AppID, wfr.ID).First(msg).Error
	if err != nil {
		wfr = nil
	}
	return
}

func (s *WorkflowRunService) GetPaginateWorkflowRuns(app_model *models.App, args map[string]any) *response.InfiniteScrollPagination {
	/*
	   Get debug workflow run list
	   Only return triggered_from == debugging

	   :param app_model: app model
	   :param args: request args
	*/
	limit := 20
	if _, ok := args["limit"]; ok {
		if _, ok := args["limit"].(int); ok {
			limit = args["limit"].(int)
		}
	}
	last_id := ""
	if _, ok := args["last_id"]; ok {
		if _, ok := args["last_id"].(string); ok {
			last_id = args["last_id"].(string)
		}
	}
	db := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("tenant_id = ? and app_id = ? and triggered_from = ?", app_model.TenantID, app_model.ID, models.WorkflowRunTriggeredFrom_DEBUGGING)

	var workflow_runs []*models.WorkflowRun
	if last_id != "" {
		last_workflow_run := new(models.WorkflowRun)
		err := db.Where("id = ?", last_id).First(last_workflow_run).Error
		if err != nil {
			mlog.Errorf("get WorkflowRun(%s) failed:%v", last_id, err)
			last_workflow_run = nil
		}

		if last_workflow_run == nil {
			panic(exceptions.NewValueError("Last workflow run not exists"))
		}

		err = db.Where("id != ? and created_at < ?", last_workflow_run.ID, last_workflow_run.CreatedAt).Order("created_at DESC").Limit(limit).Find(&workflow_runs).Error
		if err != nil {
			mlog.Errorf("get workflow_runs failed:%v", err)
			workflow_runs = make([]*models.WorkflowRun, 0)
		}
	} else {
		err := db.Order("created_at DESC").Limit(limit).Find(&workflow_runs).Error
		if err != nil {
			mlog.Errorf("get workflow_runs failed:%v", err)
			workflow_runs = make([]*models.WorkflowRun, 0)
		}
	}
	has_more := false
	if len(workflow_runs) == limit {
		current_page_first_workflow_run := workflow_runs[len(workflow_runs)-1]
		var rest_count int64
		db.Where("id != ? and created_at < ?", current_page_first_workflow_run.ID, current_page_first_workflow_run.CreatedAt).Count(&rest_count)

		if rest_count > 0 {
			has_more = true
		}
	}
	return &response.InfiniteScrollPagination{Data: workflow_runs, Limit: limit, HasMore: has_more}
}

// def __getattr__(self, item):
// 	return getattr(self._workflow_run, item)

func (s *WorkflowRunService) GetPaginateAdvancedChatWorkflowRuns(app_model *models.App, args map[string]any) *response.InfiniteScrollPagination {
	/*
	   Get advanced chat app workflow run list
	   Only return triggered_from == advanced_chat

	   :param app_model: app model
	   :param args: request args
	*/

	pagination := s.GetPaginateWorkflowRuns(app_model, args)

	with_message_workflow_runs := []*models.WorkflowWithMessage{}
	for _, workflow_run := range pagination.Data.([]*models.WorkflowRun) {
		message := workflow_run.Message()
		with_message_workflow_run := models.NewWorkflowWithMessage(workflow_run)
		if message != nil {
			with_message_workflow_run.MessageID = message.ID
			with_message_workflow_run.ConversationID = message.ConversationID
		}
		with_message_workflow_runs = append(with_message_workflow_runs, with_message_workflow_run)
	}
	pagination.Data = with_message_workflow_runs
	return pagination
}

func (s *WorkflowRunService) GetWorkflowRun(app_model *models.App, run_id string) *models.WorkflowRun {
	/*
	   Get workflow run detail

	   :param app_model: app model
	   :param run_id: workflow run id
	*/
	workflow_run := new(models.WorkflowRun)
	if app_model != nil {
		err := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ? and tenant_id = ? and app_id = ?", run_id, app_model.TenantID, app_model.ID).First(workflow_run).Error
		if err != nil {
			mlog.Errorf("get workflow_run failed:%v", err)
			workflow_run = nil
		}
	} else {
		err := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", run_id).First(workflow_run).Error
		if err != nil {
			mlog.Errorf("get workflow_run failed:%v", err)
			workflow_run = nil
		}
	}
	return workflow_run
}
func (s *WorkflowRunService) GetWorkflowRunNodeExecutions(app_model *models.App, run_id string) []*models.WorkflowNodeExecution {
	/*
	   Get workflow run node execution list
	*/
	workflow_run := s.GetWorkflowRun(app_model, run_id)

	if workflow_run == nil {
		return nil
	}
	var node_executions []*models.WorkflowNodeExecution
	if app_model != nil {
		err := dbengine.Instance().DB.Model(&models.WorkflowNodeExecution{}).Where("tenant_id = ? and app_id = ? and workflow_id = ? and triggered_from = ? and workflow_run_id = ?", app_model.TenantID, app_model.ID, workflow_run.WorkflowID, models.WorkflowNodeExecutionTriggeredFrom_WORKFLOW_RUN, run_id).Order("`index` DESC").Find(&node_executions).Error
		if err != nil {
			mlog.Errorf("get workflow_run failed:%v", err)
			node_executions = make([]*models.WorkflowNodeExecution, 0)
		}
	} else {
		err := dbengine.Instance().DB.Model(&models.WorkflowNodeExecution{}).Where("workflow_id = ? and triggered_from = ? and workflow_run_id = ?", workflow_run.WorkflowID, models.WorkflowNodeExecutionTriggeredFrom_WORKFLOW_RUN, run_id).Order("`index` DESC").Find(&node_executions).Error
		if err != nil {
			mlog.Errorf("get workflow_run failed:%v", err)
			node_executions = make([]*models.WorkflowNodeExecution, 0)
		}
	}

	return node_executions
}

func (s *WorkflowRunService) GetWorkflowRunsCount(appModel *models.App, status string, triggeredFrom string) map[string]int64 {
	result := map[string]int64{}

	baseQuery := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("app_id = ?", appModel.ID)
	if triggeredFrom != "" {
		baseQuery = baseQuery.Where("triggered_from = ?", triggeredFrom)
	}

	// Total
	var total int64
	baseQuery.Count(&total)
	result["total"] = total

	// By status
	statuses := []string{"running", "succeeded", "failed", "stopped"}
	for _, st := range statuses {
		var count int64
		baseQuery.Where("status = ?", st).Count(&count)
		result[st] = count
	}

	return result
}
