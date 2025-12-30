package variables

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
	"mlib.com/gofy/server/utils/mapstruct"
	"mlib.com/mlog"
)

// Segmenter is the interface for all segment types
type Segmenter interface {
	// GetValue returns the underlying value
	GetValue() any
	// SetValue sets the underlying value
	SetValue(any)
	// Text returns the text representation
	Text() string
	// Log returns the log representation
	Log() string
	// Markdown returns the markdown representation
	Markdown() string
	// Size returns the size in bytes
	Size() int
	// ToObject returns the value as an object
	ToObject() any
	// ValueType returns the segment type
	ValueType() variableenumtypes.SegmentType
}

// Segment is the base struct for all segment implementations
type Segment[T string | int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64 | float32 | float64 | struct{} | map[string]any | []any | []string | []map[string]any | []int | []int16 | []int32 | []int64 | []uint | []uint16 | []uint32 | []uint64 | []float32 | []float64] struct {
	value T
}

// 序列化：MarshalJSON 实现 json.Marshaler 接口
func (ct Segment[T]) MarshalJSON() ([]byte, error) {
	tmp := map[string]any{"value": ct.value}
	if v, ok := any(&ct).(Segmenter); ok {
		tmp["value_type"] = v.ValueType()
	}
	return json.Marshal(tmp)
}

// 反序列化：UnmarshalJSON 实现 json.Unmarshaler 接口
func (ct *Segment[T]) UnmarshalJSON(data []byte) error {
	tmp := map[string]any{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if _, ok := tmp["value"]; !ok {
		return fmt.Errorf("missing value field")
	}
	if _, ok := tmp["value_type"]; ok {
		if _, ok := tmp["value_type"].(string); !ok {
			return fmt.Errorf("value_type field must be a string")
		} else {
			if !variableenumtypes.SegmentType(tmp["value_type"].(string)).IsValid(tmp["value"], variableenumtypes.ArrayValidation_ALL) {
				return fmt.Errorf("value %#v must be segment type %v", tmp["value"], variableenumtypes.SegmentType(tmp["value_type"].(string)))
			} else {

			}
		}
	} else {
		if _, ok := tmp["value_type"].(T); !ok {
			return fmt.Errorf("value %#v must be type of(%#v)", tmp["value"], ct.value)
		} else {
			ct.value = tmp["value"].(T)
		}
		return nil
	}

	return nil
}

// GetValue returns the underlying value
func (s *Segment[T]) GetValue() any {
	return s.value
}

// SetValue sets the underlying value
func (s *Segment[T]) SetValue(val any) {
	if _, ok := val.(T); !ok {
		return
	}
	s.value = val.(T)
}

// Text returns the text representation
func (s *Segment[T]) Text() string {
	if any(s.value) == nil {
		return ""
	}
	return fmt.Sprintf("%v", s.value)
}

// Log returns the log representation
func (s *Segment[T]) Log() string {
	if any(s.value) == nil {
		return ""
	}
	return fmt.Sprintf("%v", s.value)
}

// Markdown returns the markdown representation
func (s *Segment[T]) Markdown() string {
	if any(s.value) == nil {
		return ""
	}
	return fmt.Sprintf("%v", s.value)
}

// ToObject returns the value as an object
func (s *Segment[T]) ToObject() any {
	return s.value
}

// NoneSegment represents a null/empty segment
type NoneSegment struct {
	*Segment[struct{}]
}

// NewNoneSegment creates a new none segment
func NewNoneSegment() *NoneSegment {
	return &NoneSegment{
		Segment: &Segment[struct{}]{},
	}
}
func (s *NoneSegment) Size() int {
	return 0
}

// ValueType returns the segment type
func (s *NoneSegment) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_NONE
}

// StringSegment represents a string segment
type StringSegment struct {
	*Segment[string]
}

// NewStringSegment creates a new string segment
func NewStringSegment(value string) *StringSegment {
	return &StringSegment{
		Segment: &Segment[string]{value: value},
	}
}

// ValueType returns the segment type
func (s *StringSegment) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_STRING
}
func (s *StringSegment) Size() int {
	return len(s.value)
}

