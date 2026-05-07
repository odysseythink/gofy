package http

import (
	"github.com/gin-gonic/gin"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	pbexceptions "github.com/odysseythink/gofy/backend/proto/exceptions"
)

type HTTPException interface {
	error
	Response(c *gin.Context)
	ToPbHttpException(HTTPException) *pbexceptions.HTTPException
}

type BaseHTTPException struct {
	*exceptions.ValueError
	code   string
	status int
}

func (e *BaseHTTPException) Response(c *gin.Context) {
	c.JSON(e.status, gin.H{
		"code":    e.code,
		"status":  e.status,
		"message": e.Error(),
	})
}

func (e *BaseHTTPException) ToPbHttpException(HTTPException) *pbexceptions.HTTPException {
	return &pbexceptions.HTTPException{
		Code:    e.code,
		Status:  int32(e.status),
		Message: e.Error(),
	}
}

type AlreadySetupError struct {
	*BaseHTTPException
}

func NewAlreadySetupError() *AlreadySetupError {
	return &AlreadySetupError{
		BaseHTTPException: &BaseHTTPException{
			code: "already_setup",
			ValueError: exceptions.NewValueError(
				"Gofy has been successfully installed. Please refresh the page or return to the dashboard homepage.",
			),
			status: 403,
		},
	}
}

type NotSetupError struct {
	*BaseHTTPException
}

func NewNotSetupError() *NotSetupError {
	return &NotSetupError{
		BaseHTTPException: &BaseHTTPException{
			code: "not_setup",
			ValueError: exceptions.NewValueError(
				`Gofy has not been initialized and installed yet. 
        Please proceed with the initialization and installation process first.`,
			),
			status: 401,
		},
	}
}

type NotInitValidateError struct {
	*BaseHTTPException
}

func NewNotInitValidateError() *NotInitValidateError {
	return &NotInitValidateError{
		BaseHTTPException: &BaseHTTPException{
			code: "not_init_validated",
			ValueError: exceptions.NewValueError(
				"Init validation has not been completed yet. Please proceed with the init validation process first.",
			),
			status: 401,
		},
	}
}

type InitValidateFailedError struct {
	*BaseHTTPException
}

func NewInitValidateFailedError() *InitValidateFailedError {
	return &InitValidateFailedError{
		BaseHTTPException: &BaseHTTPException{
			code: "init_validate_failed",
			ValueError: exceptions.NewValueError(
				"Init validation failed. Please check the password and try again.",
			),
			status: 401,
		},
	}
}

type AccountNotLinkTenantError struct {
	*BaseHTTPException
}

func NewAccountNotLinkTenantError(desc string) *AccountNotLinkTenantError {
	desc = "account not link tenant error:" + desc
	return &AccountNotLinkTenantError{
		BaseHTTPException: &BaseHTTPException{
			code: "account_not_link_tenant",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 403,
		},
	}
}

type AlreadyActivateError struct {
	*BaseHTTPException
}

func NewAlreadyActivateError() *AlreadyActivateError {
	return &AlreadyActivateError{
		BaseHTTPException: &BaseHTTPException{
			code: "already_activate",
			ValueError: exceptions.NewValueError(
				"Auth Token is invalid or account already activated, please check again.",
			),
			status: 403,
		},
	}
}

type NotAllowedCreateWorkspace struct {
	*BaseHTTPException
}

func NewNotAllowedCreateWorkspace() *NotAllowedCreateWorkspace {
	return &NotAllowedCreateWorkspace{
		BaseHTTPException: &BaseHTTPException{
			code: "not_allowed_create_workspace",
			ValueError: exceptions.NewValueError(
				"Workspace not found, please contact system admin to invite you to join in a workspace.",
			),
			status: 400,
		},
	}
}

type AccountBannedError struct {
	*BaseHTTPException
}

func NewAccountBannedError() *AccountBannedError {
	return &AccountBannedError{
		BaseHTTPException: &BaseHTTPException{
			code: "account_banned",
			ValueError: exceptions.NewValueError(
				"Account is banned.",
			),
			status: 400,
		},
	}
}

type AccountNotFound struct {
	*BaseHTTPException
}

