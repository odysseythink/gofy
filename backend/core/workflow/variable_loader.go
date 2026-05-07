package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/odysseythink/gofy/backend/core/variables"
	variablefactory "github.com/odysseythink/gofy/backend/factories/variable_factory"
	"github.com/odysseythink/mlog"
)

// VariableLoader loads and initializes variables for a workflow execution.
type VariableLoader struct{}

// LoadEnvironmentVariables parses environment variables from workflow config.
func (vl *VariableLoader) LoadEnvironmentVariables(envVarsJSON string) []variables.Variabler {
	if envVarsJSON == "" {
		return nil
	}

	var envVarsList []map[string]any
	if err := json.Unmarshal([]byte(envVarsJSON), &envVarsList); err != nil {
		mlog.Errorf("failed to parse environment variables: %v", err)
		return nil
	}

	var result []variables.Variabler
	for _, ev := range envVarsList {
		v := variablefactory.BuildEnvironmentVariableFromMapping(ev)
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

// LoadConversationVariables parses conversation variables from workflow config.
func (vl *VariableLoader) LoadConversationVariables(convVarsJSON string) []variables.Variabler {
	if convVarsJSON == "" {
		return nil
	}

	var convVarsList []map[string]any
	if err := json.Unmarshal([]byte(convVarsJSON), &convVarsList); err != nil {
		mlog.Errorf("failed to parse conversation variables: %v", err)
		return nil
	}

	var result []variables.Variabler
	for _, cv := range convVarsList {
		v := variablefactory.BuildConversationVariableFromMapping(cv)
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

// LoadNodeDefaultValues extracts default input values from a node's configuration.
func (vl *VariableLoader) LoadNodeDefaultValues(nodeConfig map[string]any) map[string]any {
	defaults := make(map[string]any)

	inputs, ok := nodeConfig["inputs"].([]any)
	if !ok {
		return defaults
	}

	for _, input := range inputs {
		inputMap, ok := input.(map[string]any)
		if !ok {
			continue
		}
		varName, _ := inputMap["variable"].(string)
		defaultVal := inputMap["default"]
		if varName != "" && defaultVal != nil {
			defaults[varName] = defaultVal
		}
	}
	return defaults
}

// ExtractVariableSelectorsFromGraph collects all variable selectors referenced in a graph.
func (vl *VariableLoader) ExtractVariableSelectorsFromGraph(graphJSON string) ([][]string, error) {
	var graph map[string]any
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return nil, fmt.Errorf("invalid graph JSON: %w", err)
	}

	var selectors [][]string
	nodes, _ := graph["nodes"].([]any)

	for _, n := range nodes {
		node, _ := n.(map[string]any)
		data, _ := node["data"].(map[string]any)
		vl.extractSelectorsFromMap(data, &selectors)
	}

	return selectors, nil
}

// extractSelectorsFromMap recursively extracts {{#...#}} variable selectors.
func (vl *VariableLoader) extractSelectorsFromMap(data map[string]any, selectors *[][]string) {
	for _, v := range data {
		switch val := v.(type) {
		case string:
			// Extract {{#node_id.var_name#}} patterns
			extracted := extractVariableSelectors(val)
			*selectors = append(*selectors, extracted...)
		case map[string]any:
			vl.extractSelectorsFromMap(val, selectors)
		case []any:
			for _, item := range val {
				if m, ok := item.(map[string]any); ok {
					vl.extractSelectorsFromMap(m, selectors)
				}
			}
		}
	}
}

// extractVariableSelectors extracts selector paths from template strings.
func extractVariableSelectors(template string) [][]string {
	var result [][]string
	// Match {{#selector.path#}} pattern
	i := 0
	for i < len(template) {
		start := findSubstring(template[i:], "{{#")
		if start < 0 {
			break
		}
		start += i
		end := findSubstring(template[start+3:], "#}}")
		if end < 0 {
			break
		}
		end += start + 3

		selector := template[start+3 : end]
		parts := splitSelector(selector)
		if len(parts) > 0 {
			result = append(result, parts)
		}
		i = end + 3
	}
	return result
}

func findSubstring(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func splitSelector(selector string) []string {
	var parts []string
	current := ""
	for _, ch := range selector {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
