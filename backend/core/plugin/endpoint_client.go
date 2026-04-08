package plugin

import (
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"

	uuid "github.com/satori/go.uuid"
)

// PluginEndpointClient manages custom endpoints from plugins.
type PluginEndpointClient struct{}

func NewPluginEndpointClient() *PluginEndpointClient {
	return &PluginEndpointClient{}
}

// CreateEndpoint creates a custom endpoint for a plugin.
func (ec *PluginEndpointClient) CreateEndpoint(tenantID string, pluginID string, name string, path string, method string, config map[string]any) (map[string]any, error) {
	// TODO: Create endpoint record in DB
	endpoint := map[string]any{
		"id":        uuid.NewV4().String(),
		"tenant_id": tenantID,
		"plugin_id": pluginID,
		"name":      name,
		"path":      path,
		"method":    method,
		"enabled":   true,
	}
	return endpoint, nil
}

// ListEndpoints lists all endpoints for a tenant.
func (ec *PluginEndpointClient) ListEndpoints(tenantID string) []map[string]any {
	// TODO: Query endpoint records
	return nil
}

// DeleteEndpoint removes a custom endpoint.
func (ec *PluginEndpointClient) DeleteEndpoint(tenantID string, endpointID string) error {
	// TODO: Delete endpoint record
	return nil
}

// EnableEndpoint enables or disables an endpoint.
func (ec *PluginEndpointClient) EnableEndpoint(tenantID string, endpointID string, enabled bool) error {
	return dbengine.Instance().DB.Model(&models.PluginInstallation{}).
		Where("tenant_id = ? AND id = ?", tenantID, endpointID).
		Update("endpoints_active", enabled).Error
}

// GetEndpoint retrieves a single endpoint by ID.
func (ec *PluginEndpointClient) GetEndpoint(tenantID string, endpointID string) (map[string]any, error) {
	var inst models.PluginInstallation
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND id = ?", tenantID, endpointID).First(&inst).Error; err != nil {
		return nil, fmt.Errorf("endpoint not found: %w", err)
	}
	return map[string]any{
		"id":                inst.ID,
		"tenant_id":         inst.TenantID,
		"plugin_id":         inst.PluginID,
		"unique_identifier": inst.PluginUniqueIdentifier,
		"endpoints_active":  inst.EndpointsActive,
	}, nil
}
