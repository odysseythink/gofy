package modelruntime

import (
	"fmt"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
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
)

func (t ModelType) Valid() bool {
	return t == Model_LLM ||
		t == Model_TEXT_EMBEDDING ||
		t == Model_RERANK ||
		t == Model_SPEECH2TEXT ||
		t == Model_MODERATION ||
		t == Model_TTS
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

type DefaultParameterNameType string

/*
Enum class for parameter template variable.
*/
const (
	DefaultParameterName_TEMPERATURE       DefaultParameterNameType = "temperature"
	DefaultParameterName_TOP_P             DefaultParameterNameType = "top_p"
	DefaultParameterName_TOP_K             DefaultParameterNameType = "top_k"
	DefaultParameterName_PRESENCE_PENALTY  DefaultParameterNameType = "presence_penalty"
	DefaultParameterName_FREQUENCY_PENALTY DefaultParameterNameType = "frequency_penalty"
	DefaultParameterName_MAX_TOKENS        DefaultParameterNameType = "max_tokens"
	DefaultParameterName_RESPONSE_FORMAT   DefaultParameterNameType = "response_format"
	DefaultParameterName_JSON_SCHEMA       DefaultParameterNameType = "json_schema"
)

type ParameterType string

/*
Enum class for parameter type.
*/
const (
	ParameterType_FLOAT   ParameterType = "float"
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
