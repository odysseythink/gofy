package workflow

import (
	"slices"
	"sort"

	"mlib.com/gofy/server/core/app/config_manageres/base"
	fileupload "mlib.com/gofy/server/core/app/config_manageres/features/file_upload"
	texttospeech "mlib.com/gofy/server/core/app/config_manageres/features/text_to_speech"
	sensitivewordavoidance "mlib.com/gofy/server/core/app/config_manageres/sensitive_word_avoidance"
	workflowvariables "mlib.com/gofy/server/core/app/config_manageres/workflow_variables"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	"mlib.com/gofy/server/models"
)

type WorkflowAppConfig struct {
	*appconfigentities.WorkflowUIBasedAppConfig
}

type WorkflowAppConfigManager struct {
	*base.BaseAppConfigManager
}

func New() *WorkflowAppConfigManager {
	return &WorkflowAppConfigManager{
		BaseAppConfigManager: &base.BaseAppConfigManager{},
	}
}

func (mgr *WorkflowAppConfigManager) GetAppConfig(app_model *models.App, wf *models.Workflow) *appconfigentities.WorkflowUIBasedAppConfig {
	features_map_dict := wf.FeaturesDict()

	// app_mode = AppMode.value_of(app_model.mode)
	app_config := &appconfigentities.WorkflowUIBasedAppConfig{
		AppConfig: &appconfigentities.AppConfig{
			TenantID: app_model.TenantID,
			AppID:    app_model.ID,
			AppMode:  app_model.Mode,

			SensitiveWordAvoidance: (&sensitivewordavoidance.SensitiveWordAvoidanceConfigManager{}).Convert(features_map_dict),
			Variables:              (&workflowvariables.WorkflowVariablesConfigManage{}).Convert(wf),
			AdditionalFeatures:     mgr.ConvertFeatures(features_map_dict, app_model.Mode),
		},
		WorkflowID: wf.ID,
	}

	return app_config
}

func (mgr *WorkflowAppConfigManager) ConfigValidate(tenant_id string, config map[string]any, only_structure_validate bool) map[string]any {
	/*
	   Validate for workflow app model config

	   :param tenant_id: tenant id
	   :param config: app model config args
	   :param only_structure_validate: only validate the structure of the config
	*/
	related_config_keys := []string{}

	// file upload validation
	config, current_related_config_keys := (&fileupload.FileUploadConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// text_to_speech
	config, current_related_config_keys = (&texttospeech.TextToSpeechConfigManager{}).ValidateAndSetDefaults(config)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	// moderation validation
	config, current_related_config_keys, _ = (&sensitivewordavoidance.SensitiveWordAvoidanceConfigManager{}).ValidateAndSetDefaults(
		tenant_id, config, only_structure_validate,
	)
	related_config_keys = append(related_config_keys, current_related_config_keys...)

	sort.Strings(related_config_keys)
	related_config_keys = slices.Compact(related_config_keys)

	// Filter out extra parameters
	filtered_config := map[string]any{}
	for _, v := range related_config_keys {
		filtered_config[v] = config[v]
	}

	return filtered_config
}
