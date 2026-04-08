package services

import (
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type SegmentService struct{}

func (s *SegmentService) GetSegments(documentID string, tenantID string, statusList []string, keyword string, page int, limit int) ([]*models.DocumentSegment, int64) {
	var segments []*models.DocumentSegment
	var total int64

	query := dbengine.Instance().DB.Where("document_id = ? AND tenant_id = ?", documentID, tenantID)
	if len(statusList) > 0 {
		query = query.Where("status IN ?", statusList)
	}
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	query.Model(&models.DocumentSegment{}).Count(&total)
	query.Order("position ASC").Offset((page - 1) * limit).Limit(limit).Find(&segments)
	return segments, total
}

func (s *SegmentService) GetSegmentByID(segmentID string, tenantID string) *models.DocumentSegment {
	var segment models.DocumentSegment
	if err := dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", segmentID, tenantID).First(&segment).Error; err != nil {
		return nil
	}
	return &segment
}

func (s *SegmentService) GetSegmentsByDocumentAndDataset(documentID string, datasetID string, status string, enabled *bool) []*models.DocumentSegment {
	var segments []*models.DocumentSegment
	query := dbengine.Instance().DB.Where("document_id = ? AND dataset_id = ?", documentID, datasetID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	query.Order("position ASC").Find(&segments)
	return segments
}

func (s *SegmentService) CreateSegment(args map[string]any, document *models.Document, dataset *models.Dataset) (*models.DocumentSegment, error) {
	content, ok := args["content"].(string)
	if !ok || content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// Get max position
	var maxPos *int
	dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("document_id = ?", document.ID).Select("MAX(position)").Scan(&maxPos)
	position := 0
	if maxPos != nil {
		position = *maxPos + 1
	}

	answer, _ := args["answer"].(string)
	keywords, _ := args["keywords"].([]string)
	keywordsJSON, _ := json.Marshal(keywords)

	wordCount := len([]rune(content)) + len([]rune(answer))

	segment := &models.DocumentSegment{
		ID:          uuid.NewV4().String(),
		TenantID:    dataset.TenantID,
		DatasetID:   dataset.ID,
		DocumentID:  document.ID,
		Position:    position,
		Content:     content,
		Answer:      answer,
		WordCount:   wordCount,
		Keywords:    datatypes.JSON(keywordsJSON),
		Status:      "completed",
		Enabled:     true,
		IndexNodeID: uuid.NewV4().String(),
	}

	if err := dbengine.Instance().DB.Create(segment).Error; err != nil {
		return nil, err
	}
	return segment, nil
}

func (s *SegmentService) UpdateSegment(segmentID string, content string, answer string, keywords []string, enabled *bool) (*models.DocumentSegment, error) {
	var segment models.DocumentSegment
	if err := dbengine.Instance().DB.Where("id = ?", segmentID).First(&segment).Error; err != nil {
		return nil, fmt.Errorf("segment not found")
	}

	updates := map[string]any{"updated_at": time.Now()}
	if content != "" {
		updates["content"] = content
		updates["word_count"] = len([]rune(content))
	}
	if answer != "" {
		updates["answer"] = answer
	}
	if keywords != nil {
		keywordsJSON, _ := json.Marshal(keywords)
		updates["keywords"] = datatypes.JSON(keywordsJSON)
	}
	if enabled != nil {
		updates["enabled"] = *enabled
		if !*enabled {
			now := time.Now()
			updates["disabled_at"] = &now
		} else {
			updates["disabled_at"] = nil
		}
	}

	if err := dbengine.Instance().DB.Model(&segment).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &segment, nil
}

func (s *SegmentService) DeleteSegment(segmentID string, documentID string, datasetID string) error {
	return dbengine.Instance().DB.Where("id = ? AND document_id = ? AND dataset_id = ?", segmentID, documentID, datasetID).Delete(&models.DocumentSegment{}).Error
}

func (s *SegmentService) DeleteSegments(segmentIDs []string, documentID string, datasetID string) error {
	return dbengine.Instance().DB.Where("id IN ? AND document_id = ? AND dataset_id = ?", segmentIDs, documentID, datasetID).Delete(&models.DocumentSegment{}).Error
}

func (s *SegmentService) UpdateSegmentsStatus(segmentIDs []string, action string, datasetID string, documentID string) error {
	now := time.Now()
	var updates map[string]any

	switch action {
	case "enable":
		updates = map[string]any{"enabled": true, "disabled_at": nil, "disabled_by": nil, "updated_at": now}
	case "disable":
		updates = map[string]any{"enabled": false, "disabled_at": now, "updated_at": now}
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("id IN ? AND dataset_id = ? AND document_id = ?", segmentIDs, datasetID, documentID).Updates(updates).Error
}

func (s *SegmentService) MultiCreateSegment(segmentArgs []map[string]any, document *models.Document, dataset *models.Dataset) ([]*models.DocumentSegment, error) {
	var maxPos *int
	dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("document_id = ?", document.ID).Select("MAX(position)").Scan(&maxPos)
	startPos := 0
	if maxPos != nil {
		startPos = *maxPos + 1
	}

	segments := make([]*models.DocumentSegment, 0, len(segmentArgs))
	for i, args := range segmentArgs {
		content, _ := args["content"].(string)
		if content == "" {
			continue
		}
		answer, _ := args["answer"].(string)
		keywords, _ := args["keywords"].([]string)
		keywordsJSON, _ := json.Marshal(keywords)

		segment := &models.DocumentSegment{
			ID:          uuid.NewV4().String(),
			TenantID:    dataset.TenantID,
			DatasetID:   dataset.ID,
			DocumentID:  document.ID,
			Position:    startPos + i,
			Content:     content,
			Answer:      answer,
			WordCount:   len([]rune(content)) + len([]rune(answer)),
			Keywords:    datatypes.JSON(keywordsJSON),
			Status:      "completed",
			Enabled:     true,
			IndexNodeID: uuid.NewV4().String(),
		}
		segments = append(segments, segment)
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("no valid segments to create")
	}

	if err := dbengine.Instance().DB.CreateInBatches(segments, 100).Error; err != nil {
		return nil, err
	}
	return segments, nil
}

func (s *SegmentService) GetChildChunks(segmentID string, documentID string, datasetID string, page int, limit int, keyword string) ([]*models.ChildChunk, int64) {
	var chunks []*models.ChildChunk
	var total int64
	query := dbengine.Instance().DB.Where("segment_id = ? AND document_id = ? AND dataset_id = ?", segmentID, documentID, datasetID)
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}
	query.Model(&models.ChildChunk{}).Count(&total)
	query.Order("position ASC").Offset((page - 1) * limit).Limit(limit).Find(&chunks)
	return chunks, total
}

func (s *SegmentService) GetChildChunkByID(childChunkID string, tenantID string) *models.ChildChunk {
	var chunk models.ChildChunk
	if err := dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", childChunkID, tenantID).First(&chunk).Error; err != nil {
		return nil
	}
	return &chunk
}

func (s *SegmentService) CreateChildChunk(content string, segment *models.DocumentSegment, document *models.Document, dataset *models.Dataset) (*models.ChildChunk, error) {
	var maxPos *int
	dbengine.Instance().DB.Model(&models.ChildChunk{}).Where("segment_id = ?", segment.ID).Select("MAX(position)").Scan(&maxPos)
	position := 0
	if maxPos != nil {
		position = *maxPos + 1
	}

	chunk := &models.ChildChunk{
		ID:          uuid.NewV4().String(),
		TenantID:    dataset.TenantID,
		DatasetID:   dataset.ID,
		DocumentID:  document.ID,
		SegmentID:   segment.ID,
		Position:    position,
		Content:     content,
		WordCount:   len([]rune(content)),
		Type:        "customized",
		IndexNodeID: uuid.NewV4().String(),
	}

	if err := dbengine.Instance().DB.Create(chunk).Error; err != nil {
		return nil, err
	}
	return chunk, nil
}

func (s *SegmentService) DeleteChildChunk(childChunkID string) error {
	return dbengine.Instance().DB.Where("id = ?", childChunkID).Delete(&models.ChildChunk{}).Error
}
