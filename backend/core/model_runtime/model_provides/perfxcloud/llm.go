package perfxcloud

import (
	"fmt"
	"iter"
	"strings"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

const perfxcloudDefaultEndpoint = "https://cloud.perfxlab.cn/v1"

type PerfXCloudLLM struct {
	*base.LargeLanguageModel
}

func NewPerfXCloudLLM() *PerfXCloudLLM {
	return &PerfXCloudLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *PerfXCloudLLM) ProviderName() string { return "perfxcloud" }
func (l *PerfXCloudLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *PerfXCloudLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *PerfXCloudLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *PerfXCloudLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func perfxcloudEndpoint(credentials map[string]any) string {
	baseURL, _ := credentials["endpoint_url"].(string)
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return perfxcloudDefaultEndpoint
	}
	return baseURL
}

func (l *PerfXCloudLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.Invoke(perfxcloudEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *PerfXCloudLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.InvokeStream(perfxcloudEndpoint(credentials), apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
