package services

import (
	"github.com/odysseythink/mlog"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type PluginService struct {
}

func (s *PluginService) FetchInstallTasks(tenant_id string, page int, page_size int) []*models.PluginInstallTask {
	limit := page_size
	offset := page_size * (page - 1)
	datas := []*models.PluginInstallTask{}
	err := dbengine.Instance().DB.Model(&models.PluginInstallTask{}).Where("tenant_id = ?", tenant_id).Order("created_at DESC").Offset(offset).Limit(limit).Scan(&datas).Error
	if err != nil {
		mlog.Errorf("get PluginInstallTask failed:%v", err)
		return nil
	}
	return datas
}
