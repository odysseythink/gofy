package tool

import (
	"mlib.com/gofy/server/core/exceptions"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	"mlib.com/mlog"
)

type ToolEntity struct {
	ProviderID         string         `json:"provider_id"`
	ProviderType       string         `json:"provider_type"` /*["builtin", "api", "workflow"]*/
	ProviderName       string         `json:"provider_name"` // redundancy
	ToolName           string         `json:"tool_name"`
	ToolLabel          string         `json:"tool_label"` // redundancy
	ToolConfigurations map[string]any `json:"tool_configurations"`
}

func (entity *ToolEntity) ValidateToolConfigurations() {
	for _, v := range entity.ToolConfigurations {
		switch v.(type) {
		case string:
		case int:
		case float32:
		case float64:
		case bool:
			panic(exceptions.NewValueError("value must be a string, int, float, or bool"))
		}
	}

}

// return value
type ToolInput struct {
	// TODO: check this type
	Value any    `json:"value"` //[Any, list[str]]
	Type  string `json:"type"`  //["mixed", "variable", "constant"]
}

func (ti *ToolInput) CheckType() {
	if ti.Type == "mixed" {
		if _, ok := ti.Value.(string); !ok {
			panic(exceptions.NewValueError("value must be a string"))
		}
	} else if ti.Type == "variable" {
		if _, ok := ti.Value.([]any); ok {
			for _, sv := range ti.Value.([]any) {
				if _, ok := sv.(string); !ok {
					panic(exceptions.NewValueError("value must be a list of string"))
				}
			}
		} else if _, ok := ti.Value.([]string); ok {
			mlog.Debugf("-----%#v", ti.Value.([]string))
		} else {
			panic(exceptions.NewValueError("value must be a list of string"))
		}

	} else if ti.Type == "constant" {
		switch ti.Value.(type) {
		case string:
		case int:
		case float32:
		case float64:
		case bool:
			panic(exceptions.NewValueError("value must be a string, int, float, or bool"))
		}
	}

}

// ToolNodeData represents answer node data
type ToolNodeData struct {
	*basenodesentities.BaseNodeData
	ToolParameters map[string]ToolInput `json:"tool_parameters"`
}
