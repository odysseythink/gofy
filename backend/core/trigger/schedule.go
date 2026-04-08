package trigger

import (
	"fmt"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	uuid "github.com/satori/go.uuid"
)

// ScheduleProvider handles cron-based scheduled triggers.
type ScheduleProvider struct{}

func NewScheduleProvider() *ScheduleProvider {
	return &ScheduleProvider{}
}

func (sp *ScheduleProvider) Type() TriggerType { return TriggerTypeSchedule }
func (sp *ScheduleProvider) Setup(config map[string]any) error { return nil }
func (sp *ScheduleProvider) Teardown() error { return nil }
func (sp *ScheduleProvider) HandleEvent(event *TriggerEvent) error { return nil }

// CreateSchedulePlan creates a new scheduled trigger plan.
func (sp *ScheduleProvider) CreateSchedulePlan(appID, nodeID, tenantID, cronExpr, timezone string) (*models.WorkflowSchedulePlan, error) {
	plan := &models.WorkflowSchedulePlan{
		ID:             uuid.NewV4().String(),
		AppID:          appID,
		NodeID:         nodeID,
		TenantID:       tenantID,
		CronExpression: cronExpr,
		Timezone:       timezone,
	}

	// Calculate next run
	now := time.Now()
	plan.NextRunAt = &now // TODO: Use cron parser to calculate actual next run

	if err := dbengine.Instance().DB.Create(plan).Error; err != nil {
		return nil, fmt.Errorf("failed to create schedule plan: %w", err)
	}

	return plan, nil
}

// GetSchedulePlans returns schedule plans for an app.
func (sp *ScheduleProvider) GetSchedulePlans(appID, tenantID string) []*models.WorkflowSchedulePlan {
	var plans []*models.WorkflowSchedulePlan
	dbengine.Instance().DB.Where("app_id = ? AND tenant_id = ?", appID, tenantID).Find(&plans)
	return plans
}

// DeleteSchedulePlan removes a schedule plan.
func (sp *ScheduleProvider) DeleteSchedulePlan(appID, nodeID string) error {
	return dbengine.Instance().DB.Where("app_id = ? AND node_id = ?", appID, nodeID).Delete(&models.WorkflowSchedulePlan{}).Error
}

// GetDuePlans returns plans that should be executed.
func (sp *ScheduleProvider) GetDuePlans() []*models.WorkflowSchedulePlan {
	var plans []*models.WorkflowSchedulePlan
	now := time.Now()
	dbengine.Instance().DB.Where("next_run_at <= ?", now).Find(&plans)
	return plans
}
