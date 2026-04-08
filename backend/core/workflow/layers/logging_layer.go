package layers

import (
	"time"

	"mlib.com/gofy/server/core/workflow/events"
	"mlib.com/mlog"
)

// LoggingLayer logs workflow execution events.
type LoggingLayer struct {
	BaseLayer
	startTimes map[string]time.Time
}

// NewLoggingLayer creates a logging layer.
func NewLoggingLayer() *LoggingLayer {
	return &LoggingLayer{startTimes: make(map[string]time.Time)}
}

func (l *LoggingLayer) OnGraphStart(graphID string, inputs map[string]any) {
	l.startTimes[graphID] = time.Now()
	mlog.Infof("[workflow] graph %s started", graphID)
}

func (l *LoggingLayer) OnGraphEnd(graphID string, outputs map[string]any, err error) {
	elapsed := time.Since(l.startTimes[graphID])
	if err != nil {
		mlog.Errorf("[workflow] graph %s failed after %v: %v", graphID, elapsed, err)
	} else {
		mlog.Infof("[workflow] graph %s completed in %v", graphID, elapsed)
	}
	delete(l.startTimes, graphID)
}

func (l *LoggingLayer) OnNodeRunStart(nodeID, nodeType string, inputs map[string]any) {
	l.startTimes[nodeID] = time.Now()
	mlog.Infof("[workflow] node %s (%s) started", nodeID, nodeType)
}

func (l *LoggingLayer) OnNodeRunEnd(nodeID, nodeType string, outputs map[string]any, err error) {
	elapsed := time.Since(l.startTimes[nodeID])
	if err != nil {
		mlog.Errorf("[workflow] node %s (%s) failed after %v: %v", nodeID, nodeType, elapsed, err)
	} else {
		mlog.Infof("[workflow] node %s (%s) completed in %v", nodeID, nodeType, elapsed)
	}
	delete(l.startTimes, nodeID)
}

func (l *LoggingLayer) OnEvent(event events.GraphEvent) {
	// Only log non-routine events
	switch event.GetEventType() {
	case events.EventGraphRunPaused, events.EventGraphRunResumed, events.EventGraphRunFailed:
		mlog.Infof("[workflow] event: %s", event.GetEventType())
	}
}
