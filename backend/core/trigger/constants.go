package trigger

// Trigger node type constants (matching workflow node types)
const (
	TRIGGER_WEBHOOK_NODE_TYPE  = "trigger-webhook"
	TRIGGER_SCHEDULE_NODE_TYPE = "trigger-schedule"
	TRIGGER_PLUGIN_NODE_TYPE   = "trigger-plugin"
)

// AppTriggerStatus represents the status of an app trigger
type AppTriggerStatus string

const (
	AppTriggerStatusEnabled      AppTriggerStatus = "enabled"
	AppTriggerStatusDisabled     AppTriggerStatus = "disabled"
	AppTriggerStatusUnauthorized AppTriggerStatus = "unauthorized"
	AppTriggerStatusRateLimited  AppTriggerStatus = "rate_limited"
)

// WorkflowTriggerStatus represents the execution status of a workflow trigger
type WorkflowTriggerStatus string

const (
	WorkflowTriggerStatusPending     WorkflowTriggerStatus = "pending"
	WorkflowTriggerStatusQueued      WorkflowTriggerStatus = "queued"
	WorkflowTriggerStatusRunning     WorkflowTriggerStatus = "running"
	WorkflowTriggerStatusSucceeded   WorkflowTriggerStatus = "succeeded"
	WorkflowTriggerStatusPaused      WorkflowTriggerStatus = "paused"
	WorkflowTriggerStatusFailed      WorkflowTriggerStatus = "failed"
	WorkflowTriggerStatusRateLimited WorkflowTriggerStatus = "rate_limited"
	WorkflowTriggerStatusRetrying    WorkflowTriggerStatus = "retrying"
)

// AppTriggerType represents the type of app trigger
type AppTriggerType string

const (
	AppTriggerTypeWebhook  AppTriggerType = "trigger-webhook"
	AppTriggerTypeSchedule AppTriggerType = "trigger-schedule"
	AppTriggerTypePlugin   AppTriggerType = "trigger-plugin"
)

// CredentialType represents the type of credentials for plugin triggers
type CredentialType string

const (
	CredentialTypeOAuth  CredentialType = "oauth"
	CredentialTypeAPIKey CredentialType = "api_key"
	CredentialTypeManual CredentialType = "manual"
)
