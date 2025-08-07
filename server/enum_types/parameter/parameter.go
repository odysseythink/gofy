package parameter

type CommonParameterType string

const (
	CommonParameter_SECRET_INPUT   CommonParameterType = "secret-input"
	CommonParameter_TEXT_INPUT     CommonParameterType = "text-input"
	CommonParameter_SELECT         CommonParameterType = "select"
	CommonParameter_STRING         CommonParameterType = "string"
	CommonParameter_NUMBER         CommonParameterType = "number"
	CommonParameter_FILE           CommonParameterType = "file"
	CommonParameter_FILES          CommonParameterType = "files"
	CommonParameter_SYSTEM_FILES   CommonParameterType = "system-files"
	CommonParameter_BOOLEAN        CommonParameterType = "boolean"
	CommonParameter_APP_SELECTOR   CommonParameterType = "app-selector"
	CommonParameter_MODEL_SELECTOR CommonParameterType = "model-selector"
	CommonParameter_TOOLS_SELECTOR CommonParameterType = "array[tools]"
	CommonParameter_ANY            CommonParameterType = "any"

	// Dynamic select parameter
	// Once you are not sure about the available options until authorization is done
	// eg: Select a Slack channel from a Slack workspace
	CommonParameter_DYNAMIC_SELECT CommonParameterType = "dynamic-select"

	// TOOL_SELECTOR = "tool-selector"
	// MCP object and array type parameters
	CommonParameter_ARRAY  CommonParameterType = "array"
	CommonParameter_OBJECT CommonParameterType = "object"
)

type AppSelectorScopeType string

const (
	AppSelectorScope_ALL        AppSelectorScopeType = "all"
	AppSelectorScope_CHAT       AppSelectorScopeType = "chat"
	AppSelectorScope_WORKFLOW   AppSelectorScopeType = "workflow"
	AppSelectorScope_COMPLETION AppSelectorScopeType = "completion"
)

func (scope AppSelectorScopeType) Valid() bool {
	return scope == AppSelectorScope_ALL ||
		scope == AppSelectorScope_CHAT ||
		scope == AppSelectorScope_WORKFLOW ||
		scope == AppSelectorScope_COMPLETION
}

type ModelSelectorScopeType string

const (
	ModelSelectorScope_LLM            ModelSelectorScopeType = "llm"
	ModelSelectorScope_TEXT_EMBEDDING ModelSelectorScopeType = "text-embedding"
	ModelSelectorScope_RERANK         ModelSelectorScopeType = "rerank"
	ModelSelectorScope_TTS            ModelSelectorScopeType = "tts"
	ModelSelectorScope_SPEECH2TEXT    ModelSelectorScopeType = "speech2text"
	ModelSelectorScope_MODERATION     ModelSelectorScopeType = "moderation"
	ModelSelectorScope_VISION         ModelSelectorScopeType = "vision"
)

func (scope ModelSelectorScopeType) Valid() bool {
	return scope == ModelSelectorScope_LLM ||
		scope == ModelSelectorScope_TEXT_EMBEDDING ||
		scope == ModelSelectorScope_RERANK ||
		scope == ModelSelectorScope_TTS ||
		scope == ModelSelectorScope_SPEECH2TEXT ||
		scope == ModelSelectorScope_MODERATION ||
		scope == ModelSelectorScope_VISION
}

type ToolSelectorScopeType string

const (
	ToolSelectorScope_ALL      ToolSelectorScopeType = "all"
	ToolSelectorScope_CUSTOM   ToolSelectorScopeType = "custom"
	ToolSelectorScope_BUILTIN  ToolSelectorScopeType = "builtin"
	ToolSelectorScope_WORKFLOW ToolSelectorScopeType = "workflow"
)

func (scope ToolSelectorScopeType) Valid() bool {
	return scope == ToolSelectorScope_ALL ||
		scope == ToolSelectorScope_CUSTOM ||
		scope == ToolSelectorScope_BUILTIN ||
		scope == ToolSelectorScope_WORKFLOW
}
