package services

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	datamanager "mlib.com/gofy/server/core/manageres/data_manager"
	providermanager "mlib.com/gofy/server/core/manageres/provider_manager"
	modelproviders "mlib.com/gofy/server/core/model_runtime/model_provides"
	dbengine "mlib.com/gofy/server/db_engine"
	coreentities "mlib.com/gofy/server/entities/core"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	"mlib.com/gofy/server/models"
)

type ProviderConfigurationService struct {
}

func (service *ProviderConfigurationService) EnableModelLoadBalancing(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string) (*models.ProviderModelSetting, error) {
	/*
	   Enable model load balancing.
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return nil, exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	var load_balancing_config_count int64
	err := dbengine.Instance().DB.Model(&models.LoadBalancingModelConfig{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", provider_configuration.TenantID, provider_configuration.Provider.Provider, model_type, model).Count(&load_balancing_config_count).Error
	if err != nil {
		mlog.Error("count LoadBalancingModelConfig failed:", err)
		return nil, err
	}

	if load_balancing_config_count <= 1 {
		return nil, exceptions.NewValueError("Model load balancing configuration must be more than 1.")
	}
	model_setting := new(models.ProviderModelSetting)
	err = dbengine.Instance().DB.Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", provider_configuration.TenantID, provider_configuration.Provider.Provider, model_type, model).First(model_setting).Error
	if err != nil {
		mlog.Warningf("count ProviderModelSetting failed:%v", err)
		model_setting = nil
	}
	if model_setting != nil {
		model_setting.LoadBalancingEnabled = true
		now := time.Now()
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Save(model_setting)
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:                   uuid.NewV4().String(),
			TenantID:             provider_configuration.TenantID,
			ProviderName:         provider_configuration.Provider.Provider,
			ModelType:            string(model_type),
			ModelName:            model,
			LoadBalancingEnabled: true,
		}
		dbengine.Instance().DB.Create(model_setting)
	}
	return model_setting, nil
}
func (service *ProviderConfigurationService) DisableModelLoadBalancing(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string) (*models.ProviderModelSetting, error) {
	/*
	   Disable model load balancing.
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return nil, exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	model_setting := new(models.ProviderModelSetting)
	err := dbengine.Instance().DB.Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", provider_configuration.TenantID, provider_configuration.Provider.Provider, model_type, model).First(model_setting).Error
	if err != nil {
		mlog.Warningf("count ProviderModelSetting failed:%v", err)
		model_setting = nil
	}

	if model_setting != nil {
		model_setting.LoadBalancingEnabled = false
		now := time.Now()
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Save(model_setting)
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:                   uuid.NewV4().String(),
			TenantID:             provider_configuration.TenantID,
			ProviderName:         provider_configuration.Provider.Provider,
			ModelType:            string(model_type),
			ModelName:            model,
			LoadBalancingEnabled: false,
		}
		dbengine.Instance().DB.Create(model_setting)
	}
	return model_setting, nil
}

func (service *ProviderConfigurationService) CustomModelCredentialsValidate(
	provider_configuration *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string, credentials map[string]any,
) (*models.ProviderModel, map[string]any) {
	/*
		Validate custom model credentials.

		:param model_type: model type
		:param model: model name
		:param credentials: model credentials
		:return:
	*/
	// get provider model
	provider_model_record := new(models.ProviderModel)
	err := dbengine.Instance().DB.Model(&models.ProviderModel{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", provider_configuration.TenantID, provider_configuration.Provider.Provider, model_type, model).First(provider_model_record).Error
	if err != nil {
		mlog.Warningf("count ProviderModel failed:%v", err)
		provider_model_record = nil
	}

	// Get provider credential secret variables
	// var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	// if provider_configurations.Provider.ModelCredentialSchema != nil {
	// 	credential_form_schemas = provider_configurations.Provider.ModelCredentialSchema.CredentialFormSchemas
	// }
	// provider_credential_secret_variables := provider_configurations.ExtractSecretVariables(credential_form_schemas)

	// if provider_model_record != nil {
	// 	json.Unmarshal([]byte(provider_model_record.EncryptedConfig), )
	// 	// try:
	// 		original_credentials = (
	// 			json.loads(provider_model_record.encrypted_config) if provider_model_record.encrypted_config else {}
	// 		)
	// 	// except JSONDecodeError:
	// 	// 	original_credentials = {}

	// 	// decrypt credentials
	// 	for key, value in credentials.items():
	// 		if key in provider_credential_secret_variables:
	// 			// if send [__HIDDEN__] in secret input, it will be same as original value
	// 			if value == HIDDEN_VALUE and key in original_credentials:
	// 				credentials[key] = encrypter.decrypt_token(s.tenant_id, original_credentials[key])
	// 			}
	// 		}
	// 	}
	// }
	credentials = (&modelproviders.ModelProviderFactory{}).ModelCredentialsValidate(
		provider_configuration.Provider.Provider, model_type, model, credentials,
	)

	// for key, value := range credentials{
	// 	if slices.Contains(provider_credential_secret_variables,key){
	// 		credentials[key] = encrypter.encrypt_token(s.tenant_id, value)
	// 	}
	// }
	return provider_model_record, credentials
}

func (service *ProviderConfigurationService) AddOrUpdateCustomModelCredentials(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string, credentials map[string]any) error {
	/*
	   Add or update custom model credentials.

	   :param model_type: model type
	   :param model: model name
	   :param credentials: model credentials
	   :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	// validate custom model config
	provider_model_record, credentials := service.CustomModelCredentialsValidate(provider_configuration, model_type, model, credentials)

	bindata, _ := json.Marshal(credentials)
	// save provider model
	// Note: Do not switch the preferred provider, which allows users to use quotas first
	if provider_model_record != nil {

		provider_model_record.EncryptedConfig = string(bindata)
		provider_model_record.IsValid = true
		now := time.Now()
		provider_model_record.UpdatedAt = &now
		dbengine.Instance().DB.Save(provider_model_record)
	} else {
		provider_model_record = &models.ProviderModel{
			ID:              uuid.NewV4().String(),
			TenantID:        provider_configuration.TenantID,
			ProviderName:    provider_configuration.Provider.Provider,
			ModelName:       model,
			ModelType:       string(model_type),
			EncryptedConfig: string(bindata),
			IsValid:         true,
		}
		dbengine.Instance().DB.Create(provider_model_record)
	}
	provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
		provider_configuration.TenantID,
		provider_model_record.ID,
		datamanager.ProviderCredentialsCache_MODEL,
	)

	provider_model_credentials_cache.Delete()
	return nil
}
func (service *ProviderConfigurationService) DeleteCustomModelCredentials(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string) error {
	/*
	   Delete custom model credentials.
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	// get provider model
	provider_model_record := new(models.ProviderModel)
	err := dbengine.Instance().DB.Model(&models.ProviderModel{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", provider_configuration.TenantID, provider_configuration.Provider.Provider, model_type, model).First(provider_model_record).Error
	if err != nil {
		mlog.Warningf("count ProviderModel failed:%v", err)
		provider_model_record = nil
	}

	// delete provider model
	if provider_model_record != nil {
		dbengine.Instance().DB.Delete(provider_model_record)

		provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
			provider_configuration.TenantID,
			provider_model_record.ID,
			datamanager.ProviderCredentialsCache_MODEL,
		)

		provider_model_credentials_cache.Delete()
	}
	return nil
}

func (s *ProviderConfigurationService) DeleteCustomCredentials(pc *coreentities.ProviderConfiguration) {
	/*
	   Delete custom Provider credentials.
	   :return:
	*/
	// get Provider
	provider_record := new(models.Provider)
	err := dbengine.Instance().DB.Model(&models.Provider{}).Where("tenant_id = ? and provider_name = ? and provider_type = ?", pc.TenantID, pc.Provider.Provider, providerenumtypes.Provider_CUSTOM).First(provider_record).Error
	if err != nil {
		mlog.Errorf("get TenantPreferredModelProvider failed:%v", err)
		provider_record = nil
		// return
	}

	// delete Provider
	if provider_record != nil {
		s.SwitchPreferredProviderType(pc, providerenumtypes.Provider_SYSTEM)
		dbengine.Instance().DB.Delete(provider_record)

		provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(pc.TenantID, provider_record.ID, datamanager.ProviderCredentialsCache_PROVIDER)

		provider_model_credentials_cache.Delete()
	}
}
func (s *ProviderConfigurationService) ObfuscatedCredentials(pc *coreentities.ProviderConfiguration, credentials map[string]any, credential_form_schemas []*modelruntimeentities.CredentialFormSchema) map[string]any {
	/*
	   Obfuscated credentials.

	   :param credentials: credentials
	   :param credential_form_schemas: credential form schemas
	   :return:
	*/
	// Get provider credential secret variables
	credential_secret_variables := pc.ExtractSecretVariables(credential_form_schemas)

	// Obfuscate provider credentials
	copy_credentials := maps.Clone(credentials)
	for key, value := range copy_credentials {
		if slices.Contains(credential_secret_variables, key) {
			// copy_credentials[key] = encrypter.obfuscated_token(value)
			copy_credentials[key] = value
		}
	}
	return copy_credentials
}

func (s *ProviderConfigurationService) GetCustomCredentials(tenant_id string, provider string, obfuscated bool) (map[string]any, error) {
	/*
	   Get custom credentials.

	   :param obfuscated: obfuscated secret data in credentials
	   :return:
	*/
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return nil, exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]
	mlog.Debugf("------CustomConfiguration=%#v", provider_configuration.CustomConfiguration)
	if provider_configuration.CustomConfiguration.Provider == nil {
		mlog.Errorf("Provider %s does not exist CustomConfiguration", provider)
		return nil, exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist CustomConfiguration", provider))
	}
	credentials := provider_configuration.CustomConfiguration.Provider.Credentials
	if !obfuscated {
		return credentials, nil
	}
	credential_form_schemas := []*modelruntimeentities.CredentialFormSchema{}
	if provider_configuration.Provider.ProviderCredentialSchema != nil {
		credential_form_schemas = provider_configuration.Provider.ProviderCredentialSchema.CredentialFormSchemas
	}
	// Obfuscate credentials
	return s.ObfuscatedCredentials(
		provider_configuration,
		credentials,
		credential_form_schemas,
	), nil
}

func (s *ProviderConfigurationService) CustomCredentialsValidate(pc *coreentities.ProviderConfiguration, credentials map[string]any) (*models.Provider, map[string]any) {
	/*
	   Validate custom credentials.
	   :param credentials: provider credentials
	   :return:
	*/
	// get provider
	provider_record := new(models.Provider)
	err := dbengine.Instance().DB.Model(&models.Provider{}).Where("tenant_id = ? and provider_name = ? and provider_type = ?", pc.TenantID, pc.Provider.Provider, providerenumtypes.Provider_CUSTOM).First(provider_record).Error
	if err != nil {
		mlog.Warningf("get TenantPreferredModelProvider failed:%v", err)
		provider_record = nil
		// return
	}

	// Get provider credential secret variables
	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	if pc.Provider.ProviderCredentialSchema != nil {
		credential_form_schemas = pc.Provider.ProviderCredentialSchema.CredentialFormSchemas
	}
	provider_credential_secret_variables := pc.ExtractSecretVariables(credential_form_schemas)

	if provider_record != nil {
		var original_credentials map[string]any

		if provider_record.EncryptedConfig != "" {
			if !strings.HasPrefix(provider_record.EncryptedConfig, "{") {
				original_credentials = map[string]any{"openai_api_key": provider_record.EncryptedConfig}
			} else {
				err := json.Unmarshal([]byte(provider_record.EncryptedConfig), &original_credentials)
				if err != nil {
					mlog.Warningf("json unmarshal(%s) failed:%v", provider_record.EncryptedConfig, err)
				}
			}
		} else {
			original_credentials = map[string]any{}
		}

		// encrypt credentials
		for key, value := range credentials {
			if slices.Contains(provider_credential_secret_variables, key) {
				// if send [__HIDDEN__] in secret input, it will be same as original value
				if _, ok := original_credentials[key]; ok && value == constants.HIDDEN_VALUE {
					credentials[key] = original_credentials[key]
				}
			}
		}
	}
	credentials = (&modelproviders.ModelProviderFactory{}).ProviderCredentialsValidate(pc.Provider.Provider, credentials)

	return provider_record, credentials
}

func (s *ProviderConfigurationService) AddOrUpdateCustomCredentials(tenant_id string, provider string, credentials map[string]any) error {
	/*
	   Add or update custom provider credentials.
	   :param credentials:
	   :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	// validate custom provider config
	provider_record, credentials := s.CustomCredentialsValidate(provider_configuration, credentials)

	// save provider
	// Note: Do not switch the preferred provider, which allows users to use quotas first
	bindata, _ := json.Marshal(credentials)
	if provider_record != nil {

		provider_record.EncryptedConfig = string(bindata)
		provider_record.IsValid = true
		now := time.Now()
		provider_record.UpdatedAt = &now
		dbengine.Instance().DB.Save(provider_record)
	} else {
		provider_record = &models.Provider{
			ID:              uuid.NewV4().String(),
			TenantID:        provider_configuration.TenantID,
			ProviderName:    provider_configuration.Provider.Provider,
			ProviderType:    providerenumtypes.Provider_CUSTOM,
			EncryptedConfig: string(bindata),
			IsValid:         true,
		}
		dbengine.Instance().DB.Create(provider_record)
	}
	provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
		provider_configuration.TenantID, provider_record.ID, datamanager.ProviderCredentialsCache_PROVIDER,
	)

	provider_model_credentials_cache.Delete()

	s.SwitchPreferredProviderType(provider_configuration, providerenumtypes.Provider_CUSTOM)
	return nil
}

func (s *ProviderConfigurationService) SwitchPreferredProviderType(pc *coreentities.ProviderConfiguration, provider_type providerenumtypes.ProviderType) {
	/*
	   Switch preferred provider type.
	   :param provider_type:
	   :return:
	*/
	if provider_type == pc.PreferredProviderType {
		return
	}
	if provider_type == providerenumtypes.Provider_SYSTEM && !pc.SystemConfiguration.Enabled {
		return
	}
	// get preferred provider
	preferred_model_provider := new(models.TenantPreferredModelProvider)
	err := dbengine.Instance().DB.Model(&models.TenantPreferredModelProvider{}).Where("tenant_id = ? and provider_name = ?", pc.TenantID, pc.Provider.Provider).First(preferred_model_provider).Error
	if err != nil {
		mlog.Warningf("get TenantPreferredModelProvider failed:%v", err)
		preferred_model_provider = nil
		// return
	}

	if preferred_model_provider != nil {
		preferred_model_provider.PreferredProviderType = provider_type
		dbengine.Instance().DB.Save(preferred_model_provider)
	} else {
		preferred_model_provider = &models.TenantPreferredModelProvider{
			ID:                    uuid.NewV4().String(),
			TenantID:              pc.TenantID,
			ProviderName:          pc.Provider.Provider,
			PreferredProviderType: provider_type,
		}
		dbengine.Instance().DB.Create(preferred_model_provider)
	}
}

func (s *ProviderConfigurationService) DeleteCustomCredentialsByTenantAndProvider(tenant_id string, provider string) error {
	/*
	   Delete custom provider credentials.
	   :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]
	s.DeleteCustomCredentials(provider_configuration)
	return nil
}

func (s *ProviderConfigurationService) DisableModel(tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string) (*models.ProviderModelSetting, error) {
	/*
	   Disable model.
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	// Get all provider configurations of the current workspace
	provider_configurations := (&providermanager.ProviderManager{}).GetConfigurations(tenant_id)
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		return nil, exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	model_setting := new(models.ProviderModelSetting)
	err := dbengine.Instance().DB.Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", provider_configuration.TenantID, provider_configuration.Provider.Provider, model_type, model).First(model_setting).Error
	if err != nil {
		mlog.Warningf("get ProviderModelSetting failed:%v", err)
		model_setting = nil
		// return
	}
	now := time.Now()
	if model_setting != nil {
		model_setting.Enabled = false
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Debug().Save(model_setting)
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:           uuid.NewV4().String(),
			TenantID:     provider_configuration.TenantID,
			ProviderName: provider_configuration.Provider.Provider,
			ModelType:    string(model_type),
			ModelName:    model,
			Enabled:      false,
		}
		mlog.Debugf("------model_setting=%#v", model_setting)
		dbengine.Instance().DB.Debug().Create(model_setting)
	}
	mlog.Debugf("------model_setting=%#v", model_setting)
	return model_setting, nil
}
