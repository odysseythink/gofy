package workflow

import (
	"sync"

	"mlib.com/gofy/server/core/workflow/events"
)

// ResponseCoordinator manages ordered streaming responses from parallel execution.
type ResponseCoordinator struct {
	mu       sync.Mutex
	buffer   map[string][]events.GraphEvent // nodeID -> buffered events
	order    []string                         // expected node ordering
	position int                              // current emit position
	outCh    chan events.GraphEvent
}

// NewResponseCoordinator creates a coordinator with expected node execution order.
func NewResponseCoordinator(nodeOrder []string, bufferSize int) *ResponseCoordinator {
	return &ResponseCoordinator{
		buffer:   make(map[string][]events.GraphEvent),
		order:    nodeOrder,
		position: 0,
		outCh:    make(chan events.GraphEvent, bufferSize),
	}
}

// Submit adds an event from a node. Events are buffered until the node's turn to emit.
func (rc *ResponseCoordinator) Submit(nodeID string, event events.GraphEvent) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.buffer[nodeID] = append(rc.buffer[nodeID], event)
	rc.tryFlush()
}

// MarkNodeComplete signals that a node has finished producing events.
func (rc *ResponseCoordinator) MarkNodeComplete(nodeID string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// If this is the current position node, advance
	if rc.position < len(rc.order) && rc.order[rc.position] == nodeID {
		rc.flushNode(nodeID)
		rc.position++
		rc.tryFlush()
	}
}

// Events returns the ordered event channel.
func (rc *ResponseCoordinator) Events() <-chan events.GraphEvent {
	return rc.outCh
}

// Close closes the output channel.
func (rc *ResponseCoordinator) Close() {
	close(rc.outCh)
}

func (rc *ResponseCoordinator) tryFlush() {
	for rc.position < len(rc.order) {
		nodeID := rc.order[rc.position]
		if evts, ok := rc.buffer[nodeID]; ok && len(evts) > 0 {
			rc.flushNode(nodeID)
			rc.position++
		} else {
			break
		}
	}
}

func (rc *ResponseCoordinator) flushNode(nodeID string) {
	if evts, ok := rc.buffer[nodeID]; ok {
		for _, evt := range evts {
			select {
			case rc.outCh <- evt:
			default:
			}
		}
		delete(rc.buffer, nodeID)
	}
}
