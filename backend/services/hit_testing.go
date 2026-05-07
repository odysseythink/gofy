package services

import (
	"fmt"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
)

type HitTestingService struct{}

func (s *HitTestingService) Retrieve(datasetID string, query string, topK int, scoreThreshold float64) ([]map[string]any, error) {
	// Simple keyword-based retrieval as fallback
	var segments []*models.DocumentSegment
	dbengine.Instance().DB.Where("dataset_id = ? AND enabled = ? AND status = ? AND content LIKE ?",
		datasetID, true, "completed", "%"+query+"%").
		Limit(topK).Find(&segments)

	results := make([]map[string]any, 0, len(segments))
	for _, seg := range segments {
		results = append(results, map[string]any{
			"segment_id": seg.ID,
			"content":    seg.Content,
			"score":      1.0,
			"word_count": seg.WordCount,
			"position":   seg.Position,
		})
	}
	return results, nil
}

func (s *HitTestingService) HitTestingArgsCheck(query string, topK int) error {
	if query == "" {
		return fmt.Errorf("query is required")
	}
	if topK <= 0 || topK > 100 {
		return fmt.Errorf("top_k must be between 1 and 100")
	}
	return nil
}
