package main

import (
	"context"
	"encoding/json"

	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
)

func (s *AdminService) GetMemberList(ctx context.Context, in *pbapi.GetMemberListRequest) (out *pbapi.GetMemberListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetMemberList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetMemberListReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
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
		members := services.ServiceGroupApp.Tenant.GetTenantMembers(in.TenantId)
		bindata, _ := json.Marshal(response.NewAccountWithRoleListResponse(members))
		out.MembersStr = string(bindata)
	}()
	return
}
