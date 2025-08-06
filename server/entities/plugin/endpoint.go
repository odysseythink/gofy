package plugin

// from datetime import datetime
// from typing import Optional

// from pydantic import BaseModel, Field, model_validator

// from configs import dify_config
// from core.entities.provider_entities import ProviderConfig
// from core.plugin.entities.base import BasePluginEntity

import (
	"time"

	providerentities "mlib.com/gofy/server/entities/provider"
	parameterenumtypes "mlib.com/gofy/server/enum_types/parameter"
)

type EndpointDeclaration struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Hidden bool   `json:"hidden"`
}
type EndpointProviderDeclaration[T1 parameterenumtypes.AppSelectorScopeType | parameterenumtypes.ModelSelectorScopeType | parameterenumtypes.ToolSelectorScopeType, T2 int | string] struct {
	Settings  []*providerentities.ProviderConfig[T1, T2] `json:"settings"`
	Endpoints []*EndpointDeclaration                     `json:"endpoints"`
}
type EndpointEntity[T1 parameterenumtypes.AppSelectorScopeType | parameterenumtypes.ModelSelectorScopeType | parameterenumtypes.ToolSelectorScopeType, T2 int | string] struct {
	*BasePluginEntity
	Settings    map[string]any                      `json:"settings"`
	TenantID    string                              `json:"tenant_id"`
	PluginID    string                              `json:"plugin_id"`
	ExpiredAt   time.Time                           `json:"expired_at"`
	Declaration EndpointProviderDeclaration[T1, T2] `json:"declaration"`
}
type EndpointEntityWithInstance[T1 parameterenumtypes.AppSelectorScopeType | parameterenumtypes.ModelSelectorScopeType | parameterenumtypes.ToolSelectorScopeType, T2 int | string] struct {
	*EndpointEntity[T1, T2]
	Name    string `json:"path"`
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	HookID  string `json:"hook_id"`
}
