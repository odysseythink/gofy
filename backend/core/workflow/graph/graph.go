package graph

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"

	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	answergeneraterouter "mlib.com/gofy/server/core/workflow/nodes_generate_router/answer"
	endgeneraterouter "mlib.com/gofy/server/core/workflow/nodes_generate_router/end"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	answernodesentities "mlib.com/gofy/server/entities/nodes/answer"
	endnodesentities "mlib.com/gofy/server/entities/nodes/end"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"

	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/utils/validate"
	"mlib.com/mlog"
)

// Graph represents a graph structure
type Graph struct {
	RootNodeID                 string                                         `json:"root_node_id"`
	NodeIDs                    []string                                       `json:"node_ids"`
	NodeIDConfigMapping        map[string]map[string]any                      `json:"node_id_config_mapping"`
	EdgeMapping                map[string][]*graphengineentities.GraphEdge    `json:"edge_mapping"`
	ReverseEdgeMapping         map[string][]*graphengineentities.GraphEdge    `json:"reverse_edge_mapping"`
	ParallelMapping            map[string]*graphengineentities.GraphParallel  `json:"parallel_mapping"`
	NodeParallelMapping        map[string]string                              `json:"node_parallel_mapping"`
	AnswerStreamGenerateRoutes *answernodesentities.AnswerStreamGenerateRoute `json:"answer_stream_generate_routes"`
	EndStreamParam             *endnodesentities.EndStreamParam               `json:"end_stream_param"`
}

