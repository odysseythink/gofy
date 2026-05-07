package services

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	datamanager "github.com/odysseythink/gofy/backend/core/manageres/data_manager"
	providermanager "github.com/odysseythink/gofy/backend/core/manageres/provider_manager"
	modelproviders "github.com/odysseythink/gofy/backend/core/model_runtime/model_provides"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	coreentities "github.com/odysseythink/gofy/backend/entities/core"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

type ModelLoadBalancingService struct {
}

func (service *ModelLoadBalancingService) EnableModelLoadBalancing(tenant_id string, provider string, model string, model_type modelruntimeenumtypes.ModelType) error {
	/*
	   enable model load balancing.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model: model name
	   :param model_type: model type
	   :return:
	*/

	// Enable model load balancing

	_, err := ServiceGroupApp.ProviderConfiguration.EnableModelLoadBalancing(tenant_id, provider, model_type, model)
	return err
}

func (service *ModelLoadBalancingService) DisableModelLoadBalancing(tenant_id string, provider string, model string, model_type modelruntimeenumtypes.ModelType) error {
	/*
	   disable model load balancing.

	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model: model name
	   :param model_type: model type
	   :return:
	*/

	// disable model load balancing
	_, err := ServiceGroupApp.ProviderConfiguration.DisableModelLoadBalancing(tenant_id, provider, model_type, model)
	return err
}

func (service *ModelLoadBalancingService) get_credential_schema(
	provider_configuration *coreentities.ProviderConfiguration,
) /* Union[ModelCredentialSchema, ProviderCredentialSchema]*/ (any, error) {
	/*Get form schemas.*/
	if provider_configuration.Provider.ModelCredentialSchema != nil {
		return provider_configuration.Provider.ModelCredentialSchema, nil
	} else if provider_configuration.Provider.ProviderCredentialSchema != nil {
		return provider_configuration.Provider.ProviderCredentialSchema, nil
	} else {
		return nil, exceptions.NewValueError("No credential schema found")
	}
}
func (service *ModelLoadBalancingService) clear_credentials_cache(tenant_id string, config_id string) {
	/*
		Clear credentials cache.
		:param tenant_id: workspace id
		:param config_id: load balancing config id
		:return:
	*/
	provider_model_credentials_cache := datamanager.NewProviderCredentialsCache(
		tenant_id, config_id, datamanager.ProviderCredentialsCache_LOAD_BALANCING_MODEL,
	)

	provider_model_credentials_cache.Delete()
}

func (service *ModelLoadBalancingService) custom_credentials_validate(
	tenant_id string,
	provider_configuration *coreentities.ProviderConfiguration,
	model_type modelruntimeenumtypes.ModelType,
	model string,
	credentials map[string]any,
	load_balancing_model_config *models.LoadBalancingModelConfig,
	validate bool, /* = True*/
) map[string]any {
	/*
		Validate custom credentials.
		:param tenant_id: workspace id
		:param provider_configuration: provider configuration
		:param model_type: model type
		:param model: model name
		:param credentials: credentials
		:param load_balancing_model_config: load balancing model config
		:param validate: validate credentials
		:return:
	*/
	// Get credential form schemas from model credential schema or provider credential schema
	credential_schemas, err := service.get_credential_schema(provider_configuration)
	if err != nil {
		mlog.Error("get_credential_schema failed:", err)
		return nil
	}
	// var provider_credential_secret_variables []string
	// switch real_credential_schemas := credential_schemas.(type) {
	// case *modelruntimeentities.ModelCredentialSchema:
	// 	// Get provider credential secret variables
	// 	provider_credential_secret_variables = provider_configuration.ExtractSecretVariables(real_credential_schemas.CredentialFormSchemas)
	// case *modelruntimeentities.ProviderCredentialSchema:
	// 	provider_credential_secret_variables = provider_configuration.ExtractSecretVariables(real_credential_schemas.CredentialFormSchemas)
	// }

	if load_balancing_model_config != nil {
		var original_credentials map[string]any
		if load_balancing_model_config.EncryptedConfig != "" {
			err = json.Unmarshal([]byte(load_balancing_model_config.EncryptedConfig), &original_credentials)
			if err != nil {
				mlog.Warningf("json unmarshal(%s) failed:%v", load_balancing_model_config.EncryptedConfig, err)
				original_credentials = map[string]any{}
			}
		} else {
			original_credentials = map[string]any{}
		}

		// // encrypt credentials
		// for key, value := range credentials{
		// 	if slices.Contains(provider_credential_secret_variables, key) {
		// 		// if send [__HIDDEN__] in secret input, it will be same as original value
		// 		if value == constants.HIDDEN_VALUE and key in original_credentials{
		// 			credentials[key] = encrypter.decrypt_token(tenant_id, original_credentials[key])
		// 		}
		// 	}
		// }
	}
	if validate {
		switch credential_schemas.(type) {
		case *modelruntimeentities.ModelCredentialSchema:
			credentials = (&modelproviders.ModelProviderFactory{}).ModelCredentialsValidate(
				provider_configuration.Provider.Provider,
				model_type,
				model,
				credentials,
			)

		case *modelruntimeentities.ProviderCredentialSchema:
			credentials = (&modelproviders.ModelProviderFactory{}).ProviderCredentialsValidate(
				provider_configuration.Provider.Provider, credentials,
			)
		}
	}
	// for key, value := range credentials {
	// 	if key in provider_credential_secret_variables{
	// 		credentials[key] = encrypter.encrypt_token(tenant_id, value)
	// 	}
	// }
	return credentials
}

