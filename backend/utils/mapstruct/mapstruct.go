package mapstruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"mlib.com/confy/cast"
)

const (
	DEFAULT_TIME_LAYOUT = "2006-01-02 15:04:05.000"
)

func Status2Map(obj any) map[string]any {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var result = make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		// 过滤私有的,不然会G
		if t.Field(i).PkgPath != "" {
			continue
		}
		tag, have := t.Field(i).Tag.Lookup("json")
		var realName string
		if have {
			realName = tag
		} else {
			realName = t.Field(i).Name
		}
		switch v.Field(i).Kind() {
		case reflect.Struct:
			valueValue := v.Field(i).Interface()
			result[realName] = Status2Map(valueValue)
		default:
			result[realName] = v.Field(i).Interface()
		}
	}
	return result
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: StructToMap
//@description: 利用反射将结构体转化为map
//@param: obj any
//@return: map[string]any

func StructToMap(obj any, embed bool) map[string]any {
	t := reflect.TypeOf(obj)
	if !(t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct) && !(t.Kind() == reflect.Struct) {
		log.Printf("[E]obj(%#v\n) must be struct or struct pointer", obj)
		return nil
	}
	val := reflect.ValueOf(obj)
	if t.Kind() == reflect.Pointer {
		val = val.Elem()
		t = t.Elem()
	}

	data := make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}
		var realname string
		if t.Field(i).Tag.Get("mapstructure") != "" {
			realname = t.Field(i).Tag.Get("mapstructure")
		} else if t.Field(i).Tag.Get("json") != "" {
			realname = t.Field(i).Tag.Get("json")
		} else {
			realname = t.Field(i).Name
		}
		switch val.Field(i).Kind() {
		case reflect.Bool:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			fallthrough
		case reflect.String:
			data[realname] = val.Field(i).Interface()
		case reflect.Struct:
			switch val.Field(i).Interface().(type) {
			case time.Time:
				data[realname] = val.Field(i).Interface().(time.Time).Format(DEFAULT_TIME_LAYOUT)
			default:
				data[realname] = StructToMap(val.Field(i).Interface(), embed)
			}
		case reflect.Slice:
			// valval := val.Field(i).
			if val.Field(i).Cap() == 0 {
				data[realname] = "[]"
			} else {
				tmp, _ := json.Marshal(val.Field(i).Interface())
				data[realname] = string(tmp)
			}
		case reflect.Array:
			valval := val.Field(i).Interface()
			if valval == nil {
				data[realname] = "[]"
			} else {
				tmp, _ := json.Marshal(val.Field(i).Interface())
				data[realname] = string(tmp)
			}
		case reflect.Map:
			if val.Field(i).Interface() == nil {
				data[realname] = "{}"
			} else {
				tmp, _ := json.Marshal(val.Field(i).Interface())
				data[realname] = string(tmp)
			}
		case reflect.Ptr:
			if !val.Field(i).IsNil() {
				switch val.Field(i).Interface().(type) {
				case *time.Time:
					data[realname] = val.Field(i).Interface().(*time.Time).Format(DEFAULT_TIME_LAYOUT)

				default:
					bintmp, err := json.Marshal(val.Field(i).Interface())
					if err != nil {
						log.Printf("json.Marshal(src[%s]=%#v) failed:%v", realname, val.Field(i).Interface(), err)
					} else {
						data[realname] = string(bintmp)
					}
				}
			} else {
				data[realname] = ""
			}
		default:
		}
	}
	return data
}

