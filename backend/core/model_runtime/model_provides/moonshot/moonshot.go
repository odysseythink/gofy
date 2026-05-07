package moonshot

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
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
