package conditionhandlers

import (
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"

	wfgraph "mlib.com/gofy/server/core/workflow/graph"
	"mlib.com/gofy/server/core/workflow/utils/condition"
)

type ConditionRunConditionHandlerHandler struct {
	*BaseRunConditionHandler
}

func NewConditionRunConditionHandlerHandler(init_params *graphengineentities.GraphInitParams, graph *wfgraph.Graph, c *graphengineentities.RunCondition) RunConditionHandler {
	return &ConditionRunConditionHandlerHandler{
		BaseRunConditionHandler: NewBaseRunConditionHandler(init_params, graph, c),
	}
}
func (handler *ConditionRunConditionHandlerHandler) Check(graph_runtime_state *graphengineentities.GraphRuntimeState, previous_route_node_state *graphengineentities.RouteNodeState) bool {
	/*
	   Check if the condition can be executed

	   :param graph_runtime_state: graph runtime state
	   :param previous_route_node_state: previous route node state
	   :return: bool
	*/
	if len(handler.condition.Conditions) == 0 {
		return true
	}
	// process condition
	_, _, final_result := (&condition.ConditionProcessor{}).ProcessConditions(
		graph_runtime_state.VariablePool,
		handler.condition.Conditions,
		"and",
	)

	return final_result
}
