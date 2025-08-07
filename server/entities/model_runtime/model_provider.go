package modelruntime

import (
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type ModelProvider interface {
	ProviderName() string
	ValidateProviderCredentials(credentials map[string]any)
	GetProviderSchema(string) *ProviderEntity
	Models(ModelProvider, modelruntimeenumtypes.ModelType) []*AIModelEntity
	GetModelInstance(modelruntimeenumtypes.ModelType) AIModeler
}
