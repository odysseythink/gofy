package cohere

import (
	"fmt"
	"iter"
	"strings"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

const cohereDefaultEndpoint = "https://api.cohere.com/compatibility/v1"

type CohereLLM struct {
	*base.LargeLanguageModel
}

func NewCohereLLM() *CohereLLM {
	return &CohereLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *CohereLLM) ProviderName() string { return "cohere" }
func (l *CohereLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *CohereLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *CohereLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *CohereLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func cohereEndpoint(credentials map[string]any) string {
	baseURL, _ := credentials["base_url"].(string)
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return cohereDefaultEndpoint
	}
	return baseURL
}

func (l *CohereLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.Invoke(cohereEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *CohereLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.InvokeStream(cohereEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
