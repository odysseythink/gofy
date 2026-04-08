package services

import (
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"

	uuid "github.com/satori/go.uuid"
)

type WebappAuthService struct{}

func (s *WebappAuthService) GetOrCreateEndUser(appID, sessionID, userType string) (*models.EndUser, error) {
	var endUser models.EndUser
	err := dbengine.Instance().DB.Where("app_id = ? AND session_id = ? AND type = ?", appID, sessionID, userType).First(&endUser).Error
	if err == nil {
		return &endUser, nil
	}
	endUser = models.EndUser{
		ID:        uuid.NewV4().String(),
		AppID:     appID,
		SessionID: sessionID,
		Type:      userType,
	}
	if err := dbengine.Instance().DB.Create(&endUser).Error; err != nil {
		return nil, err
	}
	return &endUser, nil
}

func (s *WebappAuthService) GetEndUserByID(tenantID, appID, endUserID string) *models.EndUser {
	var endUser models.EndUser
	if err := dbengine.Instance().DB.Where("id = ? AND app_id = ? AND tenant_id = ?", endUserID, appID, tenantID).First(&endUser).Error; err != nil {
		return nil
	}
	return &endUser
}

func (s *WebappAuthService) GetEndUserBySessionID(appID, sessionID string) *models.EndUser {
	var endUser models.EndUser
	if err := dbengine.Instance().DB.Where("app_id = ? AND session_id = ?", appID, sessionID).First(&endUser).Error; err != nil {
		return nil
	}
	return &endUser
}

func (s *WebappAuthService) GetSiteByCode(appCode string) *models.Site {
	var site models.Site
	if err := dbengine.Instance().DB.Where("code = ?", appCode).First(&site).Error; err != nil {
		return nil
	}
	return &site
}

func (s *WebappAuthService) IsAppPublished(appCode string) (bool, string) {
	var site models.Site
	if err := dbengine.Instance().DB.Where("code = ?", appCode).First(&site).Error; err != nil {
		return false, ""
	}
	return site.Status == "normal", site.AppID
}

func (s *WebappAuthService) GetAppByCode(appCode string) (*models.App, error) {
	var site models.Site
	if err := dbengine.Instance().DB.Where("code = ?", appCode).First(&site).Error; err != nil {
		return nil, fmt.Errorf("site not found for code: %s", appCode)
	}
	var app models.App
	if err := dbengine.Instance().DB.Where("id = ?", site.AppID).First(&app).Error; err != nil {
		return nil, fmt.Errorf("app not found")
	}
	return &app, nil
}
