package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"google.golang.org/grpc"
	"mlib.com/gofy/server/cache"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	_ "mlib.com/gofy/server/core/model_runtime/model_provides"
	dbengine "mlib.com/gofy/server/db_engine"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	_ "mlib.com/gofy/server/events"
	eventhandlers "mlib.com/gofy/server/events/event_handlers"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/gofy/server/utils"
	"mlib.com/gofy/server/utils/validate"
	"mlib.com/mlog"
)

type AppService struct {
	pbapi.UnimplementedAppServer
}

var (
	App AppService
)

func (s *AppService) AppRun(in *pbapi.AppRunRequest, out grpc.ServerStreamingServer[pbapi.AppRunReply]) error {
	mlog.Infof("remote app.AppRun call:%#v", in)

	out_rsp := &pbapi.AppRunReply{}
	var current_user any
	if in.UserId == "" {
		mlog.Errorf("user_id not provide")
		out_rsp.Exp = exceptions.NewUnauthorizedPbHttpExp("user_id not provide")
		return out.Send(out_rsp)
	}

	if in.AppId == "" {
		mlog.Error("missing app id")
		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return out.Send(out_rsp)
	}
	if in.ArgsStr == "" {
		mlog.Error("missing args")
		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("missing args")
		return out.Send(out_rsp)
	}
	if !slices.Contains([]string{"blocking", "streaming"}, in.ResponseMode) {
		mlog.Errorf("response_mode=%s can only be \"blocking\" or \"streaming\"", in.ResponseMode)
		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("response_mode=%s can only be \"blocking\" or \"streaming\"", in.ResponseMode))
		return out.Send(out_rsp)
	}

	args := map[string]any{}
	err := json.Unmarshal([]byte(in.ArgsStr), &args)
	if err != nil {
		mlog.Errorf("json unmarshal args_str=%s to dict failed:%v", in.ArgsStr, err)
		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal args_str=%s to dict failed:%v", in.ArgsStr, err))
		return out.Send(out_rsp)
	}

	streaming := in.ResponseMode == "streaming"
	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(*exceptions.ProviderTokenNotInitError); ok {
					exp1 := httpexceptions.NewProviderNotInitializeError(exp.Error())
					out_rsp.Exp = exp1.ToPbHttpException(exp1)
				} else if _, ok := r.(*exceptions.QuotaExceededError); ok {
					exp1 := httpexceptions.NewProviderQuotaExceededError()
					out_rsp.Exp = exp1.ToPbHttpException(exp1)
				} else if _, ok := r.(*exceptions.ModelCurrentlyNotSupportError); ok {
					exp1 := httpexceptions.NewProviderModelCurrentlyNotSupportError()
					out_rsp.Exp = exp1.ToPbHttpException(exp1)
				} else if exp, ok := r.(*modelruntimeexceptions.InvokeError); ok {
					exp1 := httpexceptions.NewCompletionRequestError(exp.Error())
					out_rsp.Exp = exp1.ToPbHttpException(exp1)
				} else if exp, ok := r.(error); ok {
					out_rsp.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
					out_rsp.Exp = exceptions.NewInternalServerPbHttpExp(fmt.Sprintf("%v", r))
				}
			}
		}()

		app_model := new(models.App)
		if err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", in.AppId).First(app_model).Error; err != nil {
			mlog.Errorf("get App=%s failed:%v", in.AppId, err)
			out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("get App=%s failed:%v", in.AppId, err))
			return
		}
		if in.IsDraft {
			current_account_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
			if exp != nil {
				mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
				out_rsp.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
				return
			}
			if !current_account_user.IsEditor() {
				mlog.Errorf("current user(%#v) is not editor", current_account_user)
				out_rsp.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_account_user))
				return
			}
			current_user = current_account_user
		} else {
			current_user = models.CreateOrGetEndUserByUserID(app_model, in.UserId)
		}
		switch app_model.Mode {
		case models.AppMode_WORKFLOW:
			err := validate.StringMapTypeVerify(args, validate.Rules{
				"inputs": {validate.RuleTypeOfField(reflect.Map), validate.NotEmpty()},
				"files":  {validate.RuleTypeOfField(reflect.Slice)},
			})
			if err != nil {
				mlog.Errorf("args validate failed:%v", err)
				out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("args validate failed:" + err.Error())
				return
			}
		case models.AppMode_ADVANCED_CHAT:
			err := validate.StringMapTypeVerify(args, validate.Rules{
				"inputs": {validate.RuleTypeOfField(reflect.Map), validate.NotEmpty()},
				"query":  {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
				"files":  {validate.RuleTypeOfField(reflect.Slice)},
			})
			if err != nil {
				mlog.Errorf("args validate failed:%v", err)
				out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("args validate failed:" + err.Error())
				return
			}
		default:
			mlog.Errorf("unsupported app mode=%s", app_model.Mode)
			out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("unsupported app mode=" + string(app_model.Mode))
			return
		}

		rsp, rspiter := services.ServiceGroupApp.AppGenerate.Generate(app_model, current_user, args, appenumtypes.InvokeFrom_SERVICE_API, streaming)
		mlog.Debugf("***********realrsp=%#v", rsp)
		if rsp != nil {
			bindata, _ := json.Marshal(rsp)

			out_rsp.DirectReplyDictStr = string(bindata)
			return
		}
		for item := range rspiter {
			mlog.Debugf("-----send:%s", string(item))
			out_rsp.StreamReplyStr = item
			err = out.Send(out_rsp)
			if err != nil {
				mlog.Errorf("send stream error:%v", err)
				out_rsp = nil
				return
			}
		}
	}()
	if err != nil {
		return err
	}
	if out_rsp != nil {
		return out.Send(out_rsp)
	}
	return nil
}

// go build -o app.so -buildmode=plugin main.go
func (s *AppService) Init(args ...any) error {
	mlog.Info("app init......")
	err := cache.Instance().Init()
	if err != nil {
		mlog.Errorf("cahe init failed:%v", err)
		return fmt.Errorf("cahe init failed:%v", err)
	}
	err = dbengine.Instance().Init()
	if err != nil {
		mlog.Errorf("mysql init failed:%v", err)
		return fmt.Errorf("mysql init failed:%v", err)
	}
	eventhandlers.Init()

	return nil
}

func (s *AppService) RunOnce(ctx context.Context) error {
	return nil
}
func (s *AppService) Destroy() {

}

func (s *AppService) UserData() any {
	return nil
}
