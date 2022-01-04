// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type Slice(KType,ParseKeyFunc,SetProviderByFieldType,StringValueDelim)

var typeNameSliceUint = ""

func init() {
	v := []uint{}
	typeNameSliceUint = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSliceUint, func(valPtr interface{}) flag.Getter {
		return NewSliceUint(valPtr.(*[]uint))
	})
}

type SliceUint struct {
	s   *[]uint
	set bool // if there a flag defined via command line, the slice will be cleared first.
}

func NewSliceUint(p *[]uint) *SliceUint {
	return &SliceUint{
		s:   p,
		set: false,
	}
}

func (s *SliceUint) Set(str string) error {
	for _, v := range strings.Split(str, StringValueDelim) {
		got, err := parseUint(v)
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

func (s *SliceUint) Get() interface{} {
	return []uint(*s.s)
}

func (s *SliceUint) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

func (s *SliceUint) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars %s,v%sv%sv", typeNameSliceUint, StringValueDelim, StringValueDelim)
}