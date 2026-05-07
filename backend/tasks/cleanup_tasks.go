package tasks

import (
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

// CleanOldMessages removes messages older than 30 days for conversations that are deleted.
func CleanOldMessages() {
	cutoff := time.Now().AddDate(0, 0, -30)
	result := dbengine.Instance().DB.Where("created_at < ? AND conversation_id NOT IN (SELECT id FROM conversations)", cutoff).Delete(&models.Message{})
	if result.Error != nil {
		mlog.Errorf("clean messages failed: %v", result.Error)
		return
	}
	if result.RowsAffected > 0 {
		mlog.Infof("cleaned %d orphaned messages", result.RowsAffected)
	}
}

// CleanOldWorkflowRuns removes workflow runs older than 30 days.
func CleanOldWorkflowRuns() {
	cutoff := time.Now().AddDate(0, 0, -30)

	// Delete old node executions first
	dbengine.Instance().DB.Where("created_at < ?", cutoff).Delete(&models.WorkflowNodeExecution{})

	// Delete old workflow runs
	result := dbengine.Instance().DB.Where("created_at < ?", cutoff).Delete(&models.WorkflowRun{})
	if result.Error != nil {
		mlog.Errorf("clean workflow runs failed: %v", result.Error)
		return
	}
	if result.RowsAffected > 0 {
		mlog.Infof("cleaned %d old workflow runs", result.RowsAffected)
	}
}

// CleanUnusedDatasets archives datasets with no documents and no app bindings.
func CleanUnusedDatasets() {
	// Find datasets with no documents and no app bindings
	var datasets []models.Dataset
	dbengine.Instance().DB.Where(
		"id NOT IN (SELECT DISTINCT dataset_id FROM documents) AND id NOT IN (SELECT DISTINCT dataset_id FROM app_dataset_joins)",
	).Where("created_at < ?", time.Now().AddDate(0, 0, -90)).Find(&datasets)

	for _, ds := range datasets {
		mlog.Infof("marking unused dataset %s for review", ds.ID)
		// Don't auto-delete, just log for review
	}
}
