package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

func (s *AdminService) GetWorkflowAppLogList(ctx context.Context, in *pbapi.GetWorkflowAppLogListRequest) (out *pbapi.GetWorkflowAppLogListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowAppLogList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowAppLogListReply{}
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
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.Status != "" {
		if !slices.Contains([]string{"succeeded", "failed", "stopped"}, in.Status) {
			mlog.Errorf("status[%s] must be in set[\"succeeded\", \"failed\", \"stopped\"]", in.Status)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("status[%s] must be in set[\"succeeded\", \"failed\", \"stopped\"]", in.Status))
			return
		}
	}
	if in.Page < 1 || in.Page > 99999 {
		mlog.Errorf("page[%d] must be in range[1, 99999]", in.Page)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("page[%d] must be in range[1, 99999]", in.Page))
		return
	}
	if in.Limit < 1 || in.Limit > 100 {
		mlog.Errorf("limit[%d] must be in range[1, 100]", in.Limit)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("limit[%d] must be in range[1, 100]", in.Limit))
		return
	}
	info := func() *response.WorkflowAppLogPaginationResponse {
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
		return services.ServiceGroupApp.WorkflowAppLog.GetPaginateWorkflowAppLogs(app_model, int(in.Page), int(in.Limit), in.Keyword, in.Status)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(info)
	out.WorkflowAppLogPaginationStr = string(bindata)

	return
}
