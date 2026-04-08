package tools

import (
	"strings"

	pluginenumtypes "mlib.com/gofy/server/enum_types/plugin"
)

type ToolLabelType string

const (
	ToolLabel_SEARCH        ToolLabelType = "search"
	ToolLabel_IMAGE         ToolLabelType = "image"
	ToolLabel_VIDEOS        ToolLabelType = "videos"
	ToolLabel_WEATHER       ToolLabelType = "weather"
	ToolLabel_FINANCE       ToolLabelType = "finance"
	ToolLabel_DESIGN        ToolLabelType = "design"
	ToolLabel_TRAVEL        ToolLabelType = "travel"
	ToolLabel_SOCIAL        ToolLabelType = "social"
	ToolLabel_NEWS          ToolLabelType = "news"
	ToolLabel_MEDICAL       ToolLabelType = "medical"
	ToolLabel_PRODUCTIVITY  ToolLabelType = "productivity"
	ToolLabel_EDUCATION     ToolLabelType = "education"
	ToolLabel_BUSINESS      ToolLabelType = "business"
	ToolLabel_ENTERTAINMENT ToolLabelType = "entertainment"
	ToolLabel_UTILITIES     ToolLabelType = "utilities"
	ToolLabel_OTHER         ToolLabelType = "other"
)

type ToolProviderType string

const (

	/*
	   Enum class for tool provider
	*/
	ToolProvider_PLUGIN            ToolProviderType = "plugin"
	ToolProvider_BUILT_IN          ToolProviderType = "builtin"
	ToolProvider_WORKFLOW          ToolProviderType = "workflow"
	ToolProvider_API               ToolProviderType = "api"
	ToolProvider_APP               ToolProviderType = "app"
	ToolProvider_DATASET_RETRIEVAL ToolProviderType = "dataset-retrieval"
	ToolProvider_MCP               ToolProviderType = "mcp"
)

type ApiProviderSchemaType string

const (

	/*
	   Enum class for api provider schema type.
	*/

	ApiProviderSchema_OPENAPI        ApiProviderSchemaType = "openapi"
	ApiProviderSchema_SWAGGER        ApiProviderSchemaType = "swagger"
	ApiProviderSchema_OPENAI_PLUGIN  ApiProviderSchemaType = "openai_plugin"
	ApiProviderSchema_OPENAI_ACTIONS ApiProviderSchemaType = "openai_actions"
)

func (e ApiProviderSchemaType) Valid() bool {
	return e == ApiProviderSchema_OPENAPI ||
		e == ApiProviderSchema_SWAGGER ||
		e == ApiProviderSchema_OPENAI_PLUGIN ||
		e == ApiProviderSchema_OPENAI_ACTIONS
}

type ApiProviderAuthType string

const (

	/*
	   Enum class for api provider auth type.
	*/

	ApiProviderAuth_NONE           = "none"
	ApiProviderAuth_API_KEY_HEADER = "api_key_header"
	ApiProviderAuth_API_KEY_QUERY  = "api_key_query"
)

type MessageType string

const (
	Message_TEXT                MessageType = "text"
	Message_IMAGE               MessageType = "image"
	Message_LINK                MessageType = "link"
	Message_BLOB                MessageType = "blob"
	Message_JSON                MessageType = "json"
	Message_IMAGE_LINK          MessageType = "image_link"
	Message_BINARY_LINK         MessageType = "binary_link"
	Message_VARIABLE            MessageType = "variable"
	Message_FILE                MessageType = "file"
	Message_LOG                 MessageType = "log"
	Message_BLOB_CHUNK          MessageType = "blob_chunk"
	Message_RETRIEVER_RESOURCES MessageType = "retriever_resources"
)

type ToolParameterFormType string

const (
	ToolParameterForm_SCHEMA ToolParameterFormType = "schema" // should be set while adding tool
	ToolParameterForm_FORM   ToolParameterFormType = "form"   // should be set before invoking tool
	ToolParameterForm_LLM    ToolParameterFormType = "llm"    // will be set by LLM
)

type CredentialsType string

const (
	Credentials_SECRET_INPUT CredentialsType = "secret-input"
	Credentials_TEXT_INPUT   CredentialsType = "text-input"
	Credentials_SELECT       CredentialsType = "select"
	Credentials_BOOLEAN      CredentialsType = "boolean"
)

