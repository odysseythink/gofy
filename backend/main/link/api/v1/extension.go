package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

type ExtensionApi struct {
}

func (api *ExtensionApi) GetCodeBasedExtension(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		mlog.Errorf("module is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "module is missing"})
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetCodeBasedExtension(c, &pbapi.GetCodeBasedExtensionRequest{Module: module})
		if err != nil {
			mlog.Errorf("remote call GetCodeBasedExtension failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetCodeBasedExtension failed",
			})
			return
		} else {
			mlog.Infof("remote call GetCodeBasedExtension return:%#v", pbrsp)
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
