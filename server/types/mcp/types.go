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

type RequestIdT interface {
	int | string
}

type MethodT string

type Result struct {
	/*Base class for JSON-RPC results.*/

	ModelConfig map[string]any `json:"model_config"`

	Meta map[string]any `json:"meta"` //Field(alias="_meta", default=None)
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

type JSONRPCResponse[T RequestIdT] struct {
	/*A successful (non-error) response to a request.*/

	Jsonrpc     string         `json:"jsonrpc"` //"2.0"
	ID          T              `json:"id"`
	Result      map[string]any `json:"result"`
	ModelConfig map[string]any `json:"model_config"`
}

const (
	// Standard JSON-RPC error codes
	PARSE_ERROR      = -32700
	INVALID_REQUEST  = -32600
	METHOD_NOT_FOUND = -32601
	INVALID_PARAMS   = -32602
	INTERNAL_ERROR   = -32603
)

type ErrorData struct {
	/*Error information for JSON-RPC error responses.*/

	Code int `json:"code"`
	/*The error type that occurred.*/

	Message string `json:"message"`
	/*
	   A short description of the error. The message SHOULD be limited to a concise single
	   sentence.
	*/

	Data any `json:"data"`
	/*
	   Additional information about the error. The value of this member is defined by the
	   sender (e.g. detailed error information, nested errors etc.).
	*/

	ModelConfig map[string]any `json:"model_config"`
}

type JSONRPCError[T string | int] struct {
	/*A response to a request that indicates an error occurred.*/

	Jsonrpc     string         `json:"jsonrpc"` //"2.0"
	ID          T              `json:"id"`
	Error       ErrorData      `json:"error"`
	ModelConfig map[string]any `json:"model_config"`
}

// type JSONRPCMessage(RootModel[JSONRPCRequest | JSONRPCNotification | JSONRPCResponse | JSONRPCError]):
//     pass

type EmptyResult struct {
	*Result
	/*A response that indicates success but carries no data.*/
}

type Implementation struct {
	/*Describes the name and version of an MCP implementation.*/

	Name        string         `json:"name"`
	Version     string         `json:"version"`
	ModelConfig map[string]any `json:"model_config"`
}

type RootsCapability struct {
	/*Capability for root operations.*/

	ListChanged bool `json:"listChanged"`
	/*Whether the client supports notifications for changes to the roots list.*/
	ModelConfig map[string]any `json:"model_config"`
}

type SamplingCapability struct {
	/*Capability for logging operations.*/

	ModelConfig map[string]any `json:"model_config"`
}

type ClientCapabilities struct {
	/*Capabilities a client may support.*/
	Experimental map[string]map[string]any `json:"experimental"`
	/*Experimental, non-standard capabilities that the client supports.*/
	Sampling *SamplingCapability `json:"sampling"`
	/*Present if the client supports sampling from an LLM.*/
	Roots *RootsCapability `json:"roots"`
	/*Present if the client supports listing roots.*/
	ModelConfig map[string]any `json:"model_config"`
}

type PromptsCapability struct {
	/*Capability for prompts operations.*/

	ListChanged bool `json:"listChanged"`
	/*Whether this server supports notifications for changes to the prompt list.*/
	ModelConfig map[string]any `json:"model_config"`
}

type ResourcesCapability struct {
	/*Capability for resources operations.*/

	Subscribe bool `json:"subscribe"`
	/*Whether this server supports subscribing to resource updates.*/
	ListChanged bool `json:"listChanged"`
	/*Whether this server supports notifications for changes to the resource list.*/
	ModelConfig map[string]any `json:"model_config"`
}

type ToolsCapability struct {
	/*Capability for tools operations.*/

	ListChanged bool `json:"listChanged"`
	/*Whether this server supports notifications for changes to the tool list.*/
	ModelConfig map[string]any `json:"model_config"`
}

type LoggingCapability struct {
	/*Capability for logging operations.*/

	ModelConfig map[string]any `json:"model_config"`
}

type ServerCapabilities struct {
	/*Capabilities that a server may support.*/

	Experimental map[string]map[string]any `json:"experimental"`
	/*Experimental, non-standard capabilities that the server supports.*/
	Logging *LoggingCapability `json:"logging"`
	/*Present if the server supports sending log messages to the client.*/
	Prompts *PromptsCapability `json:"prompts"`
	/*Present if the server offers any prompt templates.*/
	Resources *ResourcesCapability `json:"resources"`
	/*Present if the server offers any resources to read.*/
	Tools *ToolsCapability `json:"tools"`
	/*Present if the server offers any tools to call.*/
	ModelConfig map[string]any `json:"model_config"`
}

type InitializeResult[T string | int] struct {
	*Result
	/*After receiving an initialize request from the client, the server sends this.*/

	ProtocolVersion T `json:"protocolVersion"`
	/*The version of the Model Context Protocol that the server wants to use.*/
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   Implementation     `json:"serverInfo"`
	Instructions string             `json:"instructions"`
	/*Instructions describing how to use the server and its features.*/
}

type Annotations struct {
	Audience    []RoleType     `json:"audience"`
	Priority    float64        `json:"priority"` //Annotated[float, Field(ge=0.0, le=1.0)] | None = None
	ModelConfig map[string]any `json:"model_config"`
}

func (a *Annotations) ToDict() map[string]any {
	return map[string]any{
		"audience":     a.Audience,
		"priority":     a.Priority,
		"model_config": a.ModelConfig,
	}
}

type Resource struct {
	/*A known resource that the server is capable of reading.*/

	URI string `json:"uri"` //[AnyUrl, UrlConstraints(host_required=False)]
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
	Annotations *Annotations   `json:"annotations"`
	ModelConfig map[string]any `json:"model_config"`
}

type ResourceTemplate struct {
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
	Annotations *Annotations   `json:"annotations"`
	ModelConfig map[string]any `json:"model_config"`
}

type ListResourcesResult struct {
	*PaginatedResult
	/*The server's response to a resources/list request from the client.*/

	Resources []*Resource `json:"resources"`
}

type ListResourceTemplatesResult struct {
	*PaginatedResult
	/*The server's response to a resources/templates/list request from the client.*/

	ResourceTemplates []*ResourceTemplate `json:"resourceTemplates"`
}

type ReadResourceResult struct {
	*Result
	/*The server's response to a resources/read request from the client.*/

	Contents []ResourceContenter `json:"contents"` //list[TextResourceContents | BlobResourceContents]
}

type PromptArgument struct {
	/*An argument for a prompt template.*/

	Name string `json:"name"`
	/*The name of the argument.*/
	Description string `json:"description"`
	/*A human-readable description of the argument.*/
	Required bool `json:"required"`
	/*Whether this argument must be provided.*/
	ModelConfig map[string]any `json:"model_config"`
}

type Prompt struct {
	/*A prompt or prompt template that the server offers.*/

	Name string `json:"name"`
	/*The name of the prompt or prompt template.*/
	Description string `json:"description"`
	/*An optional description of what this prompt provides.*/
	Arguments []*PromptArgument `json:"arguments"`
	/*A list of arguments to use for templating the prompt.*/
	ModelConfig map[string]any `json:"model_config"`
}

type ListPromptsResult struct {
	*PaginatedResult
	/*The server's response to a prompts/list request from the client.*/

	Prompts []*Prompt `json:"prompts"`
}

type SamplingMessage struct {
	/*Describes a message issued to or received from an LLM API.*/

	Role        RoleType `json:"role"`
	Content     Contenter
	ModelConfig map[string]any `json:"model_config"`
}

func (msg *SamplingMessage) ToDict() map[string]any {
	if msg.Content != nil {
		return map[string]any{
			"role":         msg.Role,
			"content":      msg.Content.ToDict(),
			"model_config": msg.ModelConfig,
		}
	} else {
		return map[string]any{
			"role":         msg.Role,
			"content":      nil,
			"model_config": msg.ModelConfig,
		}
	}
}

type PromptMessage struct {
	/*Describes a message returned as part of a prompt.*/

	Role        RoleType       `json:"role"`
	Content     Contenter      `json:"content"`
	ModelConfig map[string]any `json:"model_config"`
}

type GetPromptResult struct {
	*Result
	/*The server's response to a prompts/get request from the client.*/

	Description string `json:"description"`
	/*An optional description for the prompt.*/
	Messages []*PromptMessage `json:"messages"`
}

type ToolAnnotations struct {
	/*
	   Additional properties describing a Tool to clients.

	   NOTE: all properties in ToolAnnotations are **hints**.
	   They are not guaranteed to provide a faithful description of
	   tool behavior (including descriptive properties like `title`).

	   Clients should never make tool use decisions based on ToolAnnotations
	   received from untrusted servers.
	*/

	Title string `json:"title"`
	/*A human-readable title for the tool.*/

	ReadOnlyHint bool `json:"readOnlyHint"`
	/*
	   If true, the tool does not modify its environment.
	   Default: false
	*/

	DestructiveHint bool `json:"destructiveHint"`
	/*
	   If true, the tool may perform destructive updates to its environment.
	   If false, the tool performs only additive updates.
	   (This property is meaningful only when `readOnlyHint == false`)
	   Default: true
	*/

	IdempotentHint bool `json:"idempotentHint"`
	/*
	   If true, calling the tool repeatedly with the same arguments
	   will have no additional effect on the its environment.
	   (This property is meaningful only when `readOnlyHint == false`)
	   Default: false
	*/

	OpenWorldHint bool `json:"openWorldHint"`
	/*
	   If true, this tool may interact with an "open world" of external
	   entities. If false, the tool's domain of interaction is closed.
	   For example, the world of a web search tool is open, whereas that
	   of a memory tool is not.
	   Default: true
	*/
	ModelConfig map[string]any `json:"model_config"`
}

type Tool struct {
	/*Definition for a tool the client can call.*/

	Name string `json:"name"`
	/*The name of the tool.*/
	Description string `json:"description"`
	/*A human-readable description of the tool.*/
	InputSchema map[string]any `json:"inputSchema"`
	/*A JSON Schema object defining the expected parameters for the tool.*/
	Annotations *ToolAnnotations `json:"annotations"`
	/*Optional additional tool information.*/
	ModelConfig map[string]any `json:"model_config"`
}

type ListToolsResult struct {
	*PaginatedResult
	/*The server's response to a tools/list request from the client.*/

	Tools []*Tool `json:"tools"`
}

type CallToolResult struct {
	*Result
	/*The server's response to a tool call.*/

	Content []Contenter `json:"content"`
	IsError bool        `json:"isError"`
}

type LoggingLevelType string

const (
	LoggingLevel_DEBUG     LoggingLevelType = "debug"
	LoggingLevel_INFO      LoggingLevelType = "info"
	LoggingLevel_NOTICE    LoggingLevelType = "notice"
	LoggingLevel_WARNING   LoggingLevelType = "warning"
	LoggingLevel_ERROR     LoggingLevelType = "error"
	LoggingLevel_CRITICAL  LoggingLevelType = "critical"
	LoggingLevel_ALERT     LoggingLevelType = "alert"
	LoggingLevel_EMERGENCY LoggingLevelType = "emergency"
)

type IncludeContextType string

const (
	IncludeContext_NONE       IncludeContextType = "none"
	IncludeContext_THISSERVER IncludeContextType = "thisServer"
	IncludeContext_ALLSERVERS IncludeContextType = "allServers"
)

type ModelHint struct {
	/*Hints to use for model selection.*/

	Name string `json:"name"`
	/*A hint for a model name.*/

	ModelConfig map[string]any `json:"model_config"`
}

type ModelPreferences struct {
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

	Hints []*ModelHint `json:"hints"`
	/*
	   Optional hints to use for model selection.

	   If multiple hints are specified, the client MUST evaluate them in order
	   (such that the first match is taken).

	   The client SHOULD prioritize these hints over the numeric priorities, but
	   MAY still use the priorities to select from ambiguous matches.
	*/

	CostPriority float64 `json:"costPriority"`
	/*
	   How much to prioritize cost when selecting a model. A value of 0 means cost
	   is not important, while a value of 1 means cost is the most important
	   factor.
	*/

	SpeedPriority float64 `json:"speedPriority"`
	/*
	   How much to prioritize sampling speed (latency) when selecting a model. A
	   value of 0 means speed is not important, while a value of 1 means speed is
	   the most important factor.
	*/

	IntelligencePriority float64 `json:"intelligencePriority"`
	/*
	   How much to prioritize intelligence and capabilities when selecting a
	   model. A value of 0 means intelligence is not important, while a value of 1
	   means intelligence is the most important factor.
	*/

	ModelConfig map[string]any `json:"model_config"`
}

type StopReasonType string

const (
	StopReason_ENDTURN      StopReasonType = "endTurn"
	StopReason_STOPSEQUENCE StopReasonType = "stopSequence"
	StopReason_MAXTOKENS    StopReasonType = "maxTokens"
)

type CreateMessageResult struct {
	*Result
	/*The client's response to a sampling/create_message request from the server.*/

	Role    RoleType  `json:"role"`
	Content Contenter `json:"content"`
	Model   string    `json:"model"`
	/*The name of the model that generated the message.*/
	StopReason StopReasonType `json:"stopReason"`
	/*The reason why sampling stopped, if known.*/
}

type CompletionArgument struct {
	/*The argument's information for completion requests.*/

	Name string `json:"name"`
	/*The name of the argument*/
	Value string `json:"value"`
	/*The value of the argument to use for completion matching.*/
	ModelConfig map[string]any `json:"model_config"`
}

type Completion struct {
	/*Completion information.*/

	Values []string `json:"values"`
	/*An array of completion values. Must not exceed 100 items.*/
	Total int `json:"total"`
	/*
	   The total number of completion options available. This can exceed the number of
	   values actually sent in the response.
	*/
	HasMore bool `json:"hasMore"`
	/*
	   Indicates whether there are additional completion options beyond those provided in
	   the current response, even if the exact total is unknown.
	*/
	ModelConfig map[string]any `json:"model_config"`
}

type CompleteResult struct {
	*Result
	/*The server's response to a completion/complete request*/

	Completion Completion `json:"completion"`
}

type Root struct {
	/*Represents a root directory or file that the server can operate on.*/

	URI string `json:"uri"`
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
}

type ListRootsResult struct {
	*Result
	/*
	   The client's response to a roots/list request from the server.
	   This result contains an array of Root objects, each representing a root directory
	   or file that the server can operate on.
	*/

	Roots []*Root `json:"roots"`
}

// type ClientRequest(
//     RootModel[
//         PingRequest
//         | InitializeRequest
//         | CompleteRequest
//         | SetLevelRequest
//         | GetPromptRequest
//         | ListPromptsRequest
//         | ListResourcesRequest
//         | ListResourceTemplatesRequest
//         | ReadResourceRequest
//         | SubscribeRequest
//         | UnsubscribeRequest
//         | CallToolRequest
//         | ListToolsRequest
//     ]
// ):
//     pass

// type ClientNotification(
//     RootModel[CancelledNotification | ProgressNotification | InitializedNotification | RootsListChangedNotification]
// ):
//     pass

// type ClientResult(RootModel[EmptyResult | CreateMessageResult | ListRootsResult]):
//     pass

// type ServerRequest(RootModel[PingRequest | CreateMessageRequest | ListRootsRequest]):
//     pass

// type ServerNotification(
//     RootModel[
//         CancelledNotification
//         | ProgressNotification
//         | LoggingMessageNotification
//         | ResourceUpdatedNotification
//         | ResourceListChangedNotification
//         | ToolListChangedNotification
//         | PromptListChangedNotification
//     ]
// ):
//     pass

// type ServerResult(
//     RootModel[
//         EmptyResult
//         | InitializeResult
//         | CompleteResult
//         | GetPromptResult
//         | ListPromptsResult
//         | ListResourcesResult
//         | ListResourceTemplatesResult
//         | ReadResourceResult
//         | CallToolResult
//         | ListToolsResult
//     ]
// ):
//     pass

// ResumptionToken = str

// ResumptionTokenUpdateCallback = Callable[[ResumptionToken], None]

type ClientMessageMetadata struct {
	/*Metadata specific to client messages.    */

	ResumptionToken         string `json:"resumption_token"`
	OnResumptionTokenUpdate func(string)
}
type ServerMessageMetadata[T int | string] struct {
	/*Metadata specific to server messages.*/

	RelatedRequestID T/* int | string */ `json:"related_request_id"`
	RequestContext   map[string]any `json:"request_context"`
}

// MessageMetadata = ClientMessageMetadata | ServerMessageMetadata | None

type SessionMessage[T1 ClientMessageMetadata | ServerMessageMetadata[int] | ServerMessageMetadata[string], T2 JSONRPCRequest[string] | JSONRPCRequest[int] | JSONRPCNotification | JSONRPCResponse[string] | JSONRPCResponse[int] | JSONRPCError[string] | JSONRPCError[int]] struct {
	/*A message with specific metadata for transport-specific features.*/

	Message  T2 `json:"message"`
	Metadata T1 `json:"metadata"`
}

type OAuthClientMetadata struct {
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	ClientURI               string   `json:"client_uri"`
	Scope                   string   `json:"scope"`
}

type OAuthClientInformation struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type OAuthClientInformationFull struct {
	*OAuthClientInformation
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	Scope                   string   `json:"scope"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
}

type OAuthTokens struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type OAuthMetadata struct {
	AuthorizationEndpoint         string   `json:"authorization_endpoint"`
	TokenEndpoint                 string   `json:"token_endpoint"`
	RegistrationEndpoint          string   `json:"registration_endpoint"`
	ResponseTypesSupported        []string `json:"response_types_supported"`
	GrantTypesSupported           []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
}
