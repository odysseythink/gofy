package layers

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/odysseythink/gofy/backend/core/workflow/events"
	"github.com/odysseythink/mlog"
)

// ExecutionLimitsLayer enforces workflow execution constraints.
type ExecutionLimitsLayer struct {
	BaseLayer
	maxNodes    int
	maxDuration time.Duration
	maxTokens   int
	nodeCount   int32
	tokenCount  int32
	startTime   time.Time
}

// NewExecutionLimitsLayer creates an execution limits layer.
func NewExecutionLimitsLayer(maxNodes int, maxDuration time.Duration, maxTokens int) *ExecutionLimitsLayer {
	if maxNodes <= 0 {
		maxNodes = 200
	}
	if maxDuration <= 0 {
		maxDuration = 30 * time.Minute
	}
	if maxTokens <= 0 {
		maxTokens = 1000000
	}
	return &ExecutionLimitsLayer{
		maxNodes:    maxNodes,
		maxDuration: maxDuration,
		maxTokens:   maxTokens,
	}
}

func (l *ExecutionLimitsLayer) OnGraphStart(graphID string, inputs map[string]any) {
	l.startTime = time.Now()
	atomic.StoreInt32(&l.nodeCount, 0)
	atomic.StoreInt32(&l.tokenCount, 0)
}

func (l *ExecutionLimitsLayer) OnNodeRunStart(nodeID, nodeType string, inputs map[string]any) {
	count := atomic.AddInt32(&l.nodeCount, 1)
	if int(count) > l.maxNodes {
		panic(fmt.Sprintf("execution limit exceeded: max %d nodes reached", l.maxNodes))
	}
	elapsed := time.Since(l.startTime)
	if elapsed > l.maxDuration {
		panic(fmt.Sprintf("execution limit exceeded: max duration %v reached", l.maxDuration))
	}
}

func (l *ExecutionLimitsLayer) OnNodeRunEnd(nodeID, nodeType string, outputs map[string]any, err error) {
	// Extract token usage from outputs if available
	if outputs != nil {
		if usage, ok := outputs["token_usage"].(map[string]any); ok {
			if total, ok := usage["total_tokens"].(float64); ok {
				newTotal := atomic.AddInt32(&l.tokenCount, int32(total))
				if int(newTotal) > l.maxTokens {
					mlog.Errorf("token limit warning: %d/%d tokens used", newTotal, l.maxTokens)
				}
			}
		}
	}
}

func (l *ExecutionLimitsLayer) OnEvent(event events.GraphEvent) {
	// Check time limit on every event
	if time.Since(l.startTime) > l.maxDuration {
		mlog.Errorf("execution duration limit exceeded")
	}
}

// NodeCount returns the current node execution count.
func (l *ExecutionLimitsLayer) NodeCount() int {
	return int(atomic.LoadInt32(&l.nodeCount))
}

// TokenCount returns the current token count.
func (l *ExecutionLimitsLayer) TokenCount() int {
	return int(atomic.LoadInt32(&l.tokenCount))
}
