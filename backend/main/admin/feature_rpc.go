package main

import (
	"context"

	"google.golang.org/grpc/peer"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
)

func (s *AdminService) GetSystemFeatures(ctx context.Context, in *pbapi.GetSystemFeaturesRequest) (out *pbapi.GetSystemFeaturesReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetSystemFeatures call:%#v", p.Addr.String(), in)
	out = &pbapi.GetSystemFeaturesReply{}
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
		out.Feature = services.ServiceGroupApp.Feature.GetSystemFeatures()
	}()
	return
}
func (s *AdminService) GetFeatures(ctx context.Context, in *pbapi.GetFeaturesRequest) (out *pbapi.GetFeaturesReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetSystemFeatures call:%#v", p.Addr.String(), in)
	out = &pbapi.GetFeaturesReply{}
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
		out.Feature = services.ServiceGroupApp.Feature.GetFeatures(in.TenantId)
	}()
	return
}
