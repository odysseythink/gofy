package base

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	basenodesexceptions "github.com/odysseythink/gofy/backend/core/exceptions/base_nodes"
	"github.com/odysseythink/gofy/backend/core/file"
	basenodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/base_nodes"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/mlog"
)

type NumberType interface {
	~int | ~float32 | float64
}

// DefaultValue represents a default value with type and key
type DefaultValue struct {
	Value any                                 `json:"value"`
	Type  basenodesenumtypes.DefaultValueType `json:"type"`
	Key   string                              `json:"key"`
}

// parseJSON is a unified JSON parsing handler
func parseJSON(value string) (any, error) {
	var result any
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil, basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Invalid JSON format for value: %v", value))
	}
	return result, nil
}

// validateArray is a unified array type validation
func validateArray(value any, elementType basenodesenumtypes.DefaultValueType) bool {
	_, ok := value.([]any)
	if !ok {
		return false
	}
	switch elementType {
	case basenodesenumtypes.DefaultValue_STRING:
		if _, ok := value.(string); ok {
			return true
		} else {
			return false
		}
	case basenodesenumtypes.DefaultValue_NUMBER:
		if _, ok := value.(int); ok {
			return true
		} else if _, ok := value.(float64); ok {
			return true
		} else {
			return false
		}
	case basenodesenumtypes.DefaultValue_OBJECT:
		if _, ok := value.(map[string]any); ok {
			return true
		} else {
			return false
		}
	case basenodesenumtypes.DefaultValue_ARRAY_NUMBER:
		if _, ok := value.([]int); ok {
			return true
		} else if _, ok := value.([]float64); ok {
			return true
		} else {
			return false
		}
	case basenodesenumtypes.DefaultValue_ARRAY_STRING:
		if _, ok := value.([]string); ok {
			return true
		} else {
			return false
		}
	case basenodesenumtypes.DefaultValue_ARRAY_OBJECT:
		if _, ok := value.([]map[string]any); ok {
			return true
		} else {
			return false
		}
	case basenodesenumtypes.DefaultValue_ARRAY_FILES:
		if _, ok := value.([]*file.File); ok {
			return true
		} else {
			return false
		}
	}
	return false
}

// convertNumber is a unified number conversion handler
func convertNumber(value string) (float64, error) {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, basenodesexceptions.NewDefaultValueTypeError(fmt.Sprintf("Cannot convert to number: %v", value))
	}
	return num, nil
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
	StartNodeID string `json:"start_node_id,omitempty"`
}

// BaseIterationState represents base iteration state
type BaseIterationState struct {
	IterationNodeID string         `json:"iteration_node_id"`
	Index           int            `json:"index"`
	Inputs          map[string]any `json:"inputs"`
	Metadata        any            `json:"metadata"`
}
