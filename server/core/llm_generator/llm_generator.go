package llmgenerator

import (
	"encoding/json"
	"regexp"
	"strings"

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
