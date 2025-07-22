package core

type QuotaUnit string

const (
	QuotaUnit_TIMES   QuotaUnit = "times"
	QuotaUnit_TOKENS  QuotaUnit = "tokens"
	QuotaUnit_CREDITS QuotaUnit = "credits"
)

type SystemConfigurationStatus string

const (
	/*
	   // Enum class for system configuration status.
	*/

	SystemConfigurationStatus_ACTIVE         SystemConfigurationStatus = "active"
	SystemConfigurationStatus_QUOTA_EXCEEDED SystemConfigurationStatus = "quota-exceeded"
	SystemConfigurationStatus_UNSUPPORTED    SystemConfigurationStatus = "unsupported"
)
