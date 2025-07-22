package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/libs/password"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type SetupApi struct {
}

// 登录
// @Tags
// @Summary
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /console/api/setup [GET]
func (api *SetupApi) GetSetupStatus(c *gin.Context) {
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetSetupStatus(c, &pbapi.GetSetupStatusRequest{})
		if err != nil {
			mlog.Errorf("remote call GetSetupStatus failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "internal_server_error",
				"status":  http.StatusInternalServerError,
				"message": "remote call GetSetupStatus failed",
			})
			return
		} else {
			mlog.Infof("remote call Login return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, pbrsp)
			}

			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "internal_server_error",
			"status":  http.StatusInternalServerError,
			"message": "get rpc client failed",
		})
		return
	}
}

func (api *SetupApi) Setup(c *gin.Context) {
	req := pbapi.SetupRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		mlog.Error("bind failed:", err)
		response.InvalidArgError(c)
		return
	}
	mlog.Debugf("---req=%#v", req)
	if req.Email == "" || req.Password == "" || req.Name == "" {
		mlog.Error("missing email or password or name")
		response.InvalidArgError(c)
		return
	}
	_, err = password.ValidPassword(req.Password)
	if err != nil {
		mlog.Error("invalid password")
		response.InvalidArgError(c)
		return
	}

	req.ClientIp = c.ClientIP()

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).Setup(c, &req)
		if err != nil {
			mlog.Errorf("remote call Setup failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "internal_server_error",
				"status":  http.StatusInternalServerError,
				"message": "remote call Setup failed",
			})
			return
		} else {
			mlog.Infof("remote call Setup return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusCreated, gin.H{"result": "success"})
			}

			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "internal_server_error",
			"status":  http.StatusInternalServerError,
			"message": "get rpc client failed",
		})
		return
	}
}
