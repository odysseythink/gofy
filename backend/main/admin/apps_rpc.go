package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/peer"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/core/ops"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	pbentities "mlib.com/gofy/server/proto/entities"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

func (s *AdminService) ListApps(ctx context.Context, in *pbapi.ListAppsRequest) (out *pbapi.ListAppsReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.ListApps call:%#v", p.Addr.String(), in)

	out = &pbapi.ListAppsReply{}
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
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Limit <= 0 {
		in.Page = 20
	}
	if in.Mode == "" {
		in.Mode = "all"
	}
	if in.Page < 1 || in.Page > 99999 {
		mlog.Errorf("page[%d] must be in range[1, 99999]", in.Page)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("page[%d] must be in range[1, 99999]", in.Page))
		return
	}
	if in.Limit < 1 || in.Limit > 100 {
		mlog.Errorf("limit[%d] must be in range[1, 100]", in.Limit)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("limit[%d] must be in range[1, 100]", in.Limit))
		return
	}
	if !slices.Contains([]string{"chat", "workflow", "agent-chat", "channel", "all"}, in.Mode) {
		mlog.Errorf("mode[%s] must be in set[\"chat\", \"workflow\", \"agent-chat\", \"channel\", \"all\"]", in.Mode)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("mode[%s] must be in set[\"chat\", \"workflow\", \"agent-chat\", \"channel\", \"all\"]", in.Mode))
		return
	}
	for idx, v := range in.TagIds {
		if _, err1 := uuid.FromString(v); err1 != nil {
			mlog.Errorf("TagIds[%d]=%s must be uuid format", idx, v)
			out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("TagIds[%d]=%s must be uuid format", idx, v))
			return
		}
	}
	if info := services.ServiceGroupApp.App.GetPaginateApps(current_user.ID, current_user.CurrentTenantID(), in); info == nil {
		mlog.Error("查询失败!")
		// return
		// c.JSON(http.StatusNoContent, gin.H{"data": []any{}, "total": 0, "page": in.Page, "limit": in.Limit, "has_more": false, "result": "failed", "code": 0})
	} else {
		bindata, _ := json.Marshal(info)
		out.PaginateAppsStr = string(bindata)
		// c.JSON(http.StatusOK, info)
	}
	return
}
func (s *AdminService) FindApp(ctx context.Context, in *pbapi.FindAppRequest) (out *pbapi.FindAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.FindApp call:%#v", p.Addr.String(), in)

	out = &pbapi.FindAppReply{}
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
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	app := func() *models.App {
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
		app := services.ServiceGroupApp.App.GetAppModel(in.AppId, current_user, nil)

		app = services.ServiceGroupApp.App.GetApp(in.AppId, app)

		return app
	}()
	if out.Exp != nil {
		return
	}
	rsp := response.NewAppDetailWithSiteResponse(app)
	bindata, _ := json.Marshal(rsp)
	out.AppDetailWithSiteStr = string(bindata)

	return
}
func (s *AdminService) UpdateApp(ctx context.Context, in *pbapi.UpdateAppRequest) (out *pbapi.UpdateAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.UpdateApp call:%#v", p.Addr.String(), in)

	out = &pbapi.UpdateAppReply{}
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
	if in.Name == "" {
		mlog.Errorf("missing name")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing name")
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
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
		dbengine.Instance().DB.Updates(&models.App{
			ID:                  in.AppId,
			Name:                in.Name,
			Description:         in.Description,
			IconType:            in.IconType,
			Icon:                in.Icon,
			IconBackground:      in.IconBackground,
			MaxActiveRequests:   int(in.MaxActiveRequests),
			UseIconAsAnswerIcon: in.UseIconAsAnswerIcon,
		})
		app_model.Name = in.Name
		app_model.Description = in.Description
		app_model.IconType = in.IconType
		app_model.Icon = in.Icon
		app_model.IconBackground = in.IconBackground
		app_model.MaxActiveRequests = int(in.MaxActiveRequests)
		app_model.UseIconAsAnswerIcon = in.UseIconAsAnswerIcon
		rsp := response.NewAppDetailWithSiteResponse(app_model)
		bindata, _ := json.Marshal(rsp)
		out.AppDetailWithSiteStr = string(bindata)
	}()

	return
}
func (s *AdminService) DeleteApp(ctx context.Context, in *pbapi.DeleteAppRequest) (out *pbapi.DeleteAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.DeleteApp call:%#v", p.Addr.String(), in)

	out = &pbapi.DeleteAppReply{}
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
	if err := services.ServiceGroupApp.App.DeleteByIds([]string{in.AppId}); err != nil {
		mlog.Error("delete failed:", err)
		out.Result = "failed"
		// c.JSON(http.StatusNoContent, map[string]any{"result": "failed"})
	} else {
		out.Result = "success"
		// c.JSON(http.StatusNoContent, map[string]any{"result": "success"})
	}
	return
}

