package models

import "time"

// HumanInputForm [...]
type HumanInputForm struct {
	ID                     string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt              *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	TenantID               string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID                  string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	WorkflowRunID          *string    `gorm:"column:workflow_run_id;type:varchar(36)" json:"workflow_run_id"`
	FormKind               string     `gorm:"column:form_kind;type:varchar(50);not null;default:'runtime'" json:"form_kind"`
	NodeID                 string     `gorm:"column:node_id;type:varchar(60);not null" json:"node_id"`
	FormDefinition         string     `gorm:"column:form_definition;type:text;not null" json:"form_definition"`
	RenderedContent        string     `gorm:"column:rendered_content;type:text;not null" json:"rendered_content"`
	Status                 string     `gorm:"column:status;type:varchar(50);not null;default:'waiting'" json:"status"`
	ExpirationTime         time.Time  `gorm:"column:expiration_time;type:timestamp;not null" json:"expiration_time"`
	SelectedActionID       *string    `gorm:"column:selected_action_id;type:varchar(200)" json:"selected_action_id"`
	SubmittedData          *string    `gorm:"column:submitted_data;type:text" json:"submitted_data"`
	SubmittedAt            *time.Time `gorm:"column:submitted_at;type:timestamp" json:"submitted_at"`
	SubmissionUserID       *string    `gorm:"column:submission_user_id;type:varchar(36)" json:"submission_user_id"`
	SubmissionEndUserID    *string    `gorm:"column:submission_end_user_id;type:varchar(36)" json:"submission_end_user_id"`
	CompletedByRecipientID *string    `gorm:"column:completed_by_recipient_id;type:varchar(36)" json:"completed_by_recipient_id"`
}

// TableName get sql table name.获取数据库表名
func (HumanInputForm) TableName() string {
	return "human_input_forms"
}

// HumanInputDelivery [...]
type HumanInputDelivery struct {
	ID                 string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt          *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	FormID             string     `gorm:"column:form_id;type:varchar(36);not null;index" json:"form_id"`
	DeliveryMethodType string     `gorm:"column:delivery_method_type;type:varchar(50);not null" json:"delivery_method_type"`
	DeliveryConfigID   *string    `gorm:"column:delivery_config_id;type:varchar(36)" json:"delivery_config_id"`
	ChannelPayload     string     `gorm:"column:channel_payload;type:text;not null" json:"channel_payload"`
}

// TableName get sql table name.获取数据库表名
func (HumanInputDelivery) TableName() string {
	return "human_input_form_deliveries"
}

// HumanInputFormRecipient [...]
type HumanInputFormRecipient struct {
	ID               string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	FormID           string     `gorm:"column:form_id;type:varchar(36);not null;index" json:"form_id"`
	DeliveryID       string     `gorm:"column:delivery_id;type:varchar(36);not null;index" json:"delivery_id"`
	RecipientType    string     `gorm:"column:recipient_type;type:varchar(50);not null" json:"recipient_type"`
	RecipientPayload string     `gorm:"column:recipient_payload;type:text;not null" json:"recipient_payload"`
	AccessToken      *string    `gorm:"column:access_token;type:varchar(32);not null;uniqueIndex" json:"access_token"`
}

// TableName get sql table name.获取数据库表名
func (HumanInputFormRecipient) TableName() string {
	return "human_input_form_recipients"
}
