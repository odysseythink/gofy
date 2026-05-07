package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"

	uuid "github.com/satori/go.uuid"
)

type MetadataService struct{}

func (s *MetadataService) CreateMetadata(datasetID, name, metadataType string) (*models.DatasetMetadata, error) {
	metadata := &models.DatasetMetadata{
		ID:        uuid.NewV4().String(),
		DatasetID: datasetID,
		Name:      name,
		Type:      metadataType,
	}
	if err := dbengine.Instance().DB.Create(metadata).Error; err != nil {
		return nil, err
	}
	return metadata, nil
}

func (s *MetadataService) UpdateMetadataName(datasetID, metadataID, name string) error {
	return dbengine.Instance().DB.Model(&models.DatasetMetadata{}).Where("id = ? AND dataset_id = ?", metadataID, datasetID).Update("name", name).Error
}

func (s *MetadataService) DeleteMetadata(datasetID, metadataID string) error {
	dbengine.Instance().DB.Where("metadata_id = ?", metadataID).Delete(&models.DatasetMetadataBinding{})
	return dbengine.Instance().DB.Where("id = ? AND dataset_id = ?", metadataID, datasetID).Delete(&models.DatasetMetadata{}).Error
}

func (s *MetadataService) GetDatasetMetadatas(datasetID string) []*models.DatasetMetadata {
	var metadatas []*models.DatasetMetadata
	dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Find(&metadatas)
	return metadatas
}

func (s *MetadataService) GetBuiltInFields() []map[string]any {
	return []map[string]any{
		{"name": "document_title", "type": "string"},
		{"name": "document_url", "type": "string"},
		{"name": "created_time", "type": "datetime"},
		{"name": "modified_time", "type": "datetime"},
		{"name": "author", "type": "string"},
		{"name": "language", "type": "string"},
	}
}

func (s *MetadataService) UpdateDocumentsMetadata(datasetID string, bindings []map[string]any) error {
	for _, b := range bindings {
		documentID, _ := b["document_id"].(string)
		metadataID, _ := b["metadata_id"].(string)
		createdBy, _ := b["created_by"].(string)
		if documentID == "" || metadataID == "" {
			continue
		}
		var existing models.DatasetMetadataBinding
		err := dbengine.Instance().DB.Where("document_id = ? AND metadata_id = ?", documentID, metadataID).First(&existing).Error
		if err != nil {
			dbengine.Instance().DB.Create(&models.DatasetMetadataBinding{
				ID:         uuid.NewV4().String(),
				DatasetID:  datasetID,
				DocumentID: documentID,
				MetadataID: metadataID,
				CreatedBy:  createdBy,
			})
		}
		// DatasetMetadataBinding has no value field; bindings are presence-based
	}
	return nil
}

func (s *MetadataService) DeleteDocumentMetadataBindings(documentID string) error {
	return dbengine.Instance().DB.Where("document_id = ?", documentID).Delete(&models.DatasetMetadataBinding{}).Error
}

func (s *MetadataService) GetDocumentMetadataBindings(documentID string) []*models.DatasetMetadataBinding {
	var bindings []*models.DatasetMetadataBinding
	dbengine.Instance().DB.Where("document_id = ?", documentID).Find(&bindings)
	return bindings
}
