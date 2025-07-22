package providermanager

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/exceptions"
	hostingconfiguration "mlib.com/gofy/server/core/hosting_configuration"
	datamanager "mlib.com/gofy/server/core/manageres/data_manager"
	modelproviders "mlib.com/gofy/server/core/model_runtime/model_provides"
	dbengine "mlib.com/gofy/server/db_engine"
	coreentities "mlib.com/gofy/server/entities/core"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	coreenumtypes "mlib.com/gofy/server/enum_types/core"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

var (
	original_provider_configurate_methods = map[string][]modelruntimeentities.ConfigurateMethod{}
)

type ProviderManager struct {
	decoding_rsa_key    string
	decoding_cipher_rsa string
}

func (pm *ProviderManager) to_model_settings(
	provider_entity *modelruntimeentities.ProviderEntity,
	provider_model_settings []*models.ProviderModelSetting,
	load_balancing_model_configs []*models.LoadBalancingModelConfig,
) []*coreentities.ModelSetting {
	/*
		Convert to model settings.
		:param provider_entity: provider entity
		:param provider_model_settings: provider model settings include enabled, load balancing enabled
		:param load_balancing_model_configs: load balancing model configs
		:return:
	*/
	// Get provider model credential secret variables
	// var model_credential_secret_variables []string
	// if slices.Contains(provider_entity.ConfigurateMethods, modelruntimeentities.ConfigurateMethod_PREDEFINED_MODEL) {
	// 	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	// 	if provider_entity.ProviderCredentialSchema != nil {
	// 		credential_form_schemas = provider_entity.ProviderCredentialSchema.CredentialFormSchemas
	// 	}
	// 	model_credential_secret_variables = pm.extract_secret_variables(credential_form_schemas)
	// } else {
	// 	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	// 	if provider_entity.ModelCredentialSchema != nil {
	// 		credential_form_schemas = provider_entity.ModelCredentialSchema.CredentialFormSchemas
	// 	}
	// 	model_credential_secret_variables = pm.extract_secret_variables(credential_form_schemas)
	// }
	model_settings := []*coreentities.ModelSetting{}
	if len(provider_model_settings) == 0 {
		return model_settings
	}
	for _, provider_model_setting := range provider_model_settings {
		load_balancing_configs := []*coreentities.ModelLoadBalancingConfiguration{}
		if provider_model_setting.LoadBalancingEnabled && len(load_balancing_model_configs) > 0 {
			for _, load_balancing_model_config := range load_balancing_model_configs {
				if load_balancing_model_config.ModelName == provider_model_setting.ModelName && load_balancing_model_config.ModelType == provider_model_setting.ModelType {
					if !load_balancing_model_config.Enabled {
						continue
					}
					if load_balancing_model_config.EncryptedConfig == "" {
						if load_balancing_model_config.Name == "__inherit__" {
							load_balancing_configs = append(load_balancing_configs, &coreentities.ModelLoadBalancingConfiguration{
								ID:          load_balancing_model_config.ID,
								Name:        load_balancing_model_config.Name,
								Credentials: map[string]any{},
							})
						}
						continue
					}
					provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
						load_balancing_model_config.TenantID,
						load_balancing_model_config.ID,
						datamanager.ProviderCredentialsCache_LOAD_BALANCING_MODEL,
					)

					// Get cached provider model credentials
					cached_provider_model_credentials := provider_model_credentials_cache.Get()
					var provider_model_credentials map[string]any
					if cached_provider_model_credentials == nil {
						// try{

						err := json.Unmarshal([]byte(load_balancing_model_config.EncryptedConfig), &provider_model_credentials)
						if err != nil {
							mlog.Errorf("json unmarshal(%s) failed:%v", load_balancing_model_config.EncryptedConfig, err)
							continue
						}
						// except JSONDecodeError{
						// 	continue

						// // Get decoding rsa key and cipher for decrypting credentials
						// if pm.decoding_rsa_key is None or pm.decoding_cipher_rsa is None{
						// 	pm.decoding_rsa_key, pm.decoding_cipher_rsa = encrypter.get_decrypt_decoding(
						// 		load_balancing_model_config.tenant_id
						// 	)
						// }
						// for variable in model_credential_secret_variables{
						// 	if variable in provider_model_credentials{
						// 		// try{
						// 			provider_model_credentials[variable] = encrypter.decrypt_token_with_decoding(
						// 				provider_model_credentials.get(variable),
						// 				pm.decoding_rsa_key,
						// 				pm.decoding_cipher_rsa,
						// 			)
						// 		// except ValueError{
						// 		// 	pass
						// 	}
						// }
						// cache provider model credentials
						provider_model_credentials_cache.Set(provider_model_credentials)
					} else {
						provider_model_credentials = cached_provider_model_credentials
					}
					load_balancing_configs = append(load_balancing_configs, &coreentities.ModelLoadBalancingConfiguration{
						ID:          load_balancing_model_config.ID,
						Name:        load_balancing_model_config.Name,
						Credentials: provider_model_credentials,
					})
				}
			}
		}
		model_settings = append(model_settings, &coreentities.ModelSetting{
			Model:                provider_model_setting.ModelName,
			ModelType:            modelruntimeentities.ModelType(provider_model_setting.ModelType),
			Enabled:              provider_model_setting.Enabled,
			LoadBalancingConfigs: load_balancing_configs,
		})
	}
	return model_settings
}

