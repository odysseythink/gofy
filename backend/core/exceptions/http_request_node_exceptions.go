package exceptions

import "fmt"

type HttpRequestNodeError struct {
	description string
}

func NewHttpRequestNodeError(desc string) *HttpRequestNodeError {
	return &HttpRequestNodeError{
		description: desc,
	}
}
func (e *HttpRequestNodeError) Error() string {
	return fmt.Sprintf("http request node error:%s", e.description)
}

func NewAuthorizationConfigError(desc string) *HttpRequestNodeError {
	desc = "authorization config error:" + desc
	return NewHttpRequestNodeError(desc)
}

func NewFileFetchError(desc string) *HttpRequestNodeError {
	desc = "file fetch error:" + desc
	return NewHttpRequestNodeError(desc)
}

func NewInvalidHttpMethodError(desc string) *HttpRequestNodeError {
	desc = "invalid http method error:" + desc
	return NewHttpRequestNodeError(desc)
}

func NewResponseSizeError(desc string) *HttpRequestNodeError {
	desc = "response size error:" + desc
	return NewHttpRequestNodeError(desc)
}

func NewRequestBodyError(desc string) *HttpRequestNodeError {
	desc = "response body error:" + desc
	return NewHttpRequestNodeError(desc)
}

func NewInvalidURLError(desc string) *HttpRequestNodeError {
	desc = "invalid url error:" + desc
	return NewHttpRequestNodeError(desc)
}
