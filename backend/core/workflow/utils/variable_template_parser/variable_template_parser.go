package variabletemplateparser

import (
	"fmt"
	"regexp"
	"strings"

	workflowentities "mlib.com/gofy/server/entities/workflow"
)

var (
	REGEX            = regexp.MustCompile(`\{\{(#([a-zA-Z0-9_]{1,50}(\.[a-zA-Z_][a-zA-Z0-9_]{0,29}){1,10})#)\}\}`)
	SELECTOR_PATTERN = regexp.MustCompile(`\{\{(#([a-zA-Z0-9_]{1,50}(?:\.[a-zA-Z_][a-zA-Z0-9_]{0,29}){1,10})#)\}\}`)
)

func ExtractSelectorsFromTemplate(template string) []*workflowentities.VariableSelector {
	parts := SELECTOR_PATTERN.Split(template, -1)
	selectors := make([]*workflowentities.VariableSelector, 0)
	for _, part := range parts {
		if part != "" && strings.Contains(part, ".") && strings.HasPrefix(part, "#") && strings.HasSuffix(part, "#") {
			valueSelector := strings.Trim(part, "#")
			valueSelectorParts := strings.Split(valueSelector, ".")
			selectors = append(selectors, &workflowentities.VariableSelector{
				Variable:      part,
				ValueSelector: valueSelectorParts,
			})
		}
	}
	return selectors
}

type VariableTemplateParser struct {
	template     string
	variableKeys []string
}

func NewVariableTemplateParser(template string) *VariableTemplateParser {
	return &VariableTemplateParser{
		template:     template,
		variableKeys: extractVariableKeys(template),
	}
}

func extractVariableKeys(template string) []string {
	matches := REGEX.FindAllStringSubmatch(template, -1)
	firstGroupMatches := make([]string, 0)
	for _, match := range matches {
		if len(match) > 0 {
			match[0] = strings.TrimPrefix(match[0], "{{")
			match[0] = strings.TrimSuffix(match[0], "}}")
			firstGroupMatches = append(firstGroupMatches, match[0])
		}
	}
	uniqueMatches := make(map[string]bool)
	for _, match := range firstGroupMatches {
		uniqueMatches[match] = true
	}
	result := make([]string, 0, len(uniqueMatches))
	for key := range uniqueMatches {
		result = append(result, key)
	}
	return result
}

func (vtp *VariableTemplateParser) ExtractVariableSelectors() []*workflowentities.VariableSelector {
	variableSelectors := make([]*workflowentities.VariableSelector, 0)
	for _, variableKey := range vtp.variableKeys {
		removeHash := strings.Trim(variableKey, "#")
		splitResult := strings.Split(removeHash, ".")
		if len(splitResult) < 2 {
			continue
		}
		variableSelectors = append(variableSelectors, &workflowentities.VariableSelector{
			Variable:      variableKey,
			ValueSelector: splitResult,
		})
	}
	return variableSelectors
}

func (vtp *VariableTemplateParser) Format(inputs map[string]any) string {
	replacer := func(match []string) string {
		key := match[1]
		value, exists := inputs[key]
		if !exists {
			return match[0]
		}
		if value == nil {
			return ""
		}
		strValue := fmt.Sprintf("%v", value)
		return removeTemplateVariables(strValue)
	}
	prompt := REGEX.ReplaceAllStringFunc(vtp.template, func(s string) string {
		match := REGEX.FindStringSubmatch(s)
		if len(match) > 1 {
			return replacer(match)
		}
		return s
	})
	return regexp.MustCompile(`<\|.*?\|>`).ReplaceAllString(prompt, "")
}

func removeTemplateVariables(text string) string {
	return REGEX.ReplaceAllString(text, "{${1}}")
}

// func main() {
// 	template := "Hello, {{#node_id.query.name#}}! Your age is {{#node_id.query.age#}}."
// 	parser := NewVariableTemplateParser(template)

// 	// Extract template variable keys
// 	variableKeys := parser.variableKeys
// 	fmt.Println("Variable Keys:", variableKeys)

// 	// Extract variable selectors
// 	variableSelectors := parser.ExtractVariableSelectors()
// 	fmt.Println("Variable Selectors:", variableSelectors)

// 	// Format the template string
// 	inputs := map[string]any{
// 		"#node_id.query.name#": "John",
// 		"#node_id.query.age#":  25,
// 	}
// 	formattedString := parser.Format(inputs)
// 	fmt.Println("Formatted String:", formattedString)
// }
