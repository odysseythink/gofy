package oaicompat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
)

// DefaultEmbeddingModelSchema returns a default AIModelEntity for customizable embedding models.
func DefaultEmbeddingModelSchema(model string, contextSize int) *modelruntimeentities.AIModelEntity {
	if contextSize <= 0 {
		contextSize = 8192
	}
	return &modelruntimeentities.AIModelEntity{
		Model:     model,
		Label:     commontypes.I18nObject{EnUS: model, ZhHans: model},
		ModelType: modelruntimeenumtypes.Model_TEXT_EMBEDDING,
		FetchFrom: modelruntimeenumtypes.FetchFrom_CUSTOMIZABLE_MODEL,
		ModelProperties: map[modelruntimeenumtypes.ModelPropertyKey]any{
			modelruntimeenumtypes.ModelPropertyKey_CONTEXT_SIZE: contextSize,
		},
	}
}

// InvokeEmbedding posts to an OpenAI-compatible embeddings endpoint and returns embeddings.
func InvokeEmbedding(endpoint, apiKey, model string, texts []string, user string) (*modelruntimeentities.TextEmbeddingResult, error) {
	endpoint = strings.TrimRight(endpoint, "/")
	url := endpoint + "/embeddings"

	reqBody := map[string]any{
		"model": model,
		"input": texts,
	}
	if user != "" {
		reqBody["user"] = user
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("embedding API error %d: %s", resp.StatusCode, string(body))
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
