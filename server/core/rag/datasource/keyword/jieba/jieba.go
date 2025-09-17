// import json
// from collections import defaultdict
// from typing import Any, Optional

// from pydantic import BaseModel

// from configs import dify_config
// from core.rag.datasource.keyword.jieba.jieba_keyword_table_handler import JiebaKeywordTableHandler
// from core.rag.datasource.keyword.keyword_base import BaseKeyword
// from core.rag.models.document import Document
// from extensions.ext_database import db
// from extensions.ext_redis import redis_client
// from extensions.ext_storage import storage
// from models.dataset import Dataset, DatasetKeywordTable, DocumentSegment
package jieba

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"mlib.com/confy"
	"mlib.com/gofy/server/core/rag/datasource/keywordor"
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/models"
	distributelock "mlib.com/gofy/server/utils/distribute_lock"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

type KeywordTableConfig struct {
	MaxKeywordsPerChunk int `json:"max_keywords_per_chunk"` // 10
}

func NewKeywordTableConfig() *KeywordTableConfig {
	return &KeywordTableConfig{
		MaxKeywordsPerChunk: 10,
	}
}

type Jieba struct {
	*keywordor.BaseKeyword
	_config *KeywordTableConfig
}

func New(dataset *models.Dataset) *Jieba {
	return &Jieba{
		BaseKeyword: keywordor.NewBaseKeyword(dataset),
		_config:     NewKeywordTableConfig(),
	}
}

func (j *Jieba) _get_dataset_keyword_table() map[string][]string {
	dataset_keyword_table := j.Dataset.DatasetKeywordTable()
	if dataset_keyword_table != nil {
		keyword_table_dict := dataset_keyword_table.KeywordTableDict()
		if len(keyword_table_dict) > 0 {
			tmps := mapstruct.Get(mapstruct.Get(keyword_table_dict, "__data__", map[string]any{}), "table", map[string][]string{})
			for k, v := range tmps {
				slices.Sort(v)
				tmps[k] = slices.Compact(v)
			}
			return tmps
		}
	} else {
		keyword_data_source_type := confy.GetWithDefault[string]("keyword_data_source_type", "database")
		dataset_keyword_table := &models.DatasetKeywordTable{
			ID:             uuid.NewV4().String(),
			DatasetID:      j.Dataset.ID,
			KeywordTable:   "",
			DataSourceType: keyword_data_source_type,
		}
		if keyword_data_source_type == "database" {
			bindata, _ := json.Marshal(map[string]any{
				"__type__": "keyword_table",
				"__data__": map[string]any{"index_id": j.Dataset.ID, "summary": nil, "table": map[string]any{}},
			})
			dataset_keyword_table.KeywordTable = string(bindata)
		}
		dbengine.Instance().DB.Create(dataset_keyword_table)
	}
	return map[string][]string{}
}
func (j *Jieba) _update_segment_keywords(dataset_id string, node_id string, keywords []string) {
	document_segment := new(models.DocumentSegment)
	err := dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("dataset_id = ? and index_node_id = ?", dataset_id, node_id).First(document_segment).Error
	if err != nil {
		mlog.Errorf("get DocumentSegment failed:%v", err)
		document_segment = nil
	}
	if document_segment != nil {
		bindata, _ := json.Marshal(keywords)
		dbengine.Instance().DB.Updates(&models.DocumentSegment{ID: document_segment.ID, Keywords: datatypes.JSON(bindata)})
	}
}
func (j *Jieba) _add_text_to_keyword_table(keyword_table map[string][]string, id string, keywords []string) map[string][]string {
	for _, keyword := range keywords {
		if _, ok := keyword_table[keyword]; !ok {
			keyword_table[keyword] = make([]string, 0)
		}
		keyword_table[keyword] = append(keyword_table[keyword], id)
	}
	for k, v := range keyword_table {
		slices.Sort(v)
		keyword_table[k] = slices.Compact(v)
	}
	return keyword_table
}

