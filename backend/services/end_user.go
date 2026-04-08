package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	uuid "github.com/satori/go.uuid"
)

type EndUserService struct{}

func (s *EndUserService) GetEndUserByID(tenantID, appID, endUserID string) *models.EndUser {
	var endUser models.EndUser
	if err := dbengine.Instance().DB.Where("id = ? AND app_id = ? AND tenant_id = ?", endUserID, appID, tenantID).First(&endUser).Error; err != nil {
		return nil
	}
	return &endUser
}

func (s *EndUserService) GetOrCreateEndUser(appModel *models.App, sessionID string) (*models.EndUser, error) {
	var endUser models.EndUser
	err := dbengine.Instance().DB.Where("app_id = ? AND session_id = ? AND type = ?", appModel.ID, sessionID, "browser").First(&endUser).Error
	if err == nil {
		return &endUser, nil
	}
	endUser = models.EndUser{
		ID:        uuid.NewV4().String(),
		TenantID:  appModel.TenantID,
		AppID:     appModel.ID,
		SessionID: sessionID,
		Type:      "browser",
	}
	if err := dbengine.Instance().DB.Create(&endUser).Error; err != nil {
		return nil, err
	}
	return &endUser, nil
}
