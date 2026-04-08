package vdb

type FieldType string

const (
	Field_CONTENT_KEY  FieldType = "page_content"
	Field_METADATA_KEY FieldType = "metadata"
	Field_GROUP_KEY    FieldType = "group_id"
	Field_VECTOR       FieldType = "vector"
	// Sparse Vector aims to support full text search
	Field_SPARSE_VECTOR FieldType = "sparse_vector"
	Field_TEXT_KEY      FieldType = "text"
	Field_PRIMARY_KEY   FieldType = "id"
	Field_DOC_ID        FieldType = "metadata.doc_id"
	Field_DOCUMENT_ID   FieldType = "metadata.document_id"
)
