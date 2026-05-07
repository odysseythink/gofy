package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	uuid "github.com/satori/go.uuid"
)

type ApiBasedExtensionService struct{}

func (s *ApiBasedExtensionService) GetAll(tenantID string) []*models.APIBasedExtension {
	var extensions []*models.APIBasedExtension
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&extensions)
	return extensions
}

func (s *ApiBasedExtensionService) GetByID(extensionID string) *models.APIBasedExtension {
	var ext models.APIBasedExtension
	if err := dbengine.Instance().DB.Where("id = ?", extensionID).First(&ext).Error; err != nil {
		return nil
	}
	return &ext
}

func (s *ApiBasedExtensionService) Create(tenantID, name, apiEndpoint, apiKey string) (*models.APIBasedExtension, error) {
	ext := &models.APIBasedExtension{
		ID:       uuid.NewV4().String(),
		TenantID: tenantID,
		Name:     name,
	}
	if err := dbengine.Instance().DB.Create(ext).Error; err != nil {
		return nil, err
	}
	return ext, nil
}

func (s *ApiBasedExtensionService) Delete(extensionID, tenantID string) error {
	return dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", extensionID, tenantID).Delete(&models.APIBasedExtension{}).Error
}

type CodeBasedExtensionService struct{}

func (s *CodeBasedExtensionService) GetCodeBasedExtensions() []map[string]any {
	return []map[string]any{
		{"name": "moderation", "label": map[string]string{"en_US": "Content Moderation", "zh_Hans": "内容审核"}},
		{"name": "external_data_tool", "label": map[string]string{"en_US": "External Data Tool", "zh_Hans": "外部数据工具"}},
	}
}
