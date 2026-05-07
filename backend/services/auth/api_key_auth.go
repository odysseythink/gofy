package auth

import (
	"fmt"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
)

// APIKeyAuthService handles API key based authentication.
type APIKeyAuthService struct{}

// ValidateAPIKey validates an API key and returns the associated app and tenant.
func (s *APIKeyAuthService) ValidateAPIKey(apiKey string) (*models.ApiToken, error) {
	var token models.ApiToken
	if err := dbengine.Instance().DB.Where("token = ?", apiKey).First(&token).Error; err != nil {
		return nil, fmt.Errorf("invalid API key")
	}
	return &token, nil
}

// GetAPIKeysByApp returns all API keys for an app.
func (s *APIKeyAuthService) GetAPIKeysByApp(appID string) []*models.ApiToken {
	var tokens []*models.ApiToken
	dbengine.Instance().DB.Where("app_id = ?", appID).Find(&tokens)
	return tokens
}

// ValidateDatasetAPIKey validates an API key for dataset access.
func (s *APIKeyAuthService) ValidateDatasetAPIKey(apiKey string) (*models.ApiToken, error) {
	var token models.ApiToken
	if err := dbengine.Instance().DB.Where("token = ? AND type = ?", apiKey, "dataset").First(&token).Error; err != nil {
		return nil, fmt.Errorf("invalid dataset API key")
	}
	return &token, nil
}
