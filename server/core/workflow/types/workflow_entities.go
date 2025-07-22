package types

import (
	"time"

	"mlib.com/gofy/server/core/workflow/nodes/base"
	nodesentities "mlib.com/gofy/server/entities/nodes"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	conditionentities "mlib.com/gofy/server/entities/workflow/condition"
	appenumtypes "mlib.com/gofy/server/enum_types/app"
	"mlib.com/gofy/server/models"
)

// WorkflowNodeAndResult represents a node and its result in a workflow
type WorkflowNodeAndResult[T nodesentities.GenericNodeData, T1 conditionentities.ConditionValueType] struct {
	Node   base.BaseNode[T]                `json:"node"`
	Result *workflowentities.NodeRunResult `json:"result,omitempty"`
}

// WorkflowRunState represents the state of a workflow run
type WorkflowRunState[T nodesentities.GenericNodeData, T1 conditionentities.ConditionValueType] struct {
	TenantID     string                  `json:"tenant_id"`
	AppID        string                  `json:"app_id"`
	WorkflowID   string                  `json:"workflow_id"`
	WorkflowType models.WorkflowType     `json:"workflow_type"`
	UserID       string                  `json:"user_id"`
	UserFrom     models.UserFrom         `json:"user_from"`
	InvokeFrom   appenumtypes.InvokeFrom `json:"invoke_from"`

	WorkflowCallDepth int `json:"workflow_call_depth"`

	StartAt      time.Time                      `json:"start_at"`
	VariablePool *workflowentities.VariablePool `json:"variable_pool"`

	TotalTokens int `json:"total_tokens"`

	WorkflowNodesAndResults []*WorkflowNodeAndResult[T, T1] `json:"workflow_nodes_and_results"`

	WorkflowNodeRuns  []*NodeRun `json:"workflow_node_runs"`
	WorkflowNodeSteps int        `json:"workflow_node_steps"`

	CurrentIterationState *basenodesentities.BaseIterationState `json:"current_iteration_state,omitempty"`
}

// NodeRun represents a node run in the workflow
type NodeRun struct {
	NodeID          string `json:"node_id"`
	IterationNodeID string `json:"iteration_node_id"`
}

// NewWorkflowRunState creates a new instance of WorkflowRunState
func NewWorkflowRunState[T nodesentities.GenericNodeData, T1 conditionentities.ConditionValueType](wf *models.Workflow, startAt time.Time, variablePool *workflowentities.VariablePool, userID string, userFrom models.UserFrom, invokeFrom appenumtypes.InvokeFrom, workflowCallDepth int) *WorkflowRunState[T, T1] {
	return &WorkflowRunState[T, T1]{
		TenantID:              wf.TenantID,
		AppID:                 wf.AppID,
		WorkflowID:            wf.ID,
		WorkflowType:          wf.Type,
		UserID:                userID,
		UserFrom:              userFrom,
		InvokeFrom:            invokeFrom,
		WorkflowCallDepth:     workflowCallDepth,
		StartAt:               startAt,
		VariablePool:          variablePool,
		TotalTokens:           0,
		WorkflowNodeSteps:     1,
		WorkflowNodeRuns:      make([]*NodeRun, 0),
		CurrentIterationState: &basenodesentities.BaseIterationState{},
	}
}
