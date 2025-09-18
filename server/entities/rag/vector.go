package rag

type IVector interface {
	GetType() string
	Create(texts []*Document, embeddings [][]float64, kwargs ...any)
	AddTexts(documents []*Document, embeddings [][]float64, kwargs ...any)
	TextExists(id string) bool
	DeleteByIDs(ids []string)
	GetIDsByMetadataField(key string, value string)
	DeleteByMetadataField(key string, value string)
	SearchByVector(query_vector []float64, kwargs ...any) []*Document
	SearchByFullText(query string, kwargs ...any) []*Document
	Delete()
	FilterDuplicateTexts(v IVector, texts []*Document) []*Document
	GetUUIDs(texts []*Document) []string
}
type BaseVector struct {
	CollectionNname string `json:"collection_name"`
}

func (bv *BaseVector) FilterDuplicateTexts(v IVector, texts []*Document) []*Document {
	new_text := []*Document{}
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

func (bv *BaseVector) GetUUIDs(texts []*Document) []string {
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
