package huggingface_hub

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type HuggingfaceHubProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&HuggingfaceHubProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *HuggingfaceHubProvider) ProviderName() string {
	return "huggingface_hub"
}

func (p *HuggingfaceHubProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Hugging Face Hub uses model-level credentials only
}

func (p *HuggingfaceHubProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewHuggingfaceHubLLM()
	default:
		return nil
	}
}
