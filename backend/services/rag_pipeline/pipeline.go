package ragpipeline

import (
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	uuid "github.com/satori/go.uuid"
)

// PipelineService manages RAG pipeline operations.
type PipelineService struct{}

// CreatePipeline creates a new RAG pipeline.
func (s *PipelineService) CreatePipeline(tenantID, name, description string, createdBy string) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{
		ID:          uuid.NewV4().String(),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		CreatedBy:   &createdBy,
	}
	if err := dbengine.Instance().DB.Create(pipeline).Error; err != nil {
		return nil, fmt.Errorf("failed to create pipeline: %w", err)
	}
	return pipeline, nil
}

// GetPipeline returns a pipeline by ID.
func (s *PipelineService) GetPipeline(pipelineID, tenantID string) *models.Pipeline {
	var pipeline models.Pipeline
	if err := dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", pipelineID, tenantID).First(&pipeline).Error; err != nil {
		return nil
	}
	return &pipeline
}

// ListPipelines returns all pipelines for a tenant.
func (s *PipelineService) ListPipelines(tenantID string, page, limit int) ([]*models.Pipeline, int64) {
	var pipelines []*models.Pipeline
	var total int64
	query := dbengine.Instance().DB.Where("tenant_id = ?", tenantID)
	query.Model(&models.Pipeline{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&pipelines)
	return pipelines, total
}

// DeletePipeline removes a pipeline.
func (s *PipelineService) DeletePipeline(pipelineID, tenantID string) error {
	return dbengine.Instance().DB.Where("id = ? AND tenant_id = ?", pipelineID, tenantID).Delete(&models.Pipeline{}).Error
}

// PublishPipeline marks a pipeline as published.
func (s *PipelineService) PublishPipeline(pipelineID, tenantID string) error {
	return dbengine.Instance().DB.Model(&models.Pipeline{}).Where("id = ? AND tenant_id = ?", pipelineID, tenantID).Update("is_published", true).Error
}

// ListBuiltInTemplates returns built-in pipeline templates.
func (s *PipelineService) ListBuiltInTemplates() []*models.PipelineBuiltInTemplate {
	var templates []*models.PipelineBuiltInTemplate
	dbengine.Instance().DB.Order("position ASC").Find(&templates)
	return templates
}

// ListCustomizedTemplates returns custom pipeline templates for a tenant.
func (s *PipelineService) ListCustomizedTemplates(tenantID string) []*models.PipelineCustomizedTemplate {
	var templates []*models.PipelineCustomizedTemplate
	dbengine.Instance().DB.Where("tenant_id = ?", tenantID).Order("position ASC").Find(&templates)
	return templates
}
