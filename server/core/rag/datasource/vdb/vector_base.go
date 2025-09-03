package vdb

import (
	ragentities "mlib.com/gofy/server/entities/rag"
)

type BaseVector struct {
	CollectionNname string `json:"collection_name"`
}

func (bv *BaseVector) _filter_duplicate_texts(v IVector, texts []*ragentities.Document) []*ragentities.Document {
	new_text := []*ragentities.Document{}
	for _, text := range texts {
		if len(text.Metadata) > 0 {
			if _, ok := text.Metadata["doc_id"]; ok {
				if _, ok := text.Metadata["doc_id"].(string); ok {
					doc_id := text.Metadata["doc_id"].(string)
					exists_duplicate_node := v.TextExists(doc_id)
					if !exists_duplicate_node {
						new_text = append(new_text, text)
					}
				}
			}
		}
	}
	return new_text
}

func (bv *BaseVector) _get_uuids(texts []*ragentities.Document) []string {
	ids := []string{}
	for _, text := range texts {
		if len(text.Metadata) > 0 {
			if _, ok := text.Metadata["doc_id"]; ok {
				if _, ok := text.Metadata["doc_id"].(string); ok {
					ids = append(ids, text.Metadata["doc_id"].(string))
				}
			}
		}
	}
	return ids
}

type IVector interface {
	GetType() string
	Create(texts []*ragentities.Document, embeddings [][]float64, kwargs ...any)
	AddTexts(documents []*ragentities.Document, embeddings [][]float64, kwargs ...any)
	TextExists(id string) bool
	DeleteByIDs(ids []string)
	GetIDsByMetadataField(key string, value string)
	DeleteByMetadataField(key string, value string)
	SearchByVector(query_vector []float64, kwargs ...any) []*ragentities.Document
	SearchByFullText(query string, kwargs ...any) []*ragentities.Document
	Delete()
}
