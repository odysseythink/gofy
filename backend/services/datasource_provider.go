package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"

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

// === Credential Encryption/Decryption ===

// EncryptCredentials encrypts sensitive credential values before storage.
func (s *DatasourceProviderService) EncryptCredentials(credentials map[string]any, secretKeys []string) map[string]any {
	encrypted := make(map[string]any)
	for k, v := range credentials {
		encrypted[k] = v
	}
	for _, key := range secretKeys {
		if val, ok := encrypted[key].(string); ok && val != "" {
			// Simple obfuscation — in production use AES encryption
			encrypted[key] = base64.StdEncoding.EncodeToString([]byte(val))
		}
	}
	return encrypted
}

// DecryptCredentials decrypts stored credential values.
func (s *DatasourceProviderService) DecryptCredentials(credentials map[string]any, secretKeys []string) map[string]any {
	decrypted := make(map[string]any)
	for k, v := range credentials {
		decrypted[k] = v
	}
	for _, key := range secretKeys {
		if val, ok := decrypted[key].(string); ok && val != "" {
			decoded, err := base64.StdEncoding.DecodeString(val)
			if err == nil {
				decrypted[key] = string(decoded)
			}
		}
	}
	return decrypted
}

// ObfuscateCredentials masks sensitive values for display.
func (s *DatasourceProviderService) ObfuscateCredentials(credentials map[string]any, secretKeys []string) map[string]any {
	obfuscated := make(map[string]any)
	for k, v := range credentials {
		obfuscated[k] = v
	}
	for _, key := range secretKeys {
		if val, ok := obfuscated[key].(string); ok && val != "" {
			if len(val) > 6 {
				obfuscated[key] = val[:3] + "***" + val[len(val)-3:]
			} else {
				obfuscated[key] = "***"
			}
		}
	}
	return obfuscated
}

// === OAuth Management ===

// AddDatasourceOAuthProvider adds an OAuth-based datasource provider credential.
func (s *DatasourceProviderService) AddDatasourceOAuthProvider(
	tenantID, provider, pluginID string,
	credentials map[string]any,
	name string,
) (*models.DataSourceOauthBinding, error) {
	accessToken, _ := credentials["access_token"].(string)
	sourceInfo, _ := json.Marshal(credentials)

	binding := &models.DataSourceOauthBinding{
		ID:          uuid.NewV4().String(),
		TenantID:    tenantID,
		Provider:    provider,
		AccessToken: accessToken,
		SourceInfo:  sourceInfo,
	}
	if err := dbengine.Instance().DB.Create(binding).Error; err != nil {
		return nil, err
	}
	return binding, nil
}

// ReauthorizeDatasourceOAuthProvider reauthorizes an existing OAuth credential.
func (s *DatasourceProviderService) ReauthorizeDatasourceOAuthProvider(
	tenantID, authID, provider, pluginID string,
	credentials map[string]any,
) error {
	accessToken, _ := credentials["access_token"].(string)
	sourceInfo, _ := json.Marshal(credentials)

	return dbengine.Instance().DB.Model(&models.DataSourceOauthBinding{}).
		Where("id = ? AND tenant_id = ?", authID, tenantID).
		Updates(map[string]any{
			"access_token": accessToken,
			"source_info":  string(sourceInfo),
		}).Error
}

// === Credential Listing ===

// ListDatasourceCredentials returns all credentials with obfuscated secrets.
func (s *DatasourceProviderService) ListDatasourceCredentials(tenantID, provider, pluginID string) []map[string]any {
	var apiBindings []*models.DataSourceApiKeyAuthBinding
	query := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	query.Find(&apiBindings)

	var oauthBindings []*models.DataSourceOauthBinding
	query2 := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	query2.Find(&oauthBindings)

	results := make([]map[string]any, 0)
	for _, b := range apiBindings {
		results = append(results, map[string]any{
			"id":         b.ID,
			"type":       "api_key",
			"provider":   b.Provider,
			"category":   b.Category,
			"created_at": b.CreatedAt,
		})
	}
	for _, b := range oauthBindings {
		results = append(results, map[string]any{
			"id":         b.ID,
			"type":       "oauth",
			"provider":   b.Provider,
			"created_at": b.CreatedAt,
		})
	}
	return results
}

// GetRealDatasourceCredentials returns unobfuscated credentials (for internal use).
func (s *DatasourceProviderService) GetRealDatasourceCredentials(tenantID, provider, pluginID string) []map[string]any {
	var apiBindings []*models.DataSourceApiKeyAuthBinding
	query := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	query.Find(&apiBindings)

	results := make([]map[string]any, 0)
	for _, b := range apiBindings {
		results = append(results, map[string]any{
			"id":          b.ID,
			"type":        "api_key",
			"provider":    b.Provider,
			"category":    b.Category,
			"credentials": b.Credentials,
			"created_at":  b.CreatedAt,
		})
	}
	return results
}

// UpdateDatasourceCredentials updates existing datasource credentials.
func (s *DatasourceProviderService) UpdateDatasourceCredentials(
	tenantID, authID, provider string,
	credentials map[string]any,
) error {
	updates := map[string]any{}
	if credentials != nil {
		credJSON, _ := json.Marshal(credentials)
		updates["credentials"] = string(credJSON)
	}
	if len(updates) == 0 {
		return nil
	}
	return dbengine.Instance().DB.Model(&models.DataSourceApiKeyAuthBinding{}).
		Where("id = ? AND tenant_id = ?", authID, tenantID).
		Updates(updates).Error
}

// ExtractSecretVariables returns the names of secret fields for a datasource provider.
func (s *DatasourceProviderService) ExtractSecretVariables(tenantID, providerID, credentialType string) []string {
	// TODO: Query plugin system for datasource provider schema
	// Default secret fields for common providers
	switch credentialType {
	case "api_key":
		return []string{"api_key", "secret_key", "token", "password"}
	case "oauth2":
		return []string{"client_secret", "access_token", "refresh_token"}
	default:
		return []string{"api_key", "secret_key"}
	}
}

// SetDefaultDatasourceProvider sets a datasource provider as the default for a tenant.
func (s *DatasourceProviderService) SetDefaultDatasourceProvider(tenantID, providerID, authID string) error {
	// TODO: Implement default provider preference storage
	return nil
}

// GenerateNextProviderName generates a sequential name for a new datasource provider.
func (s *DatasourceProviderService) GenerateNextProviderName(tenantID, provider, baseName string) string {
	var count int64
	dbengine.Instance().DB.Model(&models.DataSourceApiKeyAuthBinding{}).
		Where("tenant_id = ? AND provider = ?", tenantID, provider).Count(&count)
	if count == 0 {
		return baseName
	}
	return fmt.Sprintf("%s (%d)", baseName, count+1)
}
