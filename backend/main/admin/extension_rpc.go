package main

import (
	"context"

	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
)

func (s *AdminService) GetCodeBasedExtension(ctx context.Context, in *pbapi.GetCodeBasedExtensionRequest) (out *pbapi.GetCodeBasedExtensionReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetCodeBasedExtension call:%#v", p.Addr.String(), in)

	out = &pbapi.GetCodeBasedExtensionReply{Module: in.Module}
	if in.Module == "" {
		mlog.Errorf("missing module")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing module")
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
		out.Data = services.ServiceGroupApp.Extension.GetCodeBasedExtension(in.Module)
	}()

	return
}
