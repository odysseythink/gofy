package mcp

const (
	// Client support both version, not support 2025-06-18 yet.
	LATEST_PROTOCOL_VERSION = "2025-03-26"
	// Server support 2024-11-05 to allow claude to use.
	SERVER_LATEST_PROTOCOL_VERSION = "2024-11-05"
)

type ProgressToken interface {
	string | int
}
type Cursor string
type RoleType string

const (
	Role_USER      RoleType = "user"
	Role_ASSISTANT RoleType = "assistant"
)
type RequestIdT interface{
	int | string
}

type RequestParams  struct {
    Meta struct{
        ProgressToken any `json:"progressToken"`// ProgressToken | None = None
        /*
        If specified, the caller requests out-of-band progress notifications for
        this request (as represented by notifications/progress). The value of this
        parameter is an opaque token that will be attached to any subsequent
        notifications. The receiver is not obligated to provide these notifications.
        */

        ModelConfig map[string]any `json:"model_config"`
	}  `json:"meta"` //Field(alias="_meta", default=None)
	}

type NotificationParams  struct {
    Meta struct {
        ModelConfig map[string]any `json:"model_config"`
	 } `json:"meta"` //Field(alias="_meta", default=None)
    /*
    This parameter name is reserved by MCP to allow clients and servers to attach
    additional metadata to their notifications.
    */
	}

type RequestParamsT interface {
	RequestParams |  map[string]any
}
type NotificationParamsT  interface {
	NotificationParams |  map[string]any 
}
type MethodT string

type Request[T RequestParamsT] struct{
    /*Base class for JSON-RPC requests.*/

    Method string`json:"method"`
    Params T`json:"params"` /*RequestParams |  map[string]any*/
    ModelConfig map[string]any `json:"model_config"`
}

type PaginatedRequest[T RequestParamsT] struct{
	*Request[T] 
     Cursor `json:"cursor"`
    /*
    An opaque token representing the current pagination position.
    If provided, the server should return results starting after this cursor.
    */
}

type Notification[T NotificationParamsT] struct{
    /*Base class for JSON-RPC notifications.*/


    Method string`json:"method"`
    Params T`json:"params"` /*NotificationParams |  map[string]any*/
    ModelConfig map[string]any `json:"model_config"`
}

type Result  struct {
    /*Base class for JSON-RPC results.*/

    ModelConfig map[string]any `json:"model_config"`

    Meta map[string]any `json:"meta"`//Field(alias="_meta", default=None)
    /*
    This result property is reserved by the protocol to allow clients and servers to
    attach additional metadata to their responses.
    */
}

type PaginatedResult struct {
	*Result
    NextCursor Cursor `json:"nextCursor"`
    /*
    An opaque token representing the pagination position after the last returned result.
    If present, there may be more results available.
    */
}

type JSONRPCRequest[T RequestIdT] struct {
	*Request[map[string]any]
    /*A request that expects a response.*/

    Jsonrpc string `json:"jsonrpc"`//"2.0"
    ID T `json:"id"`
    Method string `json:"method"`
    Params map[string]any `json:"params"`
}

type JSONRPCNotification struct {
	*Notification[map[string]any]
    /*A notification which does not expect a response.*/

    Jsonrpc string `json:"jsonrpc"`//"2.0"
    Params map[string]any `json:"params"`
}

type JSONRPCResponse[T RequestIdT]  struct {
    /*A successful (non-error) response to a request.*/

    Jsonrpc string `json:"jsonrpc"`//"2.0"
    ID T `json:"id"`
    Result map[string]any`json:"result"`
    ModelConfig map[string]any `json:"model_config"`
}

const (
// Standard JSON-RPC error codes
PARSE_ERROR = -32700
INVALID_REQUEST = -32600
METHOD_NOT_FOUND = -32601
INVALID_PARAMS = -32602
INTERNAL_ERROR = -32603
)

type ErrorData  struct {
    /*Error information for JSON-RPC error responses.*/

    Code int `json:"code"`
    /*The error type that occurred.*/

    Message string `json:"message"`
    /*
    A short description of the error. The message SHOULD be limited to a concise single
    sentence.
    */

    Data any  `json:"data"`
    /*
    Additional information about the error. The value of this member is defined by the
    sender (e.g. detailed error information, nested errors etc.).
    */

    ModelConfig map[string]any `json:"model_config"`
}

