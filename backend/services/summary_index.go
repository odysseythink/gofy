package services

import (
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"

	uuid "github.com/satori/go.uuid"
)

type SummaryIndexService struct{}

func (s *SummaryIndexService) CreateSummaryRecord(datasetID, documentID, chunkID, summaryContent string, tokens int) (*models.DocumentSegmentSummary, error) {
	record := &models.DocumentSegmentSummary{
		ID:             uuid.NewV4().String(),
		DatasetID:      datasetID,
		DocumentID:     documentID,
		ChunkID:        chunkID,
		SummaryContent: &summaryContent,
		Tokens:         &tokens,
		Status:         "completed",
		Enabled:        true,
	}
	if err := dbengine.Instance().DB.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SummaryIndexService) GetSegmentSummary(segmentID, datasetID string) *models.DocumentSegmentSummary {
	var summary models.DocumentSegmentSummary
	if err := dbengine.Instance().DB.Where("chunk_id = ? AND dataset_id = ?", segmentID, datasetID).First(&summary).Error; err != nil {
		return nil
	}
	return &summary
}

func (s *SummaryIndexService) GetDocumentSummaries(documentID, datasetID string, page, limit int) ([]*models.DocumentSegmentSummary, int64) {
	var summaries []*models.DocumentSegmentSummary
	var total int64
	query := dbengine.Instance().DB.Where("document_id = ? AND dataset_id = ?", documentID, datasetID)
	query.Model(&models.DocumentSegmentSummary{}).Count(&total)
	query.Offset((page - 1) * limit).Limit(limit).Find(&summaries)
	return summaries, total
}

func (s *SummaryIndexService) DisableSummaries(segmentIDs []string, datasetID string) error {
	now := time.Now()
	return dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).
		Where("chunk_id IN ? AND dataset_id = ?", segmentIDs, datasetID).
		Updates(map[string]any{"enabled": false, "disabled_at": &now}).Error
}

func (s *SummaryIndexService) EnableSummaries(segmentIDs []string, datasetID string) error {
	return dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).
		Where("chunk_id IN ? AND dataset_id = ?", segmentIDs, datasetID).
		Updates(map[string]any{"enabled": true, "disabled_at": nil}).Error
}

func (s *SummaryIndexService) DeleteSummaries(segmentIDs []string, datasetID string) error {
	return dbengine.Instance().DB.Where("chunk_id IN ? AND dataset_id = ?", segmentIDs, datasetID).Delete(&models.DocumentSegmentSummary{}).Error
}

func (s *SummaryIndexService) UpdateSummaryError(summaryID string, errMsg string) error {
	return dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).Where("id = ?", summaryID).Updates(map[string]any{
		"status": "error",
		"error":  errMsg,
	}).Error
}

func (s *SummaryIndexService) GetDocumentSummaryStatus(documentID, datasetID string) map[string]int64 {
	result := map[string]int64{"total": 0, "completed": 0, "generating": 0, "error": 0}
	var total, completed, generating, errCount int64
	baseQuery := dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).Where("document_id = ? AND dataset_id = ?", documentID, datasetID)
	baseQuery.Count(&total)
	dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).Where("document_id = ? AND dataset_id = ? AND status = ?", documentID, datasetID, "completed").Count(&completed)
	dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).Where("document_id = ? AND dataset_id = ? AND status = ?", documentID, datasetID, "generating").Count(&generating)
	dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).Where("document_id = ? AND dataset_id = ? AND status = ?", documentID, datasetID, "error").Count(&errCount)
	result["total"] = total
	result["completed"] = completed
	result["generating"] = generating
	result["error"] = errCount
	return result
}

// BatchCreateSummaryRecords creates summary records for multiple segments.
func (s *SummaryIndexService) BatchCreateSummaryRecords(datasetID, documentID string, segmentSummaries []map[string]any) ([]*models.DocumentSegmentSummary, error) {
	var records []*models.DocumentSegmentSummary
	for _, ss := range segmentSummaries {
		chunkID, _ := ss["chunk_id"].(string)
		content, _ := ss["content"].(string)
		tokens := 0
		if t, ok := ss["tokens"].(int); ok {
			tokens = t
		} else {
			tokens = len([]rune(content)) / 2
		}
		record := &models.DocumentSegmentSummary{
			ID:             uuid.NewV4().String(),
			DatasetID:      datasetID,
			DocumentID:     documentID,
			ChunkID:        chunkID,
			SummaryContent: &content,
			Tokens:         &tokens,
			Status:         "completed",
			Enabled:        true,
		}
		records = append(records, record)
	}
	if len(records) > 0 {
		if err := dbengine.Instance().DB.CreateInBatches(records, 100).Error; err != nil {
			return nil, err
		}
	}
	return records, nil
}

// GetSegmentsSummaries returns summaries for multiple segments.
func (s *SummaryIndexService) GetSegmentsSummaries(segmentIDs []string, datasetID string) []*models.DocumentSegmentSummary {
	var summaries []*models.DocumentSegmentSummary
	dbengine.Instance().DB.Where("chunk_id IN ? AND dataset_id = ?", segmentIDs, datasetID).Find(&summaries)
	return summaries
}

// UpdateSummary updates the content of an existing summary.
func (s *SummaryIndexService) UpdateSummary(summaryID string, content string, tokens int) error {
	return dbengine.Instance().DB.Model(&models.DocumentSegmentSummary{}).Where("id = ?", summaryID).Updates(map[string]any{
		"summary_content": content,
		"tokens":          tokens,
		"status":          "completed",
		"error":           nil,
	}).Error
}

// GetDocumentsSummaryStatus returns summary status for multiple documents.
func (s *SummaryIndexService) GetDocumentsSummaryStatus(documentIDs []string, datasetID string) map[string]map[string]int64 {
	result := make(map[string]map[string]int64)
	for _, docID := range documentIDs {
		result[docID] = s.GetDocumentSummaryStatus(docID, datasetID)
	}
	return result
}

// GenerateSummaryForSegment generates a summary using LLM (placeholder).
func (s *SummaryIndexService) GenerateSummaryForSegment(datasetID, documentID, segmentID, content string) (*models.DocumentSegmentSummary, error) {
	// TODO: Call LLM to generate summary
	// For now, create a placeholder summary using first 200 chars
	summary := content
	if len([]rune(summary)) > 200 {
		summary = string([]rune(summary)[:200]) + "..."
	}
	tokens := len([]rune(summary)) / 2
	return s.CreateSummaryRecord(datasetID, documentID, segmentID, summary, tokens)
}
