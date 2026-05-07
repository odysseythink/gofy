package cohere

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type CohereProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&CohereProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *CohereProvider) ProviderName() string {
	return "cohere"
}

func (p *CohereProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Cohere uses model-level credentials only
}

func (p *CohereProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewCohereLLM()
	case modelruntimeenumtypes.Model_TEXT_EMBEDDING:
		return NewCohereEmbedding()
	default:
		return nil
	}
}
