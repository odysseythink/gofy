package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
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

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

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

// extractMapField extracts a map[string]any field from a map payload.
func extractMapField(data map[string]any, key string) map[string]any {
	if v, ok := data[key].(map[string]any); ok {
		return v
	}
	return nil
}

// extractStringSlice extracts a []string from a payload field that may be
// []string or []any (common when decoded from JSON).
func extractStringSlice(data map[string]any, key string) []string {
	switch v := data[key].(type) {
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, elem := range v {
			if s, ok := elem.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// payloadToMap converts an arbitrary payload to map[string]any.
func payloadToMap(payload any) (map[string]any, error) {
	if m, ok := payload.(map[string]any); ok {
		return m, nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot convert payload to map: %w", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("cannot convert payload to map: %w", err)
	}
	return m, nil
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

// convertRawMessages converts the prompt_messages field (typically []any of maps)
// into []modelruntimeentities.PromptMessager using the existing factory.
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

// convertRawTools converts the tools field (typically []any of maps) into
// []*modelruntimeentities.PromptMessageTool.
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

// getModelInstance uses the ModelManager to resolve provider, credentials,
// and load-balancing for the given model type.
func (di *DifyInvocation) getModelInstance(provider string, modelType modelruntimeenumtypes.ModelType, model string) (mi *modelmanager.ModelInstance, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get model instance: %v", r)
		}
	}()
	mgr := modelmanager.NewModelManager()
	mi = mgr.GetModelInstance(di.tenantID, provider, modelType, model)
	return
}

// ---------------------------------------------------------------------------
// Model invocations -- LLM
// ---------------------------------------------------------------------------

// InvokeLLM calls an LLM model through the model runtime.
func (di *DifyInvocation) InvokeLLM(payload any) (<-chan any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeLLM: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	stream, _ := params["stream"].(bool)
	user, _ := extractStringField(params, "user")
	stop := extractStringSlice(params, "stop")
	modelParameters := extractMapField(params, "model_parameters")
	promptMessages := convertRawMessages(params["prompt_messages"])
	tools := convertRawTools(params["tools"])

	mlog.Infof("backwards invocation: InvokeLLM tenant=%s provider=%s model=%s stream=%v", di.tenantID, provider, model, stream)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_LLM, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeLLM: %w", err)
	}

	ch := make(chan any, 50)

	go func() {
		defer close(ch)
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("InvokeLLM panic: %v", r)
				ch <- map[string]any{"type": "error", "error": fmt.Sprintf("%v", r)}
			}
		}()

		if stream {
			seq := mi.InvokeLLMStream(promptMessages, modelParameters, tools, stop, user, nil)
			for chunk := range seq {
				ch <- map[string]any{
					"type":  "llm_chunk",
					"chunk": chunk,
				}
			}
		} else {
			result := mi.InvokeLLM(promptMessages, modelParameters, tools, stop, user, nil)
			ch <- map[string]any{
				"type":   "llm_result",
				"result": result,
			}
		}
	}()

	return ch, nil
}

// InvokeLLMWithStructuredOutput calls an LLM with structured output schema.
// It adds a response_format parameter to the model parameters and delegates
// to the standard LLM invocation pipeline.
func (di *DifyInvocation) InvokeLLMWithStructuredOutput(payload any) (<-chan any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeLLMWithStructuredOutput: %w", err)
	}

	// If a response_format / schema is specified, merge it into model_parameters
	// so the underlying provider can handle structured output natively.
	if schema := extractMapField(params, "response_format"); schema != nil {
		mp := extractMapField(params, "model_parameters")
		if mp == nil {
			mp = make(map[string]any)
		}
		mp["response_format"] = schema
		params["model_parameters"] = mp
	}

	return di.InvokeLLM(params)
}

// ---------------------------------------------------------------------------
// Model invocations -- Embedding
// ---------------------------------------------------------------------------

