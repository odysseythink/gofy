package variables

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
	"mlib.com/mlog"
)

type Variabler interface {
	GetID() string
	GetName() string
	GetDescription() string
	GetSelector() []string

	SetID(string)
	SetName(string)
	SetDescription(string)
	SetSelector([]string)

	GetValue() any
	SetValue(any)
	Text() string
	Log() string
	Markdown() string
	Size() int
	ToObject() any
	ValueType() variableenumtypes.VariableType

	ToDict(Variabler) map[string]any
}
type BaseVariable[T bool | string | struct{} | float64 | int | map[string]any | *file.File | []any | []string | []float64 | []int | []map[string]any | []*file.File | []Variabler] struct {
	ModelConfig map[string]any `json:"model_config"`
	Value       T              `json:"value"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Selector    []string       `json:"selector"`
}

func (variable *BaseVariable[T]) ToDict(ver Variabler) map[string]any {
	bindata, err := json.Marshal(variable)
	if err != nil {
		mlog.Errorf("json marshal variable=%#v failed:%v", variable, err)
		panic(exceptions.NewValueError(fmt.Sprintf("json marshal variable=%#v failed:%v", variable, err)))
	}
	tmpmap := map[string]any{}
	err = json.Unmarshal(bindata, &tmpmap)
	if err != nil {
		mlog.Errorf("json Unmarshal variable=%#v failed:%v", variable, err)
		panic(exceptions.NewValueError(fmt.Sprintf("json Unmarshal variable=%#v failed:%v", variable, err)))
	}
	tmpmap["value_type"] = ver.ValueType()
	mlog.Debugf("------variable dict=%#v", tmpmap)

	return tmpmap
}

func (variable *BaseVariable[T]) GetID() string             { return variable.ID }
func (variable *BaseVariable[T]) GetName() string           { return variable.Name }
func (variable *BaseVariable[T]) GetDescription() string    { return variable.Description }
func (variable *BaseVariable[T]) GetSelector() []string     { return variable.Selector }
func (variable *BaseVariable[T]) SetID(val string)          { variable.ID = val }
func (variable *BaseVariable[T]) SetName(val string)        { variable.Name = val }
func (variable *BaseVariable[T]) SetDescription(val string) { variable.Description = val }
func (variable *BaseVariable[T]) SetSelector(val []string)  { variable.Selector = val }

func (variable *BaseVariable[T]) GetValue() any {
	return variable.Value
}

func (variable *BaseVariable[T]) SetValue(val any) {
	if reflect.TypeOf(variable.Value) != reflect.TypeOf(val) {
		panic(exceptions.NewValueError(fmt.Sprintf("when set Segment, target val=%#v must be the type of %#v", val, variable.Value)))
	}
	variable.Value = val.(T)
}

func (variable *BaseVariable[T]) Text() string {
	switch v := any(variable.Value).(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%v", variable.Value)
	case map[string]any:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	case []any:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	case []map[string]any:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	case []int:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	case []float64:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	case []string:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	case []*file.File:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	default:
		bindata, err := json.Marshal(variable.Value)
		if err != nil {
			return fmt.Sprintf("%v", variable.Value)
		}
		return string(bindata)
	}
}

func (variable *BaseVariable[T]) Log() string {
	bindata, err := json.MarshalIndent(variable.Value, "", " ")
	if err != nil {
		return fmt.Sprintf("%v", variable.Value)
	}
	return string(bindata)
}

func (variable *BaseVariable[T]) Markdown() string {
	return variable.Log()
}

// Size returns the size of the Value in bytes.
func (variable *BaseVariable[T]) Size() int {
	switch real_value := any(variable.Value).(type) {
	case bool:
		return 1
	case string:
		return len(real_value)
	case struct{}:
		return 0
	case float64:
		return 8
	case int:
		return 8
	case map[string]any:
		return len(real_value)
	case *file.File:
		return 1
	case []any:
		return len(real_value)
	case []string:
		return len(real_value)
	case []float64:
		return len(real_value)
	case []int:
		return len(real_value)
	case []map[string]any:
		return len(real_value)
	case []*file.File:
		return len(real_value)
	case []Variabler:
		return len(real_value)
	}
	return reflect.ValueOf(variable.Value).Len()
}

// ToObject returns the Value.
func (variable *BaseVariable[T]) ToObject() any {
	return variable.Value
}
