package groq

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type GroqProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&GroqProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *GroqProvider) ProviderName() string {
	return "groq"
}

func (p *GroqProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Groq uses model-level credentials only
}

func (p *GroqProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewGroqLLM()
	default:
		return nil
	}
}