type JSONRPCError[T string|int]  struct {
    /*A response to a request that indicates an error occurred.*/

    Jsonrpc string `json:"jsonrpc"`//"2.0"
    ID T`json:"id"`
    Error ErrorData`json:"error"`
    ModelConfig map[string]any `json:"model_config"`
}

// type JSONRPCMessage(RootModel[JSONRPCRequest | JSONRPCNotification | JSONRPCResponse | JSONRPCError]):
//     pass


type EmptyResult struct{
	*Result
    /*A response that indicates success but carries no data.*/
}

type Implementation  struct {
    /*Describes the name and version of an MCP implementation.*/

    Name string `json:"name"`
    Version string `json:"version"`
    ModelConfig map[string]any `json:"model_config"`
}

type RootsCapability  struct {
    /*Capability for root operations.*/

    ListChanged bool `json:"listChanged"`
    /*Whether the client supports notifications for changes to the roots list.*/
    ModelConfig map[string]any `json:"model_config"`
}

type SamplingCapability  struct {
    /*Capability for logging operations.*/

    ModelConfig map[string]any `json:"model_config"`
}

type ClientCapabilities  struct {
    /*Capabilities a client may support.*/
    Experimental map[string]map[string]any `json:"experimental"`
    /*Experimental, non-standard capabilities that the client supports.*/
    Sampling *SamplingCapability  `json:"sampling"`
    /*Present if the client supports sampling from an LLM.*/
    Roots *RootsCapability  `json:"roots"`
    /*Present if the client supports listing roots.*/
    ModelConfig map[string]any `json:"model_config"`
}

type PromptsCapability  struct {
    /*Capability for prompts operations.*/

    ListChanged bool `json:"listChanged"`
    /*Whether this server supports notifications for changes to the prompt list.*/
    ModelConfig map[string]any `json:"model_config"`
}

type ResourcesCapability  struct {
    /*Capability for resources operations.*/

    Subscribe bool `json:"subscribe"`
    /*Whether this server supports subscribing to resource updates.*/
    ListChanged bool `json:"listChanged"`
    /*Whether this server supports notifications for changes to the resource list.*/
    ModelConfig map[string]any `json:"model_config"`
}

type ToolsCapability  struct {
    /*Capability for tools operations.*/

    ListChanged bool `json:"listChanged"`
    /*Whether this server supports notifications for changes to the tool list.*/
    ModelConfig map[string]any `json:"model_config"`
}

type LoggingCapability  struct {
    /*Capability for logging operations.*/

    ModelConfig map[string]any `json:"model_config"`
}

type ServerCapabilities  struct {
    /*Capabilities that a server may support.*/

    Experimental map[string]map[string]any `json:"experimental"`
    /*Experimental, non-standard capabilities that the server supports.*/
    Logging *LoggingCapability `json:"logging"`
    /*Present if the server supports sending log messages to the client.*/
    Prompts *PromptsCapability `json:"prompts"`
    /*Present if the server offers any prompt templates.*/
    Resources *ResourcesCapability  `json:"resources"`
    /*Present if the server offers any resources to read.*/
    Tools *ToolsCapability  `json:"tools"`
    /*Present if the server offers any tools to call.*/
    ModelConfig map[string]any `json:"model_config"`
}

type InitializeRequestParams[T string | int] struct{
	*RequestParams
    /*Parameters for the initialize request.*/

    ProtocolVersion T`json:"protocolVersion"`
    /*The latest version of the Model Context Protocol that the client supports.*/
    Capabilities ClientCapabilities`json:"capabilities"`
    ClientInfo Implementation`json:"clientInfo"`
    ModelConfig map[string]any `json:"model_config"`
}

type InitializeRequest struct{
	*Request[InitializeRequestParams]
    /*
    This request is sent from the client to the server when it first connects, asking it
    to begin initialization.
    */

    Method string `json:"method"`//["initialize"]
    Params InitializeRequestParams`json:"params"`
}

type InitializeResult[T string | int] struct{
 *Result
    /*After receiving an initialize request from the client, the server sends this.*/

    ProtocolVersion T`json:"protocolVersion"`
    /*The version of the Model Context Protocol that the server wants to use.*/
    Capabilities ServerCapabilities`json:"capabilities"`
    ServerInfo Implementation`json:"serverInfo"`
    Instructions string `json:"instructions"`
    /*Instructions describing how to use the server and its features.*/
}

type InitializedNotification struct {
	*Notification[NotificationParams]
    /*
    This notification is sent from the client to the server after initialization has
    finished.
    */

    Method string `json:"method"`//["notifications/initialized"]
    Params *NotificationParams `json:"params"`
}

