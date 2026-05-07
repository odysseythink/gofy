package openai_api_compatible

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
)

type OpenAICompatibleEmbedding struct {
	*base.TextEmbeddingModel
}

func NewOpenAICompatibleEmbedding() *OpenAICompatibleEmbedding {
	return &OpenAICompatibleEmbedding{
		TextEmbeddingModel: &base.TextEmbeddingModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_TEXT_EMBEDDING,
			},
		},
	}
}

func (e *OpenAICompatibleEmbedding) ProviderName() string {
	return "openai_api_compatible"
}

func (e *OpenAICompatibleEmbedding) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_TEXT_EMBEDDING
}

func (e *OpenAICompatibleEmbedding) ValidateCredentials(model string, credentials map[string]any) {
	_, err := e.InvokeEmbedding(model, credentials, []string{"ping"}, "")
	if err != nil {
		panic(fmt.Errorf("credentials validation failed: %w", err))
	}
}

func (e *OpenAICompatibleEmbedding) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return &modelruntimeentities.AIModelEntity{
		Model:     model,
		Label:     commontypes.I18nObject{EnUS: model, ZhHans: model},
		ModelType: modelruntimeenumtypes.Model_TEXT_EMBEDDING,
		FetchFrom: modelruntimeenumtypes.FetchFrom_CUSTOMIZABLE_MODEL,
		ModelProperties: map[modelruntimeenumtypes.ModelPropertyKey]any{
			modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE: 8192,
		},
	}
}

// InvokeEmbedding calls the OpenAI-compatible embedding API.
func (e *OpenAICompatibleEmbedding) InvokeEmbedding(model string, credentials map[string]any, texts []string, user string) (*modelruntimeentities.TextEmbeddingResult, error) {
	apiKey, _ := credentials["api_key"].(string)
	endpoint, _ := credentials["endpoint_url"].(string)
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	url := endpoint + "embeddings"

	reqBody := map[string]any{
		"model": model,
		"input": texts,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
		} `json:"data"`
		Usage struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	embeddings := make([][]float64, len(result.Data))
	for _, d := range result.Data {
		if d.Index < len(embeddings) {
			embeddings[d.Index] = d.Embedding
		}
	}

	return &modelruntimeentities.TextEmbeddingResult{
		Model:      model,
		Embeddings: embeddings,
		Usage: &modelruntimeentities.EmbeddingUsage{
			Tokens:      result.Usage.TotalTokens,
			TotalTokens: result.Usage.TotalTokens,
		},
	}, nil
}
