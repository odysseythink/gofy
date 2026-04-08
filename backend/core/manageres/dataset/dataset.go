package dataset

import (
	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type DatasetManager struct {
}

func (mgr *DatasetManager) GetDatasetCollectionBinding(
	provider_name string, model_name string, collection_type string,
) *models.DatasetCollectionBinding {
	if collection_type == "" {
		collection_type = "dataset"
	}

	dataset_collection_binding := new(models.DatasetCollectionBinding)
	err := dbengine.Instance().DB.Model(&models.DatasetCollectionBinding{}).Where("provider_name =? and model_name =? and type = ?", provider_name, model_name, collection_type).Order("created_at DESC").First(dataset_collection_binding).Error
	if err != nil {
		mlog.Errorf("get DatasetCollectionBinding failed:%v", err)
		dataset_collection_binding = nil
	}

	if dataset_collection_binding == nil {
		dataset_collection_binding = &models.DatasetCollectionBinding{
			ID:             uuid.NewV4().String(),
			ProviderName:   provider_name,
			ModelName:      model_name,
			CollectionName: (&models.Dataset{}).GenCollectionNameByID(uuid.NewV4().String()),
			Type:           collection_type,
		}
		dbengine.Instance().DB.Create(dataset_collection_binding)
	}
	return dataset_collection_binding
}
