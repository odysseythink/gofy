package exceptions

import "fmt"

type TaskPipilineError struct {
	*ValueError
}

type RecordNotFoundError struct {
	*TaskPipilineError
}

func NewRecordNotFoundError(record_name string, record_id string) *RecordNotFoundError {
	return &RecordNotFoundError{
		TaskPipilineError: &TaskPipilineError{
			ValueError: NewValueError(fmt.Sprintf("%s with id %s not found", record_name, record_id)),
		},
	}
}

type WorkflowRunNotFoundError struct {
	*RecordNotFoundError
}

func NewWorkflowRunNotFoundError(workflow_run_id string) *WorkflowRunNotFoundError {
	return &WorkflowRunNotFoundError{
		RecordNotFoundError: NewRecordNotFoundError("WorkflowRun", workflow_run_id),
	}
}

type WorkflowNodeExecutionNotFoundError struct {
	*RecordNotFoundError
}

func NewWorkflowNodeExecutionNotFoundError(workflow_node_execution_id string) *WorkflowNodeExecutionNotFoundError {
	return &WorkflowNodeExecutionNotFoundError{
		RecordNotFoundError: NewRecordNotFoundError("WorkflowNodeExecution", workflow_node_execution_id),
	}
}
