package exceptions

type LastConversationNotExistsError struct {
	*ValueError
}

func NewLastConversationNotExistsError(desc string) *LastConversationNotExistsError {
	return &LastConversationNotExistsError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}

func (e *LastConversationNotExistsError) Error() string {
	return "last conversation not exists error:" + e.description
}

type ConversationNotExistsError struct {
	*ValueError
}

func NewConversationNotExistsError(desc string) *ConversationNotExistsError {
	return &ConversationNotExistsError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}

func (e *ConversationNotExistsError) Error() string {
	return "conversation not exists error:" + e.description
}
