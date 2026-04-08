package models

import "time"

// TriggerSubscription [...]
type TriggerSubscription struct {
	ID                  string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name                string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	TenantID            string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	UserID              string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	ProviderID          string     `gorm:"column:provider_id;type:varchar(255);not null" json:"provider_id"`
	EndpointID          string     `gorm:"column:endpoint_id;type:varchar(255);not null;uniqueIndex" json:"endpoint_id"`
	Parameters          string     `gorm:"column:parameters;type:json;not null" json:"parameters"`
	Properties          string     `gorm:"column:properties;type:json;not null" json:"properties"`
	Credentials         string     `gorm:"column:credentials;type:json;not null" json:"credentials"`
	CredentialType      string     `gorm:"column:credential_type;type:varchar(50);not null" json:"credential_type"`
	CredentialExpiresAt int        `gorm:"column:credential_expires_at;type:int;default:-1" json:"credential_expires_at"`
	ExpiresAt           int        `gorm:"column:expires_at;type:int;default:-1" json:"expires_at"`
	CreatedAt           *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (TriggerSubscription) TableName() string {
	return "trigger_subscriptions"
}

// TriggerOAuthSystemClient [...]
type TriggerOAuthSystemClient struct {
	ID                   string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	PluginID             string     `gorm:"column:plugin_id;type:varchar(255);not null" json:"plugin_id"`
	Provider             string     `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	EncryptedOAuthParams string     `gorm:"column:encrypted_oauth_params;type:text;not null" json:"encrypted_oauth_params"`
	CreatedAt            *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (TriggerOAuthSystemClient) TableName() string {
	return "trigger_oauth_system_clients"
}

// TriggerOAuthTenantClient [...]
type TriggerOAuthTenantClient struct {
	ID                   string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID             string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	PluginID             string     `gorm:"column:plugin_id;type:varchar(255);not null" json:"plugin_id"`
	Provider             string     `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	Enabled              bool       `gorm:"column:enabled;type:bool;not null;default:true" json:"enabled"`
	EncryptedOAuthParams string     `gorm:"column:encrypted_oauth_params;type:text;not null;default:'{}'" json:"encrypted_oauth_params"`
	CreatedAt            *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (TriggerOAuthTenantClient) TableName() string {
	return "trigger_oauth_tenant_clients"
}

// WorkflowTriggerLog [...]
type WorkflowTriggerLog struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID           string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	WorkflowID      string     `gorm:"column:workflow_id;type:varchar(36);not null;index" json:"workflow_id"`
	WorkflowRunID   *string    `gorm:"column:workflow_run_id;type:varchar(36);index" json:"workflow_run_id"`
	RootNodeID      *string    `gorm:"column:root_node_id;type:varchar(255)" json:"root_node_id"`
	TriggerMetadata string     `gorm:"column:trigger_metadata;type:text;not null" json:"trigger_metadata"`
	TriggerType     string     `gorm:"column:trigger_type;type:varchar(50);not null" json:"trigger_type"`
	TriggerData     string     `gorm:"column:trigger_data;type:text;not null" json:"trigger_data"`
	Inputs          string     `gorm:"column:inputs;type:text;not null" json:"inputs"`
	Outputs         *string    `gorm:"column:outputs;type:text" json:"outputs"`
	Status          string     `gorm:"column:status;type:varchar(50);not null;index" json:"status"`
	Error           *string    `gorm:"column:error;type:text" json:"error"`
	QueueName       string     `gorm:"column:queue_name;type:varchar(100);not null" json:"queue_name"`
	CeleryTaskID    *string    `gorm:"column:celery_task_id;type:varchar(255)" json:"celery_task_id"`
	CreatedByRole   string     `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy       string     `gorm:"column:created_by;type:varchar(255);not null" json:"created_by"`
	RetryCount      int        `gorm:"column:retry_count;type:int;not null;default:0" json:"retry_count"`
	ElapsedTime     *float64   `gorm:"column:elapsed_time;type:float" json:"elapsed_time"`
	TotalTokens     *int       `gorm:"column:total_tokens;type:int" json:"total_tokens"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index" json:"created_at"`
	TriggeredAt     *time.Time `gorm:"column:triggered_at;type:timestamp" json:"triggered_at"`
	FinishedAt      *time.Time `gorm:"column:finished_at;type:timestamp" json:"finished_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowTriggerLog) TableName() string {
	return "workflow_trigger_logs"
}

// WorkflowWebhookTrigger [...]
type WorkflowWebhookTrigger struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID     string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	NodeID    string     `gorm:"column:node_id;type:varchar(64);not null" json:"node_id"`
	TenantID  string     `gorm:"column:tenant_id;type:varchar(36);not null;index" json:"tenant_id"`
	WebhookID string     `gorm:"column:webhook_id;type:varchar(24);not null;uniqueIndex" json:"webhook_id"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowWebhookTrigger) TableName() string {
	return "workflow_webhook_triggers"
}

// WorkflowPluginTrigger [...]
type WorkflowPluginTrigger struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	NodeID         string     `gorm:"column:node_id;type:varchar(64);not null" json:"node_id"`
	TenantID       string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ProviderID     string     `gorm:"column:provider_id;type:varchar(512);not null" json:"provider_id"`
	EventName      string     `gorm:"column:event_name;type:varchar(255);not null" json:"event_name"`
	SubscriptionID string     `gorm:"column:subscription_id;type:varchar(255);not null" json:"subscription_id"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowPluginTrigger) TableName() string {
	return "workflow_plugin_triggers"
}

// AppTrigger [...]
type AppTrigger struct {
	ID           string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID     string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID        string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	NodeID       *string    `gorm:"column:node_id;type:varchar(64);not null" json:"node_id"`
	TriggerType  string     `gorm:"column:trigger_type;type:varchar(50);not null" json:"trigger_type"`
	Title        string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	ProviderName *string    `gorm:"column:provider_name;type:varchar(255);default:''" json:"provider_name"`
	Status       string     `gorm:"column:status;type:varchar(50);not null;default:'enabled'" json:"status"`
	CreatedAt    *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (AppTrigger) TableName() string {
	return "app_triggers"
}

// WorkflowSchedulePlan [...]
type WorkflowSchedulePlan struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	NodeID         string     `gorm:"column:node_id;type:varchar(64);not null" json:"node_id"`
	TenantID       string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	CronExpression string     `gorm:"column:cron_expression;type:varchar(255);not null" json:"cron_expression"`
	Timezone       string     `gorm:"column:timezone;type:varchar(64);not null" json:"timezone"`
	NextRunAt      *time.Time `gorm:"column:next_run_at;type:timestamp;index" json:"next_run_at"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowSchedulePlan) TableName() string {
	return "workflow_schedule_plans"
}
