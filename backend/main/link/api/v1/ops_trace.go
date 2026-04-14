package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

type OpsTraceApi struct {
}

func (api *OpsTraceApi) GetTraceAppConfig(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	type Req struct {
		TracingProvider string `json:"tracing_provider" form:"tracing_provider"`
	}
	req := Req{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if req.TracingProvider == "" {
		mlog.Error("missing tracing_provider")
		response.InvalidArgError(c)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetTraceAppConfig(c, &pbapi.GetTraceAppConfigRequest{AppId: app_id, TracingProvider: req.TracingProvider})
		if err != nil {
			mlog.Errorf("remote call GetTraceAppConfig failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetTraceAppConfig failed",
			})
			return
		} else {
			mlog.Infof("remote call GetTraceAppConfig return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.TraceConfigStr == "" {
					c.JSON(http.StatusOK, gin.H{"has_not_configured": true})
				} else {
					trace_config := map[string]any{}
					err := json.Unmarshal([]byte(pbrsp.TraceConfigStr), &trace_config)
					if err != nil {
						mlog.Errorf("json unmarshal(%s) to dict failed:%v", pbrsp.TraceConfigStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   "internal server error",
						})
					} else {
						c.JSON(http.StatusOK, trace_config)
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
