package services

import (
	"encoding/json"
	"slices"

	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type PluginAutoUpgradeService struct {
}

func (s *PluginAutoUpgradeService) GetStrategy(tenant_id string) *models.TenantPluginAutoUpgradeStrategy {
	strategy := new(models.TenantPluginAutoUpgradeStrategy)
	err := dbengine.Instance().DB.Model(&models.TenantPluginAutoUpgradeStrategy{}).Where("tenant_id = ?", tenant_id).First(strategy).Error
	if err != nil {
		mlog.Errorf("get TenantPluginAutoUpgradeStrategy failed:%v", err)
		return nil
	}
	return strategy
}

func (s *PluginAutoUpgradeService) ChangeStrategy(
	tenant_id string,
	strategy_setting models.TenantPluginAutoUpgradeStrategySettingType,
	upgrade_time_of_day int,
	upgrade_mode models.TenantPluginAutoUpgradeStrategyUpgradeModeType,
	exclude_plugins []string,
	include_plugins []string,
) bool {
	exist_strategy := new(models.TenantPluginAutoUpgradeStrategy)
	err := dbengine.Instance().DB.Model(&models.TenantPluginAutoUpgradeStrategy{}).Where("tenant_id = ?", tenant_id).First(exist_strategy).Error
	if err != nil {
		mlog.Errorf("get TenantPluginAutoUpgradeStrategy failed:%v", err)
		exist_strategy = nil
	}

	if exist_strategy == nil {
		strategy := &models.TenantPluginAutoUpgradeStrategy{
			Model: models.Model{
				ID: uuid.NewV4().String(),
			},
			TenantID:         tenant_id,
			StrategySetting:  strategy_setting,
			UpgradeTimeOfDay: upgrade_time_of_day,
			UpgradeMode:      upgrade_mode,
		}
		if exclude_plugins == nil {
			exclude_plugins = make([]string, 0)
		}
		bindata, _ := json.Marshal(exclude_plugins)
		strategy.ExcludePlugins.Scan(bindata)
		if include_plugins == nil {
			include_plugins = make([]string, 0)
		}
		bindata, _ = json.Marshal(include_plugins)
		strategy.IncludePlugins.Scan(bindata)
		dbengine.Instance().DB.Create(strategy)
	} else {
		update_strategy := &models.TenantPluginAutoUpgradeStrategy{
			Model: models.Model{
				ID: exist_strategy.ID,
			},
			StrategySetting:  strategy_setting,
			UpgradeTimeOfDay: upgrade_time_of_day,
		}
		if exclude_plugins == nil {
			exclude_plugins = make([]string, 0)
		}
		bindata, _ := json.Marshal(exclude_plugins)
		update_strategy.ExcludePlugins.Scan(bindata)
		if include_plugins == nil {
			include_plugins = make([]string, 0)
		}
		bindata, _ = json.Marshal(include_plugins)
		update_strategy.IncludePlugins.Scan(bindata)
		dbengine.Instance().DB.Updates(update_strategy)
	}
	return true

}

func (s *PluginAutoUpgradeService) ExcludePlugin(tenant_id string, plugin_id string) bool {
	exist_strategy := new(models.TenantPluginAutoUpgradeStrategy)
	err := dbengine.Instance().DB.Model(&models.TenantPluginAutoUpgradeStrategy{}).Where("tenant_id = ?", tenant_id).First(exist_strategy).Error
	if err != nil {
		mlog.Errorf("get TenantPluginAutoUpgradeStrategy failed:%v", err)
		exist_strategy = nil
	}

	if exist_strategy == nil {
		// create for this tenant
		s.ChangeStrategy(
			tenant_id,
			models.TenantPluginAutoUpgradeStrategySetting_FIX_ONLY,
			0,
			models.TenantPluginAutoUpgradeStrategyUpgradeMode_EXCLUDE,
			[]string{plugin_id},
			[]string{},
		)
		return true
	} else {
		if err := exist_strategy.ExcludePlugins.Bind(&exist_strategy.ExcludePluginList); err != nil {
			mlog.Warningf("bind ExcludePlugins to list failed:%v", err)
		}
		if err := exist_strategy.IncludePlugins.Bind(&exist_strategy.IncludePluginList); err != nil {
			mlog.Warningf("bind IncludePlugins to list failed:%v", err)
		}
		switch exist_strategy.UpgradeMode {
		case models.TenantPluginAutoUpgradeStrategyUpgradeMode_EXCLUDE:
			if !slices.Contains(exist_strategy.ExcludePluginList, plugin_id) {
				if exist_strategy.ExcludePluginList == nil {
					exist_strategy.ExcludePluginList = make([]string, 0)
				}
				exist_strategy.ExcludePluginList = append(exist_strategy.ExcludePluginList, plugin_id)
				update_strategy := &models.TenantPluginAutoUpgradeStrategy{
					Model: models.Model{
						ID: exist_strategy.ID,
					},
				}
				bindata, _ := json.Marshal(exist_strategy.ExcludePluginList)
				update_strategy.ExcludePlugins.Scan(bindata)
				dbengine.Instance().DB.Updates(update_strategy)
			}
		case models.TenantPluginAutoUpgradeStrategyUpgradeMode_PARTIAL:
			if slices.Contains(exist_strategy.IncludePluginList, plugin_id) {
				update_strategy := &models.TenantPluginAutoUpgradeStrategy{
					Model: models.Model{
						ID: exist_strategy.ID,
					},
				}
				exist_strategy.IncludePluginList = slices.DeleteFunc(exist_strategy.IncludePluginList, func(val string) bool { return val == plugin_id })
				bindata, _ := json.Marshal(exist_strategy.IncludePluginList)
				update_strategy.IncludePlugins.Scan(bindata)
				dbengine.Instance().DB.Updates(update_strategy)
			}
		case models.TenantPluginAutoUpgradeStrategyUpgradeMode_ALL:
			update_strategy := &models.TenantPluginAutoUpgradeStrategy{
				Model: models.Model{
					ID: exist_strategy.ID,
				},
				UpgradeMode: models.TenantPluginAutoUpgradeStrategyUpgradeMode_EXCLUDE,
			}
			bindata, _ := json.Marshal([]string{plugin_id})
			update_strategy.ExcludePlugins.Scan(bindata)
			dbengine.Instance().DB.Updates(update_strategy)
		}
		return true
	}
}
