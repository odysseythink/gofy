package zhinao

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type ZhinaoProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&ZhinaoProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *ZhinaoProvider) ProviderName() string {
	return "zhinao"
}

func (p *ZhinaoProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Zhinao uses model-level credentials only
}

func (p *ZhinaoProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewZhinaoLLM()
	default:
		return nil
	}
}
