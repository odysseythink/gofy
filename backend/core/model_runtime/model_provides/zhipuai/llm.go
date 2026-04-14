package zhipuai

import (
	"fmt"
	"iter"

	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	"mlib.com/gofy/server/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

const zhipuaiEndpoint = "https://open.bigmodel.cn/api/paas/v4"

type ZhipuAILLM struct {
	*base.LargeLanguageModel
}

func NewZhipuAILLM() *ZhipuAILLM {
	return &ZhipuAILLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *ZhipuAILLM) ProviderName() string                       { return "zhipuai" }
func (l *ZhipuAILLM) ModelType() modelruntimeenumtypes.ModelType { return modelruntimeenumtypes.Model_LLM }

func (l *ZhipuAILLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *ZhipuAILLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *ZhipuAILLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func (l *ZhipuAILLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.Invoke(zhipuaiEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *ZhipuAILLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.InvokeStream(zhipuaiEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
