package exceptions

import "fmt"

type BaseServiceError struct {
	*ValueError
}

func NewBaseServiceError(description string) *BaseServiceError {
	return &BaseServiceError{
		ValueError: &ValueError{
			description: description,
		},
	}
}

func (e *BaseServiceError) Error() string {
	return fmt.Sprintf("base service error:%s", e.description)
}

func NewWorkSpaceNotAllowedCreateError() *BaseServiceError {
	return NewBaseServiceError("work space not allowed create error")
}

func NewWorkSpaceNotFoundError() *BaseServiceError {
	return NewBaseServiceError("work space not found error")
}
