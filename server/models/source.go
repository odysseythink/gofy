package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// DataSourceApiKeyAuthBinding [...]
type DataSourceApiKeyAuthBinding struct {
	ID          string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID    string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Category    string     `gorm:"column:category;type:varchar(255);not null" json:"category"`
	Provider    string     `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	Credentials string     `gorm:"column:credentials;type:text" json:"credentials"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Disabled    bool       `gorm:"column:disabled;type:tinyint(1);default:0" json:"disabled"`
}

// TableName get sql table name.获取数据库表名
func (DataSourceApiKeyAuthBinding) TableName() string {
	return "data_source_api_key_auth_bindings"
}

func (self *DataSourceApiKeyAuthBinding) ToDict() map[string]any {
	res := map[string]any{
		"id":         self.ID,
		"tenant_id":  self.TenantID,
		"category":   self.Category,
		"provider":   self.Provider,
		"created_at": self.CreatedAt,
		"updated_at": self.UpdatedAt,
		"disabled":   self.Disabled,
	}
	tmp := map[string]any{}
	json.Unmarshal([]byte(self.Credentials), &tmp)
	res["credentials"] = tmp
	return res
}

// DataSourceOauthBinding [...]
type DataSourceOauthBinding struct {
	ID          string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID    string         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AccessToken string         `gorm:"column:access_token;type:varchar(255);not null" json:"access_token"`
	Provider    string         `gorm:"column:provider;type:varchar(255);not null" json:"provider"`
	SourceInfo  datatypes.JSON `gorm:"column:source_info;type:json;not null" json:"source_info"`
	CreatedAt   *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Disabled    bool           `gorm:"column:disabled;type:tinyint(1);default:0" json:"disabled"`
}

// TableName get sql table name.获取数据库表名
func (DataSourceOauthBinding) TableName() string {
	return "data_source_oauth_bindings"
}
