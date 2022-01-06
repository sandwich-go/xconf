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

// Var type
type Var KType

var typeNameVar = ""

func init() {
	var ss KType
	typeNameVar = reflect.TypeOf(ss).Name()
}

// NewVar new func
func NewVar(p *KType) *Var { return (*Var)(p) }

// Set for each of the types
func (f *Var) Set(s string) error {
	v, err := ParseKeyFunc(s)
	if err != nil {
		return err
	}
	*f = Var(v)
	return nil
}

// TypeName 类型名称
func (f *Var) TypeName() string { return typeNameVar }

// Get 返回类型值
func (f *Var) Get() interface{} { return KType(*f) }

// String 获取Set设置的字符串数据？或数据转换到的？
func (f *Var) String() string { return fmt.Sprintf("%v", *f) }

// Usage FlagSet使用
func (f *Var) Usage() string { return "xconf/xflag/vars" }
