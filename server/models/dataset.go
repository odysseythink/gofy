package models

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"gorm.io/datatypes"
	dbengine "mlib.com/gofy/server/db_engine"
	enumtypes "mlib.com/gofy/server/enum_types"
	"mlib.com/mlog"
)

// Dataset [...]
type Dataset struct {
	ID                     string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID               string         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant                 *Tenant        `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	Name                   string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description            string         `gorm:"column:description;type:text" json:"description"`
	Provider               string         `gorm:"column:provider;type:varchar(255);default:vendor" json:"provider"`
	Permission             string         `gorm:"column:permission;type:varchar(255);default:only_me" json:"permission"`
	DataSourceType         string         `gorm:"column:data_source_type;type:varchar(255)" json:"data_source_type"`
	IndexingTechnique      string         `gorm:"column:indexing_technique;type:varchar(255)" json:"indexing_technique"`
	IndexStruct            string         `gorm:"column:index_struct;type:text" json:"index_struct"`
	CreatedBy              string         `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedByAccount       *Account       `json:"created_by_account" form:"created_by_account" gorm:"foreignKey:CreatedBy;references:ID;"`
	CreatedAt              *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy              string         `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UpdatedAt              *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	EmbeddingModel         string         `gorm:"column:embedding_model;type:varchar(255);default:text-embedding-ada-002" json:"embedding_model"`
	EmbeddingModelProvider string         `gorm:"column:embedding_model_provider;type:varchar(255);default:openai" json:"embedding_model_provider"`
	CollectionBindingID    string         `gorm:"column:collection_binding_id;type:varchar(36)" json:"collection_binding_id"`
	RetrievalModel         datatypes.JSON `gorm:"column:retrieval_model;type:json" json:"retrieval_model"`
}

// TableName get sql table name.获取数据库表名
func (Dataset) TableName() string {
	return "datasets"
}
func (ds *Dataset) IndexStructDict() map[any]any {
	if ds.IndexStruct == "" {
		return nil
	}
	dict := map[any]any{}
	err := json.Unmarshal([]byte(ds.IndexStruct), &dict)
	if err != nil {
		mlog.Warningf("Dataset(%#v) IndexStruct must be dict", ds)
		return nil
	}
	return dict
}
func (ds *Dataset) ExternalRetrievalModel() map[string]any {
	default_retrieval_model := map[string]any{
		"top_k":           2,
		"score_threshold": 0.0,
	}
	if ds.RetrievalModel.String() == "" {
		return default_retrieval_model
	}
	retrieval_model := map[string]any{}
	err := json.Unmarshal([]byte(ds.RetrievalModel.String()), &retrieval_model)
	if err != nil {
		mlog.Warningf("Dataset RetrievalModel(%s) Unmarshal to dict failed:%v", ds.RetrievalModel.String(), err)
		return default_retrieval_model
	}
	return retrieval_model

}

func (ds *Dataset) RetrievalModelDict() map[string]any {
	default_retrieval_model := map[string]any{
		"search_method":           enumtypes.RetrievalMethod_SEMANTIC_SEARCH,
		"reranking_enable":        false,
		"reranking_model":         map[string]any{"reranking_provider_name": "", "reranking_model_name": ""},
		"top_k":                   2,
		"score_threshold_enabled": false,
	}

	retrieval_model := map[string]any{}
	err := json.Unmarshal([]byte(ds.RetrievalModel.String()), &retrieval_model)
	if err != nil {
		mlog.Warningf("Dataset RetrievalModel(%s) Unmarshal to dict failed:%v", ds.RetrievalModel.String(), err)
		return default_retrieval_model
	}
	return retrieval_model
	// return self.retrieval_model or default_retrieval_model

}
func (ds *Dataset) GenCollectionNameByID(dataset_id string) string {
	normalized_dataset_id := strings.ReplaceAll(dataset_id, "-", "_")
	return fmt.Sprintf("Vector_index_%s_Node", normalized_dataset_id)
}

func (ds *Dataset) ToDict() map[string]any {
	bindata, _ := json.Marshal(ds)
	var tmp_dict map[string]any
	err := json.Unmarshal(bindata, &tmp_dict)
	if err != nil {
		mlog.Errorf("json unmarshal=%s to dict failed:%v", string(bindata), err)
	}
	return tmp_dict
}

var (
	MODES                = []string{"automatic", "custom", "hierarchical"}
	PRE_PROCESSING_RULES = []string{"remove_stopwords", "remove_extra_spaces", "remove_urls_emails"}
	AUTOMATIC_RULES      = map[string]any{
		"pre_processing_rules": []map[string]any{
			{"id": "remove_extra_spaces", "enabled": true},
			{"id": "remove_urls_emails", "enabled": false},
		},
		"segmentation": map[string]any{"delimiter": "\n", "max_tokens": 500, "chunk_overlap": 50},
	}
)

// DatasetProcessRule [...]
type DatasetProcessRule struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	DatasetID string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	Mode      string     `gorm:"column:mode;type:varchar(255);default:automatic" json:"mode"`
	Rules     string     `gorm:"column:rules;type:text" json:"rules"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (rule *DatasetProcessRule) RulesDict() map[string]any {
	if rule.Rules != "" {
		var tmp_dict map[string]any
		err := json.Unmarshal([]byte(rule.Rules), &tmp_dict)
		if err != nil {
			mlog.Errorf("json unmarshal=%s to dict failed:%v", rule.Rules, err)
		}
		return tmp_dict
	}
	return nil
}

func (rule *DatasetProcessRule) ToDict() map[string]any {
	return map[string]any{
		"id":         rule.ID,
		"dataset_id": rule.DatasetID,
		"mode":       rule.Mode,
		"rules":      rule.RulesDict(),
	}
}

// TableName get sql table name.获取数据库表名
func (DatasetProcessRule) TableName() string {
	return "dataset_process_rules"
}

// Document [...]
type Document struct {
	ID                   string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID             string         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	DatasetID            string         `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	Position             int            `gorm:"column:position;type:int;not null" json:"position"`
	DataSourceType       string         `gorm:"column:data_source_type;type:varchar(255);not null" json:"data_source_type"`
	DataSourceInfo       string         `gorm:"column:data_source_info;type:text" json:"data_source_info"`
	DatasetProcessRuleID string         `gorm:"column:dataset_process_rule_id;type:varchar(36)" json:"dataset_process_rule_id"`
	Batch                string         `gorm:"column:batch;type:varchar(255);not null" json:"batch"`
	Name                 string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedFrom          string         `gorm:"column:created_from;type:varchar(255);not null" json:"created_from"`
	CreatedBy            string         `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAPIRequestID  string         `gorm:"column:created_api_request_id;type:varchar(36)" json:"created_api_request_id"`
	CreatedAt            *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ProcessingStartedAt  *time.Time     `gorm:"column:processing_started_at;type:timestamp" json:"processing_started_at"`
	FileID               string         `gorm:"column:file_id;type:text" json:"file_id"`
	WordCount            int            `gorm:"column:word_count;type:int" json:"word_count"`
	ParsingCompletedAt   *time.Time     `gorm:"column:parsing_completed_at;type:timestamp" json:"parsing_completed_at"`
	CleaningCompletedAt  *time.Time     `gorm:"column:cleaning_completed_at;type:timestamp" json:"cleaning_completed_at"`
	SplittingCompletedAt *time.Time     `gorm:"column:splitting_completed_at;type:timestamp" json:"splitting_completed_at"`
	Tokens               int            `gorm:"column:tokens;type:int" json:"tokens"`
	IndexingLatency      float64        `gorm:"column:indexing_latency;type:double" json:"indexing_latency"`
	CompletedAt          *time.Time     `gorm:"column:completed_at;type:timestamp" json:"completed_at"`
	IsPaused             bool           `gorm:"column:is_paused;type:tinyint(1);default:0" json:"is_paused"`
	PausedBy             string         `gorm:"column:paused_by;type:varchar(36)" json:"paused_by"`
	PausedAt             *time.Time     `gorm:"column:paused_at;type:timestamp" json:"paused_at"`
	Error                string         `gorm:"column:error;type:text" json:"error"`
	StoppedAt            *time.Time     `gorm:"column:stopped_at;type:timestamp" json:"stopped_at"`
	IndexingStatus       string         `gorm:"column:indexing_status;type:varchar(255);default:waiting" json:"indexing_status"`
	Enabled              bool           `gorm:"column:enabled;type:tinyint(1);not null;default:1" json:"enabled"`
	DisabledAt           *time.Time     `gorm:"column:disabled_at;type:timestamp" json:"disabled_at"`
	DisabledBy           string         `gorm:"column:disabled_by;type:varchar(36)" json:"disabled_by"`
	Archived             bool           `gorm:"column:archived;type:tinyint(1);not null;default:0" json:"archived"`
	ArchivedReason       string         `gorm:"column:archived_reason;type:varchar(255)" json:"archived_reason"`
	ArchivedBy           string         `gorm:"column:archived_by;type:varchar(36)" json:"archived_by"`
	ArchivedAt           *time.Time     `gorm:"column:archived_at;type:timestamp" json:"archived_at"`
	UpdatedAt            *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DocType              string         `gorm:"column:doc_type;type:varchar(40)" json:"doc_type"`
	DocMetadata          datatypes.JSON `gorm:"column:doc_metadata;type:json" json:"doc_metadata"`
	DocForm              string         `gorm:"column:doc_form;type:varchar(255);default:text_model" json:"doc_form"`
	DocLanguage          string         `gorm:"column:doc_language;type:varchar(255)" json:"doc_language"`
}

// TableName get sql table name.获取数据库表名
func (Document) TableName() string {
	return "documents"
}

func (doc *Document) DisplayStatus() string {
	status := ""
	if doc.IndexingStatus == "waiting" {
		status = "queuing"
	} else if !slices.Contains([]string{"completed", "error", "waiting"}, doc.IndexingStatus) && doc.IsPaused {
		status = "paused"
	} else if slices.Contains([]string{"parsing", "cleaning", "splitting", "indexing"}, doc.IndexingStatus) {
		status = "indexing"
	} else if doc.IndexingStatus == "error" {
		status = "error"
	} else if doc.IndexingStatus == "completed" && !doc.Archived && doc.Enabled {
		status = "available"
	} else if doc.IndexingStatus == "completed" && !doc.Archived && !doc.Enabled {
		status = "disabled"
	} else if doc.IndexingStatus == "completed" && doc.Archived {
		status = "archived"
	}
	return status

}
func (doc *Document) DataSourceInfoDict() map[string]any {
	if doc.DataSourceInfo != "" {

		var data_source_info_dict map[string]any
		err := json.Unmarshal([]byte(doc.DataSourceInfo), &data_source_info_dict)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) to dict failed:%v", doc.DataSourceInfo, err)
			data_source_info_dict = map[string]any{}
		}

		return data_source_info_dict
	}
	return map[string]any{}

}
func (doc *Document) DataSourceDetailDict() map[string]any {
	if doc.DataSourceInfo != "" {
		if doc.DataSourceType == "upload_file" {
			data_source_info_dict := doc.DataSourceInfoDict()
			upload_file_id := ""
			if _, ok := data_source_info_dict["upload_file_id"]; ok {
				if _, ok := data_source_info_dict["upload_file_id"].(string); ok {
					upload_file_id = data_source_info_dict["upload_file_id"].(string)
				}
			}
			file_detail := new(UploadFile)
			err := dbengine.Instance().DB.Model(&UploadFile{}).Where("id = ?", upload_file_id).First(file_detail).Error
			if err != nil {
				mlog.Errorf("get UploadFile failed:%v", err)
				file_detail = nil
			}
			if file_detail != nil {
				if file_detail.CreatedAt == nil {
					now := time.Now()
					file_detail.CreatedAt = &now
				}
				return map[string]any{
					"upload_file": map[string]any{
						"id":         file_detail.ID,
						"name":       file_detail.Name,
						"size":       file_detail.Size,
						"extension":  file_detail.Extension,
						"mime_type":  file_detail.MimeType,
						"created_by": file_detail.CreatedBy,
						"created_at": file_detail.CreatedAt.Unix(),
					},
				}
			}
		} else if slices.Contains([]string{"notion_import", "website_crawl"}, doc.DataSourceType) {
			return doc.DataSourceInfoDict()
		}
	}
	return nil

}
func (doc *Document) SegmentCount() int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&DocumentSegment{}).Where("document_id = ?", doc.ID).Count(&count).Error
	if err != nil {
		mlog.Errorf("count DocumentSegment failed:%v", err)
	}
	return count
}
func (doc *Document) AverageSegmentLength() int {
	if doc.WordCount != 0 && doc.SegmentCount() != 0 {
		return doc.WordCount // doc.SegmentCount
	}
	return 0

}
func (doc *Document) DatasetProcessRule() *DatasetProcessRule {
	if doc.DatasetProcessRuleID != "" {
		dataset_process_rule := new(DatasetProcessRule)
		err := dbengine.Instance().DB.Model(&DatasetProcessRule{}).Where("id = ?", doc.DatasetProcessRuleID).First(dataset_process_rule).Error
		if err != nil {
			mlog.Errorf("get DatasetProcessRule failed:%v", err)
			dataset_process_rule = nil
		}
		return dataset_process_rule
	}
	return nil

}
func (doc *Document) Dataset() *Dataset {
	if doc.DatasetID != "" {
		dataset := new(Dataset)
		err := dbengine.Instance().DB.Model(&Dataset{}).Where("id = ?", doc.DatasetID).First(dataset).Error
		if err != nil {
			mlog.Errorf("get Dataset failed:%v", err)
			dataset = nil
		}
		return dataset
	}
	return nil
}

