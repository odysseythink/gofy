package mapstruct

import (
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
				return default_val
			}

			// Determine the type based on default_val and cast accordingly
			switch any(default_val).(type) {
			// Basic types
			case int:
				subval, err := cast.ToIntE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int8:
				subval, err := cast.ToInt8E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int16:
				subval, err := cast.ToInt16E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int32:
				subval, err := cast.ToInt32E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case int64:
				subval, err := cast.ToInt64E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint:
				subval, err := cast.ToUintE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint8:
				subval, err := cast.ToUint8E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint16:
				subval, err := cast.ToUint16E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint32:
				subval, err := cast.ToUint32E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case uint64:
				subval, err := cast.ToUint64E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case float32:
				subval, err := cast.ToFloat32E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case float64:
				subval, err := cast.ToFloat64E(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case string:
				subval, err := cast.ToStringE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case bool:
				subval, err := cast.ToBoolE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case time.Time:
				subval, err := cast.ToTimeE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)
			case time.Duration:
				subval, err := cast.ToDurationE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subval).(T)

			// Slice types
			case []int:
				subvals, err := cast.ToIntSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int8:
				subvals, err := cast.ToInt8SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int16:
				subvals, err := cast.ToInt16SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int32:
				subvals, err := cast.ToInt32SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []int64:
				subvals, err := cast.ToInt64SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint:
				subvals, err := cast.ToUintSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint8:
				subvals, err := cast.ToUint8SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint16:
				subvals, err := cast.ToUint16SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint32:
				subvals, err := cast.ToUint32SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []uint64:
				subvals, err := cast.ToUint64SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []float32:
				subvals, err := cast.ToFloat32SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []float64:
				subvals, err := cast.ToFloat64SliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []string:
				subvals, err := cast.ToStringSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []bool:
				subvals, err := cast.ToBoolSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []time.Time:
				subvals, err := cast.ToTimeSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []time.Duration:
				subvals, err := cast.ToDurationSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []map[string]any:
				subvals, err := cast.ToStringMapSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case []any:
				subvals, err := cast.ToSliceE(current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)

			// Map types
			case map[string]int:
				subvals, err := cast.ToStrMapE[int](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int8:
				subvals, err := cast.ToStrMapE[int8](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int16:
				subvals, err := cast.ToStrMapE[int16](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int32:
				subvals, err := cast.ToStrMapE[int32](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]int64:
				subvals, err := cast.ToStrMapE[int64](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint:
				subvals, err := cast.ToStrMapE[uint](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint8:
				subvals, err := cast.ToStrMapE[uint8](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint16:
				subvals, err := cast.ToStrMapE[uint16](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint32:
				subvals, err := cast.ToStrMapE[uint32](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]uint64:
				subvals, err := cast.ToStrMapE[uint64](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]float32:
				subvals, err := cast.ToStrMapE[float32](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]float64:
				subvals, err := cast.ToStrMapE[float64](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]string:
				subvals, err := cast.ToStrMapE[string](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]bool:
				subvals, err := cast.ToStrMapE[bool](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]time.Time:
				subvals, err := cast.ToStrMapE[time.Time](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]time.Duration:
				subvals, err := cast.ToStrMapE[time.Duration](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]map[string]any:
				subvals, err := cast.ToStrMapE[map[string]any](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			case map[string]any:
				subvals, err := cast.ToStrMapE[any](current[k])
				if err != nil {
					mlog.Error("cast %#v to type(%#v) failed:%v", current[k], default_val, err)
					return default_val
				}
				return any(subvals).(T)
			default:
				mlog.Error("unsupported type(%#v)", default_val)
				return default_val
			}
		} else {
			// Not the last key - navigate to nested map
			if _, ok := current[k]; !ok {
				return default_val
			}
			next, ok := current[k].(map[string]any)
			if !ok {
				mlog.Error("key '%s' is not a map[string]any, cannot access nested field '%s'", k, key)
				return default_val
			}
			current = next
		}
	}

	return default_val
}
