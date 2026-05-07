package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/response"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
)

type ApiKeyApi struct {
}

func (api *ApiKeyApi) GetApiKeyListResource(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	method := strings.Split(c.Request.URL.Path, app_id)[0]

	in := pbapi.GetApiKeyListRequest{AccountId: acc.ID}
	if strings.HasSuffix(method, "/apps/") {
		in.ResourceId = app_id
		in.ResourceType = "app"
		in.ResourceIdField = "app_id"
		in.TokenPrefix = "app-"
	} else if strings.HasSuffix(method, "/datasets/") {
		in.ResourceId = app_id
		in.ResourceType = "dataset"
		in.ResourceIdField = "dataset_id"
		in.TokenPrefix = "ds-"
	} else {
		mlog.Errorf("unsupported method=", method)
		response.InvalidArgErrorWithDetail(c, "unsupported method="+method)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetApiKeyList(c, &in)
		if err != nil {
			mlog.Errorf("remote call GetApiKeyList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetApiKeyList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetApiKeyList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"items": pbrsp.Items})
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

func (api *ApiKeyApi) SetApiKeyListResource(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	method := strings.Split(c.Request.URL.Path, app_id)[0]

	in := pbapi.GenerateApiKeyRequest{AccountId: acc.ID}
	if strings.HasSuffix(method, "/apps/") {
		in.ResourceId = app_id
		in.ResourceType = "app"
		in.ResourceIdField = "app_id"
		in.TokenPrefix = "app-"
	} else if strings.HasSuffix(method, "/datasets/") {
		in.ResourceId = app_id
		in.ResourceType = "dataset"
		in.ResourceIdField = "dataset_id"
		in.TokenPrefix = "ds-"
	} else {
		mlog.Errorf("unsupported method=", method)
		response.InvalidArgErrorWithDetail(c, "unsupported method="+method)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GenerateApiKey(c, &in)
		if err != nil {
			mlog.Errorf("remote call GenerateApiKey failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GenerateApiKey failed",
			})
			return
		} else {
			mlog.Infof("remote call GenerateApiKey return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusCreated, pbrsp.Item)
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

func (api *ApiKeyApi) DelApiKeyResource(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	api_key_id := c.Param("api_key_id")
	if api_key_id == "" {
		mlog.Error("invalid api_key_id")
		response.InvalidArgError(c)
		return
	}
	method := strings.Split(c.Request.URL.Path, app_id)[0]

	in := pbapi.DelApiKeyRequest{AccountId: acc.ID, ApiKeyId: api_key_id, ResourceId: app_id}
	if strings.HasSuffix(method, "/apps/") {
		in.ResourceType = "app"
		in.ResourceIdField = "app_id"
		in.TokenPrefix = "app-"
	} else if strings.HasSuffix(method, "/datasets/") {
		in.ResourceType = "dataset"
		in.ResourceIdField = "dataset_id"
		in.TokenPrefix = "ds-"
	} else {
		mlog.Errorf("unsupported method=", method)
		response.InvalidArgErrorWithDetail(c, "unsupported method="+method)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).DelApiKey(c, &in)
		if err != nil {
			mlog.Errorf("remote call DelApiKey failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DelApiKey failed",
			})
			return
		} else {
			mlog.Infof("remote call DelApiKey return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusNoContent, gin.H{"result": "success"})
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