func NewAccountNotFound() *AccountNotFound {
	return &AccountNotFound{
		BaseHTTPException: &BaseHTTPException{
			code: "account_not_found",
			ValueError: exceptions.NewValueError(
				"Account not found.",
			),
			status: 400,
		},
	}
}

type EmailSendIpLimitError struct {
	*BaseHTTPException
}

func NewEmailSendIpLimitError() *EmailSendIpLimitError {
	return &EmailSendIpLimitError{
		BaseHTTPException: &BaseHTTPException{
			code: "email_send_ip_limit",
			ValueError: exceptions.NewValueError(
				"Too many emails have been sent from this IP address recently. Please try again later.",
			),
			status: 429,
		},
	}
}

type FileTooLargeError struct {
	*BaseHTTPException
}

func NewFileTooLargeError() *FileTooLargeError {
	return &FileTooLargeError{
		BaseHTTPException: &BaseHTTPException{
			code: "file_too_large",
			ValueError: exceptions.NewValueError(
				"File size exceeded. {message}",
			),
			status: 413,
		},
	}
}

type UnsupportedFileTypeError struct {
	*BaseHTTPException
}

func NewUnsupportedFileTypeError() *UnsupportedFileTypeError {
	return &UnsupportedFileTypeError{
		BaseHTTPException: &BaseHTTPException{
			code: "unsupported_file_type",
			ValueError: exceptions.NewValueError(
				"File type not allowed.",
			),
			status: 415,
		},
	}
}

type TooManyFilesError struct {
	*BaseHTTPException
}

func NewTooManyFilesError() *TooManyFilesError {
	return &TooManyFilesError{
		BaseHTTPException: &BaseHTTPException{
			code: "too_many_files",
			ValueError: exceptions.NewValueError(
				"Only one file is allowed.",
			),
			status: 400,
		},
	}
}

type NoFileUploadedError struct {
	*BaseHTTPException
}

func NewNoFileUploadedError() *NoFileUploadedError {
	return &NoFileUploadedError{
		BaseHTTPException: &BaseHTTPException{
			code: "no_file_uploaded",
			ValueError: exceptions.NewValueError(
				"Please upload your file.",
			),
			status: 400,
		},
	}
}

type UnauthorizedAndForceLogout struct {
	*BaseHTTPException
}

func NewUnauthorizedAndForceLogout(desc string) *UnauthorizedAndForceLogout {
	desc = "unauthorized and force logout:" + desc
	return &UnauthorizedAndForceLogout{
		BaseHTTPException: &BaseHTTPException{
			code: "unauthorized_and_force_logout",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 401,
		},
	}
}

type AccountInFreezeError struct {
	*BaseHTTPException
}

func NewAccountInFreezeError() *AccountInFreezeError {
	return &AccountInFreezeError{
		BaseHTTPException: &BaseHTTPException{
			code:   "account_in_freeze",
			status: 400,
			ValueError: exceptions.NewValueError(
				`This email account has been deleted within the past 30 days and is temporarily unavailable for new account registration.`,
			),
		},
	}
}

type RepeatPasswordNotMatchError struct {
	*BaseHTTPException
}

func NewRepeatPasswordNotMatchError() *RepeatPasswordNotMatchError {
	return &RepeatPasswordNotMatchError{
		BaseHTTPException: &BaseHTTPException{
			code: "repeat_password_not_match",
			ValueError: exceptions.NewValueError(
				"New password and repeat password does not match.",
			),
			status: 400,
		},
	}
}

type CurrentPasswordIncorrectError struct {
	*BaseHTTPException
}

func NewCurrentPasswordIncorrectError(desc string) *CurrentPasswordIncorrectError {
	desc = "current password incorrect error:" + desc
	return &CurrentPasswordIncorrectError{
		BaseHTTPException: &BaseHTTPException{
			code: "current_password_incorrect",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 400,
		},
	}
}

type ProviderRequestFailedError struct {
	*BaseHTTPException
}

func NewProviderRequestFailedError() *ProviderRequestFailedError {
	return &ProviderRequestFailedError{
		BaseHTTPException: &BaseHTTPException{
			code: "provider_request_failed",
			ValueError: exceptions.NewValueError(
				"",
			),
			status: 400,
		},
	}
}

