package ops

import (
	"context"
	"time"
)

// TraceProvider identifies a trace backend.
type TraceProvider string

const (
	TraceProviderLangfuse  TraceProvider = "langfuse"
	TraceProviderLangSmith TraceProvider = "langsmith"
	TraceProviderOpik      TraceProvider = "opik"
	TraceProviderOTLP      TraceProvider = "otlp"
)

// TraceInstance is the interface all trace backends implement.
type TraceInstance interface {
	// Trace records a complete trace span.
	Trace(ctx context.Context, span *TraceSpan) error
	// Flush ensures all pending traces are sent.
	Flush(ctx context.Context) error
	// Close shuts down the trace backend.
	Close() error
	// Provider returns the provider type.
	Provider() TraceProvider
}

// TraceSpan represents a single span in a trace.
type TraceSpan struct {
	TraceID    string         `json:"trace_id"`
	SpanID     string         `json:"span_id"`
	ParentID   string         `json:"parent_id,omitempty"`
	Name       string         `json:"name"`
	Kind       SpanKind       `json:"kind"`
	Status     SpanStatus     `json:"status"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Events     []SpanEvent    `json:"events,omitempty"`
	// LLM-specific fields
	Input      any            `json:"input,omitempty"`
	Output     any            `json:"output,omitempty"`
	Model      string         `json:"model,omitempty"`
	TokenUsage *TokenUsage    `json:"token_usage,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type SpanKind string

const (
	SpanKindWorkflow  SpanKind = "workflow"
	SpanKindNode      SpanKind = "node"
	SpanKindLLM       SpanKind = "llm"
	SpanKindTool      SpanKind = "tool"
	SpanKindRetrieval SpanKind = "retrieval"
	SpanKindEmbedding SpanKind = "embedding"
)

type SpanStatus string

const (
	SpanStatusOK    SpanStatus = "ok"
	SpanStatusError SpanStatus = "error"
)

type SpanEvent struct {
	Name       string         `json:"name"`
	Timestamp  time.Time      `json:"timestamp"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
