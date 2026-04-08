package services

import (
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type DatasourceProviderService struct{}

func (s *DatasourceProviderService) GetDatasourceCredentials(tenantID string, provider string) (*models.DataSourceApiKeyAuthBinding, error) {
	var binding models.DataSourceApiKeyAuthBinding
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND provider = ?", tenantID, provider).First(&binding).Error; err != nil {
		return nil, fmt.Errorf("credentials not found")
	}
	return &binding, nil
}

func (s *DatasourceProviderService) GetAllDatasourceCredentials(tenantID string) []*models.DataSourceApiKeyAuthBinding {
	var bindings []*models.DataSourceApiKeyAuthBinding
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&bindings)
	return bindings
}

func (s *DatasourceProviderService) AddDatasourceApiKeyProvider(tenantID string, provider string, category string, credentials string) (*models.DataSourceApiKeyAuthBinding, error) {
	binding := &models.DataSourceApiKeyAuthBinding{
		TenantID:    tenantID,
		Provider:    provider,
		Category:    category,
		Credentials: credentials,
	}
	if err := dbengine.Instance().DB.Create(binding).Error; err != nil {
		return nil, err
	}
	return binding, nil
}

func (s *DatasourceProviderService) RemoveDatasourceCredentials(tenantID string, authID string) error {
	return dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", authID, tenantID).Delete(&models.DataSourceApiKeyAuthBinding{}).Error
}

func (s *DatasourceProviderService) GetOAuthBindings(tenantID string) []*models.DataSourceOauthBinding {
	var bindings []*models.DataSourceOauthBinding
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&bindings)
	return bindings
}
