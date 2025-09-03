package generator

import (
	"fmt"

	"mlib.com/gofy/server/core/file"
	appconfigentities "mlib.com/gofy/server/entities/app/config"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
)

type SingleIterationRunEntity struct {
	NodeID string
	Inputs map[string]any
}

type AppGenerateEntityConfig struct {
	ArbitraryTypesAllowed bool `json:"arbitrary_types_allowed"`
}

func NewAppGenerateEntityConfig() *AppGenerateEntityConfig {
	return &AppGenerateEntityConfig{
		ArbitraryTypesAllowed: true,
	}
}

type AppGenerateEntitier interface {
	GetTaskID() string
	GetAppConfig() appconfigentities.AppConfiger
	GetFileUploadConfig() *file.FileUploadConfig
	GetInputs() map[string]any
	GetFiles() []*file.File
	GetUserID() string
	GetStream() bool
	GetInvokeFrom() appenumtypes.InvokeFrom
	GetCallDepth() int
	GetExtras() map[string]any
	SetTaskID(string)
	SetAppConfig(appconfigentities.AppConfiger)
	SetFileUploadConfig(*file.FileUploadConfig)
	SetInputs(map[string]any)
	SetFiles([]*file.File)
	SetUserID(string)
	SetStream(bool)
	SetInvokeFrom(appenumtypes.InvokeFrom)
	SetCallDepth(int)
	SetExtras(map[string]any)
}

type AppGenerateEntity[T interface {
	*appconfigentities.WorkflowUIBasedAppConfig | *appconfigentities.EasyUIBasedAppConfig | *appconfigentities.AdvancedChatAppConfig | *appconfigentities.AgentChatAppConfig
}] struct {
	TaskID           string                  `json:"task_id"`
	AppConfig        T                       `json:"app_config"`
	FileUploadConfig *file.FileUploadConfig  `json:"file_upload_config"`
	Inputs           map[string]any          `json:"inputs"`
	Files            []*file.File            `json:"files"`
	UserID           string                  `json:"user_id"`
	Stream           bool                    `json:"stream"`
	InvokeFrom       appenumtypes.InvokeFrom `json:"invoke_from"`
	CallDepth        int                     `json:"call_depth"`
	Extras           map[string]any          `json:"extras"`
	// TraceManager     *TraceQueueManager
}

type EasyUIBasedAppGenerateEntity struct {
	*AppGenerateEntity[*appconfigentities.EasyUIBasedAppConfig]
	ModelConf   *appconfigentities.ModelConfigWithCredentialsEntity `json:"model_conf"`
	Query       string                                              `json:"query"`
	ModelConfig map[string]any                                      `json:"model_config"`
}

type ConversationAppGenerateEntity[T interface {
	*appconfigentities.WorkflowUIBasedAppConfig | *appconfigentities.EasyUIBasedAppConfig | *appconfigentities.AdvancedChatAppConfig
}] struct {
	*AppGenerateEntity[T]
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

func (c *ConversationAppGenerateEntity[T]) ValidateParentMessageID() error {
	if c.InvokeFrom == appenumtypes.InvokeFrom_SERVICE_API && c.ParentMessageID != "" {
		return fmt.Errorf("parent_message_id should be UUID_NIL for service API")
	}
	return nil
}

type ChatAppGenerateEntity struct {
	*EasyUIBasedAppGenerateEntity
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

type CompletionAppGenerateEntity struct {
	*EasyUIBasedAppGenerateEntity
}

type AgentChatAppGenerateEntity struct {
	*AppGenerateEntity[*appconfigentities.AgentChatAppConfig]
	ModelConf       *appconfigentities.ModelConfigWithCredentialsEntity `json:"model_conf"`
	Query           string                                              `json:"query"`
	ModelConfig     map[string]any                                      `json:"model_config"`
	ConversationID  string                                              `json:"conversation_id"`
	ParentMessageID string                                              `json:"parent_message_id"`
}

type AdvancedChatAppGenerateEntity struct {
	*ConversationAppGenerateEntity[*appconfigentities.AdvancedChatAppConfig]
	WorkflowRunID      string                    `json:"workflow_run_id"`
	Query              string                    `json:"query"`
	SingleIterationRun *SingleIterationRunEntity `json:"single_iteration_run"`
}

type WorkflowAppGenerateEntity struct {
	*AppGenerateEntity[*appconfigentities.WorkflowUIBasedAppConfig]
	WorkflowRunID      string                    `json:"workflow_run_id"`
	SingleIterationRun *SingleIterationRunEntity `json:"single_iteration_run"`
}
