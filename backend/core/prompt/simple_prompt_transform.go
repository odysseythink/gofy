package prompt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	"mlib.com/gofy/server/core/memory"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	promptentities "mlib.com/gofy/server/entities/prompt"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
)

var (
	prompt_file_contents = map[string]map[string]any{}
)

type SimplePromptTransform struct {
	*PromptTransform
}

func (transform *SimplePromptTransform) _get_last_user_message(
	prompt string,
	// files  []*file.File,
) *modelruntimeentities.UserPromptMessage[string] {
	// if files{
	//     prompt_message_contents: list[PromptMessageContentUnionTypes] = []
	//     prompt_message_contents.append(TextPromptMessageContent(data=prompt))
	//     for file in files{
	//         prompt_message_contents.append(
	//             file_manager.to_prompt_message_content(file, image_detail_config=image_detail_config)
	//         )
	// 	}
	//     prompt_message = UserPromptMessage(content=prompt_message_contents)
	// } else {
	prompt_message := modelruntimeentities.NewUserPromptMessage(prompt, "")
	// }
	return prompt_message
}

func (transform *SimplePromptTransform) _prompt_file_name(app_mode models.AppMode, provider string, model string) string {
	// baichuan
	is_baichuan := false
	if provider == "baichuan" {
		is_baichuan = true
	} else {
		if slices.Contains([]string{"huggingface_hub", "openllm", "xinference"}, provider) && strings.Contains(strings.ToLower(model), "baichuan") {
			is_baichuan = true
		}
	}
	if is_baichuan {
		if app_mode == models.AppMode_COMPLETION {
			return "baichuan_completion"
		} else {
			return "baichuan_chat"
		}
	}
	// common
	if app_mode == models.AppMode_COMPLETION {
		return "common_completion"
	} else {
		return "common_chat"
	}
}

func (transform *SimplePromptTransform) _get_prompt_rule(app_mode models.AppMode, provider string, model string) map[string]any {
	prompt_file_name := transform._prompt_file_name(app_mode, provider, model)

	// Check if the prompt file is already loaded
	if _, ok := prompt_file_contents[prompt_file_name]; ok {
		return prompt_file_contents[prompt_file_name]
	}
	dir, _ := os.Getwd()
	// Get the absolute path of the subdirectory
	prompt_path := filepath.Join(dir, "prompt", "prompt_templates")
	json_file_path := filepath.Join(prompt_path, prompt_file_name+".json")

	// Open the JSON file and read its content
	bindata, err := os.ReadFile(json_file_path)
	if err != nil {
		mlog.Errorf("read file=%s failed:%v", json_file_path, err)
		panic(exceptions.NewValueError(fmt.Sprintf("read file=%s failed:%v", json_file_path, err)))
	}
	var content map[string]any
	err = json.Unmarshal(bindata, &content)
	if err != nil {
		mlog.Errorf("json unmarshal failed:%v", err)
		panic(exceptions.NewValueError(fmt.Sprintf("json unmarshal failed:%v", err)))
	}

	// Store the content of the prompt file
	prompt_file_contents[prompt_file_name] = content

	return content
}
func (transform *SimplePromptTransform) GetPromptTemplate(
	app_mode models.AppMode,
	provider string,
	model string,
	pre_prompt string,
	has_context bool,
	query_in_prompt bool,
	with_memory_prompt bool,
) map[string]any {
	prompt_rules := transform._get_prompt_rule(app_mode, provider, model)

	custom_variable_keys := []string{}
	special_variable_keys := []string{}

	prompt := ""
	for _, order := range mapstruct.Get(prompt_rules, "system_prompt_orders", []string{}) {
		if order == "context_prompt" && has_context {
			prompt += mapstruct.Get(prompt_rules, "context_prompt", "")
			special_variable_keys = append(special_variable_keys, "#context#")
		} else if order == "pre_prompt" && pre_prompt != "" {
			prompt += pre_prompt + "\n"
			pre_prompt_template := promptutils.NewPromptTemplateParser(pre_prompt, false)
			custom_variable_keys = pre_prompt_template.VariableKeys
		} else if order == "histories_prompt" && with_memory_prompt {
			prompt += mapstruct.Get(prompt_rules, "histories_prompt", "")
			special_variable_keys = append(special_variable_keys, "#histories#")
		}
	}
	if query_in_prompt {
		prompt += mapstruct.Get(prompt_rules, "query_prompt", "{{#query#}}")
		special_variable_keys = append(special_variable_keys, "#query#")
	}
	return map[string]any{
		"prompt_template":       promptutils.NewPromptTemplateParser(prompt, false),
		"custom_variable_keys":  custom_variable_keys,
		"special_variable_keys": special_variable_keys,
		"prompt_rules":          prompt_rules,
	}
}
func (transform *SimplePromptTransform) _get_prompt_str_and_rules(
	app_mode models.AppMode,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
	pre_prompt string,
	inputs map[string]string,
	query string,
	context string,
	histories string,
) (string, map[string]any) {
	// get prompt template
	prompt_template_config := transform.GetPromptTemplate(
		app_mode,
		model_config.Provider,
		model_config.Model,
		pre_prompt,
		context != "",
		query != "",
		histories != "",
	)
	variables := map[string]string{}
	for _, k := range prompt_template_config["custom_variable_keys"].([]string) {
		if _, ok := inputs[k]; ok {
			variables[k] = inputs[k]
		}
	}

	for _, v := range prompt_template_config["special_variable_keys"].([]string) {
		// support #context#, #query# and #histories#
		if v == "#context#" {
			variables["#context#"] = context
		} else if v == "#query#" {
			variables["#query#"] = query
		} else if v == "#histories#" {
			variables["#histories#"] = histories
		}
	}
	prompt_template := prompt_template_config["prompt_template"].(*promptutils.PromptTemplateParser)
	prompt := prompt_template.Format(variables, false)

	return prompt, prompt_template_config["prompt_rules"].(map[string]any)
}