func StructToStrStrMap(obj any) map[string]string {
	t := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	if t.Kind() == reflect.Pointer {
		val = val.Elem()
		t = t.Elem()
	}

	data := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}
		var realname string
		if t.Field(i).Tag.Get("mapstructure") != "" {
			realname = t.Field(i).Tag.Get("mapstructure")
		} else if t.Field(i).Tag.Get("json") != "" {
			realname = t.Field(i).Tag.Get("json")
		} else {
			realname = t.Field(i).Name
		}
		switch val.Field(i).Kind() {
		case reflect.Bool:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			fallthrough
		case reflect.String:
			data[realname] = cast.ToString(val.Field(i).Interface())
		case reflect.Struct:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Ptr:
			if !val.Field(i).IsNil() {
				switch val.Field(i).Interface().(type) {
				case *time.Time:
					data[realname] = val.Field(i).Interface().(*time.Time).Format(DEFAULT_TIME_LAYOUT)

				default:
					bintmp, err := json.Marshal(val.Field(i).Interface())
					if err != nil {
						log.Printf("json.Marshal(src[%s]=%#v) failed:%v", realname, val.Field(i).Interface(), err)
					} else {
						data[realname] = string(bintmp)
					}
				}
			} else {
				data[realname] = ""
			}
		default:
		}
	}
	return data
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ArrayToString
//@description: 将数组格式化为字符串
//@param: array []any
//@return: string

func ArrayToString(array []any) string {
	if array == nil {
		return "[]"
	} else {
		return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
	}
}

func StrStrMapToStruct(src map[string]string, dst any) error {
	if src == nil {
		log.Printf("empty src")
		return errors.New("empty src")
	}
	t := reflect.TypeOf(dst)
	val := reflect.ValueOf(dst)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		log.Printf("dst must be a struct pointer")
		return errors.New("dst must be a struct pointer")
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		val = val.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}
		var realname string
		if t.Field(i).Tag.Get("mapstructure") != "" {
			realname = t.Field(i).Tag.Get("mapstructure")
		} else if t.Field(i).Tag.Get("json") != "" {
			realname = t.Field(i).Tag.Get("json")
		} else {
			realname = t.Field(i).Name
		}
		switch val.Field(i).Kind() {
		case reflect.Bool:
			if _, ok := src[realname]; ok {
				fmt.Printf("src[%s]=%s cast.ToBool(src[realname])=%v\n", realname, src[realname], cast.ToBool(src[realname]))
				val.Field(i).SetBool(cast.ToBool(src[realname]))
			}
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			if _, ok := src[realname]; ok {
				val.Field(i).SetInt(int64(cast.ToInt(src[realname])))
			}
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			if _, ok := src[realname]; ok {
				val.Field(i).SetUint(uint64(cast.ToUint(src[realname])))
			}
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			if _, ok := src[realname]; ok {
				val.Field(i).SetFloat(float64(cast.ToFloat64(src[realname])))
			}
		case reflect.String:
			if _, ok := src[realname]; ok {
				val.Field(i).SetString(cast.ToString(src[realname]))
			}
		case reflect.Struct:
			switch val.Field(i).Interface().(type) {
			case time.Time:
				if tmp, ok := src[realname]; ok {
					tm, err := cast.ToTimeE(tmp)
					if err != nil {
						log.Printf("cast.ToTimeE(src[%s]=%#v) failed:%v", realname, tmp, err)
					} else {
						val.Field(i).Set(reflect.ValueOf(tm))
					}
				}
			default:
				tmpv := reflect.New(t.Field(i).Type).Interface()
				if tmp, ok := src[realname]; ok {
					err := json.Unmarshal([]byte(tmp), tmpv)
					if err != nil {
						log.Printf("json.Unmarshal(src[%s]=%#v) failed:%v", realname, tmp, err)
					} else {
						val.Field(i).Set(reflect.ValueOf(tmpv))
					}
				}
			}
		case reflect.Ptr:
			if val.Field(i).CanSet() {
				switch val.Field(i).Interface().(type) {
				case *time.Time:
					if tmp, ok := src[realname]; ok {
						tm, err := cast.ToTimeE(tmp)
						if err != nil {
							log.Printf("cast.ToTimeE(src[%s]=%#v) failed:%v", realname, tmp, err)
						} else {
							val.Field(i).Set(reflect.ValueOf(&tm))
						}
					}
				default:
					tmpv := reflect.New(t.Field(i).Type).Interface()
					if tmp, ok := src[realname]; ok {
						err := json.Unmarshal([]byte(tmp), tmpv)
						if err != nil {
							log.Printf("json.Unmarshal(src[%s]=%#v) failed:%v", realname, tmp, err)
						} else {
							val.Field(i).Set(reflect.ValueOf(tmpv))
						}
					}
				}
			}

		case reflect.Slice:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Map:
			tmpv := reflect.New(t.Field(i).Type).Interface()
			if tmp, ok := src[realname]; ok {
				err := json.Unmarshal([]byte(tmp), tmpv)
				if err != nil {
					log.Printf("json.Unmarshal(src[%s]=%#v) failed:%v", realname, tmp, err)
				} else {
					val.Field(i).Set(reflect.ValueOf(tmpv))
				}
			}
		default:
		}
	}
	return nil
}

