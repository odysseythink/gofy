package yi

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type YiProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&YiProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *YiProvider) ProviderName() string {
	return "yi"
}

func (p *YiProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Yi uses model-level credentials only
}

func (p *YiProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewYiLLM()
	default:
		return nil
	}
}
