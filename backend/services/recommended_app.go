package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
)

type RecommendedAppService struct{}

func (s *RecommendedAppService) GetRecommendedApps(language string) []*models.RecommendedApp {
	var apps []*models.RecommendedApp
	query := dbengine.Instance().DB.Where("is_listed = ?", true)
	if language != "" {
		query = query.Where("language = ?", language)
	}
	query.Order("position ASC").Find(&apps)
	return apps
}

func (s *RecommendedAppService) GetRecommendedAppDetail(appID string) *models.RecommendedApp {
	var app models.RecommendedApp
	if err := dbengine.Instance().DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil
	}
	return &app
}

func (s *RecommendedAppService) GetInstalledApps(tenantID string) []*models.InstalledApp {
	var apps []*models.InstalledApp
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&apps)
	return apps
}
