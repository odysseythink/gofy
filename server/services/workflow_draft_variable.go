package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mlib.com/gofy/server/constants"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

// Constants from the Python code
const (
	MIN_SELECTORS_LENGTH          = 2
	CONVERSATION_VARIABLE_NODE_ID = "conversation"
	SYSTEM_VARIABLE_NODE_ID       = "system"
)

// WorkflowDraftVariableList represents a list of workflow draft variables with total count
type WorkflowDraftVariableList struct {
	Variables []*models.WorkflowDraftVariable `json:"variables"`
	Total     *int64                          `json:"total"`
}

// WorkflowDraftVariableService handles workflow draft variable operations
type WorkflowDraftVariableService struct {
}

// GetVariable retrieves a workflow draft variable by ID
func (s *WorkflowDraftVariableService) GetVariable(variableID string) (*models.WorkflowDraftVariable, error) {
	var variable models.WorkflowDraftVariable
	err := dbengine.Instance().DB.Model(&models.WorkflowDraftVariable{}).
		Where("id = ?", variableID).
		First(&variable).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &variable, nil
}

// GetDraftVariablesBySelectors retrieves draft variables by selectors
func (s *WorkflowDraftVariableService) GetDraftVariablesBySelectors(appID string, selectors [][]string) ([]*models.WorkflowDraftVariable, error) {
	if len(selectors) == 0 {
		return []*models.WorkflowDraftVariable{}, nil
	}

	// Build OR conditions for each selector
	var conditions []string
	var args []interface{}

	for _, selector := range selectors {
		if len(selector) < MIN_SELECTORS_LENGTH {
			return nil, fmt.Errorf("invalid selector to get: %v", selector)
		}
		nodeID := selector[0]
		name := selector[1]
		conditions = append(conditions, "(node_id = ? AND name = ?)")
		args = append(args, nodeID, name)
	}

	// Build the complete query
	query := dbengine.Instance().DB.Model(&models.WorkflowDraftVariable{}).
		Where("app_id = ?", appID)

	if len(conditions) > 0 {
		orCondition := "(" + strings.Join(conditions, " OR ") + ")"
		query = query.Where(orCondition, args...)
	}

	var variables []*models.WorkflowDraftVariable
	err := query.Find(&variables).Error
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// ListVariablesWithoutValues retrieves variables without loading their values
func (s *WorkflowDraftVariableService) ListVariablesWithoutValues(appID string, page int, limit int) ([]*models.WorkflowDraftVariable, int64, error) {
	// Build base query
	query := dbengine.Instance().DB.Model(&models.WorkflowDraftVariable{}).
		Where("app_id = ?", appID)

	var total int64
	if page == 1 {
		if err := query.Count(&total).Error; err != nil {
			return nil, total, err
		}
	}

	// Execute query with pagination and ordering
	var variables []*models.WorkflowDraftVariable
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&variables).Error

	if err != nil {
		return nil, 0, err
	}

	return variables,
		total,
		nil
}

// ListNodeVariables lists variables for a specific node
func (s *WorkflowDraftVariableService) ListNodeVariables(appID string, nodeID string) ([]*models.WorkflowDraftVariable, error) {
	return s._listNodeVariables(appID, nodeID)
}

// ListConversationVariables lists conversation variables
func (s *WorkflowDraftVariableService) ListConversationVariables(appID string) ([]*models.WorkflowDraftVariable, error) {
	return s._listNodeVariables(appID, CONVERSATION_VARIABLE_NODE_ID)
}

// ListSystemVariables lists system variables
func (s *WorkflowDraftVariableService) ListSystemVariables(appID string) ([]*models.WorkflowDraftVariable, error) {
	return s._listNodeVariables(appID, SYSTEM_VARIABLE_NODE_ID)
}

// GetConversationVariable gets a conversation variable by name
func (s *WorkflowDraftVariableService) GetConversationVariable(appID string, name string) (*models.WorkflowDraftVariable, error) {
	return s._getVariable(appID, CONVERSATION_VARIABLE_NODE_ID, name)
}

// GetSystemVariable gets a system variable by name
func (s *WorkflowDraftVariableService) GetSystemVariable(appID string, name string) (*models.WorkflowDraftVariable, error) {
	return s._getVariable(appID, SYSTEM_VARIABLE_NODE_ID, name)
}

