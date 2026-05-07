package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
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
