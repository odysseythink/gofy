package trigger_schedule

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
)

// TriggerScheduleNodeData represents the schedule trigger node configuration
type TriggerScheduleNodeData struct {
	*basenodesentities.BaseNodeData
	Mode           string        `json:"mode"`
	Frequency      string        `json:"frequency,omitempty"`
	CronExpression string        `json:"cron_expression,omitempty"`
	VisualConfig   *VisualConfig `json:"visual_config,omitempty"`
	Timezone       string        `json:"timezone"`
}

// VisualConfig represents visual configuration for schedule trigger
type VisualConfig struct {
	OnMinute    int      `json:"on_minute"`
	Time        string   `json:"time"`
	Weekdays    []string `json:"weekdays,omitempty"`
	MonthlyDays []any    `json:"monthly_days,omitempty"`
}

func New() *TriggerScheduleNodeData {
	return &TriggerScheduleNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
		Mode:         "visual",
		Timezone:     "UTC",
	}
}
