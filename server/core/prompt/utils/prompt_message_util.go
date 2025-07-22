package utils

import (
	"fmt"

	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

func PromptMessagesToPromptForSaving(modelMode modelruntimeentities.LLMMode, prompt_messages []modelruntimeentities.PromptMessager) []map[string]any {
	prompts := make([]map[string]any, 0)

	if modelMode == modelruntimeentities.LLMMode_CHAT {
		toolCalls := make([]map[string]any, 0)

		for _, prompt_message := range prompt_messages {
			var role string
			switch prompt_message.Role() {
			case modelruntimeentities.PromptMessageRole_USER:
				role = "user"
			case modelruntimeentities.PromptMessageRole_ASSISTANT:
				role = "assistant"
				if assistantMsg, ok := any(prompt_message).(*modelruntimeentities.AssistantPromptMessage); ok {
					for _, toolCall := range assistantMsg.ToolCalls {
						toolCalls = append(toolCalls, map[string]any{
							"id":   toolCall.ID,
							"type": "function",
							"function": map[string]any{
								"name":      toolCall.Function.Name,
								"arguments": toolCall.Function.Arguments,
							},
						})
					}
				}
			case modelruntimeentities.PromptMessageRole_SYSTEM:
				role = "system"
			case modelruntimeentities.PromptMessageRole_TOOL:
				role = "tool"
			default:
				continue
			}

			text := ""
			files := make([]map[string]any, 0)

			switch content := prompt_message.GetContent().(type) {
			case []modelruntimeentities.PromptMessageContenter:
				for _, item := range content {
					switch itemContent := any(item).(type) {
					case *modelruntimeentities.TextPromptMessageContent:
						text += itemContent.Data()
					}
				}
			default:
				text = fmt.Sprintf("%v", content)
			}

			prompt := map[string]any{
				"role":  role,
				"text":  text,
				"files": files,
			}

			if len(toolCalls) > 0 {
				prompt["tool_calls"] = toolCalls
			}

			prompts = append(prompts, prompt)
		}
	} else {
		if len(prompt_messages) == 0 {
			return prompts
		}

		prompt_message := prompt_messages[0]
		text := ""
		files := make([]map[string]any, 0)

		switch content := prompt_message.GetContent().(type) {
		case []modelruntimeentities.PromptMessageContenter:
			for _, item := range content {
				switch itemContent := item.(type) {
				case *modelruntimeentities.TextPromptMessageContent:
					text += itemContent.Data()
				}
			}
		default:
			text = fmt.Sprintf("%v", content)
		}

		params := map[string]any{
			"role": "user",
			"text": text,
		}

		if len(files) > 0 {
			params["files"] = files
		}

		prompts = append(prompts, params)
	}

	return prompts
}
