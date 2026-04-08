package plugin

import "context"

// Session represents a plugin execution session.
// It carries the identity and context needed for a single plugin invocation.
type Session struct {
	ID                     string `json:"id"`
	TenantID               string `json:"tenant_id"`
	UserID                 string `json:"user_id"`
	PluginUniqueIdentifier string `json:"plugin_unique_identifier"`
	ClusterID              string `json:"cluster_id"`
	InvokeFrom             string `json:"invoke_from"`
	Action                 string `json:"action"`

	// Information about the incoming request
	AppID          string `json:"app_id"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	EndpointID     string `json:"endpoint_id"`

	// Arbitrary context data passed to the plugin
	ContextData map[string]any `json:"context"`

	// Internal references (not serialized)
	runtime             PluginRuntimeSessionIOInterface `json:"-"`
	backwardsInvocation BackwardsInvocation             `json:"-"`
	requestContext      context.Context                 `json:"-"`
}

// BindRuntime attaches a runtime I/O interface to this session.
func (s *Session) BindRuntime(runtime PluginRuntimeSessionIOInterface) {
	s.runtime = runtime
}

// Runtime returns the bound runtime I/O interface.
func (s *Session) Runtime() PluginRuntimeSessionIOInterface {
	return s.runtime
}

// BindBackwardsInvocation attaches a backwards invocation handler to this session.
func (s *Session) BindBackwardsInvocation(bi BackwardsInvocation) {
	s.backwardsInvocation = bi
	if s.requestContext != nil {
		s.backwardsInvocation.SetContext(s.requestContext)
	}
}

// BackwardsInvocation returns the bound backwards invocation handler.
func (s *Session) BackwardsInvocation() BackwardsInvocation {
	return s.backwardsInvocation
}

// RequestContext returns the request context, defaulting to context.Background().
func (s *Session) RequestContext() context.Context {
	if s.requestContext == nil {
		return context.Background()
	}
	return s.requestContext
}

// SetRequestContext sets the request context for the session.
func (s *Session) SetRequestContext(ctx context.Context) {
	s.requestContext = ctx
	if s.backwardsInvocation != nil {
		s.backwardsInvocation.SetContext(ctx)
	}
}

// SessionStreamEvent represents the type of event sent to the plugin.
type SessionStreamEvent string

const (
	SESSION_STREAM_EVENT_REQUEST  SessionStreamEvent = "request"
	SESSION_STREAM_EVENT_RESPONSE SessionStreamEvent = "backwards_response"
)
