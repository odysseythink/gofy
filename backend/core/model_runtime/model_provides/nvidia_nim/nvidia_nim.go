package nvidia_nim

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type NvidiaNimProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&NvidiaNimProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *NvidiaNimProvider) ProviderName() string {
	return "nvidia_nim"
}

func (p *NvidiaNimProvider) ValidateProviderCredentials(credentials map[string]any) {
	// NIM uses model-level credentials only
}

func (p *NvidiaNimProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewNvidiaNimLLM()
	default:
		return nil
	}
}
