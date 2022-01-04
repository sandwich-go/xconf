// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

type Uint uint

func NewUint(p *uint) *Uint { return (*Uint)(p) }

// Setters for each of the types
func (f *Uint) Set(s string) error {
	v, err := parseUint(s)
	if err != nil {
		return err
	}
	*f = Uint(v)
	return nil
}
func (f *Uint) Get() interface{} { return uint(*f) }
func (f *Uint) String() string   { return fmt.Sprintf("%v", *f) }
func (f *Uint) Usage() string    { return fmt.Sprintf("xconf/xflag/vars %s", reflect.TypeOf(*f)) }