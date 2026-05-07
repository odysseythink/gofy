package advancedchat

import (
	"slices"
	"sort"

	"github.com/odysseythink/gofy/backend/core/app/config_manageres/base"
	fileupload "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/file_upload"
	openingstatement "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/opening_statement"
	speechtotext "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/speech_to_text"
	suggestedquestionsafteranswer "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/suggested_questions_after_answer"
	texttospeech "github.com/odysseythink/gofy/backend/core/app/config_manageres/features/text_to_speech"
	sensitivewordavoidance "github.com/odysseythink/gofy/backend/core/app/config_manageres/sensitive_word_avoidance"
	wfvariables "github.com/odysseythink/gofy/backend/core/app/config_manageres/workflow_variables"
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	"github.com/odysseythink/gofy/backend/models"
)

type AdvancedChatAppConfigManager struct {
	*base.BaseAppConfigManager
}

func New() *AdvancedChatAppConfigManager {
	return &AdvancedChatAppConfigManager{
		BaseAppConfigManager: &base.BaseAppConfigManager{},
	}
}

func (mgr *AdvancedChatAppConfigManager) GetAppConfig(app_model *models.App, wf *models.Workflow) *appconfigentities.AdvancedChatAppConfig {
	features_dict := wf.FeaturesDict()

	app_config := &appconfigentities.AdvancedChatAppConfig{
		WorkflowUIBasedAppConfig: &appconfigentities.WorkflowUIBasedAppConfig{
			AppConfig: &appconfigentities.AppConfig{
				TenantID: app_model.TenantID,
				AppID:    app_model.ID,
				AppMode:  app_model.Mode,

				SensitiveWordAvoidance: (&sensitivewordavoidance.SensitiveWordAvoidanceConfigManager{}).Convert(features_dict),
				Variables:              (&wfvariables.WorkflowVariablesConfigManager{}).Convert(wf),
				AdditionalFeatures:     mgr.ConvertFeatures(features_dict, app_model.Mode),
			},
			WorkflowID: wf.ID,
		},
	}

	return app_config
}

func (mgr *AdvancedChatAppConfigManager) ConfigValidate(tenant_id string, config map[string]any, only_structure_validate bool) map[string]any {
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
	config, current_related_config_keys, _ = (&sensitivewordavoidance.SensitiveWordAvoidanceConfigManager{}).ValidateAndSetDefaults(
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
