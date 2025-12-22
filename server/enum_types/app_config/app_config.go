package appconfig

type VariableEntityType string

const (
	VariableEntity_TEXT_INPUT         VariableEntityType = "text-input"
	VariableEntity_SELECT             VariableEntityType = "select"
	VariableEntity_PARAGRAPH          VariableEntityType = "paragraph"
	VariableEntity_NUMBER             VariableEntityType = "number"
	VariableEntity_EXTERNAL_DATA_TOOL VariableEntityType = "external_data_tool"
	VariableEntity_FILE               VariableEntityType = "file"
	VariableEntity_FILE_LIST          VariableEntityType = "file-list"
)

// PromptType represents prompt type
type PromptType string

const (
	Prompt_SIMPLE PromptType = "simple"
)

func ValidatePromptType(val string) bool {
	return PromptType(val) == Prompt_SIMPLE
}

// RetrieveStrategy represents retrieve strategy
type RetrieveStrategy string

const (
	RetrieveStrategy_SINGLE   RetrieveStrategy = "single"
	RetrieveStrategy_MULTIPLE RetrieveStrategy = "multiple"
)

// EasyUIBasedAppModelConfigFrom represents app model config from
type EasyUIBasedAppModelConfigFrom string

const (
	EasyUIBasedAppModelConfigFrom_ARGS                         EasyUIBasedAppModelConfigFrom = "args"
	EasyUIBasedAppModelConfigFrom_APP_LATEST_CONFIG            EasyUIBasedAppModelConfigFrom = "app-latest-config"
	EasyUIBasedAppModelConfigFrom_CONVERSATION_SPECIFIC_CONFIG EasyUIBasedAppModelConfigFrom = "conversation-specific-config"
)

type MetadataFilteringModeType string

const (
	MetadataFilteringMode_Disabled  MetadataFilteringModeType = "disabled"
	MetadataFilteringMode_Automatic MetadataFilteringModeType = "automatic"
	MetadataFilteringMode_Manual    MetadataFilteringModeType = "manual"
)
