package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/tasks"
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

// SaveDocumentWithDatasetID creates documents for a dataset and triggers indexing.
func (s *DocumentService) SaveDocumentWithDatasetID(
	dataset *models.Dataset,
	dataSourceType string,
	dataSourceInfoList []map[string]any,
	processRuleID string,
	documentForm string,
	documentLanguage string,
	createdFrom string,
	account *models.Account,
) ([]*models.Document, string, error) {
	// Generate batch ID
	batch := fmt.Sprintf("%s%06d", time.Now().Format("20060102150405"), rand.Intn(1000000))

	// Get next position
	position := s.GetDocumentsPosition(dataset.ID) + 1

	var documents []*models.Document

	for i, sourceInfo := range dataSourceInfoList {
		name := ""
		switch dataSourceType {
		case "upload_file":
			fileID, _ := sourceInfo["upload_file_id"].(string)
			var uploadFile models.UploadFile
			if err := dbengine.Instance().DB.Where("id = ?", fileID).First(&uploadFile).Error; err == nil {
				name = uploadFile.Name
			}
		case "notion_import":
			name, _ = sourceInfo["page_name"].(string)
			if name == "" {
				name, _ = sourceInfo["title"].(string)
			}
		case "website_crawl":
			name, _ = sourceInfo["url"].(string)
		}
		if name == "" {
			name = fmt.Sprintf("document_%d", i+1)
		}

		dataSourceInfoJSON, _ := json.Marshal(sourceInfo)

		doc := s.BuildDocument(
			dataset, processRuleID, dataSourceType,
			documentForm, documentLanguage,
			string(dataSourceInfoJSON), createdFrom,
			position+i, account, name, batch,
		)
		documents = append(documents, doc)
	}

	// Batch create
	if len(documents) > 0 {
		if err := dbengine.Instance().DB.CreateInBatches(documents, 50).Error; err != nil {
			return nil, "", fmt.Errorf("failed to create documents: %w", err)
		}
	}

	// Trigger indexing tasks
	for _, doc := range documents {
		tasks.EnqueueTask(tasks.TaskDocumentIndexing, tasks.DocumentIndexingPayload{
			DatasetID:  dataset.ID,
			DocumentID: doc.ID,
		})
	}

	return documents, batch, nil
}

// BuildDocument creates a Document model instance (not persisted).
func (s *DocumentService) BuildDocument(
	dataset *models.Dataset,
	processRuleID string,
	dataSourceType string,
	documentForm string,
	documentLanguage string,
	dataSourceInfo string,
	createdFrom string,
	position int,
	account *models.Account,
	name string,
	batch string,
) *models.Document {
	doc := &models.Document{
		ID:                   uuid.NewV4().String(),
		TenantID:             dataset.TenantID,
		DatasetID:            dataset.ID,
		Batch:                batch,
		Name:                 name,
		DataSourceType:       dataSourceType,
		DataSourceInfo:       dataSourceInfo,
		DatasetProcessRuleID: processRuleID,
		DocForm:              documentForm,
		DocLanguage:          documentLanguage,
		IndexingStatus:       "waiting",
		Position:             position,
		Enabled:              true,
		CreatedFrom:          createdFrom,
		CreatedBy:            account.ID,
		CreatedAPIRequestID:  "",
	}
	return doc
}

// UpdateDocumentWithDatasetID updates an existing document and triggers re-indexing.
func (s *DocumentService) UpdateDocumentWithDatasetID(
	dataset *models.Dataset,
	documentID string,
	dataSourceType string,
	dataSourceInfo string,
	processRuleID string,
	documentForm string,
	name string,
	account *models.Account,
) (*models.Document, error) {
	var doc models.Document
	if err := dbengine.Instance().DB.Where("id = ? AND dataset_id = ?", documentID, dataset.ID).First(&doc).Error; err != nil {
		return nil, fmt.Errorf("document not found")
	}

	updates := map[string]any{
		"indexing_status":        "waiting",
		"processing_started_at": nil,
		"completed_at":          nil,
		"error":                 nil,
		"stopped_at":            nil,
	}

	if dataSourceType != "" {
		updates["data_source_type"] = dataSourceType
	}
	if dataSourceInfo != "" {
		updates["data_source_info"] = dataSourceInfo
	}
	if processRuleID != "" {
		updates["dataset_process_rule_id"] = processRuleID
	}
	if documentForm != "" {
		updates["doc_form"] = documentForm
	}
	if name != "" {
		updates["name"] = name
	}

	if err := dbengine.Instance().DB.Model(&doc).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reset segments status
	dbengine.Instance().DB.Model(&models.DocumentSegment{}).
		Where("document_id = ?", documentID).
		Update("status", "re_segment")

	// Trigger re-indexing
	tasks.EnqueueTask(tasks.TaskDocumentIndexingUpdate, tasks.DocumentIndexingPayload{
		DatasetID:  dataset.ID,
		DocumentID: documentID,
	})

	dbengine.Instance().DB.Where("id = ?", documentID).First(&doc)
	return &doc, nil
}

