package workflow

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
)

// WorkflowExecutionRepository manages workflow run persistence.
type WorkflowExecutionRepository struct{}

// CreateWorkflowRun creates a new workflow run record.
func (r *WorkflowExecutionRepository) CreateWorkflowRun(
	tenantID, appID, workflowID, triggeredFrom, createdBy string,
	inputs map[string]any,
) (*models.WorkflowRun, error) {
	inputsJSON, _ := json.Marshal(inputs)
	run := &models.WorkflowRun{
		ID:            uuid.NewV4().String(),
		TenantID:      tenantID,
		AppID:         appID,
		WorkflowID:    workflowID,
		TriggeredFrom: triggeredFrom,
		Status:        "running",
		Inputs:        string(inputsJSON),
		CreatedByRole: models.CreatedByRole_ACCOUNT,
		CreatedBy:     createdBy,
	}
	if err := dbengine.Instance().DB.Create(run).Error; err != nil {
		return nil, err
	}
	return run, nil
}

// UpdateWorkflowRunStatus updates a run's status and result.
func (r *WorkflowExecutionRepository) UpdateWorkflowRunStatus(
	runID string, status string, outputs map[string]any,
	totalSteps, totalTokens int, elapsedTime float64, errMsg string,
) error {
	updates := map[string]any{
		"status":       status,
		"total_steps":  totalSteps,
		"total_tokens": totalTokens,
		"elapsed_time": elapsedTime,
	}
	if outputs != nil {
		outputsJSON, _ := json.Marshal(outputs)
		updates["outputs"] = string(outputsJSON)
	}
	if errMsg != "" {
		updates["error"] = errMsg
	}
	now := time.Now()
	updates["finished_at"] = &now
	return dbengine.Instance().DB.Model(&models.WorkflowRun{}).Where("id = ?", runID).Updates(updates).Error
}

// GetWorkflowRun retrieves a workflow run by ID.
func (r *WorkflowExecutionRepository) GetWorkflowRun(runID string) *models.WorkflowRun {
	var run models.WorkflowRun
	if err := dbengine.Instance().DB.Where("id = ?", runID).First(&run).Error; err != nil {
		return nil
	}
	return &run
}

// NodeExecutionRepository manages node execution persistence.
type NodeExecutionRepository struct{}

// CreateNodeExecution creates a node execution record.
func (r *NodeExecutionRepository) CreateNodeExecution(
	tenantID, appID, workflowID, workflowRunID, triggeredFrom string,
	nodeID string, nodeType nodesenumtypes.NodeType, title string,
	index int,
	predecessorNodeID string,
	inputs map[string]any,
) (*models.WorkflowNodeExecution, error) {
	inputsJSON, _ := json.Marshal(inputs)
	exec := &models.WorkflowNodeExecution{
		ID:                uuid.NewV4().String(),
		TenantID:          tenantID,
		AppID:             appID,
		WorkflowID:        workflowID,
		WorkflowRunID:     workflowRunID,
		TriggeredFrom:     triggeredFrom,
		NodeID:            nodeID,
		NodeType:          nodeType,
		Title:             title,
		Index:             index,
		PredecessorNodeID: predecessorNodeID,
		Inputs:            string(inputsJSON),
		Status:            "running",
	}
	if err := dbengine.Instance().DB.Create(exec).Error; err != nil {
		return nil, err
	}
	return exec, nil
}

// CompleteNodeExecution updates a node execution with results.
func (r *NodeExecutionRepository) CompleteNodeExecution(
	executionID string, status string,
	outputs, processData, metadata map[string]any,
	elapsedTime float64, errMsg string, tokens int,
) error {
	updates := map[string]any{
		"status":       status,
		"elapsed_time": elapsedTime,
	}
	if outputs != nil {
		outputsJSON, _ := json.Marshal(outputs)
		updates["outputs"] = string(outputsJSON)
	}
	if processData != nil {
		pdJSON, _ := json.Marshal(processData)
		updates["process_data"] = string(pdJSON)
	}
	if metadata != nil {
		mdJSON, _ := json.Marshal(metadata)
		updates["execution_metadata"] = string(mdJSON)
	}
	if errMsg != "" {
		updates["error"] = errMsg
	}
	now := time.Now()
	updates["finished_at"] = &now
	return dbengine.Instance().DB.Model(&models.WorkflowNodeExecution{}).Where("id = ?", executionID).Updates(updates).Error
}

// GetNodeExecutionsByRun returns all node executions for a run.
func (r *NodeExecutionRepository) GetNodeExecutionsByRun(workflowRunID string) []*models.WorkflowNodeExecution {
	var executions []*models.WorkflowNodeExecution
	dbengine.Instance().DB.Where("workflow_run_id = ?", workflowRunID).Order("`index` ASC").Find(&executions)
	return executions
}

// GetNodeExecution retrieves a specific node execution.
func (r *NodeExecutionRepository) GetNodeExecution(executionID string) *models.WorkflowNodeExecution {
	var exec models.WorkflowNodeExecution
	if err := dbengine.Instance().DB.Where("id = ?", executionID).First(&exec).Error; err != nil {
		return nil
	}
	return &exec
}

// DraftVariableRepository manages draft variable persistence.
type DraftVariableRepository struct{}

// GetDraftVariables returns draft variables for an app.
func (r *DraftVariableRepository) GetDraftVariables(appID string) []*models.WorkflowDraftVariable {
	var vars []*models.WorkflowDraftVariable
	dbengine.Instance().DB.Where("app_id = ?", appID).Order("node_id, name").Find(&vars)
	return vars
}

// SaveDraftVariable saves or updates a draft variable.
func (r *DraftVariableRepository) SaveDraftVariable(v *models.WorkflowDraftVariable) error {
	var existing models.WorkflowDraftVariable
	err := dbengine.Instance().DB.Where("app_id = ? AND node_id = ? AND name = ?", v.AppID, v.NodeID, v.Name).First(&existing).Error
	if err != nil {
		v.ID = uuid.NewV4().String()
		return dbengine.Instance().DB.Create(v).Error
	}
	return dbengine.Instance().DB.Model(&existing).Updates(map[string]any{
		"value":       v.Value,
		"value_type":  v.ValueType,
		"description": v.Description,
		"selector":    v.Selector,
		"visible":     v.Visible,
		"editable":    v.Editable,
	}).Error
}

// DeleteDraftVariables deletes all draft variables for an app.
func (r *DraftVariableRepository) DeleteDraftVariables(appID string) error {
	return dbengine.Instance().DB.Where("app_id = ?", appID).Delete(&models.WorkflowDraftVariable{}).Error
}

// HumanInputFormRepository manages human input form persistence.
type HumanInputFormRepository struct{}

// GetFormByRunAndNode returns the form for a specific workflow run and node.
func (r *HumanInputFormRepository) GetFormByRunAndNode(workflowRunID, nodeID string) *models.HumanInputForm {
	var form models.HumanInputForm
	if err := dbengine.Instance().DB.Where("workflow_run_id = ? AND node_id = ?", workflowRunID, nodeID).First(&form).Error; err != nil {
		return nil
	}
	return &form
}

// GetPendingForms returns all waiting forms for a workflow run.
func (r *HumanInputFormRepository) GetPendingForms(workflowRunID string) []*models.HumanInputForm {
	var forms []*models.HumanInputForm
	dbengine.Instance().DB.Where("workflow_run_id = ? AND status = ?", workflowRunID, "waiting").Find(&forms)
	return forms
}
