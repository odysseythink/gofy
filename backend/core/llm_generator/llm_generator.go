package llmgenerator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/odysseythink/mlog"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	"mlib.com/gofy/server/core/manageres"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type LLMGenerator struct{}

func (g *LLMGenerator) GenerateConversationName(
	tenant_id, query, conversation_id, app_id string,
) string {
	prompt := CONVERSATION_TITLE_PROMPT

	if len(query) > 2000 {
		query = query[:300] + "...[TRUNCATED]..." + query[len(query)-300:]
	}
	query = strings.ReplaceAll(query, "\n", " ")

	prompt += query + "\n"

	model_instance := manageres.Instance.Model.GetDefaultModelInstance(
		tenant_id,
		modelruntimeenumtypes.Model_LLM,
	)

	prompts := []modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(prompt, "")}
	response := model_instance.InvokeLLM(prompts, map[string]any{"max_tokens": 100, "temperature": 1}, nil, nil, "", nil)

	answer := response.Message.Content
	re := regexp.MustCompile(`^.*(\{.*\}).*$`)
	cleaned_answer := re.ReplaceAllString(answer, "$1")
	if cleaned_answer == "" {
		return ""
	}
	var result_dict map[string]any
	err := json.Unmarshal([]byte(cleaned_answer), &result_dict)
	if err != nil {
		mlog.Errorf("json.Unmarshal(%s) failed:%v", cleaned_answer, err)
		return ""
	}
	answer = ""
	if _, ok := result_dict["Your Output"]; ok {
		if _, ok := result_dict["Your Output"].(string); ok {
			answer = result_dict["Your Output"].(string)
		}
	}
	name := strings.TrimSpace(answer)

	if len(name) > 75 {
		name = name[:75] + "..."
	}

	return name
}
func (g *LLMGenerator) GenerateSuggestedQuestionsAfterAnswer(tenant_id, histories string) []string {
	output_parser := &SuggestedQuestionsAfterAnswerOutputParser{}
	format_instructions := output_parser.GetFormatInstructions()

	prompt_template := promptutils.NewPromptTemplateParser("{{histories}}\n{{format_instructions}}\nquestions:\n", false)

	prompt := prompt_template.Format(map[string]string{"histories": histories, "format_instructions": format_instructions}, true)
	var questions []string
	model_instance := func() *modelmanager.ModelInstance {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(*modelruntimeexceptions.InvokeAuthorizationError); ok {
					questions = []string{}
				} else {
					panic(r)
				}
			}
		}()
		return manageres.Instance.Model.GetDefaultModelInstance(tenant_id, modelruntimeenumtypes.Model_LLM)
	}()
	if questions != nil {
		return questions
	}

	prompt_messages := []modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(prompt, "")}
	questions = func() (ret []string) {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(*modelruntimeexceptions.InvokeError); ok {
					ret = []string{}
				} else {
					mlog.Error("Failed to generate suggested questions after answer")
					ret = []string{}
				}
			}
		}()
		response := model_instance.InvokeLLM(prompt_messages, map[string]any{"max_tokens": 256, "temperature": 0}, nil, nil, "", nil)

		ret = output_parser.Parse(response.Message.Content)
		return
	}()

	return questions
}

