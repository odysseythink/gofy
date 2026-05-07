package core

import (
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	providerentities "github.com/odysseythink/gofy/backend/entities/provider"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	providerenumtypes "github.com/odysseythink/gofy/backend/enum_types/provider"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/utils/crypt"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

var (
	original_provider_configurate_methods = map[string][]modelruntimeentities.ConfigurateMethod{}
)

type ProviderConfiguration struct {
	/*
	   Model class for provider configuration.
	*/

	TenantID              string                                `json:"tenant_id"`
	Provider              *modelruntimeentities.ProviderEntity  `json:"provider"`
	PreferredProviderType providerenumtypes.ProviderType        `json:"preferred_provider_type"`
	UsingProviderType     providerenumtypes.ProviderType        `json:"using_provider_type"`
	SystemConfiguration   *providerentities.SystemConfiguration `json:"system_configuration"`
	CustomConfiguration   *providerentities.CustomConfiguration `json:"custom_configuration"`
	ModelSettings         []*providerentities.ModelSetting      `json:"model_settings"`

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}

func NewProviderConfiguration() *ProviderConfiguration {
	pc := &ProviderConfiguration{}
	if _, ok := original_provider_configurate_methods[pc.Provider.Provider]; !ok {
		// if pc.provider.provider not in original_provider_configurate_methods:
		original_provider_configurate_methods[pc.Provider.Provider] = make([]modelruntimeentities.ConfigurateMethod, 0)
		original_provider_configurate_methods[pc.Provider.Provider] = append(original_provider_configurate_methods[pc.Provider.Provider], pc.Provider.ConfigurateMethods...)
	}
	if len(original_provider_configurate_methods[pc.Provider.Provider]) == 1 && original_provider_configurate_methods[pc.Provider.Provider][0] == modelruntimeentities.ConfigurateMethod_CUSTOMIZABLE_MODEL {
		matched := false
		for _, quota_configuration := range pc.SystemConfiguration.QuotaConfigurations {
			if len(quota_configuration.RestrictModels) > 0 {
				matched = true
				break
			}
		}
		if matched && !slices.Contains(pc.Provider.ConfigurateMethods, modelruntimeentities.ConfigurateMethod_PREDEFINED_MODEL) {
			pc.Provider.ConfigurateMethods = append(pc.Provider.ConfigurateMethods, modelruntimeentities.ConfigurateMethod_PREDEFINED_MODEL)
		}

	}
	return pc
}

func (pc *ProviderConfiguration) GetCurrentCredentials(model_type modelruntimeenumtypes.ModelType, model string) map[string]any {
	/*
		Get current credentials.

		:param model_type: model type
		:param model: model name
		:return:
	*/
	if len(pc.ModelSettings) > 0 {
		// check if model is disabled by admin
		for _, model_setting := range pc.ModelSettings {
			if model_setting.ModelType == model_type && model_setting.Model == model {
				if !model_setting.Enabled {
					panic(exceptions.NewValueError(fmt.Sprintf("Model %s is disabled.", model)))
				}
			}
		}
	}
	if pc.UsingProviderType == providerenumtypes.Provider_SYSTEM {
		var restrict_models []*providerentities.RestrictModel
		for _, quota_configuration := range pc.SystemConfiguration.QuotaConfigurations {
			if pc.SystemConfiguration.CurrentQuotaType != quota_configuration.QuotaType {
				continue
			}
			restrict_models = quota_configuration.RestrictModels
		}
		if pc.SystemConfiguration.Credentials == nil {
			return nil
		}
		copy_credentials := make(map[string]any)
		for k, v := range pc.SystemConfiguration.Credentials {
			copy_credentials[k] = v
		}
		// copy_credentials = pc.SystemConfiguration.Credentials.copy()
		if len(restrict_models) > 0 {
			for _, restrict_model := range restrict_models {
				if restrict_model.ModelType == model_type && restrict_model.Model == model && restrict_model.BaseModelName != "" {
					copy_credentials["base_model_name"] = restrict_model.BaseModelName
				}
			}
		}
		return copy_credentials
	} else {
		var credentials map[string]any
		if len(pc.CustomConfiguration.Models) > 0 {
			for _, model_configuration := range pc.CustomConfiguration.Models {
				if model_configuration.ModelType == model_type && model_configuration.Model == model {
					credentials = model_configuration.Credentials
					break
				}
			}
		}
		if credentials == nil && pc.CustomConfiguration.Provider != nil {
			credentials = pc.CustomConfiguration.Provider.Credentials
		}
		return credentials
	}
}

func (pc *ProviderConfiguration) GetSystemConfigurationStatus() providerenumtypes.SystemConfigurationStatusType {
	/*
		Get system configuration status.
		:return:
	*/
	if !pc.SystemConfiguration.Enabled {
		return providerenumtypes.SystemConfigurationStatus_UNSUPPORTED
	}
	current_quota_type := pc.SystemConfiguration.CurrentQuotaType
	var current_quota_configuration *providerentities.QuotaConfiguration
	for _, q := range pc.SystemConfiguration.QuotaConfigurations {
		if q.QuotaType == current_quota_type {
			current_quota_configuration = q
			break
		}
	}

	if current_quota_configuration == nil {
		return providerenumtypes.SystemConfigurationStatusType("")
	}
	if current_quota_configuration.IsValid {
		return providerenumtypes.SystemConfigurationStatus_ACTIVE
	} else {
		return providerenumtypes.SystemConfigurationStatus_QUOTA_EXCEEDED
	}

}

func (pc *ProviderConfiguration) IsCustomConfigurationAvailable() bool {
	/*
		Check custom configuration available.
		:return:
	*/

	return pc.CustomConfiguration.Provider != nil || len(pc.CustomConfiguration.Models) > 0
}

func (pc *ProviderConfiguration) ExtractSecretVariables(credential_form_schemas []*modelruntimeentities.CredentialFormSchema) []string {
	/*
		Extract secret input form variables.

		:param credential_form_schemas:
		:return:
	*/
	var secret_input_form_variables []string
	for _, credential_form_schema := range credential_form_schemas {
		if credential_form_schema.Type == modelruntimeentities.Form_SECRET_INPUT {
			if secret_input_form_variables == nil {
				secret_input_form_variables = make([]string, 0)
			}
			secret_input_form_variables = append(secret_input_form_variables, credential_form_schema.Variable)
		}
	}
	return secret_input_form_variables

}

func (pc *ProviderConfiguration) ObfuscatedCredentials(credentials map[string]any, credential_form_schemas []*modelruntimeentities.CredentialFormSchema) map[string]any {
	/*
		Obfuscated credentials.

		:param credentials: credentials
		:param credential_form_schemas: credential form schemas
		:return:
	*/
	// Get Provider credential secret variables
	credential_secret_variables := pc.ExtractSecretVariables(credential_form_schemas)

	// Obfuscate Provider credentials
	for key, value := range credentials {
		if slices.Contains(credential_secret_variables, key) {
			if _, ok := value.(string); ok {
				credentials[key] = crypt.ObfuscatedToken(value.(string))
			}
		}
	}
	return credentials

}

func (pc *ProviderConfiguration) GetCustomCredentials(obfuscated bool) map[string]any {
	/*
		Get custom credentials.

		:param obfuscated: obfuscated secret data in credentials
		:return:
	*/
	mlog.Debugf("------pc.CustomConfiguration=%#v", pc.CustomConfiguration)
	if pc.CustomConfiguration.Provider == nil {
		return nil
	}
	mlog.Debugf("------pc.CustomConfiguration.Provider=%#v", pc.CustomConfiguration.Provider)
	credentials := pc.CustomConfiguration.Provider.Credentials
	if !obfuscated {
		return credentials
	}
	mlog.Debugf("------pc.CustomConfiguration.Provider.Credentials=%#v", pc.CustomConfiguration.Provider.Credentials)
	var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
	if pc.Provider.ProviderCredentialSchema != nil {
		credential_form_schemas = pc.Provider.ProviderCredentialSchema.CredentialFormSchemas
	}
	// Obfuscate credentials
	return pc.ObfuscatedCredentials(credentials, credential_form_schemas)

}

func (pc *ProviderConfiguration) GetCustomModelCredentials(model_type modelruntimeenumtypes.ModelType, model string, obfuscated bool) map[string]any {
	/*
		Get custom model credentials.

		:param model_type: model type
		:param model: model name
		:param obfuscated: obfuscated secret data in credentials
		:return:
	*/
	if len(pc.CustomConfiguration.Models) == 0 {
		return nil
	}
	for _, model_configuration := range pc.CustomConfiguration.Models {
		if model_configuration.ModelType == model_type && model_configuration.Model == model {
			credentials := model_configuration.Credentials
			if !obfuscated {
				return credentials
			}
			var credential_form_schemas []*modelruntimeentities.CredentialFormSchema
			if pc.Provider.ModelCredentialSchema != nil {
				credential_form_schemas = pc.Provider.ModelCredentialSchema.CredentialFormSchemas
			}
			// Obfuscate credentials
			return pc.ObfuscatedCredentials(credentials, credential_form_schemas)
		}
	}
	return nil
}
func (pc *ProviderConfiguration) EnableModel(model_type modelruntimeenumtypes.ModelType, model string) *models.ProviderModelSetting {
	/*
	   Enable model.
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	model_setting := new(models.ProviderModelSetting)
	if err := dbengine.Instance().DB.Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type.ToOriginModelType(), model).First(model_setting).Error; err != nil {
		mlog.Errorf("get ProviderModelSetting failed:%v", err)
		model_setting = nil
	}
	now := time.Now()
	if model_setting != nil {
		model_setting.Enabled = true
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Updates(&models.ProviderModelSetting{ID: model_setting.ID, Enabled: true, UpdatedAt: &now})
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:           uuid.NewV4().String(),
			TenantID:     pc.TenantID,
			ProviderName: pc.Provider.Provider,
			ModelName:    model,
			ModelType:    model_type.ToOriginModelType(),
			Enabled:      true,
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}
		dbengine.Instance().DB.Create(model_setting)
	}
	return model_setting
}

func (pc *ProviderConfiguration) DisableModel(model_type modelruntimeenumtypes.ModelType, model string) *models.ProviderModelSetting {
	/*
	   Disable model.
	   :param model_type: model type
	   :param model: model name
	   :return:
	*/
	model_setting := new(models.ProviderModelSetting)
	if err := dbengine.Instance().DB.Model(&models.ProviderModelSetting{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", pc.TenantID, pc.Provider.Provider, model_type.ToOriginModelType(), model).First(model_setting).Error; err != nil {
		mlog.Errorf("get ProviderModelSetting failed:%v", err)
		model_setting = nil
	}
	now := time.Now()
	if model_setting != nil {
		model_setting.Enabled = true
		model_setting.UpdatedAt = &now
		dbengine.Instance().DB.Updates(&models.ProviderModelSetting{ID: model_setting.ID, Enabled: false, UpdatedAt: &now})
	} else {
		model_setting = &models.ProviderModelSetting{
			ID:           uuid.NewV4().String(),
			TenantID:     pc.TenantID,
			ProviderName: pc.Provider.Provider,
			ModelName:    model,
			ModelType:    model_type.ToOriginModelType(),
			Enabled:      false,
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}
		dbengine.Instance().DB.Create(model_setting)
	}
	return model_setting
}

type ProviderConfigurations struct {
	/*
	   Model class for provider configuration.
	*/

	TenantID       string                            `json:"tenant_id"`
	Configurations map[string]*ProviderConfiguration `json:"configurations"`
}

type ProviderModelBundle struct {
	/*
	   Provider model bundle.
	*/

	Configuration     *ProviderConfiguration
	ProviderInstance  modelruntimeentities.ModelProvider
	ModelTypeInstance modelruntimeentities.AIModeler

	// pydantic configs
	ModelConfig map[string]any // = ConfigDict(arbitrary_types_allowed=True, protected_namespaces=())
}
