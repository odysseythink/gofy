package v1

type ApiGroup struct {
	AppApi
	AccountApi
	FeatureApi
	WorkspaceApi
	ModelProvideApi
	DatasetApi
	TagApi
	DraftWorkflowApi
	LoginApi
	SetupApi
	FileApi
	MemberApi
	WorkflowRunApi
	ApiKeyApi
	ExtensionApi
	ChatApi
	MessageApi
	ConversationApi
	StatisticApi
	ToolsApi
	VersionApi
	OpsTraceApi
	ModelsApi
	PluginApi
	RuleGenerateApi
	WebApi
}

var ApiGroupApp = new(ApiGroup)
