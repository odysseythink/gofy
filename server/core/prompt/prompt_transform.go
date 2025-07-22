package prompt

import (
	"math"

	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	"mlib.com/gofy/server/core/memory"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	promptentities "mlib.com/gofy/server/entities/prompt"
)

type PromptTransform struct {
}

func (transform *PromptTransform) _append_chat_histories(
	mem *memory.TokenBufferMemory,
	memory_config *promptentities.MemoryConfig,
	prompt_messages []modelruntimeentities.PromptMessager,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) []modelruntimeentities.PromptMessager {
	rest_tokens := transform._calculate_rest_token(prompt_messages, model_config)
	histories := transform._get_history_messages_list_from_memory(mem, memory_config, rest_tokens)
	prompt_messages = append(prompt_messages, histories...)

	return prompt_messages
}
func (transform *PromptTransform) _calculate_rest_token(
	prompt_messages []modelruntimeentities.PromptMessager, model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) int {
	rest_tokens := 2000

	if _, ok := model_config.ModelSchema.ModelProperties[modelruntimeentities.ModelPropertyKey_CONTEXT_SIZE]; ok {
		if model_context_tokens, ok := model_config.ModelSchema.ModelProperties[modelruntimeentities.ModelPropertyKey_CONTEXT_SIZE].(int); ok {
			model_instance := modelmanager.NewModelInstance(
				model_config.ProviderModelBundle, model_config.Model,
			)

			curr_message_tokens := model_instance.GetLLMNumTokens(prompt_messages, nil)

			max_tokens := 0
			for _, parameter_rule := range model_config.ModelSchema.ParameterRules {
				if parameter_rule.Name == "max_tokens" ||
					(parameter_rule.UseTemplate != "" && parameter_rule.UseTemplate == "max_tokens") {
					if _, ok := model_config.Parameters[parameter_rule.Name]; ok {
						if _, ok := model_config.Parameters[parameter_rule.Name].(int); ok {
							max_tokens = model_config.Parameters[parameter_rule.Name].(int)
						}
					}
					if max_tokens == 0 {
						if _, ok := model_config.Parameters[parameter_rule.UseTemplate]; ok {
							if _, ok := model_config.Parameters[parameter_rule.UseTemplate].(int); ok {
								max_tokens = model_config.Parameters[parameter_rule.UseTemplate].(int)
							}
						}
					}
				}
			}
			rest_tokens = model_context_tokens - max_tokens - curr_message_tokens
			rest_tokens = int(math.Max(float64(rest_tokens), float64(0)))
		}
	}
	return rest_tokens
}
func (transform *PromptTransform) _get_history_messages_from_memory(
	mem *memory.TokenBufferMemory,
	memory_config *promptentities.MemoryConfig,
	max_token_limit int,
	human_prefix string,
	ai_prefix string,
) string {
	kwargs := map[string]any{"max_token_limit": max_token_limit}

	if human_prefix != "" {
		kwargs["human_prefix"] = human_prefix
	}

	if ai_prefix != "" {
		kwargs["ai_prefix"] = ai_prefix
	}
	message_limit := 0
	if memory_config.Window.Enabled && memory_config.Window.Size > 0 {
		message_limit = memory_config.Window.Size
	}

	return mem.GetHistoryPromptText(human_prefix, ai_prefix, max_token_limit, message_limit)
}
func (transform *PromptTransform) _get_history_messages_list_from_memory(
	mem *memory.TokenBufferMemory, memory_config *promptentities.MemoryConfig, max_token_limit int,
) []modelruntimeentities.PromptMessager {
	message_limit := 0
	if memory_config.Window.Enabled && memory_config.Window.Size > 0 {
		message_limit = memory_config.Window.Size
	}
	return mem.GetHistoryPromptMessages(
		max_token_limit,
		message_limit,
	)
}
