package tongyi

import (
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/tongyi/llm"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/global"
)

func init() {
	global.RegisgterModelProvider(&TongyiProvide{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

type TongyiProvide struct {
	*base.BaseModelProvide
}

func (provider *TongyiProvide) ProviderName() string {
	return "tongyi"
}

func (provider *TongyiProvide) ValidateProviderCredentials(credentials map[string]any) {
	/*
	   Validate provider credentials

	   if validate failed, raise exception

	   :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
	*/
	// try:
	model_instance := provider.GetModelInstance(modelruntimeenumtypes.Model_LLM)

	// Use `qwen-turbo` model for validate,
	model_instance.ValidateCredentials("qwen-turbo", credentials)

	// except CredentialsValidateFailedError as ex:
	//     raise ex
	// except Exception as ex:
	//     logger.exception(f"{self.get_provider_schema().provider} credentials validate failed")
	//     raise ex
}

func (provider *TongyiProvide) GetModelInstance(model_type modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	/*
	   Get model instance
	   :param model_type: model type defined in `ModelType`
	   :return:
	*/
	// get dirname of the current path

	switch model_type {
	case modelruntimeenumtypes.Model_LLM:
		return &llm.TongyiLargeLanguageModel{
			LargeLanguageModel: &base.LargeLanguageModel{
				BaseAIModel: &base.BaseAIModel{
					ModeType: modelruntimeenumtypes.Model_LLM,
				},
			},
		}
	}
	// get the path of the model type classes
	return nil
}
