package schemavalidators

import (
	"fmt"
	"slices"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

type CommonValidator struct{}

func (cv *CommonValidator) ValidateAndFilterCredentialFormSchemas(credential_form_schemas []*modelruntimeentities.CredentialFormSchema, credentials map[string]any) map[string]any {
	mlog.Debugf("------credential_form_schemas=%#v", credential_form_schemas)
	mlog.Debugf("------credentials=%#v", credentials)
	need_validate_credential_form_schema_map := make(map[string]*modelruntimeentities.CredentialFormSchema)

	for _, credential_form_schema := range credential_form_schemas {
		mlog.Debugf("------credential_form_schema=%#v", credential_form_schema)
		if len(credential_form_schema.ShowOn) == 0 {
			need_validate_credential_form_schema_map[credential_form_schema.Variable] = credential_form_schema
			continue
		}

		all_show_on_match := true
		for _, show_on_object := range credential_form_schema.ShowOn {
			if _, ok := credentials[show_on_object.Variable]; !ok {
				all_show_on_match = false
				break
			} else {
				if _, ok := credentials[show_on_object.Variable].(string); !ok {
					all_show_on_match = false
					break
				}
			}

			if credentials[show_on_object.Variable].(string) != show_on_object.Value {
				all_show_on_match = false
				break
			}
		}
		if all_show_on_match {
			need_validate_credential_form_schema_map[credential_form_schema.Variable] = credential_form_schema
		}
	}

	validated_credentials := make(map[string]any)
	for _, credential_form_schema := range need_validate_credential_form_schema_map {
		mlog.Debugf("------credential_form_schema=%#v", credential_form_schema)
		result := cv.ValidateCredentialFormSchema(credential_form_schema, credentials)
		if result != nil {
			validated_credentials[credential_form_schema.Variable] = result
		}
	}

	return validated_credentials
}

func (cv *CommonValidator) ValidateCredentialFormSchema(credential_form_schema *modelruntimeentities.CredentialFormSchema, credentials map[string]any) any {
	if _, ok := credentials[credential_form_schema.Variable]; !ok || credentials[credential_form_schema.Variable] == nil {
		if credential_form_schema.Required {
			mlog.Errorf("Variable %s is required", credential_form_schema.Variable)
			panic(exceptions.NewValueError(fmt.Sprintf("Variable %s is required", credential_form_schema.Variable)))
		}
		if credential_form_schema.Default != "" {
			return credential_form_schema.Default
		}
		return nil
	}

	if _, ok := credentials[credential_form_schema.Variable].(string); !ok {
		mlog.Errorf("Variable %s should be string", credential_form_schema.Variable)
		panic(exceptions.NewValueError(fmt.Sprintf("Variable %s should be string", credential_form_schema.Variable)))
	}
	value := credentials[credential_form_schema.Variable].(string)
	if credential_form_schema.MaxLength > 0 && len(value) > credential_form_schema.MaxLength {
		mlog.Errorf("Variable %s length should not greater than %d", credential_form_schema.Variable, credential_form_schema.MaxLength)
		panic(exceptions.NewValueError(fmt.Sprintf("Variable %s length should not greater than %d", credential_form_schema.Variable, credential_form_schema.MaxLength)))
	}

	switch credential_form_schema.Type {
	case modelruntimeentities.Form_SELECT, modelruntimeentities.Form_RADIO:
		if len(credential_form_schema.Options) > 0 {
			if !slices.ContainsFunc(credential_form_schema.Options, func(option *modelruntimeentities.FormOption) bool {
				return value == option.Value
			}) {
				mlog.Errorf("Variable %s is not in options", credential_form_schema.Variable)
				panic(exceptions.NewValueError(fmt.Sprintf("Variable %s is not in options", credential_form_schema.Variable)))
			}
		}
	case modelruntimeentities.Form_SWITCH:
		if !slices.Contains([]string{"true", "false"}, strings.ToLower(value)) {
			mlog.Errorf("Variable %s is not in options", credential_form_schema.Variable)
			panic(exceptions.NewValueError(fmt.Sprintf("Variable %s should be true or false", credential_form_schema.Variable)))
		}
		return strings.ToLower(value) == "true"
	}

	return value
}
