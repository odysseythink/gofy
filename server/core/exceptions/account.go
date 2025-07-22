package exceptions

type AccountNotFoundError struct {
	*BaseServiceError
}

type AccountRegisterError struct {
	*BaseServiceError
}
type AccountLoginError struct {
	*BaseServiceError
}

type AccountPasswordError struct {
	*BaseServiceError
}

type LinkAccountIntegrateError struct {
	*BaseServiceError
}

type TenantNotFoundError struct {
	*BaseServiceError
}

type AccountAlreadyInTenantError struct {
	*BaseServiceError
}

type InvalidActionError struct {
	*BaseServiceError
}

type CannotOperateSelfError struct {
	*BaseServiceError
}

type NoPermissionError struct {
	*BaseServiceError
}

type MemberNotInTenantError struct {
	*BaseServiceError
}

type RoleAlreadyAssignedError struct {
	*BaseServiceError
}

type RateLimitExceededError struct {
	*BaseServiceError
}

func NewAccountNotFoundError(desc string) *AccountNotFoundError {
	desc = "account not found error:" + desc
	return &AccountNotFoundError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewAccountRegisterError(desc string) *AccountRegisterError {
	desc = "account register error:" + desc
	return &AccountRegisterError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewAccountLoginError(desc string) *AccountLoginError {
	desc = "account login error:" + desc
	return &AccountLoginError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewAccountPasswordError(desc string) *AccountPasswordError {
	desc = "account password error:" + desc
	return &AccountPasswordError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}

func NewLinkAccountIntegrateError(desc string) *LinkAccountIntegrateError {
	desc = "link account integrate error:" + desc
	return &LinkAccountIntegrateError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewTenantNotFoundError(desc string) *TenantNotFoundError {
	desc = "tenant not found error:" + desc
	return &TenantNotFoundError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewAccountAlreadyInTenantError(desc string) *AccountAlreadyInTenantError {
	desc = "account already intenant error:" + desc
	return &AccountAlreadyInTenantError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewInvalidActionError(desc string) *InvalidActionError {
	desc = "invalid action error:" + desc
	return &InvalidActionError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewCannotOperateSelfError(desc string) *CannotOperateSelfError {
	desc = "cannot operate self error:" + desc
	return &CannotOperateSelfError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewNoPermissionError(desc string) *NoPermissionError {
	desc = "no permission error:" + desc
	return &NoPermissionError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewMemberNotInTenantError(desc string) *MemberNotInTenantError {
	desc = "member not in tenant error:" + desc
	return &MemberNotInTenantError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewRoleAlreadyAssignedError(desc string) *RoleAlreadyAssignedError {
	desc = "role already assigned error:" + desc
	return &RoleAlreadyAssignedError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
func NewRateLimitExceededError(desc string) *RateLimitExceededError {
	desc = "rate limit exceeded error:" + desc
	return &RateLimitExceededError{
		BaseServiceError: NewBaseServiceError(desc),
	}
}
