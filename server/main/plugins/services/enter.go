package services

type ServiceGroup struct {
	PluginPermission  *PluginPermissionService
	PluginAutoUpgrade *PluginAutoUpgradeService
	Plugin            *PluginService
}

var ServiceGroupApp = ServiceGroup{
	PluginPermission:  &PluginPermissionService{},
	PluginAutoUpgrade: &PluginAutoUpgradeService{},
	Plugin:            &PluginService{},
}