// FloatSegment represents a float segment
type FloatSegment[T float32 | float64] struct {
	*Segment[T]
}

// NewFloatSegment creates a new float segment
func NewFloatSegment[T float32 | float64](value T) *FloatSegment[T] {
	return &FloatSegment[T]{
		Segment: &Segment[T]{value},
	}
}

// ValueType returns the segment type
func (s *FloatSegment[T]) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_FLOAT
}
func (s *FloatSegment[T]) Size() int {
	return len(fmt.Sprintf("%v", s.value))
}

// IntegerSegment represents an integer segment
type IntegerSegment[T int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64] struct {
	*Segment[T]
}

// NewIntegerSegment creates a new integer segment
func NewIntegerSegment[T int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64](value T) *IntegerSegment[T] {
	return &IntegerSegment[T]{
		Segment: &Segment[T]{value},
	}
}

// ValueType returns the segment type
func (s *IntegerSegment[T]) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_INTEGER
}
func (s *IntegerSegment[T]) Size() int {
	return strconv.IntSize
}

// ObjectSegment represents an object segment
type ObjectSegment struct {
	*Segment[map[string]any]
}

// NewObjectSegment creates a new object segment
func NewObjectSegment(value map[string]any) *ObjectSegment {
	return &ObjectSegment{
		Segment: &Segment[map[string]any]{value},
	}
}
func (s *ObjectSegment) Size() int {
	return len(s.value)
}

// Text returns JSON representation for object segment
func (s *ObjectSegment) Text() string {
	data, err := json.Marshal(s.value)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// Log returns formatted JSON representation for object segment
func (s *ObjectSegment) Log() string {
	data, err := json.MarshalIndent(s.value, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

// Markdown returns formatted JSON representation for object segment
func (s *ObjectSegment) Markdown() string {
	return s.Log()
}

// ValueType returns the segment type
func (s *ObjectSegment) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_OBJECT
}

// ArraySegment is the base for array segments
type ArraySegment[T []any | []map[string]any | []string | []int | []int16 | []int32 | []int64 | []uint | []uint16 | []uint32 | []uint64 | []float32 | []float64] struct {
	*Segment[T]
}

// NewArraySegment creates a new array segment
func NewArraySegment[T []any | []map[string]any | []string | []int | []int16 | []int32 | []int64 | []uint | []uint16 | []uint32 | []uint64 | []float32 | []float64](value T) *ArraySegment[T] {
	return &ArraySegment[T]{
		Segment: &Segment[T]{value: value},
	}
}
func (s *ArraySegment[T]) Size() int {
	return len(s.value)
}

// ArrayAnySegment represents an array of any type segment
type ArrayAnySegment struct {
	*ArraySegment[[]any]
}

// NewArrayAnySegment creates a new array any segment
func NewArrayAnySegment(value []any) *ArrayAnySegment {
	return &ArrayAnySegment{
		ArraySegment: NewArraySegment(value),
	}
}
func (s *ArrayAnySegment) Markdown() string {
	items := []string{}
	for _, item := range s.value {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return strings.Join(items, "\n")
}

// ValueType returns the segment type
func (s *ArrayAnySegment) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_ARRAY_ANY
}

// ArrayStringSegment represents an array of strings segment
type ArrayStringSegment struct {
	*ArraySegment[[]string]
}

// NewArrayStringSegment creates a new array string segment
func NewArrayStringSegment(value []string) *ArrayStringSegment {
	return &ArrayStringSegment{
		ArraySegment: NewArraySegment(value),
	}
}
func (s *ArrayStringSegment) Markdown() string {
	return strings.Join(s.value, "\n")
}

// Text returns JSON representation for array string segment
func (s *ArrayStringSegment) Text() string {
	data, err := json.Marshal(s.value)
	if err != nil {
		return "[]"
	}
	return string(data)
}

// ValueType returns the segment type
func (s *ArrayStringSegment) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_ARRAY_STRING
}

// ArrayNumberSegment represents an array of numbers segment
type ArrayNumberSegment[T []int | []int16 | []int32 | []int64 | []uint | []uint16 | []uint32 | []uint64 | []float32 | []float64] struct {
	*ArraySegment[T]
}

// NewArrayNumberSegment creates a new array number segment
func NewArrayNumberSegment[T []int | []int16 | []int32 | []int64 | []uint | []uint16 | []uint32 | []uint64 | []float32 | []float64](value T) *ArrayNumberSegment[T] {
	return &ArrayNumberSegment[T]{
		ArraySegment: NewArraySegment(value),
	}
}

// ValueType returns the segment type
func (s *ArrayNumberSegment[T]) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_ARRAY_NUMBER
}