func (pm *ProviderManager) get_all_providers(tenant_id string) map[string][]*models.Provider {
	/*
	   Get all provider records of the workspace.

	   :param tenant_id: workspace id
	   :return:
	*/
	var providers []*models.Provider
	err := dbengine.Instance().DB.Model(&models.Provider{}).Where("tenant_id = ? and is_valid = ?", tenant_id, 1).Find(&providers).Error
	if err != nil {
		mlog.Error("get provider failed:", err)
		return nil
	}

	provider_name_to_provider_records_dict := map[string][]*models.Provider{}
	for _, provider := range providers {
		if _, ok := provider_name_to_provider_records_dict[provider.ProviderName]; !ok {
			provider_name_to_provider_records_dict[provider.ProviderName] = make([]*models.Provider, 0)
		}
		provider_name_to_provider_records_dict[provider.ProviderName] = append(provider_name_to_provider_records_dict[provider.ProviderName], provider)
	}
	return provider_name_to_provider_records_dict
}

func (pm *ProviderManager) init_trial_provider_records(
	tenant_id string, provider_name_to_provider_records_dict map[string][]*models.Provider,
) map[string][]*models.Provider {
	/*
		Initialize trial provider records if not exists.

		:param tenant_id: workspace id
		:param provider_name_to_provider_records_dict: provider name to provider records dict
		:return:
	*/
	// Get hosting configuration
	for provider_name, configuration := range hostingconfiguration.Instance().ProviderMap {
		if !configuration.Enabled {
			continue
		}
		provider_records, ok := provider_name_to_provider_records_dict[provider_name]
		if !ok || provider_records == nil {
			provider_records = []*models.Provider{}
		}
		provider_quota_to_provider_record_dict := map[models.ProviderQuotaType]*models.Provider{}
		for _, provider_record := range provider_records {
			if provider_record.ProviderType != models.Provider_SYSTEM {
				continue
			}
			provider_quota_to_provider_record_dict[models.ProviderQuotaType(provider_record.QuotaType)] = provider_record
		}
		for _, quota := range configuration.Quotas {
			if quota.Type() == models.ProviderQuota_TRIAL {
				// Init trial provider records if not exists
				if _, ok := provider_quota_to_provider_record_dict[models.ProviderQuota_TRIAL]; !ok {
					realquota := any(quota).(*hostingconfiguration.TrialHostingQuota)
					// try{
					// FIXME ignore the type errork, onyl TrialHostingQuota has limit need to change the logic
					provider_record := &models.Provider{
						ID:           uuid.NewV4().String(),
						TenantID:     tenant_id,
						ProviderName: provider_name,
						ProviderType: models.Provider_SYSTEM,
						QuotaType:    string(models.ProviderQuota_TRIAL),
						QuotaLimit:   int64(realquota.QuotaLimit), // type: ignore
						QuotaUsed:    0,
						IsValid:      true,
					}
					dbengine.Instance().DB.Create(provider_record)
					// except IntegrityError{
					// 	db.session.rollback()
					// 	provider_record = (
					// 		db.session.query(Provider)
					// 		.filter(
					// 			Provider.tenant_id == tenant_id,
					// 			Provider.Provider_name == provider_name,
					// 			Provider.ProviderType == ProviderType.SYSTEM.value,
					// 			Provider.quota_type == ProviderQuotaType.TRIAL.value,
					// 		)
					// 		.first()
					// 	)

					if !provider_record.IsValid {
						provider_record.IsValid = true
						dbengine.Instance().DB.Save(provider_record)
					}

					if provider_name_to_provider_records_dict == nil {
						provider_name_to_provider_records_dict = make(map[string][]*models.Provider)
					}
					if provider_name_to_provider_records_dict[provider_name] == nil {
						provider_name_to_provider_records_dict[provider_name] = make([]*models.Provider, 0)
					}
					provider_name_to_provider_records_dict[provider_name] = append(provider_name_to_provider_records_dict[provider_name], provider_record)
				}
			}
		}
	}

	return provider_name_to_provider_records_dict
}

func (pm *ProviderManager) get_all_provider_models(tenant_id string) map[string][]*models.ProviderModel {
	/*
	   Get all provider model records of the workspace.

	   :param tenant_id: workspace id
	   :return:
	*/
	// Get all provider model records of the workspace
	var provider_models []*models.ProviderModel
	err := dbengine.Instance().DB.Model(&models.ProviderModel{}).Where("tenant_id = ? and is_valid = ?", tenant_id, 1).Find(&provider_models).Error
	if err != nil {
		mlog.Error("get provider model failed:", err)
		return nil
	}

	provider_name_to_provider_model_records_dict := map[string][]*models.ProviderModel{}
	for _, provider_model := range provider_models {
		if provider_name_to_provider_model_records_dict[provider_model.ProviderName] == nil {
			provider_name_to_provider_model_records_dict[provider_model.ProviderName] = make([]*models.ProviderModel, 0)
		}
		provider_name_to_provider_model_records_dict[provider_model.ProviderName] = append(provider_name_to_provider_model_records_dict[provider_model.ProviderName], provider_model)
	}
	return provider_name_to_provider_model_records_dict
}

