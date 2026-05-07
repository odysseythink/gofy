package fireworks

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type FireworksProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&FireworksProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *FireworksProvider) ProviderName() string {
	return "fireworks"
}

func (p *FireworksProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Fireworks uses model-level credentials only
}

func (p *FireworksProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewFireworksLLM()
	default:
		return nil
	}
}
