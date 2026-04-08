package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/peer"
	"mlib.com/confy"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/global"
	"mlib.com/gofy/server/libs/password"
	pbexceptions "mlib.com/gofy/server/proto/exceptions"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

func (s *AdminService) GetSetupStatus(ctx context.Context, in *pbapi.GetSetupStatusRequest) (out *pbapi.GetSetupStatusReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetSetupStatus call:%#v", p.Addr.String(), in)
	out = &pbapi.GetSetupStatusReply{
		Step: "finished",
	}
	if confy.Get[string]("EDITION") == "SELF_HOSTED" {
		setup_status, _ := services.ServiceGroupApp.DifySetup.Get()
		if setup_status != nil {
			if setup_status.SetupAt == nil {
				now := time.Now()
				setup_status.SetupAt = &now
				services.ServiceGroupApp.DifySetup.Update(setup_status)
			}
			out.SetupAt = setup_status.SetupAt.Format("2006-01-02T15:04:05.999999")
		} else {
			out.Step = "not_started"
		}
	}
	return
}
func (s *AdminService) Setup(ctx context.Context, in *pbapi.SetupRequest) (out *pbapi.SetupReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.Setup call:%#v", p.Addr.String(), in)
	out = &pbapi.SetupReply{}
	setup_status, _ := services.ServiceGroupApp.DifySetup.Get()
	if setup_status == nil {
		exp := httpexceptions.NewAlreadySetupError()
		out.Exp = exp.ToPbHttpException(exp)
		return
	}
	tenant_count := services.ServiceGroupApp.Tenant.GetTenantCount()
	if tenant_count > 0 {
		exp := httpexceptions.NewAlreadySetupError()
		out.Exp = exp.ToPbHttpException(exp)
		return
	}
	if !s.get_init_validate_status() {
		exp := httpexceptions.NewNotInitValidateError()
		out.Exp = exp.ToPbHttpException(exp)
		return
	}
	if in.Email == "" {
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Code:    "missing_email",
			Message: "missing email",
		}
		return
	}
	if in.Name == "" {
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Code:    "missing_name",
			Message: "missing name",
		}
		return
	}
	if in.Password == "" {
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Code:    "missing_password",
			Message: "missing password",
		}
		return
	}
	_, err = password.ValidPassword(in.Password)
	if err != nil {
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Code:    "invalid_password",
			Message: "invalid password:" + err.Error(),
		}
		return
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				out.Exp = &pbexceptions.HTTPException{
					Status:  http.StatusInternalServerError,
					Code:    "internal_server_error",
					Message: fmt.Sprintf("internal server error:%v", r),
				}
			}
		}()
		services.ServiceGroupApp.Register.Setup(in.Email, in.Name, in.Password, in.ClientIp)
	}()
	return
}

func (s *AdminService) get_init_validate_status() bool {
	if confy.Get[string]("EDITION") == "SELF_HOSTED" {
		if os.Getenv("INIT_PASSWORD") == "on" {
			setup_status, _ := services.ServiceGroupApp.DifySetup.Get()
			return global.IsInitValidated || setup_status != nil
		}
	}

	return true
}