func (pm *ProviderManager) get_all_preferred_model_providers(tenant_id string) map[string]*models.TenantPreferredModelProvider {
	/*
	   Get All preferred provider types of the workspace.

	   :param tenant_id: workspace id
	   :return:
	*/
	var preferred_provider_types []*models.TenantPreferredModelProvider
	err := dbengine.Instance().DB.Model(&models.TenantPreferredModelProvider{}).Where("tenant_id = ?", tenant_id).Find(&preferred_provider_types).Error
	if err != nil {
		mlog.Error("get Tenant Preferred Model Provider failed:", err)
		return nil
	}

	provider_name_to_preferred_provider_type_records_dict := map[string]*models.TenantPreferredModelProvider{}
	for _, preferred_provider_type := range preferred_provider_types {
		provider_name_to_preferred_provider_type_records_dict[preferred_provider_type.ProviderName] = preferred_provider_type
	}

	return provider_name_to_preferred_provider_type_records_dict
}

func (pm *ProviderManager) get_all_provider_model_settings(tenant_id string) map[string][]*models.ProviderModelSetting {
	/*
	   Get All provider model settings of the workspace.

	   :param tenant_id: workspace id
	   :return:
	*/
	var provider_model_settings []*models.ProviderModelSetting
	err := dbengine.Instance().DB.Model(&models.ProviderModelSetting{}).Where("tenant_id = ?", tenant_id).Find(&provider_model_settings).Error
	if err != nil {
		mlog.Error("get provider model settings failed:", err)
		return nil
	}

	provider_name_to_provider_model_settings_dict := map[string][]*models.ProviderModelSetting{}
	for _, provider_model_setting := range provider_model_settings {
		if _, ok := provider_name_to_provider_model_settings_dict[provider_model_setting.ProviderName]; !ok {
			provider_name_to_provider_model_settings_dict[provider_model_setting.ProviderName] = make([]*models.ProviderModelSetting, 0)
		}
		provider_name_to_provider_model_settings_dict[provider_model_setting.ProviderName] = append(provider_name_to_provider_model_settings_dict[provider_model_setting.ProviderName], provider_model_setting)
	}
	return provider_name_to_provider_model_settings_dict
}

