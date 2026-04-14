package services

import (
	"fmt"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/odysseythink/mlog"
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

// === Marketplace Integration ===

// SearchMarketplacePlugins searches for plugins in the marketplace.
func (s *PluginService) SearchMarketplacePlugins(query string, category string, page, pageSize int) ([]map[string]any, int64, error) {
	// TODO: Call marketplace API (env.VITE_MARKETPLACE_API_PREFIX)
	return []map[string]any{}, 0, nil
}

// GetMarketplacePluginDetail returns plugin details from marketplace.
func (s *PluginService) GetMarketplacePluginDetail(pluginID string) (map[string]any, error) {
	// TODO: Call marketplace API
	return nil, fmt.Errorf("marketplace integration not yet implemented")
}

// InstallFromMarketplace installs a plugin from the marketplace.
func (s *PluginService) InstallFromMarketplace(tenantID, pluginID, version string) (*models.PluginInstallation, error) {
	// TODO: Download from marketplace, then install
	return s.InstallPlugin(tenantID, pluginID, fmt.Sprintf("%s@%s", pluginID, version), "marketplace")
}

// === Plugin Migration ===

// MigratePluginData migrates plugin data between versions.
func (s *PluginService) MigratePluginData(tenantID, uniqueIdentifier, fromVersion, toVersion string) error {
	mlog.Infof("migrating plugin %s from %s to %s for tenant %s", uniqueIdentifier, fromVersion, toVersion, tenantID)
	// TODO: Run migration scripts from plugin package
	return nil
}

// CheckPluginUpgrade checks if a plugin has available upgrades.
func (s *PluginService) CheckPluginUpgrade(tenantID, uniqueIdentifier string) (string, bool, error) {
	installation := plugin.Manager().GetPlugin(tenantID, uniqueIdentifier)
	if installation == nil {
		return "", false, fmt.Errorf("plugin not installed")
	}
	// TODO: Check marketplace for newer version
	return "", false, nil
}

// AutoUpgradePlugins upgrades all plugins for a tenant based on strategy.
func (s *PluginService) AutoUpgradePlugins(tenantID string) (int, error) {
	var strategy models.TenantPluginAutoUpgradeStrategy
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).First(&strategy)

	if strategy.StrategySetting == "disabled" || strategy.StrategySetting == "" {
		return 0, nil
	}

	installations, _ := plugin.Manager().ListPlugins(tenantID, 1, 1000)
	upgraded := 0
	for _, inst := range installations {
		newVersion, hasUpgrade, _ := s.CheckPluginUpgrade(tenantID, inst.PluginUniqueIdentifier)
		if hasUpgrade {
			if strategy.StrategySetting == "all" {
				s.UpgradePlugin(tenantID, inst.PluginUniqueIdentifier, newVersion)
				upgraded++
			}
		}
	}
	return upgraded, nil
}

// UpgradePlugin upgrades a specific plugin.
func (s *PluginService) UpgradePlugin(tenantID, uniqueIdentifier, newVersion string) error {
	mlog.Infof("upgrading plugin %s to %s", uniqueIdentifier, newVersion)
	// Uninstall old, install new
	if err := plugin.Manager().UninstallPlugin(tenantID, uniqueIdentifier); err != nil {
		return fmt.Errorf("failed to uninstall old version: %w", err)
	}
	newID := fmt.Sprintf("%s@%s", uniqueIdentifier, newVersion)
	_, err := plugin.Manager().InstallPlugin(tenantID, uniqueIdentifier, newID, "local", "upgrade")
	return err
}

// === Debugging ===

// GetDebuggingKey generates a debugging key for plugin development.
func (s *PluginService) GetDebuggingKey(tenantID string) (string, error) {
	var perm models.TenantPluginPermission
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).First(&perm)
	if perm.DebugPermission != "everyone" && perm.DebugPermission != "admins" {
		return "", fmt.Errorf("debugging not enabled for this tenant")
	}
	// Generate a temporary debugging token
	key := uuid.NewV4().String()
	// TODO: Store in Redis with TTL
	return key, nil
}

// === Advanced Installation ===

// InstallFromPackage installs a plugin from an uploaded package file.
func (s *PluginService) InstallFromPackage(tenantID string, packageData []byte, filename string) (*models.PluginInstallation, error) {
	// Upload package
	_, err := plugin.Manager().UploadPluginPackage(filename, packageData)
	if err != nil {
		return nil, err
	}

	// TODO: Extract package, parse manifest, validate
	pluginID := strings.TrimSuffix(filename, filepath.Ext(filename))
	uniqueIdentifier := fmt.Sprintf("local/%s", pluginID)

	return s.InstallPlugin(tenantID, pluginID, uniqueIdentifier, "upload")
}

// GetPluginLogs returns recent logs for a plugin runtime.
func (s *PluginService) GetPluginLogs(tenantID, uniqueIdentifier string, lines int) ([]string, error) {
	runtime, ok := plugin.Manager().GetRuntime(uniqueIdentifier)
	if !ok {
		return nil, fmt.Errorf("plugin runtime not active")
	}
	// TODO: Implement log buffering in runtime
	_ = runtime
	return []string{}, nil
}

// GetPluginHealth checks if a plugin runtime is healthy.
func (s *PluginService) GetPluginHealth(tenantID, uniqueIdentifier string) map[string]any {
	runtime, ok := plugin.Manager().GetRuntime(uniqueIdentifier)
	if !ok {
		return map[string]any{"status": "not_running"}
	}
	if lr, ok := runtime.(*plugin.LocalRuntime); ok {
		return map[string]any{
			"status":    lr.RuntimeState().Status,
			"restarts":  lr.RuntimeState().Restarts,
			"active_at": lr.RuntimeState().ActiveAt,
		}
	}
	return map[string]any{"status": "unknown"}
}

// BatchInstallPlugins installs multiple plugins.
func (s *PluginService) BatchInstallPlugins(tenantID string, plugins []map[string]string) ([]*models.PluginInstallation, []error) {
	var installations []*models.PluginInstallation
	var errors []error
	for _, p := range plugins {
		pluginID := p["plugin_id"]
		uniqueID := p["unique_identifier"]
		source := p["source"]
		if source == "" {
			source = "marketplace"
		}
		inst, err := s.InstallPlugin(tenantID, pluginID, uniqueID, source)
		installations = append(installations, inst)
		errors = append(errors, err)
	}
	return installations, errors
}
