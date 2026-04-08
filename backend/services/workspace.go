package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type WorkspaceService struct{}

func (s *WorkspaceService) GetWorkspaces(accountID string) []*models.Tenant {
	var tenants []*models.Tenant
	dbengine.Instance().DB.Joins("JOIN tenant_account_joins ON tenant_account_joins.tenant_id = tenants.id").
		Where("tenant_account_joins.account_id = ?", accountID).Find(&tenants)
	return tenants
}