func (pm *ProviderManager) _extract_secret_variables(credential_form_schemas []*modelruntimeentities.CredentialFormSchema) []string {
	/*
	   Extract secret input form variables.

	   :param credential_form_schemas:
	   :return:
	*/
	secret_input_form_variables := []string{}
	for _, credential_form_schema := range credential_form_schemas {
		if credential_form_schema.Type == modelruntimeentities.Form_SECRET_INPUT {
			secret_input_form_variables = append(secret_input_form_variables, credential_form_schema.Variable)
		}
	}
	return secret_input_form_variables

}
func (pm *ProviderManager) to_custom_configuration(
	tenant_id string,
	provider_entity *modelruntimeentities.ProviderEntity,
	provider_records []*models.Provider,
	provider_model_records []*models.ProviderModel,
) *coreentities.CustomConfiguration {
	/*
		Convert to custom configuration.

		:param tenant_id: workspace id
		:param provider_entity: provider entity
		:param provider_records: provider records
		:param provider_model_records: provider model records
		:return
	*/
	// Get provider credential secret variables
	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	if provider_entity.ProviderCredentialSchema != nil {
		credential_form_schemas = provider_entity.ProviderCredentialSchema.CredentialFormSchemas
	}
	provider_credential_secret_variables := pm._extract_secret_variables(credential_form_schemas)
	mlog.Debugf("------provider_credential_secret_variables=%#v", provider_credential_secret_variables)

	// Get custom provider record
	mlog.Debugf("------provider_records=%#v", provider_records)
	var custom_provider_record *models.Provider
	for _, provider_record := range provider_records {
		mlog.Debugf("------provider_record=%#v", provider_record)
		if provider_record.ProviderType == models.Provider_SYSTEM {
			continue
		}
		if provider_record.EncryptedConfig == "" {
			continue
		}
		custom_provider_record = provider_record
	}
	// Get custom provider credentials
	var custom_provider_configuration *coreentities.CustomProviderConfiguration
	mlog.Debugf("------custom_provider_record=%#v", custom_provider_record)
	if custom_provider_record != nil {
		provider_credentials_cache := datamanager.NewProviderCredentialsCache(
			tenant_id,
			custom_provider_record.ID,
			datamanager.ProviderCredentialsCache_PROVIDER,
		)

		// Get cached provider credentials
		cached_provider_credentials := provider_credentials_cache.Get()
		var provider_credentials map[string]any
		if len(cached_provider_credentials) == 0 {
			// fix origin data
			if custom_provider_record.EncryptedConfig != "" && !strings.HasPrefix(custom_provider_record.EncryptedConfig, "{") {
				provider_credentials = map[string]any{"openai_api_key": custom_provider_record.EncryptedConfig}
			} else {
				err := json.Unmarshal([]byte(custom_provider_record.EncryptedConfig), &provider_credentials)
				if err != nil {
					mlog.Errorf("json unmarshal(%s) failed:%v", custom_provider_record.EncryptedConfig, err)
					provider_credentials = map[string]any{}
				}
			}

			// Get decoding rsa key and cipher for decrypting credentials
			// if pm.decoding_rsa_key == "" || pm.decoding_cipher_rsa == ""{
			// 	pm.decoding_rsa_key, pm.decoding_cipher_rsa = encrypter.get_decrypt_decoding(tenant_id)
			// }
			// for _,  variable := range provider_credential_secret_variables{
			// 	if _, ok := provider_credentials[variable]; ok {
			// 		// try{
			// 			provider_credentials[variable] =
			// 			provider_credentials[variable] = encrypter.decrypt_token_with_decoding(
			// 				provider_credentials.get(variable) or "",  // type: ignore
			// 				pm.decoding_rsa_key,
			// 				pm.decoding_cipher_rsa,
			// 			)
			// 		// except ValueError{
			// 		// 	pass
			// 	}
			// }
			// cache provider credentials
			provider_credentials_cache.Set(provider_credentials)
		} else {
			provider_credentials = cached_provider_credentials
		}
		custom_provider_configuration = &coreentities.CustomProviderConfiguration{Credentials: provider_credentials}
	}
	// if provider_entity.ModelCredentialSchema != nil {
	// 	credential_form_schemas = provider_entity.ModelCredentialSchema.CredentialFormSchemas
	// }
	// // Get provider model credential secret variables
	// model_credential_secret_variables := pm.extract_secret_variables(credential_form_schemas)

	// Get custom provider model credentials
	custom_model_configurations := []*coreentities.CustomModelConfiguration{}
	for _, provider_model_record := range provider_model_records {
		if provider_model_record.EncryptedConfig == "" {
			continue
		}
		provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
			tenant_id, provider_model_record.ID, datamanager.ProviderCredentialsCache_MODEL,
		)

		// Get cached provider model credentials
		cached_provider_model_credentials := provider_model_credentials_cache.Get()
		var provider_model_credentials map[string]any
		if len(cached_provider_model_credentials) == 0 {
			// try{

			err := json.Unmarshal([]byte(provider_model_record.EncryptedConfig), &provider_model_credentials)
			if err != nil {
				continue
			}
			// except JSONDecodeError{
			// 	continue

			// Get decoding rsa key and cipher for decrypting credentials
			// if pm.decoding_rsa_key is None or pm.decoding_cipher_rsa is None{
			// 	pm.decoding_rsa_key, pm.decoding_cipher_rsa = encrypter.get_decrypt_decoding(tenant_id)
			// }
			// for variable := range model_credential_secret_variables{
			// 	if variable := range provider_model_credentials{
			// 		// try{
			// 			provider_model_credentials[variable] = encrypter.decrypt_token_with_decoding(
			// 				provider_model_credentials.get(variable),
			// 				pm.decoding_rsa_key,
			// 				pm.decoding_cipher_rsa,
			// 			)
			// 		// except ValueError{
			// 		// 	pass
			// 	}
			// }
			// cache provider model credentials
			provider_model_credentials_cache.Set(provider_model_credentials)
		} else {
			provider_model_credentials = cached_provider_model_credentials
		}
		custom_model_configurations = append(custom_model_configurations, &coreentities.CustomModelConfiguration{
			Model:       provider_model_record.ModelName,
			ModelType:   modelruntimeentities.ModelType(provider_model_record.ModelType),
			Credentials: provider_model_credentials,
		})
	}
	return &coreentities.CustomConfiguration{Provider: custom_provider_configuration, Models: custom_model_configurations}
}
func (pm *ProviderManager) choice_current_using_quota_type(quota_configurations []*coreentities.QuotaConfiguration) (models.ProviderQuotaType, error) {
	/*
	   Choice current using quota type.
	   paid quotas > provider free quotas > hosting trial quotas
	   If there is still quota for the corresponding quota type according to the sorting,

	   :param quota_configurations:
	   :return:
	*/
	// convert to dict
	quota_type_to_quota_configuration_dict := map[models.ProviderQuotaType]*coreentities.QuotaConfiguration{}
	for _, quota_configuration := range quota_configurations {
		quota_type_to_quota_configuration_dict[quota_configuration.QuotaType] = quota_configuration
	}

	var last_quota_configuration *coreentities.QuotaConfiguration
	for _, quota_type := range []models.ProviderQuotaType{models.ProviderQuota_PAID, models.ProviderQuota_FREE, models.ProviderQuota_TRIAL} {
		if _, ok := quota_type_to_quota_configuration_dict[quota_type]; ok {
			last_quota_configuration = quota_type_to_quota_configuration_dict[quota_type]
			if last_quota_configuration.IsValid {
				return quota_type, nil
			}
		}
	}

	if last_quota_configuration != nil {
		return last_quota_configuration.QuotaType, nil
	}
	return models.ProviderQuotaType(""), exceptions.NewValueError("No quota type available")
}