// GetNodeVariable gets a node variable by name
func (s *WorkflowDraftVariableService) GetNodeVariable(appID string, nodeID string, name string) (*models.WorkflowDraftVariable, error) {
	return s._getVariable(appID, nodeID, name)
}

// _listNodeVariables is the internal method to list variables for a node
func (s *WorkflowDraftVariableService) _listNodeVariables(appID string, nodeID string) ([]*models.WorkflowDraftVariable, error) {
	var variables []*models.WorkflowDraftVariable
	err := dbengine.Instance().DB.Model(&models.WorkflowDraftVariable{}).
		Where("app_id = ? AND node_id = ?", appID, nodeID).
		Order("created_at DESC").
		Find(&variables).Error

	if err != nil {
		return nil, err
	}

	return variables,
		nil
}

// _getVariable is the internal method to get a variable
func (s *WorkflowDraftVariableService) _getVariable(appID string, nodeID string, name string) (*models.WorkflowDraftVariable, error) {
	var variable models.WorkflowDraftVariable
	err := dbengine.Instance().DB.Model(&models.WorkflowDraftVariable{}).
		Where("app_id = ? AND node_id = ? AND name = ?", appID, nodeID, name).
		First(&variable).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &variable, nil
}

// UpdateVariable updates a workflow draft variable
func (s *WorkflowDraftVariableService) UpdateVariable(variable *models.WorkflowDraftVariable, name *string, value *string) (*models.WorkflowDraftVariable, error) {
	if !variable.Editable {
		return nil, exceptions.NewValueError(fmt.Sprintf("variable not support updating, id=%s", variable.ID))
	}

	now := time.Now()
	if name != nil {
		variable.Name = *name
	}
	if value != nil {
		variable.Value = *value
	}
	variable.LastEditedAt = &now

	if err := dbengine.Instance().DB.Save(variable).Error; err != nil {
		return nil, err
	}

	return variable, nil
}

// DeleteVariable deletes a workflow draft variable
func (s *WorkflowDraftVariableService) DeleteVariable(variable *models.WorkflowDraftVariable) error {
	return dbengine.Instance().DB.Delete(variable).Error
}

// ResetVariable resets a workflow draft variable
func (s *WorkflowDraftVariableService) ResetVariable(workflow *models.Workflow, variable *models.WorkflowDraftVariable) (*models.WorkflowDraftVariable, error) {
	variableType := s.getVariableType(variable)

	// Check if system variable is editable
	if variableType == "SYS" && !s.isSystemVariableEditable(variable.Name) {
		return nil, exceptions.NewValueError(fmt.Sprintf("cannot reset system variable, variable_id=%s", variable.ID))
	}

	if variableType == "CONVERSATION" {
		return s.resetConvVar(workflow, variable)
	} else {
		return s.resetNodeVarOrSysVar(workflow, variable)
	}
}

// getVariableType determines the type of variable
func (s *WorkflowDraftVariableService) getVariableType(variable *models.WorkflowDraftVariable) string {
	if variable.NodeID == CONVERSATION_VARIABLE_NODE_ID {
		return "CONVERSATION"
	} else if variable.NodeID == SYSTEM_VARIABLE_NODE_ID {
		return "SYS"
	}
	return "NODE"
}

// isSystemVariableEditable checks if a system variable is editable
func (s *WorkflowDraftVariableService) isSystemVariableEditable(name string) bool {
	// This is a simplified version - in the original Python code, this would check against a list
	// of editable system variables. For now, we'll assume most system variables are not editable.
	nonEditableVars := []string{"conversation_id", "app_id", "workflow_id"}
	for _, varName := range nonEditableVars {
		if name == varName {
			return false
		}
	}
	return true
}

