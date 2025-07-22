package core

type ModelStatus string

/*
Enum class for model status.
*/
const (
	ModelStatus_ACTIVE         ModelStatus = "active"
	ModelStatus_NO_CONFIGURE   ModelStatus = "no-configure"
	ModelStatus_QUOTA_EXCEEDED ModelStatus = "quota-exceeded"
	ModelStatus_NO_PERMISSION  ModelStatus = "no-permission"
	ModelStatus_DISABLED       ModelStatus = "disabled"
)
