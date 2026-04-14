package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

var (
	builtin_providers = ``
)

type WorkspaceApi struct {
}

func (api *WorkspaceApi) List(c *gin.Context) {
	in := pbapi.GetWorkspaceListRequest{Page: 1, Limit: 20}
	if err := c.ShouldBindQuery(&in); err != nil {
		mlog.Errorf("bind query failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkspaceList(c, &in)
		if err != nil {
			mlog.Errorf("remote call GetWorkspaceList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkspaceList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkspaceList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.Code != 0 {
					c.JSON(http.StatusBadRequest, pbrsp)
				} else {
					c.JSON(http.StatusOK, pbrsp)
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

func (api *WorkspaceApi) All(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetTenantList(c, &pbapi.GetTenantListRequest{AccountId: acc.ID})
		if err != nil {
			mlog.Errorf("remote call GetTenantList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetTenantList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetTenantList return:%#v", pbrsp)
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

func (api *WorkspaceApi) GetCurrent(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetCurrentTenant(c, &pbapi.GetCurrentTenantRequest{AccountId: acc.ID})
		if err != nil {
			mlog.Errorf("remote call GetCurrentTenant failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetCurrentTenant failed",
			})
			return
		} else {
			mlog.Infof("remote call GetCurrentTenant return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ResultDictStr == "" {
					mlog.Errorf("result dict string is empty")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					tmp := map[string]any{}
					if err1 := json.Unmarshal([]byte(pbrsp.ResultDictStr), &tmp); err1 != nil {
						mlog.Errorf("json unmarshal result_dict_str=%s to dict failed:%v", pbrsp.ResultDictStr, err1)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   "internal server error",
						})
					} else {
						c.JSON(http.StatusOK, tmp)
					}
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
