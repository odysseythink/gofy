package modelruntime

type AIModeler interface {
	ProviderName() string
	ModelType() ModelType
	ValidateCredentials(string, map[string]any)
	// InvokeErrorMapping() map[*modelruntimeexceptions.InvokeError][]error
	GetCustomizableModelSchema(string, map[string]any) *AIModelEntity
	PredefinedModels(AIModeler) []*AIModelEntity
	GetCustomizableModelSchemaFromCredentials(modeler AIModeler, model string, credentials map[string]any) *AIModelEntity
	GetModelSchema(modeler AIModeler, model string, credentials map[string]any) *AIModelEntity
}
