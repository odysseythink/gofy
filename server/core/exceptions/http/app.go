package http

import "mlib.com/gofy/server/core/exceptions"

type AppNotFoundError struct {
	*BaseHTTPException
}

func NewAppNotFoundError(desc string) *AppNotFoundError {
	if desc == "" {
		desc = "App not found."
	}
	return &AppNotFoundError{
		BaseHTTPException: &BaseHTTPException{
			code: "app_not_found",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 404,
		},
	}
}

type ProviderNotInitializeError struct {
	*BaseHTTPException
}

func NewProviderNotInitializeError(desc string) *ProviderNotInitializeError {
	if desc == "" {
		desc = `No valid model provider credentials found. Please go to Settings -> Model Provider to complete your provider credentials.`
	} else {
		desc = `No valid model provider credentials found:` + desc
	}
	return &ProviderNotInitializeError{
		BaseHTTPException: &BaseHTTPException{
			code: "provider_not_initialize",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 400,
		},
	}
}

type ProviderQuotaExceededError struct {
	*BaseHTTPException
}

func NewProviderQuotaExceededError() *ProviderQuotaExceededError {
	return &ProviderQuotaExceededError{
		BaseHTTPException: &BaseHTTPException{
			code: "provider_quota_exceeded",
			ValueError: exceptions.NewValueError(
				"Your quota for Dify Hosted Model Provider has been exhausted. Please go to Settings -> Model Provider to complete your own provider credentials.",
			),
			status: 400,
		},
	}
}

type ProviderModelCurrentlyNotSupportError struct {
	*BaseHTTPException
}

func NewProviderModelCurrentlyNotSupportError() *ProviderModelCurrentlyNotSupportError {
	return &ProviderModelCurrentlyNotSupportError{
		BaseHTTPException: &BaseHTTPException{
			code: "model_currently_not_support",
			ValueError: exceptions.NewValueError(
				"Dify Hosted OpenAI trial currently not support the GPT-4 model.",
			),
			status: 400,
		},
	}
}

type ConversationCompletedError struct {
	*BaseHTTPException
}

func NewConversationCompletedError(desc string) *ConversationCompletedError {
	if desc == "" {
		desc = "The conversation has ended. Please start a new conversation."
	} else {
		desc = "conversation completed error:" + desc
	}
	return &ConversationCompletedError{
		BaseHTTPException: &BaseHTTPException{
			code: "conversation_completed",
			ValueError: exceptions.NewValueError(
				"The conversation has ended. Please start a new conversation.",
			),
			status: 400,
		},
	}
}

type AppUnavailableError struct {
	*BaseHTTPException
}

func NewAppUnavailableError() *AppUnavailableError {
	return &AppUnavailableError{
		BaseHTTPException: &BaseHTTPException{
			code: "app_unavailable",
			ValueError: exceptions.NewValueError(
				"App unavailable, please check your app configurations.",
			),
			status: 400,
		},
	}
}

type CompletionRequestError struct {
	*BaseHTTPException
}

func NewCompletionRequestError(desc string) *CompletionRequestError {
	return &CompletionRequestError{
		BaseHTTPException: &BaseHTTPException{
			code: "completion_request_error",
			ValueError: exceptions.NewValueError(
				"Completion request failed:" + desc,
			),
			status: 400,
		},
	}
}

type AppMoreLikeThisDisabledError struct {
	*BaseHTTPException
}

func NewAppMoreLikeThisDisabledError() *AppMoreLikeThisDisabledError {
	return &AppMoreLikeThisDisabledError{
		BaseHTTPException: &BaseHTTPException{
			code: "app_more_like_this_disabled",
			ValueError: exceptions.NewValueError(
				"The 'More like this' feature is disabled. Please refresh your page.",
			),
			status: 403,
		},
	}
}

type NoAudioUploadedError struct {
	*BaseHTTPException
}

