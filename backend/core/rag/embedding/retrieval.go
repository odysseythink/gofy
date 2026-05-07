package embedding

import "github.com/odysseythink/gofy/backend/models"

type RetrievalChildChunk struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Score    float64 `json:"score"`
	Position int     `json:"position"`
}

type RetrievalSegments struct {
	ModelConfig map[string]any         `json:"model_config"` //{"arbitrary_types_allowed": True}
	Segment     models.DocumentSegment `json:"segment"`
	ChildChunks []*RetrievalChildChunk `json:"child_chunks"`
	Score       float64                `json:"score"`
}

func NewRetrievalSegments() *RetrievalSegments {
	return &RetrievalSegments{}
}
