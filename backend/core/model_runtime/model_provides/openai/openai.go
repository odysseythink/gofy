package openai

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
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