func NewNoAudioUploadedError() *NoAudioUploadedError {
	return &NoAudioUploadedError{
		BaseHTTPException: &BaseHTTPException{
			code: "no_audio_uploaded",
			ValueError: exceptions.NewValueError(
				"Please upload your audio.",
			),
			status: 400,
		},
	}
}

type AudioTooLargeError struct {
	*BaseHTTPException
}

func NewAudioTooLargeError() *AudioTooLargeError {
	return &AudioTooLargeError{
		BaseHTTPException: &BaseHTTPException{
			code: "audio_too_large",
			ValueError: exceptions.NewValueError(
				"Audio size exceeded. {message}",
			),
			status: 413,
		},
	}
}

type UnsupportedAudioTypeError struct {
	*BaseHTTPException
}

func NewUnsupportedAudioTypeError() *UnsupportedAudioTypeError {
	return &UnsupportedAudioTypeError{
		BaseHTTPException: &BaseHTTPException{
			code: "unsupported_audio_type",
			ValueError: exceptions.NewValueError(
				"Audio type not allowed.",
			),
			status: 415,
		},
	}
}

type ProviderNotSupportSpeechToTextError struct {
	*BaseHTTPException
}

func NewProviderNotSupportSpeechToTextError() *ProviderNotSupportSpeechToTextError {
	return &ProviderNotSupportSpeechToTextError{
		BaseHTTPException: &BaseHTTPException{
			code: "provider_not_support_speech_to_text",
			ValueError: exceptions.NewValueError(
				"Provider not support speech to text.",
			),
			status: 400,
		},
	}
}

type DraftWorkflowNotExist struct {
	*BaseHTTPException
}

func NewDraftWorkflowNotExist() *DraftWorkflowNotExist {
	return &DraftWorkflowNotExist{
		BaseHTTPException: &BaseHTTPException{
			code: "draft_workflow_not_exist",
			ValueError: exceptions.NewValueError(
				"Draft workflow need to be initialized.",
			),
			status: 400,
		},
	}
}

type DraftWorkflowNotSync struct {
	*BaseHTTPException
}

func NewDraftWorkflowNotSync() *DraftWorkflowNotSync {
	return &DraftWorkflowNotSync{
		BaseHTTPException: &BaseHTTPException{
			code: "draft_workflow_not_sync",
			ValueError: exceptions.NewValueError(
				"Workflow graph might have been modified, please refresh and resubmit.",
			),
			status: 400,
		},
	}
}

type TracingConfigNotExist struct {
	*BaseHTTPException
}

func NewTracingConfigNotExist() *TracingConfigNotExist {
	return &TracingConfigNotExist{
		BaseHTTPException: &BaseHTTPException{
			code: "trace_config_not_exist",
			ValueError: exceptions.NewValueError(
				"Trace config not exist.",
			),
			status: 400,
		},
	}
}

type TracingConfigIsExist struct {
	*BaseHTTPException
}

func NewTracingConfigIsExist() *TracingConfigIsExist {
	return &TracingConfigIsExist{
		BaseHTTPException: &BaseHTTPException{
			code: "trace_config_is_exist",
			ValueError: exceptions.NewValueError(
				"Trace config is exist.",
			),
			status: 400,
		},
	}
}

type TracingConfigCheckError struct {
	*BaseHTTPException
}

func NewTracingConfigCheckError() *TracingConfigCheckError {
	return &TracingConfigCheckError{
		BaseHTTPException: &BaseHTTPException{
			code: "trace_config_check_error",
			ValueError: exceptions.NewValueError(
				"Invalid Credentials.",
			),
			status: 400,
		},
	}
}

type InternalServerError struct {
	*BaseHTTPException
}

func NewInternalServerError(desc string) *InternalServerError {
	if desc == "" {
		desc = "The server encountered an internal error and was unable to complete your request. Either the server is overloaded or there is an error in the application."
	}
	return &InternalServerError{
		BaseHTTPException: &BaseHTTPException{
			code: "internal_server_error",
			ValueError: exceptions.NewValueError(
				desc,
			),
			status: 500,
		},
	}
}
