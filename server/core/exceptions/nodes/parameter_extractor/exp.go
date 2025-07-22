package parameterextractor

import "fmt"

type ParameterExtractorNodeError struct {
	desc string
}

func NewParameterExtractorNodeError(desc string) *ParameterExtractorNodeError {
	return &ParameterExtractorNodeError{
		desc: desc,
	}
}
func (e *ParameterExtractorNodeError) Error() string {
	return fmt.Sprintf("parameter extractor node error:%s", e.desc)
}

func NewInvalidModelTypeError(desc string) *ParameterExtractorNodeError {
	desc = "invalid model type error:" + desc
	return NewParameterExtractorNodeError(desc)
}
func NewModelSchemaNotFoundError(desc string) *ParameterExtractorNodeError {
	desc = "model schema not found error:" + desc
	return NewParameterExtractorNodeError(desc)
}
func NewInvalidInvokeResultError(desc string) *ParameterExtractorNodeError {
	desc = "invalid invoke result:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidTextContentTypeError(desc string) *ParameterExtractorNodeError {
	desc = "invalid text content type:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidNumberOfParametersError(desc string) *ParameterExtractorNodeError {
	desc = "invalid number of parameters:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewRequiredParameterMissingError(desc string) *ParameterExtractorNodeError {
	desc = "required parameter is missing:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidSelectValueError(desc string) *ParameterExtractorNodeError {
	desc = "invalid select value:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidNumberValueError(desc string) *ParameterExtractorNodeError {
	desc = "invalid number value:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidBoolValueError(desc string) *ParameterExtractorNodeError {
	desc = "invalid bool value:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidStringValueError(desc string) *ParameterExtractorNodeError {
	desc = "invalid string value:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidArrayValueError(desc string) *ParameterExtractorNodeError {
	desc = "invalid array value:" + desc
	return NewParameterExtractorNodeError(desc)
}

func NewInvalidModelModeError(desc string) *ParameterExtractorNodeError {
	desc = "invalid model mode:" + desc
	return NewParameterExtractorNodeError(desc)
}