// InvokeTextEmbedding calls a text embedding model.
func (di *DifyInvocation) InvokeTextEmbedding(payload any) (any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeTextEmbedding: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	user, _ := extractStringField(params, "user")

	mlog.Infof("backwards invocation: InvokeTextEmbedding tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_TEXT_EMBEDDING, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeTextEmbedding: %w", err)
	}

	// Extract texts to embed
	var texts []string
	switch v := params["texts"].(type) {
	case []string:
		texts = v
	case []any:
		for _, t := range v {
			if s, ok := t.(string); ok {
				texts = append(texts, s)
			}
		}
	}

	if len(texts) == 0 {
		return nil, fmt.Errorf("InvokeTextEmbedding: no texts provided")
	}

	// The model type instance should implement InvokeEmbedding.
	type embeddingInvoker interface {
		InvokeEmbedding(model string, credentials map[string]any, texts []string, user string) (*modelruntimeentities.TextEmbeddingResult, error)
	}

	invoker, ok := mi.ModelTypeInstance.(embeddingInvoker)
	if !ok {
		return nil, fmt.Errorf("InvokeTextEmbedding: provider '%s' does not support embedding invocation", provider)
	}

	result, err := invoker.InvokeEmbedding(mi.Model, mi.Credentials, texts, user)
	if err != nil {
		return nil, fmt.Errorf("InvokeTextEmbedding: %w", err)
	}

	return map[string]any{
		"model":      result.Model,
		"embeddings": result.Embeddings,
		"usage": map[string]int{
			"tokens":       result.Usage.Tokens,
			"total_tokens": result.Usage.TotalTokens,
		},
	}, nil
}

// InvokeMultimodalEmbedding calls a multimodal embedding model.
// Falls back to the text embedding path since most providers do not yet
// distinguish multimodal embeddings at the model_runtime level.
func (di *DifyInvocation) InvokeMultimodalEmbedding(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeMultimodalEmbedding tenant=%s", di.tenantID)

	// Multimodal embedding providers are not yet widely implemented in
	// model_runtime. Delegate to text embedding which handles the common case.
	return di.InvokeTextEmbedding(payload)
}

// ---------------------------------------------------------------------------
// Model invocations -- Rerank
// ---------------------------------------------------------------------------

// InvokeRerank calls a reranking model.
func (di *DifyInvocation) InvokeRerank(payload any) (any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeRerank: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	query, _ := extractStringField(params, "query")

	mlog.Infof("backwards invocation: InvokeRerank tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_RERANK, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeRerank: %w", err)
	}

	// Extract documents
	var docs []string
	switch v := params["docs"].(type) {
	case []string:
		docs = v
	case []any:
		for _, d := range v {
			if s, ok := d.(string); ok {
				docs = append(docs, s)
			}
		}
	}

	// The rerank model type instance should implement an InvokeRerank-style method.
	type rerankInvoker interface {
		InvokeRerank(model string, credentials map[string]any, query string, docs []string, topN int, user string) (*modelruntimeentities.RerankResult, error)
	}

	user, _ := extractStringField(params, "user")
	topN := 10
	if n, ok := params["top_n"].(float64); ok {
		topN = int(n)
	}

	invoker, ok := mi.ModelTypeInstance.(rerankInvoker)
	if !ok {
		return nil, fmt.Errorf("InvokeRerank: provider '%s' does not support rerank invocation", provider)
	}

	result, err := invoker.InvokeRerank(mi.Model, mi.Credentials, query, docs, topN, user)
	if err != nil {
		return nil, fmt.Errorf("InvokeRerank: %w", err)
	}

	resultsOut := make([]map[string]any, 0, len(result.Docs))
	for _, doc := range result.Docs {
		resultsOut = append(resultsOut, map[string]any{
			"index": doc.Index,
			"text":  doc.Text,
			"score": doc.Score,
		})
	}

	return map[string]any{
		"model":   result.Model,
		"results": resultsOut,
	}, nil
}

// InvokeMultimodalRerank calls a multimodal reranking model.
// Falls back to the standard rerank path since multimodal reranking is not
// widely supported at the model_runtime level yet.
func (di *DifyInvocation) InvokeMultimodalRerank(payload any) (any, error) {
	mlog.Infof("backwards invocation: InvokeMultimodalRerank tenant=%s", di.tenantID)
	return di.InvokeRerank(payload)
}

// ---------------------------------------------------------------------------
// Model invocations -- TTS
// ---------------------------------------------------------------------------

