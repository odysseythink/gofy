package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	enumtypes "mlib.com/gofy/server/enum_types"
	"mlib.com/gofy/server/models"
	pbexceptions "mlib.com/gofy/server/proto/exceptions"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

func (s *AdminService) GetWorkspaceList(ctx context.Context, in *pbapi.GetWorkspaceListRequest) (out *pbapi.GetWorkspaceListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkspaceList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkspaceListReply{}
	limit := int(in.Page)
	offset := int(in.Limit * (in.Page - 1))
	// 创建db
	db := dbengine.Instance().DB.Model(&models.Tenant{})
	var ts []models.Tenant

	err = db.Count(&out.Total).Error
	if err != nil {
		mlog.Errorf("count tenant failed:%v", err)
		out.Result = "invalid data"
		out.Code = 7
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&ts).Error
	if err != nil {
		mlog.Errorf("get tenant failed:%v", err)
		out.Result = "invalid data"
		out.Code = 7
		return
	}

	if len(ts) == int(in.Limit) {
		current_page_first_tenant := ts[int(in.Limit)-1]

		rest_count := services.ServiceGroupApp.Tenant.RestCount(&current_page_first_tenant)

		if rest_count > 0 {
			out.HasMore = true
		}
	}
	for _, v := range ts {
		if out.Data == nil {
			out.Data = make([]*pbapi.WorkspaceResponse, 0)
		}
		ws := &pbapi.WorkspaceResponse{
			Id:     v.ID,
			Name:   v.Name,
			Status: string(v.Status),
		}
		if v.CreatedAt != nil {
			ws.CreatedAt = v.CreatedAt.Unix()
		}
		out.Data = append(out.Data, ws)
	}
	out.Result = "success"
	out.Code = 0
	return
}

func (s *AdminService) GetTenantList(ctx context.Context, in *pbapi.GetTenantListRequest) (out *pbapi.GetTenantListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetTenantList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetTenantListReply{}
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
	ts := services.ServiceGroupApp.Tenant.GetJoinTenants(current_user)
	for idx, v := range ts {
		if v.ID == current_user.CurrentTenantID() {
			ts[idx].Current = true
		}
		if out.Workspaces == nil {
			out.Workspaces = make([]*pbapi.TenantResponse, 0)
		}
		tmp := &pbapi.TenantResponse{
			Id:      v.ID,
			Name:    v.Name,
			Plan:    v.Plan,
			Status:  string(v.Status),
			Current: v.Current,
		}
		if v.CreatedAt != nil {
			tmp.CreatedAt = v.CreatedAt.Unix()
		}
		out.Workspaces = append(out.Workspaces, tmp)
	}
	return
}

func (s *AdminService) GetCurrentTenant(ctx context.Context, in *pbapi.GetCurrentTenantRequest) (out *pbapi.GetCurrentTenantReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetCurrentTenant call:%#v", p.Addr.String(), in)

	out = &pbapi.GetCurrentTenantReply{}
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
	if current_user.CurrentTenant().Status == enumtypes.TenantStatus_ARCHIVE {
		tenants := services.ServiceGroupApp.Tenant.GetJoinTenants(current_user)
		// # if there is any tenant, switch to the first one
		if len(tenants) > 0 {
			services.ServiceGroupApp.Tenant.SwitchTenant(current_user, tenants[0].ID)
			// tenant = tenants[0]
			services.ServiceGroupApp.Account.SetCurrentTenant(current_user, tenants[0])
			// # else, raise Unauthorized
		} else {
			mlog.Errorf("workspace is archived")
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusUnauthorized,
				Message: "workspace is archived",
			}
			return
		}
	}
	result := services.ServiceGroupApp.Tenant.GetTenantInfo(current_user.CurrentTenant(), current_user)
	if result == nil {
		result = make(map[string]any)
	}
	bindata, _ := json.Marshal(result)
	out.ResultDictStr = string(bindata)
	return
}
