package dbengine

import (
	"fmt"

	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

// AutoMigrate runs GORM auto-migration for all registered models.
// This creates tables if they don't exist and adds missing columns.
// It does NOT delete columns or change column types.
func AutoMigrate() error {
	db := Instance().DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	mlog.Info("running database auto-migration...")

	// Core models
	err := db.AutoMigrate(
		// Account/Auth
		&models.Account{},
		&models.Tenant{},
		&models.TenantAccountJoin{},
		&models.AccountIntegrate{},
		&models.InvitationCode{},
		&models.TenantPluginPermission{},
		&models.TenantPluginAutoUpgradeStrategy{},

		// App
		&models.App{},
		&models.AppModelConfig{},
		&models.DifySetup{},
		&models.RecommendedApp{},
		&models.InstalledApp{},
		&models.Site{},
		&models.ApiToken{},

		// Conversation & Message
		&models.Conversation{},
		&models.Message{},
		&models.MessageFeedback{},
		&models.MessageFile{},
		&models.MessageAnnotation{},
		&models.AppAnnotationSetting{},
		&models.AppAnnotationHitHistory{},
		&models.MessageChain{},
		&models.MessageAgentThought{},
		&models.EndUser{},
		&models.OperationLog{},

		// Workflow
		&models.Workflow{},
		&models.WorkflowRun{},
		&models.WorkflowNodeExecution{},
		&models.WorkflowAppLog{},
		&models.ConversationVariable{},
		&models.WorkflowNodeExecutionOffload{},
		&models.WorkflowArchiveLog{},
		&models.WorkflowDraftVariable{},
		&models.WorkflowDraftVariableFile{},
		&models.WorkflowPause{},
		&models.WorkflowPauseReason{},

		// Dataset
		&models.Dataset{},
		&models.DatasetProcessRule{},
		&models.Document{},
		&models.DocumentSegment{},
		&models.ChildChunk{},
		&models.AppDatasetJoin{},
		&models.DatasetQuery{},
		&models.DatasetKeywordTable{},
		&models.Embedding{},
		&models.DatasetCollectionBinding{},
		&models.DatasetPermission{},
		&models.ExternalKnowledgeApi{},
		&models.ExternalKnowledgeBinding{},
		&models.DatasetAutoDisableLog{},
		&models.RateLimitLog{},
		&models.DatasetMetadata{},
		&models.DatasetMetadataBinding{},

		// Pipeline
		&models.PipelineBuiltInTemplate{},
		&models.PipelineCustomizedTemplate{},
		&models.Pipeline{},
		&models.DocumentPipelineExecutionLog{},
		&models.PipelineRecommendedPlugin{},
		&models.SegmentAttachmentBinding{},
		&models.DocumentSegmentSummary{},

		// Tools
		&models.ToolOAuthSystemClient{},
		&models.ToolOAuthTenantClient{},
		&models.BuiltinToolProvider{},
		&models.ApiToolProvider{},
		&models.ToolLabelBinding{},
		&models.WorkflowToolProvider{},
		&models.MCPToolProvider{},
		&models.ToolModelInvoke{},
		&models.ToolConversationVariable{},
		&models.ToolFile{},

		// Provider
		&models.Provider{},
		&models.ProviderModel{},
		&models.TenantDefaultModel{},
		&models.TenantPreferredModelProvider{},
		&models.ProviderOrder{},
		&models.ProviderModelSetting{},
		&models.LoadBalancingModelConfig{},
		&models.ProviderCredential{},
		&models.ProviderModelCredential{},

		// Trigger
		&models.TriggerSubscription{},
		&models.TriggerOAuthSystemClient{},
		&models.TriggerOAuthTenantClient{},
		&models.WorkflowTriggerLog{},
		&models.WorkflowWebhookTrigger{},
		&models.WorkflowPluginTrigger{},
		&models.AppTrigger{},
		&models.WorkflowSchedulePlan{},

		// Human Input
		&models.HumanInputForm{},
		&models.HumanInputDelivery{},
		&models.HumanInputFormRecipient{},

		// Plugin
		&models.Plugin{},
		&models.PluginInstallation{},
		&models.PluginDeclaration{},

		// Other
		&models.UploadFile{},
		&models.ApiRequest{},
		&models.DatasetRetrieverResource{},
		&models.Tag{},
		&models.TagBinding{},
		&models.TraceAppConfig{},
		&models.TenantCreditPool{},
		&models.OAuthProviderApp{},
		&models.AppMCPServer{},
		&models.SavedMessage{},
		&models.PinnedConversation{},
		&models.DataSourceApiKeyAuthBinding{},
		&models.DataSourceOauthBinding{},
		&models.APIBasedExtension{},
		&models.CeleryTask{},
		&models.CeleryTaskSet{},
		&models.PluginInstallTask{},
	)

	if err != nil {
		mlog.Errorf("auto-migration failed: %v", err)
		return err
	}

	mlog.Info("database auto-migration completed successfully")
	return nil
}
