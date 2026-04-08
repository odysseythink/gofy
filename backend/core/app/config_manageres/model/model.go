package model

import (
	"fmt"
	"slices"
	"strings"

	"mlib.com/gofy/server/core/exceptions"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	modelproviders "mlib.com/gofy/server/core/model_runtime/model_provides"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	pluginentities "mlib.com/gofy/server/entities/plugin"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/mlog"
)

type ModelConfigManager struct {
}

func (mgr *ModelConfigManager) Convert(config map[string]any) appconfigentities.ModelConfigEntity {
	/*
	   Convert model config to model config

	   :param config: model config args
	*/
	//model config
	var model_config map[string]any
	if _, ok := config["model"]; ok {
		if _, ok := config["model"].(map[string]any); ok {
			model_config = config["model"].(map[string]any)
		}
	}

	if len(model_config) == 0 {
		panic(exceptions.NewValueError("model is required"))
	}
	var completion_params map[string]any
	if _, ok := model_config["completion_params"]; ok {
		if _, ok := model_config["completion_params"].(map[string]any); ok {
			completion_params = model_config["completion_params"].(map[string]any)
		}
	}
	provider := ""
	if _, ok := model_config["provider"]; ok {
		if _, ok := model_config["provider"].(string); ok {
			provider = model_config["provider"].(string)
		}
	}
	model_name := ""
	if _, ok := model_config["name"]; ok {
		if _, ok := model_config["name"].(string); ok {
			model_name = model_config["name"].(string)
		}
	}
	model_mode := ""
	if _, ok := model_config["mode"]; ok {
		if _, ok := model_config["mode"].(string); ok {
			model_mode = model_config["mode"].(string)
		}
	}
	stop := []string{}
	if _, ok := completion_params["stop"]; ok {
		if _, ok := completion_params["stop"].([]string); ok {
			stop = completion_params["stop"].([]string)
			delete(completion_params, "stop")
		} else if _, ok := completion_params["stop"].([]any); ok {
			for _, v := range completion_params["stop"].([]any) {
				if _, ok := v.(string); ok {
					stop = append(stop, v.(string))
				} else {
					mlog.Errorf("completion_params=%#v stop have invalid format", completion_params)
					panic(exceptions.NewValueError("stop has invalid format"))
				}
			}
			delete(completion_params, "stop")
		}
	}

	return appconfigentities.ModelConfigEntity{
		Provider:   provider,
		Model:      model_name,
		Mode:       model_mode,
		Parameters: completion_params,
		Stop:       stop,
	}

}
func (mgr *ModelConfigManager) ValidateAndSetDefaults(tenant_id string, config map[string]any) (map[string]any, []string) {
	/*
	   Validate and set defaults for model config

	   :param tenant_id: tenant id
	   :param config: app model config args
	*/
	if _, ok := config["model"]; !ok {
		panic(exceptions.NewValueError("model is required"))
	}
	if _, ok := config["model"].(map[string]any); !ok {
		panic(exceptions.NewValueError("model must be of object type"))
	}
	model_config_dict := config["model"].(map[string]any)

	//model.provider
	provider_entities := (&modelproviders.ModelProviderFactory{}).GetProviders()
	model_provider_names := []string{}
	for _, provider := range provider_entities {
		model_provider_names = append(model_provider_names, provider.Provider)
	}
	if _, ok := model_config_dict["provider"]; !ok {
		panic(exceptions.NewValueError(fmt.Sprintf("model.provider is required and must be in {%v}", model_provider_names)))
	}
	if _, ok := model_config_dict["provider"].(string); !ok {
		panic(exceptions.NewValueError("model.provider is required and must be string"))
	}
	if !strings.Contains(model_config_dict["provider"].(string), "/") {
		config["model"].(map[string]any)["provider"] = pluginentities.NewModelProviderID(model_config_dict["provider"].(string), false).String()
	}
	if !slices.Contains(model_provider_names, model_config_dict["provider"].(string)) {
		panic(exceptions.NewValueError(fmt.Sprintf("model.provider is required and must be in {%v}", model_provider_names)))
	}
	//model.name
	if _, ok := model_config_dict["name"]; !ok {
		panic(exceptions.NewValueError("model.name is required"))
	}
	if _, ok := model_config_dict["name"].(string); !ok {
		panic(exceptions.NewValueError("model.name must be string"))
	}
	models := (&providermanager.ProviderConfigurationManager{}).GetModels(
		(&providermanager.ProviderManager{}).GetConfigurations(tenant_id),
		model_config_dict["provider"].(string),
		modelruntimeenumtypes.Model_LLM,
		false,
	)

	if len(models) == 0 {
		panic(exceptions.NewValueError("model.name must be in the specified model list"))
	}
	model_ids := []string{}
	for _, m := range models {
		model_ids = append(model_ids, m.Model)
	}

	if !slices.Contains(model_ids, model_config_dict["name"].(string)) {
		panic(exceptions.NewValueError("model.name must be in the specified model list"))
	}
	model_mode := ""
	for _, model := range models {
		if model.Model == model_config_dict["name"].(string) {
			model_mode = model.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_MODE].(string)
			break
		}
	}
	//model.mode
	if model_mode != "" {
		config["model"].(map[string]any)["mode"] = model_mode
	} else {
		config["model"].(map[string]any)["mode"] = "completion"
	}

	//model.completion_params
	if _, ok := model_config_dict["completion_params"]; !ok {
		panic(exceptions.NewValueError("model.completion_params is required"))
	}

	config["model"].(map[string]any)["completion_params"] = mgr.ValidateModelCompletionParams(
		model_config_dict["completion_params"],
	)

	return config, []string{"model"}

}
func (mgr *ModelConfigManager) ValidateModelCompletionParams(cp any) map[string]any {
	//model.completion_params
	if _, ok := cp.(map[string]any); !ok {
		panic(exceptions.NewValueError("model.completion_params must be of object type"))
	}
	//stop

	if _, ok := cp.(map[string]any)["stop"]; !ok {
		cp.(map[string]any)["stop"] = []string{}
	} else if _, ok := cp.(map[string]any)["stop"].([]any); !ok {
		if _, ok := cp.(map[string]any)["stop"].([]string); !ok {
			panic(exceptions.NewValueError("stop in model.completion_params must be of list type"))
		}
	}
	stop := []string{}
	if _, ok := cp.(map[string]any)["stop"].([]any); ok {
		for _, v := range cp.(map[string]any)["stop"].([]any) {
			if _, ok := v.(string); !ok {
				panic(exceptions.NewValueError("stop in model.completion_params must be of string list type"))
			}
			stop = append(stop, v.(string))
		}
	} else if _, ok := cp.(map[string]any)["stop"].([]string); ok {
		stop = cp.(map[string]any)["stop"].([]string)
	}
	if len(stop) > 4 {
		panic(exceptions.NewValueError("stop sequences must be less than 4"))
	}
	return cp.(map[string]any)
}
