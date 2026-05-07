package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/gofy/backend/cluster"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/response"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
)

type RuleGenerateApi struct {
}

func (api *RuleGenerateApi) RuleGenerate(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	var args struct {
		Instruction string         `json:"instruction" binding:"required"`
		ModelConfig map[string]any `json:"model_config" binding:"required"`
		NoVariable  bool           `json:"no_variable"`
	}
	err := c.ShouldBindJSON(&args)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	bindata, _ := json.Marshal(args.ModelConfig)

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).RuleGenerate(c, &pbapi.RuleGenerateRequest{TenantId: acc.CurrentTenantID(), Instruction: args.Instruction, ModelConfigStr: string(bindata), NoVariable: args.NoVariable})
		if err != nil {
			mlog.Errorf("remote call RuleGenerate failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call RuleGenerate failed",
			})
			return
		} else {
			mlog.Infof("remote call RuleGenerate return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.RulesDictStr == "" {
					mlog.Error("RuleGenerate failed")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					rules := map[string]any{}
					err := json.Unmarshal([]byte(pbrsp.RulesDictStr), &rules)
					if err != nil {
						mlog.Errorf("json unmarshal %s failed:%v", pbrsp.RulesDictStr, err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   "internal server error",
						})
					} else {
						c.JSON(http.StatusOK, rules)
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
