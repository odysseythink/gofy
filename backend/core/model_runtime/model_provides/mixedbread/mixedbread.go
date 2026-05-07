package mixedbread

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

const mixedbreadEndpoint = "https://api.mixedbread.ai/v1"

type MixedbreadProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&MixedbreadProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *MixedbreadProvider) ProviderName() string {
	return "mixedbread"
}

func (p *MixedbreadProvider) ValidateProviderCredentials(credentials map[string]any) {
	// mixedbread uses model-level credentials only
}

func (p *MixedbreadProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_TEXT_EMBEDDING:
		return oaicompat.NewSimpleEmbedding(
			"mixedbread",
			512,
			func(map[string]any) string { return mixedbreadEndpoint },
			func(credentials map[string]any) string {
				apiKey, _ := credentials["api_key"].(string)
				return apiKey
			},
		)
	default:
		return nil
	}
}
