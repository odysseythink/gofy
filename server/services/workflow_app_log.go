package services

import (
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/mlog"
)

type WorkflowAppLogService struct {
}

func (s *WorkflowAppLogService) Get(id string) (wfr *models.WorkflowAppLog, err error) {
	wfr = &models.WorkflowAppLog{}
	err = dbengine.Instance().DB.Model(&models.WorkflowAppLog{}).Where("id = ?", id).Preload("Tenant").Preload("App").Preload("CreatedByAccount").Preload("Workflow").Preload("WorkflowRun").First(wfr).Error
	if err != nil {
		wfr = nil
	}
	return
}
func (s *WorkflowAppLogService) CreatedByAccount(id string) bool {
	wfr, err := s.Get(id)
	if err != nil {
		mlog.Errorf("Get workflow run error: %v", err)
		return false
	}
	return wfr.CreatedByRole == models.CreatedByRole_ACCOUNT
}
func (s *WorkflowAppLogService) CreatedByEndUser(id string) bool {
	wfr, err := s.Get(id)
	if err != nil {
		mlog.Errorf("Get workflow run error: %v", err)
		return false
	}
	return wfr.CreatedByRole == models.CreatedByRole_END_USER
}

func (s *WorkflowAppLogService) GetPaginateWorkflowAppLogs(app_model *models.App, page, limit int, keyword, status string) *response.WorkflowAppLogPaginationResponse {
	db := dbengine.Instance().DB.Model(&models.WorkflowAppLog{}).Joins("JOIN workflow_runs ON workflow_runs.id = workflow_app_logs.workflow_run_id").Where(
		"workflow_app_logs.tenant_id = ? and workflow_app_logs.app_id = ?",
		app_model.TenantID,
		app_model.ID,
	)

	if keyword != "" {
		keyword_like_val := ""
		if len(keyword) > 30 {
			keyword_like_val = `%` + keyword[:30] + `%`
		} else {
			keyword_like_val = `%` + keyword + `%`
		}
		db = db.Where(
			"workflow_runs.inputs like ? or workflow_runs.outputs like ?",
			keyword_like_val,
			keyword_like_val,
		)
		db = db.Joins("left outer join end_users ON workflow_runs.created_by_role = ? and workflow_runs.created_by = end_users.id and end_users.session_id like ?", models.CreatedByRole_END_USER, keyword_like_val)
	}
	if status != "" {
		// join with workflow_run and filter by status
		db = db.Where(
			"workflow_runs.status = ?",
			status,
		)
	}

	db.Select(
		"workflow_app_logs.id as workflow_app_logs_id",
		"workflow_app_logs.created_from",
		"workflow_app_logs.created_by_role",
		"workflow_app_logs.created_at as workflow_app_logs_created_at",
		"workflow_app_logs.created_by",
		"workflow_runs.id as workflow_runs_id",
		"workflow_runs.version",
		"workflow_runs.status",
		"workflow_runs.error",
		"workflow_runs.elapsed_time",
		"workflow_runs.total_tokens",
		"workflow_runs.total_steps",
		"workflow_runs.created_at as workflow_runs_created_at",
		"workflow_runs.finished_at",
		"workflow_runs.exceptions_count",
	)

	db = db.Order("workflow_app_logs.created_at DESC")
	rsp := &response.WorkflowAppLogPaginationResponse{Limit: limit, Page: page, Data: make([]*response.WorkflowAppLogPartialResponse, 0)}
	err := db.Count(&rsp.Total).Error
	if err != nil {
		mlog.Errorf("get workflow app log failed:%v", err)
		return rsp
	}
	offset := limit * (page - 1)
	db = db.Limit(limit).Offset(offset)
	type Data struct {
		ID                    string     `gorm:"column:workflow_app_logs_id"`
		CreatedFrom           string     `gorm:"column:created_from"`
		CreatedByRole         string     `gorm:"column:created_by_role"`
		CreatedAt             *time.Time `gorm:"column:workflow_app_logs_created_at"`
		CreatedBy             string     `gorm:"column:created_by"`
		WorkflowRunID         string     `gorm:"column:workflow_runs_id"`
		Version               string     `gorm:"column:version"`
		Status                string     `gorm:"column:status"`
		Error                 string     `gorm:"column:error"`
		ElapsedTime           float64    `gorm:"column:elapsed_time"`
		TotalTokens           int        `gorm:"column:total_tokens"`
		TotalSteps            int        `gorm:"column:total_steps"`
		WorkflowRunCreatedAt  *time.Time `gorm:"column:workflow_runs_created_at"`
		WorkflowRunFinishedAt *time.Time `gorm:"column:finished_at"`
		ExceptionsCount       int        `gorm:"column:exceptions_count"`
	}
	datas := []Data{}
	err = db.Scan(&datas).Error
	if err != nil {
		mlog.Errorf("get workflow app log failed:%v", err)
		return rsp
	}
	for _, v := range datas {
		mlog.Debugf("------v=%#v", v)
		sv := &response.WorkflowAppLogPartialResponse{
			ID:            v.ID,
			CreatedFrom:   v.CreatedFrom,
			CreatedByRole: v.CreatedByRole,
			// CreatedByAccount: v.CreatedByAccount,
			// CreatedByEndUser: v.CreatedByEndUser,
			// CreatedAt: v.CreatedAt.Unix(),
		}
		if v.CreatedAt != nil {
			sv.CreatedAt = v.CreatedAt.Unix()
		}
		if v.WorkflowRunID != "" {
			sv.WorkflowRun = &response.WorkflowRunForLogResponse{
				ID:              v.WorkflowRunID,
				Version:         v.Version,
				Status:          v.Status,
				Error:           v.Error,
				ElapsedTime:     v.ElapsedTime,
				TotalTokens:     v.TotalTokens,
				TotalSteps:      v.TotalSteps,
				CreatedAt:       v.WorkflowRunCreatedAt.Unix(),
				FinishedAt:      v.WorkflowRunCreatedAt.Unix(),
				ExceptionsCount: v.ExceptionsCount,
			}
		}
		if v.CreatedByRole == "account" {
			acc := new(models.Account)
			err := dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", v.CreatedBy).First(acc).Error
			if err != nil {
				mlog.Errorf("get account failed:%v", err)
				acc = nil
			}
			if acc != nil {
				sv.CreatedByAccount = &response.SimpleAccountResponse{
					ID:    acc.ID,
					Name:  acc.Name,
					Email: acc.Email,
				}
			}
		} else if v.CreatedByRole == "end_user" {
			user := new(models.EndUser)
			err := dbengine.Instance().DB.Model(&models.EndUser{}).Where("id = ?", v.CreatedBy).First(user).Error
			if err != nil {
				mlog.Errorf("get account failed:%v", err)
				user = nil
			}
			if user != nil {
				sv.CreatedByEndUser = &response.SimpleEndUserResponse{
					ID:          user.ID,
					Type:        user.Type,
					IsAnonymous: user.IsAnonymous,
					SessionID:   user.SessionID,
				}
			}
		}
		rsp.Data = append(rsp.Data, sv)
	}

	return rsp
}
