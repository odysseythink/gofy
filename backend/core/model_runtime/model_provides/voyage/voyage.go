package voyage

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

const voyageEndpoint = "https://api.voyageai.com/v1"

type VoyageProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&VoyageProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *VoyageProvider) ProviderName() string {
	return "voyage"
}

func (p *VoyageProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Voyage uses model-level credentials only
}

func (p *VoyageProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_TEXT_EMBEDDING:
		return oaicompat.NewSimpleEmbedding(
			"voyage",
			16000,
			func(map[string]any) string { return voyageEndpoint },
			func(credentials map[string]any) string {
				apiKey, _ := credentials["api_key"].(string)
				return apiKey
			},
		)
	default:
		return nil
	}
}
