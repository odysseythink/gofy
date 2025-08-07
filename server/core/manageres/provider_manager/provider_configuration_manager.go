package providermanager

import (
	"encoding/json"
	"maps"
	"slices"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	datamanager "mlib.com/gofy/server/core/manageres/data_manager"
	modelproviders "mlib.com/gofy/server/core/model_runtime/model_provides"
	dbengine "mlib.com/gofy/server/db_engine"
	coreentities "mlib.com/gofy/server/entities/core"
	modelentities "mlib.com/gofy/server/entities/model"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	providerentities "mlib.com/gofy/server/entities/provider"
	coreenumtypes "mlib.com/gofy/server/enum_types/core"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

type ProviderConfigurationManager struct {
}

func (mgr *ProviderConfigurationManager) CustomCredentialsValidate(pc *coreentities.ProviderConfiguration, credentials map[string]any) (*models.Provider, map[string]any) {
	/*
	   Validate custom credentials.
	   :param credentials: provider credentials
	   :return:
	*/
	// get provider
	provider_record := new(models.Provider)
	err := dbengine.Instance().DB.Debug().Model(&models.Provider{}).Where("tenant_id = ? and provider_name = ? and provider_type = ?", pc.TenantID, pc.Provider.Provider, providerenumtypes.Provider_CUSTOM).First(provider_record).Error
	if err != nil {
		mlog.Warningf("get Provider failded:%v", err)
		provider_record = nil
	}

	// Get provider credential secret variables
	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	if pc.Provider.ProviderCredentialSchema != nil {
		credential_form_schemas = pc.Provider.ProviderCredentialSchema.CredentialFormSchemas
	}
	provider_credential_secret_variables := pc.ExtractSecretVariables(credential_form_schemas)

	if provider_record != nil {
		var original_credentials map[string]any
		// try:
		// fix origin data
		if provider_record.EncryptedConfig != "" {
			if !strings.HasPrefix(provider_record.EncryptedConfig, "{") {
				original_credentials = map[string]any{"openai_api_key": provider_record.EncryptedConfig}
			} else {
				err := json.Unmarshal([]byte(provider_record.EncryptedConfig), &original_credentials)
				if err != nil {
					mlog.Errorf("json unmarshal(%s) failed:%v", provider_record.EncryptedConfig, err)
					original_credentials = map[string]any{}
				}
			}
		} else {
			original_credentials = map[string]any{}
		}
		// except JSONDecodeError:
		//     original_credentials = {}

		// encrypt credentials
		for key, value := range credentials {
			if slices.Contains(provider_credential_secret_variables, key) {
				// if send [__HIDDEN__] in secret input, it will be same as original value
				if _, ok := value.(string); ok && value.(string) == constants.HIDDEN_VALUE {
					if _, ok := original_credentials[key]; ok {
						credentials[key] = original_credentials[key]
					}
				}
			}
		}
	}
	credentials = (&modelproviders.ModelProviderFactory{}).ProviderCredentialsValidate(
		pc.Provider.Provider, credentials,
	)

	// for key, value := range credentials{
	//     if slices.Contains(provider_credential_secret_variables, key){
	//         credentials[key] = encrypter.encrypt_token(self.tenant_id, value)
	// 	}
	// }
	return provider_record, credentials
}

func (mgr *ProviderConfigurationManager) AddOrUpdateCustomCredentials(pc *coreentities.ProviderConfiguration, credentials map[string]any) {
	/*
	   Add or update custom provider credentials.
	   :param credentials:
	   :return:
	*/
	// validate custom provider config
	provider_record, credentials := mgr.CustomCredentialsValidate(pc, credentials)

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
			TenantID:        pc.TenantID,
			ProviderName:    pc.Provider.Provider,
			ProviderType:    providerenumtypes.Provider_CUSTOM,
			EncryptedConfig: string(bindata),
			IsValid:         true,
		}
		dbengine.Instance().DB.Create(provider_record)
	}

	provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
		pc.TenantID,
		provider_record.ID,
		datamanager.ProviderCredentialsCache_PROVIDER,
	)

	provider_model_credentials_cache.Delete()

	mgr.switch_preferred_provider_type(pc, providerenumtypes.Provider_CUSTOM)
}

