package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type OperationService struct{}

func (s *OperationService) GetOperationLogs(tenantID string, page, limit int) ([]*models.OperationLog, int64) {
	var logs []*models.OperationLog
	var total int64
	query := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	query.Model(&models.OperationLog{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&logs)
	return logs, total
}
