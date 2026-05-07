package layers

import (
	"encoding/json"
	"time"

	"github.com/odysseythink/gofy/backend/core/workflow/events"
	"github.com/odysseythink/mlog"
)

// DebugLayer provides detailed debugging output for workflow execution.
type DebugLayer struct {
	BaseLayer
	enabled    bool
	eventLog   []DebugEntry
	startTimes map[string]time.Time
}

// DebugEntry records a debug event.
type DebugEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	EventType string         `json:"event_type"`
	NodeID    string         `json:"node_id,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}

// NewDebugLayer creates a debug layer.
func NewDebugLayer(enabled bool) *DebugLayer {
	return &DebugLayer{
		enabled:    enabled,
		startTimes: make(map[string]time.Time),
	}
}

func (l *DebugLayer) OnGraphStart(graphID string, inputs map[string]any) {
	if !l.enabled {
		return
	}
	l.addEntry("graph_start", "", map[string]any{
		"graph_id":    graphID,
		"input_count": len(inputs),
	})
	mlog.Infof("[DEBUG] Graph %s starting with %d inputs", graphID, len(inputs))
}

func (l *DebugLayer) OnGraphEnd(graphID string, outputs map[string]any, err error) {
	if !l.enabled {
		return
	}
	details := map[string]any{
		"graph_id":     graphID,
		"output_count": len(outputs),
		"total_events": len(l.eventLog),
	}
	if err != nil {
		details["error"] = err.Error()
	}
	l.addEntry("graph_end", "", details)
	mlog.Infof("[DEBUG] Graph %s ended. Total events: %d", graphID, len(l.eventLog))
}

func (l *DebugLayer) OnNodeRunStart(nodeID, nodeType string, inputs map[string]any) {
	if !l.enabled {
		return
	}
	l.startTimes[nodeID] = time.Now()
	inputSummary := summarizeMap(inputs, 200)
	l.addEntry("node_start", nodeID, map[string]any{
		"node_type": nodeType,
		"inputs":    inputSummary,
	})
	mlog.Infof("[DEBUG] Node %s (%s) starting", nodeID, nodeType)
}

func (l *DebugLayer) OnNodeRunEnd(nodeID, nodeType string, outputs map[string]any, err error) {
	if !l.enabled {
		return
	}
	elapsed := time.Since(l.startTimes[nodeID])
	outputSummary := summarizeMap(outputs, 200)
	details := map[string]any{
		"node_type":  nodeType,
		"elapsed_ms": elapsed.Milliseconds(),
		"outputs":    outputSummary,
	}
	if err != nil {
		details["error"] = err.Error()
	}
	l.addEntry("node_end", nodeID, details)
	delete(l.startTimes, nodeID)
	mlog.Infof("[DEBUG] Node %s (%s) finished in %v", nodeID, nodeType, elapsed)
}

func (l *DebugLayer) OnEvent(event events.GraphEvent) {
	if !l.enabled {
		return
	}
	l.addEntry(string(event.GetEventType()), "", nil)
}

// GetEventLog returns the debug event log.
func (l *DebugLayer) GetEventLog() []DebugEntry {
	return l.eventLog
}

// GetEventLogJSON returns the event log as JSON.
func (l *DebugLayer) GetEventLogJSON() string {
	data, _ := json.MarshalIndent(l.eventLog, "", "  ")
	return string(data)
}

func (l *DebugLayer) addEntry(eventType, nodeID string, details map[string]any) {
	l.eventLog = append(l.eventLog, DebugEntry{
		Timestamp: time.Now(),
		EventType: eventType,
		NodeID:    nodeID,
		Details:   details,
	})
}

func summarizeMap(m map[string]any, maxLen int) string {
	if m == nil {
		return "{}"
	}
	data, _ := json.Marshal(m)
	s := string(data)
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
