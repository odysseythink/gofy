package code

import "fmt"

type CodeNodeError struct {
	description string
}

func NewCodeNodeError(desc string) *CodeNodeError {
	return &CodeNodeError{
		description: desc,
	}
}
func (e *CodeNodeError) Error() string {
	return fmt.Sprintf("invalid code node value:%s", e.description)
}

func (e *CodeNodeError) Description() string {
	return e.description
}

func NewOutputValidationError(desc string) *CodeNodeError {
	desc = "output validation error:" + desc
	return NewCodeNodeError(desc)
}

func NewDepthLimitError(desc string) *CodeNodeError {
	desc = "depth limit error:" + desc
	return NewCodeNodeError(desc)
}

type CodeExecutionError struct {
	description string
}

func NewCodeExecutionError(desc string) *CodeExecutionError {
	return &CodeExecutionError{
		description: desc,
	}
}
func (e *CodeExecutionError) Error() string {
	return fmt.Sprintf("code execution error:%s", e.description)
}

func (e *CodeExecutionError) Description() string {
	return e.description
}
