package variables

import (
	"encoding/json"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/file"
	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
)

type BooleanVariable struct {
	*BaseVariable[bool]
}

func (variable *BooleanVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_BOOLEAN
}

type StringVariable struct {
	*BaseVariable[string]
}

func (variable *StringVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_STRING
}

type FloatVariable struct {
	*BaseVariable[float64]
}

func (variable *FloatVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_NUMBER
}

type IntegerVariable struct {
	*BaseVariable[int]
}

func (variable *IntegerVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_NUMBER
}

type ObjectVariable struct {
	*BaseVariable[map[string]any]
}

func (variable *ObjectVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_OBJECT
}

func (variable *ObjectVariable) Text() string {
	bindata, err := json.Marshal(variable.Value)
	if err != nil {
		mlog.Errorf("json marshal(%#v) to string failed:%v", variable.Value, err)
		return "{}"
	}
	return string(bindata)
}

func (variable *ObjectVariable) Log() string {
	bindata, err := json.MarshalIndent(variable.Value, "", " ")
	if err != nil {
		mlog.Errorf("json MarshalIndent(%#v) to string failed:%v", variable.Value, err)
		return "{}"
	}
	return string(bindata)
}

func (variable *ObjectVariable) Markdown() string {
	return variable.Log()
}

type ArrayAnyVariable struct {
	*BaseVariable[[]any]
}

func (variable *ArrayAnyVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_ANY
}

type ArrayStringVariable struct {
	*BaseVariable[[]string]
}

func (variable *ArrayStringVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_STRING
}

type ArrayIntegerVariable struct {
	*BaseVariable[[]int]
}

func (variable *ArrayIntegerVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_NUMBER
}

type ArrayFloatVariable struct {
	*BaseVariable[[]float64]
}

func (variable *ArrayFloatVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_NUMBER
}

type ArrayObjectVariable struct {
	*BaseVariable[[]map[string]any]
}

func (variable *ArrayObjectVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_OBJECT
}

type NoneVariable struct {
	*BaseVariable[struct{}]
}

func (variable *NoneVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_NONE
}

func (variable *NoneVariable) Text() string {
	return ""
}

func (variable *NoneVariable) Log() string {
	return ""
}

func (variable *NoneVariable) Markdown() string {
	return ""
}

type FileVariable struct {
	*BaseVariable[*file.File]
}

func (variable *FileVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_FILE
}

func (variable *FileVariable) Markdown() string {
	return variable.Value.Markdown()
}

func (variable *FileVariable) Log() string {
	return ""
}

func (variable *FileVariable) Text() string {
	return ""
}

type ArrayFileVariable struct {
	*BaseVariable[[]*file.File]
}

func (variable *ArrayFileVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_FILE
}

type SecretVariable struct {
	*BaseVariable[string]
}

func (variable *SecretVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_SECRET
}

func (variable *SecretVariable) Text() string {
	return strings.Repeat("*", len(variable.Value))
}

func (variable *SecretVariable) Log() string {
	return strings.Repeat("*", len(variable.Value))
}

func (variable *SecretVariable) Markdown() string {
	return strings.Repeat("*", len(variable.Value))
}

type ArrayBooleanVariable struct {
	*BaseVariable[[]bool]
}

func (variable *ArrayBooleanVariable) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_ARRAY_BOOLEAN
}

type VariableGroup struct {
	*BaseVariable[[]Variabler]
}

func (vg *VariableGroup) ValueType() variableenumtypes.VariableType {
	return variableenumtypes.Variable_GROUP
}

func (vg *VariableGroup) Text() string {
	var list []string
	for _, v := range vg.Value {
		if list == nil {
			list = make([]string, 0)
		}
		list = append(list, v.Text())
	}
	return strings.Join(list, "")
}
func (vg *VariableGroup) Log() string {
	var list []string
	for _, v := range vg.Value {
		if list == nil {
			list = make([]string, 0)
		}
		list = append(list, v.Log())
	}
	return strings.Join(list, "")
}
func (vg *VariableGroup) Markdown() string {
	var list []string
	for _, v := range vg.Value {
		if list == nil {
			list = make([]string, 0)
		}
		list = append(list, v.Markdown())
	}
	return strings.Join(list, "")
}
func (vg *VariableGroup) ToObject() any {
	var list []any
	for _, v := range vg.Value {
		if list == nil {
			list = make([]any, 0)
		}
		list = append(list, v.ToObject())
	}
	return list
}