func (pm *ProviderManager) to_system_configuration(
	tenant_id string, provider_entity *modelruntimeentities.ProviderEntity, provider_records []*models.Provider,
) *coreentities.SystemConfiguration {
	/*
		Convert to system configuration.

		:param tenant_id: workspace id
		:param provider_entity: provider entity
		:param provider_records: provider records
		:return
	*/
	// Get hosting configuration
	provider_hosting_configuration := hostingconfiguration.Instance().ProviderMap[provider_entity.Provider]
	if provider_hosting_configuration == nil || !provider_hosting_configuration.Enabled {
		return &coreentities.SystemConfiguration{Enabled: false}
	}
	// Convert provider_records to dict
	quota_type_to_provider_records_dict := map[models.ProviderQuotaType]*models.Provider{}
	for _, provider_record := range provider_records {
		if provider_record.ProviderType != models.Provider_SYSTEM {
			continue
		}
		quota_type_to_provider_records_dict[models.ProviderQuotaType(provider_record.QuotaType)] = provider_record
	}
	quota_configurations := []*coreentities.QuotaConfiguration{}
	for _, provider_quota := range provider_hosting_configuration.Quotas {
		var quota_configuration *coreentities.QuotaConfiguration
		if _, ok := quota_type_to_provider_records_dict[provider_quota.Type()]; !ok {
			if provider_quota.Type() == models.ProviderQuota_FREE {
				quota_unit := coreenumtypes.QuotaUnit_TOKENS
				if string(provider_hosting_configuration.QuotaUnit) != "" {
					quota_unit = provider_hosting_configuration.QuotaUnit
				}
				quota_configuration = &coreentities.QuotaConfiguration{
					QuotaType:      provider_quota.Type(),
					QuotaUnit:      quota_unit,
					QuotaUsed:      0,
					QuotaLimit:     0,
					IsValid:        false,
					RestrictModels: provider_quota.GetRestrictModels(),
				}
			} else {
				continue
			}
		} else {
			provider_record := quota_type_to_provider_records_dict[provider_quota.Type()]
			quota_unit := coreenumtypes.QuotaUnit_TOKENS
			if string(provider_hosting_configuration.QuotaUnit) != "" {
				quota_unit = provider_hosting_configuration.QuotaUnit
			}
			quota_configuration = &coreentities.QuotaConfiguration{
				QuotaType:  provider_quota.Type(),
				QuotaUnit:  quota_unit,
				QuotaUsed:  int(provider_record.QuotaUsed),
				QuotaLimit: int(provider_record.QuotaLimit),
				IsValid:    provider_record.QuotaLimit > provider_record.QuotaUsed || provider_record.QuotaLimit == -1,
			}
		}

		quota_configurations = append(quota_configurations, quota_configuration)
	}
	if len(quota_configurations) == 0 {
		return &coreentities.SystemConfiguration{Enabled: false}
	}
	current_quota_type, _ := pm.choice_current_using_quota_type(quota_configurations)

	current_using_credentials := provider_hosting_configuration.Credentials
	if current_quota_type == models.ProviderQuota_FREE {
		provider_record_quota_free := quota_type_to_provider_records_dict[current_quota_type]

		if provider_record_quota_free != nil {
			provider_credentials_cache := datamanager.NewProviderCredentialsCache(
				tenant_id,
				provider_record_quota_free.ID,
				datamanager.ProviderCredentialsCache_PROVIDER,
			)

			// Get cached provider credentials
			cached_provider_credentials := provider_credentials_cache.Get()

			if len(cached_provider_credentials) == 0 {
				var provider_credentials map[string]any
				// try{
				err := json.Unmarshal([]byte(provider_record_quota_free.EncryptedConfig), &provider_credentials)
				if err != nil {
					mlog.Errorf("json unmarshal(%s) failed:%v", provider_record_quota_free.EncryptedConfig, err)
				}
				// except JSONDecodeError{
				// 	provider_credentials = {}

				// Get provider credential secret variables
				// var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
				// if provider_entity.ProviderCredentialSchema != nil {
				// 	credential_form_schemas = provider_entity.ProviderCredentialSchema.CredentialFormSchemas
				// }
				// provider_credential_secret_variables := pm.extract_secret_variables(credential_form_schemas)

				// Get decoding rsa key and cipher for decrypting credentials
				// if pm.decoding_rsa_key is None or pm.decoding_cipher_rsa is None{
				// 	pm.decoding_rsa_key, pm.decoding_cipher_rsa = encrypter.get_decrypt_decoding(tenant_id)
				// }
				// for variable := range provider_credential_secret_variables{
				// 	if variable in provider_credentials{
				// 		// try{
				// 			provider_credentials[variable] = encrypter.decrypt_token_with_decoding(
				// 				provider_credentials.get(variable), pm.decoding_rsa_key, pm.decoding_cipher_rsa
				// 			)
				// 		// except ValueError{
				// 		// 	pass
				// 	}
				// }
				current_using_credentials := provider_credentials

				// cache provider credentials
				provider_credentials_cache.Set(current_using_credentials)
			} else {
				current_using_credentials = cached_provider_credentials
			}
		} else {
			current_using_credentials = map[string]any{}
			quota_configurations = []*coreentities.QuotaConfiguration{}
		}
	}
	return &coreentities.SystemConfiguration{
		Enabled:             true,
		CurrentQuotaType:    current_quota_type,
		QuotaConfigurations: quota_configurations,
		Credentials:         current_using_credentials,
	}
}

