package retrieval

import "slices"

type RetrievalMethodType string

const (
	RetrievalMethod_SEMANTIC_SEARCH  RetrievalMethodType = "semantic_search"
	RetrievalMethod_FULL_TEXT_SEARCH RetrievalMethodType = "full_text_search"
	RetrievalMethod_HYBRID_SEARCH    RetrievalMethodType = "hybrid_search"
)

func (m RetrievalMethodType) IsSupportSemanticSearch() bool {
	return slices.Contains([]RetrievalMethodType{RetrievalMethod_SEMANTIC_SEARCH, RetrievalMethod_HYBRID_SEARCH}, m)
}
func (m RetrievalMethodType) IsSupportFulltextSearch() bool {
	return slices.Contains([]RetrievalMethodType{RetrievalMethod_FULL_TEXT_SEARCH, RetrievalMethod_HYBRID_SEARCH}, m)
}
