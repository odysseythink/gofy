package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
)

func (s *AdminService) SetDefaultModel(ctx context.Context, in *pbapi.SetDefaultModelRequest) (out *pbapi.SetDefaultModelReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.SetDefaultModel call:%#v", p.Addr.String(), in)

	out = &pbapi.SetDefaultModelReply{}
	if in.TenantId == "" {
		mlog.Errorf("TenantId[%s] not provide", in.TenantId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("TenantId[%s] not provide", in.TenantId))
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
		for _, model_setting := range in.ModelSettings {
			if !modelruntimeenumtypes.ModelType(model_setting.ModelType).Valid() {
				panic(exceptions.NewValueError("invalid model type"))
			}
			if model_setting.Provider == "" || model_setting.Model == "" {
				continue
			}

			services.ServiceGroupApp.ModelProvide.UpdateDefaultModelOfModelType(
				in.TenantId,
				model_setting.Provider,
				model_setting.Model,
				modelruntimeenumtypes.ModelType(model_setting.ModelType),
			)
		}
	}()

	return
}
func (s *AdminService) GetDefaultModel(ctx context.Context, in *pbapi.GetDefaultModelRequest) (out *pbapi.GetDefaultModelReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetDefaultModel call:%#v", p.Addr.String(), in)

	out = &pbapi.GetDefaultModelReply{}
	if !modelruntimeenumtypes.ModelType(in.ModelType).Valid() {
		mlog.Errorf("invalid model type=%s", in.ModelType)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("invalid model type=%s", in.ModelType))
		return
	}
	if in.TenantId == "" {
		mlog.Errorf("TenantId[%s] not provide", in.TenantId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("TenantId[%s] not provide", in.TenantId))
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
		default_model_entity := services.ServiceGroupApp.ModelProvide.GetDefaultModelOfModelType(in.TenantId, modelruntimeenumtypes.ModelType(in.ModelType))
		bindata, _ := json.Marshal(default_model_entity)
		out.DefaultModelResponseStr = string(bindata)
	}()
	return
}
func (s *AdminService) GetAvailableModelProvider(ctx context.Context, in *pbapi.GetAvailableModelProviderRequest) (out *pbapi.GetAvailableModelProviderReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetAvailableModelProvider call:%#v", p.Addr.String(), in)

	out = &pbapi.GetAvailableModelProviderReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if !modelruntimeenumtypes.ModelType(in.ModelType).Valid() {
		mlog.Errorf("model_type=%s is invalid", in.ModelType)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("model_type=%s is invalid", in.ModelType))
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
		models := services.ServiceGroupApp.ModelProvide.GetModelsByModelType(in.TenantId, modelruntimeenumtypes.ModelType(in.ModelType))
		bindata, _ := json.Marshal(models)
		out.ModelsStr = string(bindata)
	}()

	return
}

