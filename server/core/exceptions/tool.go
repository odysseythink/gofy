package exceptions

type ToolProviderNotFoundError struct {
	*ValueError
}

func NewToolProviderNotFoundError(desc string) *ToolProviderNotFoundError {
	if desc == "" {
		desc = "Tool Provider Not Found"
	}
	return &ToolProviderNotFoundError{
		ValueError: NewValueError(desc),
	}
}

type ToolNotFoundError struct {
	*ValueError
}

func NewToolNotFoundError(desc string) *ToolNotFoundError {
	if desc == "" {
		desc = "Tool Not Found"
	}
	return &ToolNotFoundError{
		ValueError: NewValueError(desc),
	}
}

type ToolParameterValidationError struct {
	*ValueError
}

func NewToolParameterValidationError(desc string) *ToolParameterValidationError {
	if desc == "" {
		desc = "Tool Parameter Validation Error"
	}
	return &ToolParameterValidationError{
		ValueError: NewValueError(desc),
	}
}

type ToolProviderCredentialValidationError struct {
	*ValueError
}

func NewToolProviderCredentialValidationError(desc string) *ToolProviderCredentialValidationError {
	if desc == "" {
		desc = "Tool Provider Credential Validation Error"
	}
	return &ToolProviderCredentialValidationError{
		ValueError: NewValueError(desc),
	}
}

type ToolNotSupportedError struct {
	*ValueError
}

func NewToolNotSupportedError(desc string) *ToolNotSupportedError {
	if desc == "" {
		desc = "Tool Not Supported Error"
	}
	return &ToolNotSupportedError{
		ValueError: NewValueError(desc),
	}
}

type ToolInvokeError struct {
	*ValueError
}

func NewToolInvokeError(desc string) *ToolInvokeError {
	if desc == "" {
		desc = "Tool Invoke Error"
	}
	return &ToolInvokeError{
		ValueError: NewValueError(desc),
	}
}

type ToolApiSchemaError struct {
	*ValueError
}

func NewToolApiSchemaError(desc string) *ToolApiSchemaError {
	if desc == "" {
		desc = "Tool Api Schema Error"
	}
	return &ToolApiSchemaError{
		ValueError: NewValueError(desc),
	}
}

// class ToolEngineInvokeError(Exception):
//     meta: ToolInvokeMeta

//     def __init__(self, meta, **kwargs):
//         self.meta = meta
//         super().__init__(**kwargs)
