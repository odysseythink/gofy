package plugin

import "context"

// InvokeType represents the type of backwards invocation.
type InvokeType string

const (
	INVOKE_TYPE_LLM                      InvokeType = "llm"
	INVOKE_TYPE_LLM_STRUCTURED_OUTPUT    InvokeType = "llm_structured_output"
	INVOKE_TYPE_TEXT_EMBEDDING           InvokeType = "text_embedding"
	INVOKE_TYPE_MULTIMODAL_EMBEDDING     InvokeType = "multimodal_embedding"
	INVOKE_TYPE_RERANK                   InvokeType = "rerank"
	INVOKE_TYPE_MULTIMODAL_RERANK        InvokeType = "multimodal_rerank"
	INVOKE_TYPE_TTS                      InvokeType = "tts"
	INVOKE_TYPE_SPEECH2TEXT              InvokeType = "speech2text"
	INVOKE_TYPE_MODERATION               InvokeType = "moderation"
	INVOKE_TYPE_TOOL                     InvokeType = "tool"
	INVOKE_TYPE_NODE_PARAMETER_EXTRACTOR InvokeType = "node_parameter_extractor"
	INVOKE_TYPE_NODE_QUESTION_CLASSIFIER InvokeType = "node_question_classifier"
	INVOKE_TYPE_APP                      InvokeType = "app"
	INVOKE_TYPE_STORAGE                  InvokeType = "storage"
	INVOKE_TYPE_ENCRYPT                  InvokeType = "encrypt"
	INVOKE_TYPE_SYSTEM_SUMMARY           InvokeType = "system_summary"
	INVOKE_TYPE_UPLOAD_FILE              InvokeType = "upload_file"
	INVOKE_TYPE_FETCH_APP                InvokeType = "fetch_app"
)

// BaseInvokeDifyRequest is the common base for all backwards invocation requests.
type BaseInvokeDifyRequest struct {
	TenantID string     `json:"tenant_id"`
	UserID   string     `json:"user_id"`
	Type     InvokeType `json:"type"`
}

// InvokeNodeResponse is the common response for parameter extractor / question classifier nodes.
type InvokeNodeResponse struct {
	ProcessData map[string]any `json:"process_data"`
	Outputs     map[string]any `json:"outputs"`
	Inputs      map[string]any `json:"inputs"`
}

// InvokeSummaryResponse is the response for summary invocations.
type InvokeSummaryResponse struct {
	Summary string `json:"summary"`
}

// UploadFileResponse is the response for file upload invocations.
type UploadFileResponse struct {
	URL string `json:"url"`
}

// BackwardsInvocation defines the interface for plugins to call back into Dify.
// Each method corresponds to a backwards invocation type that a plugin may trigger.
//
// For streaming responses, methods return a read-only channel. Callers should range
// over the channel until it is closed. For non-streaming responses, methods return
// a concrete result or map.
//
// Request types referenced here (e.g. InvokeLLMRequest) are intentionally left as
// `any` placeholders because the full request/response structs live in entities/plugin/request.go
// and entities/model_runtime/. They will be refined as the plugin system matures.
type BackwardsInvocation interface {
	SetContext(ctx context.Context)
	Context() context.Context

	// Model invocations
	InvokeLLM(payload any) (<-chan any, error)
	InvokeLLMWithStructuredOutput(payload any) (<-chan any, error)
	InvokeTextEmbedding(payload any) (any, error)
	InvokeMultimodalEmbedding(payload any) (any, error)
	InvokeRerank(payload any) (any, error)
	InvokeMultimodalRerank(payload any) (any, error)
	InvokeTTS(payload any) (<-chan any, error)
	InvokeSpeech2Text(payload any) (any, error)
	InvokeModeration(payload any) (any, error)

	// Tool invocation
	InvokeTool(payload any) (<-chan ToolResponseChunk, error)

	// App invocation
	InvokeApp(payload any) (<-chan map[string]any, error)

	// Workflow node invocations
	InvokeParameterExtractor(payload any) (*InvokeNodeResponse, error)
	InvokeQuestionClassifier(payload any) (*InvokeNodeResponse, error)

	// Encryption
	InvokeEncrypt(payload any) (map[string]any, error)

	// Summary
	InvokeSummary(payload any) (*InvokeSummaryResponse, error)

	// File operations
	UploadFile(payload any) (*UploadFileResponse, error)

	// App metadata
	FetchApp(payload any) (map[string]any, error)
}
