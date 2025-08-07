package plugin

type InstallPluginMessageEventType string

const (
	InstallPluginMessageEvent_Info  InstallPluginMessageEventType = "info"
	InstallPluginMessageEvent_Done  InstallPluginMessageEventType = "done"
	InstallPluginMessageEvent_Error InstallPluginMessageEventType = "error"
)

type PluginInstallTaskStatusType string

const (
	PluginInstallTaskStatus_Pending PluginInstallTaskStatusType = "pending"
	PluginInstallTaskStatus_Running PluginInstallTaskStatusType = "running"
	PluginInstallTaskStatus_Success PluginInstallTaskStatusType = "success"
	PluginInstallTaskStatus_Failed  PluginInstallTaskStatusType = "failed"
)

type AuthorizedCategoryType string

const (
	AuthorizedCategory_Langgenius AuthorizedCategoryType = "langgenius"
	AuthorizedCategory_Partner    AuthorizedCategoryType = "partner"
	AuthorizedCategory_Community  AuthorizedCategoryType = "community"
)
