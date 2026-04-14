package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odysseythink/mlog"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

// SegmentTaskPayload is the payload for segment-related tasks.
type SegmentTaskPayload struct {
	DatasetID  string   `json:"dataset_id"`
	DocumentID string   `json:"document_id"`
	SegmentIDs []string `json:"segment_ids"`
}

// HandleSegmentCreate indexes newly created segments.
func HandleSegmentCreate(ctx context.Context, task *Task) error {
	var payload SegmentTaskPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("indexing %d segments for document %s", len(payload.SegmentIDs), payload.DocumentID)

	// Load segments
	var segments []models.DocumentSegment
	dbengine.Instance().DB.Where("id IN ?", payload.SegmentIDs).Find(&segments)

	for _, seg := range segments {
		// TODO: Generate embedding vector using model_runtime
		// For now just mark as completed with index node
		dbengine.Instance().DB.Model(&seg).Updates(map[string]any{
			"status":        "completed",
			"enabled":       true,
			"index_node_id": seg.IndexNodeID,
			"tokens":        len([]rune(seg.Content)) / 2,
		})
	}

	return nil
}

// HandleSegmentDelete removes segments from the vector store.
func HandleSegmentDelete(ctx context.Context, task *Task) error {
	var payload SegmentTaskPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	mlog.Infof("removing %d segments from index", len(payload.SegmentIDs))

	// TODO: Remove from vector store by index_node_id
	// The actual vector store deletion depends on the vector DB implementation

	return nil
}

// HandleSegmentEnable enables segments in the vector store.
func HandleSegmentEnable(ctx context.Context, task *Task) error {
	var payload SegmentTaskPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	// TODO: Re-add to vector store
	dbengine.Instance().DB.Model(&models.DocumentSegment{}).
		Where("id IN ?", payload.SegmentIDs).
		Updates(map[string]any{"enabled": true, "disabled_at": nil})

	return nil
}

// HandleSegmentDisable disables segments in the vector store.
func HandleSegmentDisable(ctx context.Context, task *Task) error {
	var payload SegmentTaskPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	// TODO: Remove from vector store (but keep in DB)
	return nil
}
