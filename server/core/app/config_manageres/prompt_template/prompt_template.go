package prompttemplate

import (
	"mlib.com/gofy/server/core/exceptions"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	"mlib.com/gofy/server/models"
)

type PromptTemplateConfigManager struct {
}

func (mgr *PromptTemplateConfigManager) Convert(config map[string]any) appconfigentities.PromptTemplateEntity {
	if _, ok := config["prompt_type"]; !ok {
		panic(exceptions.NewValueError("prompt_type is required"))
	}
	if _, ok := config["prompt_type"].(string); !ok {
		panic(exceptions.NewValueError("prompt_type must be type of string"))
	}
	prompt_type := appconfigenumtypes.PromptType(config["prompt_type"].(string))
	if prompt_type == appconfigenumtypes.Prompt_SIMPLE {
		simple_prompt_template := ""
		if _, ok := config["pre_prompt"]; ok {
			if _, ok := config["pre_prompt"].(string); ok {
				simple_prompt_template = config["pre_prompt"].(string)
			}
		}
		return appconfigentities.PromptTemplateEntity{
			PromptType:           prompt_type,
			SimplePromptTemplate: simple_prompt_template,
		}
	} else {

		return appconfigentities.PromptTemplateEntity{
			PromptType: prompt_type,
		}
	}
}
func (mgr *PromptTemplateConfigManager) ValidatePostPromptAndSetDefaults(config map[string]any) map[string]any {
	/*
	   Validate post_prompt and set defaults for prompt feature

	   :param config: app model config args
	*/
	// post_prompt
	if _, ok := config["post_prompt"]; !ok {
		config["post_prompt"] = ""
	} else {
		if _, ok := config["post_prompt"].(string); !ok {
			panic(exceptions.NewValueError("post_prompt must be of string type"))
		}
	}

	return config
}
func (mgr *PromptTemplateConfigManager) validate_and_set_defaults(app_mode models.AppMode, config map[string]any) (map[string]any, []string) {
	/*
	   Validate pre_prompt and set defaults for prompt feature
	   depending on the config['model']

	   :param app_mode: app mode
	   :param config: app model config args
	*/
	if _, ok := config["prompt_type"]; !ok {
		config["prompt_type"] = string(appconfigenumtypes.Prompt_SIMPLE)
	}
	if _, ok := config["prompt_type"].(string); !ok {
		panic(exceptions.NewValueError("prompt_type must be type of string"))
	}

	if !appconfigenumtypes.ValidatePromptType(config["prompt_type"].(string)) {
		panic(exceptions.NewValueError("prompt_type must be PromptType"))
	}
	// chat_prompt_config
	if _, ok := config["chat_prompt_config"]; !ok {
		config["chat_prompt_config"] = map[string]any{}
	}
	if _, ok := config["chat_prompt_config"].(map[string]any); !ok {
		panic(exceptions.NewValueError("chat_prompt_config must be of object type"))
	}

	// completion_prompt_config
	if _, ok := config["completion_prompt_config"]; !ok {
		config["completion_prompt_config"] = map[string]any{}
	}
	if _, ok := config["completion_prompt_config"].(map[string]any); !ok {
		panic(exceptions.NewValueError("completion_prompt_config must be of object type"))
	}

	// pre_prompt, for simple mode
	if _, ok := config["pre_prompt"]; !ok {
		config["pre_prompt"] = ""
	}

	if _, ok := config["pre_prompt"].(string); !ok {
		panic(exceptions.NewValueError("pre_prompt must be of string type"))
	}
	return config, []string{"prompt_type", "pre_prompt", "chat_prompt_config", "completion_prompt_config"}
}
