package mcp

type RequestParams[T string | int] struct {
	Meta struct {
		ProgressToken T `json:"progressToken"` // ProgressToken | None = None
		/*
		   If specified, the caller requests out-of-band progress notifications for
		   this request (as represented by notifications/progress). The value of this
		   parameter is an opaque token that will be attached to any subsequent
		   notifications. The receiver is not obligated to provide these notifications.
		*/

		ModelConfig map[string]any `json:"model_config"`
	} `json:"meta"` //Field(alias="_meta", default=None)
}

type RequestParamsT interface {
	RequestParams[string] |
		RequestParams[int] |
		map[string]any |
		InitializeRequestParams[string, int] |
		InitializeRequestParams[string, string] |
		InitializeRequestParams[int, int] |
		InitializeRequestParams[int, string] |
		ReadResourceRequestParams[string] |
		ReadResourceRequestParams[int] |
		SubscribeRequestParams[string] |
		SubscribeRequestParams[int] |
		UnsubscribeRequestParams[string] |
		UnsubscribeRequestParams[int] |
		GetPromptRequestParams[string] |
		GetPromptRequestParams[int] |
		CallToolRequestParams[string] |
		CallToolRequestParams[int] |
		SetLevelRequestParams[string] |
		SetLevelRequestParams[int] |
		CreateMessageRequestParams[string] |
		CreateMessageRequestParams[int] |
		CompleteRequestParams[string] |
		CompleteRequestParams[int]
}

type Request[T RequestParamsT] struct {
	/*Base class for JSON-RPC requests.*/

	Method      string         `json:"method"`
	Params      T              `json:"params"` /*RequestParams |  map[string]any*/
	ModelConfig map[string]any `json:"model_config"`
}

type PaginatedRequest[T RequestParamsT] struct {
	*Request[T]
	Cursor `json:"cursor"`
	/*
	   An opaque token representing the current pagination position.
	   If provided, the server should return results starting after this cursor.
	*/
}

