package nomic

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

const nomicEndpoint = "https://api-atlas.nomic.ai/v1"

type NomicProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&NomicProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *NomicProvider) ProviderName() string {
	return "nomic"
}

func (p *NomicProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Nomic uses model-level credentials only
}

func (p *NomicProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_TEXT_EMBEDDING:
		return oaicompat.NewSimpleEmbedding(
			"nomic",
			8192,
			func(map[string]any) string { return nomicEndpoint },
			func(credentials map[string]any) string {
				apiKey, _ := credentials["api_key"].(string)
				return apiKey
			},
		)
	default:
		return nil
	}
}
