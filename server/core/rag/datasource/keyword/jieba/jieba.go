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
	"maps"
	"regexp"
	"slices"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
	"mlib.com/confy"
	"mlib.com/gofy/server/core/rag/datasource/keywordor"
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/models"
	distributelock "mlib.com/gofy/server/utils/distribute_lock"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)
type KeywordTableConfig struct{
    MaxKeywordsPerChunk int `json:"max_keywords_per_chunk"` // 10
}

func NewKeywordTableConfig()*KeywordTableConfig{
    return &KeywordTableConfig{
        MaxKeywordsPerChunk: 10,
    }
}

type Jieba struct{
    *keywordor.BaseKeyword
    _config *KeywordTableConfig
}
func New(dataset *models.Dataset) *Jieba{
    return &Jieba{
        BaseKeyword: keywordor.NewBaseKeyword(dataset),
        _config: NewKeywordTableConfig(),
    }
}
func (j *Jieba) _get_dataset_keyword_table() map[string][]string{
        dataset_keyword_table := j.Dataset.DatasetKeywordTable()
        if dataset_keyword_table != nil{
            keyword_table_dict := dataset_keyword_table.KeywordTableDict()
            if len(keyword_table_dict) > 0{
                return mapstruct.Get(mapstruct.Get(keyword_table_dict, "__data__", map[string]any{}), "table",map[string][]string{})
            }
        }else{
            keyword_data_source_type := confy.GetWithDefault[string]("keyword_data_source_type", "database")
            dataset_keyword_table := &models.DatasetKeywordTable{
                ID: uuid.NewV4().String(),
                DatasetID:j.Dataset.ID,
                KeywordTable:"",
                DataSourceType:keyword_data_source_type,
            }
            if keyword_data_source_type == "database"{
                bindata, _ := json.Marshal(map[string]any{
                     "__type__": "keyword_table",
                        "__data__": map[string]any{"index_id": j.Dataset.ID, "summary": nil, "table": map[string]any{}},
                    })
                    dataset_keyword_table.KeywordTable= string(bindata)
                }
            dbengine.Instance().DB.Create(dataset_keyword_table)
            }
        return map[string]any{}
}
func (j *Jieba) _update_segment_keywords(dataset_id  string, node_id  string, keywords  []string){
        document_segment := new(models.DocumentSegment)
        err := dbengine.Instance().DB.Model(&models.DocumentSegment{}).Where("dataset_id = ? and index_node_id = ?", dataset_id, node_id).First(document_segment).Error
        if err != nil {
            mlog.Errorf("get DocumentSegment failed:%v", err)
            document_segment = nil
        }
        if document_segment != nil{
            dbengine.Instance().DB.Updates(&models.DocumentSegment{ID: document_segment.ID, Keywords: datatypes.JSON(keywords)})
        }
}
func (j *Jieba) _add_text_to_keyword_table(keyword_table  map[string]map[string]struct{}, id  string, keywords  []string)   map[string]map[string]struct{}{
        for _,  keyword := range keywords{
            if _, ok := keyword_table[keyword]; !ok {
                keyword_table[keyword] = make(map[string]struct{})
            }
            keyword_table[keyword][id]=struct{}{}
        }
        return keyword_table

}
func (j *Jieba) Create(texts []*ragentities.Document, args ...any) keywordor.Keywordor{
        lock_name := "keyword_indexing_lock_"+j.Dataset.ID
        
        if distributelock.Instance().TryLock(lock_name, 600*time.Miliseconds){
            keyword_table_handler := NewJiebaKeywordTableHandler()
            keyword_table := j._get_dataset_keyword_table()
            for _, text := range texts{
                keywords := keyword_table_handler.ExtractKeywords(
                    text.PageContent, j._config.MaxKeywordsPerChunk,
                )
                if len(text.Metadata) > 0{
                   
                    j._update_segment_keywords(j.Dataset.ID, mapstruct.Get(text.Metadata,"doc_id",""), slices.Sorted( maps.Keys(keywords)))
                    keyword_table = j._add_text_to_keyword_table(
                        keyword_table or {}, text.Metadata["doc_id"], list(keywords)
                    )
                }
            }

            j._save_dataset_keyword_table(keyword_table)

            return j
        }

}
// func (j *Jieba) add_texts(texts: list[Document], **kwargs){
//         lock_name := "keyword_indexing_lock_"+j.Dataset.ID
//         with redis_client.lock(lock_name, timeout=600){
//             keyword_table_handler = JiebaKeywordTableHandler()

//             keyword_table = j._get_dataset_keyword_table()
//             keywords_list = kwargs.get("keywords_list")
//             for i in range(len(texts)){
//                 text = texts[i]
//                 if keywords_list{
//                     keywords = keywords_list[i]
//                     if not keywords{
//                         keywords = keyword_table_handler.extract_keywords(
//                             text.page_content, j._config.max_keywords_per_chunk
//                         )
//                 else{
//                     keywords = keyword_table_handler.extract_keywords(
//                         text.page_content, j._config.max_keywords_per_chunk
//                     )
//                 if text.metadata is not None{
//                     j._update_segment_keywords(j.Dataset.ID, text.metadata["doc_id"], list(keywords))
//                     keyword_table = j._add_text_to_keyword_table(
//                         keyword_table or {}, text.metadata["doc_id"], list(keywords)
//                     )

