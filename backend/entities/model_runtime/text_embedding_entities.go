package modelruntime

import "github.com/shopspring/decimal"

// EmbeddingUsage represents usage information for embeddings.
type EmbeddingUsage struct {
	*ModelUsage
	Tokens      int
	TotalTokens int
	UnitPrice   decimal.Decimal
	PriceUnit   decimal.Decimal
	TotalPrice  decimal.Decimal
	Currency    string
	Latency     float64
}

// TextEmbeddingResult represents the result of text embedding.
type TextEmbeddingResult struct {
	Model      string
	Embeddings [][]float64
	Usage      *EmbeddingUsage
}
