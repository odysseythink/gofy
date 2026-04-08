package model

type ModelStatusType string

/*
Enum class for model status.
*/
const (
	ModelStatus_ACTIVE         ModelStatusType = "active"
	ModelStatus_NO_CONFIGURE   ModelStatusType = "no-configure"
	ModelStatus_QUOTA_EXCEEDED ModelStatusType = "quota-exceeded"
	ModelStatus_NO_PERMISSION  ModelStatusType = "no-permission"
	ModelStatus_DISABLED       ModelStatusType = "disabled"
)
