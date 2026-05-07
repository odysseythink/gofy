package variablefactory

import (
	"encoding/json"
	"fmt"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/constants"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/file"
	"github.com/odysseythink/gofy/backend/core/variables"
	variableenumtypes "github.com/odysseythink/gofy/backend/enum_types/variable"
	"github.com/odysseythink/mlog"
)

type InvalidSelectorError struct {
	desc string
}

func NewInvalidSelectorError(desc string) *InvalidSelectorError {
	return &InvalidSelectorError{
		desc: desc,
	}
}
func (e *InvalidSelectorError) Error() string {
	return fmt.Sprintf("invalid selector:%s", e.desc)
}

type UnsupportedSegmentTypeError struct {
	desc string
}

func NewUnsupportedSegmentTypeError(desc string) *UnsupportedSegmentTypeError {
	return &UnsupportedSegmentTypeError{
		desc: desc,
	}
}
func (e *UnsupportedSegmentTypeError) Error() string {
	return fmt.Sprintf("unsupported segment type:%s", e.desc)
}

// BuildConversationVariableFromMapping builds a conversation variable from a mapping
func BuildConversationVariableFromMapping(mapping map[string]any) variables.Variabler {
	if _, ok := mapping["name"]; !ok {
		panic(variables.NewVariableError("missing name"))
	}
	return _build_variable_from_mapping(mapping, []string{constants.CONVERSATION_VARIABLE_NODE_ID, mapping["name"].(string)})
}

// buildEnvironmentVariableFromMapping builds an environment variable from a mapping
func BuildEnvironmentVariableFromMapping(mapping map[string]any) variables.Variabler {
	if _, ok := mapping["name"]; !ok {
		panic(variables.NewVariableError("missing name"))
	}
	return _build_variable_from_mapping(mapping, []string{constants.ENVIRONMENT_VARIABLE_NODE_ID, mapping["name"].(string)})
}

// _build_variable_from_mapping is a factory function to create variables
func _build_variable_from_mapping(mapping map[string]any, selector []string) variables.Variabler {
	mlog.Debugf("------mapping=%#v", mapping)
	if _, ok := mapping["value_type"]; !ok || mapping["value_type"] == nil {
		mlog.Error("missing value type")
		panic(variables.NewVariableError("missing value type"))
	}
	if _, ok := mapping["value_type"].(string); !ok || mapping["value_type"].(string) == "" {
		mlog.Error("value type must be string")
		panic(variables.NewVariableError("value type must be string"))
	}
	valueType := mapping["value_type"].(string)

	value, ok := mapping["value"]
	if !ok {
		mlog.Error("value type must be string")
		panic(variables.NewVariableError("missing value"))
	}

	switch variableenumtypes.VariableType(valueType) {
	case variableenumtypes.Variable_STRING:
		result := new(variables.BaseVariable[string])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}

		return &variables.StringVariable{BaseVariable: result}
	case variableenumtypes.Variable_BOOLEAN:
		result := new(variables.BaseVariable[bool])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}

		return &variables.BooleanVariable{BaseVariable: result}
	case variableenumtypes.Variable_NUMBER:
		if _, ok := value.(int); ok {
			result := new(variables.BaseVariable[int])
			bindata, _ := json.Marshal(mapping)
			err := json.Unmarshal(bindata, result)
			if err != nil {
				mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
				panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
			}
			if result.Selector == nil {
				result.Selector = selector
			}
			if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
				panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
			}
			return &variables.IntegerVariable{BaseVariable: result}
		} else if _, ok := value.(float64); ok {
			result := new(variables.BaseVariable[float64])
			bindata, _ := json.Marshal(mapping)
			err := json.Unmarshal(bindata, result)
			if err != nil {
				mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
				panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
			}
			if result.Selector == nil {
				result.Selector = selector
			}
			if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
				panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
			}
			return &variables.FloatVariable{BaseVariable: result}
		} else {
			mlog.Error("number value must be int or float64")
			panic(variables.NewVariableError("number value must be int or float64"))
		}
	case variableenumtypes.Variable_OBJECT:
		result := new(variables.BaseVariable[map[string]any])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.ObjectVariable{BaseVariable: result}
	case variableenumtypes.Variable_SECRET:
		result := new(variables.BaseVariable[string])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.SecretVariable{BaseVariable: result}
	case variableenumtypes.Variable_ARRAY_ANY:
		result := new(variables.BaseVariable[[]any])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.ArrayAnyVariable{BaseVariable: result}
	case variableenumtypes.Variable_ARRAY_STRING:
		result := new(variables.BaseVariable[[]string])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.ArrayStringVariable{BaseVariable: result}
	case variableenumtypes.Variable_ARRAY_NUMBER:
		if _, ok := value.([]int); ok {
			result := new(variables.BaseVariable[[]int])
			bindata, _ := json.Marshal(mapping)
			err := json.Unmarshal(bindata, result)
			if err != nil {
				mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
				panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
			}
			if result.Selector == nil {
				result.Selector = selector
			}
			if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
				panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
			}
			return &variables.ArrayIntegerVariable{BaseVariable: result}
		} else if _, ok := value.([]float64); ok {
			result := new(variables.BaseVariable[[]float64])
			bindata, _ := json.Marshal(mapping)
			err := json.Unmarshal(bindata, result)
			if err != nil {
				mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
				panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
			}
			if result.Selector == nil {
				result.Selector = selector
			}
			if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
				panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
			}
			return &variables.ArrayFloatVariable{BaseVariable: result}
		} else {
			mlog.Error("array number value must be []int or []float64")
			panic(variables.NewVariableError("array number value must be []int or []float64"))
		}
	case variableenumtypes.Variable_ARRAY_BOOLEAN:
		result := new(variables.BaseVariable[[]bool])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.ArrayBooleanVariable{BaseVariable: result}
	case variableenumtypes.Variable_ARRAY_OBJECT:
		result := new(variables.BaseVariable[[]map[string]any])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.ArrayObjectVariable{BaseVariable: result}
	case variableenumtypes.Variable_ARRAY_FILE:
		result := new(variables.BaseVariable[[]*file.File])
		bindata, _ := json.Marshal(mapping)
		err := json.Unmarshal(bindata, result)
		if err != nil {
			mlog.Errorf("mapping(%#v) json.Unmarshal to Variable failed:%v", mapping, err)
			panic(variables.NewVariableError("json.Unmarshal to Variable failed:" + err.Error()))
		}
		if result.Selector == nil {
			result.Selector = selector
		}
		if result.Size() > confy.GetWithDefault[int]("workflow.max_variable_size", 204800) {
			panic(variables.NewVariableError(fmt.Sprintf("variable size %d exceeds limit %d", result.Size(), confy.GetWithDefault[int]("workflow.max_variable_size", 204800))))
		}
		return &variables.ArrayFileVariable{BaseVariable: result}
	}

	panic(variables.NewVariableError(fmt.Sprintf("not supported value type %s", valueType)))
}

