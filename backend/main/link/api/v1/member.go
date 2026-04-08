package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type MemberApi struct {
}

func (api *MemberApi) List(c *gin.Context) {
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
		pbrsp, err := pbapi.NewAdminClient(conn).GetMemberList(c, &pbapi.GetMemberListRequest{TenantId: tenant_id})
		if err != nil {
			mlog.Errorf("remote call GetMemberList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetMemberList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetMemberList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.MembersStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					rsp := response.NewAccountWithRoleListResponse(pbrsp.MembersStr)
					if rsp.Accounts == nil {
						rsp.Accounts = make([]*response.AccountWithRoleResponse, 0)
					}
					c.JSON(http.StatusOK, gin.H{"result": "success", "accounts": rsp.Accounts})
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
