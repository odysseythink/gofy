package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/response"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
)

type AccountApi struct {
}

// Profile
// @Tags Account
// @Summary Profile
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body models.Account true "Profile"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登出成功"}"
// @Router /console/api/account/profile [GET]
func (api *AccountApi) Profile(c *gin.Context) {
	rawuser, _ := c.Get("current_user")

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetAccountProfile(c, &pbapi.GetAccountProfileRequest{AccountId: rawuser.(*models.Account).ID})
		if err != nil {
			mlog.Errorf("remote call GetAccountProfile failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetAccountProfile failed",
			})
			return
		} else {
			mlog.Infof("remote call GetAccountProfile return:%#v", pbrsp)
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
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *AccountApi) Update(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	var args map[string]string
	err := c.ShouldBindJSON(&args)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	if len(args) == 0 {
		mlog.Error("no update account field")
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).UpdateAccount(c, &pbapi.UpdateAccountRequest{AccountId: acc.ID, Infos: args})
		if err != nil {
			mlog.Errorf("remote call UpdateAccount failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call UpdateAccount failed",
			})
			return
		} else {
			mlog.Infof("remote call UpdateAccount return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AccountResponseStr == "" {
					mlog.Error("update failed")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewAccountResponse(pbrsp.AccountResponseStr))
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
