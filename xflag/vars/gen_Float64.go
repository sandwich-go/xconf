// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

// ParseKeyFunc 默认解析函数，替换

// KType 默认类型，替换

// Var 类型
type Float64 float64

var typeNameFloat64 = ""

func init() {
	var ss float64
	typeNameFloat64 = reflect.TypeOf(ss).Name()
}
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
func (f *Float64) TypeName() string { return typeNameFloat64 }
func (f *Float64) Get() interface{} { return float64(*f) }
func (f *Float64) String() string   { return fmt.Sprintf("%v", *f) }
func (f *Float64) Usage() string    { return "xconf/xflag/vars" }
