package vertex_ai

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type VertexAIProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&VertexAIProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *VertexAIProvider) ProviderName() string {
	return "vertex_ai"
}

func (p *VertexAIProvider) ValidateProviderCredentials(credentials map[string]any) {
	project, _ := credentials["vertex_project_id"].(string)
	location, _ := credentials["vertex_location"].(string)
	key, _ := credentials["vertex_service_account_key"].(string)
	if project == "" || location == "" || key == "" {
		panic(fmt.Errorf("vertex_project_id, vertex_location and vertex_service_account_key are required"))
	}
}

func (p *VertexAIProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewVertexAILLM()
	default:
		return nil
	}
}
