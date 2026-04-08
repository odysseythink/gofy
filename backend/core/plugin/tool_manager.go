package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"mlib.com/mlog"
)

// PluginToolManager manages tool-type plugins.
type PluginToolManager struct{}

// NewPluginToolManager creates a new tool manager.
func NewPluginToolManager() *PluginToolManager {
	return &PluginToolManager{}
}

// InvokeTool invokes a tool from a plugin.
func (tm *PluginToolManager) InvokeTool(ctx context.Context, tenantID string, provider string, toolName string, toolParameters map[string]any, credentials map[string]any) (<-chan *ToolResponseChunk, error) {
	mgr := Manager()
	if mgr == nil {
		return nil, fmt.Errorf("plugin manager not initialized")
	}

	// Find the plugin that provides this tool
	runtime, ok := mgr.GetRuntime(provider)
	if !ok {
		return nil, fmt.Errorf("plugin runtime not found for provider: %s", provider)
	}

	ch := make(chan *ToolResponseChunk, 10)

	go func() {
		defer close(ch)

		// Create session
		session := &Session{
			TenantID:               tenantID,
			PluginUniqueIdentifier: provider,
			Action:                 "invoke_tool",
		}

		// Build invocation payload
		payload := map[string]any{
			"tool_name":       toolName,
			"tool_parameters": toolParameters,
			"credentials":     credentials,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			mlog.Errorf("failed to marshal tool invocation payload: %v", err)
			ch <- &ToolResponseChunk{
				Type:    ToolResponseChunkTypeText,
				Message: map[string]any{"text": fmt.Sprintf("error: %v", err)},
			}
			return
		}

		// Send invocation request to plugin via session I/O
		if err := runtime.Write(session.ID, "invoke_tool", data); err != nil {
			mlog.Errorf("failed to send tool invocation: %v", err)
			ch <- &ToolResponseChunk{
				Type:    ToolResponseChunkTypeText,
				Message: map[string]any{"text": fmt.Sprintf("error: %v", err)},
			}
			return
		}

		// Listen for responses from the plugin
		msgCh, err := runtime.Listen(session.ID)
		if err != nil {
			mlog.Errorf("failed to listen for tool responses: %v", err)
			ch <- &ToolResponseChunk{
				Type:    ToolResponseChunkTypeText,
				Message: map[string]any{"text": fmt.Sprintf("error: %v", err)},
			}
			return
		}

		for msg := range msgCh {
			var responseData map[string]any
			if err := json.Unmarshal(msg.Data, &responseData); err != nil {
				mlog.Errorf("failed to unmarshal tool response: %v", err)
				continue
			}
			chunkType, _ := responseData["type"].(string)
			ch <- &ToolResponseChunk{
				Type:    ToolResponseChunkType(chunkType),
				Message: responseData,
			}
		}
	}()

	return ch, nil
}

// ValidateCredentials validates tool credentials by attempting a test invocation.
func (tm *PluginToolManager) ValidateCredentials(ctx context.Context, tenantID string, provider string, toolName string, credentials map[string]any) error {
	mlog.Infof("validating credentials for tool %s/%s", provider, toolName)

	mgr := Manager()
	if mgr == nil {
		return fmt.Errorf("plugin manager not initialized")
	}

	// Check if the plugin provides this tool
	installation := mgr.GetPlugin(tenantID, provider)
	if installation == nil {
		return fmt.Errorf("tool provider '%s' not installed for tenant", provider)
	}

	// Validate required credential fields from declaration
	decl := mgr.GetPluginDeclaration(provider)
	if decl == nil {
		return fmt.Errorf("plugin declaration not found for provider '%s'", provider)
	}

	// Parse declaration to check credential requirements
	var declaration map[string]any
	if err := json.Unmarshal([]byte(decl.Declaration), &declaration); err != nil {
		return fmt.Errorf("failed to parse plugin declaration: %w", err)
	}

	// Check credential_schema requirements if present
	if credSchema, ok := declaration["credentials_for_provider"].(map[string]any); ok {
		if required, ok := credSchema["required"].([]any); ok {
			for _, r := range required {
				key, _ := r.(string)
				if key == "" {
					continue
				}
				val, exists := credentials[key]
				if !exists || val == nil || val == "" {
					return fmt.Errorf("credential '%s' is required but missing or empty", key)
				}
			}
		}
	}

	// Validate non-empty credential values
	for key, val := range credentials {
		if val == nil || val == "" {
			return fmt.Errorf("credential '%s' is required but empty", key)
		}
	}

	return nil
}

