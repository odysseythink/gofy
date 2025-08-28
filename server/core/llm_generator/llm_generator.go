package llmgenerator

import (
	"encoding/json"
	"regexp"
	"strings"
modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeexceptions "mlib.com/gofy/server/core/exceptions/model_runtime"
	"mlib.com/gofy/server/core/manageres"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/mlog"
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

    func (g *LLMGenerator) GenerateRuleConfig(tenant_id  string, instruction  string, model_config map[string]any, no_variable bool)  map[string]any{
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

        if no_variable{
            prompt_template := promptutils.NewPromptTemplateParser(WORKFLOW_RULE_CONFIG_PROMPT_GENERATE_TEMPLATE, false)

            prompt_generate := prompt_template.Format(
                map[string]any{
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
		response := model_instance.InvokeLLM(prompts, map[string]any{"max_tokens": 100, "temperature": 1}, nil, nil, "", nil)


            try{
                response = cast(
                    LLMResult,
                    model_instance.invoke_llm(
                        prompt_messages=list(prompt_messages), model_parameters=model_parameters, stream=False
                    ),
                )

                rule_config["prompt"] = cast(str, response.message.content)

            except InvokeError as e{
                error = str(e)
                error_step = "generate rule config"
            except Exception as e{
                logging.exception("Failed to generate rule config, model: %s", model_config.get("name"))
                rule_config["error"] = str(e)

            rule_config["error"] = f"Failed to {error_step}. Error: {error}" if error else ""

            return rule_config
		}
        // get rule config prompt, parameter and statement
        prompt_generate, parameter_generate, statement_generate = output_parser.get_format_instructions()

        prompt_template = PromptTemplateParser(prompt_generate)

        parameter_template = PromptTemplateParser(parameter_generate)

        statement_template = PromptTemplateParser(statement_generate)

        // format the prompt_generate_prompt
        prompt_generate_prompt = prompt_template.format(
            inputs={
                "TASK_DESCRIPTION": instruction,
            },
            remove_template_variables=False,
        )
        prompt_messages = [UserPromptMessage(content=prompt_generate_prompt)]

        // get model instance
        model_manager = ModelManager()
        model_instance = model_manager.get_model_instance(
            tenant_id=tenant_id,
            model_type=ModelType.LLM,
            provider=model_config.get("provider", ""),
            model=model_config.get("name", ""),
        )

        try{
            try{
                // the first step to generate the task prompt
                prompt_content = cast(
                    LLMResult,
                    model_instance.invoke_llm(
                        prompt_messages=list(prompt_messages), model_parameters=model_parameters, stream=False
                    ),
                )
            except InvokeError as e{
                error = str(e)
                error_step = "generate prefix prompt"
                rule_config["error"] = f"Failed to {error_step}. Error: {error}" if error else ""

                return rule_config

            rule_config["prompt"] = cast(str, prompt_content.message.content)

            if not isinstance(prompt_content.message.content, str){
                raise NotImplementedError("prompt content is not a string")
            parameter_generate_prompt = parameter_template.format(
                inputs={
                    "INPUT_TEXT": prompt_content.message.content,
                },
                remove_template_variables=False,
            )
            parameter_messages = [UserPromptMessage(content=parameter_generate_prompt)]

            // the second step to generate the task_parameter and task_statement
            statement_generate_prompt = statement_template.format(
                inputs={
                    "TASK_DESCRIPTION": instruction,
                    "INPUT_TEXT": prompt_content.message.content,
                },
                remove_template_variables=False,
            )
            statement_messages = [UserPromptMessage(content=statement_generate_prompt)]

            try{
                parameter_content = cast(
                    LLMResult,
                    model_instance.invoke_llm(
                        prompt_messages=list(parameter_messages), model_parameters=model_parameters, stream=False
                    ),
                )
                rule_config["variables"] = re.findall(r'"\s*([^"]+)\s*"', cast(str, parameter_content.message.content))
            except InvokeError as e{
                error = str(e)
                error_step = "generate variables"

            try{
                statement_content = cast(
                    LLMResult,
                    model_instance.invoke_llm(
                        prompt_messages=list(statement_messages), model_parameters=model_parameters, stream=False
                    ),
                )
                rule_config["opening_statement"] = cast(str, statement_content.message.content)
            except InvokeError as e{
                error = str(e)
                error_step = "generate conversation opener"

        except Exception as e{
            logging.exception("Failed to generate rule config, model: %s", model_config.get("name"))
            rule_config["error"] = str(e)

        rule_config["error"] = f"Failed to {error_step}. Error: {error}" if error else ""

        return rule_config
}