func (pm *ProviderManager) get_all_provider_load_balancing_configs(tenant_id string) map[string][]*models.LoadBalancingModelConfig {
	/*
	   Get All provider load balancing configs of the workspace.

	   :param tenant_id: workspace id
	   :return:
	*/
	model_load_balancing_enabled := false
	// cache_key := fmt.Sprintf("tenant:%s:model_load_balancing_enabled", tenant_id)
	// if !cache.Instance().ExistsKey(cache_key) {
	// 	model_load_balancing_enabled = services.ServiceGroupApp.Feature.GetFeatures(tenant_id).ModelLoadBalancingEnabled
	// 	cache.Instance().SetEx(cache_key, model_load_balancing_enabled, 120*time.Second)
	// } else {
	// 	model_load_balancing_enabled = cache.Instance().GetBool(cache_key)
	// }
	if !model_load_balancing_enabled {
		return nil
	}

	var provider_load_balancing_configs []*models.LoadBalancingModelConfig
	err := dbengine.Instance().DB.Model(&models.LoadBalancingModelConfig{}).Where("tenant_id = ?", tenant_id).Find(&provider_load_balancing_configs).Error
	if err != nil {
		mlog.Error("get provider load balancing configs failed:", err)
		return nil
	}

	provider_name_to_provider_load_balancing_model_configs_dict := map[string][]*models.LoadBalancingModelConfig{}
	for _, provider_load_balancing_config := range provider_load_balancing_configs {
		if _, ok := provider_name_to_provider_load_balancing_model_configs_dict[provider_load_balancing_config.ProviderName]; !ok {
			provider_name_to_provider_load_balancing_model_configs_dict[provider_load_balancing_config.ProviderName] = make([]*models.LoadBalancingModelConfig, 0)
		}
		provider_name_to_provider_load_balancing_model_configs_dict[provider_load_balancing_config.ProviderName] = append(provider_name_to_provider_load_balancing_model_configs_dict[provider_load_balancing_config.ProviderName], provider_load_balancing_config)
	}
	return provider_name_to_provider_load_balancing_model_configs_dict

}
func (pm *ProviderManager) GetConfigurations(tenant_id string) *coreentities.ProviderConfigurations {
	/*
	   Get model provider configurations.
	   Construct ProviderConfiguration objects for each provider
	   Including:
	   1. Basic information of the provider
	   2. Hosting configuration information, including:
	     (1. Whether to enable (support) hosting type, if enabled, the following information exists
	     (2. List of hosting type provider configurations
	         (including quota type, quota limit, current remaining quota, etc.)
	     (3. The current hosting type in use (whether there is a quota or not)
	         paid quotas > provider free quotas > hosting trial quotas
	     (4. Unified credentials for hosting providers
	   3. Custom configuration information, including:
	     (1. Whether to enable (support) custom type, if enabled, the following information exists
	     (2. Custom provider configuration (including credentials)
	     (3. List of custom provider model configurations (including credentials)
	   4. Hosting/custom preferred provider type.
	   Provide methods:
	   - Get the current configuration (including credentials)
	   - Get the availability and status of the hosting configuration: active available,
	     quota_exceeded insufficient quota, unsupported hosting
	   - Get the availability of custom configuration
	     Custom provider available conditions:
	     (1. custom provider credentials available
	     (2. at least one custom model credentials available
	   - Verify, update, and delete custom provider configuration
	   - Verify, update, and delete custom provider model configuration
	   - Get the list of available models (optional provider filtering, model type filtering)
	     Append custom provider models to the list
	   - Get provider instance
	   - Switch selection priority
	   :param tenant_id:
	   :return:
	*/
	// Get all provider records of the workspace
	provider_name_to_provider_records_dict := pm.get_all_providers(tenant_id)
	mlog.Debugf("------provider_name_to_provider_records_dict=%#v", provider_name_to_provider_records_dict)
	// Initialize trial provider records if not exist
	provider_name_to_provider_records_dict = pm.init_trial_provider_records(tenant_id, provider_name_to_provider_records_dict)
	// Get all provider model records of the workspace
	provider_name_to_provider_model_records_dict := pm.get_all_provider_models(tenant_id)
	// Get all provider entities
	provider_entities := (&modelproviders.ModelProviderFactory{}).GetProviders()
	// Get All preferred provider types of the workspace
	provider_name_to_preferred_model_provider_records_dict := pm.get_all_preferred_model_providers(tenant_id)
	// Get All provider model settings
	provider_name_to_provider_model_settings_dict := pm.get_all_provider_model_settings(tenant_id)
	// Get All load balancing configs
	provider_name_to_provider_load_balancing_model_configs_dict := pm.get_all_provider_load_balancing_configs(tenant_id)
	provider_configurations := &coreentities.ProviderConfigurations{TenantID: tenant_id, Configurations: make(map[string]*coreentities.ProviderConfiguration)}
	// Construct ProviderConfiguration objects for each provider
	for _, provider_entity := range provider_entities {
		// handle include, exclude
		// if is_filtered(
		//     include_set=cast(set[str], dify_config.POSITION_PROVIDER_INCLUDES_SET),
		//     exclude_set=cast(set[str], dify_config.POSITION_PROVIDER_EXCLUDES_SET),
		//     data=provider_entity,
		//     name_func=lambda x: x.Provider,
		// ){
		//     continue
		// }
		provider_name := provider_entity.Provider
		mlog.Debugf("------provider_name_to_provider_records_dict=%#v", provider_name_to_provider_records_dict)
		provider_records := provider_name_to_provider_records_dict[provider_entity.Provider]

		provider_model_records := provider_name_to_provider_model_records_dict[provider_entity.Provider]
		// Convert to custom configuration
		custom_configuration := pm.to_custom_configuration(
			tenant_id, provider_entity, provider_records, provider_model_records,
		)
		mlog.Debugf("------custom_configuration=%#v", custom_configuration)
		mlog.Debugf("------custom_configuration.Provider=%#v", custom_configuration.Provider)
		// Convert to system configuration
		system_configuration := pm.to_system_configuration(tenant_id, provider_entity, provider_records)
		// Get preferred provider type
		preferred_provider_type_record := provider_name_to_preferred_model_provider_records_dict[provider_name]
		var preferred_provider_type models.ProviderType
		if preferred_provider_type_record != nil {
			preferred_provider_type = preferred_provider_type_record.PreferredProviderType
		} else if custom_configuration.Provider != nil || len(custom_configuration.Models) > 0 {
			preferred_provider_type = models.Provider_CUSTOM
		} else if system_configuration.Enabled {
			preferred_provider_type = models.Provider_SYSTEM
		} else {
			preferred_provider_type = models.Provider_CUSTOM
		}
		using_provider_type := preferred_provider_type
		has_valid_quota := true
		for _, quota_conf := range system_configuration.QuotaConfigurations {
			if !quota_conf.IsValid {
				has_valid_quota = false
				break
			}
		}

		if preferred_provider_type == models.Provider_SYSTEM {
			if !system_configuration.Enabled || !has_valid_quota {
				using_provider_type = models.Provider_CUSTOM
			}
		} else {
			if custom_configuration.Provider == nil && len(custom_configuration.Models) == 0 {
				if system_configuration.Enabled && has_valid_quota {
					using_provider_type = models.Provider_SYSTEM
				}
			}
		}
		// Get provider load balancing configs
		provider_model_settings := provider_name_to_provider_model_settings_dict[provider_name]
		// Get provider load balancing configs
		provider_load_balancing_configs := provider_name_to_provider_load_balancing_model_configs_dict[provider_name]
		// Convert to model settings
		model_settings := pm.to_model_settings(
			provider_entity,
			provider_model_settings,
			provider_load_balancing_configs,
		)
		provider_configuration := &coreentities.ProviderConfiguration{
			TenantID:              tenant_id,
			Provider:              provider_entity,
			PreferredProviderType: preferred_provider_type,
			UsingProviderType:     using_provider_type,
			SystemConfiguration:   system_configuration,
			CustomConfiguration:   custom_configuration,
			ModelSettings:         model_settings,
		}
		provider_configurations.Configurations[provider_name] = provider_configuration
	}
	// Return the encapsulated object
	return provider_configurations
}

