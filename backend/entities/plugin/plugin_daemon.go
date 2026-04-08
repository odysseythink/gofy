package plugin

// from collections.abc import Mapping, Sequence
// from datetime import datetime
// from enum import StrEnum
// from typing import Any, Generic, Optional, TypeVar

// from pydantic import BaseModel, ConfigDict, Field

// from core.agent.plugin_entities import AgentProviderEntityWithPlugin
// from core.model_runtime.entities.model_entities import AIModelEntity
// from core.model_runtime.entities.provider_entities import ProviderEntity
// from core.plugin.entities.base import BasePluginEntity
// from core.plugin.entities.parameters import PluginParameterOption
// from core.plugin.entities.plugin import PluginDeclaration, PluginEntity
// from core.tools.entities.common_entities import I18nObject
// from core.tools.entities.tool_entities import ToolProviderEntityWithPlugin

// T = TypeVar("T", bound=(BaseModel | dict | list | bool | str))
import (
	"time"

	agententities "mlib.com/gofy/server/entities/agent"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	parametertentities "mlib.com/gofy/server/entities/plugin/parameter"
	toolsentities "mlib.com/gofy/server/entities/tools"
	pluginenumtypes "mlib.com/gofy/server/enum_types/plugin"
	commontypes "mlib.com/gofy/server/types/common"
)

type PluginDaemonBasicResponse[T map[string]any | []any | bool | string] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
type InstallPluginMessage struct {
	Event pluginenumtypes.InstallPluginMessageEventType `json:"event"`
	Data  string                                        `json:"data"`
}
type PluginToolProviderEntity struct {
	Provider               string                                     `json:"provider"`
	PluginUniqueIdentifier string                                     `json:"plugin_unique_identifier"`
	PluginID               string                                     `json:"plugin_id"`
	Declaration            toolsentities.ToolProviderEntityWithPlugin `json:"declaration"`
}
type PluginAgentProviderEntity[T1 float64 | int | string, T2 float64 | int] struct {
	Provider               string                                              `json:"provider"`
	PluginUniqueIdentifier string                                              `json:"plugin_unique_identifier"`
	PluginID               string                                              `json:"plugin_id"`
	Declaration            agententities.AgentProviderEntityWithPlugin[T1, T2] `json:"declaration"`
	Meta                   struct {
		MinimumDifyVersion string `json:"minimum_dify_version"` // pattern=r"^\d{1,4}(\.\d{1,4}){1,3}(-\w{1,16})?$")
		Version            string `json:"version"`
	} `json:"meta"`
}
type PluginBasicBooleanResponse struct {
	Result      bool           `json:"result"`
	Credentials map[string]any `json:"credentials"`
}
type PluginModelSchemaEntity struct {
	ModelSchema modelruntimeentities.AIModelEntity `json:"model_schema"` //description="The model schema.")

	// pydantic configs
	ModelConfig map[string]any `json:"model_config"`
}
type PluginModelProviderEntity struct {
	ID                     string                              `json:"id"`                       //description="ID")
	CreatedAt              time.Time                           `json:"created_at"`               //"The created at time of the model provider.")
	UpdatedAt              time.Time                           `json:"updated_at"`               //"The updated at time of the model provider.")
	Provider               string                              `json:"provider"`                 //description="The provider of the model.")
	TenantID               string                              `json:"tenant_id"`                //description="The tenant ID.")
	PluginUniqueIdentifier string                              `json:"plugin_unique_identifier"` //description="The plugin unique identifier.")
	PluginID               string                              `json:"plugin_id"`                //description="The plugin ID.")
	Declaration            modelruntimeentities.ProviderEntity `json:"declaration"`              //description="The declaration of the model provider.")
}
type PluginTextEmbeddingNumTokensResponse struct {
	NumTokens []int `json:"num_tokens"` //description="The number of tokens.")
}
type PluginLLMNumTokensResponse struct {
	NumTokens []int `json:"num_tokens"` //description="The number of tokens.")
}
type PluginStringResultResponse struct {
	Result string `json:"result"` //description="The result of the string.")
}
type PluginVoiceEntity struct {
	Name  string `json:"name"`  //description="The name of the voice.")
	Value string `json:"value"` //description="The value of the voice.")
}
type PluginVoicesResponse struct {
	Voices []*PluginVoiceEntity `json:"voices"` //description="The result of the voices.")
}
type PluginDaemonError struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}

type PluginInstallTaskPluginStatus struct {
	PluginUniqueIdentifier string                                      `json:"plugin_unique_identifier"` //description="The plugin unique identifier of the install task.")
	PluginID               string                                      `json:"plugin_id"`                //description="The plugin ID of the install task.")
	Status                 pluginenumtypes.PluginInstallTaskStatusType `json:"status"`                   //description="The status of the install task.")
	Message                string                                      `json:"message"`                  //description="The message of the install task.")
	Icon                   string                                      `json:"icon"`                     //description="The icon of the plugin.")
	Labels                 commontypes.I18nObject                      `json:"labels"`                   //description="The labels of the plugin.")
}
type PluginInstallTask struct {
	*BasePluginEntity
	Status           pluginenumtypes.PluginInstallTaskStatusType `json:"status"`            //description="The status of the install task.")
	TotalPlugins     int                                         `json:"total_plugins"`     //description="The total number of plugins to be installed.")
	CompletedPlugins int                                         `json:"completed_plugins"` //description="The number of plugins that have been installed.")
	Plugins          []*PluginInstallTaskPluginStatus            `json:"plugins"`           //description="The status of the plugins.")
}
type PluginInstallTaskStartResponse struct {
	AllInstalled bool   `json:"all_installed"` //description="Whether all plugins are installed.")
	TaskID       string `json:"task_id"`       //description="The ID of the install task.")
}
type PluginVerification struct {
	AuthorizedCategory pluginenumtypes.AuthorizedCategoryType `json:"authorized_category"` //description="The authorized category of the plugin.")
}
type PluginDecodeResponse struct {
	UniqueIdentifier string              `json:"unique_identifier"` //description="The unique identifier of the plugin.")
	Manifest         PluginDeclaration   `json:"manifest"`
	Verification     *PluginVerification `json:"verification"` //default=None, description="Basic verification information")
}
type PluginOAuthAuthorizationUrlResponse struct {
	AuthorizationURL string `json:"authorization_url"` //description="The URL of the authorization.")
}
type PluginOAuthCredentialsResponse struct {
	Metadata    map[string]any `json:"metadata"`    //description="The metadata of the OAuth, like avatar url, name, etc."
	ExpiresAt   int            `json:"expires_at"`  //default=-1, description="The expires at time of the credentials. UTC timestamp.")
	Credentials map[string]any `json:"credentials"` //description="The credentials of the OAuth.")
}
type PluginListResponse struct {
	List  []*PluginEntity `json:"list"`
	Total int             `json:"total"`
}
type PluginDynamicSelectOptionsResponse struct {
	Options []*parametertentities.PluginParameterOption `json:"options"` //description="The options of the dynamic select.")
}
