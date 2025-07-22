package models

import "time"

// SavedMessage [...]
type SavedMessage struct {
	ID            string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID         string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App           *App       `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	MessageID     string     `gorm:"column:message_id;type:varchar(36);not null" json:"message_id"`
	Message       *Message   `json:"message" form:"message" gorm:"foreignKey:MessageID;references:ID;"`
	CreatedBy     string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedByRole string     `gorm:"column:created_by_role;type:varchar(255);default:end_user" json:"created_by_role"`
}

// TableName get sql table name.获取数据库表名
func (SavedMessage) TableName() string {
	return "saved_messages"
}

// PinnedConversation [...]
type PinnedConversation struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	ConversationID string     `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id"`
	CreatedBy      string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedByRole  string     `gorm:"column:created_by_role;type:varchar(255);default:end_user" json:"created_by_role"`
}

// TableName get sql table name.获取数据库表名
func (PinnedConversation) TableName() string {
	return "pinned_conversations"
}
