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

// ValidateCredentials validates tool credentials.
func (tm *PluginToolManager) ValidateCredentials(ctx context.Context, tenantID string, provider string, toolName string, credentials map[string]any) error {
	mlog.Infof("validating credentials for tool %s/%s", provider, toolName)
	// TODO: Send validation request to plugin runtime
	return nil
}

// GetRuntimeParameters returns runtime parameters for a tool.
func (tm *PluginToolManager) GetRuntimeParameters(ctx context.Context, tenantID string, provider string, toolName string, credentials map[string]any) ([]map[string]any, error) {
	mlog.Infof("getting runtime parameters for tool %s/%s", provider, toolName)
	// TODO: Query plugin for runtime parameter schema
	return nil, nil
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
		// TODO: Parse declaration to extract tool providers
		providers = append(providers, map[string]any{
			"plugin_id":          inst.PluginID,
			"unique_identifier":  inst.PluginUniqueIdentifier,
			"source":             inst.Source,
		})
	}

	return providers
}
