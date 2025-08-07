package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/cluster"
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	servicesentities "mlib.com/gofy/server/entities/services"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/mlog"
)

type ModelsApi struct {
}

func (api *ModelsApi) DefaultModel(c *gin.Context) {
	model_type := c.Query("model_type")
	if !modelruntimeenumtypes.ModelType(model_type).Valid() {
		mlog.Errorf("invalid model type=%s", model_type)
		response.InvalidArgError(c)
		return
	}
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetDefaultModel(c, &pbapi.GetDefaultModelRequest{TenantId: acc.CurrentTenantID(), ModelType: model_type})
		if err != nil {
			mlog.Errorf("remote call GetDefaultModel failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetDefaultModel failed",
			})
			return
		} else {
			mlog.Infof("remote call GetDefaultModel return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.DefaultModelResponseStr == "" {
					mlog.Error("get default model failed")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					c.JSON(http.StatusOK, gin.H{"data": servicesentities.NewDefaultModelResponse(pbrsp.DefaultModelResponseStr)})
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
func (api *ModelsApi) SetDefaultModel(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	in := pbapi.SetDefaultModelRequest{TenantId: acc.CurrentTenantID()}
	err := c.ShouldBindJSON(&in)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	for _, model_setting := range in.ModelSettings {
		if !modelruntimeenumtypes.ModelType(model_setting.ModelType).Valid() {
			c.Set("http_exceptions", exceptions.NewValueError("invalid model type"))
			return
		}
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).SetDefaultModel(c, &in)
		if err != nil {
			mlog.Errorf("remote call SetDefaultModel failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call SetDefaultModel failed",
			})
			return
		} else {
			mlog.Infof("remote call SetDefaultModel return:%#v", pbrsp)
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
func (api *ModelsApi) Available(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	model_type := c.Param("model_type")
	if !modelruntimeenumtypes.ModelType(model_type).Valid() {
		mlog.Errorf("model_type=%s is invalid", model_type)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("model_type=%s is invalid", model_type))
		return
	}

	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetAvailableModelProvider(c, &pbapi.GetAvailableModelProviderRequest{TenantId: acc.CurrentTenantID(), ModelType: model_type})
		if err != nil {
			mlog.Errorf("remote call GetAvailableModelProvider failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetAvailableModelProvider failed",
			})
			return
		} else {
			mlog.Infof("remote call GetAvailableModelProvider return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ModelsStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					var tmps []*servicesentities.ProviderWithModelsResponse
					if err1 := json.Unmarshal([]byte(pbrsp.ModelsStr), &tmps); err1 != nil {
						mlog.Errorf("json unmarshal models_list_str=%s to ProviderWithModelsResponse list failed:%v", pbrsp.ModelsStr, err1)
						c.JSON(http.StatusBadRequest, gin.H{"result": "no providerr list", "code": 7})
					} else {
						c.JSON(http.StatusOK, gin.H{"data": tmps, "result": "success", "code": 0})
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

func (api *ModelsApi) GetModel(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetModelProviderModel(c, &pbapi.GetModelProviderModelRequest{TenantId: tenant_id, Provider: provider})
		if err != nil {
			mlog.Errorf("remote call GetModelProviderModel failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetModelProviderModel failed",
			})
			return
		} else {
			mlog.Infof("remote call GetModelProviderModel return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ModelsStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					var tmps []*servicesentities.ModelWithProviderEntityResponse
					if err1 := json.Unmarshal([]byte(pbrsp.ModelsStr), &tmps); err1 != nil {
						mlog.Errorf("json unmarshal models_list_str=%s to ModelWithProviderEntityResponse list failed:%v", pbrsp.ModelsStr, err1)
						c.JSON(http.StatusBadRequest, gin.H{"result": "no models list", "code": 7})
					} else {
						c.JSON(http.StatusOK, gin.H{"data": tmps, "result": "success", "code": 0})
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

func (api *ModelsApi) SetModel(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	type Req struct {
		Model         string                          `json:"model"`
		ModelType     modelruntimeenumtypes.ModelType `json:"model_type"`
		Credentials   map[string]any                  `json:"credentials"`
		LoadBalancing map[string]any                  `json:"load_balancing"`
		ConfigFrom    string                          `json:"config_from"`
	}
	var req Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	if req.Model == "" {
		mlog.Error("missing model")
		response.InvalidArgErrorWithDetail(c, "missing model")
		return
	}
	if !req.ModelType.Valid() {
		mlog.Errorf("model_type=%s is invalid", req.ModelType)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("model_type=%s is invalid", req.ModelType))
		return
	}
	in := pbapi.SetModelProviderModelRequest{TenantId: tenant_id, Provider: provider, Model: req.Model, ModelType: string(req.ModelType), ConfigFrom: req.ConfigFrom}
	if len(req.Credentials) > 0 {
		bindata, _ := json.Marshal(req.Credentials)
		in.CredentialsDictStr = string(bindata)
	}
	if len(req.LoadBalancing) > 0 {
		bindata, _ := json.Marshal(req.LoadBalancing)
		in.LoadBalancingDictStr = string(bindata)
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).SetModelProviderModel(c, &in)
		if err != nil {
			mlog.Errorf("remote call SetModelProviderModel failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call SetModelProviderModel failed",
			})
			return
		} else {
			mlog.Infof("remote call SetModelProviderModel return:%#v", pbrsp)
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

func (api *ModelsApi) EnableModel(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	type Req struct {
		Model     string                          `json:"model"`
		ModelType modelruntimeenumtypes.ModelType `json:"model_type"`
	}
	var req Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	if req.Model == "" {
		mlog.Error("missing model")
		response.InvalidArgErrorWithDetail(c, "missing model")
		return
	}
	if !req.ModelType.Valid() {
		mlog.Errorf("model_type=%s is invalid", req.ModelType)
		response.InvalidArgErrorWithDetail(c, fmt.Sprintf("model_type=%s is invalid", req.ModelType))
		return
	}

	in := pbapi.EnableModelProviderModelRequest{TenantId: tenant_id, Provider: provider, Model: req.Model, ModelType: string(req.ModelType)}
	method := strings.Split(c.Request.URL.Path, provider)[1]
	switch method {
	case "/models/enable":
		in.Enable = true
	case "/models/disable":
		in.Enable = false
	default:
		mlog.Error("unsupported request url path=", c.Request.URL.Path)
		response.InvalidArgErrorWithDetail(c, "unsupported request url path="+c.Request.URL.Path)
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).EnableModelProviderModel(c, &in)
		if err != nil {
			mlog.Errorf("remote call EnableModelProviderModel failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call EnableModelProviderModel failed",
			})
			return
		} else {
			mlog.Infof("remote call EnableModelProviderModel return:%#v", pbrsp)
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
func (api *ModelsApi) GetParameterRules(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	model := c.Query("model")
	if model == "" {
		mlog.Error("missing model")
		response.InvalidArgErrorWithDetail(c, "missing model")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetModelParameterRules(c, &pbapi.GetModelParameterRulesRequest{TenantId: tenant_id, Provider: provider, Model: model})
		if err != nil {
			mlog.Errorf("remote call GetModelParameterRules failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetModelParameterRules failed",
			})
			return
		} else {
			mlog.Infof("remote call GetModelParameterRules return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ParameterRulesStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					var params []*modelruntimeentities.ParameterRule
					if err1 := json.Unmarshal([]byte(pbrsp.ParameterRulesStr), &params); err1 != nil {
						mlog.Errorf("json unmarshal parameter_rules_str=%s to ParameterRule list failed:%v", pbrsp.ParameterRulesStr, err1)
						c.JSON(http.StatusBadRequest, gin.H{"result": "no ParameterRules", "code": 7})
					} else {
						c.JSON(http.StatusOK, gin.H{"data": params})
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
