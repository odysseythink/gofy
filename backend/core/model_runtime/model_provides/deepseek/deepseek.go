package deepseek

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
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
