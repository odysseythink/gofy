package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/response"
	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
)

func (s *AdminService) AppWorkflowRunList(ctx context.Context, in *pbapi.AppWorkflowRunListRequest) (out *pbapi.AppWorkflowRunListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AppWorkflowRunList call:%#v", p.Addr.String(), in)

	out = &pbapi.AppWorkflowRunListReply{}
	if in.AppId == "" {
		mlog.Error("app_id not provide")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("app_id not provide")
		return
	}
	if !slices.Contains([]models.AppMode{models.AppMode_ADVANCED_CHAT, models.AppMode_WORKFLOW}, models.AppMode(in.AppMode)) {
		mlog.Error("app_mode not provide")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("app_mode not provide")
		return
	}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.Limit < 1 || in.Limit > 100 {
		mlog.Errorf("limit=%d is out of range[1, 100]")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("limit=%d is out of range[1, 100]"))
		return
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		if models.AppMode(in.AppMode) == models.AppMode_ADVANCED_CHAT {
			result := services.ServiceGroupApp.WorkflowRun.GetPaginateAdvancedChatWorkflowRuns(app_model, map[string]any{"last_id": in.LastId, "limit": int(in.Limit)})
			rsp := response.NewAdvancedChatWorkflowRunPaginationResponse(result)
			bindata, _ := json.Marshal(rsp)
			out.ResponseStr = string(bindata)
		} else {
			result := services.ServiceGroupApp.WorkflowRun.GetPaginateWorkflowRuns(app_model, map[string]any{"last_id": in.LastId, "limit": int(in.Limit)})
			rsp := response.NewWorkflowRunPaginationResponse(result)
			bindata, _ := json.Marshal(rsp)
			out.ResponseStr = string(bindata)
		}
	}()

	return
}

func (s *AdminService) GetWorkflowRunDetail(ctx context.Context, in *pbapi.GetWorkflowRunDetailRequest) (out *pbapi.GetWorkflowRunDetailReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowRunDetail call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowRunDetailReply{}
	if in.AppId == "" {
		mlog.Error("app_id not provide")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("app_id not provide")
		return
	}
	if in.WorkflowRunId == "" {
		mlog.Error("workflow_run_id not provide")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("workflow_run_id not provide")
		return
	}
	var current_user *models.Account
	if in.UserId != "" {
		var exp error
		current_user, exp = services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
		if exp != nil || current_user == nil {
			mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
			out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
			return
		}
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		var app_model *models.App
		if current_user != nil {
			app_model = services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		} else {
			if app_model1, err1 := services.ServiceGroupApp.App.Get(in.AppId); err1 != nil {
				mlog.Errorf("can't get app by app_id(%s) :%v", in.AppId, err1)
				out.Exp = &pbexceptions.HTTPException{
					Code:    "app_not_found",
					Status:  http.StatusNotFound,
					Message: "app not found.",
				}
				return
			} else {
				app_model = &app_model1
			}
		}

		if !in.IsNodeExecution {
			result := services.ServiceGroupApp.WorkflowRun.GetWorkflowRun(app_model, in.WorkflowRunId)
			if result == nil {
				mlog.Errorf("can't get workflow run by app_id(%s) and run_id(%s)", app_model.ID, in.WorkflowRunId)
				out.Exp = &pbexceptions.HTTPException{
					Code:    "workflow_run_not_found",
					Status:  http.StatusNotFound,
					Message: "WorkflowRun not found.",
				}

				return
			}
			rsp := response.NewWorkflowRunDetailResponse(result)
			bindata, _ := json.Marshal(rsp)
			out.ResponseStr = string(bindata)
		} else {
			results := services.ServiceGroupApp.WorkflowRun.GetWorkflowRunNodeExecutions(app_model, in.WorkflowRunId)
			rsp := response.NewWorkflowRunNodeExecutionListResponse(results)
			bindata, _ := json.Marshal(rsp)
			out.ResponseStr = string(bindata)
		}
	}()
	return
}
