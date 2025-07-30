package services

type ServiceGroup struct {
	PluginPermission  *PluginPermissionService
	PluginAutoUpgrade *PluginAutoUpgradeService
}

var ServiceGroupApp = ServiceGroup{
	PluginPermission:  &PluginPermissionService{},
	PluginAutoUpgrade: &PluginAutoUpgradeService{},
}
