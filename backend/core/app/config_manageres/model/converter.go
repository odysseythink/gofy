package model

import (
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	coreenumtypes "mlib.com/gofy/server/enum_types/core"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type ModelConfigConverter struct {
}

func (converter *ModelConfigConverter) Convert(app_config *appconfigentities.EasyUIBasedAppConfig) *appconfigentities.ModelConfigWithCredentialsEntity {
	/*
	   Convert app model config dict to entity.
	   :param app_config: app config
	   :raises ProviderTokenNotInitError: provider token not init error
	   :return: app orchestration config entity
	*/
	model_config := app_config.Model

	provider_model_bundle := (&providermanager.ProviderManager{}).GetProviderModelBundle(
		app_config.TenantID, model_config.Provider, modelruntimeenumtypes.Model_LLM,
	)

	provider_name := provider_model_bundle.Configuration.Provider.Provider
	model_name := model_config.Model

	model_type_instance := provider_model_bundle.ModelTypeInstance
	// model_type_instance = cast(LargeLanguageModel, model_type_instance)

	// check model credentials
	model_credentials := provider_model_bundle.Configuration.GetCurrentCredentials(
		modelruntimeenumtypes.Model_LLM, model_config.Model,
	)

	if len(model_credentials) == 0 {
		panic(exceptions.NewProviderTokenNotInitError(fmt.Sprintf("Model {%s} credentials is not initialized.", model_name)))
	}
	// check model
	provider_model := (&providermanager.ProviderConfigurationManager{}).GetProviderModel(
		provider_model_bundle.Configuration,
		modelruntimeenumtypes.Model_LLM,
		model_config.Model,
		false,
	)

	if provider_model == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Model {%s} not exist.", model_name)))
	}
	if provider_model.Status == coreenumtypes.ModelStatus_NO_CONFIGURE {
		panic(exceptions.NewProviderTokenNotInitError(fmt.Sprintf("Model {%s} credentials is not initialized.", model_name)))
	} else if provider_model.Status == coreenumtypes.ModelStatus_NO_PERMISSION {
		panic(exceptions.NewModelCurrentlyNotSupportError(fmt.Sprintf("Gofy Hosted OpenAI {%s} currently not support.", model_name)))
	} else if provider_model.Status == coreenumtypes.ModelStatus_QUOTA_EXCEEDED {
		panic(exceptions.NewQuotaExceededError(fmt.Sprintf("Model provider {%s} quota exceeded.", provider_name)))
	}
	// model config
	completion_params := model_config.Parameters
	stop := []string{}
	if _, ok := completion_params["stop"]; ok {
		if _, ok := completion_params["stop"].([]any); ok {
			for _, v := range completion_params["stop"].([]any) {
				if _, ok := v.(string); !ok {
					panic(exceptions.NewValueError("stop must be of string list type"))
				}
				stop = append(stop, v.(string))
			}
		} else if _, ok := completion_params["stop"].([]string); ok {
			stop = completion_params["stop"].([]string)
		}
		delete(completion_params, "stop")
	}
	model_schema := model_type_instance.GetModelSchema(model_type_instance, model_config.Model, model_credentials)

	// get model mode
	model_mode := model_config.Mode
	if model_mode == "" {
		model_mode = string(modelruntimeentities.LLMMode_CHAT)
		if model_schema != nil {
			if _, ok := model_schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_MODE]; ok {
				model_mode = model_schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_MODE].(string)
			}
		}
	}
	if model_schema == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Model {%s} not exist.", model_name)))
	}
	return &appconfigentities.ModelConfigWithCredentialsEntity{
		Provider:            model_config.Provider,
		Model:               model_config.Model,
		ModelSchema:         model_schema,
		Mode:                modelruntimeentities.LLMMode(model_mode),
		ProviderModelBundle: provider_model_bundle,
		Credentials:         model_credentials,
		Parameters:          completion_params,
		Stop:                stop,
	}
}
