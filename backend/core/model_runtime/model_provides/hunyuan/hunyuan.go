package hunyuan

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type HunyuanProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&HunyuanProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *HunyuanProvider) ProviderName() string {
	return "hunyuan"
}

func (p *HunyuanProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Hunyuan uses model-level credentials only
}

func (p *HunyuanProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewHunyuanLLM()
	default:
		return nil
	}
}
