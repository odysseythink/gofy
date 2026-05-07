package localai

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type LocalAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&LocalAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *LocalAIProvider) ProviderName() string {
	return "localai"
}

func (p *LocalAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	// LocalAI uses model-level credentials only
}

func (p *LocalAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewLocalAILLM()
	default:
		return nil
	}
}
