package upstage

import (
	"fmt"
	"iter"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

const upstageEndpoint = "https://api.upstage.ai/v1/solar"

type UpstageLLM struct {
	*base.LargeLanguageModel
}

func NewUpstageLLM() *UpstageLLM {
	return &UpstageLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *UpstageLLM) ProviderName() string { return "upstage" }
func (l *UpstageLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *UpstageLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *UpstageLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *UpstageLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func (l *UpstageLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.Invoke(upstageEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *UpstageLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.InvokeStream(upstageEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
