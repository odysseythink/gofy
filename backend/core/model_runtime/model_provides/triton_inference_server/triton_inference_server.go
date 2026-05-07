package triton_inference_server

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type TritonInferenceServerProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&TritonInferenceServerProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *TritonInferenceServerProvider) ProviderName() string {
	return "triton_inference_server"
}

func (p *TritonInferenceServerProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Triton uses model-level credentials only
}

func (p *TritonInferenceServerProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewTritonLLM()
	default:
		return nil
	}
}
