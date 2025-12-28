package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/peer"
	"mlib.com/confy"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	appexceptions "mlib.com/gofy/server/core/exceptions/app"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/core/variables"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	variablefactory "mlib.com/gofy/server/factories/variable_factory"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

func (s *AdminService) GetWorkflowDraft(ctx context.Context, in *pbapi.GetWorkflowDraftRequest) (out *pbapi.GetWorkflowDraftReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDraft call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDraftReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	rsp := func() *response.WorkflowResponse {
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		// # fetch draft workflow by app_model
		wf := services.ServiceGroupApp.Workflow.GetDraftWorkflow(app_model)

		if wf == nil {
			mlog.Error("draft workflow not exist")
			out.Exp = exceptions.NewDraftWorkflowNotExistPbHttpExp("")
			return nil
		}
		return response.NewWorkflowResponse(wf)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.WorkflowResponseStr = string(bindata)
	return
}
func (s *AdminService) WorkflowSyncDraft(ctx context.Context, in *pbapi.WorkflowSyncDraftRequest) (out *pbapi.WorkflowSyncDraftReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.WorkflowSyncDraft call:%#v", p.Addr.String(), in)

	out = &pbapi.WorkflowSyncDraftReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.GraphStr == "" {
		mlog.Error("missing graph")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing graph")
		return
	}
	graph := map[string]any{}
	err1 := json.Unmarshal([]byte(in.GraphStr), &graph)
	if err1 != nil {
		mlog.Errorf("json unmarshal graph_str=%s to dict failed:%v", in.GraphStr, err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal graph_str=%s to dict failed:%v", in.GraphStr, err1))
		return
	}
	if in.FeaturesStr == "" {
		mlog.Error("missing features")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing features")
		return
	}
	features := map[string]any{}
	err1 = json.Unmarshal([]byte(in.FeaturesStr), &features)
	if err1 != nil {
		mlog.Errorf("json unmarshal features_str=%s to dict failed:%v", in.FeaturesStr, err1)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal features_str=%s to dict failed:%v", in.FeaturesStr, err1))
		return
	}
	environment_variables_list := []map[string]any{}
	if in.EnvironmentVariablesStr != "" {
		err1 = json.Unmarshal([]byte(in.EnvironmentVariablesStr), &environment_variables_list)
		if err1 != nil {
			mlog.Errorf("json unmarshal environment_variables_str=%s to dict failed:%v", in.EnvironmentVariablesStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal environment_variables_str=%s to dict failed:%v", in.EnvironmentVariablesStr, err1))
			return
		}
	}
	conversation_variables_list := []map[string]any{}
	if in.ConversationVariablesStr != "" {
		err1 = json.Unmarshal([]byte(in.ConversationVariablesStr), &conversation_variables_list)
		if err1 != nil {
			mlog.Errorf("json unmarshal conversation_variables_str=%s to dict failed:%v", in.ConversationVariablesStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal conversation_variables_str=%s to dict failed:%v", in.ConversationVariablesStr, err1))
			return
		}
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(*appexceptions.WorkflowHashNotEqualError); ok {
					panic(httpexceptions.NewDraftWorkflowNotSync())
				} else if exp, ok := r.(httpexceptions.HTTPException); ok {
					out.Exp = exp.ToPbHttpException(exp)
				} else if exp, ok := r.(error); ok {
					out.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
				} else {
					panic(r)
				}
			}
		}()
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		environment_variables := []variables.Variabler{}
		for _, obj := range environment_variables_list {
			environment_variables = append(environment_variables, variablefactory.BuildEnvironmentVariableFromMapping(obj))
		}
		conversation_variables := []variables.Variabler{}
		for _, obj := range conversation_variables_list {
			conversation_variables = append(conversation_variables, variablefactory.BuildConversationVariableFromMapping(obj))
		}

		wf := services.ServiceGroupApp.Workflow.SyncDraftWorkflow(
			app_model,
			graph,
			features,
			in.Hash,
			current_user,
			environment_variables,
			conversation_variables,
		)
		out.UniqueHash = in.Hash
		if wf.UpdatedAt != nil {
			out.UpdatedAt = wf.UpdatedAt.Unix()
		} else {
			if wf.CreatedAt != nil {
				out.UpdatedAt = wf.CreatedAt.Unix()
			} else {
				out.UpdatedAt = time.Now().Unix()
			}
		}
	}()
	if out.Exp != nil {
		return
	}
	return
}

// func (s *AdminService) AppRun(in *pbapi.AppRunRequest, out grpc.ServerStreamingServer[pbapi.AppRunReply]) error {
// 	mlog.Infof("remote admin.AppRun call:%#v", in)

// 	out_rsp := &pbapi.AppRunReply{}
// 	if in.UserId == "" {
// 		mlog.Errorf("Account not provide")
// 		out_rsp.Exp = exceptions.NewUnauthorizedPbHttpExp("Account not provide")
// 		return out.Send(out_rsp)
// 	}
// 	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
// 	if exp != nil {
// 		mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
// 		out_rsp.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
// 		return out.Send(out_rsp)
// 	}
// 	if !current_user.IsEditor() {
// 		mlog.Errorf("current user(%#v) is not editor", current_user)
// 		out_rsp.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
// 		return out.Send(out_rsp)
// 	}
// 	if in.AppId == "" {
// 		mlog.Error("missing app id")
// 		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
// 		return out.Send(out_rsp)
// 	}
// 	if in.ArgsStr == "" {
// 		mlog.Error("missing args")
// 		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp("missing args")
// 		return out.Send(out_rsp)
// 	}
// 	args := map[string]any{}
// 	err := json.Unmarshal([]byte(in.ArgsStr), &args)
// 	if err != nil {
// 		mlog.Errorf("json unmarshal args_str=%s to dict failed:%v", in.ArgsStr, err)
// 		out_rsp.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal args_str=%s to dict failed:%v", in.ArgsStr, err))
// 		return out.Send(out_rsp)
// 	}
// 	func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				if exp, ok := r.(httpexceptions.HTTPException); ok {
// 					out_rsp.Exp = exp.ToPbHttpException(exp)
// 				} else if exp, ok := r.(error); ok {
// 					out_rsp.Exp = exceptions.NewInternalServerPbHttpExp(exp.Error())
// 				} else {
// 					panic(r)
// 				}
// 			}
// 		}()
// 		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)

// 		rsp, rspiter := services.ServiceGroupApp.AppGenerate.Generate(app_model, current_user, args, appenumtypes.InvokeFrom_DEBUGGER, true)
// 		mlog.Debugf("***********realrsp=%#v", rsp)
// 		if rsp != nil {
// 			bindata, _ := json.Marshal(rsp)

// 			out_rsp.DirectReplyDictStr = string(bindata)
// 			return
// 		}

// 		for item := range rspiter {
// 			mlog.Debugf("-----send:%s", string(item))
// 			out_rsp.StreamReplyStr = item
// 			err = out.Send(out_rsp)
// 			if err != nil {
// 				mlog.Errorf("send stream error:%v", err)
// 				out_rsp = nil
// 				return
// 			}
// 		}
// 	}()
// 	if err != nil {
// 		return err
// 	}
// 	if out_rsp != nil {
// 		return out.Send(out_rsp)
// 	}
// 	return nil
// }

func (s *AdminService) WorkflowNodeRun(ctx context.Context, in *pbapi.WorkflowNodeRunRequest) (out *pbapi.WorkflowNodeRunReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.WorkflowNodeRun call:%#v", p.Addr.String(), in)

	out = &pbapi.WorkflowNodeRunReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.NodeId == "" {
		mlog.Error("missing node id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing node id")
		return
	}
	inputs := map[string]any{}
	if in.InputsDictStr != "" {
		if err1 := json.Unmarshal([]byte(in.InputsDictStr), &inputs); err1 != nil {
			mlog.Errorf("json unmarshal inputs_dict_str=%s to dict failed:%v", in.InputsDictStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("json unmarshal inputs_dict_str=%s to dict failed:%v", in.InputsDictStr, err1))
			return
		}
	}
	rsp := func() *response.WorkflowRunNodeExecutionResponse {
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		workflow_node_execution := services.ServiceGroupApp.Workflow.RunDraftWorkflowNode(app_model, in.NodeId, inputs, current_user)
		return response.NewWorkflowRunNodeExecutionResponse(workflow_node_execution)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.WorkflowRunNodeExecutionResponseStr = string(bindata)
	return
}
func (s *AdminService) GetWorkflowConfig(ctx context.Context, in *pbapi.GetWorkflowConfigRequest) (out *pbapi.GetWorkflowConfigReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowConfig call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowConfigReply{}
	out.ParallelDepthLimit = confy.GetWithDefault[int32]("workflow.parallel_depth_limit", 3)
	return
}
func (s *AdminService) GetWorkflowDefaultBlockConfigs(ctx context.Context, in *pbapi.GetWorkflowDefaultBlockConfigsRequest) (out *pbapi.GetWorkflowDefaultBlockConfigsReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDefaultBlockConfigs call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDefaultBlockConfigsReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	cfgs := func() []map[string]any {
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
		return services.ServiceGroupApp.Workflow.GetDefaultBlockConfigs()
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(cfgs)
	out.CfgsStr = string(bindata)
	return
}
func (s *AdminService) GetWorkflowDefaultBlockConfig(ctx context.Context, in *pbapi.GetWorkflowDefaultBlockConfigRequest) (out *pbapi.GetWorkflowDefaultBlockConfigReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDefaultBlockConfig call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDefaultBlockConfigReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.BlockType == "" {
		mlog.Error("missing block_type")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing block_type")
		return
	}
	if !nodesenumtypes.NodeType(in.BlockType).Valid() {
		mlog.Error("invalid block_type")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("invalid block_type")
		return
	}
	filter_dict := map[string]any{}
	if in.FilterDictStr != "" {
		err1 := json.Unmarshal([]byte(in.FilterDictStr), &filter_dict)
		if err1 != nil {
			mlog.Errorf("json unmarshal=%s to dict failed:%v", in.FilterDictStr, err1)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp("filters must be dict")
			return
		}
	}
	rsp := func() map[string]any {
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
		return services.ServiceGroupApp.Workflow.GetDefaultBlockConfig(nodesenumtypes.NodeType(in.BlockType), filter_dict)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(rsp)
	out.CfgStr = string(bindata)
	return
}
func (s *AdminService) GetWorkflowPublished(ctx context.Context, in *pbapi.GetWorkflowPublishedRequest) (out *pbapi.GetWorkflowPublishedReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowPublished call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowPublishedReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	rsp := func() *response.WorkflowResponse {
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		wf := services.ServiceGroupApp.Workflow.GetPublishedWorkflow(app_model)
		if wf == nil {
			return nil
		}
		return response.NewWorkflowResponse(wf)
	}()
	if out.Exp != nil {
		return
	}
	if rsp != nil {
		bindata, _ := json.Marshal(rsp)
		out.WorkflowResponseStr = string(bindata)
	}
	return
}
func (s *AdminService) WorkflowPublished(ctx context.Context, in *pbapi.WorkflowPublishedRequest) (out *pbapi.WorkflowPublishedReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.WorkflowPublished call:%#v", p.Addr.String(), in)

	out = &pbapi.WorkflowPublishedReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	created_at := func() int64 {
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)
		wf := services.ServiceGroupApp.Workflow.PublishWorkflow(app_model, current_user, nil)
		return wf.CreatedAt.Unix()
	}()
	if out.Exp != nil {
		return
	}
	out.CreatedAt = created_at
	return
}
func (s *AdminService) GetWorkflowDraftVariableList(ctx context.Context, in *pbapi.GetWorkflowDraftVariableListRequest) (out *pbapi.GetWorkflowDraftVariableListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDraftVariableList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDraftVariableListReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}

	// Get current user for authorization
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)

		// Check if workflow exists
		workflowExist := services.ServiceGroupApp.Workflow.IsWorkflowExist(app_model)
		if !workflowExist {
			mlog.Error("draft workflow not exist")
			out.Exp = exceptions.NewDraftWorkflowNotExistPbHttpExp("")
			return
		}

		// Use the new Go service to get draft variables
		workflowVars, total, err := services.ServiceGroupApp.WorkflowDraftVariable.ListVariablesWithoutValues(in.AppId, int(in.Page), int(in.Limit))
		if err != nil {
			mlog.Errorf("failed to get workflow draft variables: %v", err)
			out.Exp = exceptions.NewInternalServerPbHttpExp(fmt.Sprintf("failed to get workflow draft variables: %v", err))
			return
		}

		// Convert to proto response format
		items := make([]*pbapi.PK_WORKFLOW_DRAFT_VARIABLE_WITHOUT_VALUE_FIELDS, 0)
		for _, variable := range workflowVars {
			item := &pbapi.PK_WORKFLOW_DRAFT_VARIABLE_WITHOUT_VALUE_FIELDS{
				Id:          variable.ID,
				Type:        string(variable.GetVariableType()),
				Name:        variable.Name,
				Description: variable.Description,
				Selector:    []string{variable.NodeID, variable.Name},
				ValueType:   variable.ValueType,
				Edited:      variable.LastEditedAt != nil,
				Visible:     variable.Visible,
			}
			items = append(items, item)
		}

		out.Items = items
		out.Total = total
	}()

	if out.Exp != nil {
		return
	}

	return out, nil
}
func (s *AdminService) GetWorkflowDraftVariable(ctx context.Context, in *pbapi.GetWorkflowDraftVariableRequest) (out *pbapi.GetWorkflowDraftVariableReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDraftVariable call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDraftVariableReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	if in.VariableId == "" {
		mlog.Error("missing variable id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing variable id")
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

		variable, err := services.ServiceGroupApp.WorkflowDraftVariable.GetVariable(in.VariableId)

		if err != nil || variable == nil {
			mlog.Error("find variable %s failed:%v", in.VariableId, err)
			panic(exceptions.NewValueError(fmt.Sprintf("variable not found, id={%s}", in.VariableId)))
		}
		if variable.AppID != in.AppId {
			mlog.Error("find variable %s failed:%v", in.VariableId, err)
			panic(exceptions.NewValueError(fmt.Sprintf("variable not found, id={%s}", in.VariableId)))
		}
		var value any
		seg, _ := variable.GetValue()
		if seg != nil {
			value = seg.GetValue()
		}
		item := map[string]any{
			"id":          variable.ID,
			"type":        string(variable.GetVariableType()),
			"name":        variable.Name,
			"description": variable.Description,
			"selector":    []string{variable.NodeID, variable.Name},
			"valueType":   variable.ValueType,
			"edited":      variable.LastEditedAt != nil,
			"visible":     variable.Visible,
			"value":       value,
		}
		bindata, _ := json.Marshal(item)
		out.VarStr = string(bindata)
	}()

	if out.Exp != nil {
		return
	}
	return out, nil
}
func (s *AdminService) GetWorkflowDraftSysVariableList(ctx context.Context, in *pbapi.GetWorkflowDraftVariableListRequest) (out *pbapi.GetWorkflowDraftVariableListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDraftSysVariableList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDraftVariableListReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}

	// Get current user for authorization
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)

		variables, err := services.ServiceGroupApp.WorkflowDraftVariable.GetVariableList(app_model, constants.SYSTEM_VARIABLE_NODE_ID)
		if err != nil {
			panic(err)
		}
		items := make([]map[string]any, 0)
		for _, variable := range variables {
			var value any
			seg, _ := variable.GetValue()
			if seg != nil {
				value = seg.GetValue()
			}
			item := map[string]any{
				"id":          variable.ID,
				"type":        string(variable.GetVariableType()),
				"name":        variable.Name,
				"description": variable.Description,
				"selector":    []string{variable.NodeID, variable.Name},
				"valueType":   variable.ValueType,
				"edited":      variable.LastEditedAt != nil,
				"visible":     variable.Visible,
				"value":       value,
			}
			items = append(items, item)
		}
		bindata, _ := json.Marshal(items)
		out.ItemsStr = string(bindata)
	}()

	if out.Exp != nil {
		return
	}

	return out, nil
}

