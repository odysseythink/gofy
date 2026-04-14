package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/global"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
)

func (s *AdminService) GetAccountProfile(ctx context.Context, in *pbapi.GetAccountProfileRequest) (out *pbapi.GetAccountProfileReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetAccountProfile call:%#v", p.Addr.String(), in)

	out = &pbapi.GetAccountProfileReply{}
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
	now := time.Now()
	out.Id = current_user.ID
	out.Name = current_user.Name
	out.Avatar = current_user.Avatar
	out.Email = current_user.Email
	out.IsPasswordSet = current_user.IsPasswordSet()
	out.InterfaceLanguage = current_user.InterfaceTheme
	out.Timezone = current_user.Timezone
	if current_user.LastLoginAt != nil {
		out.LastLoginAt = current_user.LastLoginAt.Unix()
	} else {
		out.LastLoginAt = now.Unix()
	}

	out.LastLoginIp = current_user.LastLoginIP
	if current_user.CreatedAt != nil {
		out.CreatedAt = current_user.CreatedAt.Unix()
	} else {
		out.CreatedAt = now.Unix()
	}

	return
}
func (s *AdminService) UpdateAccount(ctx context.Context, in *pbapi.UpdateAccountRequest) (out *pbapi.UpdateAccountReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.UpdateAccount call:%#v", p.Addr.String(), in)

	out = &pbapi.UpdateAccountReply{}
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
	if len(in.Infos) == 0 {
		mlog.Error("missing update field")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing update field")
		return
	}
	for k, v := range in.Infos {
		if k == "name" {
			if len(v) < 3 || len(v) > 30 {
				mlog.Error("Account name must be between 3 and 30 characters.")
				out.Exp = exceptions.NewInvalidArgsPbHttpExp("Account name must be between 3 and 30 characters.")
				return
			}
			current_user.Name = v
		} else if k == "avatar" {
			if v == "" {
				mlog.Error("missing avatar")
				out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing avatar")
				return
			}
			current_user.Avatar = v
		} else if k == "interface_language" {
			if !slices.Contains(global.LANGUAGES, v) {
				mlog.Error("invalid interface_language=", v)
				out.Exp = exceptions.NewInvalidArgsPbHttpExp("invalid interface_language=" + v)
				return
			}
			current_user.InterfaceLanguage = v
		} else if k == "interface_theme" {
			if !slices.Contains([]string{"light", "dark"}, v) {
				mlog.Error("invalid interface_theme=", v)
				out.Exp = exceptions.NewInvalidArgsPbHttpExp("invalid interface_theme=" + v)
				return
			}
			current_user.InterfaceTheme = v
		} else if k == "timezone" {
			if !slices.Contains(global.ALL_TIMEZONES, v) {
				mlog.Error("invalid timezone=", v)
				out.Exp = exceptions.NewInvalidArgsPbHttpExp("invalid timezone=" + v)
				return
			}
			current_user.Timezone = v
		} else {
			mlog.Error("unsupported account field=", k)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp("unsupported account field=" + k)
			return
		}
	}

	db := dbengine.Instance().DB.Model(&models.Account{})
	for k, v := range in.Infos {
		db = db.Update(k, v)
	}
	err1 := db.Where("id = ?", current_user.ID).Error
	if err1 != nil {
		mlog.Errorf("update account failed:%v", err1)
		out.Exp = exceptions.NewAccountUpdatePbHttpExp(fmt.Sprintf("update account failed:%v", err1))
		return
	}

	bindata, _ := json.Marshal(response.NewAccountResponse(current_user))
	out.AccountResponseStr = string(bindata)
	return
}
