package plugin

import (
	"encoding/json"

	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
)

// convertRawMessages coerces the "messages" payload (typically []any of maps)
// into []modelruntimeentities.PromptMessager for downstream LLM invocation.
func convertRawMessages(raw any) []modelruntimeentities.PromptMessager {
	var msgs []modelruntimeentities.PromptMessager
	switch v := raw.(type) {
	case []any:
		for _, item := range v {
			msg := modelruntimeentities.NewPromptMessager(item)
			if msg != nil {
				msgs = append(msgs, msg)
			}
		}
	case []map[string]any:
		for _, item := range v {
			msg := modelruntimeentities.NewPromptMessager(item)
			if msg != nil {
				msgs = append(msgs, msg)
			}
		}
	case []modelruntimeentities.PromptMessager:
		msgs = v
	}
	return msgs
}

// convertRawTools converts the "tools" field (typically []any of maps) into
// []*modelruntimeentities.PromptMessageTool by JSON round-trip.
func convertRawTools(raw any) []*modelruntimeentities.PromptMessageTool {
	if raw == nil {
		return nil
	}
	var tools []*modelruntimeentities.PromptMessageTool
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(data, &tools); err != nil {
		return nil
	}
	return tools
}
