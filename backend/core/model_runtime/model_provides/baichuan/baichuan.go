package baichuan

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type BaichuanProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&BaichuanProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *BaichuanProvider) ProviderName() string {
	return "baichuan"
}

func (p *BaichuanProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Baichuan uses model-level credentials only
}

func (p *BaichuanProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewBaichuanLLM()
	default:
		return nil
	}
}