func (g *LLMGenerator) GenerateRuleConfig(tenant_id string, instruction string, model_config map[string]any, no_variable bool) map[string]any {
	output_parser := &RuleConfigGeneratorOutputParser{}

	errormsg := ""
	error_step := ""
	rule_config := map[string]any{"prompt": "", "variables": []any{}, "opening_statement": "", "error": ""}
	model_parameters := map[string]any{}
	if _, ok := model_config["completion_params"]; ok {
		if _, ok := model_config["completion_params"].(map[string]any); ok {
			model_parameters = model_config["completion_params"].(map[string]any)
		}
	}

	if no_variable {
		prompt_template := promptutils.NewPromptTemplateParser(WORKFLOW_RULE_CONFIG_PROMPT_GENERATE_TEMPLATE, false)

		prompt_generate := prompt_template.Format(
			map[string]string{
				"TASK_DESCRIPTION": instruction,
			},
			false,
		)
		provider := ""
		if _, ok := model_config["provider"]; ok {
			if _, ok := model_config["provider"].(string); ok {
				provider = model_config["provider"].(string)
			}
		}
		name := ""
		if _, ok := model_config["name"]; ok {
			if _, ok := model_config["name"].(string); ok {
				name = model_config["name"].(string)
			}
		}
		prompt_messages := []modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(prompt_generate, "")}
		model_instance := manageres.Instance.Model.GetModelInstance(
			tenant_id,
			provider,
			modelruntimeenumtypes.Model_LLM,
			name,
		)
		func() {
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(*modelruntimeexceptions.InvokeError); ok {
						errormsg = e.Error()
						error_step = "generate rule config"
					} else {
						mlog.Errorf("Failed to generate rule config, model: %s", name)
						rule_config["error"] = fmt.Sprintf("%v", r)
					}
				}
			}()
			response := model_instance.InvokeLLM(prompt_messages, model_parameters, nil, nil, "", nil)

			rule_config["prompt"] = response.Message.Content
		}()

		if errormsg != "" {
			rule_config["error"] = fmt.Sprintf("Failed to {%v}. Error: {%s}", error_step, errormsg)
		} else {
			rule_config["error"] = ""
		}

		return rule_config
	}
	// get rule config prompt, parameter and statement
	prompt_generate, parameter_generate, statement_generate := output_parser.GetFormatInstructions()

	prompt_template := promptutils.NewPromptTemplateParser(prompt_generate, false)

	parameter_template := promptutils.NewPromptTemplateParser(parameter_generate, false)

	statement_template := promptutils.NewPromptTemplateParser(statement_generate, false)

	// format the prompt_generate_prompt
	prompt_generate_prompt := prompt_template.Format(
		map[string]string{
			"TASK_DESCRIPTION": instruction,
		},
		false,
	)
	prompt_messages := []modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(prompt_generate_prompt, "")}
	provider := ""
	if _, ok := model_config["provider"]; ok {
		if _, ok := model_config["provider"].(string); ok {
			provider = model_config["provider"].(string)
		}
	}
	name := ""
	if _, ok := model_config["name"]; ok {
		if _, ok := model_config["name"].(string); ok {
			name = model_config["name"].(string)
		}
	}
	// get model instance
	model_instance := manageres.Instance.Model.GetModelInstance(
		tenant_id,
		provider,
		modelruntimeenumtypes.Model_LLM,
		name,
	)

	direct_return := func() bool {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					mlog.Errorf("Failed to generate rule config, model: %s", name)
					rule_config["error"] = e.Error()
				}
			}
		}()
		prompt_content, direct_return := func() (prompt_content *modelruntimeentities.LLMResult, direct_return bool) {
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(*modelruntimeexceptions.InvokeError); ok {
						errormsg = e.Error()
						error_step = "generate prefix prompt"
						if errormsg != "" {
							rule_config["error"] = fmt.Sprintf("Failed to {%v}. Error: {%s}", error_step, errormsg)
						} else {
							rule_config["error"] = ""
						}
						direct_return = true
					}
				}
			}()
			prompt_content = model_instance.InvokeLLM(prompt_messages, model_parameters, nil, nil, "", nil)
			return
		}()
		if direct_return {
			return direct_return
		}

		rule_config["prompt"] = prompt_content.Message.Content

		parameter_generate_prompt := parameter_template.Format(
			map[string]string{
				"INPUT_TEXT": prompt_content.Message.Content,
			},
			false,
		)
		parameter_messages := []modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(parameter_generate_prompt, "")}

		// the second step to generate the task_parameter and task_statement
		statement_generate_prompt := statement_template.Format(
			map[string]string{
				"TASK_DESCRIPTION": instruction,
				"INPUT_TEXT":       prompt_content.Message.Content,
			},
			false,
		)
		statement_messages := []modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(statement_generate_prompt, "")}

		func() {
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(*modelruntimeexceptions.InvokeError); ok {
						errormsg = e.Error()
						error_step = "generate variables"
					}
				}
			}()
			parameter_content := model_instance.InvokeLLM(parameter_messages, model_parameters, nil, nil, "", nil)
			re := regexp.MustCompile(`"\s*([^"]+)\s*"`)
			matches := re.FindAllStringSubmatch(parameter_content.Message.Content, -1)
			variables := []string{}
			for _, m := range matches {
				if len(m) > 1 {
					variables = append(variables, m[1])
				}
			}
			rule_config["variables"] = variables
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(*modelruntimeexceptions.InvokeError); ok {
						errormsg = e.Error()
						error_step = "generate conversation opener"
					}
				}
			}()
			statement_content := model_instance.InvokeLLM(statement_messages, model_parameters, nil, nil, "", nil)
			rule_config["opening_statement"] = statement_content.Message.Content
		}()
		return false
	}()
	if direct_return {
		return rule_config
	}
	if errormsg != "" {
		rule_config["error"] = fmt.Sprintf("Failed to {%v}. Error: {%s}", error_step, errormsg)
	} else {
		rule_config["error"] = ""
	}

	return rule_config
}
