package dataset

import (
	"maps"
	"slices"

	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	"mlib.com/gofy/server/utils/mapstruct"
)

type DatasetConfigManager struct {
}

func (mgr *DatasetConfigManager) Convert(config map[string]any) *appconfigentities.DatasetEntity {
	dataset_ids := []string{}
	dataset_configs := map[string]any{}
	if _, ok := config["dataset_configs"]; ok {
		if _, ok := config["dataset_configs"].(map[string]any); ok {
			dataset_configs = config["dataset_configs"].(map[string]any)
		}
	}
	if _, ok := dataset_configs["datasets"]; ok {
		datasets := mapstruct.Get(dataset_configs, "datasets", map[string]any{"strategy": "router", "datasets": []any{}})

		datasets_datasets_list := mapstruct.Get(datasets, "datasets", []map[string]any{})
		for _, dataset := range datasets_datasets_list {
			keys := slices.Sorted(maps.Keys(dataset))
			if len(keys) == 0 {
				continue
			}
			if _, ok := dataset["dataset"]; !ok {
				continue
			}
			new_dataset := mapstruct.Get(dataset, "dataset", map[string]any{})
			if _, ok := new_dataset["enabled"]; !ok || new_dataset["enabled"] == nil {
				continue
			}
			if _, ok := new_dataset["enabled"].(bool); !ok || !new_dataset["enabled"].(bool) {
				continue
			}

			dataset_id := mapstruct.Get(new_dataset, "id", "")
			if dataset_id != "" {
				dataset_ids = append(dataset_ids, dataset_id)
			}
		}
	}
	agent_dict := mapstruct.Get(config, "agent_mode", map[string]any{})
	enabled := mapstruct.Get(agent_dict, "enabled", false)
	if len(agent_dict) > 0 && enabled {

		for _, tool := range mapstruct.Get(agent_dict, "tools", []map[string]any{}) {
			keys := slices.Sorted(maps.Keys(tool))
			if len(tool) == 1 {
				// old standard
				key := keys[0]

				if key != "dataset" {
					continue
				}
				tool_item := map[string]any{}
				if _, ok := tool[key].(map[string]any); ok {
					tool_item = tool[key].(map[string]any)
				} else {
					continue
				}

				if !mapstruct.Get(tool_item, "enabled", false) {
					continue
				}
				dataset_id := mapstruct.Get(tool_item, "id", "")
				dataset_ids = append(dataset_ids, dataset_id)
			}
		}
	}
	if len(dataset_ids) == 0 {
		return nil
	}
	// dataset configs
	dataset_configs = mapstruct.Get(config, "dataset_configs", map[string]any{"retrieval_model": "multiple"})
	if len(dataset_configs) == 0 {
		return nil
	}
	query_variable := mapstruct.Get(config, "dataset_query_variable", "")
	retrieval_model := mapstruct.Get(dataset_configs, "retrieval_model", "")
	if retrieval_model == "single" {
		entity := &appconfigentities.DatasetEntity{
			DatasetIDs: dataset_ids,
			RetrieveConfig: appconfigentities.DatasetRetrieveConfigEntity{
				QueryVariable:         query_variable,
				RetrieveStrategy:      appconfigenumtypes.RetrieveStrategy("single"),
				MetadataFilteringMode: appconfigenumtypes.MetadataFilteringModeType(mapstruct.Get(dataset_configs, "metadata_filtering_mode", "disabled")),
			},
		}
		if _, ok := dataset_configs["metadata_model_config"]; ok {
			if _, ok := dataset_configs["metadata_model_config"].(map[string]any); ok {
				entity.RetrieveConfig.MetadataModelConfig = appconfigentities.NewModelConfig(dataset_configs["metadata_model_config"].(map[string]any))
			}
		}
		if _, ok := dataset_configs["metadata_filtering_conditions"]; ok {
			if _, ok := dataset_configs["metadata_filtering_conditions"].(map[string]any); ok {
				entity.RetrieveConfig.MetadataFilteringConditions = appconfigentities.NewMetadataFilteringCondition(dataset_configs["metadata_filtering_conditions"].(map[string]any))
			}
		}
		return entity
	} else {
		entity := &appconfigentities.DatasetEntity{
			DatasetIDs: dataset_ids,
			RetrieveConfig: appconfigentities.DatasetRetrieveConfigEntity{
				QueryVariable:         query_variable,
				RetrieveStrategy:      appconfigenumtypes.RetrieveStrategy(retrieval_model),
				MetadataFilteringMode: appconfigenumtypes.MetadataFilteringModeType(mapstruct.Get(dataset_configs, "metadata_filtering_mode", "disabled")),
				TopK:                  mapstruct.Get(dataset_configs, "top_k", 4),
				ScoreThreshold:        mapstruct.Get(dataset_configs, "score_threshold", 0.0),
				RerankingModel:        mapstruct.Get(dataset_configs, "reranking_model", map[string]any{}),
				Weights:               mapstruct.Get(dataset_configs, "weights", map[string]any{}),
				RerankingEnabled:      mapstruct.Get(dataset_configs, "reranking_enabled", true),
				RerankMode:            mapstruct.Get(dataset_configs, "reranking_mode", "reranking_model"),
			},
		}

		if _, ok := dataset_configs["metadata_model_config"]; ok {
			if _, ok := dataset_configs["metadata_model_config"].(map[string]any); ok {
				entity.RetrieveConfig.MetadataModelConfig = appconfigentities.NewModelConfig(dataset_configs["metadata_model_config"].(map[string]any))
			}
		}
		if _, ok := dataset_configs["metadata_filtering_conditions"]; ok {
			if _, ok := dataset_configs["metadata_filtering_conditions"].(map[string]any); ok {
				entity.RetrieveConfig.MetadataFilteringConditions = appconfigentities.NewMetadataFilteringCondition(dataset_configs["metadata_filtering_conditions"].(map[string]any))
			}
		}
		return entity
	}
}
