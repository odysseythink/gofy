package services

import (
	"fmt"

	"mlib.com/gofy/server/core/plugin"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type PluginService struct{}

// InstallPlugin installs a plugin for a tenant.
func (s *PluginService) InstallPlugin(tenantID string, pluginID string, uniqueIdentifier string, source string) (*models.PluginInstallation, error) {
	mgr := plugin.Manager()
	if mgr == nil {
		return nil, fmt.Errorf("plugin manager not initialized")
	}
	return mgr.InstallPlugin(tenantID, pluginID, uniqueIdentifier, "local", source)
}

// UninstallPlugin removes a plugin for a tenant.
func (s *PluginService) UninstallPlugin(tenantID string, uniqueIdentifier string) error {
	mgr := plugin.Manager()
	if mgr == nil {
		return fmt.Errorf("plugin manager not initialized")
	}
	return mgr.UninstallPlugin(tenantID, uniqueIdentifier)
}

// ListPlugins lists installed plugins for a tenant with pagination.
func (s *PluginService) ListPlugins(tenantID string, page, pageSize int) ([]*models.PluginInstallation, int64) {
	mgr := plugin.Manager()
	if mgr == nil {
		return nil, 0
	}
	return mgr.ListPlugins(tenantID, page, pageSize)
}

// GetPlugin returns a specific plugin installation.
func (s *PluginService) GetPlugin(tenantID string, uniqueIdentifier string) *models.PluginInstallation {
	mgr := plugin.Manager()
	if mgr == nil {
		return nil
	}
	return mgr.GetPlugin(tenantID, uniqueIdentifier)
}

// GetPluginDeclaration returns the declaration/manifest for a plugin.
func (s *PluginService) GetPluginDeclaration(uniqueIdentifier string) *models.PluginDeclaration {
	mgr := plugin.Manager()
	if mgr == nil {
		return nil
	}
	return mgr.GetPluginDeclaration(uniqueIdentifier)
}

// UploadPluginPackage uploads a plugin package file.
func (s *PluginService) UploadPluginPackage(filename string, data []byte) (string, error) {
	mgr := plugin.Manager()
	if mgr == nil {
		return "", fmt.Errorf("plugin manager not initialized")
	}
	return mgr.UploadPluginPackage(filename, data)
}

// GetPluginAsset retrieves a plugin asset.
func (s *PluginService) GetPluginAsset(uniqueIdentifier string, assetPath string) ([]byte, error) {
	mgr := plugin.Manager()
	if mgr == nil {
		return nil, fmt.Errorf("plugin manager not initialized")
	}
	return mgr.GetPluginAsset(uniqueIdentifier, assetPath)
}

// FetchInstallTasks returns paginated install tasks for a tenant.
func (s *PluginService) FetchInstallTasks(tenantID string, page, pageSize int) ([]*models.PluginInstallTask, int64) {
	var tasks []*models.PluginInstallTask
	var total int64
	query := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	query.Model(&models.PluginInstallTask{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks)
	return tasks, total
}

// FetchPreferences returns plugin preferences for a tenant.
func (s *PluginService) FetchPreferences(tenantID string) map[string]any {
	// Get permission settings
	var perm models.TenantPluginPermission
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).First(&perm)

	// Get auto-upgrade strategy
	var strategy models.TenantPluginAutoUpgradeStrategy
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).First(&strategy)

	return map[string]any{
		"install_permission": perm.InstallPermission,
		"debug_permission":   perm.DebugPermission,
		"auto_upgrade":       strategy.StrategySetting,
	}
}

// GetEndpoints returns endpoints for a plugin.
func (s *PluginService) GetEndpoints(tenantID string, pluginID string) []map[string]any {
	var installations []*models.PluginInstallation
	dbengine.Instance().DB.Where("tenant_id = ? AND plugin_id = ?", tenantID, pluginID).Find(&installations)

	endpoints := make([]map[string]any, 0)
	for _, inst := range installations {
		endpoints = append(endpoints, map[string]any{
			"plugin_id":         inst.PluginID,
			"unique_identifier": inst.PluginUniqueIdentifier,
			"endpoints_active":  inst.EndpointsActive,
			"runtime_type":      inst.RuntimeType,
		})
	}
	return endpoints
}

// EnableEndpoints enables/disables endpoints for a plugin installation.
func (s *PluginService) EnableEndpoints(tenantID string, uniqueIdentifier string, active bool) error {
	return dbengine.Instance().DB.Model(&models.PluginInstallation{}).
		Where("tenant_id = ? AND plugin_unique_identifier = ?", tenantID, uniqueIdentifier).
		Update("endpoints_active", active).Error
}
