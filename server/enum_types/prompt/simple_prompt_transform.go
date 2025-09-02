package prompt

type ModelModeType string

const (
	ModelMode_COMPLETION ModelModeType = "completion"
	ModelMode_CHAT       ModelModeType = "chat"
)

func ValidateModelModeType(val string) bool {
	return ModelModeType(val) == ModelMode_COMPLETION || ModelModeType(val) == ModelMode_CHAT
}
