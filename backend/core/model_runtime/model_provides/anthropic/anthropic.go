package anthropic

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
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