type ToolRuntimeVariableType string

const (
	ToolRuntimeVariable_TEXT  ToolRuntimeVariableType = "text"
	ToolRuntimeVariable_IMAGE ToolRuntimeVariableType = "image"
)

type ModelToolPropertyKey string

const (
	ModelToolPropertyKey_IMAGE_PARAMETER_NAME ModelToolPropertyKey = "image_parameter_name"
)

type ToolInvokeFromType string

const (

	// """
	// Enum class for tool invoke
	// """

	ToolInvokeFrom_WORKFLOW ToolInvokeFromType = "workflow"
	ToolInvokeFrom_AGENT    ToolInvokeFromType = "agent"
	ToolInvokeFrom_PLUGIN   ToolInvokeFromType = "plugin"
)

type ToolVariableKey string

const (
	ToolVariableKey_IMAGE    ToolVariableKey = "image"
	ToolVariableKey_DOCUMENT ToolVariableKey = "document"
	ToolVariableKey_VIDEO    ToolVariableKey = "video"
	ToolVariableKey_AUDIO    ToolVariableKey = "audio"
	ToolVariableKey_CUSTOM   ToolVariableKey = "custom"
)

type CredentialType string

const (
	Credential_API_KEY CredentialType = "api-key"
	Credential_OAUTH2  CredentialType = "oauth2"
)

func (c CredentialType) GetName() string {
	if c == Credential_API_KEY {
		return "API KEY"
	} else if c == Credential_OAUTH2 {
		return "AUTH"
	} else {
		return strings.ToUpper(strings.ReplaceAll(string(c), "-", " "))
	}
}
func (c CredentialType) IsEditable() bool {
	return c == Credential_API_KEY
}
func (c CredentialType) IsValidateAllowed() bool {
	return c == Credential_API_KEY
}

type ToolParameterType string

const (
	ToolParameter_STRING         = ToolParameterType(pluginenumtypes.PluginParameter_STRING)
	ToolParameter_NUMBER         = ToolParameterType(pluginenumtypes.PluginParameter_NUMBER)
	ToolParameter_BOOLEAN        = ToolParameterType(pluginenumtypes.PluginParameter_BOOLEAN)
	ToolParameter_SELECT         = ToolParameterType(pluginenumtypes.PluginParameter_SELECT)
	ToolParameter_SECRET_INPUT   = ToolParameterType(pluginenumtypes.PluginParameter_SECRET_INPUT)
	ToolParameter_FILE           = ToolParameterType(pluginenumtypes.PluginParameter_FILE)
	ToolParameter_FILES          = ToolParameterType(pluginenumtypes.PluginParameter_FILES)
	ToolParameter_APP_SELECTOR   = ToolParameterType(pluginenumtypes.PluginParameter_APP_SELECTOR)
	ToolParameter_MODEL_SELECTOR = ToolParameterType(pluginenumtypes.PluginParameter_MODEL_SELECTOR)
	ToolParameter_ANY            = ToolParameterType(pluginenumtypes.PluginParameter_ANY)
	ToolParameter_DYNAMIC_SELECT = ToolParameterType(pluginenumtypes.PluginParameter_DYNAMIC_SELECT)

	// MCP object and array type parameters
	ToolParameter_ARRAY  = ToolParameterType(pluginenumtypes.MCPServerParameter_ARRAY)
	ToolParameter_OBJECT = ToolParameterType(pluginenumtypes.MCPServerParameter_OBJECT)

	// deprecated, should not use.
	ToolParameter_SYSTEM_FILES = ToolParameterType(pluginenumtypes.PluginParameter_SYSTEM_FILES)
)

func (tp ToolParameterType) CastValue(value any) any {
	return pluginenumtypes.CastParameterValue(string(tp), value)
}
func (tp ToolParameterType) AsNormalType() string {
	return pluginenumtypes.AsNormalType(string(tp))
}

type LogStatusType string

const (
	LogStatus_START   LogStatusType = "start"
	LogStatus_ERROR   LogStatusType = "error"
	LogStatus_SUCCESS LogStatusType = "success"
)
