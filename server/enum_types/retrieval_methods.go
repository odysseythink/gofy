package enumtypes

import "slices"

type RetrievalMethod string

const (
	RetrievalMethod_SEMANTIC_SEARCH  RetrievalMethod = "semantic_search"
	RetrievalMethod_FULL_TEXT_SEARCH RetrievalMethod = "full_text_search"
	RetrievalMethod_HYBRID_SEARCH    RetrievalMethod = "hybrid_search"
)

func (m RetrievalMethod) IsSupportSemanticSearch() bool {
	return slices.Contains([]RetrievalMethod{RetrievalMethod_SEMANTIC_SEARCH, RetrievalMethod_HYBRID_SEARCH}, m)
}
func (m RetrievalMethod) IsSupportFulltextSearch() bool {
	return slices.Contains([]RetrievalMethod{RetrievalMethod_FULL_TEXT_SEARCH, RetrievalMethod_HYBRID_SEARCH}, m)
}
