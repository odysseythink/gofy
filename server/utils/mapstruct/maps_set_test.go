package mapstruct

import (
	"reflect"
	"testing"
)

// TestSet_SimpleKeys tests setting simple key-value pairs
func TestSet_SimpleKeys(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]any
		key      string
		value    any
		expected map[string]any
	}{
		{
			name:     "set string value",
			initial:  map[string]any{},
			key:      "name",
			value:    "John",
			expected: map[string]any{"name": "John"},
		},
		{
			name:     "set int value",
			initial:  map[string]any{},
			key:      "age",
			value:    30,
			expected: map[string]any{"age": 30},
		},
		{
			name:     "set bool value",
			initial:  map[string]any{},
			key:      "active",
			value:    true,
			expected: map[string]any{"active": true},
		},
		{
			name:     "overwrite existing value",
			initial:  map[string]any{"name": "Old"},
			key:      "name",
			value:    "New",
			expected: map[string]any{"name": "New"},
		},
		{
			name:     "set nil value",
			initial:  map[string]any{},
			key:      "value",
			value:    nil,
			expected: map[string]any{"value": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.initial, tt.key, tt.value)
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("Set(%q, %v) = %v, want %v", tt.key, tt.value, tt.initial, tt.expected)
			}
		})
	}
}

// TestSet_NestedKeys tests setting values in nested maps
func TestSet_NestedKeys(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]any
		key      string
		value    any
		expected map[string]any
	}{
		{
			name:     "create nested map and set value",
			initial:  map[string]any{},
			key:      "user.name",
			value:    "John",
			expected: map[string]any{"user": map[string]any{"name": "John"}},
		},
		{
			name:     "set value in existing nested map",
			initial:  map[string]any{"user": map[string]any{}},
			key:      "user.name",
			value:    "John",
			expected: map[string]any{"user": map[string]any{"name": "John"}},
		},
		{
			name:     "set value in deeply nested map",
			initial:  map[string]any{},
			key:      "level1.level2.level3.value",
			value:    "deep",
			expected: map[string]any{"level1": map[string]any{"level2": map[string]any{"level3": map[string]any{"value": "deep"}}}},
		},
		{
			name:     "overwrite nested value",
			initial:  map[string]any{"user": map[string]any{"name": "Old"}},
			key:      "user.name",
			value:    "New",
			expected: map[string]any{"user": map[string]any{"name": "New"}},
		},
		{
			name:     "set multiple nested values",
			initial:  map[string]any{},
			key:      "config.settings.timeout",
			value:    30,
			expected: map[string]any{"config": map[string]any{"settings": map[string]any{"timeout": 30}}},
		},
		{
			name:     "set multiple nested values",
			initial:  map[string]any{},
			key:      "config.settings.profile",
			value:    map[string]any{"a": 1, "b": "c"},
			expected: map[string]any{"config": map[string]any{"settings": map[string]any{"profile": map[string]any{"a": 1, "b": "c"}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.initial, tt.key, tt.value)
			t.Logf("tt.initial=%#v", tt.initial)
			t.Logf("tt.expected=%#v", tt.expected)
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("Set(%q, %v) = %v, want %v", tt.key, tt.value, tt.initial, tt.expected)
			}
		})
	}
}

