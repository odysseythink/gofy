package wenxin

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type WenxinProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&WenxinProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *WenxinProvider) ProviderName() string {
	return "wenxin"
}

func (p *WenxinProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Wenxin uses model-level credentials only
}

func (p *WenxinProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewWenxinLLM()
	default:
		return nil
	}
}
