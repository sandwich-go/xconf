package xmap

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type MapKTypeVType(KType,VType,ParseKeyFunc,ParseValFunc,SetProviderByFieldType,StringValueDelim)

// KType 默认key类型
type KType int

// VType 默认val类型
type VType int

// StringValueDelim 数据分割符
var StringValueDelim = ","

// SetProviderByFieldType 替换
var SetProviderByFieldType = func(v interface{}, flagValue interface{}) { panic(1) }

// ParseKeyFunc key解析，替换
var ParseKeyFunc = func(s string) (KType, error) { panic(1) }

// ParseValFunc val解析，替换
var ParseValFunc = func(s string) (VType, error) { panic(1) }

var typeNameMapKTypeVType = ""

func init() {
	v := map[KType]VType{}
	typeNameMapKTypeVType = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapKTypeVType, func(valPtr interface{}, stringAlias func(s string) string) flag.Getter {
		return NewMapKTypeVType(valPtr, stringAlias)
	})
}

// MapKTypeVType new func
type MapKTypeVType struct {
	stringAlias func(s string) string
	s           string
	set         bool
	val         *map[KType]VType
}

// NewMapKTypeVType 创建指定类型
func NewMapKTypeVType(valPtr interface{}, stringAlias func(s string) string) *MapKTypeVType {
	return &MapKTypeVType{
		val:         valPtr.(*map[KType]VType),
		stringAlias: stringAlias,
	}
}

// Get 返回数据，必须返回map[string]interface{}类型
func (e *MapKTypeVType) Get() interface{} {
	vv := make(map[string]interface{})
	for k, v := range *e.val {
		vv[fmt.Sprintf("%v", k)] = v
	}
	return vv
}

// Usage  usage info for FlagSet
func (e *MapKTypeVType) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, key and value split by %s", StringValueDelim)
}

// TypeName type name for vars FlagValue provider
func (e *MapKTypeVType) TypeName() string { return typeNameMapKTypeVType }

// String 获取Set设置的字符串数据？或数据转换到的？
func (e *MapKTypeVType) String() string {
	if e.val == nil || len(*e.val) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", *e.val)
}

// Set 解析时由FlagSet设定而来，进行解析
func (e *MapKTypeVType) Set(s string) error {
	s = e.stringAlias(s)
	e.s = s
	kv := strings.Split(s, StringValueDelim)
	if len(kv)%2 == 1 {
		// 设定了default标签或者空的字符串
		if len(kv) == 1 && kv[0] == "" {
			return nil
		}
		return errors.New("got the odd number of input pairs")
	}
	if !e.set {
		e.set = true
		*e.val = make(map[KType]VType)
	}
	var key string
	for i, val := range kv {
		if i%2 == 0 {
			key = val
			continue
		}
		keyAlias := e.stringAlias(key)
		keyVal, err := ParseKeyFunc(keyAlias)
		if err != nil {
			return fmt.Errorf("got err:%s while parse key:%s alias:%s raw:%s", err.Error(), key, keyAlias, s)
		}
		valAlias := e.stringAlias(val)
		valVal, err := ParseValFunc(valAlias)
		if err != nil {
			return fmt.Errorf("got err:%s while parse val:%s alias:%s raw:%s", err.Error(), val, valAlias, s)
		}
		(*e.val)[keyVal] = valVal
	}
	return nil
}
