package base

import (
	"fmt"

	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type TextEmbeddingModel struct {
	*BaseAIModel
}

func (m *TextEmbeddingModel) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_TEXT_EMBEDDING
}

// InvokeSync calls the embedding model synchronously (to be overridden by providers).
func (m *TextEmbeddingModel) InvokeSync(model string, credentials map[string]any, texts []string, user string) (*modelruntimeentities.TextEmbeddingResult, error) {
	return nil, fmt.Errorf("InvokeSync not implemented for base TextEmbeddingModel")
}

// GetNumTokens estimates token count for texts.
func (m *TextEmbeddingModel) GetNumTokens(model string, credentials map[string]any, texts []string) int {
	total := 0
	for _, text := range texts {
		// Simple approximation: 1 token ≈ 4 chars for English, 1 token ≈ 2 chars for CJK
		total += len([]rune(text)) / 2
	}
	return total
}
