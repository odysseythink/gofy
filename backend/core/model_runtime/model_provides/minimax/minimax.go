package minimax

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
)

type MiniMaxProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&MiniMaxProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *MiniMaxProvider) ProviderName() string {
	return "minimax"
}

func (p *MiniMaxProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("abab6.5s-chat", credentials)
}

func (p *MiniMaxProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewMiniMaxLLM()
	default:
		return nil
	}
}
