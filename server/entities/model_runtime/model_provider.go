package modelruntime

type ModelProvider interface {
	ProviderName() string
	ValidateProviderCredentials(credentials map[string]any)
	GetProviderSchema(string) *ProviderEntity
	Models(ModelProvider, ModelType) []*AIModelEntity
	GetModelInstance(ModelType) AIModeler
}
