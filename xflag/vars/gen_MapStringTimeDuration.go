// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"
)

//template type MapKTypeVType(KType,VType,ParseKeyFunc,ParseValFunc,SetProviderByFieldType,StringValueDelim)

// KType 默认key类型

// VType 默认val类型

// StringValueDelim 数据分割符

// SetProviderByFieldType 替换

// ParseKeyFunc key解析，替换

// ParseValFunc val解析，替换

var typeNameMapStringTimeDuration = ""

func init() {
	v := map[string]time.Duration{}
	typeNameMapStringTimeDuration = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapStringTimeDuration, func(valPtr interface{}, stringAlias func(s string) string) flag.Getter {
		return NewMapStringTimeDuration(valPtr, stringAlias)
	})
}

// MapKTypeVType new func
type MapStringTimeDuration struct {
	stringAlias func(s string) string
	s           string
	set         bool
	val         *map[string]time.Duration
}

// NewMapKTypeVType 创建指定类型
func NewMapStringTimeDuration(valPtr interface{}, stringAlias func(s string) string) *MapStringTimeDuration {
	return &MapStringTimeDuration{
		val:         valPtr.(*map[string]time.Duration),
		stringAlias: stringAlias,
	}
}

// Get 返回数据，必须返回map[string]interface{}类型
func (e *MapStringTimeDuration) Get() interface{} {
	vv := make(map[string]interface{})
	for k, v := range *e.val {
		vv[fmt.Sprintf("%v", k)] = v
	}
	return vv
}

// Usage  usage info for FlagSet
func (e *MapStringTimeDuration) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, key and value split by %s", StringValueDelim)
}

// TypeName type name for vars FlagValue provider
func (e *MapStringTimeDuration) TypeName() string { return typeNameMapStringTimeDuration }

// String 获取Set设置的字符串数据？或数据转换到的？
func (e *MapStringTimeDuration) String() string {
	if e.val == nil || len(*e.val) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", *e.val)
}

// Set 解析时由FlagSet设定而来，进行解析
func (e *MapStringTimeDuration) Set(s string) error {
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
		*e.val = make(map[string]time.Duration)
	}
	var key string
	for i, val := range kv {
		if i%2 == 0 {
			key = val
			continue
		}
		keyAlias := e.stringAlias(key)
		keyVal, err := parseString(keyAlias)
		if err != nil {
			return fmt.Errorf("got err:%s while parse key:%s alias:%s raw:%s", err.Error(), key, keyAlias, s)
		}
		valAlias := e.stringAlias(val)
		valVal, err := parseTimeDuration(valAlias)
		if err != nil {
			return fmt.Errorf("got err:%s while parse val:%s alias:%s raw:%s", err.Error(), val, valAlias, s)
		}
		(*e.val)[keyVal] = valVal
	}
	return nil
}
