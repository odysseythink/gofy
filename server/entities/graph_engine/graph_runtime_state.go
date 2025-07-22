package graphengine

import (
	"time"

	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	workflowentities "mlib.com/gofy/server/entities/workflow"
)

// GraphRuntimeState represents the runtime state of a graph.
type GraphRuntimeState struct {
	VariablePool *workflowentities.VariablePool `json:"variable_pool"`
	// variable pool

	StartAt time.Time `json:"start_at"`
	// start time
	TotalTokens int `json:"total_tokens"`
	// total tokens
	LLMUsage *modelruntimeentities.LLMUsage `json:"llm_usage"`
	// llm usage info
	Outputs map[string]any `json:"outputs"`
	// outputs

	NodeRunSteps int `json:"node_run_steps"`
	// node run steps
	NodeRunState *RuntimeRouteState `json:"node_run_state"`
	// node run state
}
