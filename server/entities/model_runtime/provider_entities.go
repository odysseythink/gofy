package modelruntime

import (
	"fmt"

	"mlib.com/confy"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/gofy/server/utils"
)

type ConfigurateMethod string

const (
	/*
	   Enum class for configurate method of provider model.
	*/

	ConfigurateMethod_PREDEFINED_MODEL   ConfigurateMethod = "predefined-model"
	ConfigurateMethod_CUSTOMIZABLE_MODEL ConfigurateMethod = "customizable-model"
)

type FormType string

const (
	/*
	   Enum class for form type.
	*/

	Form_TEXT_INPUT   FormType = "text-input"
	Form_SECRET_INPUT FormType = "secret-input"
	Form_SELECT       FormType = "select"
	Form_RADIO        FormType = "radio"
	Form_SWITCH       FormType = "switch"
)

// FormShowOnObject represents conditions for showing a form field.
type FormShowOnObject struct {
	Variable string `json:"variable" yaml:"variable"`
	Value    string `json:"value" yaml:"value"`
}

// FormOption represents an option in a form field.
type FormOption struct {
	Label  commontypes.I18nObject `json:"label" yaml:"label"`
	Value  string                 `json:"value" yaml:"value"`
	ShowOn []*FormShowOnObject    `json:"show_on" yaml:"show_on"`
}

// CredentialFormSchema represents the schema for a credential form.
type CredentialFormSchema struct {
	Variable    string                 `json:"variable" yaml:"variable"`
	Label       commontypes.I18nObject `json:"label" yaml:"label"`
	Type        FormType               `json:"type" yaml:"type"`
	Required    bool                   `json:"required" yaml:"required"`
	Default     string                 `json:"default" yaml:"default"`
	Options     []*FormOption          `json:"options" yaml:"options"`
	Placeholder commontypes.I18nObject `json:"placeholder" yaml:"placeholder"`
	MaxLength   int                    `json:"max_length" yaml:"max_length"`
	ShowOn      []*FormShowOnObject    `json:"show_on" yaml:"show_on"`
}

func NewCredentialFormSchema() *CredentialFormSchema {
	return &CredentialFormSchema{
		Required: true,
	}
}

// ProviderCredentialSchema represents the schema for provider credentials.
type ProviderCredentialSchema struct {
	CredentialFormSchemas []*CredentialFormSchema `json:"credential_form_schemas" yaml:"credential_form_schemas"`
}

// FieldModelSchema represents the schema for a model field.
type FieldModelSchema struct {
	Label       commontypes.I18nObject `json:"label" yaml:"label"`
	Placeholder commontypes.I18nObject `json:"placeholder" yaml:"placeholder"`
}

// ModelCredentialSchema represents the schema for model credentials.
type ModelCredentialSchema struct {
	Model                 *FieldModelSchema       `json:"model" yaml:"model"`
	CredentialFormSchemas []*CredentialFormSchema `json:"credential_form_schemas" yaml:"credential_form_schemas"`
}

// SimpleProviderEntity represents a simplified provider entity.
type SimpleProviderEntity struct {
	Provider            string                            `json:"provider" yaml:"provider"`
	Label               commontypes.I18nObject            `json:"label" yaml:"label"`
	IconSmall           *commontypes.I18nObject           `json:"icon_small" yaml:"icon_small"`
	IconLarge           *commontypes.I18nObject           `json:"icon_large" yaml:"icon_large"`
	SupportedModelTypes []modelruntimeenumtypes.ModelType `json:"supported_model_types" yaml:"supported_model_types"`
	Models              []*AIModelEntity/*ProviderModel*/ `json:"models" yaml:"models"`
}

// ProviderHelpEntity represents help information for a provider.
type ProviderHelpEntity struct {
	Title commontypes.I18nObject `json:"title" yaml:"title"`
	URL   commontypes.I18nObject `json:"url" yaml:"url"`
}