func (s *AdminService) GetWorkflowDraftConversationVariableList(ctx context.Context, in *pbapi.GetWorkflowDraftVariableListRequest) (out *pbapi.GetWorkflowDraftVariableListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDraftConversationVariableList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDraftVariableListReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}

	// Get current user for authorization
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)

		draftWf := services.ServiceGroupApp.Workflow.GetDraftWorkflow(app_model)
		if draftWf == nil {
			mlog.Errorf("draft workflow not found, id={%s}", app_model.ID)
			panic(exceptions.NewValueError(fmt.Sprintf("draft workflow not found, id={%s}", app_model.ID)))
		}
		err := services.ServiceGroupApp.WorkflowDraftVariable.PrefillConversationVariableDefaultValues(draftWf)
		if err != nil {
			panic(err)
		}
		variables, err := services.ServiceGroupApp.WorkflowDraftVariable.GetVariableList(app_model, constants.CONVERSATION_VARIABLE_NODE_ID)
		if err != nil {
			panic(err)
		}
		items := make([]map[string]any, 0)
		for _, variable := range variables {
			var value any
			seg, _ := variable.GetValue()
			if seg != nil {
				value = seg.GetValue()
			}
			item := map[string]any{
				"id":          variable.ID,
				"type":        string(variable.GetVariableType()),
				"name":        variable.Name,
				"description": variable.Description,
				"selector":    []string{variable.NodeID, variable.Name},
				"valueType":   variable.ValueType,
				"edited":      variable.LastEditedAt != nil,
				"visible":     variable.Visible,
				"value":       value,
			}
			items = append(items, item)
		}
		bindata, _ := json.Marshal(items)
		out.ItemsStr = string(bindata)
	}()

	if out.Exp != nil {
		return
	}

	return out, nil
}
func (s *AdminService) GetWorkflowDraftEnvVariableList(ctx context.Context, in *pbapi.GetWorkflowDraftVariableListRequest) (out *pbapi.GetWorkflowDraftVariableListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetWorkflowDraftEnvVariableList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetWorkflowDraftVariableListReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}

	// Get current user for authorization
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.UserId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.UserId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.UserId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
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
		app_model := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)

		draftWf := services.ServiceGroupApp.Workflow.GetDraftWorkflow(app_model)
		if draftWf == nil {
			mlog.Errorf("draft workflow not found, id={%s}", app_model.ID)
			panic(exceptions.NewValueError(fmt.Sprintf("draft workflow not found, id={%s}", app_model.ID)))
		}
		wf := services.ServiceGroupApp.Workflow.GetDraftWorkflow(app_model)
		if wf == nil {
			panic(httpexceptions.NewDraftWorkflowNotExist())
		}
		env_vars := wf.GetEnvironmentVariables()
		env_vars_list := []map[string]any{}
		for _, v := range env_vars {
			env_vars_list = append(env_vars_list, map[string]any{
				"id":          v.GetID(),
				"type":        "env",
				"name":        v.GetName(),
				"description": v.GetDescription(),
				"selector":    v.GetSelector(),
				"value_type":  v.ValueType(),
				"value":       v.GetValue(),
				// Do not track edited for env vars.
				"edited":   false,
				"visible":  true,
				"editable": true,
			})
		}

		bindata, _ := json.Marshal(env_vars_list)
		out.ItemsStr = string(bindata)
	}()

	if out.Exp != nil {
		return
	}

	return out, nil
}
