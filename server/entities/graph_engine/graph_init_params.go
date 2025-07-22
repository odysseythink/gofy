package graphengine

import (
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
)

// GraphInitParams represents the initialization parameters for a graph.
type GraphInitParams struct {
	TenantID     string                  `json:"tenant_id"`
	AppID        string                  `json:"app_id"`
	WorkflowType models.WorkflowType     `json:"workflow_type"`
	WorkflowID   string                  `json:"workflow_id"`
	GraphConfig  map[string]any          `json:"graph_config"`
	UserID       string                  `json:"user_id"`
	UserFrom     models.UserFrom         `json:"user_from"`
	InvokeFrom   appenumtypes.InvokeFrom `json:"invoke_from"`
	CallDepth    int                     `json:"call_depth"`
}
