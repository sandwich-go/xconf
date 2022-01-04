// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type Slice(KType,ParseKeyFunc,SetProviderByFieldType,StringValueDelim)

var typeNameSliceUint16 = ""

func init() {
	v := []uint16{}
	typeNameSliceUint16 = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceUint16, func(valPtr interface{}) flag.Getter {
		return NewSliceUint16(valPtr.(*[]uint16))
	})
}

type SliceUint16 struct {
	s   *[]uint16
	set bool // if there a flag defined via command line, the slice will be cleared first.
}

func NewSliceUint16(p *[]uint16) *SliceUint16 {
	return &SliceUint16{
		s:   p,
		set: false,
	}
}

func (s *SliceUint16) Set(str string) error {
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseUint16(v)
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

func (s *SliceUint16) Get() interface{} {
	return []uint16(*s.s)
}

func (s *SliceUint16) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

func (s *SliceUint16) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars %s,v%sv%sv", typeNameSliceUint16, StringValueDelim, StringValueDelim)
}