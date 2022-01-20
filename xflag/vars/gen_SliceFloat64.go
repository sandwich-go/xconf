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

var typeNameSliceFloat64 = ""

// StringValueDelim 数据分割符

func init() {
	v := []float64{}
	typeNameSliceFloat64 = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceFloat64, func(valPtr interface{}, stringAlias func(s string) string) flag.Getter {
		return NewSliceFloat64(valPtr.(*[]float64), stringAlias)
	})
}

// Slice struct
type SliceFloat64 struct {
	stringAlias func(s string) string
	s           *[]float64
	set         bool // if there a flag defined via command line, the slice will be cleared first.
}

// NewSlice new func
func NewSliceFloat64(p *[]float64, stringAlias func(s string) string) *SliceFloat64 {
	return &SliceFloat64{
		stringAlias: stringAlias,
		s:           p,
		set:         false,
	}
}

// Set 解析时由FlagSet设定而来，进行解析
func (s *SliceFloat64) Set(str string) error {
	str = s.stringAlias(str)
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseFloat64(s.stringAlias(v))
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
func (s *SliceFloat64) Get() interface{} {
	return []float64(*s.s)
}

// TypeName type name for vars FlagValue provider
func (s *SliceFloat64) TypeName() string { return typeNameSliceFloat64 }

// String 获取Set设置的字符串数据？或数据转换到的？
func (s *SliceFloat64) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

// Usage  usage info for FlagSet
func (s *SliceFloat64) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, value split by %s", StringValueDelim)
}
