package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"mlib.com/gofy/server/services"
	"mlib.com/mlog"
)

// DifyInvocation implements BackwardsInvocation by calling gofy services directly.
// Since gofy is a monolith (not a separate daemon), backwards invocations can
// bypass HTTP and invoke services in-process.
type DifyInvocation struct {
	tenantID string
	userID   string
	appID    string
	ctx      context.Context
}

// NewDifyInvocation creates a new backwards invocation handler.
func NewDifyInvocation(tenantID, userID, appID string) *DifyInvocation {
	return &DifyInvocation{
		tenantID: tenantID,
		userID:   userID,
		appID:    appID,
		ctx:      context.Background(),
	}
}

// SetContext sets the context for subsequent invocations.
func (di *DifyInvocation) SetContext(ctx context.Context) {
	di.ctx = ctx
}

// Context returns the current context.
func (di *DifyInvocation) Context() context.Context {
	if di.ctx == nil {
		return context.Background()
	}
	return di.ctx
}

// InvokeLLM calls an LLM model through the model runtime.
func (di *DifyInvocation) InvokeLLM(payload any) (<-chan any, error) {
	ch := make(chan any, 10)

	go func() {
		defer close(ch)

		// TODO: Route through model_runtime to invoke the actual LLM.
		// The payload contains provider, model, prompt_messages, model_parameters, stop, stream, user.
		// Should delegate to core/model_runtime LLM invocation once that layer is wired up.
		mlog.Infof("backwards invocation: InvokeLLM tenant=%s payload=%v", di.tenantID, summarizePayload(payload))

		ch <- map[string]any{
			"type":    "message",
			"content": "[backwards invocation] LLM called (stub)",
		}
	}()

	return ch, nil
}

// InvokeLLMWithStructuredOutput calls an LLM with structured output schema.
func (di *DifyInvocation) InvokeLLMWithStructuredOutput(payload any) (<-chan any, error) {
	// TODO: Forward to model_runtime with structured output constraint.
	// For now, delegate to the regular LLM invocation.
	return di.InvokeLLM(payload)
}

// InvokeTextEmbedding calls a text embedding model.
func (di *DifyInvocation) InvokeTextEmbedding(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeTextEmbedding tenant=%s", di.tenantID)
	// TODO: Route through model_runtime text embedding provider.
	return map[string]any{
		"embeddings": [][]float64{},
		"model":      "",
		"usage":      map[string]int{"total_tokens": 0},
	}, nil
}

// InvokeMultimodalEmbedding calls a multimodal embedding model.
func (di *DifyInvocation) InvokeMultimodalEmbedding(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeMultimodalEmbedding tenant=%s", di.tenantID)
	// TODO: Route through model_runtime multimodal embedding provider.
	return map[string]any{
		"embeddings": [][]float64{},
		"model":      "",
		"usage":      map[string]int{"total_tokens": 0},
	}, nil
}

// InvokeRerank calls a reranking model.
func (di *DifyInvocation) InvokeRerank(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeRerank tenant=%s", di.tenantID)
	// TODO: Route through model_runtime rerank provider.
	return map[string]any{"results": []map[string]any{}}, nil
}

// InvokeMultimodalRerank calls a multimodal reranking model.
func (di *DifyInvocation) InvokeMultimodalRerank(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeMultimodalRerank tenant=%s", di.tenantID)
	// TODO: Route through model_runtime multimodal rerank provider.
	return map[string]any{"results": []map[string]any{}}, nil
}

// InvokeTTS calls a text-to-speech model.
func (di *DifyInvocation) InvokeTTS(payload any) (<-chan any, error) {
	ch := make(chan any, 10)

	go func() {
		defer close(ch)
		mlog.Infof("backwards invocation: InvokeTTS tenant=%s", di.tenantID)
		// TODO: Route through model_runtime TTS provider.
		ch <- map[string]any{"audio": []byte{}}
	}()

	return ch, nil
}

// InvokeSpeech2Text calls a speech-to-text model.
func (di *DifyInvocation) InvokeSpeech2Text(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeSpeech2Text tenant=%s", di.tenantID)
	// TODO: Route through model_runtime STT provider.
	return map[string]any{"text": ""}, nil
}

// InvokeModeration calls a content moderation service.
func (di *DifyInvocation) InvokeModeration(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeModeration tenant=%s", di.tenantID)
	// TODO: Route through moderation provider.
	return map[string]any{"flagged": false}, nil
}

// InvokeTool calls another tool through the tool engine.
func (di *DifyInvocation) InvokeTool(payload any) (<-chan ToolResponseChunk, error) {
	ch := make(chan ToolResponseChunk, 10)

	go func() {
		defer close(ch)
		mlog.Infof("backwards invocation: InvokeTool tenant=%s payload=%v", di.tenantID, summarizePayload(payload))
		// TODO: Route through core/tools tool engine to invoke the requested tool.
		ch <- ToolResponseChunk{
			Type:    ToolResponseChunkTypeText,
			Message: map[string]any{"text": "[tool invoked (stub)]"},
		}
	}()

	return ch, nil
}

