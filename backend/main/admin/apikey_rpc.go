package main

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"google.golang.org/grpc/peer"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

// func _get_resource(resource_id, tenant_id string){
//     resource = resource_model.query.filter_by(id=resource_id, tenant_id=tenant_id).first()

//     if resource is None:
//         flask_restful.abort(404, message=f"{resource_model.__name__} not found.")
// }
//     return resource
// }

func (s *AdminService) GetApiKeyList(ctx context.Context, in *pbapi.GetApiKeyListRequest) (out *pbapi.GetApiKeyListReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetApiKeyList call:%#v", p.Addr.String(), in)

	out = &pbapi.GetApiKeyListReply{}
	if in.ResourceIdField == "" {
		mlog.Error("resource_id_field must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_id_field must be set")
		return
	}
	if in.ResourceId == "" {
		mlog.Error("resource_id must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_id must be set")
		return
	}
	if in.ResourceType == "" {
		mlog.Error("resource_type must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_type must be set")
		return
	}
	if !slices.Contains([]string{"app", "dataset"}, in.ResourceType) {
		mlog.Errorf("resource_type=%s must be \"app\" or \"dataset\"", in.ResourceType)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("resource_type=%s must be \"app\" or \"dataset\"", in.ResourceType))
		return
	}
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
	switch in.ResourceType {
	case "app":
		resource := new(models.App)
		if err1 := dbengine.Instance().DB.Model(&models.App{}).Where("id= ? and tenant_id= ?", in.ResourceId, current_user.CurrentTenantID()).First(resource).Error; err1 != nil {
			mlog.Errorf("get app failed:%v", err1)
			resource = nil
		}

		if resource == nil {
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusNotFound,
				Message: "App not found.",
			}
			return
		}
	case "dataset":
		resource := new(models.Dataset)
		if err1 := dbengine.Instance().DB.Model(&models.Dataset{}).Where("id= ? and tenant_id= ?", in.ResourceId, current_user.CurrentTenantID()).First(resource).Error; err1 != nil {
			mlog.Errorf("get Dataset failed:%v", err1)
			resource = nil
		}

		if resource == nil {
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusNotFound,
				Message: "Dataset not found.",
			}
			return
		}
	}
	var keys []*models.ApiToken
	sql := fmt.Sprintf("select * from api_tokens where type = '%s' and %s = '%s';", in.ResourceType, in.ResourceIdField, in.ResourceId)
	if err1 := dbengine.Instance().DB.Raw(sql).Scan(&keys); err1 != nil {
		mlog.Errorf("get api_tokens failed:%v", err1)
	}
	for _, key := range keys {
		if out.Items == nil {
			out.Items = make([]*pbapi.ApiKeyResponse, 0)
		}
		item := &pbapi.ApiKeyResponse{
			Id:    key.ID,
			Type:  key.Type,
			Token: key.Token,
		}
		if key.LastUsedAt != nil {
			item.LastUsedAt = key.LastUsedAt.Unix()
		}
		if key.CreatedAt != nil {
			item.CreatedAt = key.CreatedAt.Unix()
		}
		out.Items = append(out.Items, item)
	}
	return
}

