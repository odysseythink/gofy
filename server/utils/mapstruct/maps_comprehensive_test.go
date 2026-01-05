package mapstruct

import (
	"testing"
	"time"
)

// TestGet_EdgeCases tests edge cases and potential issues
func TestGet_EdgeCases(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		emptyMap := map[string]any{}
		result := Get(emptyMap, "key", "default")
		if result != "default" {
			t.Errorf("Expected default value from empty map, got %v", result)
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var nilMap map[string]any
		result := Get(nilMap, "key", "default")
		if result != "default" {
			t.Errorf("Expected default value from nil map, got %v", result)
		}
	})

	t.Run("empty key", func(t *testing.T) {
		testMap := map[string]any{"key": "value"}
		result := Get(testMap, "", "default")
		if result != "default" {
			t.Errorf("Expected default value for empty key, got %v", result)
		}
	})

	t.Run("key with dots but no nesting", func(t *testing.T) {
		testMap := map[string]any{"key.with.dots": "value"}
		result := Get(testMap, "key.with.dots", "default")
		if result != "value" {
			t.Errorf("Expected value for key with dots, got %v", result)
		}
	})

	t.Run("nested map with missing intermediate", func(t *testing.T) {
		testMap := map[string]any{
			"level1": map[string]any{
				"level2": "value",
			},
		}
		result := Get(testMap, "level1.missing.level3", "default")
		if result != "default" {
			t.Errorf("Expected default value for missing intermediate key, got %v", result)
		}
	})

	t.Run("non-map value in nested path", func(t *testing.T) {
		testMap := map[string]any{
			"level1": "not a map",
		}
		result := Get(testMap, "level1.level2", "default")
		if result != "default" {
			t.Errorf("Expected default value when trying to access nested field on non-map, got %v", result)
		}
	})

	t.Run("zero value as default", func(t *testing.T) {
		testMap := map[string]any{}
		result := Get(testMap, "missing", 0)
		if result != 0 {
			t.Errorf("Expected 0 as default, got %v", result)
		}
	})

	t.Run("nil slice as default", func(t *testing.T) {
		testMap := map[string]any{}
		result := Get(testMap, "missing", []int(nil))
		if result != nil {
			t.Errorf("Expected nil as default, got %v", result)
		}
	})
}

