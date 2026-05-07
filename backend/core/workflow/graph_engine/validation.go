package graphengine

import (
	"fmt"
	"slices"

	"github.com/odysseythink/gofy/backend/core/workflow/graph"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
)

// GraphValidator validates graph structure before execution.
type GraphValidator struct{}

// getNodeType extracts the node type string from a node config map.
func getNodeType(nodeConfig map[string]any) string {
	if nodeConfig == nil {
		return ""
	}
	data, ok := nodeConfig["data"].(map[string]any)
	if !ok {
		return ""
	}
	t, ok := data["type"].(string)
	if !ok {
		return ""
	}
	return t
}

// ValidateGraph checks the graph for structural issues.
func (v *GraphValidator) ValidateGraph(g *graph.Graph) []error {
	var errs []error

	if g == nil {
		return []error{fmt.Errorf("graph is nil")}
	}

	if len(g.NodeIDs) == 0 {
		errs = append(errs, fmt.Errorf("graph has no nodes"))
		return errs
	}

	// Check for start node
	startNodes := v.getNodeIDsByType(g, string(nodesenumtypes.Node_START))
	if len(startNodes) == 0 {
		errs = append(errs, fmt.Errorf("graph must have at least one start node"))
	}
	if len(startNodes) > 1 {
		errs = append(errs, fmt.Errorf("graph must have exactly one start node, found %d", len(startNodes)))
	}

	// Check for end node (or answer node)
	endNodes := v.getNodeIDsByType(g, string(nodesenumtypes.Node_END))
	answerNodes := v.getNodeIDsByType(g, string(nodesenumtypes.Node_ANSWER))
	if len(endNodes) == 0 && len(answerNodes) == 0 {
		errs = append(errs, fmt.Errorf("graph must have at least one end or answer node"))
	}

	// Check for orphan nodes (no incoming edges, except start/root)
	for _, nodeID := range g.NodeIDs {
		nodeType := getNodeType(g.NodeIDConfigMapping[nodeID])
		if nodeType == string(nodesenumtypes.Node_START) || nodeID == g.RootNodeID {
			continue
		}
		incoming := g.ReverseEdgeMapping[nodeID]
		if len(incoming) == 0 {
			errs = append(errs, fmt.Errorf("node %s has no incoming edges (orphan)", nodeID))
		}
	}

	// Check for cycles (DFS-based)
	if v.detectCycle(g) {
		errs = append(errs, fmt.Errorf("graph contains a cycle"))
	}

	// Check edge validity: all edge source/target must reference existing nodes
	for sourceID, edges := range g.EdgeMapping {
		if !slices.Contains(g.NodeIDs, sourceID) {
			errs = append(errs, fmt.Errorf("edge references non-existent source node: %s", sourceID))
		}
		for _, edge := range edges {
			if !slices.Contains(g.NodeIDs, edge.TargetNodeID) {
				errs = append(errs, fmt.Errorf("edge references non-existent target node: %s", edge.TargetNodeID))
			}
		}
	}

	return errs
}

// ValidateNodeConnections validates that each node has valid connections
// based on its type.
func (v *GraphValidator) ValidateNodeConnections(g *graph.Graph) []error {
	var errs []error

	for _, nodeID := range g.NodeIDs {
		nodeType := getNodeType(g.NodeIDConfigMapping[nodeID])

		switch nodeType {
		case string(nodesenumtypes.Node_IF_ELSE):
			// Must have at least 2 outgoing edges (true/false branches)
			outgoing := g.EdgeMapping[nodeID]
			if len(outgoing) < 2 {
				errs = append(errs, fmt.Errorf("if_else node %s must have at least 2 outgoing edges", nodeID))
			}
		case string(nodesenumtypes.Node_QUESTION_CLASSIFIER):
			// Must have at least 2 outgoing edges (one per class)
			outgoing := g.EdgeMapping[nodeID]
			if len(outgoing) < 2 {
				errs = append(errs, fmt.Errorf("question_classifier node %s must have at least 2 outgoing edges", nodeID))
			}
		case string(nodesenumtypes.Node_ITERATION):
			// Must have body nodes connected
			outgoing := g.EdgeMapping[nodeID]
			if len(outgoing) == 0 {
				errs = append(errs, fmt.Errorf("iteration node %s must have body nodes", nodeID))
			}
		case string(nodesenumtypes.Node_LOOP):
			// Must have body nodes connected
			outgoing := g.EdgeMapping[nodeID]
			if len(outgoing) == 0 {
				errs = append(errs, fmt.Errorf("loop node %s must have body nodes", nodeID))
			}
		}
	}

	return errs
}

// getNodeIDsByType returns all node IDs matching the given type string.
func (v *GraphValidator) getNodeIDsByType(g *graph.Graph, nodeType string) []string {
	var result []string
	for _, nodeID := range g.NodeIDs {
		if getNodeType(g.NodeIDConfigMapping[nodeID]) == nodeType {
			result = append(result, nodeID)
		}
	}
	return result
}

// detectCycle uses DFS to detect cycles in the graph.
func (v *GraphValidator) detectCycle(g *graph.Graph) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, nodeID := range g.NodeIDs {
		if !visited[nodeID] {
			if v.dfsDetectCycle(g, nodeID, visited, recStack) {
				return true
			}
		}
	}
	return false
}

func (v *GraphValidator) dfsDetectCycle(g *graph.Graph, nodeID string, visited, recStack map[string]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	for _, edge := range g.EdgeMapping[nodeID] {
		neighbor := edge.TargetNodeID
		if !visited[neighbor] {
			if v.dfsDetectCycle(g, neighbor, visited, recStack) {
				return true
			}
		} else if recStack[neighbor] {
			// Allow cycles in iteration/loop nodes
			neighborType := getNodeType(g.NodeIDConfigMapping[neighbor])
			if neighborType != string(nodesenumtypes.Node_ITERATION) && neighborType != string(nodesenumtypes.Node_LOOP) {
				return true
			}
		}
	}

	recStack[nodeID] = false
	return false
}