func (j *Jieba) _save_dataset_keyword_table(keyword_table map[string][]string) {
	keyword_table_dict := map[string]any{
		"__type__": "keyword_table",
		"__data__": map[string]any{"index_id": j.Dataset.ID, "summary": nil, "table": keyword_table},
	}
	dataset_keyword_table := j.Dataset.DatasetKeywordTable()
	keyword_data_source_type := dataset_keyword_table.DataSourceType
	bindata, _ := json.Marshal(keyword_table_dict)
	if keyword_data_source_type == "database" {

		dbengine.Instance().DB.Updates(&models.DatasetKeywordTable{ID: dataset_keyword_table.ID, KeywordTable: string(bindata)})
	} else {
		file_key := filepath.Join("keyword_files", j.Dataset.TenantID, j.Dataset.ID+".txt")
		os.WriteFile(file_key, bindata, 0644)
	}
}
func (j *Jieba) _delete_ids_from_keyword_table(keyword_table map[string][]string, ids []string) map[string][]string {
	// get set of ids that correspond to node
	node_idxs_to_delete := ids

	// delete node_idxs from keyword to node idxs mapping
	keywords_to_delete := map[string]struct{}{}
	for keyword, node_idxs := range keyword_table {
		tmps_all := slices.Concat(node_idxs_to_delete, node_idxs)
		slices.Sort(tmps_all)
		tmps_all = slices.Compact(tmps_all)
		intersection := []string{}
		node_idxs_difference := []string{}
		for _, v := range tmps_all {
			if slices.Contains(node_idxs_to_delete, v) && slices.Contains(node_idxs, v) {
				intersection = append(intersection, v)
			}
			if slices.Contains(node_idxs, v) && !slices.Contains(node_idxs_to_delete, v) {
				node_idxs_difference = append(node_idxs_difference, v)
			}
		}
		if len(intersection) > 0 {
			keyword_table[keyword] = node_idxs_difference
			if len(node_idxs_difference) == 0 {
				keywords_to_delete[keyword] = struct{}{}
			}
		}
	}
	for keyword := range keywords_to_delete {
		delete(keyword_table, keyword)
	}
	return keyword_table
}
func (j *Jieba) _retrieve_ids_by_query(keyword_table map[string][]string, query string, k int) []string {
	if k <= 0 {
		k = 4
	}
	keyword_table_handler := NewJiebaKeywordTableHandler()
	keywords := keyword_table_handler.ExtractKeywords(query, 0)

	// go through text chunks in order of most matching keywords
	chunk_indices_count := map[string]int{}
	keywords_list := []string{}
	for _, keyword := range keywords {
		if _, ok := keyword_table[keyword]; ok {
			keywords_list = append(keywords_list, keyword)
		}
	}
	for _, keyword := range keywords_list {
		for _, node_id := range keyword_table[keyword] {
			if _, ok := chunk_indices_count[node_id]; !ok {
				chunk_indices_count[node_id] = 0
			}
			chunk_indices_count[node_id] += 1
		}
	}
	sorted_chunk_indices := slices.Sorted(maps.Keys(chunk_indices_count))
	sort.Slice(chunk_indices_count, func(i, j int) bool {
		return chunk_indices_count[sorted_chunk_indices[i]] > chunk_indices_count[sorted_chunk_indices[j]]
	})

	return sorted_chunk_indices[:k]
}
func (j *Jieba) Create(texts []*ragentities.Document, args map[string]any) keywordor.Keywordor {
	lock_name := "keyword_indexing_lock_" + j.Dataset.ID

	if distributelock.Instance().TryLock(lock_name, 600*time.Millisecond) {
		keyword_table_handler := NewJiebaKeywordTableHandler()

		keyword_table := j._get_dataset_keyword_table()
		for _, text := range texts {
			keywords := keyword_table_handler.ExtractKeywords(
				text.PageContent, j._config.MaxKeywordsPerChunk,
			)
			if len(text.Metadata) > 0 {
				j._update_segment_keywords(j.Dataset.ID, mapstruct.Get(text.Metadata, "doc_id", ""), keywords)
				keyword_table = j._add_text_to_keyword_table(
					keyword_table, mapstruct.Get(text.Metadata, "doc_id", ""), keywords,
				)
			}
		}

		j._save_dataset_keyword_table(keyword_table)
		distributelock.Instance().UnLock(lock_name)
	}
	return j
}

func (j *Jieba) AddTexts(texts []*ragentities.Document, args map[string]any) {
	lock_name := "keyword_indexing_lock_" + j.Dataset.ID
	if distributelock.Instance().TryLock(lock_name, 600*time.Millisecond) {
		keyword_table_handler := NewJiebaKeywordTableHandler()

		keyword_table := j._get_dataset_keyword_table()
		keywords_list := mapstruct.Get(args, "keywords_list", [][]string{})
		for i := range len(texts) {
			text := texts[i]
			var keywords []string
			if len(keywords_list) > 0 {
				if len(keywords_list)-1 >= i {
					keywords = keywords_list[i]
				}

				if keywords == nil {
					keywords = keyword_table_handler.ExtractKeywords(
						text.PageContent, j._config.MaxKeywordsPerChunk,
					)
				}
			} else {
				keywords = keyword_table_handler.ExtractKeywords(
					text.PageContent, j._config.MaxKeywordsPerChunk,
				)
			}
			if len(text.Metadata) > 0 {
				j._update_segment_keywords(j.Dataset.ID, mapstruct.Get(text.Metadata, "doc_id", ""), keywords)
				keyword_table = j._add_text_to_keyword_table(
					keyword_table, mapstruct.Get(text.Metadata, "doc_id", ""), keywords,
				)
			}
		}
		j._save_dataset_keyword_table(keyword_table)
		distributelock.Instance().UnLock(lock_name)
	}
}

