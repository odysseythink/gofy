package services

import (
	pluginentities "mlib.com/gofy/server/entities/plugin"
)

type PluginService struct {
}

func (s *PluginService) FetchInstallTasks(tenant_id string, page int, page_size int) []*pluginentities.PluginInstallTask {
	manager = PluginInstaller()
	return manager.fetch_plugin_installation_tasks(tenant_id, page, page_size)
}
