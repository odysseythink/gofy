package workflow

import (
	"fmt"

	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	workflowenumtypes "github.com/odysseythink/gofy/backend/enum_types/workflow"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

// ErrorStrategy defines how to handle node execution errors.
type ErrorStrategy string

const (
	ErrorStrategyFailBranch   ErrorStrategy = "fail-branch"
	ErrorStrategyDefaultValue ErrorStrategy = "default-value"
	ErrorStrategyRetry        ErrorStrategy = "retry"
	ErrorStrategyContinue     ErrorStrategy = "continue"
)

// ErrorHandler manages error handling during workflow execution.
type ErrorHandler struct {
	maxRetries int
}

// NewErrorHandler creates a new error handler.
func NewErrorHandler(maxRetries int) *ErrorHandler {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	return &ErrorHandler{maxRetries: maxRetries}
}

// HandleNodeError processes a node execution error according to the configured strategy.
func (h *ErrorHandler) HandleNodeError(
	nodeID string,
	strategy ErrorStrategy,
	err error,
	retryCount int,
	defaultOutput map[string]any,
) (*workflowentities.NodeRunResult, bool) {
	switch strategy {
	case ErrorStrategyRetry:
		if retryCount < h.maxRetries {
			mlog.Infof("retrying node %s (attempt %d/%d)", nodeID, retryCount+1, h.maxRetries)
			return nil, true // signal: should retry
		}
		mlog.Errorf("node %s exceeded max retries (%d), failing", nodeID, h.maxRetries)
		return h.failResult(nodeID, err), false

	case ErrorStrategyDefaultValue:
		mlog.Infof("node %s failed, using default output", nodeID)
		return &workflowentities.NodeRunResult{
			Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
			Outputs: defaultOutput,
			Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
				workflowenumtypes.NodeRunMetadataKey_ERROR_STRATEGY: string(ErrorStrategyDefaultValue),
			},
		}, false

	case ErrorStrategyContinue:
		mlog.Infof("node %s failed, continuing execution", nodeID)
		return &workflowentities.NodeRunResult{
			Status:  models.WorkflowNodeExecutionStatus_SUCCEEDED,
			Outputs: map[string]any{},
			Metadata: map[workflowenumtypes.NodeRunMetadataKey]any{
				workflowenumtypes.NodeRunMetadataKey_ERROR_STRATEGY: string(ErrorStrategyContinue),
			},
		}, false

	case ErrorStrategyFailBranch:
		fallthrough
	default:
		return h.failResult(nodeID, err), false
	}
}

func (h *ErrorHandler) failResult(nodeID string, err error) *workflowentities.NodeRunResult {
	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_FAILED,
		Error:  fmt.Sprintf("node %s failed: %v", nodeID, err),
	}
}
