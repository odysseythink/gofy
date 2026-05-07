package bedrock

import (
	"fmt"

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

type BedrockProvider struct {
	*base.BaseModelProvide
}

func init() {
	global.RegisgterModelProvider(&BedrockProvider{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

func (p *BedrockProvider) ProviderName() string {
	return "bedrock"
}

func (p *BedrockProvider) ValidateProviderCredentials(credentials map[string]any) {
	region, _ := credentials["aws_region"].(string)
	ak, _ := credentials["aws_access_key_id"].(string)
	sk, _ := credentials["aws_secret_access_key"].(string)
	if region == "" || ak == "" || sk == "" {
		panic(fmt.Errorf("aws_region, aws_access_key_id and aws_secret_access_key are required"))
	}
}

func (p *BedrockProvider) GetModelInstance(modelType modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	switch modelType {
	case modelruntimeenumtypes.Model_LLM:
		return NewBedrockLLM()
	default:
		return nil
	}
}
