package modelmanager

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	providermanager "github.com/odysseythink/gofy/backend/core/manageres/provider_manager"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

type ModelManager struct {
	*providermanager.ProviderManager
}

func NewModelManager() *ModelManager {
	return &ModelManager{
		ProviderManager: &providermanager.ProviderManager{},
	}
}

func (mgr *ModelManager) GetDefaultProviderModelName(tenant_id string, model_type modelruntimeenumtypes.ModelType) (string, string) {
	/*
	   Return first provider and the first model in the provider
	   :param tenant_id: tenant id
	   :param model_type: model type
	   :return: provider name, model name
	*/
	return mgr.ProviderManager.GetFirstProviderFirstModel(tenant_id, model_type)
}

func (mgr *ModelManager) GetDefaultModelInstance(tenant_id string, model_type modelruntimeenumtypes.ModelType) *ModelInstance {
	/*
	   Get default model instance
	   :param tenant_id: tenant id
	   :param model_type: model type
	   :return:
	*/
	default_model_entity := mgr.ProviderManager.GetDefaultModel(tenant_id, model_type)

	if default_model_entity == nil {
		panic(exceptions.NewProviderTokenNotInitError(fmt.Sprintf("Default model not found for %v", model_type)))
	}
	return mgr.GetModelInstance(
		tenant_id,
		default_model_entity.Provider.Provider,
		model_type,
		default_model_entity.Model,
	)
}
func (mgr *ModelManager) GetModelInstance(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string) *ModelInstance {
	/*
	   Get model instance
	   :param tenant_id: tenant id
	   :param provider: provider name
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	if provider == "" {
		default_model_entity := mgr.ProviderManager.GetDefaultModel(tenant_id, model_type)

		if default_model_entity == nil {
			panic(exceptions.NewProviderTokenNotInitError(fmt.Sprintf("Default model not found for %v", model_type)))
		}
		provider = default_model_entity.Provider.Provider
		// return mgr.get_default_model_instance(tenant_id, model_type)
		if provider == "" {
			panic(exceptions.NewProviderTokenNotInitError("can't find provider"))
		}
	}

	provider_model_bundle := mgr.ProviderManager.GetProviderModelBundle(
		tenant_id, provider, model_type,
	)

	return NewModelInstance(provider_model_bundle, model)
}
