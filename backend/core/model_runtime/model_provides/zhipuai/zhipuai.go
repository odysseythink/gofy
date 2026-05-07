package zhipuai

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
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
