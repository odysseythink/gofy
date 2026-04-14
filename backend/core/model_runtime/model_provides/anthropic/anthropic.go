package anthropic

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
)

type AnthropicProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&AnthropicProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *AnthropicProvider) ProviderName() string {
	return "anthropic"
}

func (p *AnthropicProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("claude-3-5-haiku-20241022", credentials)
}

func (p *AnthropicProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewAnthropicLLM()
	default:
		return nil
	}
}
