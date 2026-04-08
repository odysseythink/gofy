package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
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

	// TODO: Call actual indexing pipeline
	// 1. Extract text from document (using extractors)
	// 2. Split into segments (using splitter)
	// 3. Generate embeddings (using embedding model)
	// 4. Store in vector database
	// For now, mark as completed

	err = dbengine.Instance().DB.Model(&models.Document{}).
		Where("id = ?", payload.DocumentID).
		Updates(map[string]any{
			"indexing_status": "completed",
			"completed_at":   &now,
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
			"completed_at":   &now,
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
