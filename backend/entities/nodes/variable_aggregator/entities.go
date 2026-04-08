package variableaggregator

import (
	"encoding/json"

	"mlib.com/gofy/server/core/exceptions"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	"mlib.com/mlog"
)

type Group struct {
	OutputType string     `json:"output_type"` /*Literal["string", "number", "object", "array[string]", "array[number]", "array[object]"]*/
	Variables  [][]string `json:"output_variablestype"`
	GroupName  string     `json:"group_name"`
}

type AdvancedSettings struct {
	GroupEnabled bool    `json:"group_enabled"`
	Groups       []Group `json:"groups"`
}

// VariableAggregatorNodeData represents answer node data
type VariableAggregatorNodeData struct {
	*basenodesentities.BaseNodeData
	Type             string            `json:"type"` /*"variable-assigner"*/
	OutputType       string            `json:"output_type"`
	Variables        [][]string        `json:"variables"`
	AdvancedSettings *AdvancedSettings `json:"advanced_settings"`
}

func (data *VariableAggregatorNodeData) Marshal(config map[string]any, dest any) error {
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
	if data.Type != "variable-aggregator" {
		mlog.Errorf("invalid config type=%s field", data.Type)
		return exceptions.NewValueError("invalid config type field")
	}
	if data.AdvancedSettings != nil {
		for _, group := range data.AdvancedSettings.Groups {
			if group.OutputType != "string" &&
				group.OutputType != "number" &&
				group.OutputType != "object" &&
				group.OutputType != "array[string]" &&
				group.OutputType != "array[number]" &&
				group.OutputType != "array[object]" {
				mlog.Errorf("invalid config group output_type(%s) field", group.OutputType)
				return exceptions.NewValueError("invalid config group output_type field")
			}
		}
	}

	mlog.Debugf("---data=%#v", data)
	return nil
}
