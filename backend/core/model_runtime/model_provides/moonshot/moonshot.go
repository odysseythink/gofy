package moonshot

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
)

type MoonshotProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&MoonshotProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *MoonshotProvider) ProviderName() string {
	return "moonshot"
}

func (p *MoonshotProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("moonshot-v1-8k", credentials)
}

func (p *MoonshotProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewMoonshotLLM()
	default:
		return nil
	}
}