func (doc *Document) HitCount() int {
	count := 0
	err := dbengine.Instance().DB.Model(&DocumentSegment{}).Select("sum(hit_count)").Where("document_id = ?", doc.ID).Scan(&count).Error
	if err != nil {
		mlog.Errorf("sum  DocumentSegment hit_count failed:%v", err)
	}
	return count
}
func (doc *Document) ProcessRuleDict() map[string]any {
	if doc.DatasetProcessRuleID != "" {
		return doc.DatasetProcessRule().ToDict()
	}
	return nil
}
func (doc *Document) ToDict() map[string]any {
	dataset_process_rule := doc.DatasetProcessRule()
	rsp := map[string]any{
		"id":                      doc.ID,
		"tenant_id":               doc.TenantID,
		"dataset_id":              doc.DatasetID,
		"position":                doc.Position,
		"data_source_type":        doc.DataSourceType,
		"data_source_info":        doc.DataSourceInfo,
		"dataset_process_rule_id": doc.DatasetProcessRuleID,
		"batch":                   doc.Batch,
		"name":                    doc.Name,
		"created_from":            doc.CreatedFrom,
		"created_by":              doc.CreatedBy,
		"created_api_request_id":  doc.CreatedAPIRequestID,
		"created_at":              doc.CreatedAt,
		"processing_started_at":   doc.ProcessingStartedAt,
		"file_id":                 doc.FileID,
		"word_count":              doc.WordCount,
		"parsing_completed_at":    doc.ParsingCompletedAt,
		"cleaning_completed_at":   doc.CleaningCompletedAt,
		"splitting_completed_at":  doc.SplittingCompletedAt,
		"tokens":                  doc.Tokens,
		"indexing_latency":        doc.IndexingLatency,
		"completed_at":            doc.CompletedAt,
		"is_paused":               doc.IsPaused,
		"paused_by":               doc.PausedBy,
		"paused_at":               doc.PausedAt,
		"error":                   doc.Error,
		"stopped_at":              doc.StoppedAt,
		"indexing_status":         doc.IndexingStatus,
		"enabled":                 doc.Enabled,
		"disabled_at":             doc.DisabledAt,
		"disabled_by":             doc.DisabledBy,
		"archived":                doc.Archived,
		"archived_reason":         doc.ArchivedReason,
		"archived_by":             doc.ArchivedBy,
		"archived_at":             doc.ArchivedAt,
		"updated_at":              doc.UpdatedAt,
		"doc_type":                doc.DocType,
		"doc_metadata":            doc.DocMetadata,
		"doc_form":                doc.DocForm,
		"doc_language":            doc.DocLanguage,
		"display_status":          doc.DisplayStatus(),
		"data_source_info_dict":   doc.DataSourceInfoDict(),
		"average_segment_length":  doc.AverageSegmentLength(),
		"dataset_process_rule":    "",
		"dataset":                 nil,
		".SegmentCount":           doc.SegmentCount(),
		"hit_count":               doc.HitCount(),
	}
	if dataset_process_rule != nil {
		rsp["dataset_process_rule"] = dataset_process_rule.ToDict()
	}
	dataset := doc.Dataset()
	if dataset != nil {
		rsp["dataset"] = dataset.ToDict()
	}

	return rsp

}
func NewDocumentFromDict(data map[string]any) *Document {
	bindata, _ := json.Marshal(data)
	doc := new(Document)
	err := json.Unmarshal(bindata, doc)
	if err != nil {
		mlog.Errorf("json unmarshal=%s to Document failed:%v", string(bindata), err)
		doc = nil
	}
	return doc
}

