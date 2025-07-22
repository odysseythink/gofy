package ops

import (
	"encoding/json"
	"time"
)

// BaseModel is the base model for all trace info models
type BaseModel[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	MessageID   string         `json:"message_id"`
	MessageData any            `json:"message_data"`
	Inputs      T1             `json:"inputs"`
	Outputs     T2             `json:"outputs"`
	StartTime   *time.Time     `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	Metadata    map[string]any `json:"metadata"`
}

// WorkflowTraceInfo represents workflow trace information
type WorkflowTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	WorkflowData           any            `json:"workflow_data"`
	TenantID               string         `json:"tenant_id"`
	WorkflowID             string         `json:"workflow_id"`
	WorkflowRunID          string         `json:"workflow_run_id"`
	WorkflowRunElapsedTime float64        `json:"workflow_run_elapsed_time"`
	WorkflowRunStatus      string         `json:"workflow_run_status"`
	WorkflowRunInputs      map[string]any `json:"workflow_run_inputs"`
	WorkflowRunOutputs     map[string]any `json:"workflow_run_outputs"`
	WorkflowVersionRun     string         `json:"workflow_run_version"`
	ConversationID         string         `json:"conversation_id"`
	WorkflowAppLogID       string         `json:"workflow_app_log_id"`
	Error                  string         `json:"error"`
	TotalTokens            int            `json:"total_tokens"`
	FileList               []string       `json:"file_list"`
	Query                  string         `json:"query"`
}

// MessageTraceInfo represents message trace information
type MessageTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	ConversationModel string `json:"conversation_model"`
	MessageTokens     int    `json:"message_tokens"`
	AnswerTokens      int    `json:"answer_tokens"`
	TotalTokens       int    `json:"total_tokens"`
	Error             string `json:"error"`
	FileList          any    `json:"file_list"`
	MessageFileData   any    `json:"message_file_data"`
	ConversationMode  string `json:"conversation_mode"`
}

// ModerationTraceInfo represents moderation trace information
type ModerationTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	Flagged        bool   `json:"flagged"`
	Action         string `json:"action"`
	PresetResponse string `json:"preset_response"`
	Query          string `json:"query"`
}

// SuggestedQuestionTraceInfo represents suggested question trace information
type SuggestedQuestionTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	TotalTokens       int      `json:"total_tokens"`
	Status            string   `json:"status"`
	Error             string   `json:"error"`
	FromAccountID     string   `json:"from_account_id"`
	AgentBased        *bool    `json:"agent_based"`
	FromSource        string   `json:"from_source"`
	ModelProvider     string   `json:"model_provider"`
	ModelID           string   `json:"model_id"`
	SuggestedQuestion []string `json:"suggested_question"`
	Level             string   `json:"level"`
	StatusMessage     string   `json:"status_message"`
	WorkflowRunID     string   `json:"workflow_run_id"`
}

// DatasetRetrievalTraceInfo represents dataset retrieval trace information
type DatasetRetrievalTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	Documents any `json:"documents"`
}

// ToolTraceInfo represents tool trace information
type ToolTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	ToolName       string          `json:"tool_name"`
	ToolInputs     map[string]any  `json:"tool_inputs"`
	ToolOutputs    string          `json:"tool_outputs"`
	Error          string          `json:"error"`
	ToolConfig     map[string]any  `json:"tool_config"`
	TimeCost       float64         `json:"time_cost"`
	ToolParameters map[string]any  `json:"tool_parameters"`
	FileURL        json.RawMessage `json:"file_url"`
}

// GenerateNameTraceInfo represents generate name trace information
type GenerateNameTraceInfo[T1 string | map[string]any | []any, T2 string | map[string]any | []any] struct {
	*BaseModel[T1, T2]
	TenantID string `json:"tenant_id"`
}

// TaskData represents task data
type TaskData struct {
	AppID         string `json:"app_id"`
	TraceInfoType string `json:"trace_info_type"`
	TraceInfo     any    `json:"trace_info"`
}

// var (
// 	TRACE_INFO_INFO_MAP = map[string]any{
// 		"WorkflowTraceInfo":          WorkflowTraceInfo,
// 		"MessageTraceInfo":           MessageTraceInfo,
// 		"ModerationTraceInfo":        ModerationTraceInfo,
// 		"SuggestedQuestionTraceInfo": SuggestedQuestionTraceInfo,
// 		"DatasetRetrievalTraceInfo":  DatasetRetrievalTraceInfo,
// 		"ToolTraceInfo":              ToolTraceInfo,
// 		"GenerateNameTraceInfo":      GenerateNameTraceInfo,
// 	}
// )

type TraceTaskName string

const (
	TraceTaskName_CONVERSATION_TRACE       TraceTaskName = "conversation"
	TraceTaskName_WORKFLOW_TRACE           TraceTaskName = "workflow"
	TraceTaskName_MESSAGE_TRACE            TraceTaskName = "message"
	TraceTaskName_MODERATION_TRACE         TraceTaskName = "moderation"
	TraceTaskName_SUGGESTED_QUESTION_TRACE TraceTaskName = "suggested_question"
	TraceTaskName_DATASET_RETRIEVAL_TRACE  TraceTaskName = "dataset_retrieval"
	TraceTaskName_TOOL_TRACE               TraceTaskName = "tool"
	TraceTaskName_GENERATE_NAME_TRACE      TraceTaskName = "generate_conversation_name"
)
