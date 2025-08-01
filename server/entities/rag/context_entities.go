package rag

type DocumentContext struct {
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}
