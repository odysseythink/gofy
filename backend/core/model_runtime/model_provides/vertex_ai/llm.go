package vertex_ai

import (
	"fmt"
	"iter"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
)

type VertexAILLM struct {
	*base.LargeLanguageModel
}

func NewVertexAILLM() *VertexAILLM {
	return &VertexAILLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *VertexAILLM) ProviderName() string { return "vertex_ai" }
func (l *VertexAILLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *VertexAILLM) ValidateCredentials(model string, credentials map[string]any) {
	project, _ := credentials["vertex_project_id"].(string)
	location, _ := credentials["vertex_location"].(string)
	key, _ := credentials["vertex_service_account_key"].(string)
	if project == "" || location == "" || key == "" {
		panic(fmt.Errorf("vertex_project_id, vertex_location and vertex_service_account_key are required"))
	}
}

func (l *VertexAILLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *VertexAILLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

func (l *VertexAILLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	panic(fmt.Errorf("vertex_ai invocation is not yet implemented; provider registered for configuration only (Google Cloud OAuth2 client required)"))
}

func (l *VertexAILLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	panic(fmt.Errorf("vertex_ai invocation is not yet implemented; provider registered for configuration only (Google Cloud OAuth2 client required)"))
}
