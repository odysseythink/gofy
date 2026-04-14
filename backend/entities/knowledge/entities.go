package knowledge

import (
	"encoding/json"

	"github.com/odysseythink/mlog"
)

type ParentModeType string

const (
	ParentMode_FULL_DOC  ParentModeType = "full-doc"
	ParentMode_PARAGRAPH ParentModeType = "paragraph"
)

type NotionIcon struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Emoji string `json:"emoji"`
}

type NotionPage struct {
	PageID   string      `json:"page_id"`
	PageName string      `json:"page_name"`
	PageIcon *NotionIcon `json:"page_icon"`
	Type     string      `json:"type"`
}

type NotionInfo struct {
	WorkspaceID string       `json:"workspace_id"`
	Pages       []NotionPage `json:"pages"`
}

type WebsiteInfo struct {
	Provider        string   `json:"provider"`
	JobID           string   `json:"job_id"`
	URLs            []string `json:"urls"`
	OnlyMainContent bool     `json:"only_main_content"`
}

type FileInfo struct {
	FileIDs []string `json:"file_ids"`
}

type InfoList struct {
	DataSourceType  string        `json:"data_source_type"` // upload_file | notion_import | website_crawl
	NotionInfoList  *[]NotionInfo `json:"notion_info_list"`
	FileInfoList    *FileInfo     `json:"file_info_list"`
	WebsiteInfoList *WebsiteInfo  `json:"website_info_list"`
}

type DataSource struct {
	InfoList InfoList `json:"info_list"`
}

type PreProcessingRule struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type Segmentation struct {
	Separator    string `json:"separator"` //"\n"
	MaxTokens    int    `json:"max_tokens"`
	ChunkOverlap int    `json:"chunk_overlap"`
}

type Rule struct {
	PreProcessingRules   []PreProcessingRule `json:"pre_processing_rules"`
	Segmentation         *Segmentation       `json:"segmentation"`
	ParentMode           ParentModeType      `json:"parent_mode"`
	SubchunkSegmentation *Segmentation       `json:"subchunk_segmentation"`
}

func NewRule(args map[string]any) *Rule {
	if len(args) > 0 {
		bindata, _ := json.Marshal(args)
		rule := new(Rule)
		if err := json.Unmarshal(bindata, rule); err != nil {
			mlog.Errorf("json unmarshal %s failed:%v", string(bindata), err)
			return nil
		}
		return rule
	}
	return new(Rule)
}

type ProcessRule struct {
	Mode  string `json:"mode"` // automatic | custom | hierarchical
	Rules *Rule  `json:"rules"`
}

type RerankingModel struct {
	RerankingProviderName string `json:"reranking_provider_name"`
	RerankingModelName    string `json:"reranking_model_name"`
}

type WeightVectorSetting struct {
	VectorWeight          float64 `json:"vector_weight"`
	EmbeddingProviderName string  `json:"embedding_provider_name"`
	EmbeddingModelName    string  `json:"embedding_model_name"`
}

type WeightKeywordSetting struct {
	KeywordWeight float64 `json:"keyword_weight"`
}

type WeightModel struct {
	WeightType     string                `json:"weight_type"` // semantic_first | keyword_first | customized
	VectorSetting  *WeightVectorSetting  `json:"vector_setting"`
	KeywordSetting *WeightKeywordSetting `json:"keyword_setting"`
}

type RetrievalModel struct {
	SearchMethod          string          `json:"search_method"` // hybrid_search | semantic_search | full_text_search | keyword_search
	RerankingEnable       bool            `json:"reranking_enable"`
	RerankingModel        *RerankingModel `json:"reranking_model"`
	RerankingMode         string          `json:"reranking_mode"`
	TopK                  int             `json:"top_k"`
	ScoreThresholdEnabled bool            `json:"score_threshold_enabled"`
	ScoreThreshold        float64         `json:"score_threshold"`
	Weights               *WeightModel    `json:"weights"`
}

type MetaDataConfig struct {
	DocType     string          `json:"doc_type"`
	DocMetadata json.RawMessage `json:"doc_metadata"` // 任意 JSON 对象
}

type KnowledgeConfig struct {
	OriginalDocumentID     string          `json:"original_document_id"`
	Duplicate              bool            `json:"duplicate"`
	IndexingTechnique      string          `json:"indexing_technique"` // high_quality | economy
	DataSource             *DataSource     `json:"data_source"`
	ProcessRule            *ProcessRule    `json:"process_rule"`
	RetrievalModel         *RetrievalModel `json:"retrieval_model"`
	DocForm                string          `json:"doc_form"`     //"text_model"
	DocLanguage            string          `json:"doc_language"` // "English"
	EmbeddingModel         string          `json:"embedding_model"`
	EmbeddingModelProvider string          `json:"embedding_model_provider"`
	Name                   string          `json:"name"`
}

type SegmentUpdateArgs struct {
	Content               string   `json:"content"`
	Answer                string   `json:"answer"`
	Keywords              []string `json:"keywords"`
	RegenerateChildChunks bool     `json:"regenerate_child_chunks"`
	Enabled               bool     `json:"enabled"`
}

type ChildChunkUpdateArgs struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type MetadataArgs struct {
	Type string `json:"type"` // string | number | time
	Name string `json:"name"`
}

type MetadataUpdateArgs struct {
	Name  string `json:"name"`
	Value any    `json:"value"` // str | int | float
}

type MetadataDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value any    `json:"value"` // str | int | float
}

type DocumentMetadataOperation struct {
	DocumentID   string           `json:"document_id"`
	MetadataList []MetadataDetail `json:"metadata_list"`
}

type MetadataOperationData struct {
	OperationData []DocumentMetadataOperation `json:"operation_data"`
}
