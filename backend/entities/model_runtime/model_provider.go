package modelruntime

import (
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

type ModelProvider interface {
	ProviderName() string
	ValidateProviderCredentials(credentials map[string]any)
	GetProviderSchema(string) *ProviderEntity
	Models(ModelProvider, modelruntimeenumtypes.ModelType) []*AIModelEntity
	GetModelInstance(modelruntimeenumtypes.ModelType) AIModeler
}
