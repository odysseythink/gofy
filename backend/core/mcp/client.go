package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"mlib.com/mlog"
)

// MCPClient communicates with MCP-compatible servers.
type MCPClient struct {
	mu           sync.Mutex
	serverURL    string
	httpClient   *http.Client
	sessionID    string
	capabilities map[string]any
}

// MCPClientConfig configures the MCP client.
type MCPClientConfig struct {
	ServerURL string
	Timeout   time.Duration
	AuthToken string
}

// NewMCPClient creates a new MCP client.
func NewMCPClient(cfg MCPClientConfig) *MCPClient {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &MCPClient{
		serverURL:  cfg.ServerURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Initialize performs the MCP handshake to establish capabilities.
func (c *MCPClient) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  "initialize",
		"id":      1,
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "gofy",
				"version": "1.0.0",
			},
		},
	}

	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("MCP initialize failed: %w", err)
	}

	if result, ok := resp["result"].(map[string]any); ok {
		c.capabilities = result
		if caps, ok := result["capabilities"].(map[string]any); ok {
			c.capabilities = caps
		}
	}

	// Send initialized notification
	notif := map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}
	c.sendRequest(ctx, notif)

	mlog.Infof("MCP client initialized with server %s", c.serverURL)
	return nil
}

// ListTools retrieves the list of tools from the MCP server.
func (c *MCPClient) ListTools(ctx context.Context) ([]MCPTool, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  "tools/list",
		"id":      2,
	}

	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	result, ok := resp["result"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid tools/list response")
	}

	toolsRaw, ok := result["tools"].([]any)
	if !ok {
		return nil, nil
	}

	tools := make([]MCPTool, 0, len(toolsRaw))
	for _, t := range toolsRaw {
		if tm, ok := t.(map[string]any); ok {
			tool := MCPTool{
				Name:        fmt.Sprintf("%v", tm["name"]),
				Description: fmt.Sprintf("%v", tm["description"]),
			}
			if schema, ok := tm["inputSchema"]; ok {
				tool.InputSchema = schema
			}
			tools = append(tools, tool)
		}
	}

	return tools, nil
}

// CallTool invokes a tool on the MCP server.
func (c *MCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]any) (*MCPToolResult, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"id":      3,
		"params": map[string]any{
			"name":      toolName,
			"arguments": arguments,
		},
	}

	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	result, ok := resp["result"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid tools/call response")
	}

	toolResult := &MCPToolResult{}
	if content, ok := result["content"].([]any); ok {
		for _, c := range content {
			if cm, ok := c.(map[string]any); ok {
				toolResult.Content = append(toolResult.Content, MCPContent{
					Type: fmt.Sprintf("%v", cm["type"]),
					Text: fmt.Sprintf("%v", cm["text"]),
				})
			}
		}
	}
	if isErr, ok := result["isError"].(bool); ok {
		toolResult.IsError = isErr
	}

	return toolResult, nil
}

// ListResources retrieves available resources from the MCP server.
func (c *MCPClient) ListResources(ctx context.Context) ([]MCPResource, error) {
	req := map[string]any{
		"jsonrpc": "2.0",
		"method":  "resources/list",
		"id":      4,
	}

	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	result, ok := resp["result"].(map[string]any)
	if !ok {
		return nil, nil
	}

	resourcesRaw, _ := result["resources"].([]any)
	resources := make([]MCPResource, 0)
	for _, r := range resourcesRaw {
		if rm, ok := r.(map[string]any); ok {
			resources = append(resources, MCPResource{
				URI:         fmt.Sprintf("%v", rm["uri"]),
				Name:        fmt.Sprintf("%v", rm["name"]),
				Description: fmt.Sprintf("%v", rm["description"]),
				MimeType:    fmt.Sprintf("%v", rm["mimeType"]),
			})
		}
	}
	return resources, nil
}

// Close closes the MCP session.
func (c *MCPClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.httpClient.CloseIdleConnections()
	return nil
}

func (c *MCPClient) sendRequest(ctx context.Context, payload map[string]any) (map[string]any, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.serverURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if errObj, ok := result["error"]; ok {
		return nil, fmt.Errorf("MCP error: %v", errObj)
	}

	return result, nil
}
