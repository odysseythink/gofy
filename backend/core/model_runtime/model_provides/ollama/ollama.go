package ollama

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type OllamaProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&OllamaProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *OllamaProvider) ProviderName() string {
	return "ollama"
}

func (p *OllamaProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Ollama uses model-level credentials only
}

func (p *OllamaProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewOllamaLLM()
	default:
		return nil
	}
}
