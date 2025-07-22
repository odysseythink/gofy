package llm

import "fmt"

type LLMNodeError struct {
	desc string
}

func NewLLMNodeError(desc string) *LLMNodeError {
	return &LLMNodeError{
		desc: desc,
	}
}
func (e *LLMNodeError) Error() string {
	return fmt.Sprintf("llm node error:%s", e.desc)
}

type VariableNotFoundError struct {
	*LLMNodeError
}

func NewVariableNotFoundError(desc string) *LLMNodeError {
	desc = "variable not found:" + desc
	return NewLLMNodeError(desc)
}

func NewInvalidContextStructureError(desc string) *LLMNodeError {
	desc = "invalid context structure error:" + desc
	return NewLLMNodeError(desc)
}

func NewInvalidVariableTypeError(desc string) *LLMNodeError {
	desc = "invalid variable type error:" + desc
	return NewLLMNodeError(desc)
}

func NewModelNotExistError(desc string) *LLMNodeError {
	desc = "model not exist error:" + desc
	return NewLLMNodeError(desc)
}

func NewLLMModeRequiredError(desc string) *LLMNodeError {
	desc = "llm mode required error:" + desc
	return NewLLMNodeError(desc)
}

func NewNoPromptFoundError(desc string) *LLMNodeError {
	desc = "no prompt found error:" + desc
	return NewLLMNodeError(desc)
}

func NewTemplateTypeNotSupportError(type_name string) *LLMNodeError {
	desc := fmt.Sprintf("Prompt type %s is not supported.", type_name)
	return NewLLMNodeError(desc)
}

func NewMemoryRolePrefixRequiredError(desc string) *LLMNodeError {
	desc = "memory role prefix required error:" + desc
	return NewLLMNodeError(desc)
}

func NewFileTypeNotSupportError(type_name string) *LLMNodeError {
	desc := fmt.Sprintf("%s type is not supported by this model", type_name)
	return NewLLMNodeError(desc)
}
