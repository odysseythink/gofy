package plugin

// from datetime import datetime
// from typing import Optional

// from pydantic import BaseModel, Field, model_validator

// from configs import gofy_config
// from core.entities.provider_entities import ProviderConfig
// from core.plugin.entities.base import BasePluginEntity

import (
	"time"

	providerentities "github.com/odysseythink/gofy/backend/entities/provider"
)

type EndpointDeclaration struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Hidden bool   `json:"hidden"`
}
type EndpointProviderDeclaration struct {
	Settings  []*providerentities.ProviderConfig `json:"settings"`
	Endpoints []*EndpointDeclaration             `json:"endpoints"`
}
type EndpointEntity struct {
	*BasePluginEntity
	Settings    map[string]any              `json:"settings"`
	TenantID    string                      `json:"tenant_id"`
	PluginID    string                      `json:"plugin_id"`
	ExpiredAt   time.Time                   `json:"expired_at"`
	Declaration EndpointProviderDeclaration `json:"declaration"`
}
type EndpointEntityWithInstance struct {
	*EndpointEntity
	Name    string `json:"path"`
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	HookID  string `json:"hook_id"`
}