// DocumentSegment [...]
type DocumentSegment struct {
	ID            string         `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID      string         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	DatasetID     string         `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	DocumentID    string         `gorm:"column:document_id;type:varchar(36);not null" json:"document_id"`
	Position      int            `gorm:"column:position;type:int;not null" json:"position"`
	Content       string         `gorm:"column:content;type:text;not null" json:"content"`
	WordCount     int            `gorm:"column:word_count;type:int;not null" json:"word_count"`
	Tokens        int            `gorm:"column:tokens;type:int;not null" json:"tokens"`
	Keywords      datatypes.JSON `gorm:"column:keywords;type:json" json:"keywords"`
	IndexNodeID   string         `gorm:"column:index_node_id;type:varchar(255)" json:"index_node_id"`
	IndexNodeHash string         `gorm:"column:index_node_hash;type:varchar(255)" json:"index_node_hash"`
	HitCount      int            `gorm:"column:hit_count;type:int;not null" json:"hit_count"`
	Enabled       bool           `gorm:"column:enabled;type:tinyint(1);not null;default:1" json:"enabled"`
	DisabledAt    *time.Time     `gorm:"column:disabled_at;type:timestamp" json:"disabled_at"`
	DisabledBy    string         `gorm:"column:disabled_by;type:varchar(36)" json:"disabled_by"`
	Status        string         `gorm:"column:status;type:varchar(255);default:waiting" json:"status"`
	CreatedBy     string         `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt     *time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	IndexingAt    *time.Time     `gorm:"column:indexing_at;type:timestamp" json:"indexing_at"`
	CompletedAt   *time.Time     `gorm:"column:completed_at;type:timestamp" json:"completed_at"`
	Error         string         `gorm:"column:error;type:text" json:"error"`
	StoppedAt     *time.Time     `gorm:"column:stopped_at;type:timestamp" json:"stopped_at"`
	Answer        string         `gorm:"column:answer;type:text" json:"answer"`
	UpdatedBy     string         `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (DocumentSegment) TableName() string {
	return "document_segments"
}

// ChildChunk [...]
type ChildChunk struct {
	ID            string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID      string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	DatasetID     string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	DocumentID    string     `gorm:"column:document_id;type:varchar(36);not null" json:"document_id"`
	SegmentID     string     `gorm:"column:segment_id;type:varchar(36);not null" json:"segment_id"`
	Position      int        `gorm:"column:position;type:int;not null" json:"position"`
	Content       string     `gorm:"column:content;type:text;not null" json:"content"`
	WordCount     int        `gorm:"column:word_count;type:int;not null" json:"word_count"`
	IndexNodeID   string     `gorm:"column:index_node_id;type:varchar(255)" json:"index_node_id"`
	IndexNodeHash string     `gorm:"column:index_node_hash;type:varchar(255)" json:"index_node_hash"`
	Type          string     `gorm:"column:type;type:varchar(255);default:automatic" json:"type"`
	CreatedBy     string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy     string     `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UpdatedAt     *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IndexingAt    *time.Time `gorm:"column:indexing_at;type:timestamp" json:"indexing_at"`
	CompletedAt   *time.Time `gorm:"column:completed_at;type:timestamp" json:"completed_at"`
	Error         string     `gorm:"column:error;type:text" json:"error"`
	//     @property
	//     def dataset(self):
	//         return db.session.query(Dataset).filter(Dataset.ID == self.DatasetID).first()

	//     @property
	//     def document(self):
	//         return db.session.query(Document).filter(Document.ID == self.document_id).first()

	// @property
	// def segment(self):
	//
	//	return db.session.query(DocumentSegment).filter(DocumentSegment.ID == self.segment_id).first()
}

// TableName get sql table name.获取数据库表名
func (ChildChunk) TableName() string {
	return "child_chunks"
}

// AppDatasetJoin [...]
type AppDatasetJoin struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	AppID     string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App       *App       `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	DatasetID string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	Dataset   *Dataset   `json:"dataset" form:"dataset" gorm:"foreignKey:DatasetID;references:ID;"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// @property
	// def app(self):
	//
	//	return db.session.get(App, self.app_id)
}

// TableName get sql table name.获取数据库表名
func (AppDatasetJoin) TableName() string {
	return "app_dataset_joins"
}

// DatasetQuery [...]
type DatasetQuery struct {
	ID            string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	DatasetID     string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	Content       string     `gorm:"column:content;type:text;not null" json:"content"`
	Source        string     `gorm:"column:source;type:varchar(255);not null" json:"source"`
	SourceAppID   string     `gorm:"column:source_app_id;type:varchar(36)" json:"source_app_id"`
	CreatedByRole string     `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy     string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (DatasetQuery) TableName() string {
	return "dataset_queries"
}

// DatasetKeywordTable [...]
type DatasetKeywordTable struct {
	ID             string `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	DatasetID      string `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	KeywordTable   string `gorm:"column:keyword_table;type:text;not null" json:"keyword_table"`
	DataSourceType string `gorm:"column:data_source_type;type:varchar(255);default:database" json:"data_source_type"`
	//     @property
	//     def keyword_table_dict(self):
	//         class SetDecoder(json.JSONDecoder):
	//             def __init__(self, *args, **kwargs):
	//                 super().__init__(object_hook=self.object_hook, *args, **kwargs)

	//             def object_hook(self, dct):
	//                 if isinstance(dct, dict):
	//                     for keyword, node_idxs in dct.items():
	//                         if isinstance(node_idxs, list):
	//                             dct[keyword] = set(node_idxs)
	//                 return dct

	// # get dataset
	// dataset = Dataset.query.filter_by(id=self.DatasetID).first()
	// if not dataset:
	//
	//	return None
	//
	// if self.DataSourceType == "database":
	//
	//	return json.loads(self.keyword_table, cls=SetDecoder) if self.keyword_table else None
	//
	// else:
	//
	//	file_key = "keyword_files/" + dataset.tenant_id + "/" + self.DatasetID + ".txt"
	//	try:
	//	    keyword_table_text = storage.load_once(file_key)
	//	    if keyword_table_text:
	//	        return json.loads(keyword_table_text.decode("utf-8"), cls=SetDecoder)
	//	    return None
	//	except Exception as e:
	//	    logging.exception(f"Failed to load keyword table from file: {file_key}")
	//	    return None
}

// TableName get sql table name.获取数据库表名
func (DatasetKeywordTable) TableName() string {
	return "dataset_keyword_tables"
}

// Embedding [...]
type Embedding struct {
	ID           string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Hash         string     `gorm:"column:hash;type:varchar(64);not null" json:"hash"`
	Embedding    []byte     `gorm:"column:embedding;type:blob;not null" json:"embedding"`
	CreatedAt    *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ModelName    string     `gorm:"column:model_name;type:varchar(255);default:text-embedding-ada-002" json:"model_name"`
	ProviderName string     `gorm:"column:provider_name;type:varchar(255);default:''" json:"provider_name"`

	//     def set_embedding(self, embedding_data: list[float]):
	//         self.embedding = pickle.dumps(embedding_data, protocol=pickle.HIGHEST_PROTOCOL)

	// def get_embedding(self) -> list[float]:
	//
	//	return cast(list[float], pickle.loads(self.embedding))
}

// TableName get sql table name.获取数据库表名
func (Embedding) TableName() string {
	return "embeddings"
}

// DatasetCollectionBinding [...]
type DatasetCollectionBinding struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	ProviderName   string     `gorm:"column:provider_name;type:varchar(40);not null" json:"provider_name"`
	ModelName      string     `gorm:"column:model_name;type:varchar(255);not null" json:"model_name"`
	CollectionName string     `gorm:"column:collection_name;type:varchar(64);not null" json:"collection_name"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Type           string     `gorm:"column:type;type:varchar(40);default:dataset" json:"type"`
}

// TableName get sql table name.获取数据库表名
func (DatasetCollectionBinding) TableName() string {
	return "dataset_collection_bindings"
}

// TidbAuthBinding [...]
type TidbAuthBinding struct {
	ID          string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID    string     `gorm:"column:tenant_id;type:varchar(36)" json:"tenant_id"`
	ClusterID   string     `gorm:"column:cluster_id;type:varchar(255);not null" json:"cluster_id"`
	ClusterName string     `gorm:"column:cluster_name;type:varchar(255);not null" json:"cluster_name"`
	Active      bool       `gorm:"column:active;type:tinyint(1);not null;default:0" json:"active"`
	Status      string     `gorm:"column:status;type:varchar(255);default:CREATING" json:"status"`
	Account     string     `gorm:"column:account;type:varchar(255);not null" json:"account"`
	Password    string     `gorm:"column:password;type:varchar(255);not null" json:"password"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (TidbAuthBinding) TableName() string {
	return "tidb_auth_bindings"
}

// Whitelist [...]
type Whitelist struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID  string     `gorm:"column:tenant_id;type:varchar(36)" json:"tenant_id"`
	Category  string     `gorm:"column:category;type:varchar(255);not null" json:"category"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (Whitelist) TableName() string {
	return "whitelists"
}

// DatasetPermission [...]
type DatasetPermission struct {
	ID            string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	DatasetID     string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	AccountID     string     `gorm:"column:account_id;type:varchar(36);not null" json:"account_id"`
	HasPermission bool       `gorm:"column:has_permission;type:tinyint(1);not null;default:1" json:"has_permission"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	TenantID      string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
}

// TableName get sql table name.获取数据库表名
func (DatasetPermission) TableName() string {
	return "dataset_permissions"
}

// ExternalKnowledgeBinding [...]
type ExternalKnowledgeBinding struct {
	ID                     string                `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID               string                `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	ExternalKnowledgeApiID string                `gorm:"column:external_knowledge_api_id;type:varchar(36);not null" json:"external_knowledge_api_id"`
	ExternalKnowledgeApi   *ExternalKnowledgeApi `gorm:"ForeignKey:ExternalKnowledgeApiID;AssociationForeignKey:ID" json:"external_knowledge_api"`
	DatasetID              string                `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	Dataset                *Dataset              `gorm:"ForeignKey:DatasetID;AssociationForeignKey:ID" json:"dataset"`
	ExternalKnowledgeID    string                `gorm:"column:external_knowledge_id;type:text;not null" json:"external_knowledge_id"`
	CreatedBy              string                `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt              *time.Time            `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy              string                `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UpdatedAt              *time.Time            `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ExternalKnowledgeBinding) TableName() string {
	return "external_knowledge_bindings"
}

// ExternalKnowledgeApi [...]
type ExternalKnowledgeApi struct {
	ID          string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string     `gorm:"column:description;type:varchar(255);not null" json:"description"`
	TenantID    string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Settings    string     `gorm:"column:settings;type:text" json:"settings"`
	CreatedBy   string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy   string     `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	//     def to_dict(self):
	//         return {
	//             "id": self.ID,
	//             "tenant_id": self.tenant_id,
	//             "name": self.name,
	//             "description": self.description,
	//             "settings": self.settings_dict,
	//             "dataset_bindings": self.dataset_bindings,
	//             "created_by": self.created_by,
	//             "created_at": self.created_at.isoformat(),
	//         }

	//     @property
	//     def settings_dict(self):
	//         try:
	//             return json.loads(self.settings) if self.settings else None
	//         except JSONDecodeError:
	//             return None

	//     @property
	//     def dataset_bindings(self):
	//         external_knowledge_bindings = (
	//             db.session.query(ExternalKnowledgeBindings)
	//             .filter(ExternalKnowledgeBindings.external_knowledge_api_id == self.id)
	//             .all()
	//         )
	//         dataset_ids = [binding.DatasetID for binding in external_knowledge_bindings]
	//         datasets = db.session.query(Dataset).filter(Dataset.id.in_(dataset_ids)).all()
	//         dataset_bindings = []
	//         for dataset in datasets:
	//             dataset_bindings.append({"id": dataset.id, "name": dataset.name})

	// return dataset_bindings
}

// TableName get sql table name.获取数据库表名
func (ExternalKnowledgeApi) TableName() string {
	return "external_knowledge_apis"
}

// DatasetAutoDisableLog [...]
type DatasetAutoDisableLog struct {
	ID         string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID   string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	DatasetID  string     `gorm:"column:dataset_id;type:varchar(36);not null" json:"dataset_id"`
	DocumentID string     `gorm:"column:document_id;type:varchar(36);not null" json:"document_id"`
	Notified   bool       `gorm:"column:notified;type:tinyint(1);not null;default:0" json:"notified"`
	CreatedAt  *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (DatasetAutoDisableLog) TableName() string {
	return "dataset_auto_disable_logs"
}
