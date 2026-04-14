package streamprocessor

import (
	"iter"
	"slices"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/workflow/graph"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	workflowentities "mlib.com/gofy/server/entities/workflow"
)

type StreamProcessor interface {
	Process(generator iter.Seq[graphengineentities.GraphEngineEvent]) iter.Seq[graphengineentities.GraphEngineEvent]
	Reset()
	GenerateStreamOutputsWhenNodeFinished(event *graphengineentities.NodeRunSucceededEvent) iter.Seq[graphengineentities.GraphEngineEvent]
}

type BaseStreamProcessor struct {
	Graph        *graph.Graph
	VariablePool *workflowentities.VariablePool
	RestNodeIDs  []string
}

func NewBaseStreamProcessor(gf *graph.Graph, variable_pool *workflowentities.VariablePool) *BaseStreamProcessor {
	sp := &BaseStreamProcessor{
		Graph:        gf,
		VariablePool: variable_pool,
		RestNodeIDs:  slices.Clone(gf.NodeIDs),
	}

	return sp
}

func (sp *BaseStreamProcessor) RemoveUnreachableNodes(event any) {
	var route_state *graphengineentities.RouteNodeState
	if val, ok := event.(*graphengineentities.NodeRunSucceededEvent); ok {
		route_state = val.RouteNodeState
	} else if val, ok := event.(*graphengineentities.NodeRunExceptionEvent); ok {
		route_state = val.RouteNodeState
	} else {
		mlog.Errorf("invalid event(%#v) type", event)
		return
	}
	if route_state == nil {
		mlog.Errorf("missing route_state ")
		return
	}
	finished_node_id := route_state.NodeID
	if !slices.Contains(sp.RestNodeIDs, finished_node_id) {
		mlog.Warningf("route_state(%#v)not contain in RestNodeIDs()%#v", finished_node_id, sp.RestNodeIDs)
		return
	}

	// remove finished node id
	sp.RestNodeIDs = slices.DeleteFunc(sp.RestNodeIDs, func(v string) bool {
		return finished_node_id == v
	})

	run_result := route_state.NodeRunResult
	if run_result == nil {
		return
	}

	if run_result.EdgeSourceHandle != "" {
		reachable_node_ids := []string{}
		unreachable_first_node_ids := []string{}
		if _, ok := sp.Graph.EdgeMapping[finished_node_id]; !ok {
			mlog.Warningf("node %s has no edge mapping", finished_node_id)
			return
		}
		for _, edge := range sp.Graph.EdgeMapping[finished_node_id] {
			if edge.RunCondition != nil && edge.RunCondition.BranchIdentify != "" && run_result.EdgeSourceHandle == edge.RunCondition.BranchIdentify {
				// remove unreachable nodes
				// FIXME: because of the code branch can combine directly, so for answer node
				// we remove the node maybe shortcut the answer node, so comment this code for now
				// there is not effect on the answer node and the workflow, when we have a better solution
				// we can open this code. Issues: #11542 #9560 #10638 #10564
				// ids = sp.fetchNodeIDsInReachableBranch(edge.target_node_id)
				// if "answer" in ids:
				//     continue
				// else:
				//     reachable_node_ids.extend(ids)

				// The branch_identify parameter is added to ensure that
				// only nodes in the correct logical branch are included.
				ids := sp.fetchNodeIDsInReachableBranch(edge.TargetNodeID, run_result.EdgeSourceHandle)
				reachable_node_ids = append(reachable_node_ids, ids...)
			} else {
				unreachable_first_node_ids = append(unreachable_first_node_ids, edge.TargetNodeID)
			}
			for _, node_id := range unreachable_first_node_ids {
				sp.removeNodeIDsInUnreachableBranch(node_id, reachable_node_ids)
			}
		}
	}
}

func (sp *BaseStreamProcessor) fetchNodeIDsInReachableBranch(node_id string, branch_identify string) []string {
	node_ids := []string{}
	if _, ok := sp.Graph.EdgeMapping[node_id]; ok {
		for _, edge := range sp.Graph.EdgeMapping[node_id] {
			if edge.TargetNodeID == sp.Graph.RootNodeID {
				continue
			}

			// Only follow edges that match the branch_identify or have no run_condition
			if edge.RunCondition != nil && edge.RunCondition.BranchIdentify != "" {
				if branch_identify == "" || edge.RunCondition.BranchIdentify != branch_identify {
					continue
				}
			}

			node_ids = append(node_ids, edge.TargetNodeID)
			node_ids = append(node_ids, sp.fetchNodeIDsInReachableBranch(edge.TargetNodeID, branch_identify)...)
		}
	}
	return node_ids
}

func (sp *BaseStreamProcessor) removeNodeIDsInUnreachableBranch(node_id string, reachable_node_ids []string) {
	/*
	   remove target node ids until merge
	*/
	if !slices.Contains(sp.RestNodeIDs, node_id) {
		return
	}
	sp.RestNodeIDs = slices.DeleteFunc(sp.RestNodeIDs, func(v string) bool {
		return v == node_id
	})
	if _, ok := sp.Graph.EdgeMapping[node_id]; ok {
		for _, edge := range sp.Graph.EdgeMapping[node_id] {
			if slices.Contains(reachable_node_ids, edge.TargetNodeID) {
				continue
			}

			sp.removeNodeIDsInUnreachableBranch(edge.TargetNodeID, reachable_node_ids)
		}
	}
}
