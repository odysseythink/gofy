package advancedchat

import (
	"slices"
	"sort"

	"mlib.com/gofy/server/core/app/config_manages/base"
	fileupload "mlib.com/gofy/server/core/app/config_manages/features/file_upload"
	openingstatement "mlib.com/gofy/server/core/app/config_manages/features/opening_statement"
	speechtotext "mlib.com/gofy/server/core/app/config_manages/features/speech_to_text"
	suggestedquestionsafteranswer "mlib.com/gofy/server/core/app/config_manages/features/suggested_questions_after_answer"
	texttospeech "mlib.com/gofy/server/core/app/config_manages/features/text_to_speech"
	sensitivewordavoidance "mlib.com/gofy/server/core/app/config_manages/sensitive_word_avoidance"
	wfvariables "mlib.com/gofy/server/core/app/config_manages/workflow_variables"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	"mlib.com/gofy/server/models"
)

type AgentChatAppConfigManage struct {
	*base.BaseAppConfigManage
}

func New() *AgentChatAppConfigManage {
	return &AgentChatAppConfigManage{
		BaseAppConfigManage: &base.BaseAppConfigManage{},
	}
}

func (mgr *AgentChatAppConfigManage) GetAppConfig(
	app_model *models.App,
	app_model_config *models.AppModelConfig,
	conversation *models.Conversation,
	override_config_dict map[string]any,
) *appconfigentities.AgentChatAppConfig {
	var config_from appconfigenumtypes.EasyUIBasedAppModelConfigFrom
        if len(override_config_dict) > 0 {
            config_from = appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS
        }else if conversation != nil{
            config_from = appconfigenumtypes.EasyUIBasedAppModelConfigFrom_CONVERSATION_SPECIFIC_CONFIG
        }else{
            config_from = appconfigenumtypes.EasyUIBasedAppModelConfigFrom_APP_LATEST_CONFIG
}
        if config_from != appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS{
            app_model_config_dict = app_model_config.to_dict()
            config_dict = app_model_config_dict.copy()
        }else{
            config_dict = override_config_dict or {}
}
        app_mode = AppMode.value_of(app_model.mode)
        app_config = AgentChatAppConfig(
            tenant_id=app_model.tenant_id,
            app_id=app_model.id,
            app_mode=app_mode,
            app_model_config_from=config_from,
            app_model_config_id=app_model_config.id,
            app_model_config_dict=config_dict,
            model=ModelConfigManager.convert(config=config_dict),
            prompt_template=PromptTemplateConfigManager.convert(config=config_dict),
            sensitive_word_avoidance=SensitiveWordAvoidanceConfigManager.convert(config=config_dict),
            dataset=DatasetConfigManager.convert(config=config_dict),
            agent=AgentConfigManager.convert(config=config_dict),
            additional_features=cls.convert_features(config_dict, app_mode),
        )

        app_config.variables, app_config.external_data_variables = BasicVariablesConfigManager.convert(
            config=config_dict
        )

        return app_config
	return app_config
}

func (mgr *AgentChatAppConfigManage) ConfigValidate(tenant_id string, config map[string]any, only_structure_validate bool) map[string]any {
	/*
	   Validate for advanced chat app model config

	   :param tenant_id: tenant id
	   :param config: app model config args
	   :param only_structure_validate: if True, only structure validation will be performed
	*/
	related_config_keys := []string{}

	// file upload validation
	config, current_related_config_keys := (&fileupload.FileUploadConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// opening_statement
	config, current_related_config_keys = (&openingstatement.OpeningStatementConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// suggested_questions_after_answer
	config, current_related_config_keys, _ = (&suggestedquestionsafteranswer.SuggestedQuestionsAfterAnswerConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// speech_to_text
	config, current_related_config_keys = (&speechtotext.SpeechToTextConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// text_to_speech
	config, current_related_config_keys = (&texttospeech.TextToSpeechConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// // return retriever resource
	// config, current_related_config_keys = RetrievalResourceConfigManager.ValidateAndSetDefaults(config)
	// related_config_keys = append(related_config_keys,current_related_config_keys...)

	// moderation validation
	config, current_related_config_keys, _ = (&sensitivewordavoidance.SensitiveWordAvoidanceConfigManage{}).ValidateAndSetDefaults(
		tenant_id, config, only_structure_validate,
	)
	related_config_keys = append(related_config_keys, current_related_config_keys...)
	sort.Strings(related_config_keys)
	related_config_keys = slices.Compact(related_config_keys)

	filtered_config := map[string]any{}
	// Filter out extra parameters
	for _, v := range related_config_keys {
		filtered_config[v] = config[v]
	}

	return filtered_config
}
