package mapstruct

import "github.com/odysseythink/mlog"

func Get[T int | bool | string | map[string]any | float64 | []map[string]any | []string | map[string][]string | [][]string](val map[string]any, key string, default_val T) T {
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
			case []string:
				if _, ok := val[key].([]any); ok {
					ret := []string{}
					for _, sv := range val[key].([]any) {
						if _, ok := sv.(string); ok {
							ret = append(ret, sv.(string))
						} else {
							mlog.Warningf("val[%s]=%#v is not string slice", key, val[key])
							return default_val
						}
					}
					return any(ret).(T)
				}
			case map[string][]string:
				if _, ok := val[key].(map[string]any); ok {
					ret := map[string][]string{}
					for k, sv := range val[key].(map[string]any) {
						if _, ok := sv.([]string); ok {
							ret[k] = sv.([]string)
						} else if _, ok := sv.([]any); ok {
							for _, ssv := range sv.([]any) {
								if _, ok := ssv.(string); ok {
									if _, ok := ret[k]; ok {
										ret[k] = make([]string, 0)
									}
									ret[k] = append(ret[k], ssv.(string))
								} else {
									mlog.Warningf("val[%s]=%#v is not map[string][]string", key, val[key])
									return default_val
								}
							}
						} else {
							mlog.Warningf("val[%s]=%#v is not map[string][]string", key, val[key])
							return default_val
						}
					}
					return any(ret).(T)
				} else if _, ok := val[key].(map[string][]any); ok {
					ret := map[string][]string{}
					for k, sv := range val[key].(map[string][]any) {
						for _, ssv := range sv {
							if _, ok := ssv.(string); ok {
								if _, ok := ret[k]; ok {
									ret[k] = make([]string, 0)
								}
								ret[k] = append(ret[k], ssv.(string))
							} else {
								mlog.Warningf("val[%s]=%#v is not map[string][]string", key, val[key])
								return default_val
							}
						}
					}
					return any(ret).(T)
				}
			case [][]string:
				if _, ok := val[key].([]any); ok {
					ret := make([][]string, len(val[key].([]any)))
					for idx, sv := range val[key].([]any) {
						if _, ok := sv.([]string); ok {
							ret[idx] = sv.([]string)
						} else if _, ok := sv.([]any); ok {
							ret[idx] = make([]string, 0)
							for _, ssv := range sv.([]any) {
								if _, ok := ssv.(string); ok {
									ret[idx] = append(ret[idx], ssv.(string))
								} else {
									mlog.Warningf("val[%s]=%#v is not string slice", key, val[key])
									return default_val
								}
							}
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
