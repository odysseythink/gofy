package modelruntime

import (
	"fmt"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
	commontypes "mlib.com/gofy/server/types/common"
)

type ModelType string

/*
Enum class for model type.
*/
const (
	Model_LLM            ModelType = "llm"
	Model_TEXT_EMBEDDING ModelType = "text-embedding"
	Model_RERANK         ModelType = "rerank"
	Model_SPEECH2TEXT    ModelType = "speech2text"
	Model_MODERATION     ModelType = "moderation"
	Model_TTS            ModelType = "tts"
	Model_TEXT2IMG       ModelType = "text2img"
)

func (t ModelType) Valid() bool {
	return t == Model_LLM ||
		t == Model_TEXT_EMBEDDING ||
		t == Model_RERANK ||
		t == Model_SPEECH2TEXT ||
		t == Model_MODERATION ||
		t == Model_TTS ||
		t == Model_TEXT2IMG
}

func Str2ModelType(origin_model_type string) ModelType {
	/*
		Get model type from origin model type.

		:return: model type
	*/
	if slices.Contains([]string{"text-generation", string(Model_LLM)}, origin_model_type) {
		return Model_LLM
	} else if slices.Contains([]string{"embeddings", string(Model_TEXT_EMBEDDING)}, origin_model_type) {
		return Model_TEXT_EMBEDDING
	} else if slices.Contains([]string{"reranking", string(Model_RERANK)}, origin_model_type) {
		return Model_RERANK
	} else if slices.Contains([]string{"speech2text", string(Model_SPEECH2TEXT)}, origin_model_type) {
		return Model_SPEECH2TEXT
	} else if slices.Contains([]string{"tts", string(Model_TTS)}, origin_model_type) {
		return Model_TTS
	} else if slices.Contains([]string{"text2img", string(Model_TEXT2IMG)}, origin_model_type) {
		return Model_TEXT2IMG
	} else if origin_model_type == string(Model_MODERATION) {
		return Model_MODERATION
	} else {
		panic(exceptions.NewValueError(fmt.Sprintf("invalid origin model type %s", origin_model_type)))
	}
}
func (t ModelType) ToOriginModelType() string {
	/*
		Get origin model type from model type.

		:return: origin model type
	*/
	if t == Model_LLM {
		return "text-generation"
	} else if t == Model_TEXT_EMBEDDING {
		return "embeddings"
	} else if t == Model_RERANK {
		return "reranking"
	} else if t == Model_SPEECH2TEXT {
		return "speech2text"
	} else if t == Model_TTS {
		return "tts"
	} else if t == Model_MODERATION {
		return "moderation"
	} else if t == Model_TEXT2IMG {
		return "text2img"
	} else {
		panic(exceptions.NewValueError(fmt.Sprintf("invalid model type %v", t)))
	}
}

type FetchFrom string

/*
Enum class for fetch from.
*/
const (
	FetchFrom_PREDEFINED_MODEL   FetchFrom = "predefined-model"
	FetchFrom_CUSTOMIZABLE_MODEL FetchFrom = "customizable-model"
)

type ModelFeature string

/*
Enum class for llm feature.
*/
const (
	ModelFeature_TOOL_CALL        ModelFeature = "tool-call"
	ModelFeature_MULTI_TOOL_CALL  ModelFeature = "multi-tool-call"
	ModelFeature_AGENT_THOUGHT    ModelFeature = "agent-thought"
	ModelFeature_VISION           ModelFeature = "vision"
	ModelFeature_STREAM_TOOL_CALL ModelFeature = "stream-tool-call"
	ModelFeature_DOCUMENT         ModelFeature = "document"
	ModelFeature_VIDEO            ModelFeature = "video"
	ModelFeature_AUDIO            ModelFeature = "audio"
)

type DefaultParameterName string

/*
Enum class for parameter template variable.
*/
const (
	DefaultParameterName_TEMPERATURE       DefaultParameterName = "temperature"
	DefaultParameterName_TOP_P             DefaultParameterName = "top_p"
	DefaultParameterName_TOP_K             DefaultParameterName = "top_k"
	DefaultParameterName_PRESENCE_PENALTY  DefaultParameterName = "presence_penalty"
	DefaultParameterName_FREQUENCY_PENALTY DefaultParameterName = "frequency_penalty"
	DefaultParameterName_MAX_TOKENS        DefaultParameterName = "max_tokens"
	DefaultParameterName_RESPONSE_FORMAT   DefaultParameterName = "response_format"
	DefaultParameterName_JSON_SCHEMA       DefaultParameterName = "json_schema"
)

type ParameterType string

