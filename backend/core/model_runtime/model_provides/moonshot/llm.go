package moonshot

import (
	"fmt"
	"iter"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

const moonshotEndpoint = "https://api.moonshot.cn/v1"

type MoonshotLLM struct {
	*base.LargeLanguageModel
}

func NewMoonshotLLM() *MoonshotLLM {
	return &MoonshotLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *MoonshotLLM) ProviderName() string { return "moonshot" }
func (l *MoonshotLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *MoonshotLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *MoonshotLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *MoonshotLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func (l *MoonshotLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.Invoke(moonshotEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}

func (l *MoonshotLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	apiKey, _ := credentials["api_key"].(string)
	return oaicompat.InvokeStream(moonshotEndpoint, apiKey, model, promptMessages, modelParameters, tools, stop, user)
}