// TestSet_SliceIndex tests setting values in slices using bracket notation
func TestSet_SliceIndex(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]any
		key      string
		value    any
		expected map[string]any
	}{
		{
			name:     "set element in []any slice",
			initial:  map[string]any{"items": []any{"a", "b", "c"}},
			key:      "items[1]",
			value:    "new",
			expected: map[string]any{"items": []any{"a", "new", "c"}},
		},
		{
			name:     "set element in []int slice",
			initial:  map[string]any{"numbers": []int{1, 2, 3}},
			key:      "numbers[0]",
			value:    10,
			expected: map[string]any{"numbers": []int{10, 2, 3}},
		},
		{
			name:     "set element in []string slice",
			initial:  map[string]any{"names": []string{"a", "b", "c"}},
			key:      "names[2]",
			value:    "z",
			expected: map[string]any{"names": []string{"a", "b", "z"}},
		},
		{
			name:     "set element in []bool slice",
			initial:  map[string]any{"flags": []bool{true, false, true}},
			key:      "flags[1]",
			value:    true,
			expected: map[string]any{"flags": []bool{true, true, true}},
		},
		{
			name:     "set element in []float64 slice",
			initial:  map[string]any{"values": []float64{1.1, 2.2, 3.3}},
			key:      "values[0]",
			value:    9.9,
			expected: map[string]any{"values": []float64{9.9, 2.2, 3.3}},
		},
		{
			name:     "set element in []map[string]any slice",
			initial:  map[string]any{"items": []map[string]any{{"id": 1}, {"id": 2}}},
			key:      "items[0]",
			value:    map[string]any{"id": 10, "name": "test"},
			expected: map[string]any{"items": []map[string]any{{"id": 10, "name": "test"}, {"id": 2}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.initial, tt.key, tt.value)
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("Set(%q, %v) = %v, want %v", tt.key, tt.value, tt.initial, tt.expected)
			}
		})
	}
}

