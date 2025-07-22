package enumtypes

import "slices"

type AccountStatus string

const (
	AccountStatusPENDING       AccountStatus = "pending"
	AccountStatusUNINITIALIZED AccountStatus = "uninitialized"
	AccountStatusACTIVE        AccountStatus = "active"
	AccountStatusBANNED        AccountStatus = "banned"
	AccountStatusCLOSED        AccountStatus = "closed"
)

type TenantStatus string

const (
	TenantStatus_NORMAL  TenantStatus = "normal"
	TenantStatus_ARCHIVE TenantStatus = "archive"
)

type TenantAccountRole string

const (
	TenantAccountRole_OWNER            TenantAccountRole = "owner"
	TenantAccountRole_ADMIN            TenantAccountRole = "admin"
	TenantAccountRole_EDITOR           TenantAccountRole = "editor"
	TenantAccountRole_NORMAL           TenantAccountRole = "normal"
	TenantAccountRole_DATASET_OPERATOR TenantAccountRole = "dataset_operator"
)

func (role TenantAccountRole) IsValidRole() bool {
	return slices.Contains([]TenantAccountRole{
		TenantAccountRole_OWNER,
		TenantAccountRole_ADMIN,
		TenantAccountRole_EDITOR,
		TenantAccountRole_NORMAL,
		TenantAccountRole_DATASET_OPERATOR,
	}, role)
}
func (role TenantAccountRole) IsPrivilegedRole() bool {
	return slices.Contains([]TenantAccountRole{TenantAccountRole_OWNER, TenantAccountRole_ADMIN}, role)
}
func (role TenantAccountRole) IsAdminRole() bool {
	return role == TenantAccountRole_ADMIN
}
func (role TenantAccountRole) IsNonOwnerRole() bool {
	return slices.Contains([]TenantAccountRole{
		TenantAccountRole_ADMIN,
		TenantAccountRole_EDITOR,
		TenantAccountRole_NORMAL,
		TenantAccountRole_DATASET_OPERATOR,
	}, role)
}
func (role TenantAccountRole) IsEditingRole() bool {
	return slices.Contains([]TenantAccountRole{TenantAccountRole_OWNER, TenantAccountRole_ADMIN, TenantAccountRole_EDITOR}, role)
}
func (role TenantAccountRole) IsDatasetEditRole() bool {
	return slices.Contains([]TenantAccountRole{
		TenantAccountRole_OWNER,
		TenantAccountRole_ADMIN,
		TenantAccountRole_EDITOR,
		TenantAccountRole_DATASET_OPERATOR,
	}, role)
}

type TenantAccountJoinRole string

var (
	TenantAccountJoinRole_OWNER            TenantAccountJoinRole = "owner"
	TenantAccountJoinRole_ADMIN            TenantAccountJoinRole = "admin"
	TenantAccountJoinRole_NORMAL           TenantAccountJoinRole = "normal"
	TenantAccountJoinRole_DATASET_OPERATOR TenantAccountJoinRole = "dataset_operator"
)

func (role TenantAccountJoinRole) Valid() bool {
	return slices.Contains([]TenantAccountJoinRole{
		TenantAccountJoinRole_OWNER,
		TenantAccountJoinRole_ADMIN,
		TenantAccountJoinRole_NORMAL,
		TenantAccountJoinRole_DATASET_OPERATOR,
	}, role)
}