// resetConvVar resets a conversation variable
func (s *WorkflowDraftVariableService) resetConvVar(workflow *models.Workflow, variable *models.WorkflowDraftVariable) (*models.WorkflowDraftVariable, error) {
	// Get conversation variables from workflow
	convVars := workflow.GetConversationVariables()

	// Find the matching conversation variable
	var matchingConvVar interface{}
	for _, convVar := range convVars {
		if convVar != nil && convVar.GetName() == variable.Name {
			matchingConvVar = convVar
			break
		}
	}

	if matchingConvVar == nil {
		// Conversation variable not found, delete the draft variable
		if err := s.DeleteVariable(variable); err != nil {
			return nil, err
		}
		mlog.Warningf("Conversation variable not found for draft variable, id=%s, name=%s", variable.ID, variable.Name)
		return nil, nil
	}

	// Convert variable to JSON string for storage
	valueStr := "{}"
	if matchingConvVar != nil {
		valueDict := matchingConvVar.(interface {
			ToDict(interface{}) map[string]interface{}
		}).ToDict(matchingConvVar)
		if valueBytes, err := json.Marshal(valueDict); err == nil {
			valueStr = string(valueBytes)
		}
	}

	// Reset the variable
	variable.Value = valueStr
	variable.LastEditedAt = nil

	if err := dbengine.Instance().DB.Save(variable).Error; err != nil {
		return nil, err
	}

	return variable, nil
}

// resetNodeVarOrSysVar resets a node or system variable
func (s *WorkflowDraftVariableService) resetNodeVarOrSysVar(workflow *models.Workflow, variable *models.WorkflowDraftVariable) (*models.WorkflowDraftVariable, error) {
	// If a variable does not allow updating, it makes no sense to reset it
	if !variable.Editable {
		return variable, nil
	}

	// No execution record for this variable, delete the variable instead
	if variable.NodeExecutionID == nil || *variable.NodeExecutionID == "" {
		if err := s.DeleteVariable(variable); err != nil {
			return nil, err
		}
		mlog.Warningf("draft variable has no node_execution_id, id=%s, name=%s", variable.ID, variable.Name)
		return nil, nil
	}

	// Get node execution
	var nodeExec models.WorkflowNodeExecution
	err := dbengine.Instance().DB.Model(&models.WorkflowNodeExecution{}).
		Where("id = ?", variable.NodeExecutionID).
		First(&nodeExec).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := s.DeleteVariable(variable); err != nil {
				return nil, err
			}
			mlog.Warningf("Node execution not found for draft variable, id=%s, name=%s, node_execution_id=%s",
				variable.ID, variable.Name, variable.NodeExecutionID)
			return nil, nil
		}
		return nil, err
	}

	outputsDict := nodeExec.OutputsDict()
	if outputsDict == nil {
		outputsDict = map[string]interface{}{}
	}

	var outputValue interface{}

	if s.getVariableType(variable) == "NODE" {
		outputValue = outputsDict[variable.Name]
	} else {
		outputValue = outputsDict[fmt.Sprintf("sys.%s", variable.Name)]
	}

	// Check if the variable was found in outputs
	if outputValue == nil {
		// If variable not found in execution data, delete the variable
		if err := s.DeleteVariable(variable); err != nil {
			return nil, err
		}
		return nil, nil
	}

	// Convert output value to JSON string
	valueStr := "{}"
	if outputValue != nil {
		if valueBytes, err := json.Marshal(outputValue); err == nil {
			valueStr = string(valueBytes)
		}
	}

	// Update the variable
	variable.Value = valueStr
	variable.LastEditedAt = nil

	if err := dbengine.Instance().DB.Save(variable).Error; err != nil {
		return nil, err
	}

	return variable, nil
}

// DeleteWorkflowVariables deletes all workflow variables for an app
func (s *WorkflowDraftVariableService) DeleteWorkflowVariables(appID string) error {
	return dbengine.Instance().DB.
		Where("app_id = ?", appID).
		Delete(&models.WorkflowDraftVariable{}).Error
}

// DeleteNodeVariables deletes all node variables for an app and node
func (s *WorkflowDraftVariableService) DeleteNodeVariables(appID string, nodeID string) error {
	return s._deleteNodeVariables(appID, nodeID)
}

// _deleteNodeVariables is the internal method to delete node variables
func (s *WorkflowDraftVariableService) _deleteNodeVariables(appID string, nodeID string) error {
	return dbengine.Instance().DB.
		Where("app_id = ? AND node_id = ?", appID, nodeID).
		Delete(&models.WorkflowDraftVariable{}).Error
}

// GetConversationIDFromDraftVariable gets conversation ID from draft variable
func (s *WorkflowDraftVariableService) GetConversationIDFromDraftVariable(appID string) (*string, error) {
	draftVar, err := s._getVariable(appID, SYSTEM_VARIABLE_NODE_ID, "conversation_id")
	if err != nil {
		return nil, err
	}
	if draftVar == nil {
		return nil, nil
	}

	// Parse the value as string - this is a simplified version
	// In the original Python code, this would parse the segment value
	// For now, we'll assume it's stored as a simple string
	if draftVar.Value == "" {
		mlog.Warningf("sys.conversation_id variable has empty value: app_id=%s, id=%s", appID, draftVar.ID)
		return nil, nil
	}

	return &draftVar.Value, nil
}