type InvalidInvitationCodeError struct {
	*BaseHTTPException
}

func NewInvalidInvitationCodeError() *InvalidInvitationCodeError {
	return &InvalidInvitationCodeError{
		BaseHTTPException: &BaseHTTPException{
			code: "invalid_invitation_code",
			ValueError: exceptions.NewValueError(
				"Invalid invitation code.",
			),
			status: 400,
		},
	}
}

type AccountAlreadyInitedError struct {
	*BaseHTTPException
}

func NewAccountAlreadyInitedError() *AccountAlreadyInitedError {
	return &AccountAlreadyInitedError{
		BaseHTTPException: &BaseHTTPException{
			code: "account_already_inited",
			ValueError: exceptions.NewValueError(
				"The account has been initialized. Please refresh the page.",
			),
			status: 400,
		},
	}
}

type AccountNotInitializedError struct {
	*BaseHTTPException
}

func NewAccountNotInitializedError() *AccountNotInitializedError {
	return &AccountNotInitializedError{
		BaseHTTPException: &BaseHTTPException{
			code: "account_not_initialized",
			ValueError: exceptions.NewValueError(
				"The account has not been initialized yet. Please proceed with the initialization process first.",
			),
			status: 400,
		},
	}
}

type InvalidAccountDeletionCodeError struct {
	*BaseHTTPException
}

func NewInvalidAccountDeletionCodeError() *InvalidAccountDeletionCodeError {
	return &InvalidAccountDeletionCodeError{
		BaseHTTPException: &BaseHTTPException{
			code: "invalid_account_deletion_code",
			ValueError: exceptions.NewValueError(
				"Invalid account deletion code.",
			),
			status: 400,
		},
	}
}

type NotCompletionAppError struct {
	*BaseHTTPException
}

func NewNotCompletionAppError() *NotCompletionAppError {
	return &NotCompletionAppError{
		BaseHTTPException: &BaseHTTPException{
			code: "not_completion_app",
			ValueError: exceptions.NewValueError(
				"Please check if your Completion app mode matches the right API route.",
			),
			status: 400,
		},
	}
}

type NotChatAppError struct {
	*BaseHTTPException
}

func NewNotChatAppError() *NotChatAppError {
	return &NotChatAppError{
		BaseHTTPException: &BaseHTTPException{
			code: "not_chat_app",
			ValueError: exceptions.NewValueError(
				"Please check if your app mode matches the right API route.",
			),
			status: 400,
		},
	}
}

type NotWorkflowAppError struct {
	*BaseHTTPException
}

func NewNotWorkflowAppError() *NotWorkflowAppError {
	return &NotWorkflowAppError{
		BaseHTTPException: &BaseHTTPException{
			code: "not_workflow_app",
			ValueError: exceptions.NewValueError(
				"Please check if your app mode matches the right API route.",
			),
			status: 400,
		},
	}
}

type BadRequest struct {
	*BaseHTTPException
}

func NewBadRequest(desc string) *BadRequest {
	if desc == "" {
		desc = "The browser (or proxy) sent a request that this server could not understand."
	}
	return &BadRequest{
		BaseHTTPException: &BaseHTTPException{
			code: "bad_request",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 400,
		},
	}
}
func NewNotFound(desc string) *BaseHTTPException {
	desc = `The requested URL was not found on the server. If you entered"
 the URL manually please check your spelling and try again.` + desc
	return &BaseHTTPException{
		code: "not_found",
		ValueError: exceptions.NewValueError(
			desc,
		),
		status: 404,
	}
}

func NewAppSuggestedQuestionsAfterAnswerDisabledError() *BaseHTTPException {
	return &BaseHTTPException{
		code: "app_suggested_questions_after_answer_disabled",
		ValueError: exceptions.NewValueError(
			"Function Suggested questions after answer disabled.",
		),
		status: 403,
	}
}

func NewPluginDaemonInnerError(desc string) *BaseHTTPException {
	desc = `plugin daemon inner error:` + desc
	return &BaseHTTPException{
		code: "plugin_daemon_inner_error",
		ValueError: exceptions.NewValueError(
			desc,
		),
		status: 500,
	}
}
