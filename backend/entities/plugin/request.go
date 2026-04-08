package plugin

// from typing import Any, Literal, Optional

// from pydantic import BaseModel, ConfigDict, Field, field_validator

// from core.entities.provider_entities import BasicProviderConfig
// from core.model_runtime.entities.message_entities import (
//     AssistantPromptMessage,
//     PromptMessage,
//     PromptMessageRole,
//     PromptMessageTool,
//     SystemPromptMessage,
//     ToolPromptMessage,
//     UserPromptMessage,
// )
// from core.model_runtime.entities.model_entities import ModelType
// from core.workflow.nodes.parameter_extractor.entities import (
//     ModelConfig as ParameterExtractorModelConfig,
// )
// from core.workflow.nodes.parameter_extractor.entities import (
//     ParameterConfig,
// )
// from core.workflow.nodes.question_classifier.entities import (
//     ClassConfig,
// )
// from core.workflow.nodes.question_classifier.entities import (
//     ModelConfig as QuestionClassifierModelConfig,
// )

import (
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	llmnodesentities "mlib.com/gofy/server/entities/nodes/llm"
	parameterextractornodesentities "mlib.com/gofy/server/entities/nodes/parameter_extractor"
	questionclassifiernodesentities "mlib.com/gofy/server/entities/nodes/question_classifier"
	providerentities "mlib.com/gofy/server/entities/provider"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
)

type InvokeCredentials struct {
	ToolCredentials map[string]string `json:"tool_credentials"` //description="Map of tool provider to credential id, used to store the credential id for the tool provider.",
}
type PluginInvokeContext struct {
	Credentials *InvokeCredentials `json:"credentials"` //description="Credentials context for the plugin invocation or backward invocation.",
}
type RequestInvokeTool struct {
	ToolType       string         `json:"tool_type"` // Literal["builtin", "workflow", "api", "mcp"]
	Provider       string         `json:"provider"`
	Tool           string         `json:"tool"`
	ToolParameters map[string]any `json:"tool_parameters"`
	CredentialID   *string        `json:"credential_id"`
}
type RequestInvokeModeler interface {
	ModelType() modelruntimeenumtypes.ModelType
}
type BaseRequestInvokeModel struct {
	Provider    string         `json:"tool_credentials"`
	Model       string         `json:"model"`
	ModelConfig map[string]any `json:"model_config"`
}

type RequestInvokeLLM struct {
	*BaseRequestInvokeModel
	Mode             string                                    `json:"mode"`
	CompletionParams map[string]any                            `json:"completion_params"`
	PromptMessages   []modelruntimeentities.PromptMessager     `json:"prompt_messages"`
	Tools            []*modelruntimeentities.PromptMessageTool `json:"tools"`
	Stop             []string                                  `json:"stop"`
	Stream           bool                                      `json:"stream"`

	ModelConfig map[string]any `json:"model_config"`
}

func (model *RequestInvokeLLM) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_LLM
}

type RequestInvokeLLMWithStructuredOutput struct {
	*RequestInvokeLLM
	StructuredOutputSchema map[string]any `json:"structured_output_schema"` //description="The schema of the structured output in JSON schema format"
}
type RequestInvokeTextEmbedding struct {
	*BaseRequestInvokeModel
	Texts []string `json:"texts"`
}

func (model *RequestInvokeTextEmbedding) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_TEXT_EMBEDDING
}

type RequestInvokeRerank struct {
	*BaseRequestInvokeModel
	Query          string   `json:"query"`
	Docs           []string `json:"docs"`
	ScoreThreshold float64  `json:"score_threshold"`
	TopN           int      `json:"top_n"`
}

func (model *RequestInvokeRerank) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_RERANK
}

type RequestInvokeTTS struct {
	*BaseRequestInvokeModel
	ContentText string `json:"content_text"`
	Voice       string `json:"voice"`
}

func (model *RequestInvokeTTS) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_TTS
}

type RequestInvokeSpeech2Text struct {
	*BaseRequestInvokeModel
	File []byte `json:"file"`
}

func (model *RequestInvokeSpeech2Text) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_SPEECH2TEXT
}

type RequestInvokeModeration struct {
	*BaseRequestInvokeModel
	Text string `json:"text"`
}

func (model *RequestInvokeModeration) ModelType() modelruntimeenumtypes.ModelType {
	return modelruntimeenumtypes.Model_MODERATION
}

type RequestInvokeParameterExtractorNode struct {
	Parameters  []*parameterextractornodesentities.ParameterConfig `json:"parameters"`
	Model       *llmnodesentities.ModelConfig                      `json:"model"`
	Instruction string                                             `json:"instruction"`
	Query       string                                             `json:"query"`
}
type RequestInvokeQuestionClassifierNode struct {
	Query       string                                         `json:"query"`
	Model       *llmnodesentities.ModelConfig                  `json:"model"`
	Classes     []*questionclassifiernodesentities.ClassConfig `json:"classes"`
	Instruction string                                         `json:"instruction"`
}
type RequestInvokeApp struct {
	AppID          string           `json:"app_id"`
	Inputs         map[string]any   `json:"inputs"`
	Query          *string          `json:"query"`
	ResponseMode   string           `json:"response_mode"` // Literal["blocking", "streaming"]
	ConversationID *string          `json:"conversation_id"`
	User           *string          `json:"user"`
	Files          []map[string]any `json:"files"`
}
type RequestInvokeEncrypt struct {
	Opt       string                                  `json:"opt"`       // Literal["encrypt", "decrypt", "clear"]
	Namespace string                                  `json:"namespace"` //Literal["endpoint"]
	Identity  string                                  `json:"identity"`
	Data      map[string]any                          `json:"data"`
	Config    []*providerentities.BasicProviderConfig `json:"config"`
}
type RequestInvokeSummary struct {
	Text        string `json:"text"`
	Instruction string `json:"instruction"`
}
type RequestRequestUploadFile struct {
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
}
type RequestFetchAppInfo struct {
	AppID string `json:"app_id"`
}
