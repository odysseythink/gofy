package trigger

import (
	"sync"
	"time"

	"github.com/odysseythink/mlog"
)

// TriggerEvent represents a trigger event.
type TriggerEvent struct {
	ID          string
	TenantID    string
	AppID       string
	WorkflowID  string
	TriggerType string
	NodeID      string
	Data        map[string]any
	Metadata    map[string]any
	CreatedAt   time.Time
}

// EventHandler processes trigger events.
type EventHandler func(event *TriggerEvent) error

// EventBus distributes events to registered handlers.
type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
	eventCh  chan *TriggerEvent
	quit     chan struct{}
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		eventCh:  make(chan *TriggerEvent, 100),
		quit:     make(chan struct{}),
	}
}

// Subscribe registers a handler for a trigger type.
func (eb *EventBus) Subscribe(triggerType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[triggerType] = append(eb.handlers[triggerType], handler)
}

// Publish sends an event to the bus.
func (eb *EventBus) Publish(event *TriggerEvent) {
	select {
	case eb.eventCh <- event:
	default:
		mlog.Errorf("event bus full, dropping event for app %s", event.AppID)
	}
}

// Start begins processing events.
func (eb *EventBus) Start() {
	go func() {
		for {
			select {
			case <-eb.quit:
				return
			case event := <-eb.eventCh:
				eb.dispatch(event)
			}
		}
	}()
}

func (eb *EventBus) dispatch(event *TriggerEvent) {
	eb.mu.RLock()
	handlers := eb.handlers[event.TriggerType]
	allHandlers := eb.handlers["*"] // wildcard handlers
	eb.mu.RUnlock()

	for _, h := range append(handlers, allHandlers...) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					mlog.Errorf("event handler panicked: %v", r)
				}
			}()
			if err := h(event); err != nil {
				mlog.Errorf("event handler error: %v", err)
			}
		}()
	}
}

func (eb *EventBus) Stop() {
	close(eb.quit)
}
