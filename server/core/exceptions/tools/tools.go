package tools

import "fmt"

type InvokeModelError struct {
	description string
}

func NewInvokeModelError(desc string) *InvokeModelError {
	return &InvokeModelError{
		description: desc,
	}
}
func (e *InvokeModelError) Error() string {
	return fmt.Sprintf("invoke model error:%s", e.description)
}
