package variable

import (
	"fmt"
	"maps"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
)

type VariableType string

const (
	Variable_NUMBER  VariableType = "number"
	Variable_STRING  VariableType = "string"
	Variable_BOOLEAN VariableType = "boolean"
	Variable_OBJECT  VariableType = "object"
	Variable_SECRET  VariableType = "secret"
	Variable_FILE    VariableType = "file"

	Variable_ARRAY_ANY    VariableType = "array[any]"
	Variable_ARRAY_STRING VariableType = "array[string]"
	Variable_ARRAY_NUMBER VariableType = "array[number]"
	Variable_ARRAY_OBJECT VariableType = "array[object]"
	Variable_ARRAY_FILE   VariableType = "array[file]"

	Variable_NONE  VariableType = "none"
	Variable_GROUP VariableType = "group"
)

func (s VariableType) Valid() bool {
	return s == Variable_NUMBER ||
		s == Variable_STRING ||
		s == Variable_OBJECT ||
		s == Variable_SECRET ||
		s == Variable_FILE ||
		s == Variable_ARRAY_ANY ||
		s == Variable_ARRAY_STRING ||
		s == Variable_ARRAY_NUMBER ||
		s == Variable_ARRAY_OBJECT ||
		s == Variable_ARRAY_FILE ||
		s == Variable_NONE ||
		s == Variable_GROUP
}

var (
	ENVIRONMENT_VARIABLE_SUPPORTED_TYPES = []VariableType{Variable_STRING, Variable_NUMBER, Variable_SECRET}
)

// ArrayValidation represents the strategy for validating array elements
type ArrayValidationType string

const (
	// ArrayValidationNone skips element validation (only check array container)
	ArrayValidation_NONE ArrayValidationType = "none"

	// ArrayValidationFirst validates the first element (if array is non-empty)
	ArrayValidation_FIRST ArrayValidationType = "first"

	// ArrayValidationAll validates all elements in the array
	ArrayValidation_ALL ArrayValidationType = "all"
)

// SegmentType represents the type of a segment
type SegmentType string

const (
	Segment_NUMBER       SegmentType = "number"
	Segment_INTEGER      SegmentType = "integer"
	Segment_FLOAT        SegmentType = "float"
	Segment_STRING       SegmentType = "string"
	Segment_OBJECT       SegmentType = "object"
	Segment_SECRET       SegmentType = "secret"
	Segment_FILE         SegmentType = "file"
	Segment_ARRAY_ANY    SegmentType = "array[any]"
	Segment_ARRAY_STRING SegmentType = "array[string]"
	Segment_ARRAY_NUMBER SegmentType = "array[number]"
	Segment_ARRAY_OBJECT SegmentType = "array[object]"
	Segment_ARRAY_FILE   SegmentType = "array[file]"
	Segment_NONE         SegmentType = "none"
	Segment_GROUP        SegmentType = "group"
)

// IsArrayType checks if the segment type is an array type
func (s SegmentType) IsArrayType() bool {
	return slices.Contains(_ARRAY_TYPES, s)
}
func (s SegmentType) Validate() bool {
	return s != Segment_NUMBER &&
		s != Segment_INTEGER &&
		s != Segment_FLOAT &&
		s != Segment_STRING &&
		s != Segment_OBJECT &&
		s != Segment_SECRET &&
		s != Segment_FILE &&
		s != Segment_ARRAY_ANY &&
		s != Segment_ARRAY_STRING &&
		s != Segment_ARRAY_NUMBER &&
		s != Segment_ARRAY_OBJECT &&
		s != Segment_ARRAY_FILE &&
		s != Segment_NONE &&
		s != Segment_GROUP
}

