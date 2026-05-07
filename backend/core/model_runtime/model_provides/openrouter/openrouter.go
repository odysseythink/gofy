package openrouter

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type OpenRouterProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&OpenRouterProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *OpenRouterProvider) ProviderName() string {
	return "openrouter"
}

func (p *OpenRouterProvider) ValidateProviderCredentials(credentials map[string]any) {
	// OpenRouter uses model-level credentials only
}

func (p *OpenRouterProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewOpenRouterLLM()
	default:
		return nil
	}
}
