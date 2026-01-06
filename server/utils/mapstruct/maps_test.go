package mapstruct

import (
	"reflect"
	"testing"
	"time"
)

func TestGet_BasicTypes(t *testing.T) {
	testMap := map[string]any{
		"int_val":     42,
		"int8_val":    int8(8),
		"int16_val":   int16(16),
		"int32_val":   int32(32),
		"int64_val":   int64(64),
		"uint_val":    uint(100),
		"uint8_val":   uint8(8),
		"uint16_val":  uint16(16),
		"uint32_val":  uint32(32),
		"uint64_val":  uint64(64),
		"float32_val": float32(3.14),
		"float64_val": 3.14159,
		"string_val":  "hello",
		"bool_val":    true,
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"int", "int_val", 0, 42},
		{"int8", "int8_val", int8(0), int8(8)},
		{"int16", "int16_val", int16(0), int16(16)},
		{"int32", "int32_val", int32(0), int32(32)},
		{"int64", "int64_val", int64(0), int64(64)},
		{"uint", "uint_val", uint(0), uint(100)},
		{"uint8", "uint8_val", uint8(0), uint8(8)},
		{"uint16", "uint16_val", uint16(0), uint16(16)},
		{"uint32", "uint32_val", uint32(0), uint32(32)},
		{"uint64", "uint64_val", uint64(0), uint64(64)},
		{"float32", "float32_val", float32(0), float32(3.14)},
		{"float64", "float64_val", 0.0, 3.14159},
		{"string", "string_val", "", "hello"},
		{"bool", "bool_val", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_NestedKeys(t *testing.T) {
	testMap := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": "deep_value",
			},
			"simple": "mid_value",
		},
	}

	tests := []struct {
		name       string
		key        string
		defaultVal string
		expected   string
	}{
		{"single level", "level1.simple", "", "mid_value"},
		{"nested levels", "level1.level2.level3", "", "deep_value"},
		{"missing nested key", "level1.missing", "", ""},
		{"missing intermediate", "missing.level2", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_SliceTypes(t *testing.T) {
	testMap := map[string]any{
		"int_slice":    []int{1, 2, 3},
		"string_slice": []string{"a", "b", "c"},
		"bool_slice":   []bool{true, false, true},
		"float_slice":  []float64{1.1, 2.2, 3.3},
		"any_slice":    []any{1, "two", true},
		"map_slice":    []map[string]any{{"key": "value"}},
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"int slice", "int_slice", []int(nil), []int{1, 2, 3}},
		{"string slice", "string_slice", []string(nil), []string{"a", "b", "c"}},
		{"bool slice", "bool_slice", []bool(nil), []bool{true, false, true}},
		{"float slice", "float_slice", []float64(nil), []float64{1.1, 2.2, 3.3}},
		{"any slice", "any_slice", []any(nil), []any{1, "two", true}},
		{"map slice", "map_slice", []map[string]any(nil), []map[string]any{{"key": "value"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			if !sliceEqual(result, tt.expected) {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_MapTypes(t *testing.T) {
	testMap := map[string]any{
		"string_map": map[string]string{"key1": "value1", "key2": "value2"},
		"int_map":    map[string]int{"a": 1, "b": 2},
		"any_map":    map[string]any{"a": 1, "b": "two"},
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"string map", "string_map", map[string]string(nil), map[string]string{"key1": "value1", "key2": "value2"}},
		{"int map", "int_map", map[string]int(nil), map[string]int{"a": 1, "b": 2}},
		{"any map", "any_map", map[string]any(nil), map[string]any{"a": 1, "b": "two"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			t.Logf("testMap[%s]=%#v", tt.key, testMap[tt.key])
			t.Logf("testMap[%s]=%#v", tt.key, result)
			t.Logf("tt.expected=%#v", tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_TimeTypes(t *testing.T) {
	now := time.Now()
	duration := time.Hour

	testMap := map[string]any{
		"time_val":     now,
		"duration_val": duration,
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"time", "time_val", time.Time{}, now},
		{"duration", "duration_val", time.Duration(0), duration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_DefaultValues(t *testing.T) {
	testMap := map[string]any{
		"existing": "value",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"missing key returns default", "missing", "default", "default"},
		{"existing key returns value", "existing", "default", "value"},
		{"missing nested returns default", "missing.nested", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_TypeConversion(t *testing.T) {
	testMap := map[string]any{
		"string_num":  "123",
		"string_bool": "true",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"string to int", "string_num", 0, 123},
		{"string to bool", "string_bool", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(testMap, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Get(%q, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestGet_InvalidNestedAccess(t *testing.T) {
	testMap := map[string]any{
		"not_a_map": "string_value",
	}

	// Try to access nested field on a non-map value
	result := Get(testMap, "not_a_map.nested", "default")
	if result != "default" {
		t.Errorf("Get(%q, %v) = %v, want %v", "not_a_map.nested", "default", result, "default")
	}
}

func TestGet_EmptyKey(t *testing.T) {
	testMap := map[string]any{
		"key": "value",
	}

	// Empty key should return default
	result := Get(testMap, "", "default")
	if result != "default" {
		t.Errorf("Get(%q, %v) = %v, want %v", "", "default", result, "default")
	}
}

func TestGet_NilMap(t *testing.T) {
	var nilMap map[string]any

	result := Get(nilMap, "key", "default")
	if result != "default" {
		t.Errorf("Get(nil, %q, %v) = %v, want %v", "key", "default", result, "default")
	}
}

// Helper functions for comparison
func sliceEqual(a, b interface{}) bool {
	switch a := a.(type) {
	case []int:
		bSlice, ok := b.([]int)
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if a[i] != bSlice[i] {
				return false
			}
		}
		return true
	case []string:
		bSlice, ok := b.([]string)
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if a[i] != bSlice[i] {
				return false
			}
		}
		return true
	case []bool:
		bSlice, ok := b.([]bool)
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if a[i] != bSlice[i] {
				return false
			}
		}
		return true
	case []float64:
		bSlice, ok := b.([]float64)
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if a[i] != bSlice[i] {
				return false
			}
		}
		return true
	case []any:
		bSlice, ok := b.([]any)
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if a[i] != bSlice[i] {
				return false
			}
		}
		return true
	case []map[string]any:
		bSlice, ok := b.([]map[string]any)
		if !ok {
			return false
		}
		if len(a) != len(bSlice) {
			return false
		}
		for i := range a {
			if !reflect.DeepEqual(a[i], bSlice[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
