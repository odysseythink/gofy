package end

import (
	"encoding/json"
	"maps"
	"slices"

	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	endnodesentities "mlib.com/gofy/server/entities/nodes/end"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/mlog"
)

type EndStreamGeneratorRouter struct {
}

func (r *EndStreamGeneratorRouter) Init(
	node_id_config_mapping map[string]map[string]any,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge, // type: ignore[name-defined]
	node_parallel_mapping map[string]string,
) *endnodesentities.EndStreamParam {
	/*
	   Get stream generate routes.
	   :return:
	*/
	// parse stream output node value selector of end nodes
	end_stream_variable_selectors_mapping := map[string][][]string{}
	for end_node_id, node_config := range node_id_config_mapping {
		if _, ok := node_config["data"]; !ok || node_config["data"] == nil {
			continue
		}
		if _, ok := node_config["data"].(map[string]any); !ok {
			continue
		}
		if _, ok := node_config["data"].(map[string]any)["type"]; !ok {
			continue
		}
		if v, ok := node_config["data"].(map[string]any)["type"].(string); !ok || nodesenumtypes.ParseNodeTypeWithoutError(v) != nodesenumtypes.Node_END {
			continue
		}

		// skip end node in parallel
		if _, ok := node_parallel_mapping[end_node_id]; ok {
			continue
		}
		// get generate route for stream output
		stream_variable_selectors := r.extractStreamVariableSelector(node_id_config_mapping, node_config)
		end_stream_variable_selectors_mapping[end_node_id] = stream_variable_selectors
	}
	// fetch end dependencies
	end_node_ids := slices.Sorted(maps.Keys(end_stream_variable_selectors_mapping))
	end_dependencies := r.fetchEndsDependencies(
		end_node_ids,
		reverse_edge_mapping,
		node_id_config_mapping,
	)

	return &endnodesentities.EndStreamParam{
		EndStreamVariableSelectorMapping: end_stream_variable_selectors_mapping,
		EndDependencies:                  end_dependencies,
	}
}
func (r *EndStreamGeneratorRouter) extractStreamVariableSelector(
	node_id_config_mapping map[string]map[string]any,
	config map[string]any,
) [][]string {
	/*
	   Extract stream variable selector from node config
	   :param node_id_config_mapping: node id config mapping
	   :param config: node config
	   :return:
	*/
	mlog.Debugf("-----------config=%#v", config)
	var node_data *endnodesentities.EndNodeData
	if _, ok := config["data"]; ok && config["data"] != nil {
		bindata, _ := json.Marshal(config["data"])
		node_data = new(endnodesentities.EndNodeData)
		err := json.Unmarshal(bindata, node_data)
		if err != nil {
			mlog.Errorf("json unmarshal failed:%v", err)
		}
	}
	return r.ExtractStreamVariableSelectorFromNodeData(node_id_config_mapping, node_data)
}

func (r *EndStreamGeneratorRouter) fetchEndsDependencies(
	end_node_ids []string,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge, // type: ignore[name-defined]
	node_id_config_mapping map[string]map[string]any,
) map[string][]string {
	/*
	   Fetch end dependencies
	   :param end_node_ids: end node ids
	   :param reverse_edge_mapping: reverse edge mapping
	   :param node_id_config_mapping: node id config mapping
	   :return:
	*/
	end_dependencies := map[string][]string{}
	for _, end_node_id := range end_node_ids {
		if _, ok := end_dependencies[end_node_id]; !ok {
			end_dependencies[end_node_id] = make([]string, 0)
		}
		r.recursiveFetchEndDependencies(
			end_node_id,
			end_node_id,
			node_id_config_mapping,
			reverse_edge_mapping,
			end_dependencies,
		)
	}
	return end_dependencies
}

func (r *EndStreamGeneratorRouter) recursiveFetchEndDependencies(
	current_node_id string,
	end_node_id string,
	node_id_config_mapping map[string]map[string]any,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge, // type: ignore[name-defined]
	end_dependencies map[string][]string,
) {
	/*
	   Recursive fetch end dependencies
	   :param current_node_id: current node id
	   :param end_node_id: end node id
	   :param node_id_config_mapping: node id config mapping
	   :param reverse_edge_mapping: reverse edge mapping
	   :param end_dependencies: end dependencies
	   :return:
	*/
	var reverse_edges []*graphengineentities.GraphEdge
	if _, ok := reverse_edge_mapping[current_node_id]; ok {
		reverse_edges = reverse_edge_mapping[current_node_id]
	}
	for _, edge := range reverse_edges {
		source_node_id := edge.SourceNodeID
		if _, ok := node_id_config_mapping[source_node_id]; !ok {
			continue
		}
		source_node_type := ""
		if _, ok := node_id_config_mapping[source_node_id]["data"]; ok {
			if _, ok := node_id_config_mapping[source_node_id]["data"].(map[string]any); ok {
				if _, ok := node_id_config_mapping[source_node_id]["data"].(map[string]any)["type"]; ok {
					if _, ok := node_id_config_mapping[source_node_id]["data"].(map[string]any)["type"].(string); ok {
						source_node_type = node_id_config_mapping[source_node_id]["data"].(map[string]any)["type"].(string)
					}
				}
			}
		}
		if slices.Contains([]nodesenumtypes.NodeType{nodesenumtypes.Node_IF_ELSE, nodesenumtypes.Node_QUESTION_CLASSIFIER}, nodesenumtypes.ParseNodeTypeWithoutError(source_node_type)) {
			end_dependencies[end_node_id] = append(end_dependencies[end_node_id], source_node_id)
		} else {
			r.recursiveFetchEndDependencies(
				source_node_id,
				end_node_id,
				node_id_config_mapping,
				reverse_edge_mapping,
				end_dependencies,
			)
		}
	}
}
func (r *EndStreamGeneratorRouter) ExtractStreamVariableSelectorFromNodeData(
	node_id_config_mapping map[string]map[string]any,
	node_data *endnodesentities.EndNodeData,
) [][]string {
	/*
	   Extract stream variable selector from node data
	   :param node_id_config_mapping: node id config mapping
	   :param node_data: node data object
	   :return:
	*/
	variable_selectors := node_data.Outputs

	value_selectors := [][]string{}
	for _, variable_selector := range variable_selectors {
		if len(variable_selector.ValueSelector) == 0 {
			continue
		}
		node_id := variable_selector.ValueSelector[0]
		if node_id != "sys" && slices.Contains(slices.Sorted(maps.Keys(node_id_config_mapping)), node_id) {
			node := node_id_config_mapping[node_id]
			node_type := ""
			if _, ok := node["data"]; ok && node["data"] != nil {
				if _, ok := node["data"].(map[string]any); ok && node["data"].(map[string]any) != nil {
					if _, ok := node["data"].(map[string]any)["type"]; ok && node["data"].(map[string]any)["type"] != nil {
						if _, ok := node["data"].(map[string]any)["type"].(string); ok {
							node_type = node["data"].(map[string]any)["type"].(string)
						}
					}
				}
			}
			if !slices.ContainsFunc(value_selectors, func(a []string) bool {
				return slices.Compare(variable_selector.ValueSelector, a) == 0
			}) &&
				nodesenumtypes.ParseNodeTypeWithoutError(node_type) == nodesenumtypes.Node_LLM &&
				variable_selector.ValueSelector[1] == "text" {
				value_selectors = append(value_selectors, variable_selector.ValueSelector)
			}
		}
	}
	return value_selectors
}
