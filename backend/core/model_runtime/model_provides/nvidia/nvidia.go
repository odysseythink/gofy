package nvidia

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type NvidiaProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&NvidiaProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *NvidiaProvider) ProviderName() string {
	return "nvidia"
}

func (p *NvidiaProvider) ValidateProviderCredentials(credentials map[string]any) {
	// NVIDIA uses model-level credentials only
}

func (p *NvidiaProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewNvidiaLLM()
	default:
		return nil
	}
}
