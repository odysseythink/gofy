package cohere

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
)

func NewCohereEmbedding() *oaicompat.SimpleEmbedding {
	return oaicompat.NewSimpleEmbedding(
		"cohere",
		512,
		cohereEmbeddingURL,
		func(credentials map[string]any) string {
			apiKey, _ := credentials["api_key"].(string)
			return apiKey
		},
	)
}

func cohereEmbeddingURL(credentials map[string]any) string {
	return cohereEndpoint(credentials)
}
