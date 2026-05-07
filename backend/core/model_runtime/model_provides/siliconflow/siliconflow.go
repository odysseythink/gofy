package siliconflow

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type SiliconFlowProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&SiliconFlowProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *SiliconFlowProvider) ProviderName() string {
	return "siliconflow"
}

func (p *SiliconFlowProvider) ValidateProviderCredentials(credentials map[string]any) {
	// SiliconFlow uses model-level credentials only
}

func (p *SiliconFlowProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewSiliconFlowLLM()
	default:
		return nil
	}
}
