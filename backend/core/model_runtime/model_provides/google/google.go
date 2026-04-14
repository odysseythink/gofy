package google

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
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
