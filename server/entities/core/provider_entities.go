package core

import (
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	coreenumtypes "mlib.com/gofy/server/enum_types/core"
	"mlib.com/gofy/server/models"
)

type RestrictModel struct {
	Model                          string `json:"model"`
	BaseModelName                  string `json:"base_model_name"`
	modelruntimeentities.ModelType `json:"model_type"`

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}

type QuotaConfiguration struct {
	/*
	   Model class for provider quota configuration.
	*/

	QuotaType      models.ProviderQuotaType `json:"quota_type"`
	QuotaUnit      coreenumtypes.QuotaUnit  `json:"quota_unit"`
	QuotaLimit     int                      `json:"quota_limit"`
	QuotaUsed      int                      `json:"quota_used"`
	IsValid        bool                     `json:"is_valid"`
	RestrictModels []*RestrictModel         `json:"restrict_models"`
}

type SystemConfiguration struct {
	/*
	   Model class for provider system configuration.
	*/

	Enabled             bool                     `json:"enabled"`
	CurrentQuotaType    models.ProviderQuotaType `json:"current_quota_type"`
	QuotaConfigurations []*QuotaConfiguration    `json:"quota_configurations"`
	Credentials         map[string]any           `json:"credentials"`
}

type CustomProviderConfiguration struct {
	/*
	   Model class for provider custom configuration.
	*/

	Credentials map[string]any `json:"credentials"`
}

type CustomModelConfiguration struct {
	/*
	   Model class for provider custom model configuration.
	*/

	Model                          string `json:"model"`
	modelruntimeentities.ModelType `json:"model_type"`
	Credentials                    map[string]any `json:"credentials"`

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

	Model                          string `json:"model"`
	modelruntimeentities.ModelType `json:"model_type"`
	Enabled                        bool                               `json:"enabled"`
	LoadBalancingConfigs           []*ModelLoadBalancingConfiguration `json:"load_balancing_configs"`

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}

func NewModelSetting() *ModelSetting {
	return &ModelSetting{
		Enabled: true,
	}
}
