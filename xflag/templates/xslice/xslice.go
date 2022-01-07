package vars

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type Slice(KType,ParseKeyFunc,SetProviderByFieldType,StringValueDelim)

// KType 默认key类型
type KType int

// SetProviderByFieldType 替换
var SetProviderByFieldType = func(v interface{}, flagValue interface{}) { panic(1) }

// ParseKeyFunc val解析，替换
var ParseKeyFunc = func(s string) (KType, error) { panic(1) }
var typeNameSlice = ""

// StringValueDelim 数据分割符
var StringValueDelim = ""

func init() {
	v := []KType{}
	typeNameSlice = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSlice, func(valPtr interface{}) flag.Getter {
		return NewSlice(valPtr.(*[]KType))
	})
}

// Slice struct
type Slice struct {
	s   *[]KType
	set bool // if there a flag defined via command line, the slice will be cleared first.
}

// NewSlice new func
func NewSlice(p *[]KType) *Slice {
	return &Slice{
		s:   p,
		set: false,
	}
}

// Set 解析时由FlagSet设定而来，进行解析
func (s *Slice) Set(str string) error {
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := ParseKeyFunc(v)
		if err != nil {
			return err
		}
		if !s.set {
			*s.s = (*s.s)[:0]
			s.set = true
		}
		*s.s = append(*s.s, got)
	}
	return nil
}

// Get 返回数据
func (s *Slice) Get() interface{} {
	return []KType(*s.s)
}

// TypeName type name for vars FlagValue provider
func (s *Slice) TypeName() string { return typeNameSlice }

// String 获取Set设置的字符串数据？或数据转换到的？
func (s *Slice) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

// Usage  usage info for FlagSet
func (s *Slice) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, value split by %s", StringValueDelim)
}
