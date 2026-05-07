package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
)

func (s *AdminService) Statistic(ctx context.Context, in *pbapi.StatisticRequest) (out *pbapi.StatisticReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.WorkflowDailyRunsStatistic call:%#v", p.Addr.String(), in)

	out = &pbapi.StatisticReply{}
	if in.Method == "" {
		mlog.Errorf("method not provide")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("method not provide")
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
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	start, err := time.Parse("2006-01-02 15:04", in.Start)
	if err != nil {
		mlog.Errorf("parse start time(%s) failed:%v", in.Start, err)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("parse start time(%s) failed:%v", in.Start, err))
		return
	}
	end, err := time.Parse("2006-01-02 15:04", in.End)
	if err != nil {
		mlog.Errorf("parse end time(%s) failed:%v", in.End, err)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("parse start time(%s) failed:%v", in.Start, err))
		return
	}

	rsp := func() []map[string]any {
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
		return services.ServiceGroupApp.Statistic.Statistic(in.Method, app_model, start, end)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.ResultStr = string(bindata)
	return
}
