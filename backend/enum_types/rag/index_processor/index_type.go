package indexprocessor

type IndexType string

const (
	Index_PARAGRAPH_INDEX    IndexType = "text_model"
	Index_QA_INDEX           IndexType = "qa_model"
	Index_PARENT_CHILD_INDEX IndexType = "hierarchical_model"
)
