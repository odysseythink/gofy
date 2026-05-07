package trigger_webhook

import (
	basenodesentities "github.com/odysseythink/gofy/backend/entities/nodes/base"
)

// Method represents HTTP methods for webhook
type Method string

const (
	MethodGET    Method = "get"
	MethodPOST   Method = "post"
	MethodHEAD   Method = "head"
	MethodPATCH  Method = "patch"
	MethodPUT    Method = "put"
	MethodDELETE Method = "delete"
)

// ContentType represents the content type for webhook body
type ContentType string

const (
	ContentTypeJSON           ContentType = "application/json"
	ContentTypeFormData       ContentType = "multipart/form-data"
	ContentTypeFormURLEncoded ContentType = "application/x-www-form-urlencoded"
	ContentTypeText           ContentType = "text/plain"
	ContentTypeBinary         ContentType = "application/octet-stream"
)

// WebhookParameter represents a header or query parameter definition
type WebhookParameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// WebhookBodyParameter represents a body parameter definition
type WebhookBodyParameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// TriggerWebhookNodeData represents the webhook trigger node configuration
type TriggerWebhookNodeData struct {
	*basenodesentities.BaseNodeData
	Method       Method                  `json:"method"`
	ContentType  ContentType             `json:"content_type"`
	Headers      []*WebhookParameter     `json:"headers"`
	Params       []*WebhookParameter     `json:"params"`
	Body         []*WebhookBodyParameter `json:"body"`
	StatusCode   int                     `json:"status_code"`
	ResponseBody string                  `json:"response_body"`
	WebhookID    string                  `json:"webhook_id,omitempty"`
	Timeout      int                     `json:"timeout"`
}

func New() *TriggerWebhookNodeData {
	return &TriggerWebhookNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
		Method:       MethodGET,
		ContentType:  ContentTypeJSON,
		Headers:      make([]*WebhookParameter, 0),
		Params:       make([]*WebhookParameter, 0),
		Body:         make([]*WebhookBodyParameter, 0),
		StatusCode:   200,
		Timeout:      30,
	}
}
