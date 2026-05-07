package openai_api_compatible

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

	"github.com/odysseythink/gofy/backend/core/model_runtime/model_provides/base"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
	"github.com/odysseythink/mlog"
)

type OpenAICompatibleLLM struct {
	*base.LargeLanguageModel
}

func NewOpenAICompatibleLLM() *OpenAICompatibleLLM {
	return &OpenAICompatibleLLM{
		LargeLanguageModel: &base.LargeLanguageModel{
			BaseAIModel: &base.BaseAIModel{
				ModeType: modelruntimeenumtypes.Model_LLM,
			},
		},
	}
}

func (l *OpenAICompatibleLLM) ProviderName() string {
	return "openai_api_compatible"
}

func (l *OpenAICompatibleLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

func (l *OpenAICompatibleLLM) ValidateCredentials(model string, credentials map[string]any) {
	// Invoke panics on failure; if it returns, credentials are valid.
	_ = l.Invoke(model, credentials,
		[]modelruntimeentities.PromptMessager{modelruntimeentities.NewUserPromptMessage("ping", "")},
		map[string]any{"max_tokens": 10},
		nil, nil, "",
	)
}

func (l *OpenAICompatibleLLM) GetCustomizableModelSchema(model string, credentials map[string]any) *modelruntimeentities.AIModelEntity {
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

func (l *OpenAICompatibleLLM) GetNumTokens(
	model string,
	credentials map[string]any,
	promptMessages []modelruntimeentities.PromptMessager,
	tools []*modelruntimeentities.PromptMessageTool,
) int {
	// Simple token estimation for OpenAI-compatible models
	total := 0
	for _, msg := range promptMessages {
		if content, ok := msg.GetContent().(string); ok {
			total += len([]rune(content)) / 2
		}
	}
	return total
}

// InvokeStream calls an OpenAI-compatible API with streaming.
func (l *OpenAICompatibleLLM) InvokeStream(
	model string,
	credentials map[string]any,
	promptMessages []modelruntimeentities.PromptMessager,
	modelParameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
) iter.Seq[*modelruntimeentities.LLMResultChunk] {
	return func(yield func(*modelruntimeentities.LLMResultChunk) bool) {
		apiKey, _ := credentials["api_key"].(string)
		endpoint, _ := credentials["endpoint_url"].(string)
		if !strings.HasSuffix(endpoint, "/") {
			endpoint += "/"
		}
		url := endpoint + "chat/completions"

		messages := convertPromptMessages(promptMessages)

		reqBody := map[string]any{
			"model":    model,
			"messages": messages,
			"stream":   true,
		}
		if len(stop) > 0 {
			reqBody["stop"] = stop
		}
		if len(tools) > 0 {
			toolDefs := convertTools(tools)
			reqBody["tools"] = toolDefs
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
			if delta.ReasoningContent != "" {
				// Wrap reasoning content using the base LLM helper
				content = delta.ReasoningContent
			}

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

// Invoke calls the API without streaming.
func (l *OpenAICompatibleLLM) Invoke(
	model string,
	credentials map[string]any,
	promptMessages []modelruntimeentities.PromptMessager,
	modelParameters map[string]any,
	tools []*modelruntimeentities.PromptMessageTool,
	stop []string,
	user string,
) *modelruntimeentities.LLMResult {
	apiKey, _ := credentials["api_key"].(string)
	endpoint, _ := credentials["endpoint_url"].(string)
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	url := endpoint + "chat/completions"

	messages := convertPromptMessages(promptMessages)

	reqBody := map[string]any{
		"model":    model,
		"messages": messages,
		"stream":   false,
	}
	if len(stop) > 0 {
		reqBody["stop"] = stop
	}
	if len(tools) > 0 {
		toolDefs := convertTools(tools)
		reqBody["tools"] = toolDefs
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

// convertPromptMessages converts PromptMessagers to OpenAI message format.
func convertPromptMessages(messages []modelruntimeentities.PromptMessager) []map[string]any {
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
		// Include tool_calls for assistant messages
		if assistant, ok := msg.(*modelruntimeentities.AssistantPromptMessage); ok && len(assistant.ToolCalls) > 0 {
			toolCalls := make([]map[string]any, 0, len(assistant.ToolCalls))
			for _, tc := range assistant.ToolCalls {
				toolCalls = append(toolCalls, tc.ModelDump())
			}
			m["tool_calls"] = toolCalls
		}
		// Include tool_call_id for tool messages
		if toolMsg, ok := msg.(*modelruntimeentities.ToolPromptMessage[string]); ok {
			m["tool_call_id"] = toolMsg.ToolCallID
		}
		result = append(result, m)
	}
	return result
}

// convertTools converts PromptMessageTool to OpenAI tools format.
func convertTools(tools []*modelruntimeentities.PromptMessageTool) []map[string]any {
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
