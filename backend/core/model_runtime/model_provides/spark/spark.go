package spark

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type SparkProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&SparkProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *SparkProvider) ProviderName() string {
	return "spark"
}

func (p *SparkProvider) ValidateProviderCredentials(credentials map[string]any) {
	// Spark uses model-level credentials only
}

func (p *SparkProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewSparkLLM()
	default:
		return nil
	}
}
