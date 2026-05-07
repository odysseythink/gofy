package plugin

import (
	"fmt"
	"sync"
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/storage"
	"github.com/odysseythink/mlog"

	uuid "github.com/satori/go.uuid"
)

// PluginManager manages the lifecycle of plugins including installation,
// uninstallation, runtime registration, and asset storage.
type PluginManager struct {
	mu       sync.RWMutex
	runtimes map[string]PluginLifetime // key: plugin unique identifier
}

var globalManager *PluginManager

// InitManager initializes the global plugin manager.
func InitManager() {
	globalManager = &PluginManager{
		runtimes: make(map[string]PluginLifetime),
	}
	mlog.Info("plugin manager initialized")
}

// Manager returns the global plugin manager.
func Manager() *PluginManager {
	return globalManager
}

// InstallPlugin installs a plugin for a tenant. It creates a Plugin record
// (or increments the reference count) and a PluginInstallation record.
func (pm *PluginManager) InstallPlugin(tenantID, pluginID, uniqueIdentifier, runtimeType, source string) (*models.PluginInstallation, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	db := dbengine.Instance().DB

	// Check if already installed for this tenant
	var existing models.PluginInstallation
	err := db.Where("tenant_id = ? AND plugin_unique_identifier = ?", tenantID, uniqueIdentifier).First(&existing).Error
	if err == nil {
		return &existing, fmt.Errorf("plugin %s already installed for tenant %s", uniqueIdentifier, tenantID)
	}

	// Ensure plugin record exists; create or increment reference count
	var pluginRecord models.Plugin
	err = db.Where("plugin_unique_identifier = ?", uniqueIdentifier).First(&pluginRecord).Error
	if err != nil {
		pluginRecord = models.Plugin{
			ID:                     uuid.NewV4().String(),
			PluginUniqueIdentifier: uniqueIdentifier,
			PluginID:               pluginID,
			Refers:                 1,
			InstallType:            "marketplace",
			Source:                 source,
		}
		if createErr := db.Create(&pluginRecord).Error; createErr != nil {
			return nil, fmt.Errorf("failed to create plugin record: %w", createErr)
		}
	} else {
		db.Model(&pluginRecord).Update("refers", pluginRecord.Refers+1)
	}

	// Create installation record
	installation := &models.PluginInstallation{
		ID:                     uuid.NewV4().String(),
		TenantID:               tenantID,
		PluginID:               pluginRecord.ID,
		PluginUniqueIdentifier: uniqueIdentifier,
		RuntimeType:            runtimeType,
		EndpointsActive:        true,
		Source:                 source,
	}
	if createErr := db.Create(installation).Error; createErr != nil {
		return nil, fmt.Errorf("failed to create installation: %w", createErr)
	}

	mlog.Infof("plugin %s installed for tenant %s", uniqueIdentifier, tenantID)
	return installation, nil
}

// UninstallPlugin removes a plugin installation for a tenant. When the last
// installation is removed the plugin record and its runtime are cleaned up.
func (pm *PluginManager) UninstallPlugin(tenantID, uniqueIdentifier string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	db := dbengine.Instance().DB

	// Find installation
	var installation models.PluginInstallation
	if err := db.Where("tenant_id = ? AND plugin_unique_identifier = ?", tenantID, uniqueIdentifier).First(&installation).Error; err != nil {
		return fmt.Errorf("installation not found for tenant %s plugin %s", tenantID, uniqueIdentifier)
	}

	// Delete installation
	if err := db.Delete(&installation).Error; err != nil {
		return fmt.Errorf("failed to delete installation: %w", err)
	}

	// Decrement reference count on the plugin record
	var pluginRecord models.Plugin
	if err := db.Where("plugin_unique_identifier = ?", uniqueIdentifier).First(&pluginRecord).Error; err == nil {
		pluginRecord.Refers--
		if pluginRecord.Refers <= 0 {
			// No more installations — remove plugin record and stop runtime
			db.Delete(&pluginRecord)
			pm.stopRuntime(uniqueIdentifier)
		} else {
			db.Save(&pluginRecord)
		}
	}

	mlog.Infof("plugin %s uninstalled from tenant %s", uniqueIdentifier, tenantID)
	return nil
}

// ListPlugins returns all installed plugins for a tenant with pagination.
func (pm *PluginManager) ListPlugins(tenantID string, page, pageSize int) ([]*models.PluginInstallation, int64) {
	var installations []*models.PluginInstallation
	var total int64

	query := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	query.Model(&models.PluginInstallation{}).Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&installations)

	return installations, total
}

// GetPlugin returns a specific plugin installation, or nil if not found.
func (pm *PluginManager) GetPlugin(tenantID, uniqueIdentifier string) *models.PluginInstallation {
	var installation models.PluginInstallation
	if err := dbengine.Instance().DB.Where("tenant_id = ? AND plugin_unique_identifier = ?", tenantID, uniqueIdentifier).First(&installation).Error; err != nil {
		return nil
	}
	return &installation
}

// GetPluginDeclaration retrieves the cached declaration for a plugin.
func (pm *PluginManager) GetPluginDeclaration(uniqueIdentifier string) *models.PluginDeclaration {
	var decl models.PluginDeclaration
	if err := dbengine.Instance().DB.Where("plugin_unique_identifier = ?", uniqueIdentifier).First(&decl).Error; err != nil {
		return nil
	}
	return &decl
}

// UploadPluginPackage stores a plugin package file in object storage
// and returns the storage key.
func (pm *PluginManager) UploadPluginPackage(filename string, data []byte) (string, error) {
	key := fmt.Sprintf("plugins/packages/%s/%s", time.Now().Format("20060102"), filename)
	if err := storage.Save(key, data); err != nil {
		return "", fmt.Errorf("failed to upload plugin package: %w", err)
	}
	return key, nil
}

// GetPluginAsset retrieves a plugin asset (icon, etc.) from object storage.
func (pm *PluginManager) GetPluginAsset(uniqueIdentifier, assetPath string) ([]byte, error) {
	key := fmt.Sprintf("plugins/assets/%s/%s", uniqueIdentifier, assetPath)
	return storage.LoadOnce(key)
}

// RegisterRuntime registers an active plugin runtime.
func (pm *PluginManager) RegisterRuntime(uniqueIdentifier string, runtime PluginLifetime) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.runtimes[uniqueIdentifier] = runtime
	mlog.Infof("plugin runtime registered: %s", uniqueIdentifier)
}

// GetRuntime returns the runtime for a plugin, if one is active.
func (pm *PluginManager) GetRuntime(uniqueIdentifier string) (PluginLifetime, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	rt, ok := pm.runtimes[uniqueIdentifier]
	return rt, ok
}

// stopRuntime stops and removes the runtime for a plugin. Caller must hold pm.mu.
func (pm *PluginManager) stopRuntime(uniqueIdentifier string) {
	if rt, ok := pm.runtimes[uniqueIdentifier]; ok {
		rt.Stop()
		delete(pm.runtimes, uniqueIdentifier)
		mlog.Infof("plugin runtime stopped: %s", uniqueIdentifier)
	}
}

// StopAll stops all active plugin runtimes.
func (pm *PluginManager) StopAll() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for id, rt := range pm.runtimes {
		rt.Stop()
		delete(pm.runtimes, id)
	}
	mlog.Info("all plugin runtimes stopped")
}
