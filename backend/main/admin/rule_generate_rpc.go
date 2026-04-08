package main

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	llmgenerator "mlib.com/gofy/server/core/llm_generator"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

func (s *AdminService) RuleGenerate(ctx context.Context, in *pbapi.RuleGenerateRequest) (out *pbapi.RuleGenerateReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.RuleGenerate call:%#v", p.Addr.String(), in)
	out = &pbapi.RuleGenerateReply{}
	if in.TenantId == "" {
		mlog.Error("TenantId not provide")
		out.Exp = exceptions.NewUnauthorizedPbHttpExp("TenantId not provide")
		return
	}
	if in.ModelConfigStr == "" {
		mlog.Error("model_config not provide")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("model_config not provide")
		return
	}
	model_config := map[string]any{}

	if err1 := json.Unmarshal([]byte(in.ModelConfigStr), &model_config); err1 != nil {
		mlog.Errorf("json unmarshal %s to dict failed:%v", in.ModelConfigStr, err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("model_config_str must be dict")
		return
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				if exp, ok := r.(*exceptions.ProviderTokenNotInitError); ok {
					httpexp := httpexceptions.NewProviderNotInitializeError(exp.Error())
					out.Exp = httpexp.ToPbHttpException(httpexp)
				} else if _, ok := r.(*exceptions.QuotaExceededError); ok {
					httpexp := httpexceptions.NewProviderQuotaExceededError()
					out.Exp = httpexp.ToPbHttpException(httpexp)
				} else if _, ok := r.(*exceptions.ModelCurrentlyNotSupportError); ok {
					httpexp := httpexceptions.NewProviderModelCurrentlyNotSupportError()
					out.Exp = httpexp.ToPbHttpException(httpexp)
				} else if exp, ok := r.(*modelruntimeexceptions.InvokeError); ok {
					httpexp := httpexceptions.NewCompletionRequestError(exp.Description())
					out.Exp = httpexp.ToPbHttpException(httpexp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()

		rules := (&llmgenerator.LLMGenerator{}).GenerateRuleConfig(in.TenantId, in.Instruction, model_config, in.NoVariable)
		bindata, _ := json.Marshal(rules)
		out.RulesDictStr = string(bindata)
	}()
	return
}
