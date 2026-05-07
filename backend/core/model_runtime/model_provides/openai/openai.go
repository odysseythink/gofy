package openai

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type OpenAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&OpenAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *OpenAIProvider) ProviderName() string {
	return "openai"
}

func (p *OpenAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("gpt-4o-mini", credentials)
}

func (p *OpenAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewOpenAILLM()
	default:
		return nil
	}
}
