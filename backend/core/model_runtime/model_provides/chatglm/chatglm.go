package chatglm

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type ChatGLMProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&ChatGLMProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *ChatGLMProvider) ProviderName() string {
	return "chatglm"
}

func (p *ChatGLMProvider) ValidateProviderCredentials(credentials map[string]any) {
	// ChatGLM uses model-level credentials only
}

func (p *ChatGLMProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewChatGLMLLM()
	default:
		return nil
	}
}
