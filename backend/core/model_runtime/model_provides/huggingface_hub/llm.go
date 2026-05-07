package huggingface_hub

import (
	"fmt"
	"iter"
	"strings"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

const huggingfaceDefaultRouter = "https://router.huggingface.co/v1"

type HuggingfaceHubLLM struct {
	*base.LargeLanguageModel
}

func NewHuggingfaceHubLLM() *HuggingfaceHubLLM {
	return &HuggingfaceHubLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *HuggingfaceHubLLM) ProviderName() string { return "huggingface_hub" }
func (l *HuggingfaceHubLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *HuggingfaceHubLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *HuggingfaceHubLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *HuggingfaceHubLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func huggingfaceEndpoint(credentials map[string]any) string {
	baseURL, _ := credentials["huggingfacehub_endpoint_url"].(string)
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return huggingfaceDefaultRouter
	}
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	return baseURL
}

func (l *HuggingfaceHubLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["huggingfacehub_api_token"].(string)
	return oaicompat.Invoke(huggingfaceEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *HuggingfaceHubLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["huggingfacehub_api_token"].(string)
	return oaicompat.InvokeStream(huggingfaceEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
