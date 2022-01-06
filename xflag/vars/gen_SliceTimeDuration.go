// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"
)

//template type Slice(KType,ParseKeyFunc,SetProviderByFieldType,StringValueDelim)
// KType 默认key类型

// SetProviderByFieldType 替换

// ParseKeyFunc val解析，替换

var typeNameSliceTimeDuration = ""

func init() {
	v := []time.Duration{}
	typeNameSliceTimeDuration = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceTimeDuration, func(valPtr interface{}) flag.Getter {
		return NewSliceTimeDuration(valPtr.(*[]time.Duration))
	})
}

// Slice
type SliceTimeDuration struct {
	s   *[]time.Duration
	set bool // if there a flag defined via command line, the slice will be cleared first.
}

// NewSlice 创建指定类型
func NewSliceTimeDuration(p *[]time.Duration) *SliceTimeDuration {
	return &SliceTimeDuration{
		s:   p,
		set: false,
	}
}

// Set 解析时由FlagSet设定而来，进行解析
func (s *SliceTimeDuration) Set(str string) error {
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseTimeDuration(v)
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
func (s *SliceTimeDuration) Get() interface{} {
	return []time.Duration(*s.s)
}

// TypeName type name for vars FlagValue provider
func (e *SliceTimeDuration) TypeName() string { return typeNameSliceTimeDuration }

// String 获取Set设置的字符串数据？或数据转换到的？
func (s *SliceTimeDuration) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

// Usage  usage info for FlagSet
func (s *SliceTimeDuration) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, value split by %s", StringValueDelim)
}
