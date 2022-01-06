// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type MapKTypeVType(KType,VType,ParseKeyFunc,ParseValFunc,SetProviderByFieldType,StringValueDelim)
// KType 默认key类型

// VType 默认val类型

// StringValueDelim 数据分割符

// SetProviderByFieldType 替换

// ParseKeyFunc key解析，替换

// ParseValFunc val解析，替换

var typeNameMapStringString = ""

func init() {
	v := map[string]string{}
	typeNameMapStringString = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapStringString, func(valPtr interface{}) flag.Getter {
		return NewMapStringString(valPtr)
	})
}

// MapKTypeVType
type MapStringString struct {
	s   string
	set bool
	val *map[string]string
}

// NewMapKTypeVType 创建指定类型
func NewMapStringString(valPtr interface{}) *MapStringString {
	return &MapStringString{
		val: valPtr.(*map[string]string),
	}
}

// Get 返回数据，必须返回map[string]interface{}类型
func (e *MapStringString) Get() interface{} {
	vv := make(map[string]interface{})
	for k, v := range *e.val {
		vv[fmt.Sprintf("%v", k)] = v
	}
	return vv
}

// Usage  usage info for FlagSet
func (e *MapStringString) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, key and value split by %s", StringValueDelim)
}

// TypeName type name for vars FlagValue provider
func (e *MapStringString) TypeName() string { return typeNameMapStringString }

// String 获取Set设置的字符串数据？或数据转换到的？
func (e *MapStringString) String() string { return e.s }

// Set 解析时由FlagSet设定而来，进行解析
func (e *MapStringString) Set(s string) error {
	e.s = s
	kv := strings.Split(s, StringValueDelim)
	if len(kv)%2 == 1 {
		return errors.New("got the odd number of input pairs")
	}
	if !e.set {
		e.set = true
		*e.val = make(map[string]string)
	}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s
			continue
		}
		keyVal, err := parseString(key)
		if err != nil {
			return err
		}
		val, err := parseString(s)
		if err != nil {
			return err
		}
		(*e.val)[keyVal] = val
	}
	return nil
}