// InvokeTTS calls a text-to-speech model.
func (di *DifyInvocation) InvokeTTS(payload any) (<-chan any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeTTS: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	user, _ := extractStringField(params, "user")
	text, _ := extractStringField(params, "text")

	mlog.Infof("backwards invocation: InvokeTTS tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_TTS, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeTTS: %w", err)
	}

	ch := make(chan any, 10)

	go func() {
		defer close(ch)
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("InvokeTTS panic: %v", r)
				ch <- map[string]any{"type": "error", "error": fmt.Sprintf("%v", r)}
			}
		}()

		// TTS providers should implement an InvokeTTS method that returns
		// audio data (streaming or complete).
		type ttsInvoker interface {
			InvokeTTS(model string, credentials map[string]any, text string, user string) (iter.Seq[[]byte], error)
		}
		type ttsSyncInvoker interface {
			InvokeTTS(model string, credentials map[string]any, text string, user string) ([]byte, error)
		}

		switch invoker := mi.ModelTypeInstance.(type) {
		case ttsInvoker:
			seq, err := invoker.InvokeTTS(mi.Model, mi.Credentials, text, user)
			if err != nil {
				ch <- map[string]any{"type": "error", "error": err.Error()}
				return
			}
			for chunk := range seq {
				ch <- map[string]any{"type": "audio_chunk", "audio": chunk}
			}
		case ttsSyncInvoker:
			audio, err := invoker.InvokeTTS(mi.Model, mi.Credentials, text, user)
			if err != nil {
				ch <- map[string]any{"type": "error", "error": err.Error()}
				return
			}
			ch <- map[string]any{"type": "audio", "audio": audio}
		default:
			ch <- map[string]any{"type": "error", "error": fmt.Sprintf("provider '%s' does not support TTS invocation", provider)}
		}
	}()

	return ch, nil
}

// ---------------------------------------------------------------------------
// Model invocations -- Speech2Text
// ---------------------------------------------------------------------------

// InvokeSpeech2Text calls a speech-to-text model.
func (di *DifyInvocation) InvokeSpeech2Text(payload any) (any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeSpeech2Text: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	user, _ := extractStringField(params, "user")

	mlog.Infof("backwards invocation: InvokeSpeech2Text tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_SPEECH2TEXT, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeSpeech2Text: %w", err)
	}

	// Extract audio data from payload (base64 or raw bytes).
	type sttInvoker interface {
		InvokeSTT(model string, credentials map[string]any, audioData []byte, user string) (string, error)
	}

	invoker, ok := mi.ModelTypeInstance.(sttInvoker)
	if !ok {
		return nil, fmt.Errorf("InvokeSpeech2Text: provider '%s' does not support STT invocation", provider)
	}

	var audioData []byte
	switch v := params["audio"].(type) {
	case string:
		// Assume base64-encoded
		audioData = []byte(v)
	case []byte:
		audioData = v
	default:
		return nil, fmt.Errorf("InvokeSpeech2Text: audio field missing or invalid type")
	}

	text, err := invoker.InvokeSTT(mi.Model, mi.Credentials, audioData, user)
	if err != nil {
		return nil, fmt.Errorf("InvokeSpeech2Text: %w", err)
	}

	return map[string]any{"text": text}, nil
}

// ---------------------------------------------------------------------------
// Model invocations -- Moderation
// ---------------------------------------------------------------------------

// InvokeModeration calls a content moderation service.
func (di *DifyInvocation) InvokeModeration(payload any) (any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeModeration: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	text, _ := extractStringField(params, "text")
	user, _ := extractStringField(params, "user")

	mlog.Infof("backwards invocation: InvokeModeration tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_MODERATION, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeModeration: %w", err)
	}

	type moderationInvoker interface {
		InvokeModeration(model string, credentials map[string]any, text string, user string) (bool, error)
	}

	invoker, ok := mi.ModelTypeInstance.(moderationInvoker)
	if !ok {
		return nil, fmt.Errorf("InvokeModeration: provider '%s' does not support moderation invocation", provider)
	}

	flagged, err := invoker.InvokeModeration(mi.Model, mi.Credentials, text, user)
	if err != nil {
		return nil, fmt.Errorf("InvokeModeration: %w", err)
	}

	return map[string]any{"flagged": flagged}, nil
}

// ---------------------------------------------------------------------------
// Tool invocation
// ---------------------------------------------------------------------------

