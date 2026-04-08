package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type CreditPoolService struct{}

func (s *CreditPoolService) GetTenantCreditPool(tenantID string, poolType string) *models.TenantCreditPool {
	var pool models.TenantCreditPool
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND pool_type = ?", tenantID, poolType).First(&pool).Error; err != nil {
		return nil
	}
	return &pool
}

func (s *CreditPoolService) HasSufficientCredits(tenantID string, poolType string, required int64) bool {
	pool := s.GetTenantCreditPool(tenantID, poolType)
	if pool == nil {
		return false
	}
	return (pool.QuotaLimit - pool.QuotaUsed) >= required
}

func (s *CreditPoolService) DeductCredits(tenantID string, poolType string, amount int64) error {
	return dbengine.Instance().DB.Model(&models.TenantCreditPool{}).
		Where("tenant_id = ? AND pool_type = ?", tenantID, poolType).
		Update("quota_used", dbengine.Instance().DB.Raw("quota_used + ?", amount)).Error
}
