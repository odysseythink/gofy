package tool

import (
	toolsentities "mlib.com/gofy/server/entities/tools"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
)

type Tooler interface {
	ToolProviderType() toolsenumtypes.ToolProviderType
	RealInvoke(user_id string, tool_parameters map[string]any) (*toolsentities.ToolInvokeMessage, []*toolsentities.ToolInvokeMessage)
	Invoke(user_id string, tool_parameters map[string]any, tooler Tooler) []*toolsentities.ToolInvokeMessage
	ValidateCredentials(credentials map[string]any, parameters map[string]any, format_only bool) string
	GetRuntimeParameters() []*toolsentities.ToolParameter
}