func (s *AdminService) GenerateApiKeyList(ctx context.Context, in *pbapi.GenerateApiKeyRequest) (out *pbapi.GenerateApiKeyReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GenerateApiKey call:%#v", p.Addr.String(), in)

	out = &pbapi.GenerateApiKeyReply{}
	if in.ResourceIdField == "" {
		mlog.Error("resource_id_field must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_id_field must be set")
		return
	}
	if in.ResourceId == "" {
		mlog.Error("resource_id must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_id must be set")
		return
	}
	if in.ResourceType == "" {
		mlog.Error("resource_type must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_type must be set")
		return
	}
	if !slices.Contains([]string{"app", "dataset"}, in.ResourceType) {
		mlog.Errorf("resource_type=%s must be \"app\" or \"dataset\"", in.ResourceType)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("resource_type=%s must be \"app\" or \"dataset\"", in.ResourceType))
		return
	}
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
	switch in.ResourceType {
	case "app":
		resource := new(models.App)
		if err1 := dbengine.Instance().DB.Model(&models.App{}).Where("id= ? and tenant_id= ?", in.ResourceId, current_user.CurrentTenantID()).First(resource).Error; err1 != nil {
			mlog.Errorf("get app failed:%v", err1)
			resource = nil
		}

		if resource == nil {
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusNotFound,
				Message: "App not found.",
			}
			return
		}
	case "dataset":
		resource := new(models.Dataset)
		if err1 := dbengine.Instance().DB.Model(&models.Dataset{}).Where("id= ? and tenant_id= ?", in.ResourceId, current_user.CurrentTenantID()).First(resource).Error; err1 != nil {
			mlog.Errorf("get Dataset failed:%v", err1)
			resource = nil
		}

		if resource == nil {
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusNotFound,
				Message: "Dataset not found.",
			}
			return
		}
	}
	var current_key_count int64
	err = dbengine.Instance().DB.Model(&models.ApiToken{}).Where("type = ?", in.ResourceType).Where(in.ResourceIdField+"=?", in.ResourceId).Count(&current_key_count).Error
	if err != nil {
		mlog.Errorf("count api key failed:%v", err)
	}

	if current_key_count >= confy.GetWithDefault[int64]("api_key.max_keys", 10) {
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusBadRequest,
			Code:    "max_keys_exceeded",
			Message: fmt.Sprintf("Cannot create more than %d API keys for this resource type.", confy.GetWithDefault[int64]("api_key.max_keys", 10)),
		}
		return
	}
	now := time.Now()
	key := (models.ApiToken{}).GenerateApiKey(in.TokenPrefix, 24)
	api_token := &models.ApiToken{
		ID:         uuid.NewV4().String(),
		AppID:      in.ResourceId,
		Type:       in.ResourceType,
		Token:      key,
		LastUsedAt: &now,
		CreatedAt:  &now,
		TenantID:   current_user.CurrentTenantID(),
	}
	dbengine.Instance().DB.Create(api_token)
	out.Item = &pbapi.ApiKeyResponse{
		Id:    api_token.ID,
		Type:  api_token.Type,
		Token: api_token.Token,
	}
	if api_token.LastUsedAt != nil {
		out.Item.LastUsedAt = api_token.LastUsedAt.Unix()
	}
	if api_token.CreatedAt != nil {
		out.Item.CreatedAt = api_token.CreatedAt.Unix()
	}
	return
}

func (s *AdminService) DelApiKey(ctx context.Context, in *pbapi.DelApiKeyRequest) (out *pbapi.DelApiKeyReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.DelApiKey call:%#v", p.Addr.String(), in)

	out = &pbapi.DelApiKeyReply{}
	if in.ResourceIdField == "" {
		mlog.Error("resource_id_field must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_id_field must be set")
		return
	}
	if in.ResourceId == "" {
		mlog.Error("resource_id must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_id must be set")
		return
	}
	if in.ApiKeyId == "" {
		mlog.Error("api_key_id must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("api_key_id must be set")
		return
	}
	if in.ResourceType == "" {
		mlog.Error("resource_type must be set")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("resource_type must be set")
		return
	}
	if !slices.Contains([]string{"app", "dataset"}, in.ResourceType) {
		mlog.Errorf("resource_type=%s must be \"app\" or \"dataset\"", in.ResourceType)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("resource_type=%s must be \"app\" or \"dataset\"", in.ResourceType))
		return
	}
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
	if !current_user.IsAdminOrOwner() {
		mlog.Errorf("current user(%#v) is not admin or owner", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not admin or owner", current_user))
		return
	}
	switch in.ResourceType {
	case "app":
		resource := new(models.App)
		if err1 := dbengine.Instance().DB.Model(&models.App{}).Where("id= ? and tenant_id= ?", in.ResourceId, current_user.CurrentTenantID()).First(resource).Error; err1 != nil {
			mlog.Errorf("get app failed:%v", err1)
			resource = nil
		}

		if resource == nil {
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusNotFound,
				Message: "App not found.",
			}
			return
		}
	case "dataset":
		resource := new(models.Dataset)
		if err1 := dbengine.Instance().DB.Model(&models.Dataset{}).Where("id= ? and tenant_id= ?", in.ResourceId, current_user.CurrentTenantID()).First(resource).Error; err1 != nil {
			mlog.Errorf("get Dataset failed:%v", err1)
			resource = nil
		}

		if resource == nil {
			out.Exp = &pbexceptions.HTTPException{
				Status:  http.StatusNotFound,
				Message: "Dataset not found.",
			}
			return
		}
	}
	key := new(models.ApiToken)
	err = dbengine.Instance().DB.Model(&models.ApiToken{}).Where("type = ? and id = ?", in.ResourceType, in.ApiKeyId).Where(in.ResourceIdField+"=?", in.ResourceId).First(key).Error
	if err != nil {
		mlog.Errorf("get api key failed:%v", err)
		key = nil
	}
	if key == nil {
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusNotFound,
			Message: "API key not found",
		}
		return
	}
	dbengine.Instance().DB.Delete(key)
	return
}
