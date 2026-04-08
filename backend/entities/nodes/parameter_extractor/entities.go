package parameterextractor

import (
	"encoding/json"
	"slices"
	"strings"

	"mlib.com/gofy/server/core/exceptions"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	"mlib.com/gofy/server/entities/nodes/llm"
	"mlib.com/gofy/server/entities/prompt"
	"mlib.com/mlog"
)

type ParameterConfig struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` /*["string", "number", "bool", "select", "array[string]", "array[number]", "array[object]"]*/
	Options     []string `json:"options"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
}

func (pc *ParameterConfig) ValidateName(value string) string {
	if value == "" {
		panic(exceptions.NewValueError("Parameter name is required"))
	}
	if slices.Contains([]string{"__reason", "__is_success"}, value) {
		panic(exceptions.NewValueError("Invalid parameter name, __reason and __is_success are reserved"))
	}
	return value
}

// ParameterExtractorNodeData represents answer node data
type ParameterExtractorNodeData struct {
	*basenodesentities.BaseNodeData
	Model         llm.ModelConfig      `json:"model"`
	Query         []string             `json:"query"`
	Parameters    []ParameterConfig    `json:"parameters"`
	Instruction   string               `json:"instruction"`
	Memory        *prompt.MemoryConfig `json:"memory"`
	ReasoningMode string               `json:"reasoning_mode"` /* Literal["function_call", "prompt"]*/
	Vision        *llm.VisionConfig    `json:"vision"`
}

func New() *ParameterExtractorNodeData {
	return &ParameterExtractorNodeData{
		Vision: llm.NewVisionConfig(),
	}
}

func (data *ParameterExtractorNodeData) GetParameterJsonSchema() map[string]any {
	parameters := map[string]any{"type": "object", "properties": map[string]any{}, "required": []string{}}

	for _, parameter := range data.Parameters {
		parameter_schema := map[string]any{"description": parameter.Description}

		if slices.Contains([]string{"string", "select"}, parameter.Type) {
			parameter_schema["type"] = "string"
		} else if strings.HasPrefix(parameter.Type, "array") {
			parameter_schema["type"] = "array"
			nested_type := parameter.Type[6:]
			parameter_schema["items"] = map[string]any{"type": nested_type}
		} else {
			parameter_schema["type"] = parameter.Type
		}
		if parameter.Type == "select" {
			parameter_schema["enum"] = parameter.Options
		}
		parameters["properties"].(map[string]any)[parameter.Name] = parameter_schema

		if parameter.Required {
			required := parameters["required"].([]string)
			required = append(required, parameter.Name)
			parameters["required"] = required
		}
	}
	return parameters
}
func (data *ParameterExtractorNodeData) Marshal(config map[string]any, dest any) error {
	mlog.Debugf("---config=%#v", config)
	if config == nil {
		return exceptions.NewValueError("invalid config")
	}
	bindata, _ := json.Marshal(config)
	if dest != data {
		err := json.Unmarshal(bindata, dest)
		if err != nil {
			return err
		}
		return nil
	}
	err := json.Unmarshal(bindata, data)
	if err != nil {
		return err
	}
	if data.ReasoningMode != "function_call" && data.ReasoningMode != "prompt" {
		mlog.Errorf("invalid config reasoning_mode=%s field", data.ReasoningMode)
		return exceptions.NewValueError("invalid config reasoning_mode field")
	}

	mlog.Debugf("---data=%#v", data)
	return nil
}
