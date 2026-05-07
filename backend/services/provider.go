package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	coreentities "github.com/odysseythink/gofy/backend/entities/core"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type ProviderService struct {
}

func (s *ProviderService) getAllProviders(tenant_id string) map[string][]*models.Provider {
	/*
	   Get all provider records of the workspace.

	   :param tenant_id: workspace id
	   :return:
	*/
	providers := []*models.Provider{}
	err := dbengine.Instance().DB.Model(&models.Provider{}).Where("tenant_id = ? and is_valid = ?", tenant_id, true).Preload("Tenant").Find(providers).Error
	if err != nil {
		mlog.Warningf("get Provider by(tenant_id = %s) failed:%v", tenant_id, err)
		return nil
	}
	// providers = db.session.query(Provider).filter(Provider.tenant_id == tenant_id, Provider.is_valid == True).all()

	var provider_name_to_provider_records_dict map[string][]*models.Provider
	for _, provider := range providers {
		if provider_name_to_provider_records_dict == nil {
			provider_name_to_provider_records_dict = make(map[string][]*models.Provider)
		}
		if _, ok := provider_name_to_provider_records_dict[provider.ProviderName]; !ok {
			provider_name_to_provider_records_dict[provider.ProviderName] = make([]*models.Provider, 0)
		}

		provider_name_to_provider_records_dict[provider.ProviderName] = append(provider_name_to_provider_records_dict[provider.ProviderName], provider)
	}
	return provider_name_to_provider_records_dict
}
func (s *ProviderService) GetConfigurations(tenantId string) *coreentities.ProviderConfigurations {
	/*
	   Get model provider configurations.

	   Construct ProviderConfiguration objects for each provider
	*/

	// // Get all provider records of the workspace
	// providerNameToProviderRecordsDict := s.getAllProviders(tenantId)

	// // Initialize trial provider records if not exist
	// providerNameToProviderRecordsDict = s.initTrialProviderRecords(tenantId, providerNameToProviderRecordsDict)

	// // Get all provider model records of the workspace
	// providerNameToProviderModelRecordsDict := s.getAllProviderModels(tenantId)

	// // Get all provider entities
	// providerEntities := modelProviderFactory.GetProviders()

	// // Get All preferred provider types of the workspace
	// providerNameToPreferredModelProviderRecordsDict := s.getAllPreferredModelProviders(tenantId)

	// // Get All provider model settings
	// providerNameToProviderModelSettingsDict := s.getAllProviderModelSettings(tenantId)

	// // Get All load balancing configs
	// providerNameToProviderLoadBalancingModelConfigsDict := s.getAllProviderLoadBalancingConfigs(tenantId)

	// providerConfigurations := &entities.ProviderConfigurations{
	// 	TenantID:                  tenantId,
	// 	ProviderConfigurationsMap: make(map[string]entities.ProviderConfigurations),
	// }

	// // Construct ProviderConfiguration objects for each provider
	// for _, providerEntity := range providerEntities {
	// 	// handle include, exclude
	// 	if isFiltered(
	// 		cast.ToStringSet(gofyConfig.POSITION_PROVIDER_INCLUDES_SET),
	// 		cast.ToStringSet(gofyConfig.POSITION_PROVIDER_EXCLUDES_SET),
	// 		providerEntity,
	// 		func(x any) string {
	// 			if p, ok := x.(*ProviderEntity); ok {
	// 				return p.provider
	// 			}
	// 			return ""
	// 		},
	// 	) {
	// 		continue
	// 	}

	// 	providerName := providerEntity.provider
	// 	providerRecords := providerNameToProviderRecordsDict[providerName]
	// 	providerModelRecords := providerNameToProviderModelRecordsDict[providerName]

	// 	// Convert to custom configuration
	// 	customConfiguration := pm.toCustomConfiguration(tenantId, providerEntity, providerRecords, providerModelRecords)

	// 	// Convert to system configuration
	// 	systemConfiguration := pm.toSystemConfiguration(tenantId, providerEntity, providerRecords)

	// 	// Get preferred provider type
	// 	preferredProviderTypeRecord := providerNameToPreferredModelProviderRecordsDict[providerName]

	// 	var preferredProviderType ProviderType
	// 	if preferredProviderTypeRecord != nil {
	// 		preferredProviderType = ProviderType.valueOf(preferredProviderTypeRecord.preferred_provider_type)
	// 	} else if customConfiguration.provider != nil || len(customConfiguration.models) > 0 {
	// 		preferredProviderType = ProviderTypeCUSTOM
	// 	} else if systemConfiguration.enabled {
	// 		preferredProviderType = ProviderTypeSYSTEM
	// 	} else {
	// 		preferredProviderType = ProviderTypeCUSTOM
	// 	}

	// 	usingProviderType := preferredProviderType
	// 	hasValidQuota := false
	// 	for _, quotaConf := range systemConfiguration.quotaConfigurations {
	// 		if quotaConf.isValid {
	// 			hasValidQuota = true
	// 			break
	// 		}
	// 	}

	// 	if preferredProviderType == ProviderTypeSYSTEM {
	// 		if !systemConfiguration.enabled || !hasValidQuota {
	// 			usingProviderType = ProviderTypeCUSTOM
	// 		}
	// 	} else {
	// 		if customConfiguration.provider == nil && len(customConfiguration.models) == 0 {
	// 			if systemConfiguration.enabled && hasValidQuota {
	// 				usingProviderType = ProviderTypeSYSTEM
	// 			}
	// 		}
	// 	}

	// 	// Get provider load balancing configs
	// 	providerModelSettings := providerNameToProviderModelSettingsDict[providerName]

	// 	// Get provider load balancing configs
	// 	providerLoadBalancingConfigs := providerNameToProviderLoadBalancingModelConfigsDict[providerName]

	// 	// Convert to model settings
	// 	modelSettings := pm.toModelSettings(
	// 		providerEntity,
	// 		providerModelSettings,
	// 		providerLoadBalancingConfigs,
	// 	)

	// 	providerConfiguration := ProviderConfiguration{
	// 		tenantId:              tenantId,
	// 		provider:              providerEntity,
	// 		preferredProviderType: preferredProviderType,
	// 		usingProviderType:     usingProviderType,
	// 		systemConfiguration:   systemConfiguration,
	// 		customConfiguration:   customConfiguration,
	// 		modelSettings:         modelSettings,
	// 	}

	// 	providerConfigurations.providerConfigurationsMap[providerName] = providerConfiguration
	// }

	// // Return the encapsulated object
	// return providerConfigurations
	return nil
}
