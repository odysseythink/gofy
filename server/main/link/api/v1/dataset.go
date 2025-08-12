package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
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
	in := &pbapi.DatasetListRequest{
		Page:  1,
		Limit: 20,
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