func (service *ModelLoadBalancingService) UpdateLoadBalancingConfigs(
	tenant_id string, provider string, model string, model_type modelruntimeenumtypes.ModelType, configs []map[string]any,
) {
	/*
		Update load balancing configurations.
		:param tenant_id: workspace id
		:param provider: provider name
		:param model: model name
		:param model_type: model type
		:param configs: load balancing configs
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

	var current_load_balancing_configs []*models.LoadBalancingModelConfig
	err := dbengine.Instance().DB.Model(&models.LoadBalancingModelConfig{}).Where("tenant_id = ? and provider_name = ? and model_type = ? and model_name = ?", tenant_id, provider_configuration.Provider.Provider, model_type, model).Find(&current_load_balancing_configs).Error
	if err != nil {
		mlog.Errorf("get LoadBalancingModelConfig failed:%v", err)
		current_load_balancing_configs = make([]*models.LoadBalancingModelConfig, 0)
	}

	// id as key, config as value
	current_load_balancing_configs_dict := map[string]*models.LoadBalancingModelConfig{}
	for _, config := range current_load_balancing_configs {
		current_load_balancing_configs_dict[config.ID] = config
	}

	updated_config_ids := []string{}

	for _, config := range configs {
		config_id := ""
		if _, ok := config["id"]; ok {
			if _, ok := config["id"].(string); ok {
				config_id = config["id"].(string)
			}
		}
		name := ""
		if _, ok := config["name"]; ok {
			if _, ok := config["name"].(string); ok {
				name = config["name"].(string)
			} else {
				mlog.Error("Invalid load balancing config name")
				panic(exceptions.NewValueError("Invalid load balancing config name"))
			}
		} else {
			mlog.Error("Invalid load balancing config name")
			panic(exceptions.NewValueError("Invalid load balancing config name"))
		}
		enabled := false
		if _, ok := config["enabled"]; ok {
			if _, ok := config["enabled"].(bool); ok {
				enabled = config["enabled"].(bool)
			} else {
				mlog.Error("Invalid load balancing config enabled")
				panic(exceptions.NewValueError("Invalid load balancing config enabled"))
			}
		} else {
			mlog.Error("Invalid load balancing config enabled")
			panic(exceptions.NewValueError("Invalid load balancing config enabled"))
		}
		credentials := map[string]any{}
		if _, ok := config["credentials"]; ok {
			if _, ok := config["credentials"].(map[string]any); ok {
				credentials = config["credentials"].(map[string]any)
			} else {
				mlog.Error("Invalid load balancing config credentials")
				panic(exceptions.NewValueError("Invalid load balancing config credentials"))
			}
		}

		// is config exists
		if config_id != "" {

			if _, ok := current_load_balancing_configs_dict[config_id]; !ok {
				mlog.Error("Invalid load balancing config id: " + config_id)
				panic(exceptions.NewValueError("Invalid load balancing config id: " + config_id))
			}
			updated_config_ids = append(updated_config_ids, config_id)
			slices.Sort(updated_config_ids)
			updated_config_ids = slices.Compact(updated_config_ids)

			load_balancing_config := current_load_balancing_configs_dict[config_id]

			// check duplicate name
			for _, current_load_balancing_config := range current_load_balancing_configs {
				if current_load_balancing_config.ID != config_id && current_load_balancing_config.Name == name {
					mlog.Errorf("Load balancing config name %s already exists", name)
					panic(exceptions.NewValueError(fmt.Sprintf("Load balancing config name %s already exists", name)))
				}
			}
			if len(credentials) > 0 {
				// validate custom provider config
				credentials = service.custom_credentials_validate(
					tenant_id,
					provider_configuration,
					model_type,
					model,
					credentials,
					load_balancing_config,
					false,
				)

				// update load balancing config
				bindata, _ := json.Marshal(credentials)
				load_balancing_config.EncryptedConfig = string(bindata)
			}
			load_balancing_config.Name = name
			load_balancing_config.Enabled = enabled
			now := time.Now()
			load_balancing_config.UpdatedAt = &now
			dbengine.Instance().DB.Save(load_balancing_config)
			service.clear_credentials_cache(tenant_id, config_id)
		} else {
			// create load balancing config
			if name == "__inherit__" {
				panic(exceptions.NewValueError("Invalid load balancing config name"))
			}
			// check duplicate name
			for _, current_load_balancing_config := range current_load_balancing_configs {
				if current_load_balancing_config.Name == name {
					panic(exceptions.NewValueError(fmt.Sprintf("Load balancing config name %s already exists", name)))
				}
			}
			if len(credentials) == 0 {
				panic(exceptions.NewValueError("Invalid load balancing config credentials"))
			}

			// validate custom provider config
			credentials = service.custom_credentials_validate(
				tenant_id,
				provider_configuration,
				model_type,
				model,
				credentials,
				nil,
				false,
			)
			bindata, _ := json.Marshal(credentials)
			// create load balancing config
			load_balancing_model_config := &models.LoadBalancingModelConfig{
				ID:              uuid.NewV4().String(),
				TenantID:        tenant_id,
				ProviderName:    provider_configuration.Provider.Provider,
				ModelType:       string(model_type),
				ModelName:       model,
				Name:            name,
				EncryptedConfig: string(bindata),
			}
			dbengine.Instance().DB.Create(load_balancing_model_config)
		}
	}
	// get deleted config ids
	deleted_config_ids := []string{}
	for k := range current_load_balancing_configs_dict {
		if !slices.Contains(updated_config_ids, k) {
			deleted_config_ids = append(deleted_config_ids, k)
		}
	}
	for _, config_id := range deleted_config_ids {
		dbengine.Instance().DB.Delete(current_load_balancing_configs_dict[config_id])
		service.clear_credentials_cache(tenant_id, config_id)
	}
}
