package schemavalidators

import (
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	"mlib.com/mlog"
)

type ModelCredentialSchemaValidator struct {
	*CommonValidator
	modelruntimeentities.ModelType
	ModelCredentialSchema *modelruntimeentities.ModelCredentialSchema
}

func NewModelCredentialSchemaValidator(model_type modelruntimeentities.ModelType, model_credential_schema *modelruntimeentities.ModelCredentialSchema) *ModelCredentialSchemaValidator {
	return &ModelCredentialSchemaValidator{
		ModelType:             model_type,
		ModelCredentialSchema: model_credential_schema,
	}
}

func (csv *ModelCredentialSchemaValidator) ValidateAndFilter(credentials map[string]any) map[string]any {
	/*
	   Validate model credentials

	   :param credentials: model credentials
	   :return: filtered credentials
	*/

	if csv.ModelCredentialSchema == nil {
		mlog.Error("Model credential schema is None")
		panic(exceptions.NewValueError("Model credential schema is None"))
	}
	// get the credential_form_schemas in provider_credential_schema
	credential_form_schemas := csv.ModelCredentialSchema.CredentialFormSchemas

	credentials["__model_type"] = csv.ModelType

	return csv.ValidateAndFilterCredentialFormSchemas(credential_form_schemas, credentials)
}
