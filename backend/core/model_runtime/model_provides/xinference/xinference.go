package xinference

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type XinferenceProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&XinferenceProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *XinferenceProvider) ProviderName() string {
	return "xinference"
}

func (p *XinferenceProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Xinference uses model-level credentials only
}

func (p *XinferenceProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewXinferenceLLM()
	default:
		return nil
	}
}
