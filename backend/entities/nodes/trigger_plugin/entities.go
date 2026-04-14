package trigger_plugin

import (
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
)

// TriggerEventInput represents a plugin trigger event input
type TriggerEventInput struct {
	Value any    `json:"value"`
	Type  string `json:"type"` // "mixed", "variable", "constant"
}

// TriggerEventNodeData represents the plugin trigger node configuration
type TriggerEventNodeData struct {
	*basenodesentities.BaseNodeData
	PluginID               string                        `json:"plugin_id"`
	ProviderID             string                        `json:"provider_id"`
	EventName              string                        `json:"event_name"`
	SubscriptionID         string                        `json:"subscription_id"`
	PluginUniqueIdentifier string                        `json:"plugin_unique_identifier"`
	EventParameters        map[string]*TriggerEventInput `json:"event_parameters"`
}

// ResolveParameters resolves constant parameters from the event configuration
func (d *TriggerEventNodeData) ResolveParameters() map[string]any {
	result := make(map[string]any)
	for name, input := range d.EventParameters {
		if input.Type == "constant" {
			result[name] = input.Value
		}
	}
	return result
}

func New() *TriggerEventNodeData {
	return &TriggerEventNodeData{
		BaseNodeData:    &basenodesentities.BaseNodeData{},
		EventParameters: make(map[string]*TriggerEventInput),
	}
}
