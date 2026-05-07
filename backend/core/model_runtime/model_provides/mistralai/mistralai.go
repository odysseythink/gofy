package mistralai

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type MistralAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&MistralAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *MistralAIProvider) ProviderName() string {
	return "mistralai"
}

func (p *MistralAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Mistral uses model-level credentials only
}

func (p *MistralAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewMistralAILLM()
	default:
		return nil
	}
}
