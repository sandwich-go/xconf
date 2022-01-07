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

var typeNameSliceUint8 = ""

func init() {
	v := []uint8{}
	typeNameSliceUint8 = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceUint8, func(valPtr interface{}) flag.Getter {
		return NewSliceUint8(valPtr.(*[]uint8))
	})
}

// Slice struct
type SliceUint8 struct {
	s   *[]uint8
	set bool // if there a flag defined via command line, the slice will be cleared first.
}

// NewSlice new func
func NewSliceUint8(p *[]uint8) *SliceUint8 {
	return &SliceUint8{
		s:   p,
		set: false,
	}
}

// Set 解析时由FlagSet设定而来，进行解析
func (s *SliceUint8) Set(str string) error {
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseUint8(v)
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
func (s *SliceUint8) Get() interface{} {
	return []uint8(*s.s)
}

// TypeName type name for vars FlagValue provider
func (e *SliceUint8) TypeName() string { return typeNameSliceUint8 }

// String 获取Set设置的字符串数据？或数据转换到的？
func (s *SliceUint8) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

// Usage  usage info for FlagSet
func (s *SliceUint8) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, value split by %s", StringValueDelim)
}
