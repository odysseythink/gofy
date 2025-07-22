package exceptions

type AppModelConfigBrokenError struct {
	*ValueError
}

func NewAppModelConfigBrokenError(desc string) *AppModelConfigBrokenError {
	return &AppModelConfigBrokenError{
		ValueError: &ValueError{
			description: desc,
		},
	}
}

func (e *AppModelConfigBrokenError) Error() string {
	return "app model config broken error:" + e.description
}
