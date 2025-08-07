package provider

import (
	parameterenumtypes "mlib.com/gofy/server/enum_types/parameter"
)

type ProviderQuotaType string

const (
	ProviderQuota_PAID ProviderQuotaType = "paid"
	// """hosted paid quota"""

	ProviderQuota_FREE ProviderQuotaType = "free"
	// """third-party free quota"""

	ProviderQuota_TRIAL ProviderQuotaType = "trial"
	// """hosted trial quota"""
)

type QuotaUnitType string

const (
	QuotaUnit_TIMES   QuotaUnitType = "times"
	QuotaUnit_TOKENS  QuotaUnitType = "tokens"
	QuotaUnit_CREDITS QuotaUnitType = "credits"
)

type SystemConfigurationStatusType string

const (
	/*
	   // Enum class for system configuration status.
	*/

	SystemConfigurationStatus_ACTIVE         SystemConfigurationStatusType = "active"
	SystemConfigurationStatus_QUOTA_EXCEEDED SystemConfigurationStatusType = "quota-exceeded"
	SystemConfigurationStatus_UNSUPPORTED    SystemConfigurationStatusType = "unsupported"
)

type BasicProviderConfigType string

const (
	BasicProviderConfig_SECRET_INPUT   = BasicProviderConfigType(parameterenumtypes.CommonParameter_SECRET_INPUT)
	BasicProviderConfig_TEXT_INPUT     = BasicProviderConfigType(parameterenumtypes.CommonParameter_TEXT_INPUT)
	BasicProviderConfig_SELECT         = BasicProviderConfigType(parameterenumtypes.CommonParameter_SELECT)
	BasicProviderConfig_BOOLEAN        = BasicProviderConfigType(parameterenumtypes.CommonParameter_BOOLEAN)
	BasicProviderConfig_APP_SELECTOR   = BasicProviderConfigType(parameterenumtypes.CommonParameter_APP_SELECTOR)
	BasicProviderConfig_MODEL_SELECTOR = BasicProviderConfigType(parameterenumtypes.CommonParameter_MODEL_SELECTOR)
	BasicProviderConfig_TOOLS_SELECTOR = BasicProviderConfigType(parameterenumtypes.CommonParameter_TOOLS_SELECTOR)
)

type ProviderType string

const (
	Provider_CUSTOM ProviderType = "custom"
	Provider_SYSTEM ProviderType = "system"
)
