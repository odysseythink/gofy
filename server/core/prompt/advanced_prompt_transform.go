package prompt

import (
	"maps"
	"slices"
	"strings"

	"mlib.com/gofy/server/core/file"
	"mlib.com/gofy/server/core/memory"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	promptentities "mlib.com/gofy/server/entities/prompt"
	workflowentities "mlib.com/gofy/server/entities/workflow"
)

type AdvancedPromptTransform[T []*promptentities.ChatModelMessage | *promptentities.CompletionModelPromptTemplate] struct {
	*PromptTransform
	WithVariableTmpl  bool
	ImageDetailConfig modelruntimeentities.ImagePromptMessageContentDETAIL
}

func NewAdvancedPromptTransform[T []*promptentities.ChatModelMessage | *promptentities.CompletionModelPromptTemplate](
	with_variable_tmpl bool,
	image_detail_config modelruntimeentities.ImagePromptMessageContentDETAIL,
) *AdvancedPromptTransform[T] {
	if string(image_detail_config) == "" {
		image_detail_config = modelruntimeentities.ImagePromptMessageContentDETAIL_LOW
	}
	return &AdvancedPromptTransform[T]{
		WithVariableTmpl:  with_variable_tmpl,
		ImageDetailConfig: image_detail_config,
	}
}

