package plugin

import (
	"context"
	"fmt"
	"iter"

	"github.com/odysseythink/mlog"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	modelproviders "mlib.com/gofy/server/core/model_runtime/model_provides"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

// PluginModelClient manages model-type plugins (LLM, embedding, rerank, TTS, STT, moderation).
type PluginModelClient struct{}

func NewPluginModelClient() *PluginModelClient {
	return &PluginModelClient{}
}

// getModelInstance resolves a ModelInstance via the ModelManager, recovering
// from panics that the manager may raise for missing providers/credentials.
func (mc *PluginModelClient) getModelInstance(tenantID, provider string, modelType modelruntimeenumtypes.ModelType, model string) (mi *modelmanager.ModelInstance, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get model instance: %v", r)
		}
	}()
	mgr := modelmanager.NewModelManager()
	mi = mgr.GetModelInstance(tenantID, provider, modelType, model)
	return
}

// InvokeLLM invokes an LLM through the model runtime.
func (mc *PluginModelClient) InvokeLLM(ctx context.Context, tenantID string, provider string, model string, promptMessages any, modelParameters map[string]any, stop []string, stream bool, user string) (<-chan any, error) {
	mlog.Infof("plugin model client: InvokeLLM provider=%s model=%s stream=%v", provider, model, stream)

	mi, err := mc.getModelInstance(tenantID, provider, modelruntimeenumtypes.Model_LLM, model)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM instance: %w", err)
	}

	// Convert raw prompt messages to typed PromptMessager slice.
	messages := convertRawMessages(promptMessages)

	// Convert raw tools if present in model parameters.
	var tools []*modelruntimeentities.PromptMessageTool
	if rawTools, ok := modelParameters["tools"]; ok {
		tools = convertRawTools(rawTools)
		delete(modelParameters, "tools")
	}

	ch := make(chan any, 50)
	go func() {
		defer close(ch)
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("plugin model client LLM panic: %v", r)
				ch <- map[string]any{"type": "error", "error": fmt.Sprintf("%v", r)}
			}
		}()

		if stream {
			seq := mi.InvokeLLMStream(messages, modelParameters, tools, stop, user, nil)
			for chunk := range seq {
				ch <- map[string]any{
					"type":  "llm_chunk",
					"chunk": chunk,
				}
			}
		} else {
			result := mi.InvokeLLM(messages, modelParameters, tools, stop, user, nil)
			ch <- map[string]any{
				"type":   "llm_result",
				"result": result,
			}
		}
	}()
	return ch, nil
}

