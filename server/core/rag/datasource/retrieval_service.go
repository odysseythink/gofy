package datasource

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
	"mlib.com/gofy/server/core/exceptions"
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/models"
)

type RetrievalService struct {
}

func (s *RetrievalService) keyword_search(
        dataset *models.Dataset,
        query string,
        top_k int,
        all_documents []*ragentities.Document,
        exps []error,
        document_ids_filter []string,
    ){
                keyword = Keyword(dataset=dataset)

                documents = keyword.search(
                    cls.escape_query_for_search(query), top_k=top_k, document_ids_filter=document_ids_filter
                )
                all_documents.extend(documents)
            except Exception as e:
                exceptions.append(str(e))
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
			docs, e := s.keywordSearch(ctx, dataset_id, query, topK, documentIDsFilter)
			if e != nil {
				return e
			}
			mu.Lock()
			all_documents = append(all_documents, docs...)
			mu.Unlock()
			return nil
		})
	}

	// 2. semantic search
	if isSupportSemanticSearch(retrieval_method) {
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
