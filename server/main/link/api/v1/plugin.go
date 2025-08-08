package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type PluginApi struct {
}

func (api *PluginApi) PluginFetchPreferences(c *gin.Context) {
	rawuser, _ := c.Get("current_user")

	conn := cluster.Instance().GetRpcClientByModule("plugins")
	if conn != nil {
		pbrsp, err := pbapi.NewPluginsClient(conn).FetchPreferences(c, &pbapi.FetchPreferencesRequest{TenantId: rawuser.(*models.Account).CurrentTenantID()})
		if err != nil {
			mlog.Errorf("remote call FetchPreferences failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "internal_server_error",
				"status":  http.StatusInternalServerError,
				"message": "remote call FetchPreferences failed",
			})
			return
		} else {
			mlog.Infof("remote call FetchPreferences return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				var auto_upgrade_dict map[string]any
				json.Unmarshal([]byte(pbrsp.AutoUpgradeDictStr), &auto_upgrade_dict)
				c.JSON(http.StatusOK, gin.H{"permission": pbrsp.PermissionDict, "auto_upgrade": auto_upgrade_dict})
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

func (api *PluginApi) PluginFetchInstallTasks(c *gin.Context) {
	rawuser, _ := c.Get("current_user")

	in := &pbapi.FetchInstallTasksRequest{
		TenantId: rawuser.(*models.Account).CurrentTenantID(),
		Page:     1,
		PageSize: 20,
	}
	if err := c.ShouldBindQuery(in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("plugins")
	if conn != nil {
		pbrsp, err := pbapi.NewPluginsClient(conn).FetchInstallTasks(c, in)
		if err != nil {
			mlog.Errorf("remote call FetchInstallTasks failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "internal_server_error",
				"status":  http.StatusInternalServerError,
				"message": "remote call FetchInstallTasks failed",
			})
			return
		} else {
			mlog.Infof("remote call FetchInstallTasks return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				var tasks []*models.PluginInstallTask
				err = json.Unmarshal([]byte(pbrsp.TasksStr), &tasks)
				if tasks == nil {
					mlog.Errorf("TasksStr {%s} json unmarshal failed:%v", pbrsp.TasksStr, err)
					tasks = []*models.PluginInstallTask{}
				}
				c.JSON(http.StatusOK, gin.H{"tasks": tasks})
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
