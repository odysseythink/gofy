package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/odysseythink/mlog"
)

// IOTunnel manages bidirectional communication between gofy and plugins.
type IOTunnel struct {
	mu       sync.Mutex
	sessions map[string]*TunnelSession
	timeout  time.Duration
}

// TunnelSession represents an active communication session with a plugin.
type TunnelSession struct {
	SessionID        string
	PluginIdentifier string
	RequestCh        chan *TunnelMessage
	ResponseCh       chan *TunnelMessage
	CreatedAt        time.Time
	Ctx              context.Context
	Cancel           context.CancelFunc
}

// TunnelMessage is a message exchanged through the tunnel.
type TunnelMessage struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"` // request, response, event, error
	Method    string         `json:"method,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	RawData   []byte         `json:"raw_data,omitempty"`
	Error     string         `json:"error,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// NewIOTunnel creates an IO tunnel.
func NewIOTunnel(timeout time.Duration) *IOTunnel {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &IOTunnel{
		sessions: make(map[string]*TunnelSession),
		timeout:  timeout,
	}
}

// CreateSession starts a new communication session.
func (t *IOTunnel) CreateSession(sessionID, pluginIdentifier string) *TunnelSession {
	t.mu.Lock()
	defer t.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	session := &TunnelSession{
		SessionID:        sessionID,
		PluginIdentifier: pluginIdentifier,
		RequestCh:        make(chan *TunnelMessage, 10),
		ResponseCh:       make(chan *TunnelMessage, 10),
		CreatedAt:        time.Now(),
		Ctx:              ctx,
		Cancel:           cancel,
	}
	t.sessions[sessionID] = session
	return session
}

// GetSession returns an active session.
func (t *IOTunnel) GetSession(sessionID string) *TunnelSession {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.sessions[sessionID]
}

// CloseSession terminates a session.
func (t *IOTunnel) CloseSession(sessionID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if session, ok := t.sessions[sessionID]; ok {
		session.Cancel()
		close(session.RequestCh)
		close(session.ResponseCh)
		delete(t.sessions, sessionID)
	}
}

// SendRequest sends a request to a plugin and waits for a response.
func (t *IOTunnel) SendRequest(sessionID string, method string, data map[string]any) (*TunnelMessage, error) {
	session := t.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	req := &TunnelMessage{
		ID:        fmt.Sprintf("%s-%d", sessionID, time.Now().UnixNano()),
		Type:      "request",
		Method:    method,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Send request
	select {
	case session.RequestCh <- req:
	case <-session.Ctx.Done():
		return nil, fmt.Errorf("session timeout")
	}

	// Wait for response
	select {
	case resp := <-session.ResponseCh:
		if resp.Error != "" {
			return nil, fmt.Errorf("plugin error: %s", resp.Error)
		}
		return resp, nil
	case <-session.Ctx.Done():
		return nil, fmt.Errorf("response timeout")
	}
}

// SendStream sends a request and returns a channel of streaming responses.
func (t *IOTunnel) SendStream(sessionID string, method string, data map[string]any) (<-chan *TunnelMessage, error) {
	session := t.GetSession(sessionID)
	if session == nil {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	req := &TunnelMessage{
		ID:        fmt.Sprintf("%s-%d", sessionID, time.Now().UnixNano()),
		Type:      "request",
		Method:    method,
		Data:      data,
		Timestamp: time.Now(),
	}

	select {
	case session.RequestCh <- req:
	case <-session.Ctx.Done():
		return nil, fmt.Errorf("session timeout")
	}

	return session.ResponseCh, nil
}

// HandlePluginResponse routes a response from the plugin runtime to the waiting session.
func (t *IOTunnel) HandlePluginResponse(sessionID string, response *TunnelMessage) error {
	session := t.GetSession(sessionID)
	if session == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	select {
	case session.ResponseCh <- response:
		return nil
	case <-session.Ctx.Done():
		return fmt.Errorf("session expired")
	}
}

// InvokeTool invokes a tool through the tunnel.
func (t *IOTunnel) InvokeTool(pluginIdentifier, toolName string, params map[string]any, credentials map[string]any) (<-chan *ToolResponseChunk, error) {
	sessionID := fmt.Sprintf("tool-%s-%d", toolName, time.Now().UnixNano())
	_ = t.CreateSession(sessionID, pluginIdentifier)

	ch := make(chan *ToolResponseChunk, 10)

	go func() {
		defer close(ch)
		defer t.CloseSession(sessionID)

		resp, err := t.SendRequest(sessionID, "invoke_tool", map[string]any{
			"tool_name":       toolName,
			"tool_parameters": params,
			"credentials":     credentials,
		})
		if err != nil {
			ch <- &ToolResponseChunk{
				Type:    ToolResponseChunkTypeText,
				Message: map[string]any{"text": fmt.Sprintf("error: %v", err)},
			}
			return
		}

		// Parse tool response
		if resp.Data != nil {
			chunkType, _ := resp.Data["type"].(string)
			ch <- &ToolResponseChunk{
				Type:    ToolResponseChunkType(chunkType),
				Message: resp.Data,
			}
		}
	}()

	return ch, nil
}

// InvokeModel invokes a model through the tunnel.
func (t *IOTunnel) InvokeModel(pluginIdentifier, modelType, model string, params map[string]any, credentials map[string]any) (<-chan any, error) {
	sessionID := fmt.Sprintf("model-%s-%d", model, time.Now().UnixNano())
	_ = t.CreateSession(sessionID, pluginIdentifier)

	ch := make(chan any, 50)

	go func() {
		defer close(ch)
		defer t.CloseSession(sessionID)

		streamCh, err := t.SendStream(sessionID, "invoke_model", map[string]any{
			"model_type":  modelType,
			"model":       model,
			"parameters":  params,
			"credentials": credentials,
		})
		if err != nil {
			ch <- map[string]any{"type": "error", "error": err.Error()}
			return
		}

		for msg := range streamCh {
			if msg.Type == "error" {
				ch <- map[string]any{"type": "error", "error": msg.Error}
				return
			}
			ch <- msg.Data
		}
	}()

	return ch, nil
}

// ActiveSessions returns the number of active sessions.
func (t *IOTunnel) ActiveSessions() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.sessions)
}

// CleanupExpiredSessions removes expired sessions.
func (t *IOTunnel) CleanupExpiredSessions() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	cleaned := 0
	for id, session := range t.sessions {
		select {
		case <-session.Ctx.Done():
			delete(t.sessions, id)
			cleaned++
		default:
		}
	}
	if cleaned > 0 {
		mlog.Infof("cleaned up %d expired tunnel sessions", cleaned)
	}
	return cleaned
}
