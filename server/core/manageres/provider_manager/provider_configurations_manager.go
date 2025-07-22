package providermanager

import (
	coreentities "mlib.com/gofy/server/entities/core"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

type ProviderConfigurationsManager struct {
}

func (mgr *ProviderConfigurationsManager) GetModels(
	pcs *coreentities.ProviderConfigurations,
	provider string,
	model_type modelruntimeentities.ModelType,
	only_active bool,
) []*coreentities.ModelWithProviderEntity {
	/*
		Get available models.

		If preferred provider type is `system`:
		  Get the current **system mode** if provider supported,
		  if all system modes are not available (no quota), it is considered to be the **custom credential mode**.
		  If there is no model configured in custom mode, it is treated as no_configure.
		system > custom > no_configure

		If preferred provider type is `custom`:
		  If custom credentials are configured, it is treated as custom mode.
		  Otherwise, get the current **system mode** if supported,
		  If all system modes are not available (no quota), it is treated as no_configure.
		custom > system > no_configure

		If real mode is `system`, use system credentials to get models,
		  paid quotas > provider free quotas > system free quotas
		  include pre-defined models (exclude GPT-4, status marked as `no_permission`).
		If real mode is `custom`, use workspace custom credentials to get models,
		  include pre-defined models, custom models(manual append).
		If real mode is `no_configure`, only return pre-defined models from `model runtime`.
		  (model status marked as `no_configure` if preferred provider type is `custom` otherwise `quota_exceeded`)
		model status marked as `active` is available.

		:param provider: provider name
		:param model_type: model type
		:param only_active: only active models
		:return:
	*/
	all_models := []*coreentities.ModelWithProviderEntity{}
	for _, provider_configuration := range pcs.Configurations {
		if provider != "" && provider_configuration.Provider.Provider != provider {
			continue
		}
		provider_models := (&ProviderConfigurationManager{}).GetProviderModels(provider_configuration, model_type, only_active)

		all_models = append(all_models, provider_models...)
	}
	return all_models
}
