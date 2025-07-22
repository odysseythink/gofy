package questionclassifier

import (
	"fmt"
)

type QuestionClassifierNodeError struct {
	desc string
}

func NewQuestionClassifierNodeError(desc string) *QuestionClassifierNodeError {
	return &QuestionClassifierNodeError{
		desc: desc,
	}
}
func (e *QuestionClassifierNodeError) Error() string {
	return fmt.Sprintf("question  classifier node error:%s", e.desc)
}

func NewInvalidModelTypeError(desc string) *QuestionClassifierNodeError {
	desc = "invalid model type:" + desc
	return &QuestionClassifierNodeError{
		desc: desc,
	}
}