func (mgr *ProviderConfigurationManager) CustomModelCredentialsValidate(
	pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string, credentials map[string]any,
) (*models.ProviderModel, map[string]any, error) {
	/*
		Validate custom model credentials.

		:param model_type: model type
		:param model: model name
		:param credentials: model credentials
		:return:
	*/
	// get provider model
	provider_model_record := new(models.ProviderModel)
	err := dbengine.Instance().DB.Debug().Model(&models.ProviderModel{}).Where("tenant_id = ? and provider_name = ? and model_name = ? and model_type =?", pc.TenantID, pc.Provider.Provider, model, model_type).First(provider_model_record).Error
	if err != nil {
		mlog.Warningf("get ProviderModel failded:%v", err)
		provider_model_record = nil
	}

	// Get provider credential secret variables
	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	if pc.Provider.ModelCredentialSchema != nil {
		credential_form_schemas = pc.Provider.ModelCredentialSchema.CredentialFormSchemas
	}
	provider_credential_secret_variables := pc.ExtractSecretVariables(credential_form_schemas)

	if provider_model_record != nil {
		// try{
		var original_credentials map[string]any
		err = json.Unmarshal([]byte(provider_model_record.EncryptedConfig), &original_credentials)
		if err != nil {
			mlog.Warningf("json unmarshal failed:%v", err)
		}
		// except JSONDecodeError{
		// 	original_credentials = {}

		// decrypt credentials
		for key, value := range credentials {
			if slices.Contains(provider_credential_secret_variables, key) {
				// if send [__HIDDEN__] in secret input, it will be same as original value
				if _, ok := original_credentials[key]; ok {
					if _, ok := value.(string); ok && value.(string) == constants.HIDDEN_VALUE {
						credentials[key] = original_credentials[key]
					}
				}
			}
		}
	}
	credentials = (&modelproviders.ModelProviderFactory{}).ModelCredentialsValidate(
		pc.Provider.Provider, model_type, model, credentials,
	)

	// for key, value := range credentials{
	// 	if slices.Contains(provider_credential_secret_variables, key) {
	// 		credentials[key] = encrypter.encrypt_token(pc.TenantID, value)
	// 	}
	// }
	return provider_model_record, credentials, nil
}
func (mgr *ProviderConfigurationManager) AddOrUpdateCustomModelCredentials(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string, credentials map[string]any) error {
	/*
		Add or update custom model credentials.

		:param model_type: model type
		:param model: model name
		:param credentials: model credentials
		:return:
	*/

	// validate custom model config
	provider_model_record, credentials, err := mgr.CustomModelCredentialsValidate(pc, model_type, model, credentials)
	if err != nil {
		mlog.Errorf("CustomModelCredentialsValidate failed:%v", err)
		return err
	}
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
			TenantID:        pc.TenantID,
			ProviderName:    pc.Provider.Provider,
			ModelName:       model,
			ModelType:       string(model_type),
			EncryptedConfig: string(bindata),
			IsValid:         true,
		}
		dbengine.Instance().DB.Create(provider_model_record)
	}
	provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
		pc.TenantID,
		provider_model_record.ID,
		datamanager.ProviderCredentialsCache_MODEL,
	)

	provider_model_credentials_cache.Delete()
	return nil
}
func (mgr *ProviderConfigurationManager) DeleteCustomModelCredentials(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string) {
	/*
		Delete custom model credentials.
		:param model_type: model type
		:param model: model name
		:return:
	*/
	// get provider model
	provider_model_record := new(models.ProviderModel)
	err := dbengine.Instance().DB.Model(&models.ProviderModel{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).First(provider_model_record).Error
	if err != nil {
		mlog.Warningf("count ProviderModel failed:%v", err)
		provider_model_record = nil
	}

	// delete provider model
	if provider_model_record != nil {
		dbengine.Instance().DB.Delete(provider_model_record)

		provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
			pc.TenantID,
			provider_model_record.ID,
			datamanager.ProviderCredentialsCache_MODEL,
		)

		provider_model_credentials_cache.Delete()
	}
}
func (mgr *ProviderConfigurationManager) EnableModel(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string) *models.ProviderModelSetting {
	/*
		Enable model.
		:param model_type: model type
		:param model: model name
		:return:
	*/
	model_setting := new(models.ProviderModelSetting)
	err := dbengine.Instance().DB.Debug().Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).First(model_setting).Error
	if err != nil {
		mlog.Warningf("count ProviderModelSetting failed:%v", err)
		model_setting = nil
	}

	if model_setting != nil {
		model_setting.Enabled = true
		now := time.Now()
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Save(model_setting)
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:           uuid.NewV4().String(),
			TenantID:     pc.TenantID,
			ProviderName: pc.Provider.Provider,
			ModelType:    string(model_type),
			ModelName:    model,
			Enabled:      true,
		}
		dbengine.Instance().DB.Create(model_setting)
	}

	return model_setting
}
func (mgr *ProviderConfigurationManager) DisableModel(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string) *models.ProviderModelSetting {
	/*
		Disable model.
		:param model_type: model type
		:param model: model name
		:return:
	*/
	model_setting := new(models.ProviderModelSetting)
	err := dbengine.Instance().DB.Debug().Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).First(model_setting).Error
	if err != nil {
		mlog.Warningf("count ProviderModelSetting failed:%v", err)
		model_setting = nil
	}

	if model_setting != nil {
		model_setting.Enabled = false
		now := time.Now()
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Save(model_setting)
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:           uuid.NewV4().String(),
			TenantID:     pc.TenantID,
			ProviderName: pc.Provider.Provider,
			ModelType:    string(model_type),
			ModelName:    model,
			Enabled:      false,
		}
		dbengine.Instance().DB.Create(model_setting)
	}

	return model_setting
}
func (mgr *ProviderConfigurationManager) GetProviderModelSetting(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string) *models.ProviderModelSetting {
	/*
		Get provider model setting.
		:param model_type: model type
		:param model: model name
		:return:
	*/
	model_setting := new(models.ProviderModelSetting)
	err := dbengine.Instance().DB.Debug().Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).First(model_setting).Error
	if err != nil {
		mlog.Warningf("count ProviderModelSetting failed:%v", err)
		model_setting = nil
	}

	return model_setting
}
func (mgr *ProviderConfigurationManager) EnableModelLoadBalancing(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string) (*models.ProviderModelSetting, error) {
	/*
		Enable model load balancing.
		:param model_type: model type
		:param model: model name
		:return
	*/
	var load_balancing_config_count int64
	err := dbengine.Instance().DB.Debug().Model(&models.LoadBalancingModelConfig{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).Count(&load_balancing_config_count).Error
	if err != nil {
		mlog.Warningf("count LoadBalancingModelConfig failed:%v", err)
	}

	if load_balancing_config_count <= 1 {
		return nil, exceptions.NewValueError("Model load balancing configuration must be more than 1.")
	}
	model_setting := new(models.ProviderModelSetting)
	err = dbengine.Instance().DB.Debug().Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).First(model_setting).Error
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
			TenantID:             pc.TenantID,
			ProviderName:         pc.Provider.Provider,
			ModelType:            string(model_type),
			ModelName:            model,
			LoadBalancingEnabled: true,
		}
		dbengine.Instance().DB.Create(model_setting)
	}
	return model_setting, nil
}
func (mgr *ProviderConfigurationManager) DisableModelLoadBalancing(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string) *models.ProviderModelSetting {
	/*
		Disable model load balancing.
		:param model_type: model type
		:param model: model name
		:return:
	*/
	model_setting := new(models.ProviderModelSetting)
	err := dbengine.Instance().DB.Debug().Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type, model).First(model_setting).Error
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
			TenantID:             pc.TenantID,
			ProviderName:         pc.Provider.Provider,
			ModelType:            string(model_type),
			ModelName:            model,
			LoadBalancingEnabled: false,
		}
		dbengine.Instance().DB.Create(model_setting)
	}

	return model_setting
}
func (mgr *ProviderConfigurationManager) GetProviderInstance(pc *coreentities.ProviderConfiguration) modelruntimeentities.ModelProvider {
	/*
		Get provider instance.
		:return:
	*/
	provider_instance := (&modelproviders.ModelProviderFactory{}).GetProviderInstance(pc.Provider.Provider)
	return provider_instance
}
func (mgr *ProviderConfigurationManager) GetModelTypeInstance(pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	/*
		Get current model type instance.

		:param model_type: model type
		:return:
	*/
	// Get provider instance
	provider_instance := mgr.GetProviderInstance(pc)

	// Get model instance of LLM
	return provider_instance.GetModelInstance(model_type)
}
func (mgr *ProviderConfigurationManager) switch_preferred_provider_type(pc *coreentities.ProviderConfiguration, provider_type providerenumtypes.ProviderType) {
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
	err := dbengine.Instance().DB.Debug().Model(&models.TenantPreferredModelProvider{}).Where("tenant_id = ? and provider_name = ?", pc.TenantID, pc.Provider.Provider).First(preferred_model_provider).Error
	if err != nil {
		mlog.Warningf("count TenantPreferredModelProvider failed:%v", err)
		preferred_model_provider = nil
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

func (mgr *ProviderConfigurationManager) GetProviderModels(
	pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, only_active bool,
) []*modelentities.ModelWithProviderEntity {
	/*
		Get provider models.
		:param model_type: model type
		:param only_active: only active models
		:return:
	*/
	provider_instance := mgr.GetProviderInstance(pc)

	model_types := []modelruntimeenumtypes.ModelType{}
	if string(model_type) != "" {
		model_types = append(model_types, model_type)
	} else {
		provider_entity := provider_instance.GetProviderSchema(provider_instance.ProviderName())
		model_types = provider_entity.SupportedModelTypes
	}
	// Group model settings by model type and model
	model_setting_map := map[modelruntimeenumtypes.ModelType]map[string]*providerentities.ModelSetting{}
	for _, model_setting := range pc.ModelSettings {
		if _, ok := model_setting_map[model_setting.ModelType]; !ok {
			model_setting_map[model_setting.ModelType] = make(map[string]*providerentities.ModelSetting)
		}
		model_setting_map[model_setting.ModelType][model_setting.Model] = model_setting
	}
	var provider_models []*modelentities.ModelWithProviderEntity
	mlog.Debugf("------provider_configuration=%#v", pc)
	if pc.UsingProviderType == providerenumtypes.Provider_SYSTEM {
		provider_models = mgr._get_system_provider_models(
			pc, model_types, provider_instance, model_setting_map,
		)
	} else {
		provider_models = mgr._get_custom_provider_models(
			pc, model_types, provider_instance, model_setting_map,
		)
	}
	if only_active {
		if only_active {
			var new_provider_models []*modelentities.ModelWithProviderEntity
			for _, m := range provider_models {
				if m.Status == coreenumtypes.ModelStatus_ACTIVE {
					new_provider_models = append(new_provider_models, m)
				}
			}
			provider_models = new_provider_models
		}
	}
	// resort provider_models
	sort.Slice(provider_models, func(i int, j int) bool {
		return provider_models[i].ModelType < provider_models[j].ModelType
	})
	return provider_models
}

func (mgr *ProviderConfigurationManager) GetProviderModel(
	pc *coreentities.ProviderConfiguration, model_type modelruntimeenumtypes.ModelType, model string, only_active bool,
) *modelentities.ModelWithProviderEntity {
	/*
		Get provider model.
		:param model_type: model type
		:param model: model name
		:param only_active: return active model only
		:return:
	*/
	provider_models := mgr.GetProviderModels(pc, model_type, only_active)

	for _, provider_model := range provider_models {
		if provider_model.Model == model {
			return provider_model
		}
	}
	return nil
}
func (mgr *ProviderConfigurationManager) _get_system_provider_models(
	provider_configuration *coreentities.ProviderConfiguration,
	model_types []modelruntimeenumtypes.ModelType,
	provider_instance modelruntimeentities.ModelProvider,
	model_setting_map map[modelruntimeenumtypes.ModelType]map[string]*providerentities.ModelSetting,
) []*modelentities.ModelWithProviderEntity {
	/*
		Get system provider models.

		:param model_types: model types
		:param provider_instance: provider instance
		:param model_setting_map: model setting map
		:return:
	*/
	provider_models := []*modelentities.ModelWithProviderEntity{}
	for _, model_type := range model_types {
		for _, m := range provider_instance.Models(provider_instance, model_type) {
			status := coreenumtypes.ModelStatus_ACTIVE
			if _, ok := model_setting_map[m.ModelType]; ok {
				if _, ok := model_setting_map[m.ModelType][m.Model]; ok {
					model_setting := model_setting_map[m.ModelType][m.Model]
					if !model_setting.Enabled {
						status = coreenumtypes.ModelStatus_DISABLED
					}
				}
			}

			provider_models = append(provider_models, &modelentities.ModelWithProviderEntity{
				ProviderModelWithStatusEntity: &modelentities.ProviderModelWithStatusEntity{
					ProviderModel: &modelruntimeentities.ProviderModel{
						Model:           m.Model,
						Label:           m.Label,
						ModelType:       m.ModelType,
						Features:        m.Features,
						FetchFrom:       m.FetchFrom,
						ModelProperties: m.ModelProperties,
						Deprecated:      m.Deprecated,
					},
					Status: status,
				},
				Provider: modelentities.NewSimpleModelProviderEntity(provider_configuration.Provider),
			})
		}
	}
	if _, ok := original_provider_configurate_methods[provider_configuration.Provider.Provider]; !ok {
		original_provider_configurate_methods[provider_configuration.Provider.Provider] = make([]modelruntimeentities.ConfigurateMethod, 0)
		provider_schema := provider_instance.GetProviderSchema(provider_configuration.Provider.Provider)
		original_provider_configurate_methods[provider_configuration.Provider.Provider] = append(original_provider_configurate_methods[provider_configuration.Provider.Provider], provider_schema.ConfigurateMethods...)
	}
	should_use_custom_model := false
	if slices.Compare(original_provider_configurate_methods[provider_configuration.Provider.Provider], []modelruntimeentities.ConfigurateMethod{modelruntimeentities.ConfigurateMethod_CUSTOMIZABLE_MODEL}) == 0 {
		should_use_custom_model = true
	}
	for _, quota_configuration := range provider_configuration.SystemConfiguration.QuotaConfigurations {
		if provider_configuration.SystemConfiguration.CurrentQuotaType != quota_configuration.QuotaType {
			continue
		}
		restrict_models := quota_configuration.RestrictModels
		if len(restrict_models) == 0 {
			break
		}
		if should_use_custom_model {
			if slices.Compare(original_provider_configurate_methods[provider_configuration.Provider.Provider], []modelruntimeentities.ConfigurateMethod{modelruntimeentities.ConfigurateMethod_CUSTOMIZABLE_MODEL}) == 0 {
				// only customizable model
				for _, restrict_model := range restrict_models {
					if len(provider_configuration.SystemConfiguration.Credentials) > 0 {
						copy_credentials := maps.Clone(provider_configuration.SystemConfiguration.Credentials)
						if restrict_model.BaseModelName != "" {
							copy_credentials["base_model_name"] = restrict_model.BaseModelName
						}
						// try:
						model_instance := provider_instance.GetModelInstance(restrict_model.ModelType)
						if model_instance == nil {
							mlog.Error("get custom model schema failed")
							continue
						}
						custom_model_schema := model_instance.GetCustomizableModelSchemaFromCredentials(model_instance, restrict_model.Model, copy_credentials)
						// except Exception as ex:
						// 	logger.warning(f"get custom model schema failed, {ex}")
						// 	continue

						if custom_model_schema == nil {
							continue
						}
						if !slices.Contains(model_types, custom_model_schema.ModelType) {
							continue
						}
						status := coreenumtypes.ModelStatus_ACTIVE
						if _, ok := model_setting_map[custom_model_schema.ModelType]; ok {
							if _, ok := model_setting_map[custom_model_schema.ModelType][custom_model_schema.Model]; ok {
								model_setting := model_setting_map[custom_model_schema.ModelType][custom_model_schema.Model]
								if !model_setting.Enabled {
									status = coreenumtypes.ModelStatus_DISABLED
								}
							}
						}
						provider_models = append(provider_models, &modelentities.ModelWithProviderEntity{
							ProviderModelWithStatusEntity: &modelentities.ProviderModelWithStatusEntity{
								ProviderModel: &modelruntimeentities.ProviderModel{
									Model:           custom_model_schema.Model,
									Label:           custom_model_schema.Label,
									ModelType:       custom_model_schema.ModelType,
									Features:        custom_model_schema.Features,
									FetchFrom:       modelruntimeenumtypes.FetchFrom_PREDEFINED_MODEL,
									ModelProperties: custom_model_schema.ModelProperties,
									Deprecated:      custom_model_schema.Deprecated,
								},
								Status: status,
							},
							Provider: modelentities.NewSimpleModelProviderEntity(provider_configuration.Provider),
						})
					}
				}
			}
		}

		// if llm name not in restricted llm list, remove it
		restrict_model_names := []string{}
		for _, rm := range restrict_models {
			restrict_model_names = append(restrict_model_names, rm.Model)
		}
		for idx, model := range provider_models {
			if model.ModelType == modelruntimeenumtypes.Model_LLM && !slices.Contains(restrict_model_names, model.Model) {
				model.Status = coreenumtypes.ModelStatus_NO_PERMISSION
			} else if !quota_configuration.IsValid {
				model.Status = coreenumtypes.ModelStatus_QUOTA_EXCEEDED
			}
			provider_models[idx] = model
		}
	}

	return provider_models
}

func (pm *ProviderConfigurationManager) _get_custom_provider_models(
	provider_configuration *coreentities.ProviderConfiguration,
	model_types []modelruntimeenumtypes.ModelType,
	provider_instance modelruntimeentities.ModelProvider,
	model_setting_map map[modelruntimeenumtypes.ModelType]map[string]*providerentities.ModelSetting,
) []*modelentities.ModelWithProviderEntity {
	/*
		Get custom provider models.

		:param model_types: model types
		:param provider_instance: provider instance
		:param model_setting_map: model setting map
		:return:
	*/
	provider_models := []*modelentities.ModelWithProviderEntity{}

	credentials := map[string]any{}
	if provider_configuration.CustomConfiguration.Provider != nil {
		credentials = provider_configuration.CustomConfiguration.Provider.Credentials
	}
	for _, model_type := range model_types {
		if !slices.Contains(provider_configuration.Provider.SupportedModelTypes, model_type) {
			continue
		}
		models := provider_instance.Models(provider_instance, model_type)
		for _, m := range models {
			status := coreenumtypes.ModelStatus_NO_CONFIGURE
			if len(credentials) > 0 {
				status = coreenumtypes.ModelStatus_ACTIVE
			}
			load_balancing_enabled := false
			if _, ok := model_setting_map[m.ModelType]; ok {
				if _, ok := model_setting_map[m.ModelType][m.Model]; ok {
					model_setting := model_setting_map[m.ModelType][m.Model]
					if !model_setting.Enabled {
						status = coreenumtypes.ModelStatus_DISABLED
					}
					if len(model_setting.LoadBalancingConfigs) > 1 {
						load_balancing_enabled = true
					}
				}
			}
			provider_models = append(provider_models, &modelentities.ModelWithProviderEntity{
				ProviderModelWithStatusEntity: &modelentities.ProviderModelWithStatusEntity{
					ProviderModel: &modelruntimeentities.ProviderModel{
						Model:           m.Model,
						Label:           m.Label,
						ModelType:       m.ModelType,
						Features:        m.Features,
						FetchFrom:       m.FetchFrom,
						ModelProperties: m.ModelProperties,
						Deprecated:      m.Deprecated,
					},
					Status:               status,
					LoadBalancingEnabled: load_balancing_enabled,
				},
				Provider: modelentities.NewSimpleModelProviderEntity(provider_configuration.Provider),
			})
		}
	}
	// custom models
	for _, model_configuration := range provider_configuration.CustomConfiguration.Models {
		if !slices.Contains(model_types, model_configuration.ModelType) {
			continue
		}
		is_continue := false
		var custom_model_schema *modelruntimeentities.AIModelEntity
		func() {
			defer func() {
				if r := recover(); r != nil {
					mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
					if exp, ok := r.(error); ok {
						mlog.Warningf("get custom model schema failed, %v", exp)
						is_continue = true
					} else {
						panic(r)
					}
				}
			}()
			model_instance := provider_instance.GetModelInstance(model_configuration.ModelType)
			if model_instance == nil {
				mlog.Error("get custom model schema failed")
				is_continue = true
			}
			custom_model_schema = model_instance.GetCustomizableModelSchemaFromCredentials(model_instance, model_configuration.Model, model_configuration.Credentials)
		}()
		if is_continue {
			continue
		}

		if custom_model_schema == nil {
			continue
		}
		status := coreenumtypes.ModelStatus_ACTIVE
		load_balancing_enabled := false
		if _, ok := model_setting_map[custom_model_schema.ModelType]; ok {
			if _, ok := model_setting_map[custom_model_schema.ModelType][custom_model_schema.Model]; ok {
				model_setting := model_setting_map[custom_model_schema.ModelType][custom_model_schema.Model]
				if !model_setting.Enabled {
					status = coreenumtypes.ModelStatus_DISABLED
				}
				if len(model_setting.LoadBalancingConfigs) > 1 {
					load_balancing_enabled = true
				}
			}
		}
		provider_models = append(provider_models, &modelentities.ModelWithProviderEntity{
			ProviderModelWithStatusEntity: &modelentities.ProviderModelWithStatusEntity{
				ProviderModel: &modelruntimeentities.ProviderModel{
					Model:           custom_model_schema.Model,
					Label:           custom_model_schema.Label,
					ModelType:       custom_model_schema.ModelType,
					Features:        custom_model_schema.Features,
					FetchFrom:       custom_model_schema.FetchFrom,
					ModelProperties: custom_model_schema.ModelProperties,
					Deprecated:      custom_model_schema.Deprecated,
				},
				Status:               status,
				LoadBalancingEnabled: load_balancing_enabled,
			},
			Provider: modelentities.NewSimpleModelProviderEntity(provider_configuration.Provider),
		})
	}
	return provider_models
}