// ArrayObjectSegment represents an array of objects segment
type ArrayObjectSegment struct {
	*ArraySegment[[]map[string]any]
}

// NewArrayObjectSegment creates a new array object segment
func NewArrayObjectSegment(value []map[string]any) *ArrayObjectSegment {
	return &ArrayObjectSegment{
		ArraySegment: NewArraySegment(value),
	}
}

// ValueType returns the segment type
func (s *ArrayObjectSegment) ValueType() variableenumtypes.SegmentType {
	return variableenumtypes.Segment_ARRAY_OBJECT
}

// GetSegmentDiscriminator returns the segment type discriminator
func GetSegmentDiscriminator(v any) variableenumtypes.SegmentType {
	if segment, ok := v.(Segmenter); ok {
		segmentType := segment.ValueType()
		return segmentType
	} else if m, ok := v.(map[string]any); ok {
		if v, exist := mapstruct.GetFromMap[string](m, "value_type"); !exist {
			return variableenumtypes.SegmentType("")
		} else {
			segmentType := variableenumtypes.SegmentType(v)
			if segmentType.Validate() {
				return segmentType
			} else {
				return variableenumtypes.SegmentType("")
			}
		}
	} else {
		return variableenumtypes.SegmentType("")
	}
}

func BuildSegment(value any, segmentType variableenumtypes.SegmentType) Segmenter {
	if !segmentType.IsValid(value, variableenumtypes.ArrayValidation_ALL) {
		mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
		return nil
	}
	switch segmentType {
	case variableenumtypes.Segment_NUMBER:
		switch realVal := value.(type) {
		case int:
			return NewIntegerSegment(realVal)
		case int16:
			return NewIntegerSegment(realVal)
		case int32:
			return NewIntegerSegment(realVal)
		case int64:
			return NewIntegerSegment(realVal)
		case uint:
			return NewIntegerSegment(realVal)
		case uint16:
			return NewIntegerSegment(realVal)
		case uint32:
			return NewIntegerSegment(realVal)
		case uint64:
			return NewIntegerSegment(realVal)
		case float32:
			return NewFloatSegment(realVal)
		case float64:
			return NewFloatSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_INTEGER:
		switch realVal := value.(type) {
		case int:
			return NewIntegerSegment(realVal)
		case int16:
			return NewIntegerSegment(realVal)
		case int32:
			return NewIntegerSegment(realVal)
		case int64:
			return NewIntegerSegment(realVal)
		case uint:
			return NewIntegerSegment(realVal)
		case uint16:
			return NewIntegerSegment(realVal)
		case uint32:
			return NewIntegerSegment(realVal)
		case uint64:
			return NewIntegerSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_FLOAT:
		switch realVal := value.(type) {
		case float32:
			return NewFloatSegment(realVal)
		case float64:
			return NewFloatSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_STRING:
		switch realVal := value.(type) {
		case string:
			return NewStringSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_OBJECT:
		switch realVal := value.(type) {
		case map[string]any:
			return NewObjectSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_SECRET:
		switch realVal := value.(type) {
		case string:
			return NewStringSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_ARRAY_ANY:
		switch realVal := value.(type) {
		case []any:
			return NewArrayAnySegment(realVal)
		case []string:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []map[string]any:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []int:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []int16:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []int32:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []int64:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []uint:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []uint16:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []uint32:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []uint64:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []float32:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		case []float64:
			tmp := []any{}
			for _, item := range realVal {
				tmp = append(tmp, item)
			}
			return NewArrayAnySegment(tmp)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_ARRAY_STRING:
		switch realVal := value.(type) {
		case []any:
			tmp := []string{}
			for _, item := range realVal {
				if _, ok := item.(string); !ok {
					mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
					return nil
				}
				tmp = append(tmp, item.(string))
			}
			return NewArrayStringSegment(tmp)
		case []string:
			return NewArrayStringSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_ARRAY_NUMBER:
		switch realVal := value.(type) {
		case []any:
			if len(realVal) == 0 {
				mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
				return nil
			}
			switch realVal[0].(type) {
			case int:
				tmp := []int{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, realItem)
					case int16:
						tmp = append(tmp, int(realItem))
					case int32:
						tmp = append(tmp, int(realItem))
					case int64:
						tmp = append(tmp, int(realItem))
					case uint:
						tmp = append(tmp, int(realItem))
					case uint16:
						tmp = append(tmp, int(realItem))
					case uint32:
						tmp = append(tmp, int(realItem))
					case uint64:
						tmp = append(tmp, int(realItem))
					case float32:
						tmp = append(tmp, int(realItem))
					case float64:
						tmp = append(tmp, int(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}

				}
				return NewArrayNumberSegment(tmp)
			case int16:
				tmp := []int16{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, int16(realItem))
					case int16:
						tmp = append(tmp, realItem)
					case int32:
						tmp = append(tmp, int16(realItem))
					case int64:
						tmp = append(tmp, int16(realItem))
					case uint:
						tmp = append(tmp, int16(realItem))
					case uint16:
						tmp = append(tmp, int16(realItem))
					case uint32:
						tmp = append(tmp, int16(realItem))
					case uint64:
						tmp = append(tmp, int16(realItem))
					case float32:
						tmp = append(tmp, int16(realItem))
					case float64:
						tmp = append(tmp, int16(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case int32:
				tmp := []int32{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, int32(realItem))
					case int16:
						tmp = append(tmp, int32(realItem))
					case int32:
						tmp = append(tmp, realItem)
					case int64:
						tmp = append(tmp, int32(realItem))
					case uint:
						tmp = append(tmp, int32(realItem))
					case uint16:
						tmp = append(tmp, int32(realItem))
					case uint32:
						tmp = append(tmp, int32(realItem))
					case uint64:
						tmp = append(tmp, int32(realItem))
					case float32:
						tmp = append(tmp, int32(realItem))
					case float64:
						tmp = append(tmp, int32(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case int64:
				tmp := []int64{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, int64(realItem))
					case int16:
						tmp = append(tmp, int64(realItem))
					case int32:
						tmp = append(tmp, int64(realItem))
					case int64:
						tmp = append(tmp, realItem)
					case uint:
						tmp = append(tmp, int64(realItem))
					case uint16:
						tmp = append(tmp, int64(realItem))
					case uint32:
						tmp = append(tmp, int64(realItem))
					case uint64:
						tmp = append(tmp, int64(realItem))
					case float32:
						tmp = append(tmp, int64(realItem))
					case float64:
						tmp = append(tmp, int64(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case uint:
				tmp := []uint{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, uint(realItem))
					case int16:
						tmp = append(tmp, uint(realItem))
					case int32:
						tmp = append(tmp, uint(realItem))
					case int64:
						tmp = append(tmp, uint(realItem))
					case uint:
						tmp = append(tmp, realItem)
					case uint16:
						tmp = append(tmp, uint(realItem))
					case uint32:
						tmp = append(tmp, uint(realItem))
					case uint64:
						tmp = append(tmp, uint(realItem))
					case float32:
						tmp = append(tmp, uint(realItem))
					case float64:
						tmp = append(tmp, uint(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case uint16:
				tmp := []uint16{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, uint16(realItem))
					case int16:
						tmp = append(tmp, uint16(realItem))
					case int32:
						tmp = append(tmp, uint16(realItem))
					case int64:
						tmp = append(tmp, uint16(realItem))
					case uint:
						tmp = append(tmp, uint16(realItem))
					case uint16:
						tmp = append(tmp, realItem)
					case uint32:
						tmp = append(tmp, uint16(realItem))
					case uint64:
						tmp = append(tmp, uint16(realItem))
					case float32:
						tmp = append(tmp, uint16(realItem))
					case float64:
						tmp = append(tmp, uint16(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case uint32:
				tmp := []uint32{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, uint32(realItem))
					case int16:
						tmp = append(tmp, uint32(realItem))
					case int32:
						tmp = append(tmp, uint32(realItem))
					case int64:
						tmp = append(tmp, uint32(realItem))
					case uint:
						tmp = append(tmp, uint32(realItem))
					case uint16:
						tmp = append(tmp, uint32(realItem))
					case uint32:
						tmp = append(tmp, realItem)
					case uint64:
						tmp = append(tmp, uint32(realItem))
					case float32:
						tmp = append(tmp, uint32(realItem))
					case float64:
						tmp = append(tmp, uint32(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case uint64:
				tmp := []uint64{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, uint64(realItem))
					case int16:
						tmp = append(tmp, uint64(realItem))
					case int32:
						tmp = append(tmp, uint64(realItem))
					case int64:
						tmp = append(tmp, uint64(realItem))
					case uint:
						tmp = append(tmp, uint64(realItem))
					case uint16:
						tmp = append(tmp, uint64(realItem))
					case uint32:
						tmp = append(tmp, uint64(realItem))
					case uint64:
						tmp = append(tmp, realItem)
					case float32:
						tmp = append(tmp, uint64(realItem))
					case float64:
						tmp = append(tmp, uint64(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case float32:
				tmp := []float32{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, float32(realItem))
					case int16:
						tmp = append(tmp, float32(realItem))
					case int32:
						tmp = append(tmp, float32(realItem))
					case int64:
						tmp = append(tmp, float32(realItem))
					case uint:
						tmp = append(tmp, float32(realItem))
					case uint16:
						tmp = append(tmp, float32(realItem))
					case uint32:
						tmp = append(tmp, float32(realItem))
					case uint64:
						tmp = append(tmp, float32(realItem))
					case float32:
						tmp = append(tmp, realItem)
					case float64:
						tmp = append(tmp, float32(realItem))
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			case float64:
				tmp := []float64{}
				for _, item := range realVal {
					switch realItem := item.(type) {
					case int:
						tmp = append(tmp, float64(realItem))
					case int16:
						tmp = append(tmp, float64(realItem))
					case int32:
						tmp = append(tmp, float64(realItem))
					case int64:
						tmp = append(tmp, float64(realItem))
					case uint:
						tmp = append(tmp, float64(realItem))
					case uint16:
						tmp = append(tmp, float64(realItem))
					case uint32:
						tmp = append(tmp, float64(realItem))
					case uint64:
						tmp = append(tmp, float64(realItem))
					case float32:
						tmp = append(tmp, float64(realItem))
					case float64:
						tmp = append(tmp, realItem)
					default:
						mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
						return nil
					}
				}
				return NewArrayNumberSegment(tmp)
			default:
				mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
				return nil
			}
		case []int:
			return NewArrayNumberSegment(realVal)
		case []int16:
			return NewArrayNumberSegment(realVal)
		case []int32:
			return NewArrayNumberSegment(realVal)
		case []int64:
			return NewArrayNumberSegment(realVal)
		case []uint:
			return NewArrayNumberSegment(realVal)
		case []uint16:
			return NewArrayNumberSegment(realVal)
		case []uint32:
			return NewArrayNumberSegment(realVal)
		case []uint64:
			return NewArrayNumberSegment(realVal)
		case []float32:
			return NewArrayNumberSegment(realVal)
		case []float64:
			return NewArrayNumberSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_ARRAY_OBJECT:
		switch realVal := value.(type) {
		case []any:
			tmp := []map[string]any{}
			for _, item := range realVal {
				if _, ok := item.(map[string]any); !ok {
					mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
					return nil
				}
				tmp = append(tmp, item.(map[string]any))
			}
			return NewArrayObjectSegment(tmp)
		case []map[string]any:
			return NewArrayObjectSegment(realVal)
		default:
			mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
			return nil
		}
	case variableenumtypes.Segment_NONE:
		return NewNoneSegment()
	default:
		mlog.Error("value {%#v} is not match segment type(%s)", value, segmentType)
		return nil
	}
}
