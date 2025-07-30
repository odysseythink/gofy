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