// GetRuntimeParameters returns runtime parameters for a tool.
func (tm *PluginToolManager) GetRuntimeParameters(ctx context.Context, tenantID string, provider string, toolName string, credentials map[string]any) ([]map[string]any, error) {
	mlog.Infof("getting runtime parameters for tool %s/%s", provider, toolName)

	mgr := Manager()
	if mgr == nil {
		return nil, fmt.Errorf("plugin manager not initialized")
	}

	// Try to get parameters from plugin runtime
	runtime, ok := mgr.GetRuntime(provider)
	if !ok {
		// Return empty parameters if runtime not active
		return []map[string]any{}, nil
	}

	// Send request to runtime to get parameters
	credJSON, _ := json.Marshal(credentials)
	msg := &SessionMessage{
		Data: []byte(fmt.Sprintf(`{"type":"get_runtime_parameters","tool_name":"%s","credentials":%s}`, toolName, string(credJSON))),
	}

	if lr, ok := runtime.(*LocalRuntime); ok {
		if err := lr.SendMessage(msg); err != nil {
			return nil, fmt.Errorf("failed to send runtime parameters request: %w", err)
		}
		// Wait for response with timeout
		select {
		case response := <-lr.Messages():
			var params []map[string]any
			if err := json.Unmarshal(response.Data, &params); err != nil {
				return nil, fmt.Errorf("failed to parse runtime parameters response: %w", err)
			}
			return params, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return []map[string]any{}, nil
}

// ListToolProviders returns all available tool providers from installed plugins.
func (tm *PluginToolManager) ListToolProviders(tenantID string) []map[string]any {
	mgr := Manager()
	if mgr == nil {
		return nil
	}

	installations, _ := mgr.ListPlugins(tenantID, 1, 1000)
	providers := make([]map[string]any, 0)

	for _, inst := range installations {
		decl := mgr.GetPluginDeclaration(inst.PluginUniqueIdentifier)
		if decl == nil {
			continue
		}

		// Parse declaration to extract tool-related metadata
		var declaration map[string]any
		if err := json.Unmarshal([]byte(decl.Declaration), &declaration); err != nil {
			mlog.Errorf("failed to parse declaration for %s: %v", inst.PluginUniqueIdentifier, err)
			continue
		}

		// Extract tools list from declaration
		tools, _ := declaration["tools"].([]any)
		toolNames := make([]string, 0, len(tools))
		for _, t := range tools {
			if tool, ok := t.(map[string]any); ok {
				if name, ok := tool["name"].(string); ok {
					toolNames = append(toolNames, name)
				}
			}
		}

		providers = append(providers, map[string]any{
			"plugin_id":         inst.PluginID,
			"unique_identifier": inst.PluginUniqueIdentifier,
			"source":            inst.Source,
			"tools":             toolNames,
			"label":             declaration["label"],
			"description":       declaration["description"],
		})
	}

	return providers
}

// GetToolDeclaration returns the tool declaration from a plugin.
func (tm *PluginToolManager) GetToolDeclaration(tenantID, provider, toolName string) map[string]any {
	mgr := Manager()
	if mgr == nil {
		return nil
	}

	decl := mgr.GetPluginDeclaration(provider)
	if decl == nil {
		return nil
	}

	var declaration map[string]any
	if err := json.Unmarshal([]byte(decl.Declaration), &declaration); err != nil {
		mlog.Errorf("failed to parse declaration for %s: %v", provider, err)
		return nil
	}

	tools, _ := declaration["tools"].([]any)
	for _, t := range tools {
		tool, ok := t.(map[string]any)
		if !ok {
			continue
		}
		name, _ := tool["name"].(string)
		if name == toolName {
			return tool
		}
	}
	return nil
}

// ListToolsForProvider returns all tools from a specific provider.
func (tm *PluginToolManager) ListToolsForProvider(tenantID, provider string) []map[string]any {
	mgr := Manager()
	if mgr == nil {
		return nil
	}

	decl := mgr.GetPluginDeclaration(provider)
	if decl == nil {
		return nil
	}

	var declaration map[string]any
	if err := json.Unmarshal([]byte(decl.Declaration), &declaration); err != nil {
		mlog.Errorf("failed to parse declaration for %s: %v", provider, err)
		return nil
	}

	tools, _ := declaration["tools"].([]any)
	result := make([]map[string]any, 0, len(tools))
	for _, t := range tools {
		if tool, ok := t.(map[string]any); ok {
			result = append(result, tool)
		}
	}
	return result
}
