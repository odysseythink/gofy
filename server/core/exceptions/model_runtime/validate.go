package modelruntime

import "fmt"

type CredentialsValidateFailedError struct {
	description string
}

func NewCredentialsValidateFailedError(desc string) *CredentialsValidateFailedError {
	return &CredentialsValidateFailedError{
		description: desc,
	}
}
func (e *CredentialsValidateFailedError) Error() string {
	return fmt.Sprintf("credentials validate failed:%s", e.description)
}

func (e *CredentialsValidateFailedError) Description() string {
	return e.description
}
