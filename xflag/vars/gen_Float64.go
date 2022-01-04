// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

type Float64 float64

func NewFloat64(p *float64) *Float64 { return (*Float64)(p) }

// Setters for each of the types
func (f *Float64) Set(s string) error {
	v, err := parseFloat64(s)
	if err != nil {
		return err
	}
	*f = Float64(v)
	return nil
}
func (f *Float64) Get() interface{} { return float64(*f) }
func (f *Float64) String() string   { return fmt.Sprintf("%v", *f) }
func (f *Float64) Usage() string    { return fmt.Sprintf("xconf/xflag/vars %s", reflect.TypeOf(*f)) }