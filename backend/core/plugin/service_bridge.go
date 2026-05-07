package plugin

import (
	"iter"

	"github.com/odysseythink/gofy/backend/models"
)

// ServiceBridge abstracts the service-layer calls that backwards invocation needs.
// This breaks the import cycle between core/plugin and services.
type ServiceBridge interface {
	GetAppByIDAndTenantID(appID, tenantID string) (*models.App, error)
	GetAccountByID(accountID string) *models.Account
	AppGenerate(app *models.App, user *models.Account, args map[string]any, invokeFrom string, streaming bool) (map[string]any, iter.Seq[any])
	GetCustomCredentials(tenantID, providerName string) (map[string]any, error)
	UploadFile(tenantID, createdBy string, fileData []byte, filename, contentType string) (*models.UploadFile, error)
}

var globalServiceBridge ServiceBridge

// SetServiceBridge sets the global service bridge (called at init time by services package).
func SetServiceBridge(bridge ServiceBridge) {
	globalServiceBridge = bridge
}

// GetServiceBridge returns the global service bridge.
func GetServiceBridge() ServiceBridge {
	return globalServiceBridge
}
