package models

import "time"

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