// InvokeApp calls another Dify app.
func (di *DifyInvocation) InvokeApp(payload any) (<-chan map[string]any, error) {
	ch := make(chan map[string]any, 10)

	go func() {
		defer close(ch)
		mlog.Infof("backwards invocation: InvokeApp tenant=%s payload=%v", di.tenantID, summarizePayload(payload))
		// TODO: Route through app generator to invoke the target app.
		ch <- map[string]any{
			"type":    "message",
			"content": "[app invoked (stub)]",
		}
	}()

	return ch, nil
}

// InvokeParameterExtractor calls a parameter extraction model.
func (di *DifyInvocation) InvokeParameterExtractor(payload any) (*InvokeNodeResponse, error) {
	mlog.Infof("backwards invocation: InvokeParameterExtractor tenant=%s", di.tenantID)
	// TODO: Route through workflow node executor for parameter extraction.
	return &InvokeNodeResponse{Outputs: map[string]any{}}, nil
}

// InvokeQuestionClassifier calls a question classification model.
func (di *DifyInvocation) InvokeQuestionClassifier(payload any) (*InvokeNodeResponse, error) {
	mlog.Infof("backwards invocation: InvokeQuestionClassifier tenant=%s", di.tenantID)
	// TODO: Route through workflow node executor for question classification.
	return &InvokeNodeResponse{Outputs: map[string]any{}}, nil
}

// InvokeEncrypt encrypts sensitive data.
func (di *DifyInvocation) InvokeEncrypt(payload any) (map[string]any, error) {
	mlog.Infof("backwards invocation: InvokeEncrypt tenant=%s", di.tenantID)
	// TODO: Implement encryption service integration.
	// For now, return the payload data as-is (no-op encryption).
	if m, ok := payload.(map[string]any); ok {
		return m, nil
	}
	return map[string]any{}, nil
}

// InvokeSummary generates a summary.
func (di *DifyInvocation) InvokeSummary(payload any) (*InvokeSummaryResponse, error) {
	mlog.Infof("backwards invocation: InvokeSummary tenant=%s", di.tenantID)
	// TODO: Route through LLM-based summarization service.
	return &InvokeSummaryResponse{Summary: ""}, nil
}

// UploadFile uploads a file from a plugin.
func (di *DifyInvocation) UploadFile(payload any) (*UploadFileResponse, error) {
	mlog.Infof("backwards invocation: UploadFile tenant=%s", di.tenantID)
	// TODO: Save file via storage layer and create upload_file record.
	return &UploadFileResponse{URL: ""}, nil
}

// FetchApp retrieves app metadata by ID.
func (di *DifyInvocation) FetchApp(payload any) (map[string]any, error) {
	// Extract app_id from the payload.
	appID, err := extractStringField(payload, "app_id")
	if err != nil {
		return nil, fmt.Errorf("FetchApp: %w", err)
	}
	if appID == "" {
		appID = di.appID
	}

	mlog.Infof("backwards invocation: FetchApp tenant=%s appID=%s", di.tenantID, appID)

	app, err := services.ServiceGroupApp.App.GetByIDAndTenantID(appID, di.tenantID)
	if err != nil {
		return nil, fmt.Errorf("FetchApp: app not found: %s: %w", appID, err)
	}

	// Convert to a generic map for the plugin.
	appJSON, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("FetchApp: failed to marshal app: %w", err)
	}
	var result map[string]any
	if err := json.Unmarshal(appJSON, &result); err != nil {
		return nil, fmt.Errorf("FetchApp: failed to unmarshal app: %w", err)
	}
	return result, nil
}

// --- helpers ---

// extractStringField extracts a string field from a payload that may be a map or struct.
func extractStringField(payload any, field string) (string, error) {
	switch p := payload.(type) {
	case map[string]any:
		if v, ok := p[field]; ok {
			if s, ok := v.(string); ok {
				return s, nil
			}
		}
		return "", nil
	case map[string]string:
		return p[field], nil
	default:
		// Try JSON round-trip for struct payloads.
		data, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("cannot extract %s from payload: %w", field, err)
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			return "", fmt.Errorf("cannot extract %s from payload: %w", field, err)
		}
		if v, ok := m[field]; ok {
			if s, ok := v.(string); ok {
				return s, nil
			}
		}
		return "", nil
	}
}

// summarizePayload returns a short string summary of a payload for logging.
func summarizePayload(payload any) string {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("%T", payload)
	}
	s := string(data)
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}
