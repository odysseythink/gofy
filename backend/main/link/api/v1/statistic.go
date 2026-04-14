package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

type StatisticApi struct {
}

func (api *StatisticApi) Statistic(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid id")
		response.InvalidArgError(c)
		return
	}
	type Req struct {
		Start string `json:"start" form:"start"`
		End   string `json:"end" form:"end"`
	}
	req := Req{}
	if err := c.ShouldBindQuery(&req); err != nil {
		mlog.Errorf("ShouldBindQuery failed:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_, err := time.Parse("2006-01-02 15:04", req.Start)
	if err != nil {
		mlog.Errorf("parse time(%s) failed:%v", req.Start, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_, err = time.Parse("2006-01-02 15:04", req.End)
	if err != nil {
		mlog.Errorf("parse time(%s) failed:%v", req.End, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	method := strings.Split(c.Request.URL.Path, app_id)[1]

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).Statistic(c, &pbapi.StatisticRequest{AccountId: acc.ID, AppId: app_id, Start: req.Start, End: req.End, Method: method})
		if err != nil {
			mlog.Errorf("remote call Statistic failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call Statistic failed",
			})
			return
		} else {
			mlog.Infof("remote call Statistic return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ResultStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					data := []map[string]any{}
					err := json.Unmarshal([]byte(pbrsp.ResultStr), &data)
					if err != nil {
						mlog.Errorf("json unmarshal(%s) to object list failed:%v", pbrsp.ResultStr, err)
					}
					c.JSON(http.StatusOK, gin.H{"data": data})
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