func (j *Jieba) TextExists(id string) bool {
	keyword_table := j._get_dataset_keyword_table()
	if keyword_table == nil {
		return false
	}
	tmps := []string{}
	for v := range maps.Values(keyword_table) {
		tmps = append(tmps, v...)
	}
	slices.Sort(tmps)
	tmps = slices.Compact(tmps)
	return slices.Contains(tmps, id)
}

func (j *Jieba) DeleteByIDs(ids []string) {
	lock_name := "keyword_indexing_lock_" + j.Dataset.ID
	if distributelock.Instance().TryLock(lock_name, 600*time.Millisecond) {
		keyword_table := j._get_dataset_keyword_table()
		if len(keyword_table) > 0 {
			keyword_table = j._delete_ids_from_keyword_table(keyword_table, ids)
		}
		j._save_dataset_keyword_table(keyword_table)
		distributelock.Instance().UnLock(lock_name)
	}
}

func (j *Jieba) Search(query string, args map[string]any) []*ragentities.Document {
	keyword_table := j._get_dataset_keyword_table()

	k := mapstruct.Get(args, "top_k", 4)
	document_ids_filter := mapstruct.Get(args, "document_ids_filter", []string{})
	sorted_chunk_indices := j._retrieve_ids_by_query(keyword_table, query, k)

	documents := []*ragentities.Document{}
	for _, chunk_index := range sorted_chunk_indices {
		segment_query := dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("dataset_id = ? and index_node_id = ?", j.Dataset.ID, chunk_index)
		if len(document_ids_filter) > 0 {
			segment_query = segment_query.Where("document_id in ?", document_ids_filter)
		}
		segment := new(models.DocumentSegment)
		err := segment_query.First(segment).Error
		if err != nil {
			mlog.Errorf("get DocumentSegment failed:%v", err)
			segment = nil
		}

		if segment != nil {
			documents = append(documents, &ragentities.Document{
				PageContent: segment.Content,
				Metadata: map[string]any{
					"doc_id":      chunk_index,
					"doc_hash":    segment.IndexNodeHash,
					"document_id": segment.DocumentID,
					"dataset_id":  segment.DatasetID,
				},
			})
		}
	}
	return documents
}

func (j *Jieba) Delete() {
	lock_name := "keyword_indexing_lock_" + j.Dataset.ID
	if distributelock.Instance().TryLock(lock_name, 600*time.Millisecond) {
		dataset_keyword_table := j.Dataset.DatasetKeywordTable()
		if dataset_keyword_table != nil {
			dbengine.Instance().DB.Delete(dataset_keyword_table)
			if dataset_keyword_table.DataSourceType != "database" {
				file_key := filepath.Join("keyword_files", j.Dataset.TenantID, j.Dataset.ID+".txt")
				os.Remove(file_key)
			}
		}
		distributelock.Instance().UnLock(lock_name)
	}
}

// func (j *Jieba) create_segment_keywords(node_id  string, keywords  []string){
//         keyword_table = j._get_dataset_keyword_table()
//         j._update_segment_keywords(j.Dataset.ID, node_id, keywords)
//         keyword_table = j._add_text_to_keyword_table(keyword_table or {}, node_id, keywords)
//         j._save_dataset_keyword_table(keyword_table)
// }

// func (j *Jieba) multi_create_segment_keywords(pre_segment_data_list: list){
//         keyword_table_handler = JiebaKeywordTableHandler()
//         keyword_table = j._get_dataset_keyword_table()
//         for pre_segment_data in pre_segment_data_list{
//             segment = pre_segment_data["segment"]
//             if pre_segment_data["keywords"]{
//                 segment.keywords = pre_segment_data["keywords"]
//                 keyword_table = j._add_text_to_keyword_table(
//                     keyword_table or {}, segment.index_node_id, pre_segment_data["keywords"]
//                 )
//             else{
//                 keywords = keyword_table_handler.extract_keywords(segment.content, j._config.max_keywords_per_chunk)
//                 segment.keywords = list(keywords)
//                 keyword_table = j._add_text_to_keyword_table(
//                     keyword_table or {}, segment.index_node_id, list(keywords)
//                 )
//         j._save_dataset_keyword_table(keyword_table)
// }

// func (j *Jieba) update_segment_keywords_index(node_id  string, keywords  []string){
//         keyword_table = j._get_dataset_keyword_table()
//         keyword_table = j._add_text_to_keyword_table(keyword_table or {}, node_id, keywords)
//         j._save_dataset_keyword_table(keyword_table)// }

// class SetEncoder(json.JSONEncoder){
//     def default(obj){
//         if isinstance(obj, set){
//             return list(obj)
//         return super().default(obj)
