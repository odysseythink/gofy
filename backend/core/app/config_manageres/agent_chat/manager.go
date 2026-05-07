package advancedchat

import (
	"slices"
	"sort"

	agentcfgmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/agent"
	"github.com/odysseythink/gofy/backend/core/app/config_manageres/base"
	datasetcfgmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/dataset"
	fileupload "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/file_upload"
	openingstatement "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/opening_statement"
	speechtotext "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/speech_to_text"
	suggestedquestionsafteranswer "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/suggested_questions_after_answer"
	texttospeech "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/text_to_speech"
	modelconfigmanages "github.com/odysseythink/gofy/backend/core/app/config_manageres/model"
	prompttemplatecfgmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/prompt_template"
	sensitivewordavoidancecfgmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/sensitive_word_avoidance"
	variablescfgmgr "github.com/odysseythink/gofy/backend/core/app/config_manageres/variables"
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	appconfigenumtypes "github.com/odysseythink/gofy/backend/enum_types/app_config"
	"github.com/odysseythink/gofy/backend/models"
)

type AgentChatAppConfigManager struct {
	*base.BaseAppConfigManager
}

func New() *AgentChatAppConfigManager {
	return &AgentChatAppConfigManager{
		BaseAppConfigManager: &base.BaseAppConfigManager{},
	}
}

func (mgr *AgentChatAppConfigManager) GetAppConfig(
	app_model *models.App,
	app_model_config *models.AppModelConfig,
	conversation *models.Conversation,
	override_config_dict map[string]any,
) *appconfigentities.AgentChatAppConfig {
	var config_from appconfigenumtypes.EasyUIBasedAppModelConfigFrom
	if len(override_config_dict) > 0 {
		config_from = appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS
	} else if conversation != nil {
		config_from = appconfigenumtypes.EasyUIBasedAppModelConfigFrom_CONVERSATION_SPECIFIC_CONFIG
	} else {
		config_from = appconfigenumtypes.EasyUIBasedAppModelConfigFrom_APP_LATEST_CONFIG
	}
	var config_dict map[string]any
	if config_from != appconfigenumtypes.EasyUIBasedAppModelConfigFrom_ARGS {
		config_dict = app_model_config.ToDict()
	} else {
		if len(override_config_dict) > 0 {
			config_dict = override_config_dict
		} else {
			config_dict = map[string]any{}
		}
	}
	app_mode := app_model.Mode
	app_config := &appconfigentities.AgentChatAppConfig{
		EasyUIBasedAppConfig: &appconfigentities.EasyUIBasedAppConfig{
			AppConfig: &appconfigentities.AppConfig{
				TenantID:           app_model.TenantID,
				AppID:              app_model.ID,
				AppMode:            app_model.Mode,
				AdditionalFeatures: mgr.ConvertFeatures(config_dict, app_mode),
				// Variables             : app_model.Variables             ,
				SensitiveWordAvoidance: (&sensitivewordavoidancecfgmgr.SensitiveWordAvoidanceConfigManager{}).Convert(config_dict),
			},
			AppModelConfigFrom: config_from,
			AppModelConfigID:   app_model_config.ID,
			AppModelConfigDict: config_dict,
			Model:              (&modelconfigmanages.ModelConfigManager{}).Convert(config_dict),
			PromptTemplate:     (&prompttemplatecfgmgr.PromptTemplateConfigManager{}).Convert(config_dict),
			Dataset:            (&datasetcfgmgr.DatasetConfigManager{}).Convert(config_dict),
		},
		Agent: (&agentcfgmgr.AgentConfigManager{}).Convert(config_dict),
	}

	app_config.Variables, app_config.ExternalDataVariables = (&variablescfgmgr.BasicVariablesConfigManager{}).Convert(config_dict)

	return app_config
}

func (mgr *AgentChatAppConfigManager) ConfigValidate(tenant_id string, config map[string]any, only_structure_validate bool) map[string]any {
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
	config, current_related_config_keys, _ = (&sensitivewordavoidancecfgmgr.SensitiveWordAvoidanceConfigManager{}).ValidateAndSetDefaults(
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
