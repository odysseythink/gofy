package loop

import (
	"fmt"
	"iter"
	"maps"
	"slices"

	nodesconstants "mlib.com/gofy/server/constants/workflow/nodes"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	loopnodesentities "mlib.com/gofy/server/entities/nodes/loop"
	workflowentities "mlib.com/gofy/server/entities/workflow"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

type LoopNode struct {
	*base.BaseNode[*loopnodesentities.LoopNodeData]
}

func (n *LoopNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_LOOP
}

func (n *LoopNode) Run() (*workflowentities.NodeRunResult, iter.Seq[any]) {

	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
	}, nil
}
func (n *LoopNode) ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	node_id string,
	node_data map[string]any,
) map[string][]string {
	// Create typed NodeData from dict
	// typed_node_data = LoopNodeData.model_validate(node_data)
	typed_node_data, err := mapstruct.MapToStruct1[*loopnodesentities.LoopNodeData](node_data)
	if err != nil {
		mlog.Error("convert node data to AnswerNodeData failed:%v", err)
		panic(exceptions.NewValueError("convert node data to AnswerNodeData failed"))
	}
	mlog.Debug("node_data=", typed_node_data)

	variable_mapping := map[string][]string{}

	// Extract loop node IDs statically from graph_config

	loop_node_ids := n._extract_loop_node_ids_from_config(graph_config, node_id)

	// Get node configs from graph_config
	node_configs := map[string]map[string]any{}
	for _, node := range mapstruct.Get(graph_config, "nodes", []map[string]any{}) {
		id := mapstruct.Get(node, "id", "")
		if id != "" {
			node_configs[id] = node
		}
	}

	for sub_node_id, sub_node_config := range node_configs {
		if mapstruct.Get(mapstruct.Get(sub_node_config, "data", map[string]any{}), "loop_id", "") != node_id {
			continue
		}
		var sub_node_variable_mapping map[string][]string
		// variable selector to variable mapping
		if func() bool {
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(exceptions.NotImplementedError); ok {
						sub_node_variable_mapping = map[string][]string{}
					} else {
						panic(r)
					}
				}
			}()
			// Get node class

			node_type := nodesenumtypes.NodeType(mapstruct.Get(mapstruct.Get(sub_node_config, "data", map[string]any{}), "type", ""))
			if _, ok := nodesconstants.NODE_TYPE_CLASSES_MAPPING[node_type]; !ok {
				return true
			}

			sub_node_variable_mapping = nodesconstants.ExtractVarSelectorToVarMapping(
				graph_config, sub_node_config,
			)
			return false
		}() {
			continue
		}

		// remove loop variables
		new_sub_node_variable_mapping := map[string][]string{}
		for key, value := range sub_node_variable_mapping {
			if len(value) > 0 && value[0] != node_id {
				new_sub_node_variable_mapping[sub_node_id+"."+key] = value
			}
		}
		maps.Copy(variable_mapping, new_sub_node_variable_mapping)
	}
	for _, loop_variable := range typed_node_data.LoopVariables {
		if loop_variable.ValueType == loopnodesentities.Value_Variable {
			if loop_variable.Value == nil {
				panic(exceptions.NewValueError("Loop variable value must be provided for variable type"))
			}
			// add loop variable to variable mapping
			variable_mapping[fmt.Sprintf("%s.%s", node_id, loop_variable.Label)] = loop_variable.Value.([]string)
		}
	}
	// remove variable out from loop
	maps.DeleteFunc(variable_mapping, func(key string, value []string) bool {
		return !(len(value) > 0 && slices.Contains(loop_node_ids, value[0]))
	})

	return variable_mapping
}

func (n *LoopNode) _extract_loop_node_ids_from_config(graph_config map[string]any, loop_node_id string) []string {
	/*
	   Extract node IDs that belong to a specific loop from graph configuration.

	   This method statically analyzes the graph configuration to find all nodes
	   that are part of the specified loop, without creating actual node instances.

	   :param graph_config: the complete graph configuration
	   :param loop_node_id: the ID of the loop node
	   :return: set of node IDs that belong to the loop
	*/
	loop_node_ids := map[string]struct{}{}

	// Find all nodes that belong to this loop
	for _, node := range mapstruct.Get(graph_config, "nodes", []map[string]any{}) {
		node_data := mapstruct.Get(node, "data", map[string]any{})
		if mapstruct.Get(node_data, "loop_id", "") == loop_node_id {
			node_id := mapstruct.Get(node, "id", "")
			if node_id != "" {
				loop_node_ids[node_id] = struct{}{}
			}
		}
	}

	return slices.Sorted(maps.Keys(loop_node_ids))
}
func New() *LoopNode {
	return &LoopNode{
		BaseNode: &base.BaseNode[*loopnodesentities.LoopNodeData]{},
	}
}

func init() {
	nodesconstants.Regist(New())
}
