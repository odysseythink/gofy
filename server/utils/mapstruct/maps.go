package mapstruct

import (
	"time"

	"mlib.com/confy/cast"
	"mlib.com/mlog"
)

func Get[V cast.Basic](val map[string]any, key string, default_val V) V {
	if _, ok := val[key]; !ok {
		return default_val
	}
	subval, err := cast.ToE[V](val[key])
	if err != nil {
		mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
		return default_val
	}
	return subval
}

func GetSlice[T cast.Basic | any](val map[string]any, key string, default_val []T) []T {
	if _, ok := val[key]; !ok {
		return default_val
	}
	switch any(default_val).(type) {
	case []int:
		subvals, err := cast.ToIntSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []int8:
		subvals, err := cast.ToInt8SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []int16:
		subvals, err := cast.ToInt16SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []int32:
		subvals, err := cast.ToInt32SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []int64:
		subvals, err := cast.ToInt64SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []uint:
		subvals, err := cast.ToUintSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []uint8:
		subvals, err := cast.ToUint8SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []uint16:
		subvals, err := cast.ToUint16SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []uint32:
		subvals, err := cast.ToUint32SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []uint64:
		subvals, err := cast.ToUint64SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []float32:
		subvals, err := cast.ToFloat32SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []float64:
		subvals, err := cast.ToFloat64SliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []string:
		subvals, err := cast.ToStringSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []bool:
		subvals, err := cast.ToBoolSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []time.Time:
		subvals, err := cast.ToTimeSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []time.Duration:
		subvals, err := cast.ToDurationSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []map[string]any:
		subvals, err := cast.ToStringMapSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	case []any:
		subvals, err := cast.ToSliceE(val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).([]T)
	default:
		mlog.Error("unsupported type(%#v)", default_val)
		return default_val
	}
}

func GetStrMap[T cast.Basic | any](val map[string]any, key string, default_val map[string]T) map[string]T {
	if _, ok := val[key]; !ok {
		return default_val
	}
	switch any(default_val).(type) {
	case map[string]int:
		subvals, err := cast.ToStrMapE[int](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]int8:
		subvals, err := cast.ToStrMapE[int8](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]int16:
		subvals, err := cast.ToStrMapE[int16](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]int32:
		subvals, err := cast.ToStrMapE[int32](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]int64:
		subvals, err := cast.ToStrMapE[int64](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]uint:
		subvals, err := cast.ToStrMapE[uint](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]uint8:
		subvals, err := cast.ToStrMapE[uint8](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]uint16:
		subvals, err := cast.ToStrMapE[uint16](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]uint32:
		subvals, err := cast.ToStrMapE[uint32](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]uint64:
		subvals, err := cast.ToStrMapE[uint64](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]float32:
		subvals, err := cast.ToStrMapE[float32](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]float64:
		subvals, err := cast.ToStrMapE[float64](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]string:
		subvals, err := cast.ToStrMapE[string](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]bool:
		subvals, err := cast.ToStrMapE[bool](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]time.Time:
		subvals, err := cast.ToStrMapE[time.Time](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]time.Duration:
		subvals, err := cast.ToStrMapE[time.Duration](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]map[string]any:
		subvals, err := cast.ToStrMapE[map[string]any](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	case map[string]any:
		subvals, err := cast.ToStrMapE[any](val[key])
		if err != nil {
			mlog.Error("cast %#v to type(%#v) failed:%v", val[key], default_val, err)
			return default_val
		}
		return any(subvals).(map[string]T)
	default:
		mlog.Error("unsupported type(%#v)", default_val)
		return default_val
	}
}