// NewGraph creates a new instance of Graph
func NewGraph(graph_config map[string]any, root_node_id string) *Graph {
	// edge configs
	var edge_configs []map[string]any
	if _, ok := graph_config["edges"]; ok {
		if _, ok := graph_config["edges"].([]any); ok {
			for _, v := range graph_config["edges"].([]any) {
				if _, ok := v.(map[string]any); ok {
					edge_configs = append(edge_configs, v.(map[string]any))
				}
			}
		}
	}
	if edge_configs == nil {
		edge_configs = []map[string]any{}
	}
	fmt.Printf("------graph_config=%#v\n", graph_config)
	// node configs
	node_configs := []map[string]any{}
	if _, ok := graph_config["nodes"]; ok {
		if _, ok := graph_config["nodes"].([]any); ok {
			for _, v := range graph_config["nodes"].([]any) {
				if _, ok := v.(map[string]any); ok {
					node_configs = append(node_configs, v.(map[string]any))
				}
			}
		}
	}

	if len(node_configs) == 0 {
		panic(exceptions.NewValueError("Graph must have at least one node"))
	}
	// reorganize edges mapping
	edge_mapping := map[string][]*graphengineentities.GraphEdge{}
	reverse_edge_mapping := map[string][]*graphengineentities.GraphEdge{}
	target_edge_ids := map[string]bool{}
	fail_branch_source_node_id := []string{}
	for _, node := range node_configs {
		err := validate.StringMapTypeVerify(node, validate.Rules{
			"id":   {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
			"data": {validate.RuleTypeOfField(reflect.Map), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("node validate failed:%v", err)
			panic(exceptions.NewValueError("node validate failed:" + err.Error()))
		}
		if _, ok := node["data"].(map[string]any)["error_strategy"]; ok {
			if val, ok := node["data"].(map[string]any)["error_strategy"].(string); ok && val == "fail-branch" {
				fail_branch_source_node_id = append(fail_branch_source_node_id, node["id"].(string))
			}
		}
	}
	mlog.Debugf("--------edge_configs=%#v", edge_configs)
	for _, edge_config := range edge_configs {
		mlog.Debugf("--------edge_config=%#v", edge_config)
		source_node_id := ""
		if _, ok := edge_config["source"]; ok {
			if _, ok := edge_config["source"].(string); ok {
				source_node_id = edge_config["source"].(string)
			}
		}
		if source_node_id == "" {
			continue
		}
		if _, ok := edge_mapping[source_node_id]; !ok {
			edge_mapping[source_node_id] = make([]*graphengineentities.GraphEdge, 0)
		}

		target_node_id := ""
		if _, ok := edge_config["target"]; ok {
			if _, ok := edge_config["target"].(string); ok {
				target_node_id = edge_config["target"].(string)
			}
		}
		if target_node_id == "" {
			continue
		}
		if _, ok := reverse_edge_mapping[target_node_id]; !ok {
			reverse_edge_mapping[target_node_id] = make([]*graphengineentities.GraphEdge, 0)
		}

		target_edge_ids[target_node_id] = true

		// parse run condition
		var run_condition *graphengineentities.RunCondition
		if _, ok := edge_config["sourceHandle"]; ok {
			if sourceHandle, ok := edge_config["sourceHandle"].(string); ok && sourceHandle != "" {
				if slices.Contains(fail_branch_source_node_id, source_node_id) && sourceHandle != "fail-branch" {
					run_condition = &graphengineentities.RunCondition{Type: "branch_identify", BranchIdentify: "success-branch"}
				} else if sourceHandle != "source" {
					run_condition = &graphengineentities.RunCondition{Type: "branch_identify", BranchIdentify: sourceHandle}
				}
			}
		}

		graph_edge := &graphengineentities.GraphEdge{SourceNodeID: source_node_id, TargetNodeID: target_node_id, RunCondition: run_condition}

		edge_mapping[source_node_id] = append(edge_mapping[source_node_id], graph_edge)
		reverse_edge_mapping[target_node_id] = append(reverse_edge_mapping[target_node_id], graph_edge)
	}

	// fetch nodes that have no predecessor node
	root_node_configs := []map[string]any{}
	all_node_id_config_mapping := map[string]map[string]any{}
	for _, node_config := range node_configs {
		err := validate.StringMapTypeVerify(node_config, validate.Rules{
			"id": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("node validate failed:%v", err)
			continue
		}

		if _, ok := target_edge_ids[node_config["id"].(string)]; !ok {
			root_node_configs = append(root_node_configs, node_config)
		}
		all_node_id_config_mapping[node_config["id"].(string)] = node_config
	}
	var root_node_ids []string
	for _, node_config := range root_node_configs {
		root_node_ids = append(root_node_ids, node_config["id"].(string))
	}
	// fetch root node
	if root_node_id == "" {
		// if no root node id, use the START type node as root node
		for _, node_config := range root_node_configs {
			if data, ok := node_config["data"].(map[string]any); ok {
				if _, ok := data["type"]; ok {
					if _, ok := data["type"].(string); ok {
						if data["type"].(string) == string(nodesenumtypes.Node_START) {
							root_node_id = node_config["id"].(string)
							break
						}
					}
				}
			}
		}
	}
	if root_node_id == "" || !slices.Contains(root_node_ids, root_node_id) {
		mlog.Errorf("Root node id %s not found in the graph", root_node_id)
		panic(exceptions.NewValueError(fmt.Sprintf("Root node id %s not found in the graph", root_node_id)))
	}
	// Check whether it is connected to the previous node
	_check_connected_to_previous_node([]string{root_node_id}, edge_mapping)

	// fetch all node ids from root node
	node_ids := []string{root_node_id}
	node_ids = recursivelyAddNodeIDs(node_ids, edge_mapping, root_node_id)

	node_id_config_mapping := map[string]map[string]any{}
	for _, node_id := range node_ids {
		node_id_config_mapping[node_id] = all_node_id_config_mapping[node_id]
	}

	// init parallel mapping
	parallel_mapping := map[string]*graphengineentities.GraphParallel{}
	node_parallel_mapping := map[string]string{}
	_recursively_add_parallels(edge_mapping, reverse_edge_mapping, root_node_id, parallel_mapping, node_parallel_mapping, nil)
	// Check if it exceeds N layers of parallel
	for _, parallel := range parallel_mapping {
		if parallel.ParentParallelID != "" {
			checkExceedParallelLimit(
				parallel_mapping,
				confy.GetWithDefault[int]("workflow.parallel_depth_limit", 3),
				parallel.ParentParallelID,
				1,
			)
		}
	}

	// init answer stream generate routes
	answer_stream_generate_routes := (&answergeneraterouter.AnswerStreamGeneratorRouter{}).Init(node_id_config_mapping, reverse_edge_mapping)

	// init end stream param
	end_stream_param := (&endgeneraterouter.EndStreamGeneratorRouter{}).Init(
		node_id_config_mapping,
		reverse_edge_mapping,
		node_parallel_mapping,
	)

	return &Graph{
		RootNodeID:                 root_node_id,
		NodeIDs:                    node_ids,
		NodeIDConfigMapping:        node_id_config_mapping,
		EdgeMapping:                edge_mapping,
		ReverseEdgeMapping:         reverse_edge_mapping,
		ParallelMapping:            parallel_mapping,
		NodeParallelMapping:        node_parallel_mapping,
		AnswerStreamGenerateRoutes: answer_stream_generate_routes,
		EndStreamParam:             end_stream_param,
	}
}

// recursivelyAddNodeIDs recursively adds node IDs
func recursivelyAddNodeIDs(nodeIDs []string, edgeMapping map[string][]*graphengineentities.GraphEdge, nodeID string) []string {
	for _, graph_edge := range edgeMapping[nodeID] {
		if !slices.Contains(nodeIDs, graph_edge.TargetNodeID) {
			nodeIDs = append(nodeIDs, graph_edge.TargetNodeID)
			nodeIDs = recursivelyAddNodeIDs(nodeIDs, edgeMapping, graph_edge.TargetNodeID)
		}
	}
	return nodeIDs
}

// _recursively_add_parallels recursively adds parallel IDs
func _recursively_add_parallels(
	edge_mapping map[string][]*graphengineentities.GraphEdge,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge,
	start_node_id string,
	parallel_mapping map[string]*graphengineentities.GraphParallel,
	node_parallel_mapping map[string]string,
	parent_parallel *graphengineentities.GraphParallel,
) {
	target_node_edges := edge_mapping[start_node_id]
	var parallel *graphengineentities.GraphParallel
	if len(target_node_edges) > 1 {
		parallel_branch_node_ids := make(map[string][]string)
		condition_edge_mappings := make(map[string][]*graphengineentities.GraphEdge)
		for _, graph_edge := range target_node_edges {
			if graph_edge.RunCondition == nil {
				if _, ok := parallel_branch_node_ids["default"]; !ok {
					parallel_branch_node_ids["default"] = make([]string, 0)
				}
				parallel_branch_node_ids["default"] = append(parallel_branch_node_ids["default"], graph_edge.TargetNodeID)
			} else {
				condition_hash := graph_edge.RunCondition.Hash()
				if _, ok := condition_edge_mappings[condition_hash]; !ok {
					condition_edge_mappings[condition_hash] = make([]*graphengineentities.GraphEdge, 0)
				}
				condition_edge_mappings[condition_hash] = append(condition_edge_mappings[condition_hash], graph_edge)
			}
		}

		for condition_hash, graph_edges := range condition_edge_mappings {
			if len(graph_edges) > 1 {
				for _, graph_edge := range graph_edges {
					if _, ok := parallel_branch_node_ids[condition_hash]; !ok {
						parallel_branch_node_ids[condition_hash] = make([]string, 0)
					}
					parallel_branch_node_ids[condition_hash] = append(parallel_branch_node_ids[condition_hash], graph_edge.TargetNodeID)
				}
			}
		}

		condition_parallels := make(map[string]*graphengineentities.GraphParallel)
		for condition_hash, condition_parallel_branch_node_ids := range parallel_branch_node_ids {
			if len(condition_parallel_branch_node_ids) > 0 {
				parent_parallel_id := ""
				if parent_parallel != nil {
					parent_parallel_id = parent_parallel.ID
				}

				parallel = &graphengineentities.GraphParallel{
					ID:               uuid.NewV4().String(),
					StartFromNodeID:  start_node_id,
					ParentParallelID: parent_parallel_id,
				}
				if parent_parallel != nil {
					parallel.ParentParallelStartNodeID = parent_parallel.StartFromNodeID
				}
				parallel_mapping[parallel.ID] = parallel
				condition_parallels[condition_hash] = parallel

				in_branch_node_ids := fetchAllNodeIDsInParallels(
					edge_mapping,
					reverse_edge_mapping,
					condition_parallel_branch_node_ids,
				)
				// collect all branches node ids
				parallel_node_ids := []string{}
				for _, node_ids := range in_branch_node_ids {
					for _, node_id := range node_ids {
						in_parent_parallel := true
						if parent_parallel_id != "" {
							in_parent_parallel = false
							for parallel_node_id, parallel_id := range node_parallel_mapping {
								if parallel_id == parent_parallel_id && parallel_node_id == node_id {
									in_parent_parallel = true
									break
								}
							}
						}

						if in_parent_parallel {
							parallel_node_ids = append(parallel_node_ids, node_id)
							node_parallel_mapping[node_id] = parallel.ID
						}
					}
				}

				outside_parallel_target_node_ids := map[string]struct{}{}
				for _, node_id := range parallel_node_ids {
					if node_id == parallel.StartFromNodeID {
						continue
					}
					if _, ok := edge_mapping[node_id]; !ok {
						continue
					}
					node_edges := edge_mapping[node_id]
					if len(node_edges) > 1 {
						continue
					}
					target_node_id := node_edges[0].TargetNodeID
					if slices.Contains(parallel_node_ids, target_node_id) {
						continue
					}
					if parent_parallel_id != "" {
						if _, ok := parallel_mapping[parent_parallel_id]; !ok {
							continue
						}
						parent_parallel = parallel_mapping[parent_parallel_id]
					}
					_, ok := node_parallel_mapping[target_node_id]
					c1 := ok && node_parallel_mapping[target_node_id] == parent_parallel_id
					c2 := parent_parallel != nil && parent_parallel.EndToNodeID != "" && target_node_id == parent_parallel.EndToNodeID
					_, ok = node_parallel_mapping[target_node_id]
					c3 := !ok && parent_parallel == nil
					if c1 || c2 || c3 {
						outside_parallel_target_node_ids[target_node_id] = struct{}{}
					}
				}

				if len(outside_parallel_target_node_ids) == 1 {
					if parent_parallel != nil && parent_parallel.EndToNodeID != "" && parallel.EndToNodeID == parent_parallel.EndToNodeID {
						parallel.EndToNodeID = ""
					} else {
						for k := range outside_parallel_target_node_ids {
							parallel.EndToNodeID = k
							delete(outside_parallel_target_node_ids, k)
							break
						}
					}
				}
			}
		}

		if len(condition_edge_mappings) > 0 {
			for condition_hash, graph_edges := range condition_edge_mappings {
				for _, graph_edge := range graph_edges {
					currentParallel := getCurrentParallel(parallel_mapping, graph_edge, condition_parallels[condition_hash], parent_parallel)
					_recursively_add_parallels(edge_mapping, reverse_edge_mapping, graph_edge.TargetNodeID, parallel_mapping, node_parallel_mapping, currentParallel)
				}
			}
		} else {
			for _, graph_edge := range target_node_edges {
				currentParallel := getCurrentParallel(parallel_mapping, graph_edge, nil, parent_parallel)
				_recursively_add_parallels(edge_mapping, reverse_edge_mapping, graph_edge.TargetNodeID, parallel_mapping, node_parallel_mapping, currentParallel)
			}
		}
	} else {
		for _, graph_edge := range target_node_edges {
			currentParallel := getCurrentParallel(parallel_mapping, graph_edge, nil, parent_parallel)
			_recursively_add_parallels(edge_mapping, reverse_edge_mapping, graph_edge.TargetNodeID, parallel_mapping, node_parallel_mapping, currentParallel)
		}
	}
}

// getCurrentParallel gets the current parallel
func getCurrentParallel(parallel_mapping map[string]*graphengineentities.GraphParallel, graph_edge *graphengineentities.GraphEdge, parallel *graphengineentities.GraphParallel, parent_parallel *graphengineentities.GraphParallel) *graphengineentities.GraphParallel {
	currentParallel := parallel
	if currentParallel == nil && parent_parallel != nil {
		if parent_parallel.EndToNodeID == "" || (parent_parallel.EndToNodeID != graph_edge.TargetNodeID) {
			currentParallel = parent_parallel
		} else {
			if parent_parallel.ParentParallelID != "" {
				parentParallelParentParallel := parallel_mapping[parent_parallel.ParentParallelID]
				if parentParallelParentParallel.EndToNodeID == "" || (parentParallelParentParallel.EndToNodeID != graph_edge.TargetNodeID) {
					currentParallel = parentParallelParentParallel
				}
			}
		}
	}
	return currentParallel
}

// checkExceedParallelLimit checks if it exceeds N layers of parallel
func checkExceedParallelLimit(parallel_mapping map[string]*graphengineentities.GraphParallel, levelLimit int, parent_parallel_id string, currentLevel int) {
	parent_parallel := parallel_mapping[parent_parallel_id]
	if parent_parallel.ID == "" {
		return
	}

	currentLevel++
	if currentLevel > levelLimit {
		panic(fmt.Sprintf("Exceeds %d layers of parallel", levelLimit))
	}

	if parent_parallel.ParentParallelID != "" {
		checkExceedParallelLimit(parallel_mapping, levelLimit, parent_parallel.ParentParallelID, currentLevel)
	}
}

// recursivelyAddParallelNodeIDs recursively adds node IDs in parallels
func recursivelyAddParallelNodeIDs(branchNodeIDs []string, edgeMapping map[string][]*graphengineentities.GraphEdge, mergeNodeID string, start_node_id string) {
	for _, graph_edge := range edgeMapping[start_node_id] {
		if graph_edge.TargetNodeID != mergeNodeID && !slices.Contains(branchNodeIDs, graph_edge.TargetNodeID) {
			branchNodeIDs = append(branchNodeIDs, graph_edge.TargetNodeID)
			recursivelyAddParallelNodeIDs(branchNodeIDs, edgeMapping, mergeNodeID, graph_edge.TargetNodeID)
		}
	}
}

// fetchAllNodeIDsInParallels fetches all node IDs in parallels
func fetchAllNodeIDsInParallels(
	edge_mapping map[string][]*graphengineentities.GraphEdge,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge,
	parallel_branch_node_ids []string,
) map[string][]string {
	routes_node_ids := make(map[string][]string)
	for _, parallelBranchNodeID := range parallel_branch_node_ids {
		routes_node_ids[parallelBranchNodeID] = []string{parallelBranchNodeID}
		recursivelyFetchRoutes(edge_mapping, parallelBranchNodeID, routes_node_ids[parallelBranchNodeID])
	}

	// fetch leaf node ids from routes node ids
	leaf_node_ids := map[string][]string{}
	merge_branch_node_ids := map[string][]string{}
	for branch_node_id, node_ids := range routes_node_ids {
		for _, node_id := range node_ids {
			if _, ok := edge_mapping[node_id]; !ok || len(edge_mapping[node_id]) == 0 {
				if _, ok := leaf_node_ids[branch_node_id]; !ok {
					leaf_node_ids[branch_node_id] = make([]string, 0)
				}
				leaf_node_ids[branch_node_id] = append(leaf_node_ids[branch_node_id], node_id)
			}

			for branch_node_id2, inner_route2 := range routes_node_ids {
				if branch_node_id != branch_node_id2 && slices.Contains(inner_route2, node_id) {
					if _, ok := reverse_edge_mapping[node_id]; ok && len(reverse_edge_mapping[node_id]) > 1 {
						if isNodeInRoutes(reverse_edge_mapping, node_id, routes_node_ids) {
							all_run_condition_is_none := true
							for _, edge := range reverse_edge_mapping[node_id] {
								if edge.RunCondition != nil {
									all_run_condition_is_none = false
								}
							}
							if all_run_condition_is_none {
								if _, ok := merge_branch_node_ids[node_id]; !ok {
									merge_branch_node_ids[node_id] = make([]string, 0)
								}
								if !slices.Contains(merge_branch_node_ids[node_id], branch_node_id2) {
									merge_branch_node_ids[node_id] = append(merge_branch_node_ids[node_id], branch_node_id2)
								}
							}
						}
					}
				}

			}
		}
	}

	// sorted merge_branch_node_ids by branch_node_ids length desc
	// merge_branch_node_ids = dict(sorted(merge_branch_node_ids.items(), key=lambda x: len(x[1]), reverse=True))

	duplicate_end_node_ids := map[[2]string]any{}
	for node_id, branch_node_ids := range merge_branch_node_ids {
		for node_id2, branch_node_ids2 := range merge_branch_node_ids {
			sort.Strings(branch_node_ids)
			sort.Strings(branch_node_ids2)
			if node_id != node_id2 && slices.Compare(slices.Compact(branch_node_ids), slices.Compact(branch_node_ids2)) == 0 {
				matched := false
				for v := range maps.Keys(duplicate_end_node_ids) {
					if slices.Compare(v[:], []string{node_id, node_id2}) == 0 || slices.Compare(v[:], []string{node_id2, node_id}) == 0 {
						matched = true
						break
					}
				}
				if !matched {
					duplicate_end_node_ids[[2]string{node_id, node_id2}] = branch_node_ids
				}
			}
		}
	}

	for /*(node_id, node_id2)*/ k := range duplicate_end_node_ids {
		// check which node is after
		if isNode2AfterNode1(k[0], k[1], edge_mapping) {
			_, ok1 := merge_branch_node_ids[k[0]]
			_, ok2 := merge_branch_node_ids[k[1]]
			if ok1 && ok2 {
				delete(merge_branch_node_ids, k[1])
			}
		} else if isNode2AfterNode1(k[1], k[0], edge_mapping) {
			_, ok1 := merge_branch_node_ids[k[0]]
			_, ok2 := merge_branch_node_ids[k[1]]
			if ok1 && ok2 {
				delete(merge_branch_node_ids, k[0])
			}
		}
	}

	branches_merge_node_ids := map[string]string{}
	for node_id, branch_node_ids := range merge_branch_node_ids {
		if len(branch_node_ids) <= 1 {
			continue
		}
		for _, branch_node_id := range branch_node_ids {
			if _, ok := branches_merge_node_ids[branch_node_id]; ok {
				continue
			}
			branches_merge_node_ids[branch_node_id] = node_id
		}
	}

	in_branch_node_ids := map[string][]string{}
	for branch_node_id, node_ids := range routes_node_ids {
		in_branch_node_ids[branch_node_id] = make([]string, 0)
		if _, ok := branches_merge_node_ids[branch_node_id]; !ok {
			// all node ids in current branch is in this thread
			in_branch_node_ids[branch_node_id] = append(in_branch_node_ids[branch_node_id], branch_node_id)
			in_branch_node_ids[branch_node_id] = append(in_branch_node_ids[branch_node_id], node_ids...)
		} else {
			merge_node_id := branches_merge_node_ids[branch_node_id]
			if merge_node_id != branch_node_id {
				in_branch_node_ids[branch_node_id] = append(in_branch_node_ids[branch_node_id], branch_node_id)
			}
			// fetch all node ids from branch_node_id and merge_node_id
			recursivelyAddParallelNodeIDs(in_branch_node_ids[branch_node_id], edge_mapping, merge_node_id, branch_node_id)
		}
	}
	return in_branch_node_ids
}

// recursivelyFetchRoutes recursively fetches routes
func recursivelyFetchRoutes(edgeMapping map[string][]*graphengineentities.GraphEdge, start_node_id string, routesNodeIDs []string) []string {
	if !slices.Contains(routesNodeIDs, start_node_id) {
		return routesNodeIDs
	}
	for _, graph_edge := range edgeMapping[start_node_id] {
		if !slices.Contains(routesNodeIDs, graph_edge.TargetNodeID) {
			routesNodeIDs = append(routesNodeIDs, graph_edge.TargetNodeID)
			routesNodeIDs = recursivelyFetchRoutes(edgeMapping, graph_edge.TargetNodeID, routesNodeIDs)
		}
	}
	return routesNodeIDs
}

// isNodeInRoutes checks if a node is in routes
func isNodeInRoutes(reverseEdgeMapping map[string][]*graphengineentities.GraphEdge, start_node_id string, routesNodeIDs map[string][]string) bool {
	// Add logic for checking if node is in routes
	// ...
	return false
}

// isNode2AfterNode1 checks if node2 is after node1
func isNode2AfterNode1(node1ID string, node2ID string, edgeMapping map[string][]*graphengineentities.GraphEdge) bool {
	for _, graph_edge := range edgeMapping[node1ID] {
		if graph_edge.TargetNodeID == node2ID {
			return true
		}
		if isNode2AfterNode1(graph_edge.TargetNodeID, node2ID, edgeMapping) {
			return true
		}
	}
	return false
}

// _check_connected_to_previous_node checks whether it is connected to the previous node
func _check_connected_to_previous_node(route []string, edgeMapping map[string][]*graphengineentities.GraphEdge) {
	lastNodeID := route[len(route)-1]
	for _, graph_edge := range edgeMapping[lastNodeID] {
		if graph_edge.TargetNodeID == "" {
			continue
		}
		if slices.Contains(route, graph_edge.TargetNodeID) {
			panic(exceptions.NewValueError(fmt.Sprintf("Node %s is connected to the previous node, please check the graph.", graph_edge.SourceNodeID)))
		}

		newRoute := slices.Clone(route)
		newRoute = append(newRoute, graph_edge.TargetNodeID)
		_check_connected_to_previous_node(newRoute, edgeMapping)
	}
}
