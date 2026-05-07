package tools

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
)

// ToolManageService manages tool providers.
type ToolManageService struct{}

// ListBuiltinProviders returns all built-in tool providers.
func (s *ToolManageService) ListBuiltinProviders(tenantID string) []*models.BuiltinToolProvider {
	var providers []*models.BuiltinToolProvider
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&providers)
	return providers
}

// ListAPIProviders returns all API-based tool providers for a tenant.
func (s *ToolManageService) ListAPIProviders(tenantID string) []*models.ApiToolProvider {
	var providers []*models.ApiToolProvider
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&providers)
	return providers
}

// ListMCPProviders returns all MCP tool providers for a tenant.
func (s *ToolManageService) ListMCPProviders(tenantID string) []*models.MCPToolProvider {
	var providers []*models.MCPToolProvider
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&providers)
	return providers
}

// ListWorkflowProviders returns all workflow-based tool providers.
func (s *ToolManageService) ListWorkflowProviders(tenantID string) []*models.WorkflowToolProvider {
	var providers []*models.WorkflowToolProvider
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&providers)
	return providers
}

// GetToolLabels returns all tool labels for a tenant.
func (s *ToolManageService) GetToolLabels(tenantID string) []*models.ToolLabelBinding {
	var labels []*models.ToolLabelBinding
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Find(&labels)
	return labels
}
