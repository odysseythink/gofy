package mcp

// MCPTool represents a tool exposed by an MCP server.
type MCPTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema,omitempty"`
}

// MCPToolResult is the result of calling an MCP tool.
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError"`
}

// MCPContent represents a content item in tool results.
type MCPContent struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// MCPResource represents a resource on an MCP server.
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// MCPPrompt represents a prompt template from an MCP server.
type MCPPrompt struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Arguments   []MCPPromptArg `json:"arguments,omitempty"`
}

// MCPPromptArg is an argument for an MCP prompt.
type MCPPromptArg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}
