package workflow

import (
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
)

// SystemVariable defines a system variable with its type and description.
type SystemVariable struct {
	Key         workflowenumtypes.SystemVariableKey
	ValueType   string // string, array[file], number
	Description string
}

// AllSystemVariables returns all available system variables.
func AllSystemVariables() []SystemVariable {
	return []SystemVariable{
		{Key: workflowenumtypes.SystemVariableKey_QUERY, ValueType: "string", Description: "User query input"},
		{Key: workflowenumtypes.SystemVariableKey_FILES, ValueType: "array[file]", Description: "User uploaded files"},
		{Key: workflowenumtypes.SystemVariableKey_CONVERSATION_ID, ValueType: "string", Description: "Current conversation ID"},
		{Key: workflowenumtypes.SystemVariableKey_USER_ID, ValueType: "string", Description: "Current user ID"},
		{Key: workflowenumtypes.SystemVariableKey_DIALOGUE_COUNT, ValueType: "number", Description: "Number of dialogue turns"},
		{Key: workflowenumtypes.SystemVariableKey_APP_ID, ValueType: "string", Description: "Current app ID"},
		{Key: workflowenumtypes.SystemVariableKey_WORKFLOW_ID, ValueType: "string", Description: "Current workflow ID"},
		{Key: workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID, ValueType: "string", Description: "Current workflow run ID"},
	}
}

// GetSystemVariableType returns the value type for a system variable key.
func GetSystemVariableType(key workflowenumtypes.SystemVariableKey) string {
	for _, sv := range AllSystemVariables() {
		if sv.Key == key {
			return sv.ValueType
		}
	}
	return "string"
}

// BuildSystemVariableMap creates a map of system variables from execution context.
func BuildSystemVariableMap(
	query string,
	files []any,
	conversationID string,
	userID string,
	dialogueCount int,
	appID string,
	workflowID string,
	workflowRunID string,
) map[workflowenumtypes.SystemVariableKey]any {
	return map[workflowenumtypes.SystemVariableKey]any{
		workflowenumtypes.SystemVariableKey_QUERY:           query,
		workflowenumtypes.SystemVariableKey_FILES:           files,
		workflowenumtypes.SystemVariableKey_CONVERSATION_ID: conversationID,
		workflowenumtypes.SystemVariableKey_USER_ID:         userID,
		workflowenumtypes.SystemVariableKey_DIALOGUE_COUNT:  dialogueCount,
		workflowenumtypes.SystemVariableKey_APP_ID:          appID,
		workflowenumtypes.SystemVariableKey_WORKFLOW_ID:     workflowID,
		workflowenumtypes.SystemVariableKey_WORKFLOW_RUN_ID: workflowRunID,
	}
}