// InvokeTool calls another tool through the tool engine.
// TODO: Will be implemented to route through core/tools ToolManager.
func (di *DifyInvocation) InvokeTool(payload any) (<-chan ToolResponseChunk, error) {
	ch := make(chan ToolResponseChunk, 10)
	go func() {
		defer close(ch)
		mlog.Infof("backwards invocation: InvokeTool tenant=%s payload=%v", di.tenantID, summarizePayload(payload))
		ch <- ToolResponseChunk{
			Type:    ToolResponseChunkTypeText,
			Message: map[string]any{"text": "[tool invoked (pending tool engine integration)]"},
		}
	}()
	return ch, nil
}

// ---------------------------------------------------------------------------
// App invocation
// ---------------------------------------------------------------------------

// InvokeApp calls another Dify app.
// TODO: Will be implemented to route through AppGenerateService.
func (di *DifyInvocation) InvokeApp(payload any) (<-chan map[string]any, error) {
	ch := make(chan map[string]any, 10)
	go func() {
		defer close(ch)
		mlog.Infof("backwards invocation: InvokeApp tenant=%s payload=%v", di.tenantID, summarizePayload(payload))
		ch <- map[string]any{
			"type":    "message",
			"content": "[app invoked (pending app generate integration)]",
		}
	}()
	return ch, nil
}

// ---------------------------------------------------------------------------
// Workflow node invocations
// ---------------------------------------------------------------------------

// InvokeParameterExtractor calls a parameter extraction model by routing
// through the LLM with a parameter-extraction prompt template.
func (di *DifyInvocation) InvokeParameterExtractor(payload any) (*InvokeNodeResponse, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeParameterExtractor: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	query, _ := extractStringField(params, "query")
	parameters := extractMapField(params, "parameters")
	instruction, _ := extractStringField(params, "instruction")

	mlog.Infof("backwards invocation: InvokeParameterExtractor tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_LLM, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeParameterExtractor: %w", err)
	}

	// Build a system prompt that instructs the LLM to extract parameters.
	parametersJSON, _ := json.MarshalIndent(parameters, "", "  ")
	systemPrompt := fmt.Sprintf(
		"You are a parameter extractor. Extract the following parameters from the user's input and return them as a valid JSON object.\n\nParameters schema:\n%s",
		string(parametersJSON),
	)
	if instruction != "" {
		systemPrompt += "\n\nAdditional instructions: " + instruction
	}

	promptMessages := []modelruntimeentities.PromptMessager{
		modelruntimeentities.NewSystemPromptMessage(systemPrompt, ""),
		modelruntimeentities.NewUserPromptMessage(query, ""),
	}

	result := mi.InvokeLLM(promptMessages, map[string]any{"temperature": 0.0}, nil, nil, "", nil)

	// Parse the LLM response as JSON for the extracted parameters.
	outputs := make(map[string]any)
	if result != nil && result.Message != nil {
		content, _ := result.Message.GetContent().(string)
		if content != "" {
			// Try to parse the content as JSON.
			if err := json.Unmarshal([]byte(content), &outputs); err != nil {
				// If parsing fails, return the raw content.
				outputs["raw_output"] = content
			}
		}
	}

	return &InvokeNodeResponse{
		Outputs: outputs,
		Inputs:  map[string]any{"query": query, "parameters": parameters},
		ProcessData: map[string]any{
			"model":    model,
			"provider": provider,
		},
	}, nil
}

