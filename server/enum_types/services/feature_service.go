package services

type LicenseStatus string

const (
	LicenseStatus_NONE     LicenseStatus = "none"
	LicenseStatus_INACTIVE LicenseStatus = "inactive"
	LicenseStatus_ACTIVE   LicenseStatus = "active"
	LicenseStatus_EXPIRING LicenseStatus = "expiring"
	LicenseStatus_EXPIRED  LicenseStatus = "expired"
	LicenseStatus_LOST     LicenseStatus = "lost"
)
