package layers

import (
	"mlib.com/gofy/server/core/workflow/events"
)

// Layer is a middleware hook into the workflow engine.
type Layer interface {
	OnGraphStart(graphID string, inputs map[string]any)
	OnGraphEnd(graphID string, outputs map[string]any, err error)
	OnNodeRunStart(nodeID, nodeType string, inputs map[string]any)
	OnNodeRunEnd(nodeID, nodeType string, outputs map[string]any, err error)
	OnEvent(event events.GraphEvent)
}

// BaseLayer provides default no-op implementations.
type BaseLayer struct{}

func (l *BaseLayer) OnGraphStart(graphID string, inputs map[string]any)                    {}
func (l *BaseLayer) OnGraphEnd(graphID string, outputs map[string]any, err error)           {}
func (l *BaseLayer) OnNodeRunStart(nodeID, nodeType string, inputs map[string]any)          {}
func (l *BaseLayer) OnNodeRunEnd(nodeID, nodeType string, outputs map[string]any, err error) {}
func (l *BaseLayer) OnEvent(event events.GraphEvent)                                        {}

// LayerStack manages an ordered stack of layers.
type LayerStack struct {
	layers []Layer
}

// NewLayerStack creates a new layer stack.
func NewLayerStack() *LayerStack {
	return &LayerStack{}
}

// Add adds a layer to the stack.
func (s *LayerStack) Add(layer Layer) {
	s.layers = append(s.layers, layer)
}

// OnGraphStart calls all layers in order.
func (s *LayerStack) OnGraphStart(graphID string, inputs map[string]any) {
	for _, l := range s.layers {
		l.OnGraphStart(graphID, inputs)
	}
}

// OnGraphEnd calls all layers in reverse order.
func (s *LayerStack) OnGraphEnd(graphID string, outputs map[string]any, err error) {
	for i := len(s.layers) - 1; i >= 0; i-- {
		s.layers[i].OnGraphEnd(graphID, outputs, err)
	}
}

// OnNodeRunStart calls all layers.
func (s *LayerStack) OnNodeRunStart(nodeID, nodeType string, inputs map[string]any) {
	for _, l := range s.layers {
		l.OnNodeRunStart(nodeID, nodeType, inputs)
	}
}

// OnNodeRunEnd calls all layers in reverse order.
func (s *LayerStack) OnNodeRunEnd(nodeID, nodeType string, outputs map[string]any, err error) {
	for i := len(s.layers) - 1; i >= 0; i-- {
		s.layers[i].OnNodeRunEnd(nodeID, nodeType, outputs, err)
	}
}

// OnEvent broadcasts an event to all layers.
func (s *LayerStack) OnEvent(event events.GraphEvent) {
	for _, l := range s.layers {
		l.OnEvent(event)
	}
}
