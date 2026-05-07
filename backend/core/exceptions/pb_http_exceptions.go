package exceptions

import (
	"net/http"

	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
)

func NewInvalidArgsPbHttpExp(desc string) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "invalid_args",
		Message: desc,
	}
}

func NewUnauthorizedPbHttpExp(desc string) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Status:  http.StatusUnauthorized,
		Code:    "unauthorized",
		Message: desc,
	}
}
func NewForbiddenPbHttpExp(desc string) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Status:  http.StatusForbidden,
		Code:    "forbidden",
		Message: desc,
	}
}
func NewInternalServerPbHttpExp(desc string) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Status:  http.StatusInternalServerError,
		Code:    "internal_server_error",
		Message: desc,
	}
}

func NewDraftWorkflowNotExistPbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "Draft workflow need to be initialized."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "draft_workflow_not_exist",
		Message: desc,
	}
}

func NewAccountUpdatePbHttpExp(desc string) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "account_update_error",
		Message: desc,
	}
}

func NewAccountNotInitializedPbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "The account has not been initialized yet. Please proceed with the initialization process first."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "account_not_initialized",
		Message: desc,
	}
}

func NewRepeatPasswordNotMatchPbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "New password and repeat password does not match."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "repeat_password_not_match",
		Message: desc,
	}
}

func NewCurrentPasswordIncorrectPbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "Current password is incorrect."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "current_password_incorrect",
		Message: desc,
	}
}

func NewProviderRequestFailedPbHttpExp(desc string) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "provider_request_failed",
		Message: desc,
	}
}

func NewInvalidInvitationCodePbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "Invalid invitation code."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "invalid_invitation_code",
		Message: desc,
	}
}

func NewAccountAlreadyInitedPbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "The account has been initialized. Please refresh the page."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "account_already_inited",
		Message: desc,
	}
}

func NewInvalidAccountDeletionCodePbHttpExp(desc string) *pbexceptions.HTTPException {
	if desc == "" {
		desc = "Invalid account deletion code."
	}
	return &pbexceptions.HTTPException{
		Status:  http.StatusBadRequest,
		Code:    "invalid_account_deletion_code",
		Message: desc,
	}
}