// TestSet_NestedSliceIndex tests setting values in nested structures with slice indices
func TestSet_NestedSliceIndex(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]any
		key      string
		value    any
		expected map[string]any
	}{
		{
			name: "set nested field in slice element",
			initial: map[string]any{
				"users": []map[string]any{
					{"name": "John", "age": 30},
					{"name": "Jane", "age": 25},
				},
			},
			key:   "users[0].age",
			value: 35,
			expected: map[string]any{
				"users": []map[string]any{
					{"name": "John", "age": 35},
					{"name": "Jane", "age": 25},
				},
			},
		},
		{
			name: "set deeply nested field in slice element",
			initial: map[string]any{
				"data": []map[string]any{
					{
						"info": map[string]any{
							"details": map[string]any{
								"value": "old",
							},
						},
					},
				},
			},
			key:   "data[0].info.details.value",
			value: "new",
			expected: map[string]any{
				"data": []map[string]any{
					{
						"info": map[string]any{
							"details": map[string]any{
								"value": "new",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.initial, tt.key, tt.value)
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("Set(%q, %v) = %v, want %v", tt.key, tt.value, tt.initial, tt.expected)
			}
		})
	}
}

// TestSet_ErrorCases tests error conditions
func TestSet_ErrorCases(t *testing.T) {
	t.Run("nil destination map", func(t *testing.T) {
		var nilMap map[string]any
		Set(nilMap, "key", "value")
		// Should not panic, just return
	})

	t.Run("empty key", func(t *testing.T) {
		testMap := map[string]any{}
		Set(testMap, "", "value")
		t.Logf("testMap=%#v", testMap)
		// Should not panic, just return
		if len(testMap) == 0 {
			t.Errorf("Expected none empty map, got %v", testMap)
		}
	})

	t.Run("invalid slice index format - missing closing bracket", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b"}}
		Set(testMap, "items[0", "new")
		t.Logf("testMap=%#v", testMap)
		// Should not panic, just return
		expected := map[string]any{"items": []any{"a", "b"}, "items[0": "new"}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("invalid slice index format - missing opening bracket", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b"}}
		Set(testMap, "items0]", "new")
		// Should not panic, just return
		expected := map[string]any{"items": []any{"a", "b"}, "items0]": "new"}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("invalid slice index - non-numeric", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b"}}
		Set(testMap, "items[abc]", "new")
		// Should not panic, just return
		expected := map[string]any{"items": []any{"a", "b"}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("slice index out of bounds", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b"}}
		Set(testMap, "items[10]", "new")
		// Should not panic, just return
		expected := map[string]any{"items": []any{"a", "b"}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("negative slice index", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b"}}
		Set(testMap, "items[-1]", "new")
		// Should not panic, just return
		expected := map[string]any{"items": []any{"a", "b"}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("slice key not found", func(t *testing.T) {
		testMap := map[string]any{}
		Set(testMap, "items[0]", "new")
		// Should not panic, just return
		if len(testMap) != 0 {
			t.Errorf("Expected empty map, got %v", testMap)
		}
	})

	t.Run("slice element is not a map for nested access", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b"}}
		Set(testMap, "items[0].nested", "value")
		// Should not panic, just return
		expected := map[string]any{"items": []any{"a", "b"}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("non-map value in nested path", func(t *testing.T) {
		testMap := map[string]any{"level1": "not a map"}
		Set(testMap, "level1.level2", "value")
		// Should not panic, just return
		expected := map[string]any{"level1": "not a map"}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("type mismatch in slice element", func(t *testing.T) {
		testMap := map[string]any{"numbers": []int{1, 2, 3}}
		Set(testMap, "numbers[0]", "string")
		// Should not panic, just return
		expected := map[string]any{"numbers": []int{1, 2, 3}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})

	t.Run("unsupported slice type", func(t *testing.T) {
		testMap := map[string]any{"items": []float32{1.1, 2.2}}
		Set(testMap, "items[0]", 3.3)
		// Should not panic, just return
		expected := map[string]any{"items": []float32{1.1, 2.2}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected unchanged map, got %v", testMap)
		}
	})
}

// TestSet_EdgeCases tests edge cases
func TestSet_EdgeCases(t *testing.T) {
	t.Run("key with dots but no nesting", func(t *testing.T) {
		testMap := map[string]any{}
		Set(testMap, "key.with.dots", "value")
		// This will create nested maps
		expected := map[string]any{
			"key": map[string]any{
				"with": map[string]any{
					"dots": "value",
				},
			},
		}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected nested structure, got %v", testMap)
		}
	})

	t.Run("set value in map with existing nested structure", func(t *testing.T) {
		testMap := map[string]any{
			"user": map[string]any{
				"name": "John",
				"age":  30,
			},
		}
		Set(testMap, "user.email", "john@example.com")
		expected := map[string]any{
			"user": map[string]any{
				"name":  "John",
				"age":   30,
				"email": "john@example.com",
			},
		}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected updated map, got %v", testMap)
		}
	})

	t.Run("overwrite nested map with value", func(t *testing.T) {
		testMap := map[string]any{
			"user": map[string]any{
				"name": "John",
			},
		}
		Set(testMap, "user", "simple_value")
		expected := map[string]any{
			"user": "simple_value",
		}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected overwritten value, got %v", testMap)
		}
	})

	t.Run("set zero index in slice", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b", "c"}}
		Set(testMap, "items[0]", "first")
		expected := map[string]any{"items": []any{"first", "b", "c"}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected updated slice, got %v", testMap)
		}
	})

	t.Run("set last index in slice", func(t *testing.T) {
		testMap := map[string]any{"items": []any{"a", "b", "c"}}
		Set(testMap, "items[2]", "last")
		expected := map[string]any{"items": []any{"a", "b", "last"}}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected updated slice, got %v", testMap)
		}
	})
}

// TestSet_ComplexTypes tests setting complex types
func TestSet_ComplexTypes(t *testing.T) {
	t.Run("set nested map value", func(t *testing.T) {
		testMap := map[string]any{}
		nestedValue := map[string]any{
			"nested1": "value1",
			"nested2": 42,
		}
		Set(testMap, "config", nestedValue)
		expected := map[string]any{
			"config": map[string]any{
				"nested1": "value1",
				"nested2": 42,
			},
		}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected nested map, got %v", testMap)
		}
	})

	t.Run("set slice value", func(t *testing.T) {
		testMap := map[string]any{}
		sliceValue := []any{1, 2, 3}
		Set(testMap, "numbers", sliceValue)
		expected := map[string]any{
			"numbers": []any{1, 2, 3},
		}
		if !reflect.DeepEqual(testMap, expected) {
			t.Errorf("Expected slice, got %v", testMap)
		}
	})
}
