package datasource

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	jiebakeyword "mlib.com/gofy/server/core/rag/datasource/keyword/jieba"
	ragentities "mlib.com/gofy/server/entities/rag"
	retrievalenumtypes "mlib.com/gofy/server/enum_types/rag/retrieval"
	"mlib.com/gofy/server/models"
)

type RetrievalService struct {
}

func escape_query_for_search(query string) string {
	return strings.ReplaceAll(query, `"`, `\"`)
}
func (s *RetrievalService) KeywordSearch(
	dataset *models.Dataset,
	query string,
	top_k int,
	exps []error,
	document_ids_filter []string,
) []*ragentities.Document {
	keyword := jiebakeyword.New(dataset)

	documents := keyword.Search(
		escape_query_for_search(query), map[string]any{"top_k": top_k, "document_ids_filter": document_ids_filter},
	)
	return documents
}

func (s *RetrievalService) EmbeddingSearch(
	dataset *models.Dataset,
	query string,
	top_k int,
        score_threshold float64,
        reranking_model map[string]any,
        retrieval_method string,
        exps []error,
        document_ids_filter []string,
    ){


                vector = Vector(dataset=dataset)
                documents = vector.search_by_vector(
                    query,
                    search_type="similarity_score_threshold",
                    top_k=top_k,
                    score_threshold=score_threshold,
                    filter={"group_id": [dataset.id]},
                    document_ids_filter=document_ids_filter,
                )

                if documents:
                    if (
                        reranking_model
                        and reranking_model.get("reranking_model_name")
                        and reranking_model.get("reranking_provider_name")
                        and retrieval_method == RetrievalMethod.SEMANTIC_SEARCH.value
                    ):
                        data_post_processor = DataPostProcessor(
                            str(dataset.tenant_id), str(RerankMode.RERANKING_MODEL.value), reranking_model, None, False
                        )
                        all_documents.extend(
                            data_post_processor.invoke(
                                query=query,
                                documents=documents,
                                score_threshold=score_threshold,
                                top_n=len(documents),
                            )
                        )
                    else:
                        all_documents.extend(documents)
}
}

// Retrieve 等价于 Python retrieve
func (s *RetrievalService) Retrieve(
	ctx context.Context,
	retrieval_method string,
	dataset_id string,
	query string,
	topK int,
	score_threshold float64,
	reranking_model map[string]any,
	reranking_mode string,
	weights map[string]any,
	document_ids_filter []string,
) []*ragentities.Document {
	if reranking_mode == "" {
		reranking_mode = "reranking_model"
	}
	if query == "" {
		return []*ragentities.Document{}
	}
	dataset := models.GetDataset(dataset_id)
	if dataset == nil {
		return []*ragentities.Document{}
	}

	threshold := score_threshold

	var (
		all_documents []*ragentities.Document
		mu            sync.Mutex // 保护 all_documents
		errGroup      *errgroup.Group
	)
	errGroup, ctx = errgroup.WithContext(ctx)

	// 1. keyword search
	if retrieval_method == "keyword_search" {
		errGroup.Go(func() error {
			docs := s.KeywordSearch(dataset, query, topK, nil, document_ids_filter)
			mu.Lock()
			if all_documents == nil {
				all_documents = make([]*ragentities.Document, 0)
			}
			all_documents = append(all_documents, docs...)
			mu.Unlock()
			return nil
		})
	}

	// 2. semantic search
	if retrievalenumtypes.RetrievalMethodType(retrieval_method).IsSupportSemanticSearch() {
		errGroup.Go(func() error {
			docs, e := s.embeddingSearch(ctx, dataset_id, query, topK, threshold, rerankingModel, documentIDsFilter)
			if e != nil {
				return e
			}
			mu.Lock()
			all_documents = append(all_documents, docs...)
			mu.Unlock()
			return nil
		})
	}

	// 3. full-text search
	if isSupportFulltextSearch(retrieval_method) {
		errGroup.Go(func() error {
			docs, e := s.fullTextIndexSearch(ctx, dataset_id, query, topK, threshold, rerankingModel, documentIDsFilter)
			if e != nil {
				return e
			}
			mu.Lock()
			all_documents = append(all_documents, docs...)
			mu.Unlock()
			return nil
		})
	}

	// 等待全部完成（或任一失败）
	if err := errGroup.Wait(); err != nil {
		return nil, fmt.Errorf("retrieval error: %w", err)
	}

	// 4. 混合搜索后处理
	if retrieval_method == "hybrid_search" {
		post := &DataPostProcessor{
			tenantID:       dataset.TenantID,
			rerankingMode:  rerankingMode,
			rerankingModel: rerankingModel,
			weights:        weights,
		}
		all_documents, err = post.Invoke(ctx, query, all_documents, threshold, topK)
		if err != nil {
			return nil, err
		}
	}

	return all_documents, nil
}
