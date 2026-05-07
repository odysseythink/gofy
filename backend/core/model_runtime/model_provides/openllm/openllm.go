package openllm

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type OpenLLMProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&OpenLLMProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *OpenLLMProvider) ProviderName() string {
	return "openllm"
}

func (p *OpenLLMProvider) ValidateProviderCredentials(credentials map[string]any) {
	// OpenLLM uses model-level credentials only
}

func (p *OpenLLMProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewOpenLLMLLM()
	default:
		return nil
	}
}
