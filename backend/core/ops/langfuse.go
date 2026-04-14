package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/odysseythink/mlog"
)

// LangfuseTraceConfig configures the Langfuse backend.
type LangfuseTraceConfig struct {
	PublicKey string
	SecretKey string
	Host      string // default: https://cloud.langfuse.com
}

// LangfuseInstance implements TraceInstance for Langfuse.
type LangfuseInstance struct {
	config LangfuseTraceConfig
	client *http.Client
	mu     sync.Mutex
	buffer []*TraceSpan
}

// NewLangfuseInstance creates a new Langfuse trace backend.
func NewLangfuseInstance(cfg LangfuseTraceConfig) *LangfuseInstance {
	if cfg.Host == "" {
		cfg.Host = "https://cloud.langfuse.com"
	}
	return &LangfuseInstance{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
		buffer: make([]*TraceSpan, 0),
	}
}

func (l *LangfuseInstance) Provider() TraceProvider {
	return TraceProviderLangfuse
}

func (l *LangfuseInstance) Trace(ctx context.Context, span *TraceSpan) error {
	l.mu.Lock()
	l.buffer = append(l.buffer, span)
	shouldFlush := len(l.buffer) >= 10
	l.mu.Unlock()

	if shouldFlush {
		return l.Flush(ctx)
	}
	return nil
}

func (l *LangfuseInstance) Flush(ctx context.Context) error {
	l.mu.Lock()
	spans := l.buffer
	l.buffer = make([]*TraceSpan, 0)
	l.mu.Unlock()

	if len(spans) == 0 {
		return nil
	}

	// Convert spans to Langfuse ingestion format
	batch := make([]map[string]any, 0, len(spans))
	for _, span := range spans {
		event := l.spanToLangfuseEvent(span)
		batch = append(batch, event)
	}

	payload := map[string]any{"batch": batch}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/public/ingestion", l.config.Host)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(l.config.PublicKey, l.config.SecretKey)

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("langfuse ingestion failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("langfuse returned status %d", resp.StatusCode)
	}

	mlog.Infof("langfuse: flushed %d spans", len(spans))
	return nil
}

func (l *LangfuseInstance) Close() error {
	return l.Flush(context.Background())
}

func (l *LangfuseInstance) spanToLangfuseEvent(span *TraceSpan) map[string]any {
	eventType := "span-create"
	if span.Kind == SpanKindLLM {
		eventType = "generation-create"
	}

	event := map[string]any{
		"type": eventType,
		"body": map[string]any{
			"id":        span.SpanID,
			"traceId":   span.TraceID,
			"parentId":  span.ParentID,
			"name":      span.Name,
			"startTime": span.StartTime.Format(time.RFC3339Nano),
			"endTime":   span.EndTime.Format(time.RFC3339Nano),
			"input":     span.Input,
			"output":    span.Output,
			"metadata":  span.Metadata,
			"level":     "DEFAULT",
		},
	}

	bodyMap := event["body"].(map[string]any)

	if span.Status == SpanStatusError {
		bodyMap["level"] = "ERROR"
	}
	if span.Model != "" {
		bodyMap["model"] = span.Model
	}
	if span.TokenUsage != nil {
		bodyMap["usage"] = map[string]any{
			"input":  span.TokenUsage.PromptTokens,
			"output": span.TokenUsage.CompletionTokens,
			"total":  span.TokenUsage.TotalTokens,
		}
	}

	return event
}