type JSONRPCRequest[T RequestIdT] struct {
	*Request[map[string]any]
	/*A request that expects a response.*/

	Jsonrpc string         `json:"jsonrpc"` //"2.0"
	ID      T              `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

type InitializeRequestParams[T1 string | int, T2 string | int] struct {
	*RequestParams[T2]
	/*Parameters for the initialize request.*/

	ProtocolVersion T1 `json:"protocolVersion"`
	/*The latest version of the Model Context Protocol that the client supports.*/
	Capabilities ClientCapabilities `json:"capabilities"`
	ClientInfo   Implementation     `json:"clientInfo"`
	ModelConfig  map[string]any     `json:"model_config"`
}

type InitializeRequest[T InitializeRequestParams[string, int] | InitializeRequestParams[string, string] | InitializeRequestParams[int, int] | InitializeRequestParams[int, string]] struct {
	*Request[T]
	/*
	   This request is sent from the client to the server when it first connects, asking it
	   to begin initialization.
	*/

	Method string `json:"method"` //["initialize"]
	Params T      `json:"params"`
}

type PingRequest[T RequestParams[string] | RequestParams[int]] struct {
	*Request[T]
	/*
	   A ping, issued by either the server or the client, to check that the other party is
	   still alive.
	*/

	Method string `json:"method"` //["ping"]
	Params *T     `json:"params"`
}

type ListResourcesRequest[T RequestParams[string] | RequestParams[int]] struct {
	*PaginatedRequest[T]
	/*Sent from the client to request a list of resources the server has.*/

	Method string `json:"method"` //["resources/list"]
	Params *T     `json:"params"`
}

type ListResourceTemplatesRequest[T RequestParams[string] | RequestParams[int]] struct {
	*PaginatedRequest[T]
	/*Sent from the client to request a list of resource templates the server has.*/

	Method string `json:"method"` //["resources/templates/list"]
	Params *T     `json:"params"`
}

type ReadResourceRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for reading a resource.*/

	URI string `json:"uri"` //[AnyUrl, UrlConstraints(host_required=False)]
	/*
	   The URI of the resource to read. The URI can use any protocol; it is up to the
	   server how to interpret it.
	*/
	ModelConfig map[string]any `json:"model_config"`
}

type ReadResourceRequest[T ReadResourceRequestParams[string] | ReadResourceRequestParams[int]] struct {
	*Request[T]
	/*Sent from the client to the server, to read a specific resource URI.*/

	Method string `json:"method"` //["resources/read"]
	Params T      `json:"params"`
}

type SubscribeRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for subscribing to a resource.*/

	URI string `json:"uri"` //[AnyUrl, UrlConstraints(host_required=False)]
	/*
	   The URI of the resource to subscribe to. The URI can use any protocol; it is up to
	   the server how to interpret it.
	*/
	ModelConfig map[string]any `json:"model_config"`
}

type SubscribeRequest[T SubscribeRequestParams[string] | SubscribeRequestParams[int]] struct {
	*Request[T]
	/*
	   Sent from the client to request resources/updated notifications from the server
	   whenever a particular resource changes.
	*/

	Method string `json:"method"` //["resources/subscribe"]
	Params T      `json:"params"`
}

type UnsubscribeRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for unsubscribing from a resource.*/

	URI string `json:"uri"` //[AnyUrl, UrlConstraints(host_required=False)]
	/*The URI of the resource to unsubscribe from.*/
	ModelConfig map[string]any `json:"model_config"`
}

type UnsubscribeRequest[T UnsubscribeRequestParams[string] | UnsubscribeRequestParams[int]] struct {
	*Request[T]
	/*
	   Sent from the client to request cancellation of resources/updated notifications from
	   the server.
	*/

	Method string `json:"method"` //["resources/unsubscribe"]
	Params T      `json:"params"`
}

type ListPromptsRequest[T RequestParams[string] | RequestParams[int]] struct {
	*PaginatedRequest[T]
	/*Sent from the client to request a list of prompts and prompt templates.*/

	Method string `json:"method"` //["prompts/list"]
	Params *T     `json:"params"`
}

type GetPromptRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for getting a prompt.*/

	Name string `json:"name"`
	/*The name of the prompt or prompt template.*/
	Arguments map[string]string `json:"arguments"`
	/*Arguments to use for templating the prompt.*/
	ModelConfig map[string]any `json:"model_config"`
}

type GetPromptRequest[T GetPromptRequestParams[string] | GetPromptRequestParams[int]] struct {
	*Request[T]
	/*Used by the client to get a prompt provided by the server.*/

	Method string `json:"method"` //["prompts/get"]
	Params T      `json:"params"`
}
type ListToolsRequest[T RequestParams[string] | RequestParams[int]] struct {
	*PaginatedRequest[T]
	/*Sent from the client to request a list of tools the server has.*/

	Method string `json:"method"` //["tools/list"]
	Params *T     `json:"params"`
}

type CallToolRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for calling a tool.*/

	Name        string         `json:"name"`
	Arguments   map[string]any `json:"arguments"`
	ModelConfig map[string]any `json:"model_config"`
}

type CallToolRequest[T CallToolRequestParams[string] | CallToolRequestParams[int]] struct {
	*Request[T]
	/*Used by the client to invoke a tool provided by the server.*/

	Method string `json:"method"` //["tools/call"]
	Params T      `json:"params"`
}

type SetLevelRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for setting the logging level.*/

	Level LoggingLevelType `json:"level"`
	/*The level of logging that the client wants to receive from the server.*/
	ModelConfig map[string]any `json:"model_config"`
}

type SetLevelRequest[T SetLevelRequestParams[string] | SetLevelRequestParams[int]] struct {
	*Request[T]
	/*A request from the client to the server, to enable or adjust logging.*/

	Method string `json:"method"` //["logging/setLevel"]
	Params T      `json:"params"`
}

type CreateMessageRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for creating a message.*/

	Messages         []*SamplingMessage `json:"messages"`
	ModelPreferences *ModelPreferences  `json:"modelPreferences"`
	/*
	   The server's preferences for which model to select. The client MAY ignore
	   these preferences.
	*/
	SystemPrompt string `json:"systemPrompt"`
	/*An optional system prompt the server wants to use for sampling.*/
	IncludeContext IncludeContextType `json:"includeContext"`
	/*
	   A request to include context from one or more MCP servers (including the caller), to
	   be attached to the prompt.
	*/
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
	/*The maximum number of tokens to sample, as requested by the server.*/
	StopSequences []string       `json:"stopSequences"`
	Metadata      map[string]any `json:"metadata"`
	/*Optional metadata to pass through to the LLM provider.*/
	ModelConfig map[string]any `json:"model_config"`
}

type CreateMessageRequest[T CreateMessageRequestParams[string] | CreateMessageRequestParams[int]] struct {
	*Request[T]
	/*A request from the server to sample an LLM via the client.*/

	Method string `json:"method"` //["sampling/createMessage"]
	Params T      `json:"params"`
}

type CompleteRequestParams[T string | int] struct {
	*RequestParams[T]
	/*Parameters for completion requests.*/

	Ref         Referencer         `json:"ref"`
	Argument    CompletionArgument `json:"argument"`
	ModelConfig map[string]any     `json:"model_config"`
}

type CompleteRequest[T CompleteRequestParams[string] | CompleteRequestParams[int]] struct {
	*Request[T]
	/*A request from the client to the server, to ask for completion options.*/

	Method string `json:"method"` //["completion/complete"]
	Params T      `json:"params"`
}

type ListRootsRequest[T RequestParams[string] | RequestParams[int]] struct {
	*Request[T]
	/*
	   Sent from the server to request a list of root URIs from the client. Roots allow
	   servers to ask for specific directories or files to operate on. A common example
	   for roots is providing a set of repositories or directories a server should operate
	   on.

	   This request is typically used when the server needs to understand the file system
	   structure or access specific locations that the client has permission to read from.
	*/

	Method string `json:"method"` //["roots/list"]
	Params *T     `json:"params"`
}
