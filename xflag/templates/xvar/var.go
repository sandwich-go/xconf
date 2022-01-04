package xvar

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

var ParseKeyFunc = func(s string) (KType, error) { panic(1) }

type KType int
type Var KType

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
func (f *Var) Get() interface{} { return KType(*f) }
func (f *Var) String() string   { return fmt.Sprintf("%v", *f) }
func (f *Var) Usage() string    { return fmt.Sprintf("xconf/xflag/vars %s", reflect.TypeOf(*f)) }