func (pm *ProviderManager) UpdateDefaultModelRecord(
	tenant_id string, provider string, model string, model_type modelruntimeentities.ModelType,
) *models.TenantDefaultModel {
	/*
		Update default model record.

		:param tenant_id: workspace id
		:param model_type: model type
		:param provider: provider name
		:param model: model name
		:return:
	*/
	provider_configurations := pm.GetConfigurations(tenant_id)

	// get provider instance
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider)))
	}

	// get available models from provider_configurations
	available_models := (&ProviderConfigurationsManager{}).GetModels(provider_configurations, "", model_type, true)

	// check if the model is exist in available models
	model_names := []string{}
	for _, model := range available_models {
		model_names = append(model_names, model.Model)
	}
	if !slices.Contains(model_names, model) {
		panic(exceptions.NewValueError(fmt.Sprintf("Model %s does not exist.", model)))
	}
	// Get the list of available models from get_configurations and check if it is LLM
	default_model := new(models.TenantDefaultModel)
	err := dbengine.Instance().DB.Debug().Model(&models.TenantDefaultModel{}).Where("tenant_id = ? and model_type =?", tenant_id, model_type).First(default_model).Error
	if err != nil {
		mlog.Warningf("get TenantDefaultModel failded:%v", err)
		default_model = nil
	}

	// create or update TenantDefaultModel record
	if default_model != nil {
		// update default model
		default_model.ProviderName = provider
		default_model.ModelName = model

	} else {
		// create default model
		default_model = &models.TenantDefaultModel{
			ID:           uuid.NewV4().String(),
			TenantID:     tenant_id,
			ModelType:    model_type,
			ProviderName: provider,
			ModelName:    model,
		}
		dbengine.Instance().DB.Create(default_model)
	}
	return default_model
}

