package togetherai

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type TogetherAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&TogetherAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *TogetherAIProvider) ProviderName() string {
	return "togetherai"
}

func (p *TogetherAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	// together.ai uses model-level credentials only
}

func (p *TogetherAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewTogetherAILLM()
	default:
		return nil
	}
}
