package models

import (
	"time"

	pluginidentityentities "github.com/odysseythink/gofy/backend/entities/plugin/identity"
	pluginenumtypes "github.com/odysseythink/gofy/backend/enum_types/plugin"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
)

// CeleryTask [...]
type CeleryTask struct {
	ID        int        `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	TaskID    string     `gorm:"column:task_id;type:varchar(155)" json:"task_id"`
	Status    string     `gorm:"column:status;type:varchar(50)" json:"status"`
	Result    []byte     `gorm:"column:result;type:blob" json:"result"`
	DateDone  *time.Time `gorm:"column:date_done;type:timestamp" json:"date_done"`
	Traceback string     `gorm:"column:traceback;type:text" json:"traceback"`
	Name      string     `gorm:"column:name;type:varchar(155)" json:"name"`
	Args      []byte     `gorm:"column:args;type:blob" json:"args"`
	Kwargs    []byte     `gorm:"column:kwargs;type:blob" json:"kwargs"`
	Worker    string     `gorm:"column:worker;type:varchar(155)" json:"worker"`
	Retries   int        `gorm:"column:retries;type:int" json:"retries"`
	Queue     string     `gorm:"column:queue;type:varchar(155)" json:"queue"`
}

// TableName get sql table name.获取数据库表名
func (CeleryTask) TableName() string {
	return "celery_taskmeta"
}

// CeleryTaskSet [...]
type CeleryTaskSet struct {
	ID        int        `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	TasksetID string     `gorm:"column:taskset_id;type:varchar(155)" json:"taskset_id"`
	Result    []byte     `gorm:"column:result;type:blob" json:"result"`
	DateDone  *time.Time `gorm:"column:date_done;type:timestamp" json:"date_done"`
}

// TableName get sql table name.获取数据库表名
func (CeleryTaskSet) TableName() string {
	return "celery_tasksetmeta"
}

type InstallTaskPluginStatus struct {
	PluginUniqueIdentifier pluginidentityentities.PluginUniqueIdentifier `json:"plugin_unique_identifier"`
	Labels                 commontypes.I18nObject                        `json:"labels"`
	Icon                   string                                        `json:"icon"`
	IconDark               string                                        `json:"icon_dark"`
	PluginID               string                                        `json:"plugin_id"`
	Status                 pluginenumtypes.PluginInstallTaskStatusType   `json:"status"`
	Message                string                                        `json:"message"`
}

type PluginInstallTask struct {
	ID               string                                      `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt        time.Time                                   `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt        time.Time                                   `json:"updated_at" gorm:"column:updated_at;not null"`
	Status           pluginenumtypes.PluginInstallTaskStatusType `json:"status" gorm:"column:status;not null"`
	TenantID         string                                      `json:"tenant_id" gorm:"column:tenant_id;type:uuid;not null"`
	TotalPlugins     int                                         `json:"total_plugins" gorm:"column:total_plugins;not null"`
	CompletedPlugins int                                         `json:"completed_plugins" gorm:"column:completed_plugins;not null"`
	Plugins          []InstallTaskPluginStatus                   `json:"plugins" gorm:"column:plugins;serializer:json"`
}

func (PluginInstallTask) TableName() string {
	return "plugin_install_tasks"
}
