package xvar

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

// ParseKeyFunc 默认解析函数，替换
var ParseKeyFunc = func(s string) (KType, error) { panic(1) }

// KType 默认类型，替换
type KType int

// Var
type Var KType

var typeNameVar = ""

func init() {
	var ss KType
	typeNameVar = reflect.TypeOf(ss).Name()
}

// NewVar new
func NewVar(p *KType) *Var { return (*Var)(p) }

// Setters for each of the types
func (f *Var) Set(s string) error {
	v, err := ParseKeyFunc(s)
	if err != nil {
		return err
	}
	*f = Var(v)
	return nil
}
func (f *Var) TypeName() string { return typeNameVar }
func (f *Var) Get() interface{} { return KType(*f) }
func (f *Var) String() string   { return fmt.Sprintf("%v", *f) }

func (f *Var) Usage() string { return "xconf/xflag/vars" }