func (pm *ProviderManager) GetFirstProviderFirstModel(tenant_id string, model_type modelruntimeentities.ModelType) (string, string) {
	/*
	   Get names of first model and its provider

	   :param tenant_id: workspace id
	   :param model_type: model type
	   :return: provider name, model name
	*/
	provider_configurations := pm.GetConfigurations(tenant_id)

	// get available models from provider_configurations
	mods := (&ProviderConfigurationsManager{}).GetModels(provider_configurations, "", model_type, false)

	return mods[0].Provider.Provider, mods[0].Model
}

func (pm *ProviderManager) GetDefaultModel(tenant_id string, model_type modelruntimeentities.ModelType) *coreentities.DefaultModelEntity {
	/*
	   Get default model.

	   :param tenant_id: workspace id
	   :param model_type: model type
	   :return:
	*/
	// Get the corresponding TenantDefaultModel record
	default_model := new(models.TenantDefaultModel)
	err := dbengine.Instance().DB.Model(&models.TenantDefaultModel{}).Where("tenant_id = ? and model_type = ?", tenant_id, model_type).First(default_model).Error
	if err != nil {
		mlog.Warningf("get TenantDefaultModel by(tenant_id = %s and model_type = %s) failed:%v", tenant_id, model_type, err)
		default_model = nil
	}

	// If it does not exist, get the first available provider model from get_configurations
	// and update the TenantDefaultModel record
	if default_model == nil {
		// Get provider configurations
		provider_configurations := pm.GetConfigurations(tenant_id)

		// get available models from provider_configurations
		available_models := (&ProviderConfigurationsManager{}).GetModels(provider_configurations, "", model_type, true)

		if len(available_models) > 0 {
			var available_model *coreentities.ModelWithProviderEntity = available_models[0]
			for _, model := range available_models {
				if model.Model == "gpt-4" {
					available_model = model
					break
				}
			}

			default_model = &models.TenantDefaultModel{
				ID:           uuid.NewV4().String(),
				TenantID:     tenant_id,
				ModelType:    model_type,
				ProviderName: available_model.Provider.Provider,
				ModelName:    available_model.Model,
			}
			dbengine.Instance().DB.Create(default_model)
		}
	}
	if default_model == nil {
		return nil
	}
	provider_instance := (&modelproviders.ModelProviderFactory{}).GetProviderInstance(default_model.ProviderName)

	provider_schema := provider_instance.GetProviderSchema(provider_instance.ProviderName())

	return &coreentities.DefaultModelEntity{
		Model:     default_model.ModelName,
		ModelType: string(model_type),
		Provider: &coreentities.DefaultModelProviderEntity{
			Provider:            provider_schema.Provider,
			Label:               provider_schema.Label,
			IconSmall:           provider_schema.IconSmall,
			IconLarge:           provider_schema.IconLarge,
			SupportedModelTypes: provider_schema.SupportedModelTypes,
		},
	}
}

func (pm *ProviderManager) GetProviderModelBundle(tenant_id string, provider string, model_type modelruntimeentities.ModelType) *coreentities.ProviderModelBundle {
	/*
	   Get provider model bundle.
	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model_type: model type
	   :return:
	*/
	provider_configurations := pm.GetConfigurations(tenant_id)

	// get provider instance
	if _, ok := provider_configurations.Configurations[provider]; !ok || provider_configurations.Configurations[provider] == nil {
		mlog.Errorf("Provider %s does not exist.", provider)
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not exist.", provider)))
	}
	provider_configuration := provider_configurations.Configurations[provider]

	provider_instance := (&modelproviders.ModelProviderFactory{}).GetProviderInstance(provider_configuration.Provider.Provider)

	model_type_instance := provider_instance.GetModelInstance(model_type)
	if model_type_instance == nil {
		mlog.Errorf("provider(%s) don't exist model_type(%v) instance", provider_configuration.Provider.Provider, model_type)
		panic(exceptions.NewValueError(fmt.Sprintf("provider(%s) don't exist model_type(%v) instance", provider_configuration.Provider.Provider, model_type)))
	}

	return &coreentities.ProviderModelBundle{
		Configuration:     provider_configuration,
		ProviderInstance:  provider_instance,
		ModelTypeInstance: model_type_instance,
	}
}
