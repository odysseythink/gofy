package conditionhandlers

import (
	wfgraph "github.com/odysseythink/gofy/backend/core/workflow/graph"
	graphengineentities "github.com/odysseythink/gofy/backend/entities/graph_engine"
)

type ConditionManager struct{}

func (mgr *ConditionManager) GetConditionHandler(
	init_params *graphengineentities.GraphInitParams, graph *wfgraph.Graph, run_condition *graphengineentities.RunCondition,
) RunConditionHandler {
	/*
	   Get condition handler

	   :param init_params: init params
	   :param graph: graph
	   :param run_condition: run condition
	   :return: condition handler
	*/
	if run_condition.Type == "branch_identify" {
		return NewBranchIdentifyRunConditionHandler(init_params, graph, run_condition)
	} else {
		return NewConditionRunConditionHandlerHandler(init_params, graph, run_condition)
	}
}
