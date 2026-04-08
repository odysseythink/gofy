package exceptions

type GraphRunFailedError struct {
	err string
}

func NewGraphRunFailedError(err string) error {
	return &GraphRunFailedError{
		err: err,
	}
}
func (e *GraphRunFailedError) Error() string {
	return e.err
}

func (e *GraphRunFailedError) Description() string {
	return e.err
}
