package modelmanager

import (
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/cache"
	providerentities "mlib.com/gofy/server/entities/provider"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type LBModelManager struct {
	TenantID string
	Provider string
	modelruntimeenumtypes.ModelType
	Model                string
	LoadBalancingConfigs []*providerentities.ModelLoadBalancingConfiguration
}

func NewLBModelManager(
	tenant_id string,
	provider string,
	model_type modelruntimeenumtypes.ModelType,
	model string,
	load_balancing_configs []*providerentities.ModelLoadBalancingConfiguration,
	managed_credentials map[string]any,
) *LBModelManager {
	/*
	   Load balancing model manager
	   :param tenant_id: tenant_id
	   :param provider: provider
	   :param model_type: model_type
	   :param model: model name
	   :param load_balancing_configs: all load balancing configurations
	   :param managed_credentials: credentials if load balancing configuration name is __inherit__
	*/
	lmmm := &LBModelManager{
		TenantID:             tenant_id,
		Provider:             provider,
		ModelType:            model_type,
		Model:                model,
		LoadBalancingConfigs: load_balancing_configs,
	}

	for idx, load_balancing_config := range lmmm.LoadBalancingConfigs { // Iterate over a shallow copy of the list
		if load_balancing_config.Name == "__inherit__" {
			if len(managed_credentials) == 0 {
				// remove __inherit__ if managed credentials is not provided
				lmmm.LoadBalancingConfigs = slices.DeleteFunc(lmmm.LoadBalancingConfigs, func(val *providerentities.ModelLoadBalancingConfiguration) bool {
					return val == load_balancing_config
				})
			} else {
				load_balancing_config.Credentials = managed_credentials
				lmmm.LoadBalancingConfigs[idx] = load_balancing_config
			}
		}
	}
	return lmmm
}

func (mgr *LBModelManager) fetch_next() *providerentities.ModelLoadBalancingConfiguration {
	/*
	   Get next model load balancing config
	   Strategy: Round Robin
	   :return:
	*/
	cache_key := fmt.Sprintf("model_lb_index:%s:%s:%s:%s", mgr.TenantID, mgr.Provider, mgr.ModelType, mgr.Model)

	cooldown_load_balancing_configs := []*providerentities.ModelLoadBalancingConfiguration{}
	max_index := int64(len(mgr.LoadBalancingConfigs))

	for {
		current_index := cache.Instance().IncrKey(cache_key)
		if current_index >= 10000000 {
			current_index = 1
			cache.Instance().Set(cache_key, current_index)
		}
		cache.Instance().ExpireKey(cache_key, 3600)
		if current_index > max_index {
			current_index = current_index % max_index
		}
		real_index := current_index - 1
		if real_index > max_index {
			real_index = 0
		}
		config := mgr.LoadBalancingConfigs[real_index]

		if mgr.in_cooldown(config) {
			cooldown_load_balancing_configs = append(cooldown_load_balancing_configs, config)
			if len(cooldown_load_balancing_configs) >= len(mgr.LoadBalancingConfigs) {
				// all configs are in cooldown
				return nil
			}
			continue
		}
		mlog.Debugf("Model LB\nid: %s\nname:%s\ntenant_id: %s\nprovider: %s\nmodel_type: %s\nmodel: %s", config.ID, config.Name, mgr.TenantID, mgr.Provider, mgr.ModelType, mgr.Model)

		return config
	}
	// return nil
}
func (mgr *LBModelManager) cooldown(config *providerentities.ModelLoadBalancingConfiguration, expire int /* = 60*/) {
	/*
	   Cooldown model load balancing config
	   :param config: model load balancing config
	   :param expire: cooldown time
	   :return:
	*/
	cooldown_cache_key := fmt.Sprintf("model_lb_index:cooldown:%s:%s:%s:%s:%s", mgr.TenantID, mgr.Provider, mgr.ModelType, mgr.Model, config.ID)
	cache.Instance().SetExKey(cooldown_cache_key, "true", time.Duration(expire)*time.Second)
}
func (mgr *LBModelManager) in_cooldown(config *providerentities.ModelLoadBalancingConfiguration) bool {
	/*
	   Check if model load balancing config is in cooldown
	   :param config: model load balancing config
	   :return:
	*/
	cooldown_cache_key := fmt.Sprintf("model_lb_index:cooldown:%s:%s:%s:%s:%s", mgr.TenantID, mgr.Provider, mgr.ModelType, mgr.Model, config.ID)

	return cache.Instance().ExistsKey(cooldown_cache_key)
}
func (mgr *LBModelManager) get_config_in_cooldown_and_ttl(
	tenant_id string, provider string, model_type modelruntimeenumtypes.ModelType, model string, config_id string,
) (bool, int) {
	/*
	   Get model load balancing config is in cooldown and ttl
	   :param tenant_id: workspace id
	   :param provider: provider name
	   :param model_type: model type
	   :param model: model name
	   :param config_id: model load balancing config id
	   :return:
	*/
	cooldown_cache_key := fmt.Sprintf("model_lb_index:cooldown:%s:%s:%s:%s:%s", tenant_id, provider, model_type, model, config_id)

	ttl := cache.Instance().TTL(cooldown_cache_key)
	if ttl == -2 {
		return false, 0
	}
	return true, ttl
}
