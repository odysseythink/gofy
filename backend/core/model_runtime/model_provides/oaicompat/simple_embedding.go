package oaicompat

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

// SimpleEmbedding is a reusable OpenAI-compatible text-embedding model instance.
// Providers supply their own provider name and an endpoint resolver that reads credentials.
type SimpleEmbedding struct {
	*base.TextEmbeddingModel
	Provider       string
	ContextSize    int
	ResolveURL     func(credentials map[string]any) string
	ResolveAPIKey  func(credentials map[string]any) string
}

func NewSimpleEmbedding(provider string, contextSize int, resolveURL func(map[string]any) string, resolveAPIKey func(map[string]any) string) *SimpleEmbedding {
	return &SimpleEmbedding{
		TextEmbeddingModel: &base.TextEmbeddingModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_TEXT_EMBEDDING,
			},
		},
		Provider:      provider,
		ContextSize:   contextSize,
		ResolveURL:    resolveURL,
		ResolveAPIKey: resolveAPIKey,
	}
}

func (e *SimpleEmbedding) ProviderName() string { return e.Provider }
func (e *SimpleEmbedding) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_TEXT_EMBEDDING
}

func (e *SimpleEmbedding) ValidateCredentials(model string, credentials map[string]any) {
	_, err := e.InvokeEmbedding(model, credentials, []string{"ping"}, "")
	if err != nil {
		panic(fmt.Errorf("credentials validation failed: %w", err))
	}
}

func (e *SimpleEmbedding) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return DefaultEmbeddingModelSchema(model, e.ContextSize)
}

func (e *SimpleEmbedding) InvokeEmbedding(model string, credentials map[string]any, texts []string, user string) (*modelruntimeentities.TextEmbeddingResult, error) {
	url := ""
	if e.ResolveURL != nil {
		url = e.ResolveURL(credentials)
	}
	apiKey := ""
	if e.ResolveAPIKey != nil {
		apiKey = e.ResolveAPIKey(credentials)
	}
	return InvokeEmbedding(url, apiKey, model, texts, user)
}
