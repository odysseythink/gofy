package jinja2

import (
	codeexecutor "github.com/odysseythink/gofy/backend/core/helper/code_executor"
	codeexecutorenumtypes "github.com/odysseythink/gofy/backend/enum_types/code_executor"
)

type Jinja2Formatter struct {
}

func (formatter *Jinja2Formatter) Format(template string, inputs map[string]string) string {
	new_input := map[string]any{}
	for k, v := range inputs {
		new_input[k] = v
	}
	result := codeexecutor.ExecuteWorkflowCodeTemplate(codeexecutorenumtypes.CodeLanguage_JINJA2, template, new_input)
	str_result := ""
	if _, ok := result["result"]; ok {
		if _, ok := result["result"].(string); ok {
			str_result = result["result"].(string)
		}
	}
	return str_result
}
