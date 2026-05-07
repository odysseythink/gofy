package replicate

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type ReplicateProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&ReplicateProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *ReplicateProvider) ProviderName() string {
	return "replicate"
}

func (p *ReplicateProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Replicate uses model-level credentials only
}

func (p *ReplicateProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewReplicateLLM()
	default:
		return nil
	}
}
