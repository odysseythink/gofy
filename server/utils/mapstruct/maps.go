package mapstruct

import "mlib.com/mlog"

func Get[T int | bool | string | map[string]any | float64 | []map[string]any](val map[string]any, key string, default_val T) T {
	if _, ok := val[key]; ok {
		if _, ok := val[key].(T); ok {
			return val[key].(T)
		} else {
			switch any(default_val).(type) {
			case []map[string]any:
				if _, ok := val[key].([]any); ok {
					ret := []map[string]any{}
					for _, sv := range val[key].([]any) {
						if _, ok := sv.(map[string]any); ok {
							ret = append(ret, sv.(map[string]any))
						} else {
							mlog.Warningf("val[%s]=%#v is not map slice", key, val[key])
							return default_val
						}
					}
					return any(ret).(T)
				}
			}
			mlog.Warningf("val[%s]=%#v is not bool", key, val[key])
			return default_val
		}
	} else {
		mlog.Warningf("val dosen't exist field=%s", key)
		return default_val
	}
}
