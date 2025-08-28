package mcp

type NotificationParams struct {
	Meta struct {
		ModelConfig map[string]any `json:"model_config"`
	} `json:"meta"` //Field(alias="_meta", default=None)
	/*
	   This parameter name is reserved by MCP to allow clients and servers to attach
	   additional metadata to their notifications.
	*/
}
type NotificationParamsT interface {
	NotificationParams | map[string]any | ProgressNotificationParams | ResourceUpdatedNotificationParams | LoggingMessageNotificationParams | CancelledNotificationParams
}

type Notification[T NotificationParamsT] struct {
	/*Base class for JSON-RPC notifications.*/

	Method      string         `json:"method"`
	Params      T              `json:"params"` /*NotificationParams |  map[string]any*/
	ModelConfig map[string]any `json:"model_config"`
}

type JSONRPCNotification struct {
	*Notification[map[string]any]
	/*A notification which does not expect a response.*/

	Jsonrpc string         `json:"jsonrpc"` //"2.0"
	Params  map[string]any `json:"params"`
}

type InitializedNotification struct {
	*Notification[NotificationParams]
	/*
	   This notification is sent from the client to the server after initialization has
	   finished.
	*/

	Method string              `json:"method"` //["notifications/initialized"]
	Params *NotificationParams `json:"params"`
}

type ProgressNotificationParams struct {
	*NotificationParams
	/*Parameters for progress notifications.*/

	ProgressToken any `json:"progressToken"` // ProgressToken | None = None
	/*
	   The progress token which was given in the initial request, used to associate this
	   notification with the request that is proceeding.
	*/
	Progress float64 `json:"progress"`
	/*
	   The progress thus far. This should increase every time progress is made, even if the
	   total is unknown.
	*/
	Total float64 `json:"total"`
	/*Total number of items to process (or total progress required), if known.*/
	ModelConfig map[string]any `json:"model_config"`
}

type ProgressNotification struct {
	*Notification[ProgressNotificationParams]
	/*
	   An out-of-band notification used to inform the receiver of a progress update for a
	   long-running request.
	*/

	Method string                     `json:"method"` //["notifications/progress"]
	Params ProgressNotificationParams `json:"params"`
}
type ResourceListChangedNotification struct {
	*Notification[NotificationParams]
	/*
	   An optional notification from the server to the client, informing it that the list
	   of resources it can read from has changed.
	*/

	Method string              `json:"method"` //["notifications/resources/list_changed"]
	Params *NotificationParams `json:"params"`
}
type ResourceUpdatedNotificationParams struct {
	*NotificationParams
	/*Parameters for resource update notifications.*/

	URI string `json:"uri"` //[AnyUrl, UrlConstraints(host_required=False)]
	/*
	   The URI of the resource that has been updated. This might be a sub-resource of the
	   one that the client actually subscribed to.
	*/
	ModelConfig map[string]any `json:"model_config"`
}

type ResourceUpdatedNotification struct {
	Notification[ResourceUpdatedNotificationParams]
	/*
	   A notification from the server to the client, informing it that a resource has
	   changed and may need to be read again.
	*/

	Method string                            `json:"method"` //["notifications/resources/updated"]
	Params ResourceUpdatedNotificationParams `json:"params"`
}
type PromptListChangedNotification struct {
	Notification[NotificationParams]

	/*
	   An optional notification from the server to the client, informing it that the list
	   of prompts it offers has changed.
	*/

	Method string              `json:"method"` //["notifications/prompts/list_changed"]
	Params *NotificationParams `json:"params"`
}
type ToolListChangedNotification struct {
	Notification[NotificationParams]
	/*
	   An optional notification from the server to the client, informing it that the list
	   of tools it offers has changed.
	*/

	Method string              `json:"method"` //["notifications/tools/list_changed"]
	Params *NotificationParams `json:"params"`
}
type LoggingMessageNotificationParams struct {
	*NotificationParams
	/*Parameters for logging message notifications.*/

	Level LoggingLevelType `json:"level"`
	/*The severity of this log message.*/
	Logger string `json:"logger"`
	/*An optional name of the logger issuing this message.*/
	Data any `json:"data"`
	/*
	   The data to be logged, such as a string message or an object. Any JSON serializable
	   type is allowed here.
	*/
	ModelConfig map[string]any `json:"model_config"`
}

type LoggingMessageNotification struct {
	Notification[LoggingMessageNotificationParams]
	/*Notification of a log message passed from server to client.*/

	Method string                           `json:"method"` //["notifications/message"]
	Params LoggingMessageNotificationParams `json:"params"`
}

type RootsListChangedNotification struct {
	Notification[NotificationParams]
	/*
	   A notification from the client to the server, informing it that the list of
	   roots has changed.

	   This notification should be sent whenever the client adds, removes, or
	   modifies any root. The server should then request an updated list of roots
	   using the ListRootsRequest.
	*/

	Method string              `json:"method"` //["notifications/roots/list_changed"]
	Params *NotificationParams `json:"params"`
}
type CancelledNotificationParams struct {
	*NotificationParams
	/*Parameters for cancellation notifications.    */

	RequestId any/* int | string */ `json:"requestId"`
	/*The ID of the request to cancel.    */
	Reason string `json:"reason"`
	/*An optional string describing the reason for the cancellation.    */
	ModelConfig map[string]any `json:"model_config"`
}

type CancelledNotification struct {
	Notification[CancelledNotificationParams]
	/*
	   This notification can be sent by either side to indicate that it is canceling a
	   previously-issued request.
	*/

	Method string                      `json:"method"` //["notifications/cancelled"]
	Params CancelledNotificationParams `json:"params"`
}
