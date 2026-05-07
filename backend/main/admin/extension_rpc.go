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
