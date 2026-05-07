package jina

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

const jinaEndpoint = "https://api.jina.ai/v1"

type JinaProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&JinaProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *JinaProvider) ProviderName() string {
	return "jina"
}

func (p *JinaProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Jina uses model-level credentials only
}

func (p *JinaProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_TEXT_EMBEDDING:
		return oaicompat.NewSimpleEmbedding(
			"jina",
			8192,
			func(map[string]any) string { return jinaEndpoint },
			func(credentials map[string]any) string {
				apiKey, _ := credentials["api_key"].(string)
				return apiKey
			},
		)
	default:
		return nil
	}
}
