package answer

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/odysseythink/mlog"
	promptutils "mlib.com/gofy/server/core/prompt/utils"
	variabletemplateparser "mlib.com/gofy/server/core/workflow/utils/variable_template_parser"
	graphengineentities "mlib.com/gofy/server/entities/graph_engine"
	answernodesentities "mlib.com/gofy/server/entities/nodes/answer"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
)

type AnswerStreamGeneratorRouter struct {
}

func (r *AnswerStreamGeneratorRouter) Init(
	node_id_config_mapping map[string]map[string]any,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge, // type: ignore[name-defined]
) *answernodesentities.AnswerStreamGenerateRoute {
	/*
	   Get stream generate routes.
	   :return:
	*/
	// parse stream output node value selectors of answer nodes
	answer_generate_route := map[string][]answernodesentities.GenerateRouteChunker{}
	for answer_node_id, node_config := range node_id_config_mapping {
		if _, ok := node_config["data"]; !ok || node_config["data"] == nil {
			continue
		}
		if _, ok := node_config["data"].(map[string]any); !ok {
			continue
		}
		if _, ok := node_config["data"].(map[string]any)["type"]; !ok {
			continue
		}
		if v, ok := node_config["data"].(map[string]any)["type"].(string); !ok || nodesenumtypes.ParseNodeTypeWithoutError(v) != nodesenumtypes.Node_ANSWER {
			continue
		}

		// get generate route for stream output
		generate_route := r.extractGenerateRouteSelectors(node_config)
		answer_generate_route[answer_node_id] = generate_route
	}
	// fetch answer dependencies
	answer_node_ids := slices.Sorted(maps.Keys(answer_generate_route))
	answer_dependencies := r.fetchAnswersDependencies(
		answer_node_ids,
		reverse_edge_mapping,
		node_id_config_mapping,
	)

	return &answernodesentities.AnswerStreamGenerateRoute{
		AnswerGenerateRoute: answer_generate_route,
		AnswerDependencies:  answer_dependencies,
	}
}

func (r *AnswerStreamGeneratorRouter) ExtractGenerateRouteFromNodeData(node_data *answernodesentities.AnswerNodeData) []answernodesentities.GenerateRouteChunker {
	/*
		Extract generate route from node data
		:param node_data: node data object
		:return:
	*/
	variable_template_parser := variabletemplateparser.NewVariableTemplateParser(node_data.Answer)
	variable_selectors := variable_template_parser.ExtractVariableSelectors()
	for _, v := range variable_selectors {
		mlog.Debugf("------variable_selector=%#v", v)
	}

	// value_selector_mapping = {
	// 	variable_selector.variable: variable_selector.value_selector for variable_selector in variable_selectors
	// }
	value_selector_mapping := map[string][]string{}
	for _, variable_selector := range variable_selectors {
		value_selector_mapping[variable_selector.Variable] = variable_selector.ValueSelector
	}

	variable_keys := []string{}
	for v := range maps.Keys(value_selector_mapping) {
		variable_keys = append(variable_keys, v)
	}

	// format answer template
	template_parser := promptutils.NewPromptTemplateParser(node_data.Answer, true)
	template_variable_keys := template_parser.VariableKeys
	mlog.Debugf("------template_variable_keys=%#v", template_variable_keys)

	// Take the intersection of variable_keys and template_variable_keys
	variable_keys = append(variable_keys, template_variable_keys...)
	sort.Strings(variable_keys)
	variable_keys = slices.Compact(variable_keys)

	template := node_data.Answer
	for _, v := range variable_keys {
		mlog.Debugf("------variable_key=%#v", v)
		template = strings.ReplaceAll(template, fmt.Sprintf("{{%s}}", v), fmt.Sprintf("Ω{{%s}}Ω", v))
	}
	mlog.Debugf("------template=%#v", template)
	generate_routes := []answernodesentities.GenerateRouteChunker{}
	for part := range strings.SplitSeq(template, "Ω") {
		if part != "" {
			if r.isVariable(part, variable_keys) {
				var_key := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(part, "Ω", ""), "{{", ""), "}}", "")
				value_selector := value_selector_mapping[var_key]
				generate_routes = append(generate_routes, &answernodesentities.VarGenerateRouteChunk{ValueSelector: value_selector})
			} else {
				generate_routes = append(generate_routes, &answernodesentities.TextGenerateRouteChunk{Text: part})
			}
		}
	}
	return generate_routes
}

