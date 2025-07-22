package utils

import (
	"regexp"
)

var (
	REGEX                    = regexp.MustCompile(`\{\{([a-zA-Z_][a-zA-Z0-9_]{0,29}|#histories#|#query#|#context#)\}\}`)
	WITH_VARIABLE_TMPL_REGEX = regexp.MustCompile(`\{\{([a-zA-Z_][a-zA-Z0-9_]{0,29}|#[a-zA-Z0-9_]{1,50}\.[a-zA-Z0-9_\.]{1,100}#|#histories#|#query#|#context#)\}\}`)
)

type PromptTemplateParser struct {
	template         string
	withVariableTmpl bool
	regex            *regexp.Regexp
	VariableKeys     []string
}

func NewPromptTemplateParser(template string, withVariableTmpl bool) *PromptTemplateParser {
	regex := REGEX
	if withVariableTmpl {
		regex = WITH_VARIABLE_TMPL_REGEX
	}
	return &PromptTemplateParser{
		template:         template,
		withVariableTmpl: withVariableTmpl,
		regex:            regex,
		VariableKeys:     extractVariableKeys(template, regex),
	}
}

func extractVariableKeys(template string, regex *regexp.Regexp) []string {
	matches := regex.FindAllStringSubmatch(template, -1)
	firstGroupMatches := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			firstGroupMatches = append(firstGroupMatches, match[1])
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

func (ptp *PromptTemplateParser) Extract() []string {
	return ptp.VariableKeys
}

func (ptp *PromptTemplateParser) Format(inputs map[string]string, removeTemplateVariables bool) string {
	replacer := func(match []string) string {
		key := match[1]
		value, exists := inputs[key]
		if !exists {
			return match[0]
		}
		if removeTemplateVariables {
			return RemoveTemplateVariables(value, ptp.withVariableTmpl)
		}
		return value
	}
	prompt := ptp.regex.ReplaceAllStringFunc(ptp.template, func(s string) string {
		match := ptp.regex.FindStringSubmatch(s)
		if len(match) > 1 {
			return replacer(match)
		}
		return s
	})
	return regexp.MustCompile(`<\|.*?\|>`).ReplaceAllString(prompt, "")
}

func RemoveTemplateVariables(text string, withVariableTmpl bool) string {
	regex := REGEX
	if withVariableTmpl {
		regex = WITH_VARIABLE_TMPL_REGEX
	}
	return regex.ReplaceAllString(text, "{${1}}")
}

// func main() {
// 	template := "Hello, {{user_name}}! Your query is {{#query#}}."
// 	parser := NewPromptTemplateParser(template, false)

// 	// Extract template variable keys
// 	VariableKeys := parser.Extract()
// 	fmt.Println("Variable Keys:", VariableKeys)

// 	// Format the template string
// 	inputs := map[string]string{
// 		"user_name": "John",
// 		"#query#":   "What is the weather today?",
// 	}
// 	formattedString := parser.Format(inputs, true)
// 	fmt.Println("Formatted String:", formattedString)
// }
