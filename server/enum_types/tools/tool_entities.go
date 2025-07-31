package tools

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"mlib.com/gofy/server/core/exceptions"
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

	ToolProvider_BUILT_IN          ToolProviderType = "builtin"
	ToolProvider_WORKFLOW          ToolProviderType = "workflow"
	ToolProvider_API               ToolProviderType = "api"
	ToolProvider_APP               ToolProviderType = "app"
	ToolProvider_DATASET_RETRIEVAL ToolProviderType = "dataset-retrieval"
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

	ApiProviderAuth_NONE    = "none"
	ApiProviderAuth_API_KEY = "api_key"
)

type MessageType string

const (
	Message_TEXT       MessageType = "text"
	Message_IMAGE      MessageType = "image"
	Message_LINK       MessageType = "link"
	Message_BLOB       MessageType = "blob"
	Message_JSON       MessageType = "json"
	Message_IMAGE_LINK MessageType = "image_link"
	Message_FILE       MessageType = "file"
)

type ToolParameterType string

const (
	ToolParameter_STRING       ToolParameterType = "string"
	ToolParameter_NUMBER       ToolParameterType = "number"
	ToolParameter_BOOLEAN      ToolParameterType = "boolean"
	ToolParameter_SELECT       ToolParameterType = "select"
	ToolParameter_SECRET_INPUT ToolParameterType = "secret-input"
	ToolParameter_FILE         ToolParameterType = "file"
	ToolParameter_FILES        ToolParameterType = "files"
	// deprecated, should not use.
	ToolParameter_SYSTEM_FILES ToolParameterType = "systme-files"
)

func (tp ToolParameterType) CastValue(value any) any {
	// try:
	switch tp {
	case ToolParameter_SECRET_INPUT:
		fallthrough
	case ToolParameter_SELECT:
		fallthrough
	case ToolParameter_STRING:
		if value == nil {
			return ""
		} else {
			if _, ok := value.(string); ok {
				return value.(string)
			} else {
				return fmt.Sprintf("%v", value)
			}
		}

	case ToolParameter_BOOLEAN:
		if value == nil {
			return false
		} else if real_value, ok := value.(string); ok {
			// Allowed YAML boolean value strings: https://yaml.org/type/bool.html
			// and also '0' for False and '1' for True
			switch strings.ToLower(real_value) {
			case "true":
				fallthrough
			case "yes":
				fallthrough
			case "y":
				fallthrough
			case "1":
				return true
			case "false":
				fallthrough
			case "no":
				fallthrough
			case "n":
				fallthrough
			case "0":
				return false
			default:
				tmp, err := cast.ToBoolE(real_value)
				if err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value=%#v is not in correct type.", value)))
				}
				return tmp
			}
		} else if real_value, ok := value.(bool); ok {
			return real_value
		} else {
			tmp, err := cast.ToBoolE(value)
			if err != nil {
				panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value=%#v is not in correct type.", value)))
			}
			return tmp
		}
	case ToolParameter_NUMBER:
		if real_value, ok := value.(int); ok {
			return real_value
		} else if real_value, ok := value.(float64); ok {
			return real_value
		} else if real_value, ok := value.(string); ok {
			if strings.Contains(real_value, ".") {
				tmp, err := strconv.ParseFloat(real_value, 64)
				if err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value=%#v is not in correct type.", value)))
				}
				return tmp
			} else {
				tmp, err := strconv.Atoi(real_value)
				if err != nil {
					panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value=%#v is not in correct type.", value)))
				}
				return tmp
			}
		}
	case ToolParameter_SYSTEM_FILES:
		fallthrough
	case ToolParameter_FILES:
		fallthrough
	case ToolParameter_FILE:
		return value
	default:
		return fmt.Sprintf("%v", value)
	}
	panic(exceptions.NewValueError(fmt.Sprintf("The tool parameter value=%#v is not in correct type.", value)))
}

type ToolParameterForm string

const (
	ToolParameterForm_SCHEMA ToolParameterForm = "schema" //# should be set while adding tool
	ToolParameterForm_FORM   ToolParameterForm = "form"   //# should be set before invoking tool
	ToolParameterForm_LLM    ToolParameterForm = "llm"    //# will be set by LLM
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

type ToolInvokeFrom string

const (

	// """
	// Enum class for tool invoke
	// """

	ToolInvokeFrom_WORKFLOW ToolInvokeFrom = "workflow"
	ToolInvokeFrom_AGENT    ToolInvokeFrom = "agent"
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
