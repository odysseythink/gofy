package main

import (
	"context"
	"net/http"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	"github.com/odysseythink/gofy/backend/models"
	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
)

func (s *AdminService) Logout(ctx context.Context, in *pbapi.LogoutRequest) (out *pbapi.LogoutReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.Logout call:%#v", p.Addr.String(), in)
	out = &pbapi.LogoutReply{}
	account, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusUnauthorized,
			Code:    "unauthorized",
			Message: "Unauthorized",
		}
		return
	}
	services.ServiceGroupApp.Account.Logout(account)
	return
}

func (s *AdminService) Login(ctx context.Context, in *pbapi.LoginRequest) (out *pbapi.LoginReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.Login call:%#v", p.Addr.String(), in)
	out = &pbapi.LoginReply{Result: ""}
	if in.Email == "" || in.Password == "" {
		mlog.Error("invalid args")
		out.Result = "fail"
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Code:    "email_and_password_is_required",
			Message: "email and password is required",
		}
		return
	}

	is_login_error_rate_limit := services.ServiceGroupApp.Account.IsLoginErrorRateLimit(in.Email)
	if is_login_error_rate_limit {
		mlog.Error("Too many incorrect password attempts. Please try again later.")
		out.Result = "fail"
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusTooManyRequests,
			Code:    "email_code_login_limit",
			Message: "Too many incorrect password attempts. Please try again later.",
		}

		return
	}

	invitation := map[string]any{}
	if in.InviteToken != "" {
		invitation = services.ServiceGroupApp.Register.GetInvitationIfTokenValid("", in.Email, in.InviteToken)
	}
	// language := ""
	// if in.Language == "zh-Hans" {
	// 	// if args["language"] is not None and args["language"] == "zh-Hans":
	// 	language = "zh-Hans"
	// } else {
	// 	language = "en-US"
	// }

	account := func() *models.Account {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if _, ok := r.(*exceptions.AccountLoginError); ok {
					exp := httpexceptions.NewAccountBannedError()
					out.Result = "fail"
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.AccountPasswordError); ok {
					services.ServiceGroupApp.Account.AddLoginErrorRateLimit(in.Email)
					exp := httpexceptions.NewEmailOrPasswordMismatchError()
					out.Result = "fail"
					out.Exp = exp.ToPbHttpException(exp)
				} else if _, ok := r.(*exceptions.AccountNotFoundError); ok {
					exp := httpexceptions.NewAccountNotFound()
					out.Result = "fail"
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Result = "fail"
					out.Exp = exp.ToPbHttpException(exp)
				} else {
					out.Result = "fail"
					out.Exp = &pbexceptions.HTTPException{
						Status:  http.StatusInternalServerError,
						Code:    "internal_server_error",
						Message: "internal server error",
					}
				}
			}
		}()
		if len(invitation) > 0 {
			data := map[string]any{}
			if _, ok := invitation["data"]; ok {
				if _, ok := invitation["data"].(map[string]any); ok {
					data = invitation["data"].(map[string]any)
				}
			}

			invitee_email := ""
			if _, ok := data["email"]; ok {
				if _, ok := data["email"].(string); ok {
					invitee_email = data["email"].(string)
				}
			}
			if invitee_email != in.Email {
				panic(httpexceptions.NewInvalidEmailError())
			}
			return services.ServiceGroupApp.Account.Authenticate(in.Email, in.Password, in.InviteToken)
		} else {
			return services.ServiceGroupApp.Account.Authenticate(in.Email, in.Password, "")
		}
	}()
	if out.Exp != nil {
		return
	}

	// # SELF_HOSTED only have one workspace
	tenants := services.ServiceGroupApp.Tenant.GetJoinTenants(account)
	if len(tenants) == 0 {
		out.Errmsg = "workspace not found, please contact system admin to invite you to join in a workspace"
		out.Result = "fail"
		return
	}
	token_pair := services.ServiceGroupApp.Account.Login(account, in.ClientIp)
	services.ServiceGroupApp.Account.ResetLoginErrorRateLimit(in.Email)
	out.Result = "success"
	out.AccessToken = token_pair.AccessToken
	out.RefreshToken = token_pair.RefreshToken

	return
}
