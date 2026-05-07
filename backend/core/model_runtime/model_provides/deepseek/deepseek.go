package deepseek

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type DeepSeekProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&DeepSeekProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *DeepSeekProvider) ProviderName() string {
	return "deepseek"
}

func (p *DeepSeekProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("deepseek-chat", credentials)
}

func (p *DeepSeekProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewDeepSeekLLM()
	default:
		return nil
	}
}
