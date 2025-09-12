package rag

type RetrievalSourceMetadata struct {
	Position        int            `json:"position"`
	DatasetID       string         `json:"dataset_id"`
	DatasetName     string         `json:"dataset_name"`
	DocumentID      string         `json:"document_id"`
	DocumentName    string         `json:"document_name"`
	DataSourceType  string         `json:"data_source_type"`
	SegmentID       string         `json:"segment_id"`
	RetrieverFrom   string         `json:"retriever_from"`
	Score           float64        `json:"score"`
	HitCount        int            `json:"hit_count"`
	WordCount       int            `json:"word_count"`
	SegmentPosition int            `json:"segment_position"`
	IndexNodeHash   string         `json:"index_node_hash"`
	Content         string         `json:"content"`
	Page            int            `json:"page"`
	DocMetadata     map[string]any `json:"doc_metadata"`
	Title           string         `json:"title"`
}
