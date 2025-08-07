package builtintool

import (
	"strings"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/tools/base"
	modelinvocationutils "mlib.com/gofy/server/core/tools/utils/model_invocation_utils"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
)

const (
	_SUMMARY_PROMPT = `You are a professional language researcher, you are interested in the language
and you can quickly aimed at the main point of an webpage and reproduce it in your own words but 
retain the original meaning and keep the key points. 
however, the text you got is too long, what you got is possible a part of the text.
Please summarize the text you got.`
)

type BuiltinTool struct {
	*base.Tool
	Provider string `json:"provider"`
}

func (t *BuiltinTool) ToolProviderType() toolsenumtypes.ToolProviderType {
	return toolsenumtypes.ToolProvider_BUILT_IN
}

//	func (t *BuiltinTool) ForkToolRuntime(runtime *base.ToolRuntime) *base.Tooler {
//		return &BuiltinTool{
//			Tool: &base.Tool{
//				Entity:  toolsentities.NewToolEntity(t.Tool.Entity),
//				Runtime: base.NewToolRuntime(t.Tool.Runtime),
//			},
//			Provider: t.Provider,
//		}
//	}
func (t *BuiltinTool) InvokeModel(user_id string, prompt_messages []modelruntimeentities.PromptMessager, stop []string) *modelruntimeentities.LLMResult {
	// invoke model
	if t.Tool == nil || t.Tool.Runtime == nil || t.Tool.Entity == nil {
		panic(exceptions.NewValueError("runtime and Entity are required"))
	}
	return modelinvocationutils.Invoke(
		user_id,
		t.Tool.Runtime.TenantID,
		"builtin",
		t.Entity.Identity.Name,
		prompt_messages,
	)
}
func (t *BuiltinTool) GetMaxTokens() int {
	if t.Tool.Runtime == nil {
		panic(exceptions.NewValueError("runtime is required"))
	}
	return modelinvocationutils.GetMaxLLMContextTokens(t.Tool.Runtime.TenantID)
}
func (t *BuiltinTool) GetPromptTokens(prompt_messages []modelruntimeentities.PromptMessager) int {
	if t.Tool.Runtime == nil {
		panic(exceptions.NewValueError("runtime is required"))
	}
	return modelinvocationutils.CalculateTokens(t.Tool.Runtime.TenantID, prompt_messages)
}

func (t *BuiltinTool) _summary_get_prompt_tokens(content string) int {
	return t.GetPromptTokens(
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewSystemPromptMessage(_SUMMARY_PROMPT, ""), modelruntimeentities.NewUserPromptMessage(content, "")},
	)
}
func (t *BuiltinTool) _summary_summarize(user_id, content string) string {
	summary := t.InvokeModel(
		user_id,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewSystemPromptMessage(_SUMMARY_PROMPT, ""), modelruntimeentities.NewUserPromptMessage(content, "")},
		nil,
	)
	return summary.Message.Content
}
func (t *BuiltinTool) Summary(user_id string, content string) string {
	max_tokens := t.GetMaxTokens()

	if float64(t.GetPromptTokens([]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(content, "")})) < float64(max_tokens)*0.6 {
		return content
	}

	lines := strings.Split(content, "\n")
	new_lines := []string{}
	// split long line into multiple lines
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if float64(len(line)) < float64(max_tokens)*0.5 {
			new_lines = append(new_lines, line)
		} else if float64(t._summary_get_prompt_tokens(line)) > float64(max_tokens)*0.7 {
			for float64(t._summary_get_prompt_tokens(line)) > float64(max_tokens)*0.7 {
				new_lines = append(new_lines, line[:int(float64(max_tokens)*0.5)])
				line = line[int(float64(max_tokens)*0.5):]
			}
			new_lines = append(new_lines, line)
		} else {
			new_lines = append(new_lines, line)
		}
	}
	// merge lines into messages with max tokens
	messages := []string{}
	for _, j := range new_lines {
		if len(messages) == 0 {
			messages = append(messages, j)
		} else {
			if float64(len(messages[len(messages)-1])+len(j)) < float64(max_tokens)*0.5 {
				messages[len(messages)-1] += j
			}
			if float64(t._summary_get_prompt_tokens(messages[len(messages)-1]+j)) > float64(max_tokens)*0.7 {
				messages = append(messages, j)
			} else {
				messages[len(messages)-1] += j
			}
		}
	}
	summaries := []string{}
	for _, message := range messages {
		summary := t._summary_summarize(user_id, message)
		summaries = append(summaries, summary)
	}
	result := strings.Join(summaries, "\n")
	if float64(t.GetPromptTokens([]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage(result, "")})) > float64(max_tokens)*0.7 {
		return t.Summary(user_id, result)
	}
	return result
}
