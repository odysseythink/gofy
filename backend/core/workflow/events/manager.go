package events

import (
	"sync"

	"mlib.com/mlog"
)

// EventHandler processes workflow events.
type EventHandler func(event *GraphEvent)

// EventManager manages event registration and dispatch.
type EventManager struct {
	mu          sync.RWMutex
	handlers    map[EventType][]EventHandler
	allHandlers []EventHandler // catch-all handlers
	eventCh     chan *GraphEvent
	quit        chan struct{}
}

// NewEventManager creates a new event manager.
func NewEventManager(bufferSize int) *EventManager {
	if bufferSize <= 0 {
		bufferSize = 256
	}
	return &EventManager{
		handlers: make(map[EventType][]EventHandler),
		eventCh:  make(chan *GraphEvent, bufferSize),
		quit:     make(chan struct{}),
	}
}

// On registers a handler for a specific event type.
func (em *EventManager) On(eventType EventType, handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.handlers[eventType] = append(em.handlers[eventType], handler)
}

// OnAll registers a handler that receives all events.
func (em *EventManager) OnAll(handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.allHandlers = append(em.allHandlers, handler)
}

// Emit publishes an event.
func (em *EventManager) Emit(event *GraphEvent) {
	select {
	case em.eventCh <- event:
	default:
		mlog.Errorf("event buffer full, dropping event: %s", event.Type)
	}
}

// Start begins processing events asynchronously.
func (em *EventManager) Start() {
	go func() {
		for {
			select {
			case <-em.quit:
				return
			case event := <-em.eventCh:
				em.dispatch(event)
			}
		}
	}()
}

// Stop stops the event manager.
func (em *EventManager) Stop() {
	close(em.quit)
}

// Drain processes all remaining events synchronously.
func (em *EventManager) Drain() {
	for {
		select {
		case event := <-em.eventCh:
			em.dispatch(event)
		default:
			return
		}
	}
}

func (em *EventManager) dispatch(event *GraphEvent) {
	em.mu.RLock()
	handlers := em.handlers[event.Type]
	all := em.allHandlers
	em.mu.RUnlock()

	for _, h := range all {
		safeCall(h, event)
	}
	for _, h := range handlers {
		safeCall(h, event)
	}
}

func safeCall(handler EventHandler, event *GraphEvent) {
	defer func() {
		if r := recover(); r != nil {
			mlog.Errorf("event handler panicked on %s: %v", event.Type, r)
		}
	}()
	handler(event)
}
