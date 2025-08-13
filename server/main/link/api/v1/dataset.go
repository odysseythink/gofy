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

type DatasetApi struct {
}

func (api *DatasetApi) RetrievalSetting(c *gin.Context) {
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetDatasetRetrievalSetting(c, &pbapi.GetDatasetRetrievalSettingRequest{})
		if err != nil {
			mlog.Errorf("remote call GetDatasetRetrievalSetting failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetDatasetRetrievalSetting failed",
			})
			return
		} else {
			mlog.Infof("remote call GetDatasetRetrievalSetting return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, pbrsp.RetrievalMethod)
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
func (api *DatasetApi) DatasetList(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	in := &pbapi.DatasetListRequest{
		AccountId: acc.ID,
		Page:      1,
		Limit:     20,
	}
	if err := c.ShouldBindQuery(in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("datasets")
	if conn != nil {
		pbrsp, err := pbapi.NewDatasetsClient(conn).DatasetList(c, in)
		if err != nil {
			mlog.Errorf("remote call DatasetList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DatasetList failed",
			})
			return
		} else {
			mlog.Infof("remote call DatasetList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.DatasetsStr == "" {
					c.JSON(http.StatusOK, gin.H{"data": []*models.DatasetDetailFields{}, "has_more": pbrsp.HasMore, "limit": pbrsp.Limit, "total": pbrsp.Total, "page": pbrsp.Page})
				} else {
					var data []*models.DatasetDetailFields
					err := json.Unmarshal([]byte(pbrsp.DatasetsStr), &data)
					if err != nil {
						mlog.Errorf("json unmarshal %s failed:%v", pbrsp.DatasetsStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   "json unmarshal datasets string failed",
						})
						return
					}
					c.JSON(http.StatusOK, gin.H{"data": data, "has_more": pbrsp.HasMore, "limit": pbrsp.Limit, "total": pbrsp.Total, "page": pbrsp.Page})
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
func (api *DatasetApi) ExternalKnowledgeApiList(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	in := &pbapi.ExternalKnowledgeApiListRequest{
		CurrentTenantId: acc.CurrentTenantID(),
		Page:            1,
		Limit:           20,
	}
	if err := c.ShouldBindQuery(in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("datasets")
	if conn != nil {
		pbrsp, err := pbapi.NewDatasetsClient(conn).ExternalKnowledgeApiList(c, in)
		if err != nil {
			mlog.Errorf("remote call ExternalKnowledgeApiList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call ExternalKnowledgeApiList failed",
			})
			return
		} else {
			mlog.Infof("remote call ExternalKnowledgeApiList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.DatasStr == "" {
					c.JSON(http.StatusOK, gin.H{"data": []map[string]any{}, "has_more": pbrsp.HasMore, "limit": pbrsp.Limit, "total": pbrsp.Total, "page": pbrsp.Page})
				} else {
					var data []map[string]any
					err := json.Unmarshal([]byte(pbrsp.DatasStr), &data)
					if err != nil {
						mlog.Errorf("json unmarshal %s failed:%v", pbrsp.DatasStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   "json unmarshal datas string failed",
						})
						return
					}
					c.JSON(http.StatusOK, gin.H{"data": data, "has_more": pbrsp.HasMore, "limit": pbrsp.Limit, "total": pbrsp.Total, "page": pbrsp.Page})
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
