package retrieval

import (
	ragretrievalenumtypes "mlib.com/gofy/server/enum_types/rag/retrieval"
)

var (
	default_retrieval_model = map[string]any{
		"search_method":           ragretrievalenumtypes.RetrievalMethod_SEMANTIC_SEARCH,
		"reranking_enable":        false,
		"reranking_model":         map[string]any{"reranking_provider_name": "", "reranking_model_name": ""},
		"top_k":                   2,
		"score_threshold_enabled": false,
	}
)

type DatasetRetrieval struct {
}
