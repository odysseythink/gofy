package workflow

type SystemVariableKey string

const (
	/*
	   System Variables.
	*/

	SystemVariableKey_QUERY           SystemVariableKey = "query"
	SystemVariableKey_FILES           SystemVariableKey = "files"
	SystemVariableKey_CONVERSATION_ID SystemVariableKey = "conversation_id"
	SystemVariableKey_USER_ID         SystemVariableKey = "user_id"
	SystemVariableKey_DIALOGUE_COUNT  SystemVariableKey = "dialogue_count"
	SystemVariableKey_APP_ID          SystemVariableKey = "app_id"
	SystemVariableKey_WORKFLOW_ID     SystemVariableKey = "workflow_id"
	SystemVariableKey_WORKFLOW_RUN_ID SystemVariableKey = "workflow_run_id"
)
