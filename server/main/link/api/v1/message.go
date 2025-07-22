package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type MessageApi struct {
}

func (api *MessageApi) GetSuggestedQuestion(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	message_id := c.Param("message_id")
	if message_id == "" {
		mlog.Error("invalid message_id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetSuggestedQuestionMessage(c, &pbapi.GetSuggestedQuestionMessageRequest{
			AppId:     app_id,
			AccountId: acc.ID,
			MessageId: message_id,
		})
		if err != nil {
			mlog.Errorf("remote call GetSuggestedQuestionMessage failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetSuggestedQuestionMessage failed",
			})
			return
		} else {
			mlog.Infof("remote call GetSuggestedQuestionMessage return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"data": pbrsp.Questions})
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
func (api *MessageApi) ChatMessageList(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	conversation_id := c.Query("conversation_id")
	if conversation_id == "" {
		mlog.Error("missing conversation_id")
		response.InvalidArgError(c)
		return
	}
	first_id := c.Query("first_id")
	limit_str := c.Query("limit")
	limit := 20
	if limit_str != "" {
		var err error
		limit, err = strconv.Atoi(limit_str)
		if err != nil {
			mlog.Errorf("limit(%s) convert to int failed:%v", limit_str, err)
			response.InvalidArgError(c)
			return
		} else if limit < 1 || limit > 100 {
			mlog.Errorf("limit(%s) is not in range[1, 100]", limit_str)
			response.InvalidArgError(c)
			return
		}
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).ChatMessageList(c, &pbapi.ChatMessageListRequest{
			AppId:          app_id,
			AccountId:      acc.ID,
			ConversationId: conversation_id,
			FirstId:        first_id,
			Limit:          int32(limit),
		})
		if err != nil {
			mlog.Errorf("remote call ChatMessageList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call ChatMessageList failed",
			})
			return
		} else {
			mlog.Infof("remote call ChatMessageList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.MessageInfiniteScrollPaginationStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewMessageInfiniteScrollPaginationResponse(pbrsp.MessageInfiniteScrollPaginationStr))
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
func (api *MessageApi) MessageFeedback(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	in := pbapi.MessageFeedbackRequest{
		AppId:     app_id,
		AccountId: acc.ID,
	}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	if _, err = uuid.FromString(in.MessageId); err != nil {
		mlog.Error("message_id is invalid uuid")
		response.InvalidArgErrorWithDetail(c, "message_id is invalid uuid")
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).MessageFeedback(c, &in)
		if err != nil {
			mlog.Errorf("remote call MessageFeedback failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call MessageFeedback failed",
			})
			return
		} else {
			mlog.Infof("remote call MessageFeedback return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"result": "success"})
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
func (api *MessageApi) SetMessageAnnotation(c *gin.Context) {
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
	var args struct {
		MessageID       string         `json:"message_id"`
		Question        string         `json:"question"`
		Answer          string         `json:"answer"`
		AnnotationReply map[string]any `json:"annotation_reply"`
	}

	err := c.ShouldBindJSON(&args)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	if args.Question == "" {
		mlog.Error("missing question")
		response.InvalidArgError(c)
		return
	}
	if args.Answer == "" {
		mlog.Error("missing answer")
		response.InvalidArgError(c)
		return
	}
	annotation_reply_str := ""
	if len(args.AnnotationReply) > 0 {
		bindata, _ := json.Marshal(args.AnnotationReply)
		annotation_reply_str = string(bindata)
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).SetMessageAnnotation(c, &pbapi.SetMessageAnnotationRequest{
			AppId:              app_id,
			AccountId:          acc.ID,
			MessageId:          args.MessageID,
			Question:           args.Question,
			Answer:             args.Answer,
			AnnotationReplyStr: annotation_reply_str,
		})
		if err != nil {
			mlog.Errorf("remote call SetMessageAnnotation failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call SetMessageAnnotation failed",
			})
			return
		} else {
			mlog.Infof("remote call SetMessageAnnotation return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.AnnotationStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewAnnotationResponse(pbrsp.AnnotationStr))
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

func (api *MessageApi) MessageAnnotationCount(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).MessageAnnotationCount(c, &pbapi.MessageAnnotationCountRequest{
			AppId:     app_id,
			AccountId: acc.ID,
		})
		if err != nil {
			mlog.Errorf("remote call MessageAnnotationCount failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call MessageAnnotationCount failed",
			})
			return
		} else {
			mlog.Infof("remote call MessageAnnotationCount return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"count": pbrsp.Count})
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
func (api *MessageApi) GetMessage(c *gin.Context) {
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	message_id := c.Param("message_id")
	if message_id == "" {
		mlog.Error("invalid message_id")
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetMessage(c, &pbapi.GetMessageRequest{
			AppId:     app_id,
			AccountId: acc.ID,
			MessageId: message_id,
		})
		if err != nil {
			mlog.Errorf("remote call GetMessage failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetMessage failed",
			})
			return
		} else {
			mlog.Infof("remote call GetMessage return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.MessageDetailStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewMessageDetailResponse(pbrsp.MessageDetailStr))
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
