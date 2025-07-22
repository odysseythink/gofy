package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pbexceptions "mlib.com/gofy/server/proto/exceptions"
)

func Unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "Authorization Required")
	description := `The server could not verify that you are authorized to access
	     the URL requested. You either supplied the wrong credentials
	     (e.g. a bad password), or your browser doesn't understand
	     how to supply the credentials required.`
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": description})
}

func Forbidden(c *gin.Context) {
	// """*403* `Forbidden`

	// Raise if the user doesn't have the permission for the requested resource
	// but was authenticated.
	// """

	// code = 403
	description := `You don't have the permission to access the requested
         resource. It is either read-protected or not readable by the
         server.`
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": description})
}

func InvalidArgError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "invalid_arg",
		Status:  400,
		Data:    nil,
		Message: "invalid arg.",
	})
}

func InvalidArgErrorWithDetail(c *gin.Context, desc string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "invalid_arg",
		Status:  http.StatusBadRequest,
		Data:    nil,
		Message: desc,
	})
}

func PbHttpException(c *gin.Context, exp *pbexceptions.HTTPException) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    exp.Code,
		Status:  int(exp.Status),
		Data:    nil,
		Message: exp.Message,
	})
}

func EmailPasswordLoginLimitError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "email_code_login_limit",
		Status:  429,
		Data:    nil,
		Message: "Too many incorrect password attempts. Please try again later.",
	})
}

func InvalidEmailError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "invalid_email",
		Status:  400,
		Data:    nil,
		Message: "The email address is not valid.",
	})
}

func EmailOrPasswordMismatchError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "email_or_password_mismatch",
		Status:  400,
		Data:    nil,
		Message: "The email or password is mismatched.",
	})
}

func AppNotFoundError(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code:    "app_not_found",
		Status:  404,
		Data:    nil,
		Message: "App not found.",
	})
}

func AccountNotFoundError(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Code:    "account_not_found",
		Status:  404,
		Data:    nil,
		Message: "Account not found.",
	})
}

func ProviderNotInitializeError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "provider_not_initialize",
		Status:  400,
		Data:    nil,
		Message: "No valid model provider credentials found. \nPlease go to Settings -> Model Provider to complete your provider credentials.",
	})
}

func ProviderQuotaExceededError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "provider_quota_exceeded",
		Status:  400,
		Data:    nil,
		Message: "Your quota for Dify Hosted Model Provider has been exhausted. \nPlease go to Settings -> Model Provider to complete your own provider credentials.",
	})
}

func ProviderModelCurrentlyNotSupportError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "model_currently_not_support",
		Status:  400,
		Data:    nil,
		Message: "Dify Hosted OpenAI trial currently not support the GPT-4 model.",
	})
}

func ConversationCompletedError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "conversation_completed",
		Status:  400,
		Data:    nil,
		Message: "The conversation has ended. Please start a new conversation.",
	})
}

func AppUnavailableError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "app_unavailable",
		Status:  400,
		Data:    nil,
		Message: "App unavailable, please check your app configurations.",
	})
}

func CompletionRequestError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    "completion_request_error",
		Status:  400,
		Data:    nil,
		Message: "Completion request failed.",
	})
}

func AppMoreLikeThisDisabledError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    "app_more_like_this_disabled",
		Status:  403,
		Data:    nil,
		Message: "The 'More like this' feature is disabled. Please refresh your page.",
	})
}

func NoAudioUploadedError(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "no_audio_uploaded",
		Status:  400,
		Data:    nil,
		Message: "Please upload your audio.",
	})
}

func AudioTooLargeError(c *gin.Context) {
	c.JSON(413, Response{
		Code:    "audio_too_large",
		Status:  413,
		Data:    nil,
		Message: "Audio size exceeded. {message}",
	})
}

func UnsupportedAudioTypeError(c *gin.Context) {
	c.JSON(415, Response{
		Code:    "unsupported_audio_type",
		Status:  415,
		Data:    nil,
		Message: "Audio type not allowed.",
	})
}

func ProviderNotSupportSpeechToTextError(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "provider_not_support_speech_to_text",
		Status:  400,
		Data:    nil,
		Message: "Provider not support speech to text.",
	})
}

func NoFileUploadedError(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "no_file_uploaded",
		Status:  400,
		Data:    nil,
		Message: "Please upload your file.",
	})
}

func TooManyFilesError(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "too_many_files",
		Status:  400,
		Data:    nil,
		Message: "Only one file is allowed.",
	})
}

func DraftWorkflowNotExist(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "draft_workflow_not_exist",
		Status:  400,
		Data:    nil,
		Message: "Draft workflow need to be initialized.",
	})
}

// func DraftWorkflowNotSync(c *gin.Context) {
// 	c.JSON(400, Response{
// 		Code:    "draft_workflow_not_sync",
// 		Status:  400,
// 		Data:    nil,
// 		Message: "Workflow graph might have been modified, please refresh and resubmit.",
// 	})
// }

func TracingConfigNotExist(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "trace_config_not_exist",
		Status:  400,
		Data:    nil,
		Message: "Trace config not exist.",
	})
}

func TracingConfigIsExist(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "trace_config_is_exist",
		Status:  400,
		Data:    nil,
		Message: "Trace config is exist.",
	})
}

func TracingConfigCheckError(c *gin.Context) {
	c.JSON(400, Response{
		Code:    "trace_config_check_error",
		Status:  400,
		Data:    nil,
		Message: "Invalid Credentials.",
	})
}

func InvokeRateLimitError(c *gin.Context) {
	c.JSON(429, Response{
		Code:    "rate_limit_error",
		Status:  429,
		Data:    nil,
		Message: "Rate Limit Error",
	})
}