func (transform *SimplePromptTransform) _get_completion_model_prompt_messages(
	app_mode models.AppMode,
	pre_prompt string,
	inputs map[string]string,
	query string,
	context string,
	files []*file.File,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) ([]modelruntimeentities.PromptMessager, []string) {
	// get prompt
	prompt, prompt_rules := transform._get_prompt_str_and_rules(
		app_mode,
		model_config,
		pre_prompt,
		inputs,
		query,
		context,
		"",
	)

	if mem != nil {
		tmp_human_message := modelruntimeentities.NewUserPromptMessage(prompt, "")

		rest_tokens := transform._calculate_rest_token([]modelruntimeentities.PromptMessager{tmp_human_message}, model_config)
		histories := transform._get_history_messages_from_memory(
			mem,
			&promptentities.MemoryConfig{
				Window: promptentities.WindowConfig{
					Enabled: false,
				},
			},
			rest_tokens,
			mapstruct.Get(prompt_rules, "human_prefix", "Human"),
			mapstruct.Get(prompt_rules, "assistant_prefix", "Assistant"),
		)

		// get prompt
		prompt, prompt_rules = transform._get_prompt_str_and_rules(
			app_mode,
			model_config,
			pre_prompt,
			inputs,
			query,
			context,
			histories,
		)
	}
	stops := mapstruct.Get(prompt_rules, "stops", []string{})
	if len(stops) == 0 {
		stops = nil
	}
	return []modelruntimeentities.PromptMessager{transform._get_last_user_message(prompt)}, stops
}
func (transform *SimplePromptTransform) _get_chat_model_prompt_messages(
	app_mode models.AppMode,
	pre_prompt string,
	inputs map[string]string,
	query string,
	context string,
	files []*file.File,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) ([]modelruntimeentities.PromptMessager, []string) {
	prompt_messages := []modelruntimeentities.PromptMessager{}

	// get prompt
	prompt, _ := transform._get_prompt_str_and_rules(
		app_mode,
		model_config,
		pre_prompt,
		inputs,
		"",
		context,
		"",
	)

	if prompt != "" && query != "" {
		prompt_messages = append(prompt_messages, modelruntimeentities.NewSystemPromptMessage(prompt, ""))
	}
	if mem != nil {
		prompt_messages = transform._append_chat_histories(
			mem,
			&promptentities.MemoryConfig{
				Window: promptentities.WindowConfig{
					Enabled: false,
				},
			},
			prompt_messages,
			model_config,
		)
	}
	if query != "" {
		prompt_messages = append(prompt_messages, transform._get_last_user_message(query))
	} else {
		prompt_messages = append(prompt_messages, transform._get_last_user_message(prompt))
	}
	return prompt_messages, nil
}
func (transform *SimplePromptTransform) GetPrompt(
	app_mode models.AppMode,
	prompt_template_entity *appconfigentities.PromptTemplateEntity,
	inputs map[string]string,
	query string,
	files []*file.File,
	context string,
	mem *memory.TokenBufferMemory,
	model_config *appconfigentities.ModelConfigWithCredentialsEntity,
) ([]modelruntimeentities.PromptMessager, []string) {

	model_mode := model_config.Mode
	var prompt_messages []modelruntimeentities.PromptMessager
	var stops []string
	if model_mode == modelruntimeentities.LLMMode_CHAT {
		prompt_messages, stops = transform._get_chat_model_prompt_messages(
			app_mode,
			prompt_template_entity.SimplePromptTemplate,
			inputs,
			query,
			context,
			files,
			mem,
			model_config,
		)
	} else {
		prompt_messages, stops = transform._get_completion_model_prompt_messages(
			app_mode,
			prompt_template_entity.SimplePromptTemplate,
			inputs,
			query,
			context,
			files,
			mem,
			model_config,
		)
	}
	return prompt_messages, stops
}
