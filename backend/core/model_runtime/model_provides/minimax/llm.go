package minimax

import (
	"fmt"
	"iter"

	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	"mlib.com/gofy/server/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

const minimaxEndpoint = "https://api.minimax.chat/v1"

type MiniMaxLLM struct {
	*base.LargeLanguageModel
}

func NewMiniMaxLLM() *MiniMaxLLM {
	return &MiniMaxLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *MiniMaxLLM) ProviderName() string                       { return "minimax" }
func (l *MiniMaxLLM) ModelType() modelruntimeenumtypes.ModelType { return modelruntimeenumtypes.Model_LLM }

func (l *MiniMaxLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *MiniMaxLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *MiniMaxLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func (l *MiniMaxLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.Invoke(minimaxEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *MiniMaxLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.InvokeStream(minimaxEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
