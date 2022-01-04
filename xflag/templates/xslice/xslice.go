package vars

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type Slice(KType,ParseKeyFunc,SetProviderByFieldType,StringValueDelim)

type KType int

var SetProviderByFieldType = func(v interface{}, flagValue interface{}) { panic(1) }
var ParseKeyFunc = func(s string) (KType, error) { panic(1) }
var typeNameSlice = ""
var StringValueDelim = ""

func init() {
	v := []KType{}
	typeNameSlice = fmt.Sprintf("[]%s", reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameSlice, func(valPtr interface{}) flag.Getter {
		return NewSlice(valPtr.(*[]KType))
	})
}

type Slice struct {
	s   *[]KType
	set bool // if there a flag defined via command line, the slice will be cleared first.
}

func NewSlice(p *[]KType) *Slice {
	return &Slice{
		s:   p,
		set: false,
	}
}

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

func (s *Slice) Get() interface{} {
	return []KType(*s.s)
}

func (s *Slice) String() string {
	if s.s == nil {
		return ""
	}
	return fmt.Sprintf("%v", *s.s)
}

func (s *Slice) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars %s,v%sv%sv", typeNameSlice, StringValueDelim, StringValueDelim)
}
