package schemavalidators

import (
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

type ProviderCredentialSchemaValidator struct {
	*CommonValidator
	ProviderCredentialSchema *modelruntimeentities.ProviderCredentialSchema
}

func NewProviderCredentialSchemaValidator(provider_credential_schema *modelruntimeentities.ProviderCredentialSchema) *ProviderCredentialSchemaValidator {
	return &ProviderCredentialSchemaValidator{
		ProviderCredentialSchema: provider_credential_schema,
	}
}

func (csv *ProviderCredentialSchemaValidator) ValidateAndFilter(credentials map[string]any) map[string]any {
	/*
	   Validate provider credentials

	   :param credentials: provider credentials
	   :return: validated provider credentials
	*/
	// get the credential_form_schemas in provider_credential_schema
	credential_form_schemas := csv.ProviderCredentialSchema.CredentialFormSchemas

	return csv.ValidateAndFilterCredentialFormSchemas(credential_form_schemas, credentials)
}
