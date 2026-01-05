package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"mlib.com/gofy/server/core/exceptions"
	basenodesexceptions "mlib.com/gofy/server/core/exceptions/base_nodes"
	basenodesenumtypes "mlib.com/gofy/server/enum_types/base_nodes"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/mlog"
)

type NumberType interface {
	~int | ~float32 | float64
}

type DefaultValue struct {
	Key   string                              `json:"key"`
	Value any                                 `json:"value"`
	Type  basenodesenumtypes.DefaultValueType `json:"type"`
}

// parseJSON parses a JSON string into a Go value.
func parseJSON(value string) (any, error) {
	var result any
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Invalid JSON format for value: %s", value))
	}
	return result, nil
}

// validateArray checks if the value is a slice and each element matches the expected element type.
func validateArray(value any, elementType basenodesenumtypes.DefaultValueType) bool {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		switch elementType {
		case basenodesenumtypes.DefaultValue_STRING:
			if _, ok := elem.(string); !ok {
				return false
			}
		case basenodesenumtypes.DefaultValue_NUMBER:
			if !isNumberType(elem) {
				return false
			}
		case basenodesenumtypes.DefaultValue_OBJECT:
			if _, ok := elem.(map[string]any); !ok {
				return false
			}
		default:
			// unknown element type
			return false
		}
	}
	return true
}

// isNumberType checks if the value is int, float32, or float64.
func isNumberType(value any) bool {
	switch value.(type) {
	case int, float32, float64:
		return true
	default:
		return false
	}
}

// convertNumber converts a string to float64.
func convertNumber(value string) (float64, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Cannot convert to number: %s", value))
	}
	return f, nil
}

// Validate validates the DefaultValue according to its type.
func (dv *DefaultValue) Validate() error {
	if dv.Type == "" {
		return basenodesexceptions.NewDefaultValueTypeError("type field is required")
	}

	// Special case for ARRAY_FILES - skip validation
	if dv.Type == basenodesenumtypes.DefaultValue_ARRAY_FILES {
		return nil
	}

	// Define validation rules for each type
	switch dv.Type {
	case basenodesenumtypes.DefaultValue_STRING:
		// Value must be string
		if _, ok := dv.Value.(string); !ok {
			return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Value must be string type for %v", dv.Value))
		}

	case basenodesenumtypes.DefaultValue_NUMBER:
		// If value is string, convert it
		if strVal, ok := dv.Value.(string); ok {
			num, err := convertNumber(strVal)
			if err != nil {
				return err
			}
			dv.Value = num
		}
		if !isNumberType(dv.Value) {
			return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Value must be number type for %v", dv.Value))
		}

	case basenodesenumtypes.DefaultValue_OBJECT:
		// If value is string, parse as JSON
		if strVal, ok := dv.Value.(string); ok {
			obj, err := parseJSON(strVal)
			if err != nil {
				return err
			}
			dv.Value = obj
		}
		if _, ok := dv.Value.(map[string]any); !ok {
			return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Value must be object type for %v", dv.Value))
		}

	case basenodesenumtypes.DefaultValue_ARRAY_NUMBER:
		// If value is string, parse as JSON
		if strVal, ok := dv.Value.(string); ok {
			arr, err := parseJSON(strVal)
			if err != nil {
				return err
			}
			dv.Value = arr
		}
		if !validateArray(dv.Value, basenodesenumtypes.DefaultValue_NUMBER) {
			return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("All elements must be number for %v", dv.Value))
		}

	case basenodesenumtypes.DefaultValue_ARRAY_STRING:
		if strVal, ok := dv.Value.(string); ok {
			arr, err := parseJSON(strVal)
			if err != nil {
				return err
			}
			dv.Value = arr
		}
		if !validateArray(dv.Value, basenodesenumtypes.DefaultValue_STRING) {
			return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("All elements must be string for %v", dv.Value))
		}

	case basenodesenumtypes.DefaultValue_ARRAY_OBJECT:
		if strVal, ok := dv.Value.(string); ok {
			arr, err := parseJSON(strVal)
			if err != nil {
				return err
			}
			dv.Value = arr
		}
		if !validateArray(dv.Value, basenodesenumtypes.DefaultValue_OBJECT) {
			return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("All elements must be object for %v", dv.Value))
		}

	default:
		return basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Unsupported type: %s", dv.Type))
	}

	return nil
}

// RetryConfig represents node retry configuration
type RetryConfig struct {
	MaxRetries    int  `json:"max_retries"`    // max retry times
	RetryInterval int  `json:"retry_interval"` // retry interval in milliseconds
	RetryEnabled  bool `json:"retry_enabled"`  // whether retry is enabled
}

// RetryIntervalSeconds returns the retry interval in seconds
func (rc *RetryConfig) RetryIntervalSeconds() float64 {
	return float64(rc.RetryInterval) / 1000
}

type NodeDataer interface {
	Marshal(map[string]any, any) error
	ModelDump() map[string]any
}

// BaseNodeData represents base node data
type BaseNodeData struct {
	Title         string                       `json:"title"`
	Desc          string                       `json:"desc,omitempty"`
	ErrorStrategy nodesenumtypes.ErrorStrategy `json:"error_strategy,omitempty"`
	DefaultValue  []*DefaultValue              `json:"default_value,omitempty"`
	Version       string                       `json:"version"`
	RetryConfig   RetryConfig                  `json:"retry_config"`
}

func (data *BaseNodeData) ModelDump() map[string]any {
	ret := map[string]any{}
	bindata, _ := json.Marshal(data)
	err := json.Unmarshal(bindata, &ret)
	if err != nil {
		mlog.Errorf("json unmarshal(%s) failed:%v", string(bindata), err)
	}
	return ret
}

func (data *BaseNodeData) Marshal(config map[string]any, dest any) error {
	if config == nil {
		return exceptions.NewValueError("ivnalid config")
	}
	mlog.Debugf("------config=%#v", config)
	bindata, _ := json.Marshal(config)
	err := json.Unmarshal(bindata, dest)
	if err != nil {
		return err
	}
	mlog.Debugf("------dest=%#v", dest)
	return nil
}

// DefaultValueDict returns a dictionary of default values
func (bnd *BaseNodeData) DefaultValueDict() map[string]any {
	result := make(map[string]any)
	for _, dv := range bnd.DefaultValue {
		result[dv.Key] = dv.Value
	}
	return result
}

// BaseIterationNodeData represents base iteration node data
type BaseIterationNodeData struct {
	*BaseNodeData
	StartNodeID string `json:"start_node_id"`
}

// BaseIterationState represents base iteration state
type BaseIterationState struct {
	IterationNodeID string         `json:"iteration_node_id"`
	Index           int            `json:"index"`
	Inputs          map[string]any `json:"inputs"`
	Metadata        any            `json:"metadata"`
}

type BaseLoopNodeData struct {
	*BaseNodeData
	StartNodeID string `json:"start_node_id"`
}

type BaseLoopState struct {
	LoopNodeID string         `json:"loop_node_id"`
	Index      int            `json:"index"`
	Inputs     map[string]any `json:"inputs"`
}