//             j._save_dataset_keyword_table(keyword_table)

// }
// func (j *Jieba) text_exists(id  string) -> bool{
//         keyword_table = j._get_dataset_keyword_table()
//         if keyword_table is None{
//             return False
//         return id in set.union(*keyword_table.values())

// }
// func (j *Jieba) delete_by_ids(ids  []string) -> None{
//         lock_name := "keyword_indexing_lock_"+j.Dataset.ID
//         with redis_client.lock(lock_name, timeout=600){
//             keyword_table = j._get_dataset_keyword_table()
//             if keyword_table is not None{
//                 keyword_table = j._delete_ids_from_keyword_table(keyword_table, ids)

//             j._save_dataset_keyword_table(keyword_table)

// }
// func (j *Jieba) search(query  string, **kwargs: Any) -> list[Document]{
//         keyword_table = j._get_dataset_keyword_table()

//         k = kwargs.get("top_k", 4)
//         document_ids_filter = kwargs.get("document_ids_filter")
//         sorted_chunk_indices = j._retrieve_ids_by_query(keyword_table or {}, query, k)

//         documents = []
//         for chunk_index in sorted_chunk_indices{
//             segment_query = db.session.query(DocumentSegment).where(
//                 DocumentSegment.dataset_id == j.Dataset.ID, DocumentSegment.index_node_id == chunk_index
//             )
//             if document_ids_filter{
//                 segment_query = segment_query.where(DocumentSegment.document_id.in_(document_ids_filter))
//             segment = segment_query.first()

//             if segment{
//                 documents.append(
//                     Document(
//                         page_content=segment.content,
//                         metadata={
//                             "doc_id": chunk_index,
//                             "doc_hash": segment.index_node_hash,
//                             "document_id": segment.document_id,
//                             "dataset_id": segment.dataset_id,
//                         },
//                     )
//                 )

//         return documents

// }
// func (j *Jieba) delete(self) -> None{
//         lock_name := "keyword_indexing_lock_"+j.Dataset.ID
//         with redis_client.lock(lock_name, timeout=600){
//             dataset_keyword_table = j.Dataset.dataset_keyword_table
//             if dataset_keyword_table{
//                 db.session.delete(dataset_keyword_table)
//                 db.session.commit()
//                 if dataset_keyword_table.data_source_type != "database"{
//                     file_key = "keyword_files/" + j.Dataset.tenant_id + "/" + j.Dataset.ID + ".txt"
//                     storage.delete(file_key)

// }
// func (j *Jieba) _save_dataset_keyword_table(keyword_table){
//         keyword_table_dict = {
//             "__type__": "keyword_table",
//             "__data__": {"index_id": j.Dataset.ID, "summary": None, "table": keyword_table},
//         }
//         dataset_keyword_table = j.Dataset.dataset_keyword_table
//         keyword_data_source_type = dataset_keyword_table.data_source_type
//         if keyword_data_source_type == "database"{
//             dataset_keyword_table.keyword_table = json.dumps(keyword_table_dict, cls=SetEncoder)
//             db.session.commit()
//         else{
//             file_key = "keyword_files/" + j.Dataset.tenant_id + "/" + j.Dataset.ID + ".txt"
//             if storage.exists(file_key){
//                 storage.delete(file_key)
//             storage.save(file_key, json.dumps(keyword_table_dict, cls=SetEncoder).encode("utf-8"))

// }


// func (j *Jieba) _delete_ids_from_keyword_table(keyword_table  map[string]any, ids  []string)   map[string]any{
//         // get set of ids that correspond to node
//         node_idxs_to_delete = set(ids)

//         // delete node_idxs from keyword to node idxs mapping
//         keywords_to_delete = set()
//         for keyword, node_idxs in keyword_table.items(){
//             if node_idxs_to_delete.intersection(node_idxs){
//                 keyword_table[keyword] = node_idxs.difference(node_idxs_to_delete)
//                 if not keyword_table[keyword]{
//                     keywords_to_delete.add(keyword)

//         for keyword in keywords_to_delete{
//             del keyword_table[keyword]

//         return keyword_table

// }
// func (j *Jieba) _retrieve_ids_by_query(keyword_table  map[string]any, query  string, k: int = 4){
//         keyword_table_handler = JiebaKeywordTableHandler()
//         keywords = keyword_table_handler.extract_keywords(query)

//         // go through text chunks in order of most matching keywords
//         chunk_indices_count: dict[str, int] = defaultdict(int)
//         keywords_list = [keyword for keyword in keywords if keyword in set(keyword_table.keys())]
//         for keyword in keywords_list{
//             for node_id in keyword_table[keyword]{
//                 chunk_indices_count[node_id] += 1

//         sorted_chunk_indices = sorted(
//             chunk_indices_count.keys(),
//             key=lambda x: chunk_indices_count[x],
//             reverse=True,
//         )

//         return sorted_chunk_indices[:k]

// }

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
//         j._save_dataset_keyword_table(keyword_table)
// }

// class SetEncoder(json.JSONEncoder){
//     def default(obj){
//         if isinstance(obj, set){
//             return list(obj)
//         return super().default(obj)
