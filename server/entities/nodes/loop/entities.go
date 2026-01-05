package loop

import (
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	basenodesentities "mlib.com/gofy/server/entities/nodes/base"
	condition "mlib.com/gofy/server/entities/workflow/condition"
	varenumtypes "mlib.com/gofy/server/enum_types/variable"
)

var (
	_VALID_VAR_TYPE = []varenumtypes.SegmentType{
		varenumtypes.Segment_STRING,
		varenumtypes.Segment_NUMBER,
		varenumtypes.Segment_OBJECT,
		varenumtypes.Segment_ARRAY_STRING,
		varenumtypes.Segment_ARRAY_NUMBER,
		varenumtypes.Segment_ARRAY_OBJECT,
	}
)

func _is_valid_var_type(segType varenumtypes.SegmentType) (varenumtypes.SegmentType, error) {
	for _, valid := range _VALID_VAR_TYPE {
		if segType == valid {
			return segType, nil
		}
	}
	return "", exceptions.NewValueError(fmt.Sprintf("invalid segment type: %s", segType))
}

// ValueType defines the type of value for a loop variable.
type ValueType string

const (
	ValueTypeVariable ValueType = "variable"
	ValueTypeConstant ValueType = "constant"
)

// LoopVariableData represents loop variable data.
type LoopVariableData struct {
	Label     string                   `json:"label"`
	VarType   varenumtypes.SegmentType `json:"var_type"`
	ValueType ValueType                `json:"value_type"`      //"variable", "constant"
	Value     any                      `json:"value,omitempty"` //Any | list[str]
}

// LoopNodeData represents answer node data
type LoopNodeData struct {
	*basenodesentities.BaseLoopNodeData

	LoopCount       int                           `json:"loop_count"`
	BreakConditions []condition.Condition         `json:"break_conditions,omitempty"`
	LogicalOperator condition.LogicalOperatorType `json:"logical_operator"`
	LoopVariables   []LoopVariableData            `json:"loop_variables,omitempty"`
	Outputs         map[string]any                `json:"outputs,omitempty"`
}

// LoopStartNodeData represents loop start node data.
type LoopStartNodeData struct {
	*basenodesentities.BaseNodeData
}

// LoopEndNodeData represents loop end node data.
type LoopEndNodeData struct {
	*basenodesentities.BaseNodeData
}

// LoopStateMetaData represents metadata for loop state.
type LoopStateMetaData struct {
	LoopLength int `json:"loop_length"`
}

// LoopState represents loop state.
type LoopState struct {
	*basenodesentities.BaseLoopState
	Outputs       []any              `json:"outputs,omitempty"`
	CurrentOutput any                `json:"current_output,omitempty"`
	MetaData      *LoopStateMetaData `json:"metadata,omitempty"`
}

// GetLastOutput returns the last output if exists.
func (s *LoopState) GetLastOutput() any {
	if len(s.Outputs) > 0 {
		return s.Outputs[len(s.Outputs)-1]
	}
	return nil
}

// GetCurrentOutput returns the current output.
func (s *LoopState) GetCurrentOutput() any {
	return s.CurrentOutput
}
