package mapstruct

import (
	"fmt"
	"strings"
	"time"

	"mlib.com/confy/cast"
	"mlib.com/mlog"
)

func Get[T any](val map[string]any, key string, default_val T) T {
	// Split key by dots to support nested field access
	keys := strings.Split(key, ".")

	// Navigate through nested maps
	current := val
	for i, k := range keys {
		if i == len(keys)-1 {
			// Last key - get the value
			if _, ok := current[k]; !ok {
				fmt.Println("current[", k, "] not exist")
				return default_val
			}

			// Determine the type based on default_val and cast accordingly
			switch any(default_val).(type) {
			// Basic types
			case int:
				subval, err := cast.ToIntE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int8:
				subval, err := cast.ToInt8E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int16:
				subval, err := cast.ToInt16E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int32:
				subval, err := cast.ToInt32E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int64:
				subval, err := cast.ToInt64E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint:
				subval, err := cast.ToUintE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint8:
				subval, err := cast.ToUint8E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint16:
				subval, err := cast.ToUint16E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint32:
				subval, err := cast.ToUint32E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint64:
				subval, err := cast.ToUint64E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case float32:
				subval, err := cast.ToFloat32E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case float64:
				subval, err := cast.ToFloat64E(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case string:
				subval, err := cast.ToStringE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case bool:
				subval, err := cast.ToBoolE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case time.Time:
				subval, err := cast.ToTimeE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case time.Duration:
				subval, err := cast.ToDurationE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)

			// Slice types
			case []int:
				subvals, err := cast.ToIntSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int8:
				subvals, err := cast.ToInt8SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int16:
				subvals, err := cast.ToInt16SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int32:
				subvals, err := cast.ToInt32SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int64:
				subvals, err := cast.ToInt64SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint:
				subvals, err := cast.ToUintSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint8:
				subvals, err := cast.ToUint8SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint16:
				subvals, err := cast.ToUint16SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint32:
				subvals, err := cast.ToUint32SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint64:
				subvals, err := cast.ToUint64SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []float32:
				subvals, err := cast.ToFloat32SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []float64:
				subvals, err := cast.ToFloat64SliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []string:
				subvals, err := cast.ToStringSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []bool:
				subvals, err := cast.ToBoolSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []time.Time:
				subvals, err := cast.ToTimeSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []time.Duration:
				subvals, err := cast.ToDurationSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []map[string]any:
				subvals, err := cast.ToStringMapSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []any:
				subvals, err := cast.ToSliceE(current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)

			// Map types
			case map[string]int:
				subvals, err := cast.ToStrMapE[int](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int8:
				subvals, err := cast.ToStrMapE[int8](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int16:
				subvals, err := cast.ToStrMapE[int16](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int32:
				subvals, err := cast.ToStrMapE[int32](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int64:
				subvals, err := cast.ToStrMapE[int64](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint:
				subvals, err := cast.ToStrMapE[uint](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint8:
				subvals, err := cast.ToStrMapE[uint8](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint16:
				subvals, err := cast.ToStrMapE[uint16](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint32:
				subvals, err := cast.ToStrMapE[uint32](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint64:
				subvals, err := cast.ToStrMapE[uint64](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]float32:
				subvals, err := cast.ToStrMapE[float32](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]float64:
				subvals, err := cast.ToStrMapE[float64](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]string:
				subvals, err := cast.ToStrMapE[string](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]bool:
				subvals, err := cast.ToStrMapE[bool](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]time.Time:
				subvals, err := cast.ToStrMapE[time.Time](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]time.Duration:
				subvals, err := cast.ToStrMapE[time.Duration](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]map[string]any:
				subvals, err := cast.ToStrMapE[map[string]any](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]any:
				subvals, err := cast.ToStrMapE[any](current[k])
				if err != nil {
					mlog.Errorf("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			default:
				mlog.Errorf("unsupported type(%#v)", default_val)
				return default_val
			}
		} else {
			// Not the last key - navigate to nested map
			if _, ok := current[k]; !ok {
				return default_val
			}
			next, ok := current[k].(map[string]any)
			if !ok {
				mlog.Errorf("key '%s' is not a map[string]any, cannot access nested field '%s'", k, key)
				return default_val
			}
			current = next
		}
	}

	return default_val
}

// Set sets a value in a nested map using dot notation (e.g., "a.b.c")
// Supports slice index access using bracket notation (e.g., "a.b[0].c")
// Parameters:
//   - dst: the destination map to write to
//   - key: the key path using dot notation (e.g., "a.b.c" or "a.b[0].c")
//   - val: the value to set
func Set(dst map[string]any, key string, val any) {
	if dst == nil {
		dst = make(map[string]any)
	}

	// Split key by dots to support nested field access
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		mlog.Error("key is empty")
		return
	}

	// Navigate through nested maps, creating them if they don't exist
	current := dst
	for i, k := range keys {
		if i == len(keys)-1 {
			// Last key - set the value
			// Check if this key contains a slice index (e.g., "items[0]")
			if strings.Contains(k, "[") && strings.Contains(k, "]") {
				// Parse slice key and index
				sliceKey, index, err := parseSliceKey(k)
				if err != nil {
					mlog.Errorf("failed to parse slice key '%s': %v", k, err)
					return
				}

				// Get the slice
				slice, ok := current[sliceKey]
				if !ok {
					mlog.Errorf("slice key '%s' not found", sliceKey)
					return
				}

				// Set the element at the specified index
				if err := setSliceElement(slice, index, val); err != nil {
					mlog.Errorf("failed to set slice element at index %d: %v", index, err)
					return
				}
			} else {
				current[k] = val
			}
		} else {
			// Not the last key - navigate to or create nested map
			// Check if this key contains a slice index (e.g., "items[0]")
			if strings.Contains(k, "[") && strings.Contains(k, "]") {
				// Parse slice key and index
				sliceKey, index, err := parseSliceKey(k)
				if err != nil {
					mlog.Errorf("failed to parse slice key '%s': %v", k, err)
					return
				}

				// Get the slice
				slice, ok := current[sliceKey]
				if !ok {
					mlog.Errorf("slice key '%s' not found", sliceKey)
					return
				}

				// Get the element at the specified index
				element, err := getSliceElement(slice, index)
				if err != nil {
					mlog.Errorf("failed to get slice element at index %d: %v", index, err)
					return
				}

				// Check if the element is a map
				next, ok := element.(map[string]any)
				if !ok {
					mlog.Errorf("slice element at index %d is not a map[string]any, cannot access nested field '%s'", index, key)
					return
				}
				current = next
			} else {
				// Not the last key - navigate to or create nested map
				if _, ok := current[k]; !ok {
					// Create a new nested map
					current[k] = make(map[string]any)
				}
				next, ok := current[k].(map[string]any)
				if !ok {
					mlog.Errorf("key '%s' is not a map[string]any, cannot set nested field '%s'", k, key)
					return
				}
				current = next
			}
		}
	}
}

// parseSliceKey parses a slice key with index (e.g., "items[0]") and returns the key and index
func parseSliceKey(key string) (string, int, error) {
	openBracket := strings.Index(key, "[")
	closeBracket := strings.Index(key, "]")

	if openBracket == -1 || closeBracket == -1 || closeBracket <= openBracket {
		return "", 0, fmt.Errorf("invalid slice key format: %s", key)
	}

	sliceKey := key[:openBracket]
	indexStr := key[openBracket+1 : closeBracket]

	index, err := cast.ToIntE(indexStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid index '%s': %v", indexStr, err)
	}

	return sliceKey, index, nil
}

// setSliceElement sets an element at the specified index in a slice
func setSliceElement(slice any, index int, val any) error {
	switch s := slice.(type) {
	case []any:
		if index < 0 || index >= len(s) {
			return fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		s[index] = val
	case []int:
		if index < 0 || index >= len(s) {
			return fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		casted, err := cast.ToIntE(val)
		if err != nil {
			return fmt.Errorf("failed to cast value to int: %v", err)
		}
		s[index] = casted
	case []string:
		if index < 0 || index >= len(s) {
			return fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		casted, err := cast.ToStringE(val)
		if err != nil {
			return fmt.Errorf("failed to cast value to string: %v", err)
		}
		s[index] = casted
	case []bool:
		if index < 0 || index >= len(s) {
			return fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		casted, err := cast.ToBoolE(val)
		if err != nil {
			return fmt.Errorf("failed to cast value to bool: %v", err)
		}
		s[index] = casted
	case []float64:
		if index < 0 || index >= len(s) {
			return fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		casted, err := cast.ToFloat64E(val)
		if err != nil {
			return fmt.Errorf("failed to cast value to float64: %v", err)
		}
		s[index] = casted
	case []map[string]any:
		if index < 0 || index >= len(s) {
			return fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		casted, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("value is not a map[string]any")
		}
		s[index] = casted
	default:
		return fmt.Errorf("unsupported slice type: %T", slice)
	}
	return nil
}

// getSliceElement gets an element at the specified index from a slice
func getSliceElement(slice any, index int) (any, error) {
	switch s := slice.(type) {
	case []any:
		if index < 0 || index >= len(s) {
			return nil, fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		return s[index], nil
	case []map[string]any:
		if index < 0 || index >= len(s) {
			return nil, fmt.Errorf("index %d out of bounds for slice of length %d", index, len(s))
		}
		return s[index], nil
	default:
		return nil, fmt.Errorf("unsupported slice type for element access: %T", slice)
	}
}
