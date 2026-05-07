package model

import (
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelenumtypes "github.com/odysseythink/gofy/backend/enum_types/model"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
)

// SimpleModelProviderEntity represents a simple provider entity
type SimpleModelProviderEntity struct {
	Provider            string                            `json:"provider" yaml:"provider"`
	Label               commontypes.I18nObject            `json:"label" yaml:"label"`
	IconSmall           *commontypes.I18nObject           `json:"icon_small,omitempty" yaml:"icon_small"`
	IconLarge           *commontypes.I18nObject           `json:"icon_large,omitempty" yaml:"icon_large"`
	SupportedModelTypes []modelruntimeenumtypes.ModelType `json:"supported_model_types" yaml:"supported_model_types"`
}

// NewSimpleModelProviderEntity creates a new instance of SimpleModelProviderEntity
func NewSimpleModelProviderEntity(providerEntity *modelruntimeentities.ProviderEntity) *SimpleModelProviderEntity {
	return &SimpleModelProviderEntity{
		Provider:            providerEntity.Provider,
		Label:               providerEntity.Label,
		IconSmall:           providerEntity.IconSmall,
		IconLarge:           providerEntity.IconLarge,
		SupportedModelTypes: providerEntity.SupportedModelTypes,
	}
}

// ProviderModelWithStatusEntity represents a model with status for model response
type ProviderModelWithStatusEntity struct {
	*modelruntimeentities.ProviderModel
	Status               modelenumtypes.ModelStatusType `json:"status" yaml:"status"`
	LoadBalancingEnabled bool                           `json:"load_balancing_enabled" yaml:"load_balancing_enabled"`
}

// ModelWithProviderEntity represents a model with provider entity
type ModelWithProviderEntity struct {
	*ProviderModelWithStatusEntity
	Provider *SimpleModelProviderEntity `json:"provider" yaml:"provider"`
}

// DefaultModelProviderEntity represents a default model provider entity
type DefaultModelProviderEntity struct {
	Provider            string                            `json:"provider" yaml:"provider"`
	Label               commontypes.I18nObject            `json:"label" yaml:"label"`
	IconSmall           *commontypes.I18nObject           `json:"icon_small,omitempty" yaml:"icon_small"`
	IconLarge           *commontypes.I18nObject           `json:"icon_large,omitempty" yaml:"icon_large"`
	SupportedModelTypes []modelruntimeenumtypes.ModelType `json:"supported_model_types" yaml:"supported_model_types"`
}

// DefaultModelEntity represents a default model entity
type DefaultModelEntity struct {
	Model     string                      `json:"model" yaml:"model"`
	ModelType string                      `json:"model_type" yaml:"model_type"`
	Provider  *DefaultModelProviderEntity `json:"provider" yaml:"provider"`
}
