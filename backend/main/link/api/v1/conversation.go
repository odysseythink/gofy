package v1

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
)

type ConversationApi struct {
}

func (api *ConversationApi) ChatConversationDetail(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	conversation_id := c.Param("conversation_id")
	if conversation_id == "" {
		mlog.Error("invalid conversation_id")
		response.InvalidArgError(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).ChatConversationDetail(c, &pbapi.ChatConversationDetailRequest{AccountId: acc.ID, AppId: app_id, ConversationId: conversation_id})
		if err != nil {
			mlog.Errorf("remote call ChatConversationDetail failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call ChatConversationDetail failed",
			})
			return
		} else {
			mlog.Infof("remote call ChatConversationDetail return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ConversationDetailStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewConversationDetailResponse(pbrsp.ConversationDetailStr))
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

func (api *ConversationApi) DelChatConversation(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	conversation_id := c.Param("conversation_id")
	if conversation_id == "" {
		mlog.Error("invalid conversation_id")
		response.InvalidArgError(c)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).DelChatConversation(c, &pbapi.DelChatConversationRequest{AccountId: acc.ID, AppId: app_id, ConversationId: conversation_id})
		if err != nil {
			mlog.Errorf("remote call DelChatConversation failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DelChatConversation failed",
			})
			return
		} else {
			mlog.Infof("remote call DelChatConversation return:%#v", pbrsp)
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

func (api *ConversationApi) GetChatConversationPagination(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	if !acc.IsEditor() {
		mlog.Errorf("current user(%#v) is not editor", acc)
		c.JSON(http.StatusBadRequest, gin.H{"result": "forbidden", "code": 7})
		return
	}
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	in := pbapi.GetChatConversationPaginationRequest{
		AnnotationStatus: "all",
		Page:             1,
		Limit:            20,
		SortBy:           "-updated_at",
	}

	err := c.ShouldBindQuery(&in)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	mlog.Debugf("------args=%#v", in)
	if in.Start != "" {
		_, err = time.Parse("2006-01-02 15:04", in.Start)
		if err != nil {
			mlog.Errorf("start=%s is invalid:%v", in.Start, err)
			c.JSON(http.StatusBadRequest, map[string]any{"result": "start arg has invalid format"})
			return
		}
	}
	if in.End != "" {
		_, err = time.Parse("2006-01-02 15:04", in.End)
		if err != nil {
			mlog.Errorf("end=%s is invalid:%v", in.End, err)
			c.JSON(http.StatusBadRequest, map[string]any{"result": "end arg has invalid format"})
			return
		}
	}
	if !slices.Contains([]string{"annotated", "not_annotated", "all"}, in.AnnotationStatus) {
		mlog.Errorf("annotation_status=%s must be \"annotated\", \"not_annotated\" or \"all\"", in.AnnotationStatus)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "annotation_status must be \"annotated\", \"not_annotated\" or \"all\""})
		return
	}
	if in.MessageCountGte != 0 && (in.MessageCountGte < 1 || in.MessageCountGte > 99999) {
		mlog.Errorf("message_count_gte=%d is beyond range(1, 99999)", in.MessageCountGte)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "message_count_gte is beyond range(1, 99999)"})
		return
	}
	if in.Page < 1 || in.Page > 99999 {
		mlog.Errorf("page=%d is beyond range(1, 99999)", in.Page)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "page is beyond range(1, 99999)"})
		return
	}
	if in.Limit < 1 || in.Page > 100 {
		mlog.Errorf("limit=%d is beyond range(1, 100)", in.Limit)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "limit is beyond range(1, 100)"})
		return
	}
	if !slices.Contains([]string{"created_at", "-created_at", "updated_at", "-updated_at"}, in.SortBy) {
		mlog.Errorf("sort_by=%s must be \"created_at\", \"-created_at\", \"updated_at\", \"-updated_at\"", in.SortBy)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "sort_by must be \"created_at\", \"-created_at\", \"updated_at\", \"-updated_at\""})
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetChatConversationPagination(c, &in)
		if err != nil {
			mlog.Errorf("remote call GetChatConversationPagination failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetChatConversationPagination failed",
			})
			return
		} else {
			mlog.Infof("remote call GetChatConversationPagination return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ConversationWithSummaryPaginationStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusCreated, response.NewConversationWithSummaryPaginationResponse(pbrsp.ConversationWithSummaryPaginationStr))
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
