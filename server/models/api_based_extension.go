package models

import "time"

// APIBasedExtensions [...]
type APIBasedExtensions struct {
	ID          string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID    string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	APIEndpoint string     `gorm:"column:api_endpoint;type:varchar(255);not null" json:"api_endpoint"`
	APIKey      string     `gorm:"column:api_key;type:text;not null" json:"api_key"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (APIBasedExtensions) TableName() string {
	return "api_based_extensions"
}
