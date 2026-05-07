package conditionhandlers

import (
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
	"github.com/odysseythink/mlog"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	wfgraph "github.com/odysseythink/gofy/backend/core/workflow/graph"
)

type BranchIdentifyRunConditionHandler struct {
	*BaseRunConditionHandler
}

func NewBranchIdentifyRunConditionHandler(init_params *graphengineentities.GraphInitParams, graph *wfgraph.Graph, condition *graphengineentities.RunCondition) RunConditionHandler {
	return &BranchIdentifyRunConditionHandler{
		BaseRunConditionHandler: NewBaseRunConditionHandler(init_params, graph, condition),
	}
}

func (handler *BranchIdentifyRunConditionHandler) Check(graph_runtime_state *graphengineentities.GraphRuntimeState, previous_route_node_state *graphengineentities.RouteNodeState) bool {
	/*
	   Check if the condition can be executed

	   :param graph_runtime_state: graph runtime state
	   :param previous_route_node_state: previous route node state
	   :return: bool
	*/
	mlog.Debugf("------handler.condition=%#v", handler.condition)
	if handler.condition.BranchIdentify == "" {
		panic(exceptions.NewValueError("Branch identify is required"))
	}
	run_result := previous_route_node_state.NodeRunResult
	mlog.Debugf("------run_result=%#v", run_result)
	if run_result == nil {
		return false
	}
	if run_result.EdgeSourceHandle == "" {
		return false
	}
	return handler.condition.BranchIdentify == run_result.EdgeSourceHandle
}
