package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	pbexceptions "mlib.com/gofy/server/proto/exceptions"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

func (s *AdminService) GetModelProviderList(ctx context.Context, in *pbapi.GetModelProviderListRequest) (out *pbapi.GetModelProviderListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetModelProviderList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetModelProviderListReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if in.ModelType != "" {
		if !modelruntimeentities.ModelType(in.ModelType).Valid() {
			mlog.Errorf("model_type=%s is invalid", in.ModelType)
			out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("model_type=%s is invalid", in.ModelType))
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
		provider_list := services.ServiceGroupApp.ModelProvide.GetProviderList(in.TenantId, modelruntimeentities.ModelType(in.ModelType))
		bindata, _ := json.Marshal(provider_list)
		out.ProviderListStr = string(bindata)
	}()

	return
}
func (s *AdminService) UpdateModelProvider(ctx context.Context, in *pbapi.UpdateModelProviderRequest) (out *pbapi.UpdateModelProviderReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.UpdateModelProvider call:%#v", p.Addr.String(), in)

	out = &pbapi.UpdateModelProviderReply{}
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
	if in.CredentialsDictStr == "" {
		mlog.Error("missing credentials")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing credentials")
		return
	}
	var credentials map[string]any
	if err1 := json.Unmarshal([]byte(in.CredentialsDictStr), &credentials); err1 != nil {
		mlog.Errorf("json unmarshal credentials_dict_str=%s to dict failed:%v", in.CredentialsDictStr, err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal credentials_dict_str=%s to dict failed:%v", in.CredentialsDictStr, err1))
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
		if err1 := services.ServiceGroupApp.ModelProvide.SaveProviderCredentials(in.TenantId, in.Provider, credentials); err1 != nil {
			mlog.Error("SaveProviderCredentials failed:", err1)
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusBadRequest,
				Message: "SaveProviderCredentials failed:" + err1.Error(),
			}
			return
		}
	}()
	return
}
func (s *AdminService) DelModelProvider(ctx context.Context, in *pbapi.DelModelProviderRequest) (out *pbapi.DelModelProviderReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.DelModelProvider call:%#v", p.Addr.String(), in)

	out = &pbapi.DelModelProviderReply{}
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
		if err1 := services.ServiceGroupApp.ModelProvide.RemoveProviderCredentials(in.TenantId, in.Provider); err1 != nil {
			mlog.Error("RemoveProviderCredentials failed:", err1)
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusBadRequest,
				Message: "RemoveProviderCredentials failed:" + err1.Error(),
			}
			return
		}
	}()
	return
}

func (s *AdminService) GetModelProviderCredentials(ctx context.Context, in *pbapi.GetModelProviderCredentialsRequest) (out *pbapi.GetModelProviderCredentialsReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetModelProviderCredentials call:%#v", p.Addr.String(), in)

	out = &pbapi.GetModelProviderCredentialsReply{}
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
		credentials := services.ServiceGroupApp.ModelProvide.GetProviderCredentials(in.TenantId, in.Provider)
		bindata, _ := json.Marshal(credentials)
		out.CredentialsDictStr = string(bindata)
	}()

	return
}