func (r *AnswerStreamGeneratorRouter) extractGenerateRouteSelectors(config map[string]any) []answernodesentities.GenerateRouteChunker {
	/*
		Extract generate route selectors
		:param config: node config
		:return:
	*/
	var data map[string]any
	if _, ok := config["data"]; ok {
		if _, ok := config["data"].(map[string]any); ok {
			data = config["data"].(map[string]any)
		}
	}
	node_data := new(answernodesentities.AnswerNodeData)
	bindata, _ := json.Marshal(data)
	json.Unmarshal(bindata, node_data)

	return r.ExtractGenerateRouteFromNodeData(node_data)
}
func (r *AnswerStreamGeneratorRouter) isVariable(part string, variable_keys []string) bool {
	cleaned_part := strings.ReplaceAll(strings.ReplaceAll(part, "{{", ""), "}}", "")
	// cleaned_part = part.replace("{{", "").replace("}}", "")
	return strings.HasPrefix(part, "{{") && slices.Contains(variable_keys, cleaned_part)
}

func (r *AnswerStreamGeneratorRouter) fetchAnswersDependencies(
	answer_node_ids []string,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge, // type: ignore[name-defined]
	node_id_config_mapping map[string]map[string]any,
) map[string][]string {
	/*
		Fetch answer dependencies
		:param answer_node_ids: answer node ids
		:param reverse_edge_mapping: reverse edge mapping
		:param node_id_config_mapping: node id config mapping
		:return:
	*/
	answer_dependencies := map[string][]string{}
	for _, answer_node_id := range answer_node_ids {
		if _, ok := answer_dependencies[answer_node_id]; !ok {
			answer_dependencies[answer_node_id] = make([]string, 0)
		}
		r.recursiveFetchAnswerDependencies(
			answer_node_id,
			answer_node_id,
			node_id_config_mapping,
			reverse_edge_mapping,
			answer_dependencies,
		)
	}
	return answer_dependencies
}

func (r *AnswerStreamGeneratorRouter) recursiveFetchAnswerDependencies(
	current_node_id string,
	answer_node_id string,
	node_id_config_mapping map[string]map[string]any,
	reverse_edge_mapping map[string][]*graphengineentities.GraphEdge, // type: ignore[name-defined]
	answer_dependencies map[string][]string,
) {
	/*
		Recursive fetch answer dependencies
		:param current_node_id: current node id
		:param answer_node_id: answer node id
		:param node_id_config_mapping: node id config mapping
		:param reverse_edge_mapping: reverse edge mapping
		:param answer_dependencies: answer dependencies
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
		var source_node_data map[string]any
		if _, ok := node_id_config_mapping[source_node_id]["data"]; ok && node_id_config_mapping[source_node_id]["data"] == nil {
			if _, ok := node_id_config_mapping[source_node_id]["data"].(map[string]any); ok {
				source_node_data = node_id_config_mapping[source_node_id]["data"].(map[string]any)
			}
		}

		var source_node_type string
		if _, ok := source_node_data["type"]; ok {
			if _, ok := source_node_data["type"].(string); ok {
				source_node_type = source_node_data["type"].(string)
			}
		}
		var error_strategy string
		if _, ok := source_node_data["error_strategy"]; ok {
			if _, ok := source_node_data["error_strategy"].(string); ok {
				error_strategy = source_node_data["error_strategy"].(string)
			}
		}
		if slices.Contains([]nodesenumtypes.NodeType{
			nodesenumtypes.Node_ANSWER,
			nodesenumtypes.Node_IF_ELSE,
			nodesenumtypes.Node_QUESTION_CLASSIFIER,
			nodesenumtypes.Node_ITERATION,
			nodesenumtypes.Node_VARIABLE_ASSIGNER,
		}, nodesenumtypes.ParseNodeTypeWithoutError(source_node_type)) || error_strategy == string(nodesenumtypes.ErrorStrategy_FAIL_BRANCH) {
			answer_dependencies[answer_node_id] = append(answer_dependencies[answer_node_id], source_node_id)
		} else {
			r.recursiveFetchAnswerDependencies(
				source_node_id,
				answer_node_id,
				node_id_config_mapping,
				reverse_edge_mapping,
				answer_dependencies,
			)
		}
	}
}
