package prompttemplate

import (
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	promptenumtypes "mlib.com/gofy/server/enum_types/prompt"
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
		var advanced_chat_prompt_template *appconfigentities.AdvancedChatPromptTemplateEntity
		chat_prompt_config := map[string]any{}
		if _, ok := config["chat_prompt_config"]; ok {
			if _, ok := config["chat_prompt_config"].(map[string]any); ok {
				chat_prompt_config = config["chat_prompt_config"].(map[string]any)
			}
		}
		if len(chat_prompt_config) > 0 {
			chat_prompt_messages := []appconfigentities.AdvancedChatMessageEntity{}
			if _, ok := chat_prompt_config["prompt"]; ok {
				if _, ok := chat_prompt_config["prompt"].([]any); ok {
					for _, v := range chat_prompt_config["prompt"].([]any) {
						if _, ok := v.(map[string]any); ok {
							text := ""
							if _, ok := v.(map[string]any)["text"]; ok {
								if _, ok := v.(map[string]any)["text"].(string); ok {
									text = v.(map[string]any)["text"].(string)
								}
							}
							role := ""
							if _, ok := v.(map[string]any)["role"]; ok {
								if _, ok := v.(map[string]any)["role"].(string); ok {
									role = v.(map[string]any)["role"].(string)
								}
							}
							chat_prompt_messages = append(chat_prompt_messages, appconfigentities.AdvancedChatMessageEntity{
								Text: text,
								Role: modelruntimeentities.PromptMessageRole(role),
							})
						}
					}
				} else if _, ok := chat_prompt_config["prompt"].([]map[string]any); ok {
					for _, v := range chat_prompt_config["prompt"].([]map[string]any) {
						text := ""
						if _, ok := v["text"]; ok {
							if _, ok := v["text"].(string); ok {
								text = v["text"].(string)
							}
						}
						role := ""
						if _, ok := v["role"]; ok {
							if _, ok := v["role"].(string); ok {
								role = v["role"].(string)
							}
						}
						chat_prompt_messages = append(chat_prompt_messages, appconfigentities.AdvancedChatMessageEntity{
							Text: text,
							Role: modelruntimeentities.PromptMessageRole(role),
						})
					}
				}
			}
			advanced_chat_prompt_template = &appconfigentities.AdvancedChatPromptTemplateEntity{Messages: chat_prompt_messages}
		}
		var advanced_completion_prompt_template *appconfigentities.AdvancedCompletionPromptTemplateEntity
		completion_prompt_config := map[string]any{}
		if _, ok := config["completion_prompt_config"]; ok {
			if _, ok := config["completion_prompt_config"].(map[string]any); ok {
				completion_prompt_config = config["completion_prompt_config"].(map[string]any)
			}
		}
		if len(completion_prompt_config) > 0 {
			var prompt map[string]any
			if _, ok := completion_prompt_config["prompt"]; ok {
				if _, ok := completion_prompt_config["prompt"].(map[string]any); ok {
					prompt = completion_prompt_config["prompt"].(map[string]any)
				}
			}
			prompt_text := ""
			if _, ok := prompt["text"]; ok {
				if _, ok := prompt["text"].(string); ok {
					prompt_text = prompt["text"].(string)
				}
			}
			advanced_completion_prompt_template = &appconfigentities.AdvancedCompletionPromptTemplateEntity{
				Prompt: prompt_text,
				RolePrefix: &appconfigentities.RolePrefixEntity{
					User:      "",
					Assistant: "",
				},
			}
			if _, ok := completion_prompt_config["conversation_histories_role"]; ok {
				if _, ok := completion_prompt_config["conversation_histories_role"].(map[string]any); ok {
					conversation_histories_role := completion_prompt_config["conversation_histories_role"].(map[string]any)
					if _, ok := conversation_histories_role["user_prefix"]; ok {
						if _, ok := conversation_histories_role["user_prefix"].(string); ok {
							advanced_completion_prompt_template.RolePrefix.User = conversation_histories_role["user_prefix"].(string)
						}
					}
					if _, ok := conversation_histories_role["assistant_prefix"]; ok {
						if _, ok := conversation_histories_role["assistant_prefix"].(string); ok {
							advanced_completion_prompt_template.RolePrefix.Assistant = conversation_histories_role["assistant_prefix"].(string)
						}
					}
				}
			}
		}
		return appconfigentities.PromptTemplateEntity{
			PromptType:                       prompt_type,
			AdvancedChatPromptTemplate:       advanced_chat_prompt_template,
			AdvancedCompletionPromptTemplate: advanced_completion_prompt_template,
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

	if config["prompt_type"] == string(appconfigenumtypes.Prompt_ADVANCED) {
		// if not config["chat_prompt_config"] and not config["completion_prompt_config"]{
		// 	raise ValueError(
		// 		"chat_prompt_config or completion_prompt_config is required when prompt_type is advanced"
		// 	)
		// }
		var model_config map[string]any
		if _, ok := config["model"]; ok {
			if _, ok := config["model"].(map[string]any); ok {
				model_config = config["model"].(map[string]any)
			}
		}
		model_mode := ""
		if _, ok := model_config["mode"]; ok {
			if _, ok := model_config["mode"].(string); ok {
				model_mode = model_config["mode"].(string)
			}
		}

		if !promptenumtypes.ValidateModelModeType(model_mode) {
			panic(exceptions.NewValueError(fmt.Sprintf("model.mode must be ModelModeType when prompt_type is advanced")))
		}

		if app_mode == models.AppMode_CHAT && model_mode == string(promptenumtypes.ModelMode_COMPLETION) {
			completion_prompt_config := config["completion_prompt_config"].(map[string]any)
			if _, ok := completion_prompt_config["conversation_histories_role"]; ok {
				if _, ok := completion_prompt_config["conversation_histories_role"].(map[string]any); ok {
					conversation_histories_role := completion_prompt_config["conversation_histories_role"].(map[string]any)
					if _, ok := conversation_histories_role["user_prefix"]; ok {
						if _, ok := conversation_histories_role["user_prefix"].(string); !ok {
							conversation_histories_role["user_prefix"] = "Human"
						}
					} else {
						conversation_histories_role["user_prefix"] = "Human"
					}
					if _, ok := conversation_histories_role["assistant_prefix"]; ok {
						if _, ok := conversation_histories_role["assistant_prefix"].(string); !ok {
							conversation_histories_role["assistant_prefix"] = "Assistant"
						}
					} else {
						conversation_histories_role["assistant_prefix"] = "Assistant"
					}
					completion_prompt_config["conversation_histories_role"] = conversation_histories_role
				} else {
					completion_prompt_config["conversation_histories_role"] = map[string]any{
						"user_prefix":      "Human",
						"assistant_prefix": "Assistant",
					}
				}
			} else {
				completion_prompt_config["conversation_histories_role"] = map[string]any{
					"user_prefix":      "Human",
					"assistant_prefix": "Assistant",
				}
			}
			config["completion_prompt_config"] = completion_prompt_config
		}
		if model_mode == string(promptenumtypes.ModelMode_CHAT) {
			prompt_list := []any{}
			if _, ok := config["chat_prompt_config"].(map[string]any)["prompt"]; ok {
				if _, ok := config["chat_prompt_config"].(map[string]any)["prompt"].([]any); ok {
					prompt_list = config["chat_prompt_config"].(map[string]any)["prompt"].([]any)
				}
			}

			if len(prompt_list) > 10 {
				panic(exceptions.NewValueError("prompt messages must be less than 10"))
			}
		}
	} else {
		// pre_prompt, for simple mode
		if _, ok := config["pre_prompt"]; !ok {
			config["pre_prompt"] = ""
		}

		if _, ok := config["pre_prompt"].(string); !ok {
			panic(exceptions.NewValueError("pre_prompt must be of string type"))
		}
	}
	return config, []string{"prompt_type", "pre_prompt", "chat_prompt_config", "completion_prompt_config"}
}
