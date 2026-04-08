package dashscope

import (
	"errors"
	"fmt"
)

var (
	ErrAuthentication                    = errors.New("dashscope exception:authentication error")
	ErrInvalidParameter                  = errors.New("dashscope exception:invalid parameter")
	ErrInvalidTask                       = errors.New("dashscope exception:invalid task")
	ErrUnsupportedModel                  = errors.New("dashscope exception:unsupported model")
	ErrUnsupportedTask                   = errors.New("dashscope exception:unsupported task")
	ErrModelRequired                     = errors.New("dashscope exception:model required")
	ErrInvalidModel                      = errors.New("dashscope exception:invalid model")
	ErrInvalidInput                      = errors.New("dashscope exception:invalid input")
	ErrInvalidFileFormat                 = errors.New("dashscope exception:invalid file format")
	ErrUnsupportedApiProtocol            = errors.New("dashscope exception:unsupported api protocol")
	ErrNotImplemented                    = errors.New("dashscope exception:not implemented")
	ErrMultiInputsWithBinaryNotSupported = errors.New("dashscope exception:multi inputs with binary not supported")
	ErrUnexpectedMessageReceived         = errors.New("dashscope exception:unexpected message received")
	ErrUnsupportedData                   = errors.New("dashscope exception:unsupported data")
	ErrUnknownMessageReceived            = errors.New("dashscope exception:unknown message received")
	ErrInputDataRequired                 = errors.New("dashscope exception:input data required")
	ErrInputRequired                     = errors.New("dashscope exception:input required")
	ErrUnsupportedDataType               = errors.New("dashscope exception:unsupported data type")
	ErrServiceUnavailableError           = errors.New("dashscope exception:service unavailable error")
	ErrUnsupportedHTTPMethod             = errors.New("dashscope exception:unsupported http method")
	ErrAsyncTaskCreateFailed             = errors.New("dashscope exception:async task create failed")
	ErrUploadFileException               = errors.New("dashscope exception:upload file exception")
	ErrTimeoutException                  = errors.New("dashscope exception:timeout exception")
)

type DashScopeException struct {
	description string
}

func NewDashScopeException(desc string) *DashScopeException {
	return &DashScopeException{
		description: desc,
	}
}
func (e *DashScopeException) Error() string {
	return fmt.Sprintf("dashscope exception:%s", e.description)
}

func NewRequestFailure(
	request_id string,
	message string,
	name string,
	http_code string) *DashScopeException {
	desc := fmt.Sprintf("Request failed, request_id: %s, http_code: %s error_name: %s, error_message: %s", request_id, http_code, name, message)
	return NewDashScopeException(desc)
}

func NewAuthenticationError(desc string) *DashScopeException {
	desc = "authentication error:" + desc
	return NewDashScopeException(desc)
}
func NewInvalidParameter(desc string) *DashScopeException {
	desc = "invalid parameter:" + desc
	return NewDashScopeException(desc)
}
func NewInvalidTask(desc string) *DashScopeException {
	desc = "invalid task:" + desc
	return NewDashScopeException(desc)
}
func NewUnsupportedModel(desc string) *DashScopeException {
	desc = "unsupported model:" + desc
	return NewDashScopeException(desc)
}
func NewUnsupportedTask(desc string) *DashScopeException {
	desc = "unsupported task:" + desc
	return NewDashScopeException(desc)
}
func NewModelRequired(desc string) *DashScopeException {
	desc = "model required:" + desc
	return NewDashScopeException(desc)
}
func NewInvalidModel(desc string) *DashScopeException {
	desc = "invalid model:" + desc
	return NewDashScopeException(desc)
}
func NewInvalidInput(desc string) *DashScopeException {
	desc = "invalid input:" + desc
	return NewDashScopeException(desc)
}
func NewInvalidFileFormat(desc string) *DashScopeException {
	desc = "invalid file format:" + desc
	return NewDashScopeException(desc)
}
func NewUnsupportedApiProtocol(desc string) *DashScopeException {
	desc = "unsupported api protocol:" + desc
	return NewDashScopeException(desc)
}
func NewNotImplemented(desc string) *DashScopeException {
	desc = "not implemented:" + desc
	return NewDashScopeException(desc)
}
func NewMultiInputsWithBinaryNotSupported(desc string) *DashScopeException {
	desc = "multi inputs with binary not supported:" + desc
	return NewDashScopeException(desc)
}
func NewUnexpectedMessageReceived(desc string) *DashScopeException {
	desc = "unexpected message received:" + desc
	return NewDashScopeException(desc)
}
func NewUnsupportedData(desc string) *DashScopeException {
	desc = "unsupported data:" + desc
	return NewDashScopeException(desc)
}
func NewUnknownMessageReceived(desc string) *DashScopeException {
	desc = "unknown message received:" + desc
	return NewDashScopeException(desc)
}
func NewInputDataRequired(desc string) *DashScopeException {
	desc = "input data required:" + desc
	return NewDashScopeException(desc)
}
func NewInputRequired(desc string) *DashScopeException {
	desc = "input required:" + desc
	return NewDashScopeException(desc)
}
func NewUnsupportedDataType(desc string) *DashScopeException {
	desc = "unsupported data type:" + desc
	return NewDashScopeException(desc)
}
func NewServiceUnavailableError(desc string) *DashScopeException {
	desc = "service unavailable error:" + desc
	return NewDashScopeException(desc)
}
func NewUnsupportedHTTPMethod(desc string) *DashScopeException {
	desc = "unsupported http method:" + desc
	return NewDashScopeException(desc)
}
func NewAsyncTaskCreateFailed(desc string) *DashScopeException {
	desc = "async task create failed:" + desc
	return NewDashScopeException(desc)
}
func NewUploadFileException(desc string) *DashScopeException {
	desc = "upload file exception:" + desc
	return NewDashScopeException(desc)
}
func NewTimeoutException(desc string) *DashScopeException {
	desc = "timeout exception:" + desc
	return NewDashScopeException(desc)
}
