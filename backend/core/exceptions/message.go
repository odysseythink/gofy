package exceptions

import "fmt"

type FirstMessageNotExistsError struct {
	*ValueError
}

func NewFirstMessageNotExistsError(desc string) *FirstMessageNotExistsError {
	return &FirstMessageNotExistsError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *FirstMessageNotExistsError) Error() string {
	return fmt.Sprintf("first message not exists error:%s", e.description)
}

type LastMessageNotExistsError struct {
	*ValueError
}

func NewLastMessageNotExistsError(desc string) *LastMessageNotExistsError {
	return &LastMessageNotExistsError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *LastMessageNotExistsError) Error() string {
	return fmt.Sprintf("last message not exists error:%s", e.description)
}

type MessageNotExistsError struct {
	*ValueError
}

func NewMessageNotExistsError(desc string) *MessageNotExistsError {
	return &MessageNotExistsError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *MessageNotExistsError) Error() string {
	return fmt.Sprintf("message not exists error:%s", e.description)
}

type SuggestedQuestionsAfterAnswerDisabledError struct {
	*ValueError
}

func NewSuggestedQuestionsAfterAnswerDisabledError(desc string) *SuggestedQuestionsAfterAnswerDisabledError {
	return &SuggestedQuestionsAfterAnswerDisabledError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}
func (e *SuggestedQuestionsAfterAnswerDisabledError) Error() string {
	return fmt.Sprintf("suggested questions after answer disabled error:%s", e.description)
}
