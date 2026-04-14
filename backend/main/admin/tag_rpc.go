package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
)

func (s *AdminService) GetTagList(ctx context.Context, in *pbapi.GetTagListRequest) (out *pbapi.GetTagListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetTagList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetTagListReply{}
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
		out.Tags = services.ServiceGroupApp.Tag.GetTags(in.Type, in.TenantId, in.Keyword)
	}()

	return
}
func (s *AdminService) AddTag(ctx context.Context, in *pbapi.AddTagRequest) (out *pbapi.AddTagReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AddTag call:%#v", p.Addr.String(), in)

	out = &pbapi.AddTagReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	if len(in.Name) < 1 || len(in.Name) > 50 {
		mlog.Errorf("Name must be between 1 to 50 characters.")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("Name must be between 1 to 50 characters.")
		return
	}
	if !slices.Contains(models.TAG_TYPE_LIST, in.Type) {
		mlog.Errorf("Invalid tag type=%s.", in.Type)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("Invalid tag type=%s.", in.Type))
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !(current_user.IsEditor() || current_user.IsDatasetEditor()) {
		mlog.Errorf("current user(%#v) must be editor or dataset editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) must be editor or dataset editor", current_user))
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

		tag := services.ServiceGroupApp.Tag.SaveTag(current_user, in.Name, in.Type)
		out.Id = tag.ID
		out.Name = tag.Name
		out.Type = tag.Type
		out.BindingCount = 0
	}()

	return
}
func (s *AdminService) UpdateTag(ctx context.Context, in *pbapi.UpdateTagRequest) (out *pbapi.UpdateTagReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.UpdateTag call:%#v", p.Addr.String(), in)

	out = &pbapi.UpdateTagReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	if len(in.Name) < 1 || len(in.Name) > 50 {
		mlog.Errorf("Name must be between 1 to 50 characters.")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("Name must be between 1 to 50 characters.")
		return
	}
	if in.TagId == "" {
		mlog.Errorf("missing tag_id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing tag_id")
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !(current_user.IsEditor() || current_user.IsDatasetEditor()) {
		mlog.Errorf("current user(%#v) must be editor or dataset editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) must be editor or dataset editor", current_user))
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
		tag := services.ServiceGroupApp.Tag.UpdateTag(in.Name, in.TagId)

		binding_count := services.ServiceGroupApp.Tag.GetTagBindingCount(in.TagId)
		out.Id = tag.ID
		out.Name = tag.Name
		out.Type = tag.Type
		out.BindingCount = binding_count
	}()

	return
}
func (s *AdminService) DelTag(ctx context.Context, in *pbapi.DelTagRequest) (out *pbapi.DelTagReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.DelTag call:%#v", p.Addr.String(), in)

	out = &pbapi.DelTagReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	if in.TagId == "" {
		mlog.Errorf("missing tag_id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing tag_id")
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) must be editor or dataset editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) must be editor or dataset editor", current_user))
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
		services.ServiceGroupApp.Tag.DeleteByIds([]string{in.TagId})
	}()

	return
}
func (s *AdminService) AddTagBinding(ctx context.Context, in *pbapi.AddTagBindingRequest) (out *pbapi.AddTagBindingReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AddTagBinding call:%#v", p.Addr.String(), in)

	out = &pbapi.AddTagBindingReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	if !slices.Contains(models.TAG_TYPE_LIST, in.Type) {
		mlog.Errorf("Invalid tag type=%s.", in.Type)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("Invalid tag type=%s.", in.Type))
		return
	}
	if in.TargetId == "" {
		mlog.Errorf("missing target_id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing target_id")
		return
	}
	if len(in.TagIds) < 1 {
		mlog.Errorf("Tag IDs is required.")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("Tag IDs is required.")
		return
	}
	for _, v := range in.TagIds {
		if v == "" {
			mlog.Errorf("item of Tag IDs=%#v can't be empty", in.TagIds)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp("item of Tag IDs can't be empty")
			return
		}
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !(current_user.IsEditor() || current_user.IsDatasetEditor()) {
		mlog.Errorf("current user(%#v) must be editor or dataset editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) must be editor or dataset editor", current_user))
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
		services.ServiceGroupApp.Tag.AddTagBinding(current_user, in.Type, in.TargetId, in.TagIds)
	}()

	return
}
func (s *AdminService) DelTagBinding(ctx context.Context, in *pbapi.DelTagBindingRequest) (out *pbapi.DelTagBindingReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.DelTagBinding call:%#v", p.Addr.String(), in)

	out = &pbapi.DelTagBindingReply{}
	if in.AccountId == "" {
		mlog.Errorf("Account[%s] not provide", in.AccountId)
		out.Exp = exceptions.NewUnauthorizedPbHttpExp(fmt.Sprintf("Account[%s] not provide", in.AccountId))
		return
	}
	if !slices.Contains(models.TAG_TYPE_LIST, in.Type) {
		mlog.Errorf("Invalid tag type=%s.", in.Type)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("Invalid tag type=%s.", in.Type))
		return
	}
	if in.TargetId == "" {
		mlog.Errorf("missing target_id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing target_id")
		return
	}
	if in.TagId == "" {
		mlog.Errorf("missing tag_id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing tag_id")
		return
	}
	current_user, exp := services.ServiceGroupApp.Account.LoadLoggedInAccount(in.AccountId)
	if exp != nil {
		mlog.Errorf("load user(%s) failed:%v", in.AccountId, exp.Error())
		out.Exp = exceptions.NewAccountNotInitializedPbHttpExp(fmt.Sprintf("load user(%s) failed:%v", in.AccountId, exp.Error()))
		return
	}
	if !(current_user.IsEditor() || current_user.IsDatasetEditor()) {
		mlog.Errorf("current user(%#v) must be editor or dataset editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) must be editor or dataset editor", current_user))
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
		services.ServiceGroupApp.Tag.DeleteTagBinding(current_user, in.TagId, in.Type, in.TargetId)
	}()

	return
}
