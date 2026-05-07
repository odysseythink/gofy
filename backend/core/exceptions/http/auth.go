package http

import "github.com/odysseythink/gofy/backend/core/exceptions"

type ApiKeyAuthFailedError struct {
	*BaseHTTPException
}

func NewApiKeyAuthFailedError() *ApiKeyAuthFailedError {
	return &ApiKeyAuthFailedError{
		BaseHTTPException: &BaseHTTPException{
			code: "auth_failed",
			ValueError: exceptions.NewValueError(
				"api key auth failed",
			),
			status: 500,
		},
	}
}

type InvalidEmailError struct {
	*BaseHTTPException
}

func NewInvalidEmailError() *InvalidEmailError {
	return &InvalidEmailError{
		BaseHTTPException: &BaseHTTPException{
			code: "invalid_email",
			ValueError: exceptions.NewValueError(
				"The email address is not valid.",
			),
			status: 400,
		},
	}
}

type PasswordMismatchError struct {
	*BaseHTTPException
}

func NewPasswordMismatchError() *PasswordMismatchError {
	return &PasswordMismatchError{
		BaseHTTPException: &BaseHTTPException{
			code: "password_mismatch",
			ValueError: exceptions.NewValueError(
				"The passwords do not match.",
			),
			status: 400,
		},
	}
}

type InvalidTokenError struct {
	*BaseHTTPException
}

func NewInvalidTokenError() *InvalidTokenError {
	return &InvalidTokenError{
		BaseHTTPException: &BaseHTTPException{
			code: "invalid_or_expired_token",
			ValueError: exceptions.NewValueError(
				"The token is invalid or has expired.",
			),
			status: 400,
		},
	}
}

type PasswordResetRateLimitExceededError struct {
	*BaseHTTPException
}

func NewPasswordResetRateLimitExceededError() *PasswordResetRateLimitExceededError {
	return &PasswordResetRateLimitExceededError{
		BaseHTTPException: &BaseHTTPException{
			code: "password_reset_rate_limit_exceeded",
			ValueError: exceptions.NewValueError(
				"Too many password reset emails have been sent. Please try again in 1 minutes.",
			),
			status: 429,
		},
	}
}

type EmailCodeError struct {
	*BaseHTTPException
}

func NewEmailCodeError() *EmailCodeError {
	return &EmailCodeError{
		BaseHTTPException: &BaseHTTPException{
			code: "email_code_error",
			ValueError: exceptions.NewValueError(
				"Email code is invalid or expired.",
			),
			status: 400,
		},
	}
}

type EmailOrPasswordMismatchError struct {
	*BaseHTTPException
}

func NewEmailOrPasswordMismatchError() *EmailOrPasswordMismatchError {
	return &EmailOrPasswordMismatchError{
		BaseHTTPException: &BaseHTTPException{
			code: "email_or_password_mismatch",
			ValueError: exceptions.NewValueError(
				"The email or password is mismatched.",
			),
			status: 400,
		},
	}
}

type EmailPasswordLoginLimitError struct {
	*BaseHTTPException
}

func NewEmailPasswordLoginLimitError() *EmailPasswordLoginLimitError {
	return &EmailPasswordLoginLimitError{
		BaseHTTPException: &BaseHTTPException{
			code: "email_code_login_limit",
			ValueError: exceptions.NewValueError(
				"Too many incorrect password attempts. Please try again later.",
			),
			status: 429,
		},
	}
}

type EmailCodeLoginRateLimitExceededError struct {
	*BaseHTTPException
}

func NewEmailCodeLoginRateLimitExceededError() *EmailCodeLoginRateLimitExceededError {
	return &EmailCodeLoginRateLimitExceededError{
		BaseHTTPException: &BaseHTTPException{
			code: "email_code_login_rate_limit_exceeded",
			ValueError: exceptions.NewValueError(
				"Too many login emails have been sent. Please try again in 5 minutes.",
			),
			status: 429,
		},
	}
}

type EmailCodeAccountDeletionRateLimitExceededError struct {
	*BaseHTTPException
}

func NewEmailCodeAccountDeletionRateLimitExceededError() *EmailCodeAccountDeletionRateLimitExceededError {
	return &EmailCodeAccountDeletionRateLimitExceededError{
		BaseHTTPException: &BaseHTTPException{
			code: "email_code_account_deletion_rate_limit_exceeded",
			ValueError: exceptions.NewValueError(
				"Too many account deletion emails have been sent. Please try again in 5 minutes.",
			),
			status: 429,
		},
	}
}