type PingRequest struct{
 *Request[RequestParams]
    /*
    A ping, issued by either the server or the client, to check that the other party is
    still alive.
    */

    Method string `json:"method"`//["ping"]
    Params *RequestParams `json:"params"`
}


type ProgressNotificationParams struct{
*NotificationParams
    /*Parameters for progress notifications.*/

    ProgressToken any `json:"progressToken"`// ProgressToken | None = None
    /*
    The progress token which was given in the initial request, used to associate this
    notification with the request that is proceeding.
    */
    Progress float64 `json:"progress"`
    /*
    The progress thus far. This should increase every time progress is made, even if the
    total is unknown.
    */
    Total  float64 `json:"total"`
    /*Total number of items to process (or total progress required), if known.*/
    ModelConfig map[string]any `json:"model_config"`
}

type ProgressNotification struct{
	*Notification[ProgressNotificationParams]
    /*
    An out-of-band notification used to inform the receiver of a progress update for a
    long-running request.
    */

    Method string `json:"method"`//["notifications/progress"]
Params ProgressNotificationParams `json:"params"`
}

type ListResourcesRequest struct{
	*PaginatedRequest[RequestParams]
    /*Sent from the client to request a list of resources the server has.*/

    Method string `json:"method"`//["resources/list"]
    Params *RequestParams `json:"params"`
}


type Annotations  struct {
    Audience []RoleType`json:"audience"`
    Priority float64 `json:"priority"`//Annotated[float, Field(ge=0.0, le=1.0)] | None = None
    ModelConfig map[string]any `json:"model_config"`
}

type Resource  struct {
    /*A known resource that the server is capable of reading.*/

    URI string `json:"uri"`//[AnyUrl, UrlConstraints(host_required=False)]
    /*The URI of this resource.*/
    Name string `json:"name"`
    /*A human-readable name for this resource.*/
    Description string `json:"description"`
    /*A description of what this resource represents.*/
    MimeType string `json:"mimeType"`
    /*The MIME type of this resource, if known.*/
    Size int `json:"size"`
    /*
    The size of the raw resource content, in bytes (i.e., before base64 encoding
    or any tokenization), if known.

    This can be used by Hosts to display file sizes and estimate context window usage.
    */
    Annotations *Annotations  `json:"annotations"`
    ModelConfig map[string]any `json:"model_config"`
}

type ResourceTemplate  struct {
    /*A template description for resources available on the server.*/

    UriTemplate string `json:"uriTemplate"`
    /*
    A URI template (according to RFC 6570) that can be used to construct resource
    URIs.
    */
    Name string `json:"name"`
    /*A human-readable name for the type of resource this template refers to.*/
    Description string `json:"description"`
    /*A human-readable description of what this template is for.*/
    MimeType string `json:"mimeType"`
    /*
    The MIME type for all resources that match this template. This should only be
    included if all resources matching this template have the same type.
    */
    Annotations *Annotations  `json:"annotations"`
    ModelConfig map[string]any `json:"model_config"`
}

type ListResourcesResult struct{
	*PaginatedResult
    /*The server's response to a resources/list request from the client.*/

    Resources []*Resource`json:"resources"`
}

type ListResourceTemplatesRequest struct{
*PaginatedRequest[RequestParams]
    /*Sent from the client to request a list of resource templates the server has.*/

    Method string `json:"method"`//["resources/templates/list"]
    Params *RequestParams `json:"params"`
}


type ListResourceTemplatesResult struct{
	*PaginatedResult
    /*The server's response to a resources/templates/list request from the client.*/

    ResourceTemplates []*ResourceTemplate `json:"resourceTemplates"`
}

type ReadResourceRequestParams struct{
	*RequestParams
    /*Parameters for reading a resource.*/

    URI string `json:"uri"`//[AnyUrl, UrlConstraints(host_required=False)]
    /*
    The URI of the resource to read. The URI can use any protocol; it is up to the
    server how to interpret it.
    */
    ModelConfig map[string]any `json:"model_config"`
}

type ReadResourceRequest struct{
 *Request[ReadResourceRequestParams]
    /*Sent from the client to the server, to read a specific resource URI.*/

    Method string `json:"method"`//["resources/read"]
    Params ReadResourceRequestParams`json:"params"`
 }

