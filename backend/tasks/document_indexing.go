package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

// DocumentIndexingPayload is the payload for document indexing tasks.
type DocumentIndexingPayload struct {
	DatasetID  string `json:"dataset_id"`
	DocumentID string `json:"document_id"`
}

// HandleDocumentIndexing processes a document indexing task.
func HandleDocumentIndexing(ctx context.Context, task *Task) error {
	var payload DocumentIndexingPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("indexing document %s in dataset %s", payload.DocumentID, payload.DatasetID)

	// Update document status to indexing
	now := time.Now()
	err := dbengine.Instance().DB.Model(&models.Document{}).
		Where("id = ? AND dataset_id = ?", payload.DocumentID, payload.DatasetID).
		Updates(map[string]any{
			"indexing_status":       "indexing",
			"processing_started_at": &now,
		}).Error
	if err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// Load the document
	var doc models.Document
	if err := dbengine.Instance().DB.Where("id = ?", payload.DocumentID).First(&doc).Error; err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// 1. Load the dataset
	var dataset models.Dataset
	if err := dbengine.Instance().DB.Where("id = ?", payload.DatasetID).First(&dataset).Error; err != nil {
		return fmt.Errorf("dataset not found: %w", err)
	}

	// 2. Get the process rule
	var processRule models.DatasetProcessRule
	dbengine.Instance().DB.Where("dataset_id = ?", payload.DatasetID).Order("created_at DESC").First(&processRule)

	// 3. Extract document content based on data source type
	// Document model has no inline content field; content comes from existing segments or extraction.
	// Collect content from any pre-existing segments, or log that extraction is needed.
	var existingSegments []models.DocumentSegment
	dbengine.Instance().DB.Where("document_id = ?", payload.DocumentID).Order("position ASC").Find(&existingSegments)

	var content string
	if len(existingSegments) > 0 {
		var parts []string
		for _, seg := range existingSegments {
			parts = append(parts, seg.Content)
		}
		content = strings.Join(parts, "\n\n")
	}
	if content == "" {
		// TODO: Use appropriate extractor based on doc.DataSourceType
		mlog.Infof("document %s has no inline content, needs extraction", payload.DocumentID)
	}

	// 4. Split content into segments
	segments := splitContent(content, processRule)

	// 5. Delete old segments before creating new ones
	dbengine.Instance().DB.Where("document_id = ?", payload.DocumentID).Delete(&models.DocumentSegment{})

	// 6. Create document segments in DB
	for i, segContent := range segments {
		seg := &models.DocumentSegment{
			ID:          uuid.NewV4().String(),
			TenantID:    dataset.TenantID,
			DatasetID:   dataset.ID,
			DocumentID:  doc.ID,
			Position:    i,
			Content:     segContent,
			WordCount:   len([]rune(segContent)),
			Status:      "completed",
			Enabled:     true,
			IndexNodeID: uuid.NewV4().String(),
		}
		dbengine.Instance().DB.Create(seg)
	}

	// 7. Update document status
	completedAt := time.Now()
	err = dbengine.Instance().DB.Model(&models.Document{}).Where("id = ?", payload.DocumentID).Updates(map[string]any{
		"indexing_status": "completed",
		"completed_at":    &completedAt,
		"word_count":      len([]rune(content)),
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update completion status: %w", err)
	}

	mlog.Infof("document %s indexing completed", payload.DocumentID)
	return nil
}

// HandleDocumentIndexingUpdate re-indexes an existing document.
func HandleDocumentIndexingUpdate(ctx context.Context, task *Task) error {
	var payload DocumentIndexingPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("updating index for document %s", payload.DocumentID)

	// Mark document as re-indexing
	now := time.Now()
	dbengine.Instance().DB.Model(&models.Document{}).
		Where("id = ? AND dataset_id = ?", payload.DocumentID, payload.DatasetID).
		Updates(map[string]any{
			"indexing_status":       "indexing",
			"processing_started_at": &now,
		})

	// TODO: Delete old segments and embeddings, then re-index
	// For now, mark as completed
	dbengine.Instance().DB.Model(&models.Document{}).
		Where("id = ?", payload.DocumentID).
		Updates(map[string]any{
			"indexing_status": "completed",
			"completed_at":    &now,
		})

	return nil
}

// HandleDocumentClean removes a document's indexed data.
func HandleDocumentClean(ctx context.Context, task *Task) error {
	var payload DocumentIndexingPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("cleaning document %s data", payload.DocumentID)

	// Delete segments
	dbengine.Instance().DB.Where("document_id = ?", payload.DocumentID).Delete(&models.DocumentSegment{})

	// Delete child chunks
	dbengine.Instance().DB.Where("document_id = ?", payload.DocumentID).Delete(&models.ChildChunk{})

	// TODO: Delete from vector store

	return nil
}

// splitContent splits document content into segments based on the process rule.
func splitContent(content string, rule models.DatasetProcessRule) []string {
	if content == "" {
		return nil
	}

	// Default separator and max tokens
	separator := "\n\n"
	maxTokens := 500

	// Parse process rule
	if rule.Mode == "custom" && rule.Rules != "" {
		var rules map[string]any
		if err := json.Unmarshal([]byte(rule.Rules), &rules); err == nil {
			if seg, ok := rules["segmentation"].(map[string]any); ok {
				if sep, ok := seg["separator"].(string); ok && sep != "" {
					separator = sep
				}
				if mt, ok := seg["max_tokens"].(float64); ok && mt > 0 {
					maxTokens = int(mt)
				}
			}
		}
	}

	// Split by separator
	parts := strings.Split(content, separator)

	// Merge small segments, split large ones
	var segments []string
	var current strings.Builder
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if current.Len() > 0 && current.Len()+len(part) > maxTokens*4 { // rough char estimate
			segments = append(segments, current.String())
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteString(separator)
		}
		current.WriteString(part)
	}
	if current.Len() > 0 {
		segments = append(segments, current.String())
	}

	return segments
}