// InvokeTextEmbedding invokes a text embedding model.
func (mc *PluginModelClient) InvokeTextEmbedding(ctx context.Context, tenantID string, provider string, model string, texts []string, user string) (any, error) {
	mlog.Infof("plugin model client: InvokeTextEmbedding provider=%s model=%s", provider, model)

	mi, err := mc.getModelInstance(tenantID, provider, modelruntimeenumtypes.Model_TEXT_EMBEDDING, model)
	if err != nil {
		return nil, fmt.Errorf("embedding model not found: %w", err)
	}

	type embeddingInvoker interface {
		InvokeEmbedding(model string, credentials map[string]any, texts []string, user string) (*modelruntimeentities.TextEmbeddingResult, error)
	}

	invoker, ok := mi.ModelTypeInstance.(embeddingInvoker)
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support embedding invocation", provider)
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

// InvokeRerank invokes a reranking model.
func (mc *PluginModelClient) InvokeRerank(ctx context.Context, tenantID string, provider string, model string, query string, docs []string, topN int, user string) (any, error) {
	mlog.Infof("plugin model client: InvokeRerank provider=%s model=%s", provider, model)

	mi, err := mc.getModelInstance(tenantID, provider, modelruntimeenumtypes.Model_RERANK, model)
	if err != nil {
		return nil, fmt.Errorf("rerank model not found: %w", err)
	}

	type rerankInvoker interface {
		InvokeRerank(model string, credentials map[string]any, query string, docs []string, topN int, user string) (*modelruntimeentities.RerankResult, error)
	}

	invoker, ok := mi.ModelTypeInstance.(rerankInvoker)
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support rerank invocation", provider)
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

// InvokeTTS invokes text-to-speech.
func (mc *PluginModelClient) InvokeTTS(ctx context.Context, tenantID string, provider string, model string, text string, voice string, user string) (any, error) {
	mlog.Infof("plugin model client: InvokeTTS provider=%s model=%s", provider, model)

	mi, err := mc.getModelInstance(tenantID, provider, modelruntimeenumtypes.Model_TTS, model)
	if err != nil {
		return nil, fmt.Errorf("TTS model not found: %w", err)
	}

	// TTS providers may return streaming audio or complete audio.
	type ttsStreamInvoker interface {
		InvokeTTS(model string, credentials map[string]any, text string, user string) (iter.Seq[[]byte], error)
	}
	type ttsSyncInvoker interface {
		InvokeTTS(model string, credentials map[string]any, text string, user string) ([]byte, error)
	}

	switch invoker := mi.ModelTypeInstance.(type) {
	case ttsStreamInvoker:
		seq, err := invoker.InvokeTTS(mi.Model, mi.Credentials, text, user)
		if err != nil {
			return nil, fmt.Errorf("InvokeTTS: %w", err)
		}
		// Collect all audio chunks into a single result.
		var audioChunks [][]byte
		for chunk := range seq {
			audioChunks = append(audioChunks, chunk)
		}
		return map[string]any{"audio_chunks": audioChunks}, nil
	case ttsSyncInvoker:
		audio, err := invoker.InvokeTTS(mi.Model, mi.Credentials, text, user)
		if err != nil {
			return nil, fmt.Errorf("InvokeTTS: %w", err)
		}
		return map[string]any{"audio": audio}, nil
	default:
		return nil, fmt.Errorf("provider '%s' does not support TTS invocation", provider)
	}
}

// InvokeSpeech2Text invokes speech-to-text.
func (mc *PluginModelClient) InvokeSpeech2Text(ctx context.Context, tenantID string, provider string, model string, file any, user string) (any, error) {
	mlog.Infof("plugin model client: InvokeSpeech2Text provider=%s model=%s", provider, model)

	mi, err := mc.getModelInstance(tenantID, provider, modelruntimeenumtypes.Model_SPEECH2TEXT, model)
	if err != nil {
		return nil, fmt.Errorf("STT model not found: %w", err)
	}

	type sttInvoker interface {
		InvokeSTT(model string, credentials map[string]any, audioData []byte, user string) (string, error)
	}

	invoker, ok := mi.ModelTypeInstance.(sttInvoker)
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support STT invocation", provider)
	}

	var audioData []byte
	switch v := file.(type) {
	case string:
		audioData = []byte(v)
	case []byte:
		audioData = v
	default:
		return nil, fmt.Errorf("InvokeSpeech2Text: file must be string or []byte, got %T", file)
	}

	text, err := invoker.InvokeSTT(mi.Model, mi.Credentials, audioData, user)
	if err != nil {
		return nil, fmt.Errorf("InvokeSpeech2Text: %w", err)
	}

	return map[string]any{"text": text}, nil
}

// InvokeModeration invokes content moderation.
func (mc *PluginModelClient) InvokeModeration(ctx context.Context, tenantID string, provider string, model string, text string, user string) (any, error) {
	mlog.Infof("plugin model client: InvokeModeration provider=%s model=%s", provider, model)

	mi, err := mc.getModelInstance(tenantID, provider, modelruntimeenumtypes.Model_MODERATION, model)
	if err != nil {
		return nil, fmt.Errorf("moderation model not found: %w", err)
	}

	type moderationInvoker interface {
		InvokeModeration(model string, credentials map[string]any, text string, user string) (bool, error)
	}

	invoker, ok := mi.ModelTypeInstance.(moderationInvoker)
	if !ok {
		return nil, fmt.Errorf("provider '%s' does not support moderation invocation", provider)
	}

	flagged, err := invoker.InvokeModeration(mi.Model, mi.Credentials, text, user)
	if err != nil {
		return nil, fmt.Errorf("InvokeModeration: %w", err)
	}

	return map[string]any{"flagged": flagged}, nil
}

// ListModelProviders returns all model providers with their models for a tenant.
func (mc *PluginModelClient) ListModelProviders(tenantID string) []map[string]any {
	factory := &modelproviders.ModelProviderFactory{}
	providers := factory.GetProviders()
	result := make([]map[string]any, 0, len(providers))
	for _, p := range providers {
		result = append(result, map[string]any{
			"provider":              p.Provider,
			"label":                 p.Label,
			"supported_model_types": p.SupportedModelTypes,
			"models":                p.Models,
		})
	}
	return result
}
