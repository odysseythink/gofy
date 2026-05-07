// Package oaicompat provides shared helpers for OpenAI-compatible model providers.
package oaicompat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"strings"
	"time"

	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
	"github.com/odysseythink/mlog"
)

// ConvertPromptMessages converts PromptMessagers to OpenAI message format.
func ConvertPromptMessages(messages []modelruntimeentities.PromptMessager) []map[string]any {
	result := make([]map[string]any, 0, len(messages))
	for _, msg := range messages {
		m := map[string]any{
			"role": string(msg.Role()),
		}
		switch content := msg.GetContent().(type) {
		case string:
			m["content"] = content
		case []modelruntimeentities.PromptMessageContenter:
			parts := make([]map[string]any, 0, len(content))
			for _, c := range content {
				parts = append(parts, c.ToDict())
			}
			m["content"] = parts
		}
		if assistant, ok := msg.(*modelruntimeentities.AssistantPromptMessage); ok && len(assistant.ToolCalls) > 0 {
			toolCalls := make([]map[string]any, 0, len(assistant.ToolCalls))
			for _, tc := range assistant.ToolCalls {
				toolCalls = append(toolCalls, tc.ModelDump())
			}
			m["tool_calls"] = toolCalls
		}
		if toolMsg, ok := msg.(*modelruntimeentities.ToolPromptMessage[string]); ok {
			m["tool_call_id"] = toolMsg.ToolCallID
		}
		result = append(result, m)
	}
	return result
}

// ConvertTools converts PromptMessageTool to OpenAI tools format.
func ConvertTools(tools []*modelruntimeentities.PromptMessageTool) []map[string]any {
	result := make([]map[string]any, 0, len(tools))
	for _, tool := range tools {
		result = append(result, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        tool.Name,
				"description": tool.Description,
				"parameters":  tool.Parameters,
			},
		})
	}
	return result
}

// EstimateTokens provides a rough token count estimation.
func EstimateTokens(messages []modelruntimeentities.PromptMessager) int {
	total := 0
	for _, msg := range messages {
		if content, ok := msg.GetContent().(string); ok {
			total += len([]rune(content)) / 2
		}
	}
	return total
}

// DefaultCustomizableModelSchema returns a default AIModelEntity for customizable models.
func DefaultCustomizableModelSchema(model string) *modelruntimeentities.AIModelEntity {
	return &modelruntimeentities.AIModelEntity{
		Model:     model,
		Label:     commontypes.I18nObject{EnUS: model, ZhHans: model},
		ModelType: modelruntimeenumtypes.Model_LLM,
		Features: []modelruntimeenumtypes.ModelFeature{
			modelruntimeenumtypes.ModelFeature_TOOL_CALL,
			modelruntimeenumtypes.ModelFeature_STREAM_TOOL_CALL,
		},
		FetchFrom: modelruntimeenumtypes.FetchFrom_CUSTOMIZABLE_MODEL,
		ModelProperties: map[modelruntimeenumtypes.ModelPropertyKey]any{
			modelruntimeenumtypes.ModelPropertyKey_MODE: modelruntimeentities.LLMMode_CHAT,
		},
	}
}

// Invoke makes a non-streaming chat completion request to an OpenAI-compatible endpoint.
func Invoke(
	endpoint string,
	apiKey string,
	model string,
	promptMessages []modelruntimeentities.PromptMessager,
	modelParameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
) *modelruntimeentities.LLMResult {
	url := strings.TrimRight(endpoint, "/") + "/chat/completions"

	messages := ConvertPromptMessages(promptMessages)

	reqBody := map[string]any{
		"model":    model,
		"messages": messages,
		"stream":   false,
	}
	if len(stop) > 0 {
		reqBody["stop"] = stop
	}
	if len(tools) > 0 {
		reqBody["tools"] = ConvertTools(tools)
	}
	for k, v := range modelParameters {
		reqBody[k] = v
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		panic(fmt.Errorf("failed to create request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 300 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("API error %d: %s", resp.StatusCode, string(body)))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
				Role    string `json:"role"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(fmt.Errorf("failed to parse response: %w", err))
	}

	if len(result.Choices) == 0 {
		panic(fmt.Errorf("no choices in response"))
	}

	assistantMsg := modelruntimeentities.NewAssistantPromptMessage(result.Choices[0].Message.Content, "", nil)

	return &modelruntimeentities.LLMResult{
		Model:          model,
		PromptMessages: promptMessages,
		Message:        assistantMsg,
		Usage: &modelruntimeentities.LLMUsage{
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			TotalTokens:      result.Usage.TotalTokens,
		},
	}
}

// InvokeStream makes a streaming chat completion request to an OpenAI-compatible endpoint.
func InvokeStream(
	endpoint string,
	apiKey string,
	model string,
	promptMessages []modelruntimeentities.PromptMessager,
	modelParameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		url := strings.TrimRight(endpoint, "/") + "/chat/completions"

		messages := ConvertPromptMessages(promptMessages)

		reqBody := map[string]any{
			"model":    model,
			"messages": messages,
			"stream":   true,
		}
		if len(stop) > 0 {
			reqBody["stop"] = stop
		}
		if len(tools) > 0 {
			reqBody["tools"] = ConvertTools(tools)
		}
		for k, v := range modelParameters {
			reqBody[k] = v
		}

		bodyBytes, _ := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			mlog.Errorf("failed to create request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: 300 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			mlog.Errorf("request failed: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			mlog.Errorf("API error %d: %s", resp.StatusCode, string(body))
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		index := 0
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk struct {
				ID      string `json:"id"`
				Choices []struct {
					Delta struct {
						Content          string `json:"content"`
						Role             string `json:"role"`
						ReasoningContent string `json:"reasoning_content,omitempty"`
					} `json:"delta"`
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
				Usage *struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				} `json:"usage"`
			}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			delta := chunk.Choices[0].Delta
			content := delta.Content

			assistantMsg := modelruntimeentities.NewAssistantPromptMessage(content, "", nil)

			result := &modelruntimeentities.LLMResultChunk{
				Model:          model,
				PromptMessages: promptMessages,
				Delta: &modelruntimeentities.LLMResultChunkDelta{
					Index:   index,
					Message: assistantMsg,
				},
			}

			if chunk.Choices[0].FinishReason != nil {
				result.Delta.FinishReason = *chunk.Choices[0].FinishReason
			}

			if chunk.Usage != nil {
				result.Delta.Usage = &modelruntimeentities.LLMUsage{
					PromptTokens:     chunk.Usage.PromptTokens,
					CompletionTokens: chunk.Usage.CompletionTokens,
					TotalTokens:      chunk.Usage.TotalTokens,
				}
			}

			if !yield(result) {
				return
			}
			index++
		}
	}
}