// InferSegmentType attempts to infer the SegmentType based on the Go type of the value parameter.
// Returns nil if no appropriate SegmentType can be determined for the given value.
func InferSegmentType(value any) SegmentType {
	if value == nil {
		return SegmentType("")
	}

	switch v := value.(type) {
	case []any:
		if len(v) == 0 {
			return Segment_ARRAY_ANY
		}
		elemTypes := make(map[SegmentType]bool)
		for _, item := range v {
			segmentType := InferSegmentType(item)
			if string(segmentType) == "" {
				return SegmentType("")
			}
			elemTypes[segmentType] = true
		}

		if len(elemTypes) != 1 {
			// Check if all types are numerical
			if !slices.ContainsFunc(slices.Sorted(maps.Keys(elemTypes)), func(segType SegmentType) bool {
				return !slices.Contains(_NUMERICAL_TYPES, segType)
			}) {
				return Segment_ARRAY_NUMBER
			} else {
				return Segment_ARRAY_ANY
			}
		} else {
			switch v[0].(type) {
			case string:
				return Segment_ARRAY_STRING
			case int, int32, int64, int16, int8, uint, uint32, uint64, uint16, uint8, float32, float64:
				return Segment_ARRAY_NUMBER
			case map[string]any:
				return Segment_ARRAY_OBJECT
			default:
				// This should be unreachable
				panic(exceptions.NewValueError(fmt.Sprintf("not supported value %v", value)))
			}
		}
	case int, int32, int64, int16, int8, uint, uint32, uint64, uint16, uint8, float32, float64:
		return Segment_NUMBER
	case string:
		return Segment_STRING
	case map[string]any:
		return Segment_OBJECT
	default:
		return SegmentType("")
	}
}

// validateArray validates an array value based on the segment type and validation strategy
func (s SegmentType) validateArray(value any, arrayValidation ArrayValidationType) bool {
	arr, ok := value.([]any)
	if !ok {
		return false
	}

	// Skip element validation if array is empty
	if len(arr) == 0 {
		return true
	}

	if s == Segment_ARRAY_ANY {
		return true
	}

	elementType, exists := _ARRAY_ELEMENT_TYPES_MAPPING[s]
	if !exists {
		return false
	}

	switch arrayValidation {
	case ArrayValidation_NONE:
		return true
	case ArrayValidation_FIRST:
		if len(arr) > 0 {
			return elementType.IsValid(arr[0], ArrayValidation_NONE)
		}
		return true
	case ArrayValidation_ALL:
		for _, item := range arr {
			if !elementType.IsValid(item, ArrayValidation_NONE) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// IsValid checks if a value matches the segment type
func (s SegmentType) IsValid(value any, arrayValidation ArrayValidationType) bool {
	if s.IsArrayType() {
		return s.validateArray(value, arrayValidation)
	} else if s == Segment_NUMBER {
		switch value.(type) {
		case int, int32, int64, int16, int8, uint, uint32, uint64, uint16, uint8, float32, float64:
			return true
		default:
			return false
		}
	} else if s == Segment_INTEGER {
		switch value.(type) {
		case int, int32, int64, int16, int8, uint, uint32, uint64, uint16, uint8:
			return true
		default:
			return false
		}
	} else if s == Segment_FLOAT {
		switch value.(type) {
		case float32, float64:
			return true
		default:
			return false
		}
	} else if s == Segment_STRING {
		switch value.(type) {
		case string:
			return true
		default:
			return false
		}
	} else if s == Segment_OBJECT {
		switch value.(type) {
		case map[string]any:
			return true
		default:
			return false
		}
	} else if s == Segment_SECRET {
		switch value.(type) {
		case string:
			return true
		default:
			return false
		}
	} else if s == Segment_NONE {
		return value == nil
	} else {
		panic(exceptions.NewValueError("this statement should be unreachable."))
	}
}

// ExposedType returns the type exposed to the frontend.
// The frontend treats INTEGER and FLOAT as NUMBER, so these are returned as NUMBER here.
func (s SegmentType) ExposedType() SegmentType {
	if s == Segment_INTEGER || s == Segment_FLOAT {
		return Segment_NUMBER
	}
	return s
}

var (
	// _ARRAY_ELEMENT_TYPES_MAPPING maps array segment types to their element types
	_ARRAY_ELEMENT_TYPES_MAPPING = map[SegmentType]SegmentType{
		// ARRAYANY does not have corresponding element type
		Segment_ARRAY_STRING: Segment_STRING,
		Segment_ARRAY_NUMBER: Segment_NUMBER,
		Segment_ARRAY_OBJECT: Segment_OBJECT,
		Segment_ARRAY_FILE:   Segment_FILE,
	}

	// _ARRAY_TYPES contains all array segment types
	_ARRAY_TYPES = []SegmentType{
		Segment_ARRAY_STRING,
		Segment_ARRAY_NUMBER,
		Segment_ARRAY_OBJECT,
		Segment_ARRAY_FILE,
		Segment_ARRAY_ANY,
	}

	// _NUMERICAL_TYPES contains all numerical segment types
	_NUMERICAL_TYPES = []SegmentType{
		Segment_NUMBER,
		Segment_INTEGER,
		Segment_FLOAT,
	}
)