// GetOrCreateConversation gets or creates a conversation for debugging
func (s *WorkflowDraftVariableService) GetOrCreateConversation(accountID string, app *models.App, workflow *models.Workflow) (string, error) {
	convID, err := s.GetConversationIDFromDraftVariable(workflow.AppID)
	if err != nil {
		return "", err
	}

	if convID != nil {
		// Check if conversation exists
		var conversation models.Conversation
		err := dbengine.Instance().DB.Model(&models.Conversation{}).
			Where("id = ? AND app_id = ?", *convID, workflow.AppID).
			First(&conversation).Error
		if err == nil {
			// Conversation exists, return its ID
			return *convID, nil
		}
		if err != gorm.ErrRecordNotFound {
			return "", err
		}
	}

	// Create new conversation
	now := time.Now()
	conversation := &models.Conversation{
		ID:                      uuid.NewV4().String(),
		AppID:                   workflow.AppID,
		AppModelConfigID:        app.AppModelConfigID,
		ModelProvider:           "",
		ModelID:                 "",
		OverrideModelConfigsStr: "{}",
		Mode:                    app.Mode,
		Name:                    "Draft Debugging Conversation",
		InputsJson:              "{}",
		Introduction:            "",
		SystemInstruction:       "",
		SystemInstructionTokens: 0,
		Status:                  "normal",
		InvokeFrom:              "debugger",
		FromSource:              "console",
		FromEndUserID:           "",
		FromAccountID:           accountID,
		CreatedAt:               &now,
		UpdatedAt:               &now,
	}

	if err := dbengine.Instance().DB.Create(conversation).Error; err != nil {
		return "", err
	}

	return conversation.ID, nil
}

// PrefillConversationVariableDefaultValues prefills conversation variable default values
func (s *WorkflowDraftVariableService) PrefillConversationVariableDefaultValues(workflow *models.Workflow) error {
	draftConvVars := []*models.WorkflowDraftVariable{}

	// Get conversation variables from workflow
	convVars := workflow.GetConversationVariables()
	for _, convVar := range convVars {
		// Convert variable to JSON string for storage
		valueStr := "{}"
		if convVar != nil {
			valueDict := convVar.ToDict(convVar)
			if valueBytes, err := json.Marshal(valueDict); err == nil {
				valueStr = string(valueBytes)
			}
		}

		draftVar := &models.WorkflowDraftVariable{
			ID:          uuid.NewV4().String(),
			AppID:       workflow.AppID,
			NodeID:      CONVERSATION_VARIABLE_NODE_ID,
			Name:        convVar.GetName(),
			ValueType:   "conversation",
			Value:       valueStr,
			Description: "", // Could be extracted from convVar if available
			Editable:    true,
			Visible:     true,
			CreatedAt:   &time.Time{},
			UpdatedAt:   &time.Time{},
		}
		draftConvVars = append(draftConvVars, draftVar)
	}

	// Batch insert with ignore policy
	for _, draftVar := range draftConvVars {
		// Use FirstOrCreate to implement ignore policy
		var existing models.WorkflowDraftVariable
		err := dbengine.Instance().DB.
			Where("app_id = ? AND node_id = ? AND name = ?", draftVar.AppID, draftVar.NodeID, draftVar.Name).
			First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Record doesn't exist, create it
			if err := dbengine.Instance().DB.Create(draftVar).Error; err != nil {
				mlog.Warningf("Failed to create draft variable: %v", err)
				// Continue with other variables even if one fails
			}
		}
	}

	return nil
}

func (s *WorkflowDraftVariableService) GetVariableList(app_model *models.App, node_id string) ([]*models.WorkflowDraftVariable, error) {
	if node_id == constants.CONVERSATION_VARIABLE_NODE_ID {
		return s.ListConversationVariables(app_model.ID)
	} else if node_id == constants.SYSTEM_VARIABLE_NODE_ID {
		return s.ListSystemVariables(app_model.ID)
	} else {
		return s.ListNodeVariables(app_model.ID, node_id)
	}
}
