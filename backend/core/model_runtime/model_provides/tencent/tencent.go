package tencent

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type TencentProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&TencentProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *TencentProvider) ProviderName() string {
	return "tencent"
}

func (p *TencentProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Tencent LKEAP uses model-level credentials only
}

func (p *TencentProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewTencentLLM()
	default:
		return nil
	}
}
