package embedding

type IEmbeddings interface {
	EmbedDocuments(texts []string) [][]float64
	EmbedQuery(text string) []float64
	AEmbedDocuments(texts []string) [][]float64
	AEmbedQuery(text string) []float64
}