func (s *AdminService) CreateApp(ctx context.Context, in *pbapi.CreateAppRequest) (out *pbapi.CreateAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.CreateApp call:%#v", p.Addr.String(), in)

	out = &pbapi.CreateAppReply{}
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
	if in.Name == "" {
		mlog.Errorf("missing name")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing name")
		return
	}
	if !slices.Contains(constants.ALLOW_CREATE_APP_MODES, in.Mode) {
		mlog.Errorf("app mode[%s] must be in set %v", in.Mode, constants.ALLOW_CREATE_APP_MODES)
		out.Exp = exceptions.NewInvalidArgsPbHttpExp(fmt.Sprintf("app mode[%s] must be in set %v", in.Mode, constants.ALLOW_CREATE_APP_MODES))
		return
	}
	if !current_user.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", current_user)
		out.Exp = exceptions.NewForbiddenPbHttpExp(fmt.Sprintf("current user(%#v) is not editor", current_user))
		return
	}
	app := services.ServiceGroupApp.App.CreateApp(current_user.CurrentTenantID(), in, current_user)
	app_detail := response.NewAppDetailResponse(app)
	bindata, _ := json.Marshal(app_detail)
	out.AppDetailStr = string(bindata)
	return
}
func (s *AdminService) ImportApp(ctx context.Context, in *pbapi.ImportAppRequest) (out *pbapi.ImportAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.ImportApp call:%#v", p.Addr.String(), in)

	out = &pbapi.ImportAppReply{}
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
	if in.Mode == "" {
		mlog.Errorf("missing mode")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing mode")
		return
	}
	out.Import = func() *pbentities.Import {
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
		return services.ServiceGroupApp.AppDSL.ImportApp(
			current_user,
			services.ImportModeType(in.Mode),
			in.YamlContent,
			in.YamlUrl,
			in.Name,
			in.Description,
			in.IconType,
			in.Icon,
			in.IconBackground,
			in.AppId,
		)
	}()
	if out.Exp != nil {
		return
	}
	// Return appropriate status code based on result
	status := out.Import.Status
	if status == string(services.ImportStatus_FAILED) {
		out.HttpStatus = int32(http.StatusBadRequest)
	} else if status == string(services.ImportStatus_PENDING) {
		out.HttpStatus = int32(http.StatusAccepted)
	} else {
		out.HttpStatus = int32(http.StatusAccepted)
		// c.JSON(http.StatusOK, result)
	}
	return
}

