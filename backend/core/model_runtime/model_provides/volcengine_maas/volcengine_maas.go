package volcengine_maas

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type VolcengineMaasProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&VolcengineMaasProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *VolcengineMaasProvider) ProviderName() string {
	return "volcengine_maas"
}

func (p *VolcengineMaasProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Volcengine Ark uses model-level credentials only
}

func (p *VolcengineMaasProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewVolcengineMaasLLM()
	default:
		return nil
	}
}
