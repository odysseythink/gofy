package anthropic

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

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/model_runtime/model_provides/base"
	"mlib.com/gofy/server/core/model_runtime/model_provides/oaicompat"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

const anthropicEndpoint = "https://api.anthropic.com/v1/messages"

type AnthropicLLM struct {
	*base.LargeLanguageModel
}

func NewAnthropicLLM() *AnthropicLLM {
	return &AnthropicLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *AnthropicLLM) ProviderName() string { return "anthropic" }
func (l *AnthropicLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *AnthropicLLM) ValidateCredentials(model string, credentials map[string]any) {
	result := l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10}, nil, nil, "")
	if result == nil {
		panic(fmt.Errorf("credentials validation failed: no result"))
	}
}

func (l *AnthropicLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
	return oaicompat.DefaultCustomizableModelSchema(model)
}

func (l *AnthropicLLM) GetNumTokens(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, tools []*modelruntimeentities.PromptMessageTool) int {
	return oaicompat.EstimateTokens(promptMessages)
}

// convertMessages separates system message and converts the rest to Anthropic format.
func convertMessages(messages []modelruntimeentities.PromptMessager) (string, []map[string]any) {
	systemMsg := ""
	result := make([]map[string]any, 0, len(messages))
	for _, msg := range messages {
		role := string(msg.Role())
		if role == "system" {
			if content, ok := msg.GetContent().(string); ok {
				systemMsg = content
			}
			continue
		}
		m := map[string]any{"role": role}
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
		result = append(result, m)
	}
	return systemMsg, result
}

func (l *AnthropicLLM) Invoke(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	systemMsg, messages := convertMessages(promptMessages)

	maxTokens := 4096
	if mt, ok := modelParameters["max_tokens"]; ok {
		switch v := mt.(type) {
		case int:
			maxTokens = v
		case float64:
			maxTokens = int(v)
		}
		delete(modelParameters, "max_tokens")
	}

	reqBody := map[string]any{
		"model":      model,
		"messages":   messages,
		"max_tokens": maxTokens,
		"stream":     false,
	}
	if systemMsg != "" {
		reqBody["system"] = systemMsg
	}
	if len(stop) > 0 {
		reqBody["stop_sequences"] = stop
	}
	for k, v := range modelParameters {
		reqBody[k] = v
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", anthropicEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		panic(fmt.Errorf("failed to create request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

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
		ID      string `json:"id"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(fmt.Errorf("failed to parse response: %w", err))
	}

	text := ""
	for _, block := range result.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}

	return &modelruntimeentities.LLMResult{
		Model:          model,
		PromptMessages: promptMessages,
		Message:        modelruntimeentities.NewAssistantPromptMessage(text, "", nil),
		Usage: &modelruntimeentities.LLMUsage{
			PromptTokens:     result.Usage.InputTokens,
			CompletionTokens: result.Usage.OutputTokens,
			TotalTokens:      result.Usage.InputTokens + result.Usage.OutputTokens,
		},
	}
}

func (l *AnthropicLLM) InvokeStream(model string, credentials map[string]any, promptMessages []modelruntimeentities.PromptMessager, modelParameters map[string]any, tools []*modelruntimeentities.PromptMessageTool, stop []string, user string) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		apiKey, _ := credentials["api_key"].(string)
		systemMsg, messages := convertMessages(promptMessages)

		maxTokens := 4096
		if mt, ok := modelParameters["max_tokens"]; ok {
			switch v := mt.(type) {
			case int:
				maxTokens = v
			case float64:
				maxTokens = int(v)
			}
			delete(modelParameters, "max_tokens")
		}

		reqBody := map[string]any{
			"model":      model,
			"messages":   messages,
			"max_tokens": maxTokens,
			"stream":     true,
		}
		if systemMsg != "" {
			reqBody["system"] = systemMsg
		}
		if len(stop) > 0 {
			reqBody["stop_sequences"] = stop
		}
		for k, v := range modelParameters {
			reqBody[k] = v
		}

		bodyBytes, _ := json.Marshal(reqBody)
		req, err := http.NewRequest("POST", anthropicEndpoint, bytes.NewReader(bodyBytes))
		if err != nil {
			mlog.Errorf("failed to create request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")

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
		var inputTokens, outputTokens int

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")

			var event map[string]any
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue
			}

			eventType, _ := event["type"].(string)

			switch eventType {
			case "message_start":
				if msg, ok := event["message"].(map[string]any); ok {
					if usage, ok := msg["usage"].(map[string]any); ok {
						if v, ok := usage["input_tokens"].(float64); ok {
							inputTokens = int(v)
						}
					}
				}
			case "content_block_delta":
				delta, ok := event["delta"].(map[string]any)
				if !ok {
					continue
				}
				text, _ := delta["text"].(string)

				result := &modelruntimeentities.LLMResultChunk{
					Model:          model,
					PromptMessages: promptMessages,
					Delta: &modelruntimeentities.LLMResultChunkDelta{
						Index:   index,
						Message: modelruntimeentities.NewAssistantPromptMessage(text, "", nil),
					},
				}
				if !yield(result) {
					return
				}
				index++
			case "message_delta":
				if delta, ok := event["delta"].(map[string]any); ok {
					finishReason, _ := delta["stop_reason"].(string)
					if usage, ok := event["usage"].(map[string]any); ok {
						if v, ok := usage["output_tokens"].(float64); ok {
							outputTokens = int(v)
						}
					}
					result := &modelruntimeentities.LLMResultChunk{
						Model:          model,
						PromptMessages: promptMessages,
						Delta: &modelruntimeentities.LLMResultChunkDelta{
							Index:        index,
							Message:      modelruntimeentities.NewAssistantPromptMessage("", "", nil),
							FinishReason: finishReason,
							Usage: &modelruntimeentities.LLMUsage{
								PromptTokens:     inputTokens,
								CompletionTokens: outputTokens,
								TotalTokens:      inputTokens + outputTokens,
							},
						},
					}
					if !yield(result) {
						return
					}
				}
			case "message_stop":
				return
			}
		}
	}
}