func MapToStruct(src map[string]any, dst any) error {
	if src == nil {
		log.Printf("empty src")
		return errors.New("empty src")
	}
	t := reflect.TypeOf(dst)
	val := reflect.ValueOf(dst)
	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		log.Printf("dst must be a struct pointer")
		return errors.New("dst must be a struct pointer")
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		val = val.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			continue
		}
		var realname string
		if t.Field(i).Tag.Get("mapstructure") != "" {
			realname = t.Field(i).Tag.Get("mapstructure")
		} else if t.Field(i).Tag.Get("json") != "" {
			realname = t.Field(i).Tag.Get("json")
		} else {
			realname = t.Field(i).Name
		}
		switch val.Field(i).Kind() {
		case reflect.Bool:
			if _, ok := src[realname]; ok {
				val.Field(i).SetBool(cast.ToBool(src[realname]))
			}
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			if _, ok := src[realname]; ok {
				val.Field(i).SetInt(int64(cast.ToInt(src[realname])))
			}
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			if _, ok := src[realname]; ok {
				val.Field(i).SetUint(uint64(cast.ToUint(src[realname])))
			}
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			if _, ok := src[realname]; ok {
				val.Field(i).SetFloat(float64(cast.ToFloat64(src[realname])))
			}
		case reflect.String:
			if _, ok := src[realname]; ok {
				val.Field(i).SetString(cast.ToString(src[realname]))
			}
		case reflect.Struct:
			switch val.Field(i).Interface().(type) {
			case time.Time:
				if tmp, ok := src[realname]; ok {
					tm, err := cast.ToTimeE(tmp)
					if err != nil {
						log.Printf("cast.ToTimeE(src[%s]=%#v) failed:%v", realname, tmp, err)
					} else {
						val.Field(i).Set(reflect.ValueOf(tm))
					}
				}
			default:
				if tmp, ok := src[realname]; ok {
					bintmp, err := json.Marshal(tmp)
					if err != nil {
						log.Printf("json.Marshal(src[%s]=%#v) failed:%v", realname, tmp, err)
					} else {
						tmpv := reflect.New(t.Field(i).Type).Interface()
						err := json.Unmarshal(bintmp, tmpv)
						if err != nil {
							log.Printf("json.Unmarshal(src[%s]=%#v) failed:%v", realname, bintmp, err)
						} else {
							val.Field(i).Set(reflect.ValueOf(tmpv))
						}
					}
				}
			}
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Map:
			tmpv := reflect.New(t.Field(i).Type).Interface()
			if tmp, ok := src[realname]; ok {
				bintmp, err := json.Marshal(tmp)
				if err != nil {
					log.Printf("json.Marshal(src[%s]=%#v) failed:%v", realname, tmp, err)
				} else {
					err := json.Unmarshal([]byte(bintmp), tmpv)
					if err != nil {
						log.Printf("json.Unmarshal(src[%s]=%#v) failed:%v", realname, string(bintmp), err)
					} else {
						val.Field(i).Set(reflect.ValueOf(tmpv))
					}
				}
			}
		case reflect.Ptr:
			if val.Field(i).CanSet() {
				switch val.Field(i).Interface().(type) {
				case *time.Time:
					if tmp, ok := src[realname]; ok {
						tm, err := cast.ToTimeE(tmp)
						if err != nil {
							log.Printf("cast.ToTimeE(src[%s]=%#v) failed:%v", realname, tmp, err)
						} else {
							val.Field(i).Set(reflect.ValueOf(&tm))
						}
					}
				default:
					val.Field(i).Set(val.Field(i))
				}
			}
		default:
		}
	}
	return nil
}
