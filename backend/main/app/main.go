package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	modelruntimeexceptions "github.com/odysseythink/gofy/backend/core/exceptions/model_runtime"
	_ "github.com/odysseythink/gofy/backend/core/model_runtime/model_provides"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appenumtypes "github.com/odysseythink/gofy/backend/enum_types/app"
	_ "github.com/odysseythink/gofy/backend/events"
	eventhandlers "github.com/odysseythink/gofy/backend/events/event_handlers"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/gofy/backend/utils/validate"
	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"
	"google.golang.org/grpc"
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

func main() {
	var cfgfile string
	flag.StringVar(&cfgfile, "c", "", "choose config file.")
	flag.Parse()
	if cfgfile == "" {
		log.Println("usage: ./server -c config.yml")
		return
	}
	confy.SetConfigFile(cfgfile)
	confy.SetConfigType("yaml")
	err := confy.ReadInConfig()
	if err != nil {
		log.Printf("read config file(%s) failed: %v\n", cfgfile, err)
		return
	}
	{
		logpath := confy.GetWithDefault[string]("log.path", "logs")
		loglevel := confy.GetWithDefault[uint32]("log.log_level", 1)
		log.Println("******loglevel=", loglevel)
		if loglevel >= 4 {
			loglevel = 1
		}
		mlog.SetLogLevel(loglevel)
		mlog.SetLogDir(logpath)
	}
	defer mlog.Flush()
	confy.WatchConfig()

	mrun.Register(cluster.Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)}, []any{&App})

	err = mrun.Run(&App)
	mlog.Infof("%s Server End!:%v", os.Args[0], err)
}
