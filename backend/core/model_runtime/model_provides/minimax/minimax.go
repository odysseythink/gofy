package minimax

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
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
