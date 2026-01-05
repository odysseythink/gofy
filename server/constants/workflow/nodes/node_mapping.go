package nodes

import (
	"errors"
	"fmt"
	"log"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/workflow/nodes/base"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/utils/mapstruct"
)

const (
	LATEST_VERSION = "latest"
)

var (
	NODE_TYPE_CLASSES_MAPPING = map[nodesenumtypes.NodeType]map[string]base.Noder{}
)

func Regist(n base.Noder) error {
	if n == nil {
		log.Println("missing noder")
		return errors.New("missing noder")
	}
	if _, ok := NODE_TYPE_CLASSES_MAPPING[n.Type()]; ok {
		log.Printf("node(%s) already regist\n", n.Type())
		return fmt.Errorf("node(%s) already regist", n.Type())
	}
	NODE_TYPE_CLASSES_MAPPING[n.Type()] = map[string]base.Noder{LATEST_VERSION: n, "1": n}
	return nil
}

func ExtractVarSelectorToVarMapping(
	graph_config map[string]any,
	config map[string]any,
) map[string][]string {
	/*Extracts references variable selectors from node configuration.

	  The `config` parameter represents the configuration for a specific node type and corresponds
	  to the `data` field in the node definition object.

	  The returned mapping has the following structure:

	      {'1747829548239.#1747829667553.result#': ['1747829667553', 'result']}

	  For loop and iteration nodes, the mapping may look like this:

	      {
	          "1748332301644.input_selector": ["1748332363630", "result"],
	          "1748332325079.1748332325079.#sys.workflow_id#": ["sys", "workflow_id"],
	      }

	  where `1748332301644` is the ID of the loop / iteration node,
	  and `1748332325079` is the ID of the node inside the loop or iteration node.

	  Here, the key consists of two parts: the current node ID (provided as the `node_id`
	  parameter to `_extract_variable_selector_to_variable_mapping`) and the variable selector,
	  enclosed in `#` symbols. These two parts are separated by a dot (`.`).

	  The value is a list of string representing the variable selector, where the first element is the node ID
	  of the referenced variable, and the second element is the variable name within that node.

	  The meaning of the above response is:

	  The node with ID `1747829548239` references the variable `result` from the node with
	  ID `1747829667553`. For example, if `1747829548239` is a LLM node, its prompt may contain a
	  reference to the `result` output variable of node `1747829667553`.

	  :param graph_config: graph config
	  :param config: node config
	  :return:
	*/
	node_id, exist := mapstruct.GetFromMap[string](config, "id")
	if !exist || node_id == "" {
		panic(exceptions.NewValueError("Node ID is required when extracting variable selector to variable mapping."))
	}
	nodeType := mapstruct.Get(mapstruct.Get(config, "data", map[string]any{}), "type", "")
	if _, ok := NODE_TYPE_CLASSES_MAPPING[nodesenumtypes.NodeType(nodeType)]; !ok {
		panic(exceptions.NewValueError(fmt.Sprintf("nodeType %s is not regist.", nodeType)))
	}
	cls := NODE_TYPE_CLASSES_MAPPING[nodesenumtypes.NodeType(nodeType)][LATEST_VERSION]
	if cls == nil {
		panic(exceptions.NewValueError(fmt.Sprintf("nodeType %s's class is nil.", nodeType)))
	}
	// Pass raw dict data instead of creating NodeData instance
	return cls.ExtractVarSelectorToVarMapping(
		graph_config, node_id, mapstruct.Get(config, "data", map[string]any{}),
	)
}
