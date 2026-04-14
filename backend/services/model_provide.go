package services

import (
	"fmt"
	"slices"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	modelentities "mlib.com/gofy/server/entities/model"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	servicesentities "mlib.com/gofy/server/entities/services"
	modelenumtypes "mlib.com/gofy/server/enum_types/model"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/models"
)

type ModelProvideService struct {
	provider_manager *providermanager.ProviderManager
}

func (s *ModelProvideService) GetModelsByProvider(tenant_id, provider string) []*servicesentities.ModelWithProviderEntityResponse {
	/*
		get provider models.
		For the model provider page,
		only supports passing in a single provider to query the list of supported models.

		:param tenant_id:
		:param provider:
		:return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	// get available models from provider_configurations
	available_models := (&providermanager.ProviderConfigurationsManager{}).GetModels(provider_configurations, provider, modelruntimeenumtypes.ModelType(""), false)

	res := []*servicesentities.ModelWithProviderEntityResponse{}
	for _, model := range available_models {
		res = append(res, &servicesentities.ModelWithProviderEntityResponse{
			ModelWithProviderEntity: model,
		})
	}

	// Get provider available models
	return res
}

func (s *ModelProvideService) GetModelsByModelType(tenant_id string, model_type modelruntimeenumtypes.ModelType) []*servicesentities.ProviderWithModelsResponse {
	/*
		get models by model type.

		:param tenant_id: workspace id
		:param model_type: model type
		:return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	// get available models from provider_configurations
	available_models := (&providermanager.ProviderConfigurationsManager{}).GetModels(provider_configurations, "", model_type, false)

	// Group models by provider
	provider_models := map[string][]*modelentities.ModelWithProviderEntity{}
	for _, model := range available_models {
		if _, ok := provider_models[model.Provider.Provider]; !ok {
			provider_models[model.Provider.Provider] = make([]*modelentities.ModelWithProviderEntity, 0)
		}
		if model.Deprecated {
			continue
		}
		if model.Status != modelenumtypes.ModelStatus_ACTIVE {
			continue
		}
		provider_models[model.Provider.Provider] = append(provider_models[model.Provider.Provider], model)
	}
	// convert to ProviderWithModelsResponse list
	providers_with_models := []*servicesentities.ProviderWithModelsResponse{}
	for provider, models := range provider_models {
		if len(models) == 0 {
			continue
		}
		first_model := models[0]

		rsp := &servicesentities.ProviderWithModelsResponse{
			Provider:  provider,
			Label:     first_model.Provider.Label,
			IconSmall: first_model.Provider.IconSmall,
			IconLarge: first_model.Provider.IconLarge,
			Status:    servicesentities.CustomConfigurationStatus_ACTIVE,
			Models:    make([]*modelentities.ProviderModelWithStatusEntity, 0),
			// Models:[]*modelentities.ProviderModelWithStatusEntity{
			//
			// 		ProviderModel: &modelruntimeentities.ProviderModel{

			// 		},

			// 		)
			// 		for model in models
			// 	},
		}
		for _, model := range models {
			rsp.Models = append(rsp.Models, &modelentities.ProviderModelWithStatusEntity{
				ProviderModel: &modelruntimeentities.ProviderModel{
					Model:           model.Model,
					Label:           model.Label,
					ModelType:       model.ModelType,
					Features:        model.Features,
					FetchFrom:       model.FetchFrom,
					ModelProperties: model.ModelProperties,
				},
				Status:               model.Status,
				LoadBalancingEnabled: model.LoadBalancingEnabled,
			})
		}
		providers_with_models = append(providers_with_models, rsp)
	}
	return providers_with_models
}

/*
get default model of model type.

:param tenant_id: workspace id
:param model_type: model type
:return:
*/
func (s *ModelProvideService) GetDefaultModelOfModelType(tenant_id string, model_type modelruntimeenumtypes.ModelType) *servicesentities.DefaultModelResponse {
	result := (&providermanager.ProviderManager{}).GetDefaultModel(tenant_id, model_type)
	if result == nil {
		return nil
	}
	return &servicesentities.DefaultModelResponse{
		Model:     result.Model,
		ModelType: result.ModelType,
		Provider: &servicesentities.SimpleProviderEntityResponse{
			SimpleProviderEntity: &modelruntimeentities.SimpleProviderEntity{
				Provider:            result.Provider.Provider,
				Label:               result.Provider.Label,
				IconSmall:           result.Provider.IconSmall,
				IconLarge:           result.Provider.IconLarge,
				SupportedModelTypes: result.Provider.SupportedModelTypes,
			},
		},
	}
}

func (s *ModelProvideService) GetProviderList(tenant_id string, model_type modelruntimeenumtypes.ModelType) []*servicesentities.ProviderResponse {
	/*
		get provider list.

		:param tenant_id: workspace id
		:param model_type: model type
		:return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)

	provider_responses := []*servicesentities.ProviderResponse{}
	for _, provider_configuration := range provider_configurations.Configurations {
		if model_type != "" {
			model_type_entity := modelruntimeenumtypes.ModelType(model_type)
			if !slices.Contains(provider_configuration.Provider.SupportedModelTypes, model_type_entity) {
				continue
			}
		}
		provider_response := &servicesentities.ProviderResponse{
			Provider:                 provider_configuration.Provider.Provider,
			Label:                    provider_configuration.Provider.Label,
			Description:              provider_configuration.Provider.Description,
			IconSmall:                provider_configuration.Provider.IconSmall,
			IconLarge:                provider_configuration.Provider.IconLarge,
			Background:               provider_configuration.Provider.Background,
			Help:                     provider_configuration.Provider.Help,
			SupportedModelTypes:      provider_configuration.Provider.SupportedModelTypes,
			ConfigurateMethods:       provider_configuration.Provider.ConfigurateMethods,
			ProviderCredentialSchema: provider_configuration.Provider.ProviderCredentialSchema,
			ModelCredentialSchema:    provider_configuration.Provider.ModelCredentialSchema,
			PreferredProviderType:    provider_configuration.PreferredProviderType,
			// CustomConfiguration: &servicesentities.CustomConfigurationResponse{
			// 	status=CustomConfigurationStatus.ACTIVE
			// 	if provider_configuration.is_custom_configuration_available()
			// 	else CustomConfigurationStatus.NO_CONFIGURE
			// },
			SystemConfiguration: &servicesentities.SystemConfigurationResponse{
				Enabled:             provider_configuration.SystemConfiguration.Enabled,
				CurrentQuotaType:    provider_configuration.SystemConfiguration.CurrentQuotaType,
				QuotaConfigurations: provider_configuration.SystemConfiguration.QuotaConfigurations,
			},
		}
		if provider_configuration.IsCustomConfigurationAvailable() {
			provider_response.CustomConfiguration = &servicesentities.CustomConfigurationResponse{Status: servicesentities.CustomConfigurationStatus_ACTIVE}
		} else {
			provider_response.CustomConfiguration = &servicesentities.CustomConfigurationResponse{Status: servicesentities.CustomConfigurationStatus_NO_CONFIGURE}
		}

		provider_responses = append(provider_responses, provider_response)
	}
	return provider_responses
}

func (s *ModelProvideService) SaveModelCredentials(
	tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string, credentials map[string]any,
) error {
	/*
	   save model credentials.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model_type: model type
	   :param model: model name
	   :param credentials: model credentials
	   :return:
	*/
	// Get all provider configurations of the current workspace

	// Add or update custom model credentials
	return ServiceGroupApp.ProviderConfiguration.AddOrUpdateCustomModelCredentials(
		tenant_id, provider, model_type, model, credentials,
	)
}

func (s *ModelProvideService) RemoveModelCredentials(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string) error {
	/*
	   remove model credentials.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/

	// Remove custom model credentials
	return ServiceGroupApp.ProviderConfiguration.DeleteCustomModelCredentials(tenant_id, provider, model_type, model)
}
func (s *ModelProvideService) GetProviderCredentials(tenant_id string, provider string) map[string]any {
	/*
	   get provider credentials.
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)

	// Get provider configuration
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider)))
	}

	provider_configuration := provider_configurations.Configurations[provider]

	return provider_configuration.GetCustomCredentials(true)
}

func (s *ModelProvideService) SaveProviderCredentials(tenant_id string, provider string, credentials map[string]any) error {
	/*
	   save custom provider config.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param credentials: provider credentials
	   :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)

	// Get provider configuration
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]
	// Add or update custom provider credentials.
	(&providermanager.ProviderConfigurationManager{}).AddOrUpdateCustomCredentials(provider_configuration, credentials)
	return nil
}

func (s *ModelProvideService) RemoveProviderCredentials(tenant_id string, provider string) error {
	/*
	   remove custom provider config.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :return:
	*/
	return ServiceGroupApp.ProviderConfiguration.DeleteCustomCredentialsByTenantAndProvider(tenant_id, provider)
}

func (s *ModelProvideService) UpdateDefaultModelOfModelType(tenant_id string, provider string, model string, model_type modelruntimeenumtypes.ModelType) *models.TenantDefaultModel {
	/*
	   update default model of model type.

	   :param tenant_id: workspace id
	   :param model_type: model type
	   :param provider: provider name
	   :param model: model name
	   :return:
	*/
	return (&providermanager.ProviderManager{}).UpdateDefaultModelRecord(
		tenant_id, provider, model, model_type,
	)
}

func (s *ModelProvideService) GetModelParameterRules(tenant_id string, provider string, model string) []*modelruntimeentities.ParameterRule {
	/*
				get model parameter rules.
		        Only supports LLM.

		        :param tenant_id: workspace id
		        :param provider: provider name
		        :param model: model name
		        :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)

	// Get provider configuration
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider)))
	}

	provider_configuration := provider_configurations.Configurations[provider]

	// Get model instance of LLM
	tmp_model_type_instance := (&providermanager.ProviderConfigurationManager{}).GetModelTypeInstance(provider_configuration, modelruntimeenumtypes.Model_LLM)

	model_type_instance := tmp_model_type_instance.(modelruntimeentities.LargeLanguageModeler)

	// fetch credentials
	credentials := provider_configuration.GetCurrentCredentials(modelruntimeenumtypes.Model_LLM, model)

	if len(credentials) == 0 {
		return nil
	}

	// Call get_parameter_rules method of model instance to get model parameter rules
	return model_type_instance.GetParameterRules(model_type_instance, model, credentials)
}

func (s *ModelProvideService) EnableModel(tenant_id string, provider string, model string, model_type string) {
	/*
	   enable model.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model: model name
	   :param model_type: model type
	   :return:
	*/
	// Get all provider configurations of the current workspace
	if s.provider_manager == nil {
		s.provider_manager = &providermanager.ProviderManager{}
	}
	provider_configurations := s.provider_manager.GetConfigurations(tenant_id)

	// Get provider configuration
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider)))
	}

	provider_configuration := provider_configurations.Configurations[provider]
	// Enable model
	provider_configuration.EnableModel(modelruntimeenumtypes.ModelType(model_type), model)
}
func (s *ModelProvideService) DisableModel(tenant_id string, provider string, model string, model_type string) {
	/*
	   disable model.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model: model name
	   :param model_type: model type
	   :return:
	*/
	// Get all provider configurations of the current workspace
	if s.provider_manager == nil {
		s.provider_manager = &providermanager.ProviderManager{}
	}
	provider_configurations := s.provider_manager.GetConfigurations(tenant_id)

	// Get provider configuration
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider)))
	}
	provider_configuration := provider_configurations.Configurations[provider]
	// Enable model
	provider_configuration.DisableModel(modelruntimeenumtypes.ModelType(model_type), model)
}
