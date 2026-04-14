package services

import (
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type PluginPermissionService struct {
}

func (s *PluginPermissionService) GetPermission(tenant_id string) *models.TenantPluginPermission {
	permission := new(models.TenantPluginPermission)
	db := dbengine.Instance().DB.Model(&models.TenantPluginPermission{})
	err := db.Where("tenant_id = ?", tenant_id).First(permission).Error
	if err != nil {
		mlog.Errorf("get TenantPluginPermission failed:%v", err)
		return nil
	}
	return permission
}

func (s *PluginPermissionService) ChangePermission(
	tenant_id string,
	install_permission models.TenantPluginInstallPermissionType,
	debug_permission models.TenantPluginDebugPermissionType,
) bool {
	permission := new(models.TenantPluginPermission)
	err := dbengine.Instance().DB.Model(&models.TenantPluginPermission{}).Where("tenant_id = ?", tenant_id).First(permission).Error
	if err != nil {
		mlog.Errorf("get TenantPluginPermission failed:%v", err)
		permission = nil
	}
	if permission == nil {
		permission = &models.TenantPluginPermission{
			ID:                uuid.NewV4().String(),
			TenantID:          tenant_id,
			InstallPermission: install_permission,
			DebugPermission:   debug_permission,
		}
		dbengine.Instance().DB.Create(permission)
	} else {
		dbengine.Instance().DB.Updates(&models.TenantPluginPermission{ID: permission.ID, InstallPermission: install_permission, DebugPermission: debug_permission})
	}
	return true
}
