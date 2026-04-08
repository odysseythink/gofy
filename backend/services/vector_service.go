package services

import (
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

type VectorService struct{}

func (s *VectorService) CreateSegmentsVector(datasetID string, segmentIDs []string) error {
	// TODO: Generate embeddings and store in vector DB
	for _, segID := range segmentIDs {
		dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("id = ?", segID).Updates(map[string]any{
			"status":  "completed",
			"enabled": true,
		})
	}
	return nil
}

func (s *VectorService) UpdateSegmentVector(segmentID string, datasetID string) error {
	// TODO: Re-generate embedding for updated segment
	return nil
}

func (s *VectorService) DeleteSegmentVector(segmentID string, datasetID string) error {
	// TODO: Remove from vector DB
	return nil
}

func (s *VectorService) CreateChildChunkVector(childChunkID string, datasetID string) error {
	// TODO: Generate embedding for child chunk
	return nil
}

func (s *VectorService) DeleteChildChunkVector(childChunkID string, datasetID string) error {
	// TODO: Remove child chunk from vector DB
	return nil
}
