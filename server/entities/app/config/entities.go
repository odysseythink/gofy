package config

import (
	"mlib.com/gofy/server/core/file"
	agententities "mlib.com/gofy/server/entities/agent"
	coreentities "mlib.com/gofy/server/entities/core"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	appconfigenumtypes "mlib.com/gofy/server/enum_types/app_config"
	"mlib.com/gofy/server/models"
)

type ModelConfigWithCredentialsEntity struct {
	Provider            string
	Model               string
	ModelSchema         *modelruntimeentities.AIModelEntity
	Mode                modelruntimeentities.LLMMode
	ProviderModelBundle *coreentities.ProviderModelBundle
	Credentials         map[string]any
	Parameters          map[string]any
	Stop                []string
}

// ModelConfigEntity represents model config entity
type ModelConfigEntity struct {
	Provider   string         `json:"provider"`
	Model      string         `json:"model"`
	Mode       *string        `json:"mode,omitempty"`
	Parameters map[string]any `json:"parameters"`
	Stop       []string       `json:"stop"`
}

// AdvancedChatMessageEntity represents advanced chat message entity
type AdvancedChatMessageEntity struct {
	Text string                                 `json:"text"`
	Role modelruntimeentities.PromptMessageRole `json:"role"`
}

// AdvancedChatPromptTemplateEntity represents advanced chat prompt template entity
type AdvancedChatPromptTemplateEntity struct {
	Messages []AdvancedChatMessageEntity `json:"messages"`
}

// AdvancedCompletionPromptTemplateEntity represents advanced completion prompt template entity
type AdvancedCompletionPromptTemplateEntity struct {
	Prompt     string            `json:"prompt"`
	RolePrefix *RolePrefixEntity `json:"role_prefix,omitempty"`
}

// RolePrefixEntity represents role prefix entity
type RolePrefixEntity struct {
	User      string `json:"user"`
	Assistant string `json:"assistant"`
}

// PromptTemplateEntity represents prompt template entity
type PromptTemplateEntity struct {
	PromptType                       appconfigenumtypes.PromptType           `json:"prompt_type"`
	SimplePromptTemplate             *string                                 `json:"simple_prompt_template,omitempty"`
	AdvancedChatPromptTemplate       *AdvancedChatPromptTemplateEntity       `json:"advanced_chat_prompt_template,omitempty"`
	AdvancedCompletionPromptTemplate *AdvancedCompletionPromptTemplateEntity `json:"advanced_completion_prompt_template,omitempty"`
}

// VariableEntity represents variable entity
type VariableEntity struct {
	Variable                 string                                `json:"variable"`
	Label                    string                                `json:"label"`
	Description              string                                `json:"description"`
	Type                     appconfigenumtypes.VariableEntityType `json:"type"`
	Required                 bool                                  `json:"required"`
	MaxLength                int                                   `json:"max_length,omitempty"`
	Options                  []string                              `json:"options"`
	AllowedFileTypes         []file.FileType                       `json:"allowed_file_types"`
	AllowedFileExtensions    []string                              `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []file.FileTransferMethod             `json:"allowed_file_upload_methods"`
}

// ExternalDataVariableEntity represents external data variable entity
type ExternalDataVariableEntity struct {
	Variable string         `json:"variable"`
	Type     string         `json:"type"`
	Config   map[string]any `json:"config"`
}

// DatasetRetrieveConfigEntity represents dataset retrieve config entity
type DatasetRetrieveConfigEntity struct {
	QueryVariable    *string                             `json:"query_variable,omitempty"`
	RetrieveStrategy appconfigenumtypes.RetrieveStrategy `json:"retrieve_strategy"`
	TopK             int                                 `json:"top_k,omitempty"`
	ScoreThreshold   float64                             `json:"score_threshold"`
	RerankMode       string                              `json:"rerank_mode,omitempty"`
	RerankingModel   map[string]any                      `json:"reranking_model,omitempty"`
	Weights          map[string]any                      `json:"weights,omitempty"`
	RerankingEnabled bool                                `json:"reranking_enabled"`
}

func NewDatasetRetrieveConfigEntity() *DatasetRetrieveConfigEntity {
	return &DatasetRetrieveConfigEntity{
		RerankMode:       "reranking_model",
		RerankingEnabled: true,
	}
}

// DatasetEntity represents dataset entity
type DatasetEntity struct {
	DatasetIDs     []string                    `json:"dataset_ids"`
	RetrieveConfig DatasetRetrieveConfigEntity `json:"retrieve_config"`
}

// SensitiveWordAvoidanceEntity represents sensitive word avoidance entity
type SensitiveWordAvoidanceEntity struct {
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

// TextToSpeechEntity represents text to speech entity
type TextToSpeechEntity struct {
	Enabled  bool   `json:"enabled"`
	Voice    string `json:"voice,omitempty"`
	Language string `json:"language,omitempty"`
}

// TracingConfigEntity represents tracing config entity
type TracingConfigEntity struct {
	Enabled         bool   `json:"enabled"`
	TracingProvider string `json:"tracing_provider"`
}

// AppAdditionalFeatures represents app additional features
type AppAdditionalFeatures struct {
	FileUpload                    *file.FileUploadConfig `json:"file_upload,omitempty"`
	OpeningStatement              string                 `json:"opening_statement,omitempty"`
	SuggestedQuestions            []string               `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer bool                   `json:"suggested_questions_after_answer"`
	ShowRetrieveSource            bool                   `json:"show_retrieve_source"`
	MoreLikeThis                  bool                   `json:"more_like_this"`
	SpeechToText                  bool                   `json:"speech_to_text"`
	TextToSpeech                  *TextToSpeechEntity    `json:"text_to_speech,omitempty"`
	TraceConfig                   *TracingConfigEntity   `json:"trace_config,omitempty"`
}