func (transform *AdvancedPromptTransform[T]) GetPrompt(
	prompt_template T,
	inputs map[string]string,
	query string,
	files []*file.File,
	context string,
	memory_config *promptentities.MemoryConfig,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) []modelruntimeentities.PromptMessager {
	var prompt_messages []modelruntimeentities.PromptMessager
	if real_prompt_template, ok := any(prompt_template).(*promptentities.CompletionModelPromptTemplate); ok {
		prompt_messages = transform._get_completion_model_prompt_messages(
			real_prompt_template,
			inputs,
			query,
			files,
			context,
			memory_config,
			mem,
			model_config,
		)
	} else if real_prompt_template, ok := any(prompt_template).([]*promptentities.ChatModelMessage); ok {
		prompt_messages = transform._get_chat_model_prompt_messages(
			real_prompt_template,
			inputs,
			query,
			files,
			context,
			memory_config,
			mem,
			model_config,
		)
	}
	return prompt_messages
}
func (transform *AdvancedPromptTransform[T]) _get_completion_model_prompt_messages(
	prompt_template *promptentities.CompletionModelPromptTemplate,
	inputs map[string]string,
	query string,
	files []*file.File,
	context string,
	memory_config *promptentities.MemoryConfig,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) []modelruntimeentities.PromptMessager {
	raw_prompt := prompt_template.Text
	prompt_messages := []modelruntimeentities.PromptMessager{}

	parser := promptutils.NewPromptTemplateParser(raw_prompt, transform.WithVariableTmpl)
	prompt_inputs := map[string]string{}
	for _, k := range parser.VariableKeys {
		if _, ok := inputs[k]; ok {
			prompt_inputs[k] = inputs[k]
		}
	}
	prompt_inputs = transform._set_context_variable(context, parser, prompt_inputs)
	if mem != nil && memory_config != nil && memory_config.RolePrefix != nil {
		role_prefix := memory_config.RolePrefix
		prompt_inputs = transform._set_histories_variable(
			mem,
			memory_config,
			raw_prompt,
			role_prefix,
			parser,
			prompt_inputs,
			model_config,
		)
	}
	if query != "" {
		prompt_inputs = transform._set_query_variable(query, parser, prompt_inputs)
	}
	prompt := parser.Format(prompt_inputs, false)

	// if len(files) > 0{
	//     prompt_message_contents: list[PromptMessageContent] = []
	//     prompt_message_contents.append(TextPromptMessageContent(data=prompt))
	//     for file in files{
	//         prompt_message_contents.append(file_manager.to_prompt_message_content(file))
	// 	}
	//     prompt_messages.append(UserPromptMessage(content=prompt_message_contents))
	// } else {
	prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(prompt, ""))
	// }
	return prompt_messages
}
func (transform *AdvancedPromptTransform[T]) _get_chat_model_prompt_messages(
	prompt_template []*promptentities.ChatModelMessage,
	inputs map[string]string,
	query string,
	files []*file.File,
	context string,
	memory_config *promptentities.MemoryConfig,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) []modelruntimeentities.PromptMessager {
	prompt_messages := []modelruntimeentities.PromptMessager{}
	prompt_inputs := map[string]string{}
	for _, prompt_item := range prompt_template {
		raw_prompt := prompt_item.Text
		prompt := ""
		if transform.WithVariableTmpl {
			vp := workflowentities.NewVariablePool(nil, nil, nil, nil)
			for k, v := range inputs {
				if strings.HasPrefix(k, "#") {
					vp.Add(strings.Split(k[1:], "."), v)
				}
			}
			raw_prompt = strings.ReplaceAll(raw_prompt, "{{#context#}}", context)
			prompt = vp.ConvertTemplate(raw_prompt).Text()
		} else {
			parser := promptutils.NewPromptTemplateParser(raw_prompt, transform.WithVariableTmpl)

			for _, k := range parser.VariableKeys {
				if _, ok := inputs[k]; ok {
					prompt_inputs[k] = inputs[k]
				}
			}

			prompt_inputs = transform._set_context_variable(
				context, parser, prompt_inputs,
			)
			prompt = parser.Format(prompt_inputs, false)
		}

		if prompt_item.Role == modelruntimeentities.PromptMessageRole_USER {
			prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(prompt, ""))
		} else if prompt_item.Role == modelruntimeentities.PromptMessageRole_SYSTEM && prompt != "" {
			prompt_messages = append(prompt_messages, modelruntimeentities.NewSystemPromptMessage(prompt, ""))
		} else if prompt_item.Role == modelruntimeentities.PromptMessageRole_ASSISTANT {
			prompt_messages = append(prompt_messages, modelruntimeentities.NewAssistantPromptMessage(prompt, "", nil))
		}
	}
	if query != "" && memory_config != nil && memory_config.QueryPromptTemplate != "" {
		parser := promptutils.NewPromptTemplateParser(
			memory_config.QueryPromptTemplate, transform.WithVariableTmpl,
		)
		for _, k := range parser.VariableKeys {
			if _, ok := inputs[k]; ok {
				prompt_inputs[k] = inputs[k]
			}
		}

		prompt_inputs["#sys.query#"] = query
		prompt_inputs = transform._set_context_variable(context, parser, prompt_inputs)
		query = parser.Format(prompt_inputs, false)
	}
	if mem != nil && memory_config != nil {
		prompt_messages = transform._append_chat_histories(mem, memory_config, prompt_messages, model_config)
		prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(query, ""))
	} else if query != "" {
		prompt_messages = append(prompt_messages, modelruntimeentities.NewUserPromptMessage(query, ""))
	}
	return prompt_messages
}
func (transform *AdvancedPromptTransform[T]) _set_context_variable(
	context string, parser *promptutils.PromptTemplateParser, prompt_inputs map[string]string,
) map[string]string {
	if slices.Contains(parser.VariableKeys, "#context#") {
		if prompt_inputs == nil {
			prompt_inputs = make(map[string]string)
		}
		if context != "" {
			prompt_inputs["#context#"] = context
		} else {
			prompt_inputs["#context#"] = ""
		}
	}
	return prompt_inputs
}
func (transform *AdvancedPromptTransform[T]) _set_query_variable(
	query string, parser *promptutils.PromptTemplateParser, prompt_inputs map[string]string,
) map[string]string {
	if slices.Contains(parser.VariableKeys, "#query#") {
		if prompt_inputs == nil {
			prompt_inputs = make(map[string]string)
		}
		if query != "" {
			prompt_inputs["#query#"] = query
		} else {
			prompt_inputs["#query#"] = ""
		}
	}
	return prompt_inputs
}
func (transform *AdvancedPromptTransform[T]) _set_histories_variable(
	mem *memory.TokenBufferMemory,
	memory_config *promptentities.MemoryConfig,
	raw_prompt string,
	role_prefix *promptentities.RolePrefix,
	parser *promptutils.PromptTemplateParser,
	prompt_inputs map[string]string,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) map[string]string {
	if slices.Contains(parser.VariableKeys, "#histories#") {
		if mem != nil {
			inputs := maps.Clone(prompt_inputs)
			inputs["#histories#"] = ""
			parser := promptutils.NewPromptTemplateParser(raw_prompt, transform.WithVariableTmpl)
			for _, k := range parser.VariableKeys {
				if _, ok := inputs[k]; ok {
					prompt_inputs[k] = inputs[k]
				}
			}

			tmp_human_message := modelruntimeentities.NewUserPromptMessage(parser.Format(prompt_inputs, false), "")
			rest_tokens := transform._calculate_rest_token([]modelruntimeentities.PromptMessager{tmp_human_message}, model_config)
			histories := transform._get_history_messages_from_memory(
				mem,
				memory_config,
				rest_tokens,
				role_prefix.User,
				role_prefix.Assistant,
			)
			prompt_inputs["#histories#"] = histories
		} else {
			prompt_inputs["#histories#"] = ""
		}
	}
	return prompt_inputs
}
