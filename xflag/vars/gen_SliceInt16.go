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

var typeNameSliceInt16 = ""

// StringValueDelim 数据分割符

func init() {
	v := []int16{}
	typeNameSliceInt16 = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceInt16, func(valPtr interface{}, stringAlias func(s string) string) flag.Getter {
		return NewSliceInt16(valPtr.(*[]int16), stringAlias)
	})
}

// Slice struct
type SliceInt16 struct {
	stringAlias func(s string) string
	s           *[]int16
	set         bool // if there a flag defined via command line, the slice will be cleared first.
}

// NewSlice new func
func NewSliceInt16(p *[]int16, stringAlias func(s string) string) *SliceInt16 {
	return &SliceInt16{
		stringAlias: stringAlias,
		s:           p,
		set:         false,
	}
}

// Set 解析时由FlagSet设定而来，进行解析
func (s *SliceInt16) Set(str string) error {
	str = s.stringAlias(str)
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseInt16(s.stringAlias(v))
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
func (s *SliceInt16) Get() interface{} {
	return []int16(*s.s)
}

// TypeName type name for vars FlagValue provider
func (s *SliceInt16) TypeName() string { return typeNameSliceInt16 }

// String 获取Set设置的字符串数据？或数据转换到的？
func (s *SliceInt16) String() string {
	if s.s == nil || len(*s.s) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

// Usage  usage info for FlagSet
func (s *SliceInt16) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, value split by %s", StringValueDelim)
}
