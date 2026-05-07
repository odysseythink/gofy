package bedrock

import (
	"fmt"
	"iter"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

type BedrockLLM struct {
	*base.LargeLanguageModel
}

func NewBedrockLLM() *BedrockLLM {
	return &BedrockLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *BedrockLLM) ProviderName() string { return "bedrock" }
func (l *BedrockLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *BedrockLLM) ValidateCredentials(model string, credentials map[string]any) {
	region, _ := credentials["aws_region"].(string)
	ak, _ := credentials["aws_access_key_id"].(string)
	sk, _ := credentials["aws_secret_access_key"].(string)
	if region == "" || ak == "" || sk == "" {
		panic(fmt.Errorf("aws_region, aws_access_key_id and aws_secret_access_key are required"))
	}
}

func (l *BedrockLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *BedrockLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func (l *BedrockLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	panic(fmt.Errorf("bedrock invocation is not yet implemented; provider registered for configuration only (AWS SigV4 client required)"))
}

func (l *BedrockLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	panic(fmt.Errorf("bedrock invocation is not yet implemented; provider registered for configuration only (AWS SigV4 client required)"))
}
