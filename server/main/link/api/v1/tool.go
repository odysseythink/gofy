package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/entities/tools"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type ToolsApi struct {
}

func (api *ToolsApi) ListToolLabels(c *gin.Context) {
	conn := cluster.Instance().GetRpcClientByModule("tools")
	if conn != nil {
		pbrsp, err := pbapi.NewToolsClient(conn).GetToolLabelsList(c, &pbapi.GetToolLabelsListRequest{})
		if err != nil {
			mlog.Errorf("remote call GetToolLabelsList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetToolLabelsList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetToolLabelsList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ToolLabelsListStr != "" {
					var tmp []*tools.ToolLabel
					if err1 := json.Unmarshal([]byte(pbrsp.ToolLabelsListStr), &tmp); err1 != nil {
						mlog.Errorf("json unmarshal(%s) to ToolLabel list failed:%v", pbrsp.ToolLabelsListStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   fmt.Sprintf("json unmarshal(%s) to ToolLabel list failed:%v", pbrsp.ToolLabelsListStr, err),
						})
					} else {
						c.JSON(http.StatusOK, tmp)
					}
				} else {
					c.JSON(http.StatusOK, nil)
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
func (api *ToolsApi) ListToolProvider(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tool_type := c.Query("type")
	if tool_type != "" {
		if !slices.Contains([]string{"builtin", "model", "api", "workflow"}, tool_type) {
			mlog.Errorf("Invalid tool_type=%s.", tool_type)
			response.InvalidArgErrorWithDetail(c, fmt.Sprintf("Invalid tool_type=%s.", tool_type))
			return
		}
	}
	conn := cluster.Instance().GetRpcClientByModule("tools")
	if conn != nil {
		pbrsp, err := pbapi.NewToolsClient(conn).GetToolProviderList(c, &pbapi.GetToolProviderListRequest{Type: tool_type, AccountId: acc.ID, TenantId: acc.CurrentTenantID()})
		if err != nil {
			mlog.Errorf("remote call GetToolProviderList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetToolProviderList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetToolProviderList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ProvidersStr != "" {
					var tmp []map[string]any
					if err1 := json.Unmarshal([]byte(pbrsp.ProvidersStr), &tmp); err1 != nil {
						mlog.Errorf("json unmarshal(%s) to dict list failed:%v", pbrsp.ProvidersStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   fmt.Sprintf("json unmarshal(%s) to dict list failed:%v", pbrsp.ProvidersStr, err),
						})
					} else {
						c.JSON(http.StatusOK, tmp)
					}
				} else {
					c.JSON(http.StatusOK, []map[string]any{})
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

func (api *ToolsApi) GetToolList(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tool_type := c.Param("tool_type")
	if !slices.Contains([]string{"builtin", "api", "workflow"}, tool_type) {
		mlog.Errorf("Invalid tool_type=%s.", tool_type)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("Invalid tool_type=%s.", tool_type))
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("tools")
	if conn != nil {
		pbrsp, err := pbapi.NewToolsClient(conn).GetToolList(c, &pbapi.GetToolListRequest{Type: tool_type, AccountId: acc.ID, TenantId: acc.CurrentTenantID()})
		if err != nil {
			mlog.Errorf("remote call GetToolProviderList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetToolProviderList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetToolProviderList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ToolsStr != "" {
					var tmp []map[string]any
					if err1 := json.Unmarshal([]byte(pbrsp.ToolsStr), &tmp); err1 != nil {
						mlog.Errorf("json unmarshal(%s) to dict list failed:%v", pbrsp.ToolsStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   fmt.Sprintf("json unmarshal(%s) to dict list failed:%v", pbrsp.ToolsStr, err),
						})
					} else {
						c.JSON(http.StatusOK, tmp)
					}
				} else {
					c.JSON(http.StatusOK, []map[string]any{})
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
