package enumtypes

type APIBasedExtensionPoint string

const (
	APIBasedExtensionPoint_APP_EXTERNAL_DATA_TOOL_QUERY APIBasedExtensionPoint = "app.external_data_tool.query"
	APIBasedExtensionPoint_PING                         APIBasedExtensionPoint = "ping"
	APIBasedExtensionPoint_APP_MODERATION_INPUT         APIBasedExtensionPoint = "app.moderation.input"
	APIBasedExtensionPoint_APP_MODERATION_OUTPUT        APIBasedExtensionPoint = "app.moderation.output"
)
