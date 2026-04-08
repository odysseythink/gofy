package models

import "time"

// Plugin represents a plugin record in the database.
// It tracks plugin metadata and the number of installations (Refers).
type Plugin struct {
	ID                     string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	PluginUniqueIdentifier string     `gorm:"column:plugin_unique_identifier;type:varchar(255);not null;uniqueIndex" json:"plugin_unique_identifier"`
	PluginID               string     `gorm:"column:plugin_id;type:varchar(255);not null" json:"plugin_id"`
	Refers                 int        `gorm:"column:refers;type:int;not null;default:0" json:"refers"`
	InstallType            string     `gorm:"column:install_type;type:varchar(64);not null;default:'marketplace'" json:"install_type"`
	ManifestType           string     `gorm:"column:manifest_type;type:varchar(64)" json:"manifest_type"`
	Source                 string     `gorm:"column:source;type:text" json:"source"`
	CreatedAt              *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Plugin) TableName() string {
	return "plugins"
}

// PluginInstallation represents a per-tenant plugin installation record.
type PluginInstallation struct {
	ID                     string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID               string     `gorm:"column:tenant_id;type:varchar(36);not null;index:idx_tenant_plugin" json:"tenant_id"`
	PluginID               string     `gorm:"column:plugin_id;type:varchar(36);not null" json:"plugin_id"`
	PluginUniqueIdentifier string     `gorm:"column:plugin_unique_identifier;type:varchar(255);not null;index:idx_tenant_plugin" json:"plugin_unique_identifier"`
	RuntimeType            string     `gorm:"column:runtime_type;type:varchar(64);not null;default:'local'" json:"runtime_type"`
	EndpointsSetups        string     `gorm:"column:endpoints_setups;type:text" json:"endpoints_setups"`
	EndpointsActive        bool       `gorm:"column:endpoints_active;type:tinyint(1);not null;default:0" json:"endpoints_active"`
	Source                 string     `gorm:"column:source;type:text" json:"source"`
	Meta                   string     `gorm:"column:meta;type:text" json:"meta"`
	CreatedAt              *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PluginInstallation) TableName() string {
	return "plugin_installations"
}

// PluginDeclaration stores the cached declaration/manifest for a plugin.
type PluginDeclaration struct {
	ID                     string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	PluginUniqueIdentifier string     `gorm:"column:plugin_unique_identifier;type:varchar(255);not null;uniqueIndex" json:"plugin_unique_identifier"`
	Declaration            string     `gorm:"column:declaration;type:text" json:"declaration"`
	CreatedAt              *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PluginDeclaration) TableName() string {
	return "plugin_declarations"
}
