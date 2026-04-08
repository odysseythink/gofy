package variables

import (
	"fmt"
)

type VariableError struct {
	desc string
}

func NewVariableError(desc string) *VariableError {
	return &VariableError{
		desc: desc,
	}
}
func (e *VariableError) Error() string {
	return fmt.Sprintf("invalid variable:%s", e.desc)
}