// BuildSegment builds a segment from a value
func BuildSegment(value any) variables.Variabler {
	if value == nil {
		return &variables.NoneVariable{
			BaseVariable: &variables.BaseVariable[struct{}]{},
		}
	}

	switch v := value.(type) {
	case bool:
		return &variables.BooleanVariable{
			BaseVariable: &variables.BaseVariable[bool]{
				Value: v,
			},
		}
	case string:
		return &variables.StringVariable{
			BaseVariable: &variables.BaseVariable[string]{
				Value: v,
			},
		}
	case int:
		return &variables.IntegerVariable{
			BaseVariable: &variables.BaseVariable[int]{
				Value: v,
			},
		}
	case float64:
		return &variables.FloatVariable{
			BaseVariable: &variables.BaseVariable[float64]{
				Value: v,
			},
		}
	case map[string]any:
		return &variables.ObjectVariable{
			BaseVariable: &variables.BaseVariable[map[string]any]{
				Value: v,
			},
		}
	case map[string]string:
		newobject := map[string]any{}
		for sk, sv := range v {
			newobject[sk] = sv
		}
		return &variables.ObjectVariable{
			BaseVariable: &variables.BaseVariable[map[string]any]{
				Value: newobject,
			},
		}
	case []any:
		return &variables.ArrayAnyVariable{
			BaseVariable: &variables.BaseVariable[[]any]{
				Value: v,
			},
		}
	case []map[string]any:
		return &variables.ArrayObjectVariable{
			BaseVariable: &variables.BaseVariable[[]map[string]any]{
				Value: v,
			},
		}
	case []int:
		return &variables.ArrayIntegerVariable{
			BaseVariable: &variables.BaseVariable[[]int]{
				Value: v,
			},
		}
	case []float64:
		return &variables.ArrayFloatVariable{
			BaseVariable: &variables.BaseVariable[[]float64]{
				Value: v,
			},
		}
	case []bool:
		return &variables.ArrayBooleanVariable{
			BaseVariable: &variables.BaseVariable[[]bool]{
				Value: v,
			},
		}
	case []string:
		return &variables.ArrayStringVariable{
			BaseVariable: &variables.BaseVariable[[]string]{
				Value: v,
			},
		}
	case []*file.File:
		return &variables.ArrayFileVariable{
			BaseVariable: &variables.BaseVariable[[]*file.File]{
				Value: v,
			},
		}
	}
	panic(exceptions.NewValueError(fmt.Sprintf("not supported value %#v", value)))
}