func (s *AdminService) GetModelProviderModel(ctx context.Context, in *pbapi.GetModelProviderModelRequest) (out *pbapi.GetModelProviderModelReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetModelProviderModel call:%#v", p.Addr.String(), in)

	out = &pbapi.GetModelProviderModelReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if in.Provider == "" {
		mlog.Error("missing provider")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing provider")
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
		models := services.ServiceGroupApp.ModelProvide.GetModelsByProvider(in.TenantId, in.Provider)
		bindata, _ := json.Marshal(models)
		out.ModelsStr = string(bindata)
	}()

	return
}
func (s *AdminService) SetModelProviderModel(ctx context.Context, in *pbapi.SetModelProviderModelRequest) (out *pbapi.SetModelProviderModelReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.SetModelProviderModel call:%#v", p.Addr.String(), in)

	out = &pbapi.SetModelProviderModelReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if in.Provider == "" {
		mlog.Error("missing provider")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing provider")
		return
	}
	if !modelruntimeenumtypes.ModelType(in.ModelType).Valid() {
		mlog.Errorf("invalid model type=%s", in.ModelType)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("invalid model type=%s", in.ModelType))
		return
	}
	if in.Model == "" {
		mlog.Error("missing model")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing model")
		return
	}
	var credentials map[string]any
	if in.CredentialsDictStr != "" {
		if err1 := json.Unmarshal([]byte(in.CredentialsDictStr), &credentials); err1 != nil {
			mlog.Errorf("json unmarshal credentials_dict_str=%s to dict failed:%v", in.CredentialsDictStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal credentials_dict_str=%s to dict failed:%v", in.CredentialsDictStr, err1))
			return
		}
	}
	var load_balancing map[string]any
	if in.LoadBalancingDictStr != "" {
		if err1 := json.Unmarshal([]byte(in.LoadBalancingDictStr), &load_balancing); err1 != nil {
			mlog.Errorf("json unmarshal load_balancing_str=%s to dict failed:%v", in.LoadBalancingDictStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal load_balancing_str=%s to dict failed:%v", in.LoadBalancingDictStr, err1))
			return
		}
	}
	load_balancing_enabled := false
	if len(load_balancing) > 0 {
		if _, ok := load_balancing["enabled"]; ok {
			if _, ok := load_balancing["enabled"].(bool); ok {
				load_balancing_enabled = load_balancing["enabled"].(bool)
			}
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
		if load_balancing_enabled {
			if _, ok := load_balancing["configs"]; !ok {
				mlog.Errorf("when enabled load balancing, must provide load bancing configs")
				out.Exp = &pbexceptions.HTTPException{
					Status:  http.StatusBadRequest,
					Message: "when enabled load balancing, must provide load bancing configs",
				}
				return
			}
			if _, ok := load_balancing["configs"].([]map[string]any); !ok {
				mlog.Errorf("load bancing configs=%#v must be dict")
				out.Exp = &pbexceptions.HTTPException{
					Status:  http.StatusBadRequest,
					Message: fmt.Sprintf("load bancing configs=%#v must be dict"),
				}
				return
			}
			configs := load_balancing["configs"].([]map[string]any)
			// save load balancing configs
			services.ServiceGroupApp.ModelLoadBalancing.UpdateLoadBalancingConfigs(
				in.TenantId,
				in.Provider,
				in.Model,
				modelruntimeenumtypes.ModelType(in.ModelType),
				configs,
			)

			// enable load balancing
			if err1 := services.ServiceGroupApp.ModelLoadBalancing.EnableModelLoadBalancing(
				in.TenantId, in.Provider, in.Model, modelruntimeenumtypes.ModelType(in.ModelType),
			); err1 != nil {
				mlog.Error("EnableModelLoadBalancing failed:", err1)
				out.Exp = &pbexceptions.HTTPException{
					Status:  http.StatusBadRequest,
					Message: "EnableModelLoadBalancing failed:" + err1.Error(),
				}
				return
			}
		} else {
			// disable load balancing
			if err1 := services.ServiceGroupApp.ModelLoadBalancing.DisableModelLoadBalancing(
				in.TenantId, in.Provider, in.Model, modelruntimeenumtypes.ModelType(in.ModelType),
			); err1 != nil {
				mlog.Error("DisableModelLoadBalancing failed:", err1)
				out.Exp = &pbexceptions.HTTPException{
					Status:  http.StatusBadRequest,
					Message: "DisableModelLoadBalancing failed:" + err1.Error(),
				}
				return
			}
			if in.ConfigFrom != "predefined-model" {
				if err1 := services.ServiceGroupApp.ModelProvide.SaveModelCredentials(in.TenantId, in.Provider, modelruntimeenumtypes.ModelType(in.ModelType), in.Model, credentials); err1 != nil {
					mlog.Error("SaveModelCredentials failed:", err1)
					out.Exp = &pbexceptions.HTTPException{
						Status:  http.StatusBadRequest,
						Message: "SaveModelCredentials failed:" + err1.Error(),
					}
					return
				}

			}
		}
	}()
	return
}

func (s *AdminService) EnableModelProviderModel(ctx context.Context, in *pbapi.EnableModelProviderModelRequest) (out *pbapi.EnableModelProviderModelReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.EnableModelProviderModel call:%#v", p.Addr.String(), in)

	out = &pbapi.EnableModelProviderModelReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if in.Provider == "" {
		mlog.Error("missing provider")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing provider")
		return
	}
	if !modelruntimeenumtypes.ModelType(in.ModelType).Valid() {
		mlog.Errorf("invalid model type=%s", in.ModelType)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("invalid model type=%s", in.ModelType))
		return
	}
	if in.Model == "" {
		mlog.Error("missing model")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing model")
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
		if in.Enable {
			services.ServiceGroupApp.ModelProvide.EnableModel(in.TenantId, in.Provider, in.Model, in.ModelType)
		} else {
			services.ServiceGroupApp.ModelProvide.DisableModel(in.TenantId, in.Provider, in.Model, in.ModelType)
		}
	}()
	return
}
func (s *AdminService) GetModelParameterRules(ctx context.Context, in *pbapi.GetModelParameterRulesRequest) (out *pbapi.GetModelParameterRulesReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetModelParameterRules call:%#v", p.Addr.String(), in)

	out = &pbapi.GetModelParameterRulesReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if in.Provider == "" {
		mlog.Error("missing provider")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing provider")
		return
	}
	if in.Model == "" {
		mlog.Error("missing model")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing model")
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
		params := services.ServiceGroupApp.ModelProvide.GetModelParameterRules(in.TenantId, in.Provider, in.Model)
		bindata, _ := json.Marshal(params)
		out.ParameterRulesStr = string(bindata)
	}()
	return
}
