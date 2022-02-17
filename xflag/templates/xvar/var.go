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
type Var struct {
	v           *KType
	stringAlias func(s string) string
}

var typeNameVar = ""

func init() {
	var ss KType
	typeNameVar = reflect.TypeOf(ss).Name()
}

// NewVar new func
func NewVar(p *KType, stringAlias func(s string) string) *Var {
	return &Var{
		v:           p,
		stringAlias: stringAlias,
	}
}

// Set for each of the types
func (f *Var) Set(s string) error {
	v, err := ParseKeyFunc(f.stringAlias(s))
	if err != nil {
		return err
	}
	*f.v = v
	return nil
}

// TypeName 类型名称
func (f *Var) TypeName() string { return typeNameVar }

// Get 返回类型值
func (f *Var) Get() interface{} { return *f.v }

// String 获取Set设置的字符串数据？或数据转换到的？
func (f *Var) String() string {
	if f.v == nil {
		return ""
	}
	return fmt.Sprintf("%v", *f.v)
}

// Usage FlagSet使用
func (f *Var) Usage() string { return "xconf/xflag/vars" }

// IsBoolFlag IsBoolFlag
func (f *Var) IsBoolFlag() bool { return typeNameVar == "bool" }