/*
Enum class for parameter type.
*/
const (
	ParameterType_FLOAT   ParameterType = " float64"
	ParameterType_INT     ParameterType = "int"
	ParameterType_STRING  ParameterType = "string"
	ParameterType_BOOLEAN ParameterType = "boolean"
	ParameterType_TEXT    ParameterType = "text"
)

type ModelPropertyKey string

/*
Enum class for model property key.
*/
const (
	ModelPropertyKey_MODE                      ModelPropertyKey = "mode"
	ModelPropertyKey_CONTEXT_SIZE              ModelPropertyKey = "context_size"
	ModelPropertyKey_MAX_CHUNKS                ModelPropertyKey = "max_chunks"
	ModelPropertyKey_FILE_UPLOAD_LIMIT         ModelPropertyKey = "file_upload_limit"
	ModelPropertyKey_SUPPORTED_FILE_EXTENSIONS ModelPropertyKey = "supported_file_extensions"
	ModelPropertyKey_MAX_CHARACTERS_PER_CHUNK  ModelPropertyKey = "max_characters_per_chunk"
	ModelPropertyKey_DEFAULT_VOICE             ModelPropertyKey = "default_voice"
	ModelPropertyKey_VOICES                    ModelPropertyKey = "voices"
	ModelPropertyKey_WORD_LIMIT                ModelPropertyKey = "word_limit"
	ModelPropertyKey_AUDIO_TYPE                ModelPropertyKey = "audio_type"
	ModelPropertyKey_MAX_WORKERS               ModelPropertyKey = "max_workers"
)

// ProviderModel represents a provider model.
type ProviderModel struct {
	Model           string                   `json:"model" yaml:"model"`
	Label           commontypes.I18nObject   `json:"label" yaml:"label"`
	ModelType       ModelType                `json:"model_type" yaml:"model_type"`
	Features        []ModelFeature           `json:"features" yaml:"features"`
	FetchFrom       FetchFrom                `json:"fetch_from" yaml:"fetch_from"`
	ModelProperties map[ModelPropertyKey]any `json:"model_properties" yaml:"model_properties"`
	Deprecated      bool                     `json:"deprecated" yaml:"deprecated"`
	ModelConfig     map[string]any           `json:"model_config" yaml:"model_config"`
}

// ParameterRule represents a parameter rule.
type ParameterRule struct {
	Name        string                 `json:"name" yaml:"name"`
	UseTemplate string                 `json:"use_template" yaml:"use_template"`
	Label       commontypes.I18nObject `json:"label" yaml:"label"`
	Type        ParameterType          `json:"type" yaml:"type"`
	Help        commontypes.I18nObject `json:"help" yaml:"help"`
	Required    bool                   `json:"required" yaml:"required"`
	Default     any                    `json:"default" yaml:"default"`
	Min         float64                `json:"min" yaml:"min"`
	Max         float64                `json:"max" yaml:"max"`
	Precision   int                    `json:"precision" yaml:"precision"`
	Options     []string               `json:"options" yaml:"options"`
}

// PriceConfig represents the pricing configuration.
type PriceConfig struct {
	Input    float64 `json:"input" yaml:"input"`
	Output   float64 `json:"output" yaml:"output"`
	Unit     float64 `json:"unit" yaml:"parametuniter_rules"`
	Currency string  `json:"currency" yaml:"currency"`
}

// AIModelEntity represents an AI model entity.
type AIModelEntity struct {
	Model           string                   `json:"model" yaml:"model"`
	Label           commontypes.I18nObject   `json:"label" yaml:"label"`
	ModelType       ModelType                `json:"model_type" yaml:"model_type"`
	Features        []ModelFeature           `json:"features" yaml:"features"`
	FetchFrom       FetchFrom                `json:"fetch_from" yaml:"fetch_from"`
	ModelProperties map[ModelPropertyKey]any `json:"model_properties" yaml:"model_properties"`
	Deprecated      bool                     `json:"deprecated" yaml:"deprecated"`
	ModelConfig     map[string]any           `json:"model_config" yaml:"model_config"`
	ParameterRules  []*ParameterRule         `json:"parameter_rules" yaml:"parameter_rules"`
	Pricing         *PriceConfig             `json:"pricing" yaml:"pricing"`
}

// ModelUsage represents model usage information.
type ModelUsage struct {
}

type PriceType string

/*
Enum class for price type.
*/
const (
	PriceType_INPUT  PriceType = "input"
	PriceType_OUTPUT PriceType = "output"
)

// PriceInfo represents price information.
type PriceInfo struct {
	UnitPrice   float64
	Unit        float64
	TotalAmount float64
	Currency    string
}