// ProviderEntity represents a provider entity.
type ProviderEntity struct {
	Provider                 string                            `json:"provider" yaml:"provider"`
	Label                    commontypes.I18nObject            `json:"label" yaml:"label"`
	Description              *commontypes.I18nObject           `json:"description" yaml:"description"`
	IconSmall                *commontypes.I18nObject           `json:"icon_small" yaml:"icon_small"`
	IconLarge                *commontypes.I18nObject           `json:"icon_large" yaml:"icon_large"`
	Background               string                            `json:"background" yaml:"background"`
	Help                     *ProviderHelpEntity               `json:"help" yaml:"help"`
	SupportedModelTypes      []modelruntimeenumtypes.ModelType `json:"supported_model_types" yaml:"supported_model_types"`
	ConfigurateMethods       []ConfigurateMethod               `json:"configurate_methods" yaml:"configurate_methods"`
	Models                   []*AIModelEntity/*ProviderModel*/ `json:"models" yaml:"models"`
	ProviderCredentialSchema *ProviderCredentialSchema `json:"provider_credential_schema" yaml:"provider_credential_schema"`
	ModelCredentialSchema    *ModelCredentialSchema    `json:"model_credential_schema" yaml:"model_credential_schema"`
}

func (pe *ProviderEntity) Fullfile() {
	if pe.ModelCredentialSchema != nil {
		for idx, v := range pe.ModelCredentialSchema.CredentialFormSchemas {
			if v.Options == nil {
				v.Options = make([]*FormOption, 0)
			}
			for sidx, v1 := range v.Options {
				if v1.ShowOn == nil {
					v1.ShowOn = make([]*FormShowOnObject, 0)
				}
				v.Options[sidx] = v1
			}
			if v.ShowOn == nil {
				v.ShowOn = make([]*FormShowOnObject, 0)
			}
			pe.ModelCredentialSchema.CredentialFormSchemas[idx] = v
		}
	}
	ip := utils.GetIP()
	port := confy.Get[int]("system.addr")
	url := ""
	if port > 0 && port < 65536 {
		url = fmt.Sprintf("http://%s:%d", ip, port)
	}
	if url != "" {
		if pe.IconSmall != nil {
			if pe.IconSmall.EnUS != "" {
				pe.IconSmall.EnUS = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconSmall.EnUS
			}
			if pe.IconSmall.ZhHans != "" {
				pe.IconSmall.ZhHans = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconSmall.ZhHans
			}
			if pe.IconSmall.PtBR != "" {
				pe.IconSmall.PtBR = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconSmall.PtBR
			}
			if pe.IconSmall.JaJP != "" {
				pe.IconSmall.JaJP = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconSmall.JaJP
			}
		}
		if pe.IconLarge != nil {
			if pe.IconLarge.EnUS != "" {
				pe.IconLarge.EnUS = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconLarge.EnUS
			}
			if pe.IconLarge.ZhHans != "" {
				pe.IconLarge.ZhHans = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconLarge.ZhHans
			}
			if pe.IconLarge.PtBR != "" {
				pe.IconLarge.PtBR = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconLarge.PtBR
			}
			if pe.IconLarge.JaJP != "" {
				pe.IconLarge.JaJP = url + "/model_runtime/model_provides/" + pe.Provider + "/assets/" + pe.IconLarge.JaJP
			}
		}
	}
	if pe.ProviderCredentialSchema != nil {
		for idx, v := range pe.ProviderCredentialSchema.CredentialFormSchemas {
			if v.Options == nil {
				v.Options = make([]*FormOption, 0)
			}
			for sidx, v1 := range v.Options {
				if v1.ShowOn == nil {
					v1.ShowOn = make([]*FormShowOnObject, 0)
				}
				v.Options[sidx] = v1
			}
			if v.ShowOn == nil {
				v.ShowOn = make([]*FormShowOnObject, 0)
			}
			pe.ProviderCredentialSchema.CredentialFormSchemas[idx] = v
		}
	}
}

// ToSimpleProvider converts the ProviderEntity to a SimpleProviderEntity.
func (p *ProviderEntity) ToSimpleProvider() SimpleProviderEntity {
	return SimpleProviderEntity{
		Provider:            p.Provider,
		Label:               p.Label,
		IconSmall:           p.IconSmall,
		IconLarge:           p.IconLarge,
		SupportedModelTypes: p.SupportedModelTypes,
		Models:              p.Models,
	}
}

// ProviderConfig represents a provider configuration.
type ProviderConfig struct {
	Provider    string         `json:"provider" yaml:"provider"`
	Credentials map[string]any `json:"credentials" yaml:"credentials"`
}
