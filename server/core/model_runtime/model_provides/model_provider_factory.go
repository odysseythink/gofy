package modelproviders

import (
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	_ "mlib.com/gofy/server/core/model_runtime/model_provides/tongyi"
	schemavalidators "mlib.com/gofy/server/core/model_runtime/schema_validators"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	"mlib.com/gofy/server/global"
	"mlib.com/mlog"
)

type ModelProviderFactory struct {
}

func (factory *ModelProviderFactory) GetProviders() []*modelruntimeentities.ProviderEntity {
	/*
	   Get all providers
	   :return: list of providers
	*/

	// traverse all model_provider_extensions
	providers := []*modelruntimeentities.ProviderEntity{}
	global.AllModelProviders.Range(func(key any, value any) bool {
		provider_name := key.(string)
		model_provider_instance := value.(modelruntimeentities.ModelProvider)
		// get provider schema
		provider_schema := model_provider_instance.GetProviderSchema(provider_name)
		for _, model_type := range provider_schema.SupportedModelTypes {
			models := model_provider_instance.Models(model_provider_instance, modelruntimeentities.ModelType(model_type))
			if len(models) > 0 {
				provider_schema.Models = append(provider_schema.Models, models...)
			}
		}
		providers = append(providers, provider_schema)
		return true
	})

	// return providers
	return providers
}

func (factory *ModelProviderFactory) GetProviderInstance(provider_name string) modelruntimeentities.ModelProvider {
	/*
	   Get provider instance by provider name
	   :param provider: provider name
	   :return: provider instance
	*/
	// scan all providers
	val, ok := global.AllModelProviders.Load(provider_name)
	if !ok {
		panic(fmt.Errorf("invalid provider: %s", provider_name))
	}
	if _, ok := val.(modelruntimeentities.ModelProvider); !ok {
		global.AllModelProviders.Delete(provider_name)
		panic(fmt.Errorf("local provider: %s must be type of modelruntimeentities.ModelProvider", provider_name))
	}

	return val.(modelruntimeentities.ModelProvider)
}

func (factory *ModelProviderFactory) ProviderCredentialsValidate(provider string, credentials map[string]any) map[string]any {
	/*
	   Validate provider credentials

	   :param provider: provider name
	   :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
	   :return:
	*/
	// get the provider instance
	model_provider_instance := factory.GetProviderInstance(provider)

	// get provider schema
	provider_schema := model_provider_instance.GetProviderSchema(provider)

	// get provider_credential_schema and validate credentials according to the rules
	provider_credential_schema := provider_schema.ProviderCredentialSchema

	if provider_credential_schema == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not have provider_credential_schema", provider)))
	}
	// validate provider credential schema
	validator := schemavalidators.NewProviderCredentialSchemaValidator(provider_credential_schema)
	filtered_credentials := validator.ValidateAndFilter(credentials)

	// validate the credentials, raise exception if validation failed
	model_provider_instance.ValidateProviderCredentials(filtered_credentials)

	return filtered_credentials

}

func (factory *ModelProviderFactory) ModelCredentialsValidate(
	provider string, model_type modelruntimeentities.ModelType, model string, credentials map[string]any,
) map[string]any {
	/*
	   Validate model credentials

	   :param provider: provider name
	   :param model_type: model type
	   :param model: model name
	   :param credentials: model credentials, credentials form defined in `model_credential_schema`.
	   :return:
	*/
	// get the provider instance
	model_provider_instance := factory.GetProviderInstance(provider)

	// get provider schema
	provider_schema := model_provider_instance.GetProviderSchema(provider)

	// get model_credential_schema and validate credentials according to the rules
	model_credential_schema := provider_schema.ModelCredentialSchema

	if model_credential_schema == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("Provider %s does not have model_credential_schema", provider)))
	}
	// validate model credential schema
	validator := schemavalidators.NewModelCredentialSchemaValidator(model_type, model_credential_schema)
	mlog.Debugf("------credentials=%#v", credentials)
	filtered_credentials := validator.ValidateAndFilter(credentials)
	mlog.Debugf("------filtered_credentials=%#v", filtered_credentials)
	// get model instance of the model type
	model_instance := model_provider_instance.GetModelInstance(model_type)

	// call validate_credentials method of model type to validate credentials, raise exception if validation failed
	model_instance.ValidateCredentials(model, filtered_credentials)
	return filtered_credentials
}
