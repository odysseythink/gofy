package modelruntime

// RerankDocument represents a document for reranking.
type RerankDocument struct {
	Index int
	Text  string
	Score float64
}

// RerankResult represents the result of reranking.
type RerankResult struct {
	Model string
	Docs  []*RerankDocument
}