func (s *AdminService) SetAppName(ctx context.Context, in *pbapi.SetAppNameRequest) (out *pbapi.SetAppNameReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.SetAppName call:%#v", p.Addr.String(), in)

	out = &pbapi.SetAppNameReply{}
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
	if in.Name == "" {
		mlog.Errorf("missing name")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing name")
		return
	}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	app := func() *models.App {
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
		app_model = services.ServiceGroupApp.App.UpdateAppName(app_model, in.Name, current_user)
		return app_model
	}()
	if out.Exp != nil {
		return
	}
	rsp := response.NewAppDetailResponse(app)
	bindata, _ := json.Marshal(rsp)
	out.AppDetailStr = string(bindata)
	return
}
func (s *AdminService) SetAppIcon(ctx context.Context, in *pbapi.SetAppIconRequest) (out *pbapi.SetAppIconReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.SetAppIcon call:%#v", p.Addr.String(), in)

	out = &pbapi.SetAppIconReply{}
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
	app := func() *models.App {
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
		app_model = services.ServiceGroupApp.App.UpdateAppIcon(app_model, in.Icon, in.IconBackground, current_user)
		return app_model
	}()
	if out.Exp != nil {
		return
	}
	rsp := response.NewAppDetailResponse(app)
	bindata, _ := json.Marshal(rsp)
	out.AppDetailStr = string(bindata)
	return
}
func (s *AdminService) CopyApp(ctx context.Context, in *pbapi.CopyAppRequest) (out *pbapi.CopyAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.CopyApp call:%#v", p.Addr.String(), in)

	out = &pbapi.CopyAppReply{}
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
	app := func() *models.App {
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
		yaml_content := services.ServiceGroupApp.AppDSL.ExportDSL(app_model, true)
		result := services.ServiceGroupApp.AppDSL.ImportApp(
			current_user,
			services.ImportMode_YAML_CONTENT,
			yaml_content,
			"",
			in.Name,
			in.Description,
			in.IconType,
			in.Icon,
			in.IconBackground,
			"",
		)
		app := new(models.App)
		err = dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", result.AppId).First(app).Error
		if err != nil {
			mlog.Errorf("get App from mysql failed:%v", err)
			app = nil
			panic(exceptions.NewValueError("get app failed"))
		}
		return app
	}()
	if out.Exp != nil {
		return
	}

	rsp := response.NewAppDetailWithSiteResponse(app)
	bindata, _ := json.Marshal(rsp)
	out.AppDetailWithSiteStr = string(bindata)
	return
}
func (s *AdminService) ExportApp(ctx context.Context, in *pbapi.ExportAppRequest) (out *pbapi.ExportAppReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.ExportApp call:%#v", p.Addr.String(), in)

	out = &pbapi.ExportAppReply{}
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
	out.Dsl = func() string {
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
		return services.ServiceGroupApp.AppDSL.ExportDSL(app_model, in.IncludeSecret)
	}()
	if out.Exp != nil {
		return
	}

	return
}
func (s *AdminService) AppUpdateSiteStatus(ctx context.Context, in *pbapi.AppUpdateSiteStatusRequest) (out *pbapi.AppUpdateSiteStatusReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AppUpdateSiteStatus call:%#v", p.Addr.String(), in)

	out = &pbapi.AppUpdateSiteStatusReply{}
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
	app := func() *models.App {
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
		return services.ServiceGroupApp.App.UpdateSiteStatus(app_model, in.EnableSite, current_user)
	}()
	if out.Exp != nil {
		return
	}
	rsp := response.NewAppDetailResponse(app)
	bindata, _ := json.Marshal(rsp)
	out.AppDetailStr = string(bindata)
	return
}
func (s *AdminService) AppUpdateApiStatus(ctx context.Context, in *pbapi.AppUpdateApiStatusRequest) (out *pbapi.AppUpdateApiStatusReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AppUpdateApiStatus call:%#v", p.Addr.String(), in)

	out = &pbapi.AppUpdateApiStatusReply{}
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
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	app := func() *models.App {
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
		return services.ServiceGroupApp.App.UpdateSiteStatus(app_model, in.EnableApi, current_user)
	}()
	if out.Exp != nil {
		return
	}
	rsp := response.NewAppDetailResponse(app)
	bindata, _ := json.Marshal(rsp)
	out.AppDetailStr = string(bindata)
	return
}
func (s *AdminService) AppGetTrace(ctx context.Context, in *pbapi.AppGetTraceRequest) (out *pbapi.AppGetTraceReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AppGetTrace call:%#v", p.Addr.String(), in)

	out = &pbapi.AppGetTraceReply{}
	if in.AppId == "" {
		mlog.Error("missing app id")
		out.Exp = exceptions.NewInvalidArgsPbHttpExp("missing app id")
		return
	}
	app_trace_config := func() map[string]any {
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
		return (&ops.OpsTraceManager{}).GetAppTracingConfig(in.AppId)
	}()
	if out.Exp != nil {
		return
	}
	bindata, _ := json.Marshal(app_trace_config)
	out.AppTraceConfigStr = string(bindata)
	return
}
func (s *AdminService) AppSetTrace(ctx context.Context, in *pbapi.AppSetTraceRequest) (out *pbapi.AppSetTraceReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.AppSetTrace call:%#v", p.Addr.String(), in)

	out = &pbapi.AppSetTraceReply{}
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
		(&ops.OpsTraceManager{}).UpdateAppTracingConfig(in.AppId, in.Enabled, in.TracingProvider)
	}()

	return
}
