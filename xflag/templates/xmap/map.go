package xmap

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type MapKTypeVType(KType,VType,ParseKeyFunc,ParseValFunc,SetProviderByFieldType,StringValueDelim)

type KType int
type VType int

var StringValueDelim = ","
var SetProviderByFieldType = func(v interface{}, flagValue interface{}) { panic(1) }
var ParseKeyFunc = func(s string) (KType, error) { panic(1) }
var ParseValFunc = func(s string) (VType, error) { panic(1) }

var typeNameMapKTypeVType = ""

func init() {
	v := map[KType]VType{}
	typeNameMapKTypeVType = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapKTypeVType, func(valPtr interface{}) flag.Getter {
		return NewMapKTypeVType(valPtr)
	})
}

// MapKTypeVType
type MapKTypeVType struct {
	s   string
	set bool
	val *map[KType]VType
}

// NewMapKTypeVType 创建指定类型
func NewMapKTypeVType(valPtr interface{}) *MapKTypeVType {
	return &MapKTypeVType{
		val: valPtr.(*map[KType]VType),
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
func (e *MapKTypeVType) String() string { return e.s }

// Set 解析时由FlagSet设定而来，进行解析
func (e *MapKTypeVType) Set(s string) error {
	e.s = s
	kv := strings.Split(s, StringValueDelim)
	if len(kv)%2 == 1 {
		return errors.New("got the odd number of input pairs")
	}
	if !e.set {
		e.set = true
		*e.val = make(map[KType]VType)
	}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s
			continue
		}
		keyVal, err := ParseKeyFunc(key)
		if err != nil {
			return err
		}
		val, err := ParseValFunc(s)
		if err != nil {
			return err
		}
		(*e.val)[keyVal] = val
	}
	return nil
}
