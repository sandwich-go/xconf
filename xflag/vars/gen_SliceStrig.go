// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type Slice(KType,ParseKeyFunc,SetProviderByFieldType,StringValueDelim)

// KType 默认key类型

// SetProviderByFieldType 替换

// ParseKeyFunc val解析，替换

var typeNameSliceStrig = ""

// StringValueDelim 数据分割符

func init() {
	v := []string{}
	typeNameSliceStrig = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceStrig, func(valPtr interface{}, stringAlias func(s string) string) flag.Getter {
		return NewSliceStrig(valPtr.(*[]string), stringAlias)
	})
}

// Slice struct
type SliceStrig struct {
	stringAlias func(s string) string
	s           *[]string
	set         bool // if there a flag defined via command line, the slice will be cleared first.
}

// NewSlice new func
func NewSliceStrig(p *[]string, stringAlias func(s string) string) *SliceStrig {
	return &SliceStrig{
		stringAlias: stringAlias,
		s:           p,
		set:         false,
	}
}

// Set 解析时由FlagSet设定而来，进行解析
func (s *SliceStrig) Set(str string) error {
	str = s.stringAlias(str)
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseString(s.stringAlias(v))
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
func (s *SliceStrig) Get() interface{} {
	return []string(*s.s)
}

// TypeName type name for vars FlagValue provider
func (s *SliceStrig) TypeName() string { return typeNameSliceStrig }

// String 获取Set设置的字符串数据？或数据转换到的？
func (s *SliceStrig) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

// Usage  usage info for FlagSet
func (s *SliceStrig) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, value split by %s", StringValueDelim)
}
