package openai_api_compatible

import (
	"fmt"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
)

// OpenAICompatibleProvider works with any OpenAI-compatible API endpoint.
type OpenAICompatibleProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&OpenAICompatibleProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *OpenAICompatibleProvider) ProviderName() string {
	return "openai_api_compatible"
}

func (p *OpenAICompatibleProvider) ValidateProviderCredentials(credentials map[string]any) {
	apiKey, _ := credentials["api_key"].(string)
	if apiKey == "" {
		panic(fmt.Errorf("api_key is required"))
	}
	endpoint, _ := credentials["endpoint_url"].(string)
	if endpoint == "" {
		panic(fmt.Errorf("endpoint_url is required"))
	}
}

func (p *OpenAICompatibleProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewOpenAICompatibleLLM()
	case modelruntimeenumtypes.Model_TEXT_EMBEDDING:
		return NewOpenAICompatibleEmbedding()
	default:
		mlog.Errorf("unsupported model type: %s", modelType)
		return nil
	}
}
