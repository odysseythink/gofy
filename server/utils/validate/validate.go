package validate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"mlib.com/confy/cast"
	"mlib.com/mlog"
)

const (
	PHONENUM_REGX_STRING = `^1(3[0-9]|5[0-3,5-9]|7[1-3,5-8]|8[0-9])\d{8}$`
	IP_REGX_STRING       = `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: RegisterRule
//@description: 注册自定义规则方案建议在路由初始化层即注册
//@param: key string, rule Rules
//@return: err error

func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: NotEmpty
//@description: 非空 不能为其对应类型的0值
//@return: string

func NotEmpty() string {
	return "notEmpty"
}

// @author: [zooqkl](https://github.com/zooqkl)
// @function: RegexpMatch
// @description: 正则校验 校验输入项是否满足正则表达式
// @param:  rule string
// @return: string

func RegexpMatch(rule string) string {
	return "regexp=" + rule
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Lt
//@description: 小于入参(<) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Lt(mark string) string {
	return "lt=" + mark
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Le
//@description: 小于等于入参(<=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Le(mark string) string {
	return "le=" + mark
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Eq
//@description: 等于入参(==) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Eq(mark string) string {
	return "eq=" + mark
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Ne
//@description: 不等于入参(!=)  如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Ne(mark string) string {
	return "ne=" + mark
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Ge
//@description: 大于等于入参(>=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Ge(mark string) string {
	return "ge=" + mark
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Gt
//@description: 大于入参(>) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
//@param: mark string
//@return: string

func Gt(mark string) string {
	return "gt=" + mark
}

//
//@author: [piexlmax](https://github.com/piexlmax)
//@function: Verify
//@description: 校验方法
//@param: st any, roleMap Rules(入参实例，规则map)
//@return: err error

func Verify(st any, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		if tagVal.Type.Kind() == reflect.Struct {
			if err = Verify(val.Interface(), roleMap); err != nil {
				return err
			}
		}
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmpty":
					if isBlank(val) {
						return errors.New(tagVal.Name + "值不能为空")
					}
				case strings.Split(v, "=")[0] == "regexp":
					if !regexpMatch(strings.Split(v, "=")[1], val.String()) {
						return errors.New(tagVal.Name + "格式校验不通过")
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
					}
				}
			}
			delete(roleMap, tagVal.Name)
		}
	}
	if len(roleMap) > 0 {
		keys := make([]string, 0)
		for k := range roleMap {
			keys = append(keys, k)
		}
		return fmt.Errorf("fields(%#v) not exist in st(%#v)", keys, st)
	}
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: compareVerify
//@description: 长度和数字的校验方法 根据类型自动校验
//@param: value reflect.Value, VerifyStr string
//@return: bool

func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String:
		return compare(len([]rune(value.String())), VerifyStr)
	case reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: isBlank
//@description: 非空校验
//@param: value reflect.Value
//@return: bool

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: compare
//@description: 比较函数
//@param: value any, VerifyStr string
//@return: bool

func compare(value any, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}
}

func regexpMatch(rule, matchStr string) bool {
	return regexp.MustCompile(rule).MatchString(matchStr)
}

func IsPhoneNum(str string) bool {
	if str == "" {
		mlog.Errorf("invalid arg")
		return false
	}
	regx, err := regexp.Compile(PHONENUM_REGX_STRING)
	if err != nil {
		mlog.Errorf("regexp.Compile %s failed:%v", PHONENUM_REGX_STRING, err)
		return false
	}
	matched := regx.FindString(str)
	if matched != str {
		mlog.Errorf("str(%s) contains other things that is not a phone num", str)
		return false
	}
	return true
}

func IsValidIP(str string) bool {
	if str == "" {
		mlog.Errorf("invalid arg")
		return false
	}
	regx, err := regexp.Compile(IP_REGX_STRING)
	if err != nil {
		mlog.Errorf("regexp.Compile %s failed:%v", IP_REGX_STRING, err)
		return false
	}
	matched := regx.FindString(str)
	if matched != str {
		mlog.Errorf("str(%s) contains other things that is not a ip", str)
		return false
	}
	return true
}

func IsValidNetPort(str string) bool {
	if str == "" {
		mlog.Errorf("invalid arg")
		return false
	}
	port, err := cast.ToE[int](str)
	if err != nil {
		return false
	}
	if port < 0 || port > 65535 {
		return false
	}
	return true
}

// type TypeRules map[string]reflect.Kind

func InterfaceTypeVerify(st any, t reflect.Kind) bool {
	if st == nil {
		return false
		// return errors.New("invalid arg")
	}
	val := reflect.ValueOf(st) // 获取reflect.Type类型
	switch val.Kind() {
	case t:
		return true
	default:
		return false
	}
}

type TypeRules map[string]reflect.Kind

var (
	typesOfFieldRelationMap = map[string]reflect.Kind{
		"typeOfField=1":  reflect.Bool,
		"typeOfField=2":  reflect.Int,
		"typeOfField=3":  reflect.Int8,
		"typeOfField=4":  reflect.Int16,
		"typeOfField=5":  reflect.Int32,
		"typeOfField=6":  reflect.Int64,
		"typeOfField=7":  reflect.Uint,
		"typeOfField=8":  reflect.Uint8,
		"typeOfField=9":  reflect.Uint16,
		"typeOfField=10": reflect.Uint32,
		"typeOfField=11": reflect.Uint64,
		"typeOfField=12": reflect.Uintptr,
		"typeOfField=13": reflect.Float32,
		"typeOfField=14": reflect.Float64,
		"typeOfField=15": reflect.Complex64,
		"typeOfField=16": reflect.Complex128,
		"typeOfField=17": reflect.Array,
		"typeOfField=18": reflect.Chan,
		"typeOfField=19": reflect.Func,
		"typeOfField=20": reflect.Interface,
		"typeOfField=21": reflect.Map,
		"typeOfField=22": reflect.Pointer,
		"typeOfField=23": reflect.Slice,
		"typeOfField=24": reflect.String,
		"typeOfField=25": reflect.Struct,
		"typeOfField=26": reflect.UnsafePointer,
	}
)

func getTypeOfField(str string) reflect.Kind {
	if !strings.HasPrefix(str, "typeOfField=") {
		return reflect.Invalid
	}
	if val, ok := typesOfFieldRelationMap[str]; !ok {
		return reflect.Invalid
	} else {
		return val
	}
}

func RuleTypeOfField(t reflect.Kind) string {
	return fmt.Sprintf("typeOfField=%d", t)
}

func StringMapTypeVerify(src any, rs Rules) error {
	if src == nil || rs == nil {
		return errors.New("invalid arg")
	}
	if st, ok := src.(map[string]any); !ok {
		return errors.New("src must be a map[string]any")
	} else {
		for k, v := range rs {
			if val, ok := st[k]; !ok {
				return fmt.Errorf("field(%s) not exist in map(%#v)", k, st)
			} else {
				tmp := reflect.ValueOf(val)
				for _, v1 := range v {
					switch {
					case v1 == "notEmpty":
						if isBlank(tmp) {
							return fmt.Errorf("field(%s) is empty", k)
						}
					case strings.HasPrefix(v1, "typeOfField="):
						t := getTypeOfField(v1)
						if t == reflect.Invalid {
							return fmt.Errorf("field(%s) is invalid type(%s)", k, v1)
						}
						switch tmp.Kind() {
						case t:
							break
						default:
							return fmt.Errorf("field(%s) type not matched,wanted(%v), real(%v)", k, v, tmp.Kind())
						}
					}
				}
			}
		}
		return nil
	}
}

func IsNumber(val any) bool {
	switch val.(type) {
	case int:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	}
	return false
}
