package zhipuai

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
)

type ZhipuAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&ZhipuAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *ZhipuAIProvider) ProviderName() string {
	return "zhipuai"
}

func (p *ZhipuAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("glm-4-flash", credentials)
}

func (p *ZhipuAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewZhipuAILLM()
	default:
		return nil
	}
}
