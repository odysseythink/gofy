package basenodes

import (
	"fmt"
)

type BaseNodeError struct {
	desc string
}

func NewBaseNodeError(desc string) *BaseNodeError {
	return &BaseNodeError{
		desc: desc,
	}
}
func (e *BaseNodeError) Error() string {
	return fmt.Sprintf("base node error:%s", e.desc)
}

type DefaultValueTypeError struct {
	*BaseNodeError
}

func NewDefaultValueTypeError(desc string) *DefaultValueTypeError {
	return &DefaultValueTypeError{
		BaseNodeError: NewBaseNodeError(desc),
	}
}
