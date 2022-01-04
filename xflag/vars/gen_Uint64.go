// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

type Uint64 uint64

func NewUint64(p *uint64) *Uint64 { return (*Uint64)(p) }

// Setters for each of the types
func (f *Uint64) Set(s string) error {
	v, err := parseUint64(s)
	if err != nil {
		return err
	}
	*f = Uint64(v)
	return nil
}
func (f *Uint64) Get() interface{} { return uint64(*f) }
func (f *Uint64) String() string   { return fmt.Sprintf("%v", *f) }
func (f *Uint64) Usage() string    { return fmt.Sprintf("xconf/xflag/vars %s", reflect.TypeOf(*f)) }