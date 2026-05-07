package conditionhandlers

import (
	wfgraph "github.com/odysseythink/gofy/backend/core/workflow/graph"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
)

type RunConditionHandler interface {
	/*
		Check if the condition can be executed

		:param graph_runtime_state: graph runtime state
		:param previous_route_node_state: previous route node state
		:return: bool
	*/
	Check(graph_runtime_state *graphengineentities.GraphRuntimeState, previous_route_node_state *graphengineentities.RouteNodeState) bool
}

type BaseRunConditionHandler struct {
	initParams *graphengineentities.GraphInitParams
	graph      *wfgraph.Graph
	condition  *graphengineentities.RunCondition
}

func NewBaseRunConditionHandler(init_params *graphengineentities.GraphInitParams, graph *wfgraph.Graph, condition *graphengineentities.RunCondition) *BaseRunConditionHandler {
	return &BaseRunConditionHandler{
		initParams: init_params,
		graph:      graph,
		condition:  condition,
	}
}
