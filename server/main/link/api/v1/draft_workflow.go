package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type DraftWorkflowApi struct {
}

func (api *DraftWorkflowApi) Find(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraft(c, &pbapi.GetWorkflowDraftRequest{AccountId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraft failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowDraft failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraft return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.WorkflowResponseStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewWorkflowResponse(pbrsp.WorkflowResponseStr))
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

func (api *DraftWorkflowApi) DefaultBlockConfigs(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDefaultBlockConfigs(c, &pbapi.GetWorkflowDefaultBlockConfigsRequest{AccountId: acc.ID})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDefaultBlockConfigs failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowDefaultBlockConfigs failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDefaultBlockConfigs return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				cfgs := []map[string]any{}
				if pbrsp.CfgsStr != "" {
					err := json.Unmarshal([]byte(pbrsp.CfgsStr), &cfgs)
					if err != nil {
						mlog.Errorf("json unmarshal=%s to object list failed:%v", pbrsp.CfgsStr, err)
					}
				}
				c.JSON(http.StatusOK, cfgs)
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
func (api *DraftWorkflowApi) DefaultBlockConfig(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	block_type := c.Param("block_type")
	if block_type == "" {
		mlog.Error("invalid block_type")
		response.InvalidArgError(c)
		return
	}
	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	filter_dict_str := c.Query("q")
	if filter_dict_str != "" {
		tmp := map[string]any{}
		err := json.Unmarshal([]byte(filter_dict_str), &tmp)
		if err != nil {
			mlog.Errorf("q=%s arg must be dict string:%v", filter_dict_str, err)
			response.InvalidArgErrorWithDetail(c, fmt.Sprintf("q=%s arg must be dict string:%v", filter_dict_str, err))
			return
		}
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDefaultBlockConfig(c, &pbapi.GetWorkflowDefaultBlockConfigRequest{AccountId: acc.ID, AppId: app_id, BlockType: block_type, FilterDictStr: filter_dict_str})
		if err != nil {
			mlog.Errorf("remote call filter_dict_str failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call filter_dict_str failed",
			})
			return
		} else {
			mlog.Infof("remote call filter_dict_str return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				cfg := []map[string]any{}
				if pbrsp.CfgStr != "" {
					err := json.Unmarshal([]byte(pbrsp.CfgStr), &cfg)
					if err != nil {
						mlog.Errorf("json unmarshal=%s to object list failed:%v", pbrsp.CfgStr, err)
					}
				}
				c.JSON(http.StatusOK, cfg)
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

func (api *DraftWorkflowApi) Sync(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}

	content_type := c.GetHeader("Content-Type")
	type Request struct {
		Graph                 map[string]any   `json:"graph"`
		Features              map[string]any   `json:"features"`
		Hash                  string           `json:"hash"`
		EnvironmentVariables  []map[string]any `json:"environment_variables"`
		ConversationVariables []map[string]any `json:"conversation_variables"`
	}
	req := &Request{}
	if strings.Contains(content_type, "application/json") {
		err := c.ShouldBindJSON(req)
		if err != nil {
			mlog.Errorf("can't bind args: %v", err)
			response.InvalidArgError(c)
			return
		}
	} else if strings.Contains(content_type, "text/plain") {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			mlog.Errorf("read request body failed:%v", err)
			response.InvalidArgError(c)
			return
		}
		err = json.Unmarshal(body, req)
		if err != nil {
			mlog.Errorf("request body(%s) has invalid format:%v", string(body), err)
			response.InvalidArgError(c)
			return
		}
	} else {
		mlog.Errorf("unsurported Content-Type=%s", content_type)
		c.AbortWithStatus(http.StatusUnsupportedMediaType)
		return
	}
	in := pbapi.WorkflowSyncDraftRequest{AccountId: acc.ID, AppId: app_id, Hash: req.Hash}
	bindata, _ := json.Marshal(req.Graph)
	in.GraphStr = string(bindata)
	bindata, _ = json.Marshal(req.Features)
	in.FeaturesStr = string(bindata)
	bindata, _ = json.Marshal(req.EnvironmentVariables)
	in.EnvironmentVariablesStr = string(bindata)
	bindata, _ = json.Marshal(req.ConversationVariables)
	in.ConversationVariablesStr = string(bindata)
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).WorkflowSyncDraft(c, &in)
		if err != nil {
			mlog.Errorf("remote call WorkflowSyncDraft failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call WorkflowSyncDraft failed",
			})
			return
		} else {
			mlog.Infof("remote call WorkflowSyncDraft return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"result": "success", "hash": pbrsp.UniqueHash, "updated_at": pbrsp.UpdatedAt, "code": 0})
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

func (api *DraftWorkflowApi) Config(c *gin.Context) {
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowConfig(c, &pbapi.GetWorkflowConfigRequest{})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowConfig failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowConfig failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowConfig return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"parallel_depth_limit": pbrsp.ParallelDepthLimit})
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

func (api *DraftWorkflowApi) Run(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	method := strings.Split(c.Request.URL.Path, app_id)[1]
	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	var err error
	args_dict_str := ""
	if method == "/workflows/draft/run" {
		var req struct {
			Inputs map[string]any `json:"inputs"`
			File   []any          `json:"files"`
		}
		err = c.ShouldBindJSON(&req)
		if err != nil {
			mlog.Errorf("can't bind args: %v", err)
			response.InvalidArgError(c)
			return
		}
		bindata, _ := json.Marshal(req)
		args_dict_str = string(bindata)

	} else if method == "/advanced-chat/workflows/draft/run" {
		var req struct {
			Inputs          map[string]any `json:"inputs"`
			Query           string         `json:"query"`
			ConversationID  string         `json:"conversation_id"`
			ParentMessageID string         `json:"parent_message_id"`
		}
		err = c.ShouldBindJSON(&req)
		if err != nil {
			mlog.Errorf("can't bind args: %v", err)
			response.InvalidArgErrorWithDetail(c, "can't bind args: "+err.Error())
			return
		}
		if req.Query == "" {
			mlog.Errorf("missing query arg")
			response.InvalidArgErrorWithDetail(c, "missing query arg")
			return
		}
		bindata, _ := json.Marshal(req)
		args_dict_str = string(bindata)
	} else {
		mlog.Error("unsupported method=", method)
		response.InvalidArgErrorWithDetail(c, "unsupported method="+method)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("app")
	if conn != nil {
		stream, err := pbapi.NewAppClient(conn).AppRun(c, &pbapi.AppRunRequest{UserId: acc.ID, AppId: app_id, ArgsStr: args_dict_str, ResponseMode: "streaming"})
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
func (api *DraftWorkflowApi) NodeRun(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	node_id := c.Param("node_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	if node_id == "" {
		mlog.Error("invalid node_id")
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	type Req struct {
		Inputs map[string]any `json:"inputs"`
	}
	var req Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		mlog.Errorf("can't bind args: %v", err)
		response.InvalidArgError(c)
		return
	}
	inputs_dict_str := ""
	if len(req.Inputs) > 0 {
		bindata, _ := json.Marshal(req.Inputs)
		inputs_dict_str = string(bindata)
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).WorkflowNodeRun(c, &pbapi.WorkflowNodeRunRequest{AccountId: acc.ID, AppId: app_id, NodeId: node_id, InputsDictStr: inputs_dict_str})
		if err != nil {
			mlog.Errorf("remote call WorkflowNodeRun failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call WorkflowNodeRun failed",
			})
			return
		} else {
			mlog.Infof("remote call WorkflowNodeRun return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.WorkflowRunNodeExecutionResponseStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewWorkflowRunNodeExecutionResponse(pbrsp.WorkflowRunNodeExecutionResponseStr))
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

func (api *DraftWorkflowApi) GetPublished(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowPublished(c, &pbapi.GetWorkflowPublishedRequest{AccountId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowPublished failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowPublished failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowPublished return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.WorkflowResponseStr == "" {
					c.JSON(http.StatusOK, nil)
				} else {
					c.JSON(http.StatusOK, response.NewWorkflowResponse(pbrsp.WorkflowResponseStr))
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

func (api *DraftWorkflowApi) SetPublished(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).WorkflowPublished(c, &pbapi.WorkflowPublishedRequest{AccountId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call WorkflowPublished failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call WorkflowPublished failed",
			})
			return
		} else {
			mlog.Infof("remote call WorkflowPublished return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{"result": "success", "created_at": pbrsp.CreatedAt})
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
func (api *DraftWorkflowApi) GetVariable(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	variableID := c.Param("variable_id")
	if variableID == "" {
		mlog.Error("invalid variable_id")
		response.InvalidArgError(c)
		return
	}
	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraftVariable(c, &pbapi.GetWorkflowDraftVariableRequest{UserId: acc.ID, AppId: app_id, VariableId: variableID})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraftVariable failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    "remote call GetWorkflowDraftVariable failed",
				"code":    "internal_server_error",
				"status":  500,
				"items":   nil,
				"total":   0,
				"message": "remote call GetWorkflowDraftVariable failed.",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraftVariable return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				var varObj map[string]any
				err := json.Unmarshal([]byte(pbrsp.VarStr), &varObj)
				if err != nil {
					mlog.Errorf("json unmarshal failed:%v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "json unmarshal failed",
					})
					return
				}

				c.JSON(http.StatusOK, varObj)
				return
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
func (api *DraftWorkflowApi) UpdateVariable(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraft(c, &pbapi.GetWorkflowDraftRequest{AccountId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraft failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetWorkflowDraft failed",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraft return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.WorkflowResponseStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, response.NewWorkflowResponse(pbrsp.WorkflowResponseStr))
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

func (api *DraftWorkflowApi) ListVariable(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		mlog.Errorf("invalid page: %v", err)
		response.InvalidArgError(c)
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		mlog.Errorf("invalid limit: %v", err)
		response.InvalidArgError(c)
		return
	}

	// # The role of the current user in the ta table must be admin, owner, or editor
	if !acc.IsEditor() {
		mlog.Errorf("forbiden")
		response.Forbidden(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraftVariableList(c, &pbapi.GetWorkflowDraftVariableListRequest{UserId: acc.ID, AppId: app_id, Page: int32(page), Limit: int32(limit)})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraftVariableList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    "remote call GetWorkflowDraftVariableList failed",
				"code":    "internal_server_error",
				"status":  500,
				"items":   nil,
				"total":   0,
				"message": "remote call GetWorkflowDraftVariableList failed.",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraftVariableList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"items": pbrsp.Items,
					"total": pbrsp.Total,
				})
				return
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

func (api *DraftWorkflowApi) ListSysVariable(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		mlog.Errorf("invalid page: %v", err)
		response.InvalidArgError(c)
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		mlog.Errorf("invalid limit: %v", err)
		response.InvalidArgError(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraftSysVariableList(c, &pbapi.GetWorkflowDraftVariableListRequest{UserId: acc.ID, AppId: app_id, Page: int32(page), Limit: int32(limit)})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraftVariableList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    "remote call GetWorkflowDraftVariableList failed",
				"code":    "internal_server_error",
				"status":  500,
				"items":   nil,
				"total":   0,
				"message": "remote call GetWorkflowDraftVariableList failed.",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraftVariableList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				varObj := []map[string]any{}
				err := json.Unmarshal([]byte(pbrsp.ItemsStr), &varObj)
				if err != nil {
					mlog.Errorf("json unmarshal failed:%v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "json unmarshal failed",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"items": varObj,
				})
				return
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

func (api *DraftWorkflowApi) ListConversationVariable(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraftConversationVariableList(c, &pbapi.GetWorkflowDraftVariableListRequest{UserId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraftConversationVariableList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    "remote call GetWorkflowDraftConversationVariableList failed",
				"code":    "internal_server_error",
				"status":  500,
				"items":   nil,
				"total":   0,
				"message": "remote call GetWorkflowDraftConversationVariableList failed.",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraftConversationVariableList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				varObj := []map[string]any{}
				err := json.Unmarshal([]byte(pbrsp.ItemsStr), &varObj)
				if err != nil {
					mlog.Errorf("json unmarshal failed:%v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "json unmarshal failed",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"items": varObj,
				})
				return
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

func (api *DraftWorkflowApi) ListEnvironmentVariable(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	app_id := c.Param("app_id")
	if app_id == "" {
		mlog.Error("invalid app_id")
		response.InvalidArgError(c)
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetWorkflowDraftEnvVariableList(c, &pbapi.GetWorkflowDraftVariableListRequest{UserId: acc.ID, AppId: app_id})
		if err != nil {
			mlog.Errorf("remote call GetWorkflowDraftEnvVariableList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    "remote call GetWorkflowDraftEnvVariableList failed",
				"code":    "internal_server_error",
				"status":  500,
				"items":   nil,
				"total":   0,
				"message": "remote call GetWorkflowDraftEnvVariableList failed.",
			})
			return
		} else {
			mlog.Infof("remote call GetWorkflowDraftEnvVariableList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				varObj := []map[string]any{}
				err := json.Unmarshal([]byte(pbrsp.ItemsStr), &varObj)
				if err != nil {
					mlog.Errorf("json unmarshal failed:%v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "json unmarshal failed",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"items": varObj,
				})
				return
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
