package app

import "fmt"

type WorkflowHashNotEqualError struct {
	desc string
}

func NewWorkflowHashNotEqualError(desc string) *WorkflowHashNotEqualError {
	return &WorkflowHashNotEqualError{
		desc: desc,
	}
}
func (e *WorkflowHashNotEqualError) Error() string {
	return fmt.Sprintf("workflow hash not equal error:%s", e.desc)
}

type MoreLikeThisDisabledError struct {
	desc string
}

func NewMoreLikeThisDisabledError(desc string) *MoreLikeThisDisabledError {
	return &MoreLikeThisDisabledError{
		desc: desc,
	}
}
func (e *MoreLikeThisDisabledError) Error() string {
	return fmt.Sprintf("more like this disabled error:%s", e.desc)
}
