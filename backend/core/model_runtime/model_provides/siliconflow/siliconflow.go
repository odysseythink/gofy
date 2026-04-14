package siliconflow

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
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
