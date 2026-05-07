package upstage

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type UpstageProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&UpstageProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *UpstageProvider) ProviderName() string {
	return "upstage"
}

func (p *UpstageProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Upstage uses model-level credentials only
}

func (p *UpstageProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewUpstageLLM()
	default:
		return nil
	}
}
