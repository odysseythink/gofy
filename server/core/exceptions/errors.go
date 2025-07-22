package exceptions

import "fmt"

type ValueError struct {
	description string
}

func NewValueError(desc string) *ValueError {
	return &ValueError{
		description: desc,
	}
}
func (e *ValueError) Error() string {
	return fmt.Sprintf("invalid value:%s", e.description)
}

type PrivkeyNotFoundError struct {
	*ValueError
}

func NewPrivkeyNotFoundError(desc string) *PrivkeyNotFoundError {
	return &PrivkeyNotFoundError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *PrivkeyNotFoundError) Error() string {
	return fmt.Sprintf("privkey not found error:%s", e.description)
}

type NotImplementedError struct {
	*ValueError
}

func NewNotImplementedError(desc string) *NotImplementedError {
	return &NotImplementedError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("not implemented:%s", e.description)
}

type OutputParserError struct {
	*ValueError
}

func NewOutputParserError(desc string) *OutputParserError {
	return &OutputParserError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *OutputParserError) Error() string {
	return fmt.Sprintf("output parser error:%s", e.description)
}

/*
Custom exception raised when the quota for an app has been exceeded.
*/
type AppInvokeQuotaExceededError struct {
	*ValueError
}

func NewAppInvokeQuotaExceededError() *AppInvokeQuotaExceededError {
	return &AppInvokeQuotaExceededError{
		ValueError: NewValueError("App Invoke Quota Exceeded"),
	}
}

type LLMError struct {
	/*Base class for all LLM exceptions.*/

	*ValueError
}

func NewLLMError(desc string) *LLMError {
	return &LLMError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}

func (e *LLMError) Error() string {
	return "llm error:" + e.description
}

type LLMBadRequestError struct {
	/*Raised when the LLM returns bad request.*/
	*LLMError
	// description = "Bad Request"
}

func NewLLMBadRequestError(desc string) *LLMBadRequestError {
	if desc == "" {
		desc = "Bad Request"
	}
	return &LLMBadRequestError{
		LLMError: NewLLMError(desc),
	}
}

type ProviderTokenNotInitError struct {
	/*
	   Custom exception raised when the provider token is not initialized.
	*/
	*ValueError
}

func NewProviderTokenNotInitError(desc string) *ProviderTokenNotInitError {
	return &ProviderTokenNotInitError{
		ValueError: NewValueError(desc),
	}
}
func (e *ProviderTokenNotInitError) Error() string {
	return fmt.Sprintf("provider token not init error:%s", e.description)
}

type QuotaExceededError struct {
	/*
	   Custom exception raised when the quota for a provider has been exceeded.
	*/
	*ValueError
}

func NewQuotaExceededError(desc string) *QuotaExceededError {
	return &QuotaExceededError{
		ValueError: NewValueError(desc),
	}
}
func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("Quota Exceeded:%s", e.description)
}

type ModelCurrentlyNotSupportError struct {
	/*
	   Custom exception raised when the model not support
	*/
	*ValueError
}

func NewModelCurrentlyNotSupportError(desc string) *ModelCurrentlyNotSupportError {
	return &ModelCurrentlyNotSupportError{
		ValueError: NewValueError(desc),
	}
}
func (e *ModelCurrentlyNotSupportError) Error() string {
	return fmt.Sprintf("Model Currently Not Support:%s", e.description)
}

type RuntimeError struct {
	*ValueError
}

func NewRuntimeError(desc string) *RuntimeError {
	return &RuntimeError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}

func (e *RuntimeError) Error() string {
	return "runtime error:" + e.description
}

type GenerateTaskStoppedError struct {
}

func (e *GenerateTaskStoppedError) Error() string {
	return "generate task stopped error"
}

func (e *GenerateTaskStoppedError) Description() string {
	return "generate task stopped error"
}

type ValidationError struct {
	*ValueError
}

func NewValidationError(desc string) *ValidationError {
	return &ValidationError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}

func (e *ValidationError) Error() string {
	return "validation error:" + e.description
}
