package api

import (
	"fmt"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/extension"
	"mlib.com/gofy/server/core/external_data_tool/base"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/global"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/gofy/server/utils/mapstruct"
)

type ApiExternalDataTool struct {
	*base.ExternalDataTool
}

func init() {
	global.RegisgterExtension(&ApiExternalDataTool{
		ExternalDataTool: &base.ExternalDataTool{
			Extensible: &extension.Extensible{
				Config: make(map[string]any),
			},
		},
	}, nil, nil, false, 0)
}

func (aedt *ApiExternalDataTool) Name() string {
	return "api"
}

func (aedt *ApiExternalDataTool) ValidateConfig(tenant_id string, config map[string]any) {
	/*
	   Validate the incoming form config data.

	   :param tenant_id: the id of workspace
	   :param config: the form config data
	   :return:
	*/
	// own validation logic
	api_based_extension_id := mapstruct.Get(config, "api_based_extension_id", "")
	if api_based_extension_id == "" {
		panic(exceptions.NewValueError("api_based_extension_id is required"))
	}
	// get api_based_extension
	api_based_extension := new(models.APIBasedExtension)
	err := dbengine.Instance().DB.Model(&models.APIBasedExtension{}).Where("tenant_id = ? and id = ?", tenant_id, api_based_extension_id).First(api_based_extension).Error
	if err != nil {
		mlog.Errorf("get APIBasedExtension failed:%v", err)
		api_based_extension = nil
	}

	if api_based_extension == nil {
		panic(exceptions.NewValueError("api_based_extension_id is invalid"))
	}
}
func (aedt *ApiExternalDataTool) Query(inputs map[string]any, query string) string {
	/*
	   Query the external data tool.

	   :param inputs: user inputs
	   :param query: the query of chat app
	   :return: the tool query result
	*/
	// get params from config
	if aedt.Config == nil {
		mlog.Error("config is required")
		panic(exceptions.NewValueError("config is required"))
	}
	api_based_extension_id := mapstruct.Get(aedt.Config, "api_based_extension_id", "")
	if api_based_extension_id == "" {
		mlog.Error("api_based_extension_id is required")
		panic(exceptions.NewValueError("api_based_extension_id is required"))
	}
	// get api_based_extension
	api_based_extension := new(models.APIBasedExtension)
	err := dbengine.Instance().DB.Model(&models.APIBasedExtension{}).Where("tenant_id = ? and id = ?", aedt.TenantID, api_based_extension_id).First(api_based_extension).Error
	if err != nil {
		mlog.Errorf("get APIBasedExtension failed:%v", err)
		api_based_extension = nil
	}

	if api_based_extension == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("[External data tool] API query failed, variable: {%v}, error: api_based_extension_id is invalid", aedt.Variable)))
	}
	// decrypt api_key
	// api_key = encrypter.decrypt_token(tenant_id=aedt.tenant_id, token=api_based_extension.api_key)
	api_key := api_based_extension.APIKey
	response_json := func() map[string]any {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				if real_exp, ok := r.(error); ok {
					panic(exceptions.NewValueError(fmt.Sprintf("[External data tool] API query failed, variable: {%s}, error: {%v}", aedt.Variable, real_exp)))
				} else {
					panic(r)
				}
			}
		}()
		// request api
		return extension.NewAPIBasedExtensionRequestor(api_based_extension.APIEndpoint, api_key).Request(
			models.APIBasedExtensionPoint_APP_EXTERNAL_DATA_TOOL_QUERY,
			map[string]any{"app_id": aedt.AppID, "tool_variable": aedt.Variable, "inputs": inputs, "query": query},
		)
	}()

	if _, ok := response_json["result"]; !ok {
		panic(exceptions.NewValueError(fmt.Sprintf("[External data tool] API query failed, variable: {%s}, error: result not found in response", aedt.Variable)))
	}
	if _, ok := response_json["result"].(string); !ok {
		panic(exceptions.NewValueError(fmt.Sprintf("[External data tool] API query failed, variable: {%s}, error: result is not string", aedt.Variable)))
	}

	return response_json["result"].(string)
}
