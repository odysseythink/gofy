package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/gofy/backend/cluster"
	httpexceptions "github.com/odysseythink/gofy/backend/core/exceptions/http"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/response"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
)

type WorkflowRunApi struct {
}

func (api *WorkflowRunApi) AppList(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	method := strings.Split(c.Request.URL.Path, app_id)[1]

	type Req struct {
		LastID string `json:"last_id"`
		Limit  int    `json:"limit"`
	}

	req := Req{Limit: 20}

	err := c.ShouldBindJSON(&req)
	if err != nil && err != io.EOF {
		mlog.Errorf("can't bind args: %v", err)
		response.InvalidArgError(c)
		return
	}

	if req.Limit < 1 || req.Limit > 100 {
		mlog.Errorf("invalid limit(%d) args", req.Limit)
		response.InvalidArgError(c)
		return
	}
	in := pbapi.AppWorkflowRunListRequest{
		AccountId: acc.ID,
		AppId:     app_id,
		LastId:    req.LastID,
		Limit:     int32(req.Limit),
	}
	if method == "/advanced-chat/workflow-runs" {
		in.AppMode = string(models.AppMode_ADVANCED_CHAT)
	} else if method != "/workflow-runs" && method != "/advanced-chat/workflow-runs" {
		in.AppMode = string(models.AppMode_WORKFLOW)
	} else {
		mlog.Errorf("unsupported method=%s", method)
		response.InvalidArgErrorWithDetail(c, "unsupported method="+method)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).AppWorkflowRunList(c, &in)
		if err != nil {
			mlog.Errorf("remote call AppWorkflowRunList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AppWorkflowRunList failed",
			})
			return
		} else {
			mlog.Infof("remote call AppWorkflowRunList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ResponseStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					if in.AppMode == string(models.AppMode_ADVANCED_CHAT) {
						c.JSON(http.StatusOK, response.NewAdvancedChatWorkflowRunPaginationResponse(pbrsp.ResponseStr))
					} else {
						c.JSON(http.StatusOK, response.NewWorkflowRunPaginationResponse(pbrsp.ResponseStr))
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

func (api *WorkflowRunApi) Detail(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	workflow_run_id := c.Param("workflow_run_id")
	if app_id == "" {
		mlog.Error("invalid workflow_run_id")
		response.InvalidArgError(c)
		return
	}
	method := strings.Split(c.Request.URL.Path, app_id)[1]
	in := pbapi.GetWorkflowRunDetailRequest{
		UserId:          acc.ID,
		AppId:           app_id,
		WorkflowRunId:   workflow_run_id,
		IsNodeExecution: strings.HasSuffix(method, "/node-executions"),
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowRunDetail(c, &in)
		if err != nil {
			mlog.Errorf("remote call GetWorkflowRunDetail failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowRunDetail failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowRunDetail return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ResponseStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					if !in.IsNodeExecution {
						c.JSON(http.StatusOK, response.NewWorkflowRunDetailResponse(pbrsp.ResponseStr))
					} else {
						c.JSON(http.StatusOK, response.NewWorkflowRunNodeExecutionListResponse(pbrsp.ResponseStr))
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

// func (api *WorkflowRunApi) Run(c *gin.Context) {
// 	rawapp, _ := c.Get("app_model")
// 	app := rawapp.(*models.App)

// 	if app.Mode != models.AppMode_WORKFLOW {
// 		panic(httpexceptions.NewNotWorkflowAppError())
// 	}
// 	type Req struct {
// 		Inputs       map[string]any `json:"inputs"`
// 		ResponseMode string         `json:"response_mode"`
// 		User         string         `json:"user"`
// 		Files        []any          `json:"files"`
// 	}
// 	var req Req
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
// 		return
// 	}
// 	if len(req.Inputs) == 0 {
// 		mlog.Error("missing inputs")
// 		response.InvalidArgError(c)
// 		return
// 	}
// 	if req.User == "" {
// 		mlog.Error("missing user")
// 		response.InvalidArgError(c)
// 		return
// 	}
// 	end_user := models.CreateOrGetEndUserByUserID(app, req.User)

// 	streaming := req.ResponseMode == "streaming"

// 	defer func() {
// 		if r := recover(); r != nil {
// 			if exp, ok := r.(*exceptions.ProviderTokenNotInitError); ok {
// 				panic(httpexceptions.NewProviderNotInitializeError(exp.Error()))
// 			} else if _, ok := r.(*exceptions.QuotaExceededError); ok {
// 				panic(httpexceptions.NewProviderQuotaExceededError())
// 			} else if _, ok := r.(*exceptions.ModelCurrentlyNotSupportError); ok {
// 				panic(httpexceptions.NewProviderModelCurrentlyNotSupportError())
// 			} else if exp, ok := r.(*modelruntimeexceptions.InvokeError); ok {
// 				panic(httpexceptions.NewCompletionRequestError(exp.Error()))
// 			} else if exp, ok := r.(*exceptions.ValueError); ok {
// 				panic(exp)
// 			} else {
// 				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
// 				panic(httpexceptions.NewInternalServerError(""))
// 			}
// 		}
// 	}()
// 	rsp, rspiter := services.ServiceGroupApp.AppGenerate.Generate(app, end_user, map[string]any{"inputs": req.Inputs, "response_mode": req.ResponseMode}, appenumtypes.InvokeFrom_SERVICE_API, streaming)
// 	mlog.Debugf("***********realrsp=%#v", rsp)
// 	if rsp != nil {
// 		c.JSON(http.StatusOK, rsp)
// 		return
// 	}

// 	for item := range rspiter {
// 		mlog.Debugf("-----send:%s", string(item))
// 		c.Data(http.StatusOK, "text/event-stream", []byte(item))
// 	}
// }

func (api *WorkflowRunApi) GetDetailOnlyByApp(c *gin.Context) {
	rawapp, _ := c.Get("app_model")
	app := rawapp.(*models.App)

	if app.Mode != models.AppMode_WORKFLOW {
		panic(httpexceptions.NewNotWorkflowAppError())
	}
	workflow_run_id := c.Param("workflow_run_id")
	if workflow_run_id == "" {
		mlog.Error("invalid workflow_run_id")
		response.InvalidArgError(c)
		return
	}
	path := c.Request.URL.Path
	in := pbapi.GetWorkflowRunDetailRequest{
		AppId:           app.ID,
		WorkflowRunId:   workflow_run_id,
		IsNodeExecution: !strings.HasSuffix(path, "/workflows/run/"),
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowRunDetail(c, &in)
		if err != nil {
			mlog.Errorf("remote call GetWorkflowRunDetail failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowRunDetail failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowRunDetail return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ResponseStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					if !in.IsNodeExecution {
						c.JSON(http.StatusOK, response.NewWorkflowRunDetailResponse(pbrsp.ResponseStr))
					} else {
						c.JSON(http.StatusOK, response.NewWorkflowRunNodeExecutionListResponse(pbrsp.ResponseStr))
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

func (api *WorkflowRunApi) AppRun(c *gin.Context) {
	rawapp, _ := c.Get("app_model")
	app := rawapp.(*models.App)

	path := c.Request.URL.Path
	if !(strings.HasSuffix(path, "/workflows/run") || strings.HasSuffix(path, "/chat-messages")) {
		mlog.Error("unsupported path=", path)
		response.InvalidArgErrorWithDetail(c, "unsupported path="+path)
		return
	}
	if strings.HasSuffix(path, "/workflows/run") && app.Mode != models.AppMode_WORKFLOW {
		panic(httpexceptions.NewNotWorkflowAppError())
	}
	if strings.HasSuffix(path, "/chat-messages") && !slices.Contains([]models.AppMode{models.AppMode_ADVANCED_CHAT, models.AppMode_CHAT, models.AppMode_AGENT_CHAT}, app.Mode) {
		panic(httpexceptions.NewNotChatAppError())
	}
	in := pbapi.AppRunRequest{AppId: app.ID}
	switch app.Mode {
	case models.AppMode_WORKFLOW:
		type Req struct {
			Inputs       map[string]any `json:"inputs"`
			ResponseMode string         `json:"response_mode"`
			User         string         `json:"user"`
			Files        []any          `json:"files"`
		}
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		if len(req.Inputs) == 0 {
			mlog.Error("missing inputs")
			response.InvalidArgError(c)
			return
		}
		if req.User == "" {
			mlog.Error("missing user")
			response.InvalidArgError(c)
			return
		}
		in.UserId = req.User
		in.ResponseMode = req.ResponseMode
		bindata, _ := json.Marshal(req)
		in.ArgsStr = string(bindata)
	case models.AppMode_ADVANCED_CHAT:
		fallthrough
	case models.AppMode_CHAT:
		fallthrough
	case models.AppMode_AGENT_CHAT:
		type Req struct {
			Inputs           map[string]any `json:"inputs"`
			Query            string         `json:"query"`
			ResponseMode     string         `json:"response_mode"`
			ConversationID   string         `json:"conversation_id"`
			RetrieverFrom    string         `json:"retriever_from"`
			AutoGenerateName bool           `json:"auto_generate_name"`
			User             string         `json:"user"`
		}
		req := Req{
			RetrieverFrom:    "dev",
			AutoGenerateName: true,
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		if len(req.Inputs) == 0 {
			mlog.Error("missing inputs")
			response.InvalidArgError(c)
			return
		}
		if req.User == "" {
			mlog.Error("missing user")
			response.InvalidArgError(c)
			return
		}
		if req.Query == "" {
			mlog.Error("query user")
			response.InvalidArgError(c)
			return
		}
		in.UserId = req.User
		in.ResponseMode = req.ResponseMode
		bindata, _ := json.Marshal(req)
		in.ArgsStr = string(bindata)
	default:
		mlog.Error("unsupported app mode=", app.Mode)
		response.InvalidArgErrorWithDetail(c, "unsupported app mode="+string(app.Mode))
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("app")
	if conn != nil {
		stream, err := pbapi.NewAppClient(conn).AppRun(c, &in)
		if err != nil {
			mlog.Errorf("remote call AppRun failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call AppRun failed",
			})
			return
		} else {
			for {
				pbrsp, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						mlog.Infof("remote call AppRun end")
						break
					} else {
						mlog.Errorf("remote call AppRun stream read failed:%v", err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"result": "fail",
							"data":   fmt.Sprintf("remote call AppRun stream read failed:%v", err),
						})
						break
					}
				} else {
					mlog.Infof("remote call AppRun stream return=%#v", pbrsp)
					if pbrsp.Exp != nil {
						response.PbHttpException(c, pbrsp.Exp)
						return
					} else {
						if pbrsp.DirectReplyDictStr != "" {
							tmp_dict := map[string]any{}
							if err1 := json.Unmarshal([]byte(pbrsp.DirectReplyDictStr), &tmp_dict); err1 != nil {
								mlog.Errorf("json unmarshal direct_reply_dict_str=%s to dict failed:%v", pbrsp.DirectReplyDictStr, err)
								c.JSON(http.StatusInternalServerError, gin.H{
									"result": "fail",
									"data":   fmt.Sprintf("json unmarshal direct_reply_dict_str=%s to dict failed:%v", pbrsp.DirectReplyDictStr, err),
								})
							} else {
								c.JSON(http.StatusOK, tmp_dict)
							}
							break
						} else {
							mlog.Debugf("-----send:%s", pbrsp.StreamReplyStr)
							c.Data(http.StatusOK, "text/event-stream", []byte(pbrsp.StreamReplyStr))
						}
					}
				}
			}
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
