package google

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type GoogleProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&GoogleProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *GoogleProvider) ProviderName() string {
	return "google"
}

func (p *GoogleProvider) ValidateProviderCredentials(credentials map[string]any) {
	model_instance := p.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	model_instance.ValidateCredentials("gemini-2.0-flash", credentials)
}

func (p *GoogleProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewGoogleLLM()
	default:
		return nil
	}
}
