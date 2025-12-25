package tongyi

import (
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	"mlib.com/gofy/server/core/model_runtime/model_provides/ladder/llm"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/global"
	"mlib.com/mlog"
)

func init() {
	global.RegisgterModelProvider(&LadderProvide{
		BaseModelProvide: &base.BaseModelProvide{},
	})
}

type LadderProvide struct {
	*base.BaseModelProvide
}

func (provider *LadderProvide) ProviderName() string {
	return "ladder"
}

func (provider *LadderProvide) ValidateProviderCredentials(credentials map[string]any) {
	/*
	   Validate provider credentials

	   if validate failed, raise exception

	   :param credentials: provider credentials, credentials form defined in `provider_credential_schema`.
	*/
	// try:
	model_instance := provider.GetModelInstance(modelruntimeenumtypes.Model_LLM)
	mlog.Debug("validate credentials doubao-seed-1-6-thinking")
	// Use `doubao-seed-1-6-thinking` model for validate,
	model_instance.ValidateCredentials("doubao-seed-1-6-thinking", credentials)

	// except CredentialsValidateFailedError as ex:
	//     raise ex
	// except Exception as ex:
	//     logger.exception(f"{self.get_provider_schema().provider} credentials validate failed")
	//     raise ex
}

func (provider *LadderProvide) GetModelInstance(model_type modelruntimeenumtypes.ModelType) modelruntimeentities.AIModeler {
	/*
	   Get model instance
	   :param model_type: model type defined in `ModelType`
	   :return:
	*/
	// get dirname of the current path

	switch model_type {
	case modelruntimeenumtypes.Model_LLM:
		return &llm.LadderLargeLanguageModel{
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
