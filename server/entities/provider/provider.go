package provider

import (
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	parameterenumtypes "mlib.com/gofy/server/enum_types/parameter"
	providerenumtypes "mlib.com/gofy/server/enum_types/provider"
	commontypes "mlib.com/gofy/server/types/common"
)

type RestrictModel struct {
	Model         string                          `json:"model"`
	BaseModelName string                          `json:"base_model_name"`
	ModelType     modelruntimeenumtypes.ModelType `json:"model_type"`

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}
type QuotaConfiguration struct {
	QuotaType      models.ProviderQuotaType `json:"quota_type"`
	QuotaUnit      coreenumtypes.QuotaUnit  `json:"quota_unit"`
	QuotaLimit     int                      `json:"quota_limit"`
	QuotaUsed      int                      `json:"quota_used"`
	IsValid        bool                     `json:"is_valid"`
	RestrictModels []*RestrictModel         `json:"restrict_models"`
}
type SystemConfiguration struct {
	Enabled             bool                     `json:"enabled"`
	CurrentQuotaType    models.ProviderQuotaType `json:"current_quota_type"`
	QuotaConfigurations []*QuotaConfiguration    `json:"quota_configurations"`
	Credentials         map[string]any           `json:"credentials"`
}
type CustomProviderConfiguration struct {
	Credentials map[string]any `json:"credentials"`
}
type CustomModelConfiguration struct {
	/*
	   Model class for provider custom model configuration.
	*/

	Model                           string `json:"model"`
	modelruntimeenumtypes.ModelType `json:"model_type"`
	Credentials                     map[string]any `json:"credentials"`

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}
type CustomConfiguration struct {
	/*
	   Model class for provider custom configuration.
	*/

	Provider *CustomProviderConfiguration `json:"provider"`
	Models   []*CustomModelConfiguration  `json:"models"`
}
type ModelLoadBalancingConfiguration struct {
	/*
	   Class for model load balancing configuration.
	*/

	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Credentials map[string]any `json:"credentials"`
}
type ModelSetting struct {
	/*
	   Model class for model settings.
	*/

	Model                           string `json:"model"`
	modelruntimeenumtypes.ModelType `json:"model_type"`
	Enabled                         bool                               `json:"enabled"`
	LoadBalancingConfigs            []*ModelLoadBalancingConfiguration `json:"load_balancing_configs"`

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}

func NewModelSetting() *ModelSetting {
	return &ModelSetting{
		Enabled: true,
	}
}

type BasicProviderConfig struct {
	Type providerenumtypes.BasicProviderConfigType `json:"type"` //description="The type of the credentials")
	Name string                                    `json:"name"` //description="The name of the credentials")
}
type ProviderConfig[T1 parameterenumtypes.AppSelectorScopeType | parameterenumtypes.ModelSelectorScopeType | parameterenumtypes.ToolSelectorScopeType, T2 int | string] struct {
	*BasicProviderConfig

	Scope    T1   `json:"scope"`
	Required bool `json:"required"`
	Default  T2   `json:"default"`
	Options  []struct {
		Value string                 `json:"value"` //description="The value of the option")
		Label commontypes.I18nObject `json:"label"` //description="The label of the option")
	} `json:"options"`
	Label       *commontypes.I18nObject `json:"label"`
	Help        *commontypes.I18nObject `json:"help"`
	URL         string                  `json:"url"`
	Placeholder *commontypes.I18nObject `json:"placeholder"`
}

func (pc *ProviderConfig[T1, T2]) ToBasicProviderConfig() *BasicProviderConfig {
	return &BasicProviderConfig{Type: pc.Type, Name: pc.Name}
}
