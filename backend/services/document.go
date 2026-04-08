package services

import (
	"fmt"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type DocumentService struct{}

func (s *DocumentService) GetDocument(datasetID string, documentID string) *models.Document {
	var doc models.Document
	query := dbengine.Instance().DB.Where("dataset_id = ?", datasetID)
	if documentID != "" {
		query = query.Where("id = ?", documentID)
	}
	if err := query.First(&doc).Error; err != nil {
		return nil
	}
	return &doc
}

func (s *DocumentService) GetDocumentsByIDs(datasetID string, documentIDs []string) []*models.Document {
	var docs []*models.Document
	dbengine.Instance().DB.Where("dataset_id = ? AND id IN ?", datasetID, documentIDs).Find(&docs)
	return docs
}

func (s *DocumentService) GetDocumentByID(documentID string) *models.Document {
	var doc models.Document
	if err := dbengine.Instance().DB.Where("id = ?", documentID).First(&doc).Error; err != nil {
		return nil
	}
	return &doc
}

func (s *DocumentService) GetDocumentsByDatasetID(datasetID string) []*models.Document {
	var docs []*models.Document
	dbengine.Instance().DB.Where("dataset_id = ?", datasetID).Find(&docs)
	return docs
}

func (s *DocumentService) GetBatchDocuments(datasetID string, batch string) []*models.Document {
	var docs []*models.Document
	dbengine.Instance().DB.Where("dataset_id = ? AND batch = ?", datasetID, batch).Find(&docs)
	return docs
}

func (s *DocumentService) DeleteDocument(document *models.Document) error {
	return dbengine.Instance().DB.Delete(document).Error
}

func (s *DocumentService) DeleteDocuments(datasetID string, documentIDs []string) error {
	return dbengine.Instance().DB.Where("dataset_id = ? AND id IN ?", datasetID, documentIDs).Delete(&models.Document{}).Error
}

func (s *DocumentService) RenameDocument(datasetID string, documentID string, name string) (*models.Document, error) {
	var doc models.Document
	if err := dbengine.Instance().DB.Where("dataset_id = ? AND id = ?", datasetID, documentID).First(&doc).Error; err != nil {
		return nil, fmt.Errorf("document not found")
	}
	doc.Name = name
	if err := dbengine.Instance().DB.Save(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *DocumentService) PauseDocument(document *models.Document) error {
	if document.IndexingStatus != "indexing" {
		return fmt.Errorf("document is not in indexing status")
	}
	document.IsPaused = true
	return dbengine.Instance().DB.Save(document).Error
}

func (s *DocumentService) RecoverDocument(document *models.Document) error {
	if !document.IsPaused {
		return fmt.Errorf("document is not paused")
	}
	document.IsPaused = false
	return dbengine.Instance().DB.Save(document).Error
}

func (s *DocumentService) BatchUpdateDocumentStatus(datasetID string, documentIDs []string, action string) error {
	now := time.Now()
	var updates map[string]any

	switch action {
	case "enable":
		updates = map[string]any{"enabled": true, "disabled_at": nil, "disabled_by": nil, "updated_at": now}
	case "disable":
		updates = map[string]any{"enabled": false, "disabled_at": now, "updated_at": now}
	case "archive":
		updates = map[string]any{"archived": true, "archived_at": now, "updated_at": now}
	case "un_archive":
		updates = map[string]any{"archived": false, "archived_at": nil, "updated_at": now}
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id = ? AND id IN ?", datasetID, documentIDs).Updates(updates).Error
}

func (s *DocumentService) GetDocumentsPosition(datasetID string) int {
	var maxPos *int
	dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id = ?", datasetID).Select("MAX(position)").Scan(&maxPos)
	if maxPos == nil {
		return 0
	}
	return *maxPos
}

func (s *DocumentService) GetDocumentFileDetail(fileID string) *models.UploadFile {
	var file models.UploadFile
	if err := dbengine.Instance().DB.Where("id = ?", fileID).First(&file).Error; err != nil {
		return nil
	}
	return &file
}

func (s *DocumentService) CheckArchived(document *models.Document) error {
	if document.Archived {
		return fmt.Errorf("document is archived")
	}
	return nil
}

func (s *DocumentService) DocumentCount(datasetID string) int64 {
	var count int64
	dbengine.Instance().DB.Model(&models.Document{}).Where("dataset_id = ?", datasetID).Count(&count)
	return count
}

func (s *DocumentService) GetPaginatedDocuments(datasetID string, page, perPage int, search string, status string) ([]*models.Document, int64) {
	var docs []*models.Document
	var total int64

	query := dbengine.Instance().DB.Where("dataset_id = ?", datasetID)
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("indexing_status = ?", status)
	}

	query.Model(&models.Document{}).Count(&total)
	query.Order("position ASC, created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&docs)
	return docs, total
}
