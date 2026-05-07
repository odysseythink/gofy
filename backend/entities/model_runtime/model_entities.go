package modelruntime

import (
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
)

// ProviderModel represents a provider model.
type ProviderModel struct {
	Model           string                                         `json:"model" yaml:"model"`
	Label           commontypes.I18nObject                         `json:"label" yaml:"label"`
	ModelType       modelruntimeenumtypes.ModelType                `json:"model_type" yaml:"model_type"`
	Features        []modelruntimeenumtypes.ModelFeature           `json:"features" yaml:"features"`
	FetchFrom       modelruntimeenumtypes.FetchFrom                `json:"fetch_from" yaml:"fetch_from"`
	ModelProperties map[modelruntimeenumtypes.ModelPropertyKey]any `json:"model_properties" yaml:"model_properties"`
	Deprecated      bool                                           `json:"deprecated" yaml:"deprecated"`
	ModelConfig     map[string]any                                 `json:"model_config" yaml:"model_config"`
}

// ParameterRule represents a parameter rule.
type ParameterRule struct {
	Name        string                              `json:"name" yaml:"name"`
	UseTemplate string                              `json:"use_template" yaml:"use_template"`
	Label       commontypes.I18nObject              `json:"label" yaml:"label"`
	Type        modelruntimeenumtypes.ParameterType `json:"type" yaml:"type"`
	Help        commontypes.I18nObject              `json:"help" yaml:"help"`
	Required    bool                                `json:"required" yaml:"required"`
	Default     any                                 `json:"default" yaml:"default"`
	Min         float64                             `json:"min" yaml:"min"`
	Max         float64                             `json:"max" yaml:"max"`
	Precision   int                                 `json:"precision" yaml:"precision"`
	Options     []string                            `json:"options" yaml:"options"`
}

// PriceConfig represents the pricing configuration.
type PriceConfig struct {
	Input    float64 `json:"input" yaml:"input"`
	Output   float64 `json:"output" yaml:"output"`
	Unit     float64 `json:"unit" yaml:"parametuniter_rules"`
	Currency string  `json:"currency" yaml:"currency"`
}

// AIModelEntity represents an AI model entity.
type AIModelEntity struct {
	Model           string                                         `json:"model" yaml:"model"`
	Label           commontypes.I18nObject                         `json:"label" yaml:"label"`
	ModelType       modelruntimeenumtypes.ModelType                `json:"model_type" yaml:"model_type"`
	Features        []modelruntimeenumtypes.ModelFeature           `json:"features" yaml:"features"`
	FetchFrom       modelruntimeenumtypes.FetchFrom                `json:"fetch_from" yaml:"fetch_from"`
	ModelProperties map[modelruntimeenumtypes.ModelPropertyKey]any `json:"model_properties" yaml:"model_properties"`
	Deprecated      bool                                           `json:"deprecated" yaml:"deprecated"`
	ModelConfig     map[string]any                                 `json:"model_config" yaml:"model_config"`
	ParameterRules  []*ParameterRule                               `json:"parameter_rules" yaml:"parameter_rules"`
	Pricing         *PriceConfig                                   `json:"pricing" yaml:"pricing"`
}

// ModelUsage represents model usage information.
type ModelUsage struct {
}

type PriceType string

/*
Enum class for price type.
*/
const (
	PriceType_INPUT  PriceType = "input"
	PriceType_OUTPUT PriceType = "output"
)

// PriceInfo represents price information.
type PriceInfo struct {
	UnitPrice   float64
	Unit        float64
	TotalAmount float64
	Currency    string
}
