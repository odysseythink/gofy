package services

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	modelentities "mlib.com/gofy/server/entities/model"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	providerentities "mlib.com/gofy/server/entities/provider"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/mlog"
)

// CustomConfigurationStatus represents the status of custom configuration.
type CustomConfigurationStatus string

const (
	CustomConfigurationStatus_ACTIVE       CustomConfigurationStatus = "active"
	CustomConfigurationStatus_NO_CONFIGURE CustomConfigurationStatus = "no-configure"
)

// CustomConfigurationResponse represents the response for custom configuration.
type CustomConfigurationResponse struct {
	Status CustomConfigurationStatus `json:"status" yaml:"status"`
}

// SystemConfigurationResponse represents the response for system configuration.
type SystemConfigurationResponse struct {
	Enabled             bool                                   `json:"enabled" yaml:"enabled"`
	CurrentQuotaType    providerenumtypes.ProviderQuotaType    `json:"current_quota_type" yaml:"current_quota_type"` // Assuming ProviderQuotaType is a string
	QuotaConfigurations []*providerentities.QuotaConfiguration `json:"quota_configurations" yaml:"quota_configurations"`
}

// ProviderResponse represents the response for a provider.
type ProviderResponse struct {
	Provider                 string                                         `json:"provider" yaml:"provider"`
	Label                    commontypes.I18nObject                         `json:"label" yaml:"label"`
	Description              *commontypes.I18nObject                        `json:"description" yaml:"description"`
	IconSmall                *commontypes.I18nObject                        `json:"icon_small" yaml:"icon_small"`
	IconLarge                *commontypes.I18nObject                        `json:"icon_large" yaml:"icon_large"`
	Background               string                                         `json:"background" yaml:"background"`
	Help                     *modelruntimeentities.ProviderHelpEntity       `json:"help" yaml:"help"`
	SupportedModelTypes      []modelruntimeenumtypes.ModelType              `json:"supported_model_types" yaml:"supported_model_types"`           // Assuming ModelType is a string
	ConfigurateMethods       []modelruntimeentities.ConfigurateMethod       `json:"configurate_methods" yaml:"configurate_methods"`               // Assuming ConfigurateMethod is a string
	ProviderCredentialSchema *modelruntimeentities.ProviderCredentialSchema `json:"provider_credential_schema" yaml:"provider_credential_schema"` // Assuming ProviderCredentialSchema is a complex type
	ModelCredentialSchema    *modelruntimeentities.ModelCredentialSchema    `json:"model_credential_schema" yaml:"model_credential_schema"`       // Assuming ModelCredentialSchema is a complex type
	PreferredProviderType    providerenumtypes.ProviderType                 `json:"preferred_provider_type" yaml:"preferred_provider_type"`       // Assuming ProviderType is a string
	CustomConfiguration      *CustomConfigurationResponse                   `json:"custom_configuration" yaml:"custom_configuration"`
	SystemConfiguration      *SystemConfigurationResponse                   `json:"system_configuration" yaml:"system_configuration"`
}

func NewProviderResponse(args map[string]any) *ProviderResponse {

	rsp := &ProviderResponse{}
	if args != nil {
		bindata, _ := json.Marshal(args)
		err := json.Unmarshal(bindata, rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%#v) failed:%v", args, err)
		}
	}
	url_prefix := viper.GetString("CONSOLE_API_URL") + "/console/api/workspaces/current/model-providers/" + rsp.Provider
	if rsp.IconSmall != nil {
		rsp.IconSmall = &commontypes.I18nObject{
			EnUS: fmt.Sprintf("%s/icon_small/en_US", url_prefix), ZhHans: fmt.Sprintf("%s/icon_small/zh_Hans", url_prefix),
		}
	}
	if rsp.IconLarge != nil {
		rsp.IconLarge = &commontypes.I18nObject{
			EnUS: fmt.Sprintf("%s/icon_large/en_US", url_prefix), ZhHans: fmt.Sprintf("%s/icon_large/zh_Hans", url_prefix),
		}
	}
	return rsp
}

// ProviderWithModelsResponse represents the response for a provider with models.
type ProviderWithModelsResponse struct {
	Provider  string                                         `json:"provider" yaml:"provider"`
	Label     commontypes.I18nObject                         `json:"label" yaml:"label"`
	IconSmall *commontypes.I18nObject                        `json:"icon_small" yaml:"icon_small"`
	IconLarge *commontypes.I18nObject                        `json:"icon_large" yaml:"icon_large"`
	Status    CustomConfigurationStatus                      `json:"status" yaml:"status"`
	Models    []*modelentities.ProviderModelWithStatusEntity `json:"models" yaml:"models"`
}

