package azure_openai

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type AzureOpenAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&AzureOpenAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *AzureOpenAIProvider) ProviderName() string {
	return "azure_openai"
}

func (p *AzureOpenAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Azure OpenAI uses model-level credentials only
}

func (p *AzureOpenAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewAzureOpenAILLM()
	default:
		return nil
	}
}
