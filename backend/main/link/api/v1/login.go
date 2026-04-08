package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type LoginApi struct {
}

// 登录
// @Tags
// @Summary 登录Account
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Account true "登录Account"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /console/api/login [POST]
func (api *LoginApi) Login(c *gin.Context) {
	req := pbapi.LoginRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		mlog.Error("bind failed:", err)
		response.InvalidArgError(c)
		return
	}
	mlog.Debugf("---email=%#v", req.Email)
	if req.Email == "" || req.Password == "" {
		mlog.Error("invalid args")
		response.InvalidArgError(c)
		return
	}

	if req.Language == "" {
		req.Language = "en-US"
	}
	req.ClientIp = c.ClientIP()
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).Login(c, &req)
		if err != nil {
			mlog.Errorf("remote call Login failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call Login failed",
			})
			return
		} else {
			mlog.Infof("remote call Login return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.Errmsg != "" {
					c.JSON(http.StatusOK, gin.H{"result": pbrsp.Result, "data": pbrsp.Errmsg})
				} else {
					c.JSON(http.StatusOK, gin.H{"result": pbrsp.Result, "data": map[string]any{"access_token": pbrsp.AccessToken, "refresh_token": pbrsp.RefreshToken}})
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

// 登出
// @Tags
// @Summary 登出Account
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Account true "登出Account"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登出成功"}"
// @Router /console/api/logout [GET]
func (api *LoginApi) Logout(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).Logout(c, &pbapi.LogoutRequest{AccountId: acc.ID})
		if err != nil {
			mlog.Errorf("remote call Logout failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call Logout failed",
			})
			return
		} else {
			mlog.Infof("remote call Logout return:%#v", pbrsp)
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
