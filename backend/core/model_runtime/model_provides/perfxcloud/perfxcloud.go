package perfxcloud

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type PerfXCloudProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&PerfXCloudProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *PerfXCloudProvider) ProviderName() string {
	return "perfxcloud"
}

func (p *PerfXCloudProvider) ValidateProviderCredentials(credentials map[string]any) {
	// PerfXCloud uses model-level credentials only
}

func (p *PerfXCloudProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewPerfXCloudLLM()
	default:
		return nil
	}
}
