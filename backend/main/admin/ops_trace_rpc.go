package main

import (
	"context"
	"encoding/json"

	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
)

func (s *AdminService) GetTraceAppConfig(ctx context.Context, in *pbapi.GetTraceAppConfigRequest) (out *pbapi.GetTraceAppConfigReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetTraceAppConfig call:%#v", p.Addr.String(), in)

	out = &pbapi.GetTraceAppConfigReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.TracingProvider == "" {
		mlog.Error("missing tracing provider")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing tracing provider")
		return
	}
	trace_config := func() map[string]any {
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
		return services.ServiceGroupApp.Ops.GetTraceAppConfig(in.AppId, in.TracingProvider)
	}()
	if out.Exp != nil {
		return
	}
	if len(trace_config) > 0 {
		bindata, _ := json.Marshal(trace_config)
		out.TraceConfigStr = string(bindata)
	}
	return
}
