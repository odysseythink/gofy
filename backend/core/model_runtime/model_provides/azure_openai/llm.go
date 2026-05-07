package azure_openai

import (
	"fmt"
	"iter"
	"strings"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

type AzureOpenAILLM struct {
	*base.LargeLanguageModel
}

func NewAzureOpenAILLM() *AzureOpenAILLM {
	return &AzureOpenAILLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *AzureOpenAILLM) ProviderName() string { return "azure_openai" }
func (l *AzureOpenAILLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *AzureOpenAILLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *AzureOpenAILLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *AzureOpenAILLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func azureEndpoint(credentials map[string]any) string {
	baseURL, _ := credentials["openai_api_base"].(string)
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		panic(fmt.Errorf("openai_api_base is required"))
	}
	// Newer Azure OpenAI supports OpenAI-compatible routes under /openai/v1
	if !strings.Contains(baseURL, "/openai/v1") {
		baseURL += "/openai/v1"
	}
	return baseURL
}

func (l *AzureOpenAILLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["openai_api_key"].(string)
	return oaicompat.Invoke(azureEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *AzureOpenAILLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["openai_api_key"].(string)
	return oaicompat.InvokeStream(azureEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
