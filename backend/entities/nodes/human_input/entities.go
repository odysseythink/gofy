package humaninput

import (
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
)

// FormInput represents a single form input field
type FormInput struct {
	Variable  string   `json:"variable"`
	Label     string   `json:"label"`
	Type      string   `json:"type"` // text_input, paragraph, select, number
	Required  bool     `json:"required"`
	Default   string   `json:"default,omitempty"`
	Options   []string `json:"options,omitempty"`
	MaxLength int      `json:"max_length,omitempty"`
}

// UserAction represents an action button presented to the user
type UserAction struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Style string `json:"style"` // primary, default, danger
}

// DeliveryChannelConfig represents delivery method configuration
type DeliveryChannelConfig struct {
	Type   string         `json:"type"` // webapp, email
	Config map[string]any `json:"config,omitempty"`
}

// HumanInputNodeData represents human input node data
type HumanInputNodeData struct {
	*basenodesentities.BaseNodeData
	FormContent     string                  `json:"form_content"`
	Inputs          []FormInput             `json:"inputs"`
	UserActions     []UserAction            `json:"user_actions"`
	DeliveryMethods []DeliveryChannelConfig `json:"delivery_methods"`
	Timeout         int                     `json:"timeout"`
	TimeoutUnit     string                  `json:"timeout_unit"` // minutes, hours, days
}

func New() *HumanInputNodeData {
	return &HumanInputNodeData{
		BaseNodeData: &basenodesentities.BaseNodeData{},
	}
}