type ResourceContents  struct {
    /*The contents of a specific resource or sub-resource.*/

    URI string `json:"uri"`//[AnyUrl, UrlConstraints(host_required=False)]
    /*The URI of this resource.*/
    MimeType string `json:"mimeType"`
    /*The MIME type of this resource, if known.*/
    ModelConfig map[string]any `json:"model_config"`
}

type TextResourceContents struct{
*ResourceContents
    /*Text contents of a resource.*/

    Text string `json:"text"`
    /*
    The text of the item. This must only be set if the item can actually be represented
    as text (not binary data).
    */
}

type BlobResourceContents struct{
*ResourceContents
    /*Binary contents of a resource.*/

    Blob string `json:"blob"`
    /*A base64-encoded string representing the binary data of the item.*/
}

type ReadResourceResult struct{
 *Result
    /*The server's response to a resources/read request from the client.*/

    contents: list[TextResourceContents | BlobResourceContents]
}

type ResourceListChangedNotification(
    Notification[NotificationParams | None, Literal["notifications/resources/list_changed"]]
){
    /*
    An optional notification from the server to the client, informing it that the list
    of resources it can read from has changed.
    */

    Method string `json:"method"`//["notifications/resources/list_changed"]
    Params *NotificationParams `json:"params"`


type SubscribeRequestParams struct{
	*RequestParams
    /*Parameters for subscribing to a resource.*/

    URI string `json:"uri"`//[AnyUrl, UrlConstraints(host_required=False)]
    /*
    The URI of the resource to subscribe to. The URI can use any protocol; it is up to
    the server how to interpret it.
    */
    ModelConfig map[string]any `json:"model_config"`


type SubscribeRequest struct{
 *Request[SubscribeRequestParams, Literal["resources/subscribe"]]){
    /*
    Sent from the client to request resources/updated notifications from the server
    whenever a particular resource changes.
    */

    Method string `json:"method"`//["resources/subscribe"]
    params: SubscribeRequestParams


type UnsubscribeRequestParams struct{
	*RequestParams
    /*Parameters for unsubscribing from a resource.*/

    URI string `json:"uri"`//[AnyUrl, UrlConstraints(host_required=False)]
    /*The URI of the resource to unsubscribe from.*/
    ModelConfig map[string]any `json:"model_config"`


type UnsubscribeRequest struct{
 *Request[UnsubscribeRequestParams, Literal["resources/unsubscribe"]]){
    /*
    Sent from the client to request cancellation of resources/updated notifications from
    the server.
    */

    Method string `json:"method"`//["resources/unsubscribe"]
    params: UnsubscribeRequestParams


type ResourceUpdatedNotificationParams struct{
*NotificationParams
    /*Parameters for resource update notifications.*/

    URI string `json:"uri"`//[AnyUrl, UrlConstraints(host_required=False)]
    /*
    The URI of the resource that has been updated. This might be a sub-resource of the
    one that the client actually subscribed to.
    */
    ModelConfig map[string]any `json:"model_config"`


type ResourceUpdatedNotification(
    Notification[ResourceUpdatedNotificationParams, Literal["notifications/resources/updated"]]
){
    /*
    A notification from the server to the client, informing it that a resource has
    changed and may need to be read again.
    */

    Method string `json:"method"`//["notifications/resources/updated"]
    params: ResourceUpdatedNotificationParams


type ListPromptsRequest struct{
*PaginatedRequest[RequestParams | None, Literal["prompts/list"]]){
    /*Sent from the client to request a list of prompts and prompt templates.*/

    Method string `json:"method"`//["prompts/list"]
    Params *RequestParams `json:"params"`
}


type PromptArgument  struct {
    /*An argument for a prompt template.*/

    Name string `json:"name"`
    /*The name of the argument.*/
    Description string `json:"description"`
    /*A human-readable description of the argument.*/
    required bool `json:"listChanged"`
    /*Whether this argument must be provided.*/
    ModelConfig map[string]any `json:"model_config"`


type Prompt  struct {
    /*A prompt or prompt template that the server offers.*/

    Name string `json:"name"`
    /*The name of the prompt or prompt template.*/
    Description string `json:"description"`
    /*An optional description of what this prompt provides.*/
    arguments: list[PromptArgument] | None = None
    /*A list of arguments to use for templating the prompt.*/
    ModelConfig map[string]any `json:"model_config"`


type ListPromptsResult struct{
	*PaginatedResult
    /*The server's response to a prompts/list request from the client.*/

    prompts: list[Prompt]


type GetPromptRequestParams struct{
	*RequestParams
    /*Parameters for getting a prompt.*/

    Name string `json:"name"`
    /*The name of the prompt or prompt template.*/
    arguments: dict[str, str] | None = None
    /*Arguments to use for templating the prompt.*/
    ModelConfig map[string]any `json:"model_config"`


type GetPromptRequest struct{
 *Request[GetPromptRequestParams, Literal["prompts/get"]]){
    /*Used by the client to get a prompt provided by the server.*/

    Method string `json:"method"`//["prompts/get"]
    params: GetPromptRequestParams


type TextContent  struct {
    /*Text content for a message.*/

    type: Literal["text"]
    Text string `json:"text"`
    /*The text content of the message.*/
    Annotations *Annotations  `json:"annotations"`
    ModelConfig map[string]any `json:"model_config"`


type ImageContent  struct {
    /*Image content for a message.*/

    type: Literal["image"]
    data string `json:"code"`
    /*The base64-encoded image data.*/
    MimeType string `json:"mimeType"`
    /*
    The MIME type of the image. Different providers may support different
    image types.
    */
    Annotations *Annotations  `json:"annotations"`
    ModelConfig map[string]any `json:"model_config"`


type SamplingMessage  struct {
    /*Describes a message issued to or received from an LLM API.*/

    role: Role
    content: TextContent | ImageContent
    ModelConfig map[string]any `json:"model_config"`


type EmbeddedResource  struct {
    /*
    The contents of a resource, embedded into a prompt or tool call result.

    It is up to the client how best to render embedded resources for the benefit
    of the LLM and/or the user.
    */

    type: Literal["resource"]
    resource: TextResourceContents | BlobResourceContents
    Annotations *Annotations  `json:"annotations"`
    ModelConfig map[string]any `json:"model_config"`


type PromptMessage  struct {
    /*Describes a message returned as part of a prompt.*/

    role: Role
    content: TextContent | ImageContent | EmbeddedResource
    ModelConfig map[string]any `json:"model_config"`


type GetPromptResult struct{
 *Result
    /*The server's response to a prompts/get request from the client.*/

    Description string `json:"description"`
    /*An optional description for the prompt.*/
    messages: list[PromptMessage]


type PromptListChangedNotification(
    Notification[NotificationParams | None, Literal["notifications/prompts/list_changed"]]
){
    /*
    An optional notification from the server to the client, informing it that the list
    of prompts it offers has changed.
    */

    Method string `json:"method"`//["notifications/prompts/list_changed"]
    Params *NotificationParams `json:"params"`


type ListToolsRequest struct{
*PaginatedRequest[RequestParams | None, Literal["tools/list"]]){
    /*Sent from the client to request a list of tools the server has.*/

    Method string `json:"method"`//["tools/list"]
    Params *RequestParams `json:"params"`
}


type ToolAnnotations  struct {
    /*
    Additional properties describing a Tool to clients.

    NOTE: all properties in ToolAnnotations are **hints**.
    They are not guaranteed to provide a faithful description of
    tool behavior (including descriptive properties like `title`).

    Clients should never make tool use decisions based on ToolAnnotations
    received from untrusted servers.
    */

    title string `json:"code"`
    /*A human-readable title for the tool.*/

    readOnlyHint bool `json:"listChanged"`
    /*
    If true, the tool does not modify its environment.
    Default: false
    */

    destructiveHint bool `json:"listChanged"`
    /*
    If true, the tool may perform destructive updates to its environment.
    If false, the tool performs only additive updates.
    (This property is meaningful only when `readOnlyHint == false`)
    Default: true
    */

    idempotentHint bool `json:"listChanged"`
    /*
    If true, calling the tool repeatedly with the same arguments
    will have no additional effect on the its environment.
    (This property is meaningful only when `readOnlyHint == false`)
    Default: false
    */

    openWorldHint bool `json:"listChanged"`
    /*
    If true, this tool may interact with an "open world" of external
    entities. If false, the tool's domain of interaction is closed.
    For example, the world of a web search tool is open, whereas that
    of a memory tool is not.
    Default: true
    */
    ModelConfig map[string]any `json:"model_config"`


type Tool  struct {
    /*Definition for a tool the client can call.*/

    Name string `json:"name"`
    /*The name of the tool.*/
    Description string `json:"description"`
    /*A human-readable description of the tool.*/
    inputSchema map[string]any
    /*A JSON Schema object defining the expected parameters for the tool.*/
    annotations: ToolAnnotations | None = None
    /*Optional additional tool information.*/
    ModelConfig map[string]any `json:"model_config"`


type ListToolsResult struct{
	*PaginatedResult
    /*The server's response to a tools/list request from the client.*/

    tools: list[Tool]


type CallToolRequestParams struct{
	*RequestParams
    /*Parameters for calling a tool.*/

    Name string `json:"name"`
    arguments map[string]any | None = None
    ModelConfig map[string]any `json:"model_config"`


type CallToolRequest struct{
 *Request[CallToolRequestParams, Literal["tools/call"]]){
    /*Used by the client to invoke a tool provided by the server.*/

    Method string `json:"method"`//["tools/call"]
    params: CallToolRequestParams


type CallToolResult struct{
 *Result
    /*The server's response to a tool call.*/

    content: list[TextContent | ImageContent | EmbeddedResource]
    isError: bool = False


type ToolListChangedNotification(Notification[NotificationParams | None, Literal["notifications/tools/list_changed"]]){
    /*
    An optional notification from the server to the client, informing it that the list
    of tools it offers has changed.
    */

    Method string `json:"method"`//["notifications/tools/list_changed"]
    Params *NotificationParams `json:"params"`


LoggingLevel = Literal["debug", "info", "notice", "warning", "error", "critical", "alert", "emergency"]


type SetLevelRequestParams struct{
	*RequestParams
    /*Parameters for setting the logging level.*/

    level: LoggingLevel
    /*The level of logging that the client wants to receive from the server.*/
    ModelConfig map[string]any `json:"model_config"`


type SetLevelRequest struct{
 *Request[SetLevelRequestParams, Literal["logging/setLevel"]]){
    /*A request from the client to the server, to enable or adjust logging.*/

    Method string `json:"method"`//["logging/setLevel"]
    params: SetLevelRequestParams


type LoggingMessageNotificationParams struct{
*NotificationParams
    /*Parameters for logging message notifications.*/

    level: LoggingLevel
    /*The severity of this log message.*/
    logger string `json:"code"`
    /*An optional name of the logger issuing this message.*/
    data: Any
    /*
    The data to be logged, such as a string message or an object. Any JSON serializable
    type is allowed here.
    */
    ModelConfig map[string]any `json:"model_config"`


type LoggingMessageNotification(Notification[LoggingMessageNotificationParams, Literal["notifications/message"]]){
    /*Notification of a log message passed from server to client.*/

    Method string `json:"method"`//["notifications/message"]
    params: LoggingMessageNotificationParams


IncludeContext = Literal["none", "thisServer", "allServers"]


type ModelHint  struct {
    /*Hints to use for model selection.*/

    Name string `json:"name"`
    /*A hint for a model name.*/

    ModelConfig map[string]any `json:"model_config"`


type ModelPreferences  struct {
    /*
    The server's preferences for model selection, requested by the client during
    sampling.

    Because LLMs can vary along multiple dimensions, choosing the "best" model is
    rarely straightforward.  Different models excel in different areas—some are
    faster but less capable, others are more capable but more expensive, and so
    on. This interface allows servers to express their priorities across multiple
    dimensions to help clients make an appropriate selection for their use case.

    These preferences are always advisory. The client MAY ignore them. It is also
    up to the client to decide how to interpret these preferences and how to
    balance them against other considerations.
    */

    hints: list[ModelHint] | None = None
    /*
    Optional hints to use for model selection.

    If multiple hints are specified, the client MUST evaluate them in order
    (such that the first match is taken).

    The client SHOULD prioritize these hints over the numeric priorities, but
    MAY still use the priorities to select from ambiguous matches.
    */

    costPriority  float64 `json:"test"`
    /*
    How much to prioritize cost when selecting a model. A value of 0 means cost
    is not important, while a value of 1 means cost is the most important
    factor.
    */

    speedPriority  float64 `json:"test"`
    /*
    How much to prioritize sampling speed (latency) when selecting a model. A
    value of 0 means speed is not important, while a value of 1 means speed is
    the most important factor.
    */

    intelligencePriority  float64 `json:"test"`
    /*
    How much to prioritize intelligence and capabilities when selecting a
    model. A value of 0 means intelligence is not important, while a value of 1
    means intelligence is the most important factor.
    */

    ModelConfig map[string]any `json:"model_config"`


type CreateMessageRequestParams struct{
	*RequestParams
    /*Parameters for creating a message.*/

    messages: list[SamplingMessage]
    modelPreferences: ModelPreferences | None = None
    /*
    The server's preferences for which model to select. The client MAY ignore
    these preferences.
    */
    systemPrompt string `json:"code"`
    /*An optional system prompt the server wants to use for sampling.*/
    includeContext: IncludeContext | None = None
    /*
    A request to include context from one or more MCP servers (including the caller), to
    be attached to the prompt.
    */
    temperature  float64 `json:"test"`
    maxTokens int `json:"code"`
    /*The maximum number of tokens to sample, as requested by the server.*/
    stopSequences: list[str] | None = None
    metadata map[string]any | None = None
    /*Optional metadata to pass through to the LLM provider.*/
    ModelConfig map[string]any `json:"model_config"`


type CreateMessageRequest struct{
 *Request[CreateMessageRequestParams, Literal["sampling/createMessage"]]){
    /*A request from the server to sample an LLM via the client.*/

    Method string `json:"method"`//["sampling/createMessage"]
    params: CreateMessageRequestParams


StopReason = Literal["endTurn", "stopSequence", "maxTokens"] | str


type CreateMessageResult struct{
 *Result
    /*The client's response to a sampling/create_message request from the server.*/

    role: Role
    content: TextContent | ImageContent
    model string `json:"code"`
    /*The name of the model that generated the message.*/
    stopReason: StopReason | None = None
    /*The reason why sampling stopped, if known.*/


type ResourceReference  struct {
    /*A reference to a resource or resource template definition.*/

    type: Literal["ref/resource"]
    uri string `json:"code"`
    /*The URI or URI template of the resource.*/
    ModelConfig map[string]any `json:"model_config"`


type PromptReference  struct {
    /*Identifies a prompt.*/

    type: Literal["ref/prompt"]
    Name string `json:"name"`
    /*The name of the prompt or prompt template*/
    ModelConfig map[string]any `json:"model_config"`


type CompletionArgument  struct {
    /*The argument's information for completion requests.*/

    Name string `json:"name"`
    /*The name of the argument*/
    value string `json:"code"`
    /*The value of the argument to use for completion matching.*/
    ModelConfig map[string]any `json:"model_config"`


type CompleteRequestParams struct{
	*RequestParams
    /*Parameters for completion requests.*/

    ref: ResourceReference | PromptReference
    argument: CompletionArgument
    ModelConfig map[string]any `json:"model_config"`


type CompleteRequest struct{
 *Request[CompleteRequestParams, Literal["completion/complete"]]){
    /*A request from the client to the server, to ask for completion options.*/

    Method string `json:"method"`//["completion/complete"]
    params: CompleteRequestParams


type Completion  struct {
    /*Completion information.*/

    values: list[str]
    /*An array of completion values. Must not exceed 100 items.*/
    total: int | None = None
    /*
    The total number of completion options available. This can exceed the number of
    values actually sent in the response.
    */
    hasMore bool `json:"listChanged"`
    /*
    Indicates whether there are additional completion options beyond those provided in
    the current response, even if the exact total is unknown.
    */
    ModelConfig map[string]any `json:"model_config"`


type CompleteResult struct{
 *Result
    /*The server's response to a completion/complete request*/

    completion: Completion


type ListRootsRequest struct{
 *Request[RequestParams | None, Literal["roots/list"]]){
    /*
    Sent from the server to request a list of root URIs from the client. Roots allow
    servers to ask for specific directories or files to operate on. A common example
    for roots is providing a set of repositories or directories a server should operate
    on.

    This request is typically used when the server needs to understand the file system
    structure or access specific locations that the client has permission to read from.
    */

    Method string `json:"method"`//["roots/list"]
    Params *RequestParams `json:"params"`
}


type Root  struct {
    /*Represents a root directory or file that the server can operate on.*/

    uri: FileUrl
    /*
    The URI identifying the root. This *must* start with file:// for now.
    This restriction may be relaxed in future versions of the protocol to allow
    other URI schemes.
    */
    Name string `json:"name"`
    /*
    An optional name for the root. This can be used to provide a human-readable
    identifier for the root, which may be useful for display purposes or for
    referencing the root in other parts of the application.
    */
    ModelConfig map[string]any `json:"model_config"`


type ListRootsResult struct{
 *Result
    /*
    The client's response to a roots/list request from the server.
    This result contains an array of Root objects, each representing a root directory
    or file that the server can operate on.
    */

    roots: list[Root]


type RootsListChangedNotification(
    Notification[NotificationParams | None, Literal["notifications/roots/list_changed"]]
){
    /*
    A notification from the client to the server, informing it that the list of
    roots has changed.

    This notification should be sent whenever the client adds, removes, or
    modifies any root. The server should then request an updated list of roots
    using the ListRootsRequest.
    */

    Method string `json:"method"`//["notifications/roots/list_changed"]
    Params *NotificationParams `json:"params"`


type CancelledNotificationParams struct{
*NotificationParams
    /*Parameters for cancellation notifications.    */

    requestId: RequestId
    /*The ID of the request to cancel.    */
    reason string `json:"code"`
    /*An optional string describing the reason for the cancellation.    */
    ModelConfig map[string]any `json:"model_config"`


type CancelledNotification(Notification[CancelledNotificationParams, Literal["notifications/cancelled"]]){
    /*
    This notification can be sent by either side to indicate that it is canceling a
    previously-issued request.
    */

    Method string `json:"method"`//["notifications/cancelled"]
    params: CancelledNotificationParams


type ClientRequest(
    RootModel[
        PingRequest
        | InitializeRequest
        | CompleteRequest
        | SetLevelRequest
        | GetPromptRequest
        | ListPromptsRequest
        | ListResourcesRequest
        | ListResourceTemplatesRequest
        | ReadResourceRequest
        | SubscribeRequest
        | UnsubscribeRequest
        | CallToolRequest
        | ListToolsRequest
    ]
):
    pass


type ClientNotification(
    RootModel[CancelledNotification | ProgressNotification | InitializedNotification | RootsListChangedNotification]
):
    pass


type ClientResult(RootModel[EmptyResult | CreateMessageResult | ListRootsResult]):
    pass


type ServerRequest(RootModel[PingRequest | CreateMessageRequest | ListRootsRequest]):
    pass


type ServerNotification(
    RootModel[
        CancelledNotification
        | ProgressNotification
        | LoggingMessageNotification
        | ResourceUpdatedNotification
        | ResourceListChangedNotification
        | ToolListChangedNotification
        | PromptListChangedNotification
    ]
):
    pass


type ServerResult(
    RootModel[
        EmptyResult
        | InitializeResult
        | CompleteResult
        | GetPromptResult
        | ListPromptsResult
        | ListResourcesResult
        | ListResourceTemplatesResult
        | ReadResourceResult
        | CallToolResult
        | ListToolsResult
    ]
):
    pass


ResumptionToken = str

ResumptionTokenUpdateCallback = Callable[[ResumptionToken], None]


@dataclass
type ClientMessageMetadata{
    /*Metadata specific to client messages.    */

    resumption_token: ResumptionToken | None = None
    on_resumption_token_update: Callable[[ResumptionToken], None] | None = None


@dataclass
type ServerMessageMetadata{
    /*Metadata specific to server messages.*/

    related_request_id: RequestId | None = None
    request_context: object | None = None


MessageMetadata = ClientMessageMetadata | ServerMessageMetadata | None


@dataclass
type SessionMessage{
    /*A message with specific metadata for transport-specific features.*/

    message: JSONRPCMessage
    metadata: MessageMetadata = None


type OAuthClientMetadata  struct {
    client_name string `json:"code"`
    redirect_uris: list[str]
    grant_types: Optional[list[str]] = None
    response_types: Optional[list[str]] = None
    token_endpoint_auth_method: Optional[str] = None
    client_uri: Optional[str] = None
    scope: Optional[str] = None


type OAuthClientInformation  struct {
    client_id string `json:"code"`
    client_secret: Optional[str] = None


type OAuthClientInformationFull(OAuthClientInformation):
    client_name string `json:"code"`
    redirect_uris: list[str]
    scope: Optional[str] = None
    grant_types: Optional[list[str]] = None
    response_types: Optional[list[str]] = None
    token_endpoint_auth_method: Optional[str] = None


type OAuthTokens  struct {
    access_token string `json:"code"`
    token_type string `json:"code"`
    expires_in: Optional[int] = None
    refresh_token: Optional[str] = None
    scope: Optional[str] = None


type OAuthMetadata  struct {
    authorization_endpoint string `json:"code"`
    token_endpoint string `json:"code"`
    registration_endpoint: Optional[str] = None
    response_types_supported: list[str]
    grant_types_supported: Optional[list[str]] = None
    code_challenge_methods_supported: Optional[list[str]] = None