type AppConfiger interface {
	GetTenantID() string
	SetTenantID(string)
	GetAppID() string
	SetAppID(string)
	GetAppMode() models.AppMode
	SetAppMode(models.AppMode)
	GetAdditionalFeatures() *AppAdditionalFeatures
	SetAdditionalFeatures(*AppAdditionalFeatures)
	GetVariables() []*VariableEntity
	SetVariables([]*VariableEntity)
	GetSensitiveWordAvoidance() *SensitiveWordAvoidanceEntity
	SetSensitiveWordAvoidance(*SensitiveWordAvoidanceEntity)
}

// AppConfig represents app config entity
type AppConfig struct {
	TenantID               string                        `json:"tenant_id"`
	AppID                  string                        `json:"app_id"`
	AppMode                models.AppMode                `json:"app_mode"`
	AdditionalFeatures     *AppAdditionalFeatures        `json:"additional_features"`
	Variables              []*VariableEntity             `json:"variables"`
	SensitiveWordAvoidance *SensitiveWordAvoidanceEntity `json:"sensitive_word_avoidance,omitempty"`
}

func (cfg *AppConfig) GetTenantID() string                           { return cfg.TenantID }
func (cfg *AppConfig) GetAppID() string                              { return cfg.AppID }
func (cfg *AppConfig) GetAppMode() models.AppMode                    { return cfg.AppMode }
func (cfg *AppConfig) GetAdditionalFeatures() *AppAdditionalFeatures { return cfg.AdditionalFeatures }
func (cfg *AppConfig) GetVariables() []*VariableEntity               { return cfg.Variables }
func (cfg *AppConfig) GetSensitiveWordAvoidance() *SensitiveWordAvoidanceEntity {
	return cfg.SensitiveWordAvoidance
}

func (cfg *AppConfig) SetTenantID(val string)                           { cfg.TenantID = val }
func (cfg *AppConfig) SetAppID(val string)                              { cfg.AppID = val }
func (cfg *AppConfig) SetAppMode(val models.AppMode)                    { cfg.AppMode = val }
func (cfg *AppConfig) SetAdditionalFeatures(val *AppAdditionalFeatures) { cfg.AdditionalFeatures = val }
func (cfg *AppConfig) SetVariables(val []*VariableEntity)               { cfg.Variables = val }
func (cfg *AppConfig) SetSensitiveWordAvoidance(val *SensitiveWordAvoidanceEntity) {
	cfg.SensitiveWordAvoidance = val
}

// EasyUIBasedAppConfig represents easy UI based app config entity
type EasyUIBasedAppConfig struct {
	*AppConfig
	AppModelConfigFrom    appconfigenumtypes.EasyUIBasedAppModelConfigFrom `json:"app_model_config_from"`
	AppModelConfigID      string                                           `json:"app_model_config_id"`
	AppModelConfigDict    map[string]any                                   `json:"app_model_config_dict"`
	Model                 ModelConfigEntity                                `json:"model"`
	PromptTemplate        PromptTemplateEntity                             `json:"prompt_template"`
	Dataset               *DatasetEntity                                   `json:"dataset,omitempty"`
	ExternalDataVariables []*ExternalDataVariableEntity                    `json:"external_data_variables"`
}

// WorkflowUIBasedAppConfig represents workflow UI based app config entity
type WorkflowUIBasedAppConfig struct {
	*AppConfig
	WorkflowID string `json:"workflow_id"`
}

type AdvancedChatAppConfig struct {
	*WorkflowUIBasedAppConfig
}
type AgentChatAppConfig struct {
	*EasyUIBasedAppConfig
	Agent *agententities.AgentEntity `json:"agent"`
}