func NewProviderWithModelsResponse(args map[string]any) *ProviderWithModelsResponse {

	rsp := &ProviderWithModelsResponse{}
	if args != nil {
		bindata, _ := json.Marshal(args)
		err := json.Unmarshal(bindata, rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%#v) failed:%v", args, err)
		}
	}
	url_prefix := viper.GetString("CONSOLE_API_URL") + "/console/api/workspaces/current/model-providers/" + rsp.Provider
	if rsp.IconSmall != nil {
		rsp.IconSmall = &commontypes.I18nObject{
			EnUS: fmt.Sprintf("%s/icon_small/en_US", url_prefix), ZhHans: fmt.Sprintf("%s/icon_small/zh_Hans", url_prefix),
		}
	}
	if rsp.IconLarge != nil {
		rsp.IconLarge = &commontypes.I18nObject{
			EnUS: fmt.Sprintf("%s/icon_large/en_US", url_prefix), ZhHans: fmt.Sprintf("%s/icon_large/zh_Hans", url_prefix),
		}
	}
	return rsp
}

// SimpleProviderEntityResponse represents a simple provider entity response.
type SimpleProviderEntityResponse struct {
	*modelruntimeentities.SimpleProviderEntity
}

func NewSimpleProviderEntityResponse(args map[string]any) *SimpleProviderEntityResponse {

	rsp := &SimpleProviderEntityResponse{}
	if args != nil {
		bindata, _ := json.Marshal(args)
		err := json.Unmarshal(bindata, rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%#v) failed:%v", args, err)
		}
	}
	url_prefix := viper.GetString("CONSOLE_API_URL") + "/console/api/workspaces/current/model-providers/" + rsp.Provider
	if rsp.IconSmall != nil {
		rsp.IconSmall = &commontypes.I18nObject{
			EnUS: fmt.Sprintf("%s/icon_small/en_US", url_prefix), ZhHans: fmt.Sprintf("%s/icon_small/zh_Hans", url_prefix),
		}
	}
	if rsp.IconLarge != nil {
		rsp.IconLarge = &commontypes.I18nObject{
			EnUS: fmt.Sprintf("%s/icon_large/en_US", url_prefix), ZhHans: fmt.Sprintf("%s/icon_large/zh_Hans", url_prefix),
		}
	}
	return rsp
}

// DefaultModelResponse represents a default model response.
type DefaultModelResponse struct {
	Model     string                        `json:"model"`
	ModelType string                        `json:"model_type"` // Assuming ModelType is a string
	Provider  *SimpleProviderEntityResponse `json:"provider"`
}

func NewDefaultModelResponse(args any) *DefaultModelResponse {
	if real_args, ok := args.(string); ok && real_args != "" {
		rsp := new(DefaultModelResponse)
		err := json.Unmarshal([]byte(real_args), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to DefaultModelResponse failed:%v", real_args, err)
		}
		return rsp
	} else if real_args, ok := args.(map[string]any); ok && len(real_args) > 0 {
		rsp := new(DefaultModelResponse)
		bindata, _ := json.Marshal(real_args)
		err := json.Unmarshal([]byte(bindata), rsp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to DefaultModelResponse failed:%v", string(bindata), err)
		}
		return rsp
	} else {
		mlog.Errorf("unsupported args=%#v", args)
	}
	return new(DefaultModelResponse)
}

// ModelWithProviderEntityResponse represents a model with provider entity response.
type ModelWithProviderEntityResponse struct {
	// Assuming ModelWithProviderEntity is a complex type
	*modelentities.ModelWithProviderEntity
	Provider *SimpleProviderEntityResponse `json:"provider"`
}

func NewModelWithProviderEntityResponse(model *modelentities.ModelWithProviderEntity) *ModelWithProviderEntityResponse {

	rsp := &ModelWithProviderEntityResponse{
		ModelWithProviderEntity: model,
	}
	return rsp
}

// type SimpleProviderEntityResponse struct {
// 	*SimpleProviderEntity
// 	UrlPrefix string `json:"url_prefix"`
// }

// /*
//    Simple provider entity response.
// */

// func NewSimpleProviderEntityResponse(console_api_url string, simple_provider *SimpleProviderEntity) *SimpleProviderEntityResponse {
// 	//     super().__init__(**data)
// 	rsp := &SimpleProviderEntityResponse{
// 		SimpleProviderEntity: simple_provider,
// 	}
// 	// url_prefix = dify_config.CONSOLE_API_URL + f"/console/api/workspaces/current/model-providers/{self.provider}"
// 	rsp.UrlPrefix = console_api_url + "/console/api/workspaces/current/model-providers/" + rsp.Provider
// 	if rsp.IconSmall == nil {
// 		rsp.IconSmall = &commontypes.I18nObject{
// 			EnUS:   rsp.UrlPrefix + "/icon_small/en_US",
// 			ZhHans: rsp.UrlPrefix + "/icon_small/zh_Hans",
// 		}
// 	}
// 	if rsp.IconLarge == nil {
// 		rsp.IconLarge = &commontypes.I18nObject{
// 			EnUS:   rsp.UrlPrefix + "/icon_large/en_US",
// 			ZhHans: rsp.UrlPrefix + "/icon_large/zh_Hans",
// 		}
// 	}
// 	return rsp
// }