// TestGet_TypeConversions tests various type conversions
func TestGet_TypeConversions(t *testing.T) {
	testMap := map[string]any{
		"string_int":    "123",
		"string_float":  "3.14",
		"string_bool":   "true",
		"int_to_float":  42,
		"float_to_int":  3.14,
		"bool_to_int":   true,
		"int_to_string": 123,
	}

	t.Run("string to int", func(t *testing.T) {
		result := Get(testMap, "string_int", 0)
		if result != 123 {
			t.Errorf("Expected 123, got %v", result)
		}
	})

	t.Run("string to float64", func(t *testing.T) {
		result := Get(testMap, "string_float", 0.0)
		if result != 3.14 {
			t.Errorf("Expected 3.14, got %v", result)
		}
	})

	t.Run("string to bool", func(t *testing.T) {
		result := Get(testMap, "string_bool", false)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("int to float64", func(t *testing.T) {
		result := Get(testMap, "int_to_float", 0.0)
		if result != 42.0 {
			t.Errorf("Expected 42.0, got %v", result)
		}
	})

	t.Run("float64 to int", func(t *testing.T) {
		result := Get(testMap, "float_to_int", 0)
		if result != 3 {
			t.Errorf("Expected 3, got %v", result)
		}
	})

	t.Run("bool to int", func(t *testing.T) {
		result := Get(testMap, "bool_to_int", 0)
		if result != 1 {
			t.Errorf("Expected 1, got %v", result)
		}
	})

	t.Run("int to string", func(t *testing.T) {
		result := Get(testMap, "int_to_string", "")
		if result != "123" {
			t.Errorf("Expected '123', got %v", result)
		}
	})
}

// TestGet_AllIntegerTypes tests all integer type variants
func TestGet_AllIntegerTypes(t *testing.T) {
	testMap := map[string]any{
		"val": 42,
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"int", "val", 0, 42},
		{"int8", "val", int8(0), int8(42)},
		{"int16", "val", int16(0), int16(42)},
		{"int32", "val", int32(0), int32(42)},
		{"int64", "val", int64(0), int64(42)},
		{"uint", "val", uint(0), uint(42)},
		{"uint8", "val", uint8(0), uint8(42)},
		{"uint16", "val", uint16(0), uint16(42)},
		{"uint32", "val", uint32(0), uint32(42)},
		{"uint64", "val", uint64(0), uint64(42)},
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

// TestGet_AllFloatTypes tests all float type variants
func TestGet_AllFloatTypes(t *testing.T) {
	testMap := map[string]any{
		"val": 3.14,
	}

	tests := []struct {
		name       string
		key        string
		defaultVal interface{}
		expected   interface{}
	}{
		{"float32", "val", float32(0), float32(3.14)},
		{"float64", "val", 0.0, 3.14},
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

// TestGet_DeepNesting tests deeply nested map access
func TestGet_DeepNesting(t *testing.T) {
	testMap := map[string]any{
		"l1": map[string]any{
			"l2": map[string]any{
				"l3": map[string]any{
					"l4": map[string]any{
						"l5": "deep_value",
					},
				},
			},
		},
	}

	result := Get(testMap, "l1.l2.l3.l4.l5", "")
	if result != "deep_value" {
		t.Errorf("Expected 'deep_value', got %v", result)
	}

	// Test missing deep key
	result = Get(testMap, "l1.l2.l3.l4.missing", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got %v", result)
	}
}

// TestGet_TimeAndDuration tests time.Time and time.Duration types
func TestGet_TimeAndDuration(t *testing.T) {
	now := time.Now()
	duration := 2 * time.Hour

	testMap := map[string]any{
		"time_val":     now,
		"duration_val": duration,
	}

	t.Run("time.Time", func(t *testing.T) {
		result := Get(testMap, "time_val", time.Time{})
		if !result.Equal(now) {
			t.Errorf("Expected %v, got %v", now, result)
		}
	})

	t.Run("time.Duration", func(t *testing.T) {
		result := Get(testMap, "duration_val", time.Duration(0))
		if result != duration {
			t.Errorf("Expected %v, got %v", duration, result)
		}
	})

	t.Run("missing time returns default", func(t *testing.T) {
		defaultTime := time.Time{}
		result := Get(testMap, "missing", defaultTime)
		if !result.Equal(defaultTime) {
			t.Errorf("Expected default time, got %v", result)
		}
	})
}

// TestGet_SliceConversions tests slice type conversions
func TestGet_SliceConversions(t *testing.T) {
	testMap := map[string]any{
		"int_array":    []any{1, 2, 3},
		"string_array": []any{"a", "b", "c"},
		"mixed_array":  []any{1, "two", true},
	}

	t.Run("[]any to []int", func(t *testing.T) {
		result := Get(testMap, "int_array", []int(nil))
		expected := []int{1, 2, 3}
		if !sliceEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("[]any to []string", func(t *testing.T) {
		result := Get(testMap, "string_array", []string(nil))
		expected := []string{"a", "b", "c"}
		if !sliceEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("[]any to []any", func(t *testing.T) {
		result := Get(testMap, "mixed_array", []any(nil))
		expected := []any{1, "two", true}
		if !sliceEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestGet_MapConversions tests map type conversions
func TestGet_MapConversions(t *testing.T) {
	testMap := map[string]any{
		"string_any_map":    map[string]any{"key1": "value1", "key2": 123},
		"string_string_map": map[string]string{"key1": "value1", "key2": "value2"},
	}

	t.Run("map[string]any to map[string]string", func(t *testing.T) {
		result := Get(testMap, "string_any_map", map[string]string(nil))
		// This should fail conversion and return default
		if result != nil {
			t.Errorf("Expected nil (default) for failed conversion, got %v", result)
		}
	})

	t.Run("map[string]string", func(t *testing.T) {
		result := Get(testMap, "string_string_map", map[string]string(nil))
		expected := map[string]string{"key1": "value1", "key2": "value2"}
		if !mapEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}

// TestGet_NestedWithDifferentTypes tests nested access with different value types
func TestGet_NestedWithDifferentTypes(t *testing.T) {
	testMap := map[string]any{
		"config": map[string]any{
			"settings": map[string]any{
				"enabled": true,
				"count":   42,
				"name":    "test",
			},
		},
	}

	t.Run("nested bool", func(t *testing.T) {
		result := Get(testMap, "config.settings.enabled", false)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("nested int", func(t *testing.T) {
		result := Get(testMap, "config.settings.count", 0)
		if result != 42 {
			t.Errorf("Expected 42, got %v", result)
		}
	})

	t.Run("nested string", func(t *testing.T) {
		result := Get(testMap, "config.settings.name", "")
		if result != "test" {
			t.Errorf("Expected 'test', got %v", result)
		}
	})
}
