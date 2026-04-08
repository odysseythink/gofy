package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mlib.com/mlog"
)

// LangSmithTraceConfig configures the LangSmith backend.
type LangSmithTraceConfig struct {
	APIKey   string
	Project  string
	Endpoint string // default: https://api.smith.langchain.com
}

// LangSmithInstance implements TraceInstance for LangSmith.
type LangSmithInstance struct {
	config LangSmithTraceConfig
	client *http.Client
}

// NewLangSmithInstance creates a new LangSmith trace backend.
func NewLangSmithInstance(cfg LangSmithTraceConfig) *LangSmithInstance {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://api.smith.langchain.com"
	}
	if cfg.Project == "" {
		cfg.Project = "default"
	}
	return &LangSmithInstance{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (ls *LangSmithInstance) Provider() TraceProvider { return TraceProviderLangSmith }

func (ls *LangSmithInstance) Trace(ctx context.Context, span *TraceSpan) error {
	run := map[string]any{
		"id":             span.SpanID,
		"trace_id":       span.TraceID,
		"parent_run_id":  span.ParentID,
		"name":           span.Name,
		"run_type":       ls.kindToRunType(span.Kind),
		"start_time":     span.StartTime.Format(time.RFC3339Nano),
		"end_time":       span.EndTime.Format(time.RFC3339Nano),
		"inputs":         span.Input,
		"outputs":        span.Output,
		"extra":          span.Metadata,
		"session_name":   ls.config.Project,
	}

	if span.Status == SpanStatusError {
		run["error"] = span.Attributes["error"]
	}

	body, err := json.Marshal(run)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/runs", ls.config.Endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", ls.config.APIKey)

	resp, err := ls.client.Do(req)
	if err != nil {
		mlog.Errorf("langsmith trace failed: %v", err)
		return nil // non-blocking
	}
	defer resp.Body.Close()
	return nil
}

func (ls *LangSmithInstance) Flush(ctx context.Context) error { return nil }
func (ls *LangSmithInstance) Close() error                    { return nil }

func (ls *LangSmithInstance) kindToRunType(kind SpanKind) string {
	switch kind {
	case SpanKindLLM:
		return "llm"
	case SpanKindTool:
		return "tool"
	case SpanKindRetrieval:
		return "retriever"
	case SpanKindEmbedding:
		return "embedding"
	default:
		return "chain"
	}
}
