package modelruntime

import "fmt"

type InvokeError struct {
	description string
}

func NewInvokeError(desc string) *InvokeError {
	return &InvokeError{
		description: desc,
	}
}
func (e *InvokeError) Error() string {
	return fmt.Sprintf("invoke error:%s", e.description)
}

func (e *InvokeError) Description() string {
	return e.description
}

type InvokeConnectionError struct {
	*InvokeError
}

func NewInvokeConnectionError(desc string) *InvokeConnectionError {
	return &InvokeConnectionError{
		InvokeError: NewInvokeError("Connection Error"),
	}
}

type InvokeServerUnavailableError struct {
	*InvokeError
}

func NewInvokeServerUnavailableError(desc string) *InvokeServerUnavailableError {
	return &InvokeServerUnavailableError{
		InvokeError: NewInvokeError("Server Unavailable Error"),
	}
}

type InvokeRateLimitError struct {
	*InvokeError
}

func NewInvokeRateLimitError(desc string) *InvokeRateLimitError {
	return &InvokeRateLimitError{
		InvokeError: NewInvokeError("Rate Limit Error"),
	}
}

type InvokeAuthorizationError struct {
	*InvokeError
}

func NewInvokeAuthorizationError(desc string) *InvokeAuthorizationError {
	return &InvokeAuthorizationError{
		InvokeError: NewInvokeError("Incorrect model credentials provided, please check and try again. "),
	}
}

type InvokeBadRequestError struct {
	*InvokeError
}

func NewInvokeBadRequestError(desc string) *InvokeBadRequestError {
	return &InvokeBadRequestError{
		InvokeError: NewInvokeError("Bad Request Error"),
	}
}