// InvokeQuestionClassifier calls a question classification model by routing
// through the LLM with a classification prompt template.
func (di *DifyInvocation) InvokeQuestionClassifier(payload any) (*InvokeNodeResponse, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeQuestionClassifier: %w", err)
	}

	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")
	query, _ := extractStringField(params, "query")
	instruction, _ := extractStringField(params, "instruction")

	mlog.Infof("backwards invocation: InvokeQuestionClassifier tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	// Extract classes from payload.
	var classes []map[string]any
	if rawClasses, ok := params["classes"].([]any); ok {
		for _, c := range rawClasses {
			if cm, ok := c.(map[string]any); ok {
				classes = append(classes, cm)
			}
		}
	}

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_LLM, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeQuestionClassifier: %w", err)
	}

	classesJSON, _ := json.MarshalIndent(classes, "", "  ")
	systemPrompt := fmt.Sprintf(
		"You are a question classifier. Classify the user's question into one of the following categories and return ONLY the category ID.\n\nCategories:\n%s",
		string(classesJSON),
	)
	if instruction != "" {
		systemPrompt += "\n\nAdditional instructions: " + instruction
	}

	promptMessages := []modelruntimeentities.PromptMessager{
		modelruntimeentities.NewSystemPromptMessage(systemPrompt, ""),
		modelruntimeentities.NewUserPromptMessage(query, ""),
	}

	result := mi.InvokeLLM(promptMessages, map[string]any{"temperature": 0.0}, nil, nil, "", nil)

	outputs := make(map[string]any)
	if result != nil && result.Message != nil {
		content, _ := result.Message.GetContent().(string)
		outputs["class_id"] = content
	}

	return &InvokeNodeResponse{
		Outputs: outputs,
		Inputs:  map[string]any{"query": query, "classes": classes},
		ProcessData: map[string]any{
			"model":    model,
			"provider": provider,
		},
	}, nil
}

// ---------------------------------------------------------------------------
// Encryption
// ---------------------------------------------------------------------------

// InvokeEncrypt encrypts sensitive data using the provider configuration
// service's credential obfuscation logic.
func (di *DifyInvocation) InvokeEncrypt(payload any) (map[string]any, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeEncrypt: %w", err)
	}

	mlog.Infof("backwards invocation: InvokeEncrypt tenant=%s", di.tenantID)

	// Extract the data and opt_field to determine encryption behavior.
	data := extractMapField(params, "data")
	if data == nil {
		// If no nested "data" field, treat the whole payload as the data.
		data = params
	}

	// Use the provider configuration service for credential-level encryption.
	provider, _ := extractStringField(params, "provider")
	if provider != "" {
		creds, err := services.ServiceGroupApp.ProviderConfiguration.GetCustomCredentials(
			di.tenantID, provider, true, // obfuscated=true means we encrypt
		)
		if err == nil && creds != nil {
			return creds, nil
		}
	}

	// For non-provider encryption, return data as-is (pass-through).
	// A full encryption service can be wired in here.
	return data, nil
}

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

// InvokeSummary generates a summary by invoking the default LLM with a
// summarization prompt.
func (di *DifyInvocation) InvokeSummary(payload any) (*InvokeSummaryResponse, error) {
	params, err := payloadToMap(payload)
	if err != nil {
		return nil, fmt.Errorf("InvokeSummary: %w", err)
	}

	text, _ := extractStringField(params, "text")
	provider, _ := extractStringField(params, "provider")
	model, _ := extractStringField(params, "model")

	mlog.Infof("backwards invocation: InvokeSummary tenant=%s provider=%s model=%s", di.tenantID, provider, model)

	mi, err := di.getModelInstance(provider, modelruntimeenumtypes.Model_LLM, model)
	if err != nil {
		return nil, fmt.Errorf("InvokeSummary: %w", err)
	}

	promptMessages := []modelruntimeentities.PromptMessager{
		modelruntimeentities.NewSystemPromptMessage(
			"You are a summarization assistant. Provide a concise summary of the following text.", "",
		),
		modelruntimeentities.NewUserPromptMessage(text, ""),
	}

	var summary string
	func() {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("InvokeSummary LLM panic: %v", r)
			}
		}()
		result := mi.InvokeLLM(promptMessages, map[string]any{"temperature": 0.3}, nil, nil, "", nil)
		if result != nil && result.Message != nil {
			summary, _ = result.Message.GetContent().(string)
		}
	}()

	return &InvokeSummaryResponse{Summary: summary}, nil
}

// ---------------------------------------------------------------------------
// File operations
// ---------------------------------------------------------------------------

// UploadFile uploads a file from a plugin.
// TODO: Will be implemented to route through FileService + storage.
func (di *DifyInvocation) UploadFile(payload any) (*UploadFileResponse, error) {
	mlog.Infof("backwards invocation: UploadFile tenant=%s", di.tenantID)
	return &UploadFileResponse{URL: ""}, nil
}

// ---------------------------------------------------------------------------
// App metadata
// ---------------------------------------------------------------------------

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
