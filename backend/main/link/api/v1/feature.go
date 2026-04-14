package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

var (
	currentFeature = `{
    "billing": {
        "enabled": false,
        "subscription": {
            "plan": "sandbox",
            "interval": ""
        }
    },
    "members": {
        "size": 0,
        "limit": 1
    },
    "apps": {
        "size": 0,
        "limit": 10
    },
    "vector_space": {
        "size": 0,
        "limit": 5
    },
    "annotation_quota_limit": {
        "size": 0,
        "limit": 10
    },
    "documents_upload_quota": {
        "size": 0,
        "limit": 50
    },
    "docs_processing": "standard",
    "can_replace_logo": false,
    "model_load_balancing_enabled": false,
    "dataset_operator_enabled": false
}
`
)

type FeatureApi struct {
}

func (api *FeatureApi) ListSystem(c *gin.Context) {
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetSystemFeatures(c, &pbapi.GetSystemFeaturesRequest{})
		if err != nil {
			mlog.Errorf("remote call GetSystemFeatures failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetSystemFeatures failed",
			})
			return
		} else {
			mlog.Infof("remote call GetSystemFeatures return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, pbrsp.Feature)
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

func (api *FeatureApi) List(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetFeatures(c, &pbapi.GetFeaturesRequest{TenantId: tenant_id})
		if err != nil {
			mlog.Errorf("remote call GetFeatures failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetFeatures failed",
			})
			return
		} else {
			mlog.Infof("remote call GetFeatures return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, pbrsp.Feature)
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