// SaveDocumentWithoutDatasetID creates a new dataset and adds documents to it.
func (s *DocumentService) SaveDocumentWithoutDatasetID(
	tenantID string,
	dataSourceType string,
	dataSourceInfoList []map[string]any,
	indexingTechnique string,
	processRuleMode string,
	processRules string,
	documentForm string,
	documentLanguage string,
	account *models.Account,
) (*models.Dataset, []*models.Document, string, error) {
	// Create dataset
	dsService := &DatasetService{}
	dataset, err := dsService.CreateEmptyDataset(
		tenantID, "", "", indexingTechnique, account, "only_me", "vendor", "", "",
	)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create dataset: %w", err)
	}

	// Create process rule
	ruleID := ""
	if processRuleMode != "" {
		rule := &models.DatasetProcessRule{
			ID:        uuid.NewV4().String(),
			DatasetID: dataset.ID,
			Mode:      processRuleMode,
			Rules:     processRules,
		}
		dbengine.Instance().DB.Create(rule)
		ruleID = rule.ID
	}

	// Create documents
	docs, batch, err := s.SaveDocumentWithDatasetID(
		dataset, dataSourceType, dataSourceInfoList,
		ruleID, documentForm, documentLanguage, "web", account,
	)
	if err != nil {
		return nil, nil, "", err
	}

	// Update dataset name from first document
	if len(docs) > 0 {
		dsName := docs[0].Name
		if len([]rune(dsName)) > 18 {
			dsName = string([]rune(dsName)[:18]) + "..."
		}
		dbengine.Instance().DB.Model(dataset).Updates(map[string]any{
			"name":        dsName,
			"description": fmt.Sprintf("useful for when you want to answer queries about the %s", docs[0].Name),
		})
	}

	return dataset, docs, batch, nil
}

// RetryDocument retries indexing for failed documents.
func (s *DocumentService) RetryDocument(datasetID string, documents []*models.Document) error {
	for _, doc := range documents {
		dbengine.Instance().DB.Model(doc).Updates(map[string]any{
			"indexing_status":        "waiting",
			"processing_started_at": nil,
			"completed_at":          nil,
			"error":                 nil,
		})

		// Reset segments
		dbengine.Instance().DB.Model(&models.DocumentSegment{}).
			Where("document_id = ?", doc.ID).
			Update("status", "re_segment")

		tasks.EnqueueTask(tasks.TaskDocumentIndexing, tasks.DocumentIndexingPayload{
			DatasetID:  datasetID,
			DocumentID: doc.ID,
		})
	}
	return nil
}

// SyncWebsiteDocument re-crawls and re-indexes a website document.
func (s *DocumentService) SyncWebsiteDocument(datasetID string, document *models.Document) error {
	dbengine.Instance().DB.Model(document).Updates(map[string]any{
		"indexing_status":        "waiting",
		"processing_started_at": nil,
		"completed_at":          nil,
	})
	tasks.EnqueueTask(tasks.TaskDocumentIndexing, tasks.DocumentIndexingPayload{
		DatasetID:  datasetID,
		DocumentID: document.ID,
	})
	return nil
}

// GetDocumentDownloadURL generates a download URL for a document's source file.
func (s *DocumentService) GetDocumentDownloadURL(document *models.Document) (string, error) {
	if document.DataSourceType != "upload_file" {
		return "", fmt.Errorf("document is not an upload file")
	}
	var info map[string]any
	json.Unmarshal([]byte(document.DataSourceInfo), &info)
	fileID, _ := info["upload_file_id"].(string)
	if fileID == "" {
		return "", fmt.Errorf("upload file ID not found")
	}
	// TODO: Generate signed download URL
	return fmt.Sprintf("/files/%s/download", fileID), nil
}

// GetTenantDocumentsCount returns total document count across all datasets for a tenant.
func (s *DocumentService) GetTenantDocumentsCount(tenantID string) int64 {
	var count int64
	dbengine.Instance().DB.Model(&models.Document{}).
		Joins("JOIN datasets ON datasets.id = documents.dataset_id").
		Where("datasets.tenant_id = ?", tenantID).Count(&count)
	return count
}

// GetWorkingDocumentsByDatasetID returns documents being actively indexed.
func (s *DocumentService) GetWorkingDocumentsByDatasetID(datasetID string) []*models.Document {
	var docs []*models.Document
	dbengine.Instance().DB.Where("dataset_id = ? AND indexing_status IN ?", datasetID, []string{"waiting", "parsing", "cleaning", "splitting", "indexing"}).Find(&docs)
	return docs
}

// GetErrorDocumentsByDatasetID returns documents that failed indexing.
func (s *DocumentService) GetErrorDocumentsByDatasetID(datasetID string) []*models.Document {
	var docs []*models.Document
	dbengine.Instance().DB.Where("dataset_id = ? AND indexing_status = ?", datasetID, "error").Find(&docs)
	return docs
}
