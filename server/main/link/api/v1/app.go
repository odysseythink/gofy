package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

// var AppApiApp = new(AppApi)

type AppApi struct {
}

// Delete 删除App
// @Tags App
// @Summary 删除App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body models.App true "删除App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /console/api/apps/:id [delete]
func (api *AppApi) Delete(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		c.JSON(http.StatusNoContent, map[string]any{"result": "invalid id"})
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).DeleteApp(c, &pbapi.DeleteAppRequest{AccountId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call DeleteApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DeleteApp failed",
			})
			return
		} else {
			mlog.Infof("remote call DeleteApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusNoContent, map[string]any{"result": pbrsp.Result})
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

// Update 更新App
// @Tags App
// @Summary 更新App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body models.App true "更新App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /console/api/apps/:id [put]
func (api *AppApi) Update(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	in := pbapi.UpdateAppRequest{AccountId: acc.ID, AppId: app_id}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		mlog.Error("bind failed:", err)
		response.InvalidArgError(c)
		return
	}
	if in.Name == "" {
		mlog.Error("missing name arg")
		response.InvalidArgErrorWithDetail(c, "missing name arg")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).UpdateApp(c, &in)
		if err != nil {
			mlog.Errorf("remote call UpdateApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call UpdateApp failed",
			})
			return
		} else {
			mlog.Infof("remote call UpdateApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailWithSiteStr == "" {
					mlog.Error("查询失败!")
					response.AppNotFoundError(c)
				} else {
					c.JSON(http.StatusOK, response.NewAppDetailWithSiteResponse(pbrsp.AppDetailWithSiteStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

// Find 用id查询App
// @Tags App
// @Summary 用id查询App
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query models.App true "用id查询App"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /console/api/apps/:id [get]
func (api *AppApi) Find(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("missing app_id")
		response.InvalidArgErrorWithDetail(c, "missing app_id")
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).FindApp(c, &pbapi.FindAppRequest{
			AppId:     app_id,
			AccountId: acc.ID,
		})
		if err != nil {
			mlog.Errorf("remote call FindApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call FindApp failed",
			})
			return
		} else {
			mlog.Infof("remote call FindApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailWithSiteStr == "" {
					mlog.Error("查询失败!")
					response.AppNotFoundError(c)
				} else {
					c.JSON(http.StatusOK, response.NewAppDetailWithSiteResponse(pbrsp.AppDetailWithSiteStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AppApi) GetList(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	in := &pbapi.ListAppsRequest{
		Mode:      "all",
		Page:      1,
		Limit:     20,
		AccountId: acc.ID,
	}
	if err := c.ShouldBindQuery(in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if in.Page < 1 || in.Page > 99999 {
		mlog.Errorf("page[%d] must be in range[1, 99999]", in.Page)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("page[%d] must be in range[1, 99999]", in.Page))
		return
	}
	if in.Limit < 1 || in.Limit > 100 {
		mlog.Errorf("limit[%d] must be in range[1, 100]", in.Limit)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("limit[%d] must be in range[1, 100]", in.Limit))
		return
	}
	if !slices.Contains([]string{"chat", "workflow", "agent-chat", "channel", "all"}, in.Mode) {
		mlog.Errorf("mode[%s] must be in set[\"chat\", \"workflow\", \"agent-chat\", \"channel\", \"all\"]", in.Mode)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("mode[%s] must be in set[\"chat\", \"workflow\", \"agent-chat\", \"channel\", \"all\"]", in.Mode))
		return
	}
	for idx, v := range in.TagIds {
		if _, err1 := uuid.FromString(v); err1 != nil {
			mlog.Errorf("TagIds[%d]=%s must be uuid format", idx, v)
			response.InvalidArgErrorWithDetail(c, fmt.Sprintf("TagIds[%d]=%s must be uuid format", idx, v))
			return
		}
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).ListApps(c, in)
		if err != nil {
			mlog.Errorf("remote call ListApps failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call ListApps failed",
			})
			return
		} else {
			mlog.Infof("remote call ListApps return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.PaginateAppsStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusNoContent, gin.H{"data": []any{}, "total": 0, "page": in.Page, "limit": in.Limit, "has_more": false, "result": "failed", "code": 0})
				} else {
					c.JSON(http.StatusOK, response.NewAppPaginationResponse(pbrsp.PaginateAppsStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AppApi) Create(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	in := &pbapi.CreateAppRequest{
		AccountId: acc.ID,
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if in.Name == "" {
		mlog.Errorf("missing name")
		response.InvalidArgErrorWithDetail(c, "missing name")
		return
	}
	if !slices.Contains(constants.ALLOW_CREATE_APP_MODES, in.Mode) {
		mlog.Errorf("app mode[%s] must be in set %v", in.Mode, constants.ALLOW_CREATE_APP_MODES)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("app mode[%s] must be in set %v", in.Mode, constants.ALLOW_CREATE_APP_MODES))
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).CreateApp(c, in)
		if err != nil {
			mlog.Errorf("remote call CreateApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call CreateApp failed",
			})
			return
		} else {
			mlog.Infof("remote call CreateApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusCreated, response.NewAppDetailResponse(pbrsp.AppDetailStr))
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AppApi) SetName(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	in := pbapi.SetAppNameRequest{
		AccountId: acc.ID,
		AppId:     app_id,
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if in.Name == "" {
		mlog.Error("missing name arg")
		response.InvalidArgErrorWithDetail(c, "missing name arg")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).SetAppName(c, &in)
		if err != nil {
			mlog.Errorf("remote call SetAppName failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call SetAppName failed",
			})
			return
		} else {
			mlog.Infof("remote call SetAppName return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailStr == "" {
					mlog.Error("查询失败!")
					response.AppNotFoundError(c)
				} else {
					c.JSON(http.StatusOK, response.NewAppDetailResponse(pbrsp.AppDetailStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *AppApi) SetIcon(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	in := pbapi.SetAppIconRequest{
		AccountId: acc.ID,
		AppId:     app_id,
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).SetAppIcon(c, &in)
		if err != nil {
			mlog.Errorf("remote call SetAppIcon failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call SetAppIcon failed",
			})
			return
		} else {
			mlog.Infof("remote call SetAppIcon return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailStr == "" {
					mlog.Error("查询失败!")
					response.AppNotFoundError(c)
				} else {
					c.JSON(http.StatusOK, response.NewAppDetailResponse(pbrsp.AppDetailStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *AppApi) UpdateSiteStatus(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	type Req struct {
		EnableSite bool `json:"enable_site"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AppUpdateSiteStatus(c, &pbapi.AppUpdateSiteStatusRequest{AccountId: acc.ID, AppId: app_id, EnableSite: req.EnableSite})
		if err != nil {
			mlog.Errorf("remote call AppUpdateSiteStatus failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AppUpdateSiteStatus failed",
			})
			return
		} else {
			mlog.Infof("remote call AppUpdateSiteStatus return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailStr == "" {
					mlog.Error("查询失败!")
					response.AppNotFoundError(c)
				} else {
					c.JSON(http.StatusOK, response.NewAppDetailResponse(pbrsp.AppDetailStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *AppApi) UpdateApiStatus(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	type Req struct {
		EnableApi bool `json:"enable_api"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AppUpdateApiStatus(c, &pbapi.AppUpdateApiStatusRequest{AccountId: acc.ID, AppId: app_id, EnableApi: req.EnableApi})
		if err != nil {
			mlog.Errorf("remote call AppUpdateApiStatus failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AppUpdateApiStatus failed",
			})
			return
		} else {
			mlog.Infof("remote call AppUpdateApiStatus return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailStr == "" {
					mlog.Error("查询失败!")
					response.AppNotFoundError(c)
				} else {
					c.JSON(http.StatusOK, response.NewAppDetailResponse(pbrsp.AppDetailStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AppApi) GetWorkflowAppLogList(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("missing app_id")
		response.InvalidArgErrorWithDetail(c, "missing app_id")
		return
	}

	in := pbapi.GetWorkflowAppLogListRequest{Page: 1, Limit: 20, AccountId: acc.ID, AppId: app_id}
	if err := c.ShouldBindQuery(&in); err != nil {
		mlog.Error("bind query failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if in.Status != "" {
		if !slices.Contains([]string{"succeeded", "failed", "stopped"}, in.Status) {
			mlog.Errorf("Status=%s must be succeeded, failed or stopped", in.Status)
			response.InvalidArgErrorWithDetail(c, fmt.Sprintf("Status=%s must be succeeded, failed or stopped", in.Status))
			return
		}
	}
	if in.Limit < 1 || in.Limit > 100 {
		mlog.Errorf("limit=%d must be in range[1, 100]", in.Limit)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("limit=%d must be in range[1, 100]", in.Limit))
		return
	}
	if in.Page < 1 || in.Page > 99999 {
		mlog.Error("page=%d must be in range[1, 99999]", in.Page)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("page=%d must be in range[1, 99999]", in.Page))
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowAppLogList(c, &in)
		if err != nil {
			mlog.Errorf("remote call GetWorkflowAppLogList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowAppLogList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowAppLogList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.WorkflowAppLogPaginationStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewWorkflowAppLogPaginationResponse(pbrsp.WorkflowAppLogPaginationStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AppApi) GetTrace(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AppGetTrace(c, &pbapi.AppGetTraceRequest{AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call AppGetTrace failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AppGetTrace failed",
			})
			return
		} else {
			mlog.Infof("remote call AppGetTrace return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				app_trace_config := map[string]any{}
				err := json.Unmarshal([]byte(pbrsp.AppTraceConfigStr), &app_trace_config)
				if err != nil {
					mlog.Errorf("json unmarshal(%s) to dict failed:%v", pbrsp.AppTraceConfigStr, err)
				}
				c.JSON(http.StatusOK, app_trace_config)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *AppApi) SetTrace(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	type Req struct {
		Enabled         bool   `json:"enabled" form:"enabled"`
		TracingProvider string `json:"tracing_provider" form:"tracing_provider"`
	}
	in := pbapi.AppSetTraceRequest{AppId: app_id}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if in.TracingProvider != "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": "tracing_provider args is needed"})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AppSetTrace(c, &in)
		if err != nil {
			mlog.Errorf("remote call AppSetTrace failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AppSetTrace failed",
			})
			return
		} else {
			mlog.Infof("remote call AppSetTrace return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"result": "success"})
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AppApi) Copy(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		c.JSON(http.StatusNoContent, map[string]any{"result": "invalid id"})
		return
	}
	in := pbapi.CopyAppRequest{
		AccountId: acc.ID,
		AppId:     app_id,
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	if in.Name == "" {
		mlog.Error("missing name arg")
		response.InvalidArgErrorWithDetail(c, "missing name arg")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).CopyApp(c, &in)
		if err != nil {
			mlog.Errorf("remote call CopyApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call CopyApp failed",
			})
			return
		} else {
			mlog.Infof("remote call CopyApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AppDetailWithSiteStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusCreated, response.NewAppDetailWithSiteResponse(pbrsp.AppDetailWithSiteStr))
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *AppApi) Export(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		c.JSON(http.StatusNoContent, map[string]any{"result": "invalid id"})
		return
	}
	var args struct {
		IncludeSecret bool `json:"include_secret"`
	}
	err := c.ShouldBindQuery(&args)
	if err != nil {
		mlog.Error("bind failed:", err)
		response.InvalidArgError(c)
		return
	}
	in := pbapi.ExportAppRequest{
		AccountId:     acc.ID,
		AppId:         app_id,
		IncludeSecret: args.IncludeSecret,
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).ExportApp(c, &in)
		if err != nil {
			mlog.Errorf("remote call ExportApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call ExportApp failed",
			})
			return
		} else {
			mlog.Infof("remote call ExportApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"data": pbrsp.Dsl})
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *AppApi) Import(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	in := pbapi.ImportAppRequest{AccountId: acc.ID}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		mlog.Error("bind failed:", err)
		response.InvalidArgError(c)
		return
	}
	if in.Mode == "" {
		mlog.Error("missing mode arg")
		c.JSON(http.StatusNoContent, map[string]any{"result": "missing mode arg"})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).ImportApp(c, &in)
		if err != nil {
			mlog.Errorf("remote call ImportApp failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call ImportApp failed",
			})
			return
		} else {
			mlog.Infof("remote call ImportApp return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(int(pbrsp.HttpStatus), pbrsp.Import)
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
