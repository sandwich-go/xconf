// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"fmt"
	"reflect"
)

//template type Var(KType,ParseKeyFunc)

// ParseKeyFunc 默认解析函数，替换

// KType 默认类型，替换

// Var type
type Int64 struct {
	v           *int64
	stringAlias func(s string) string
}

var typeNameInt64 = ""

func init() {
	var ss int64
	typeNameInt64 = reflect.TypeOf(ss).Name()
}

// NewVar new func
func NewInt64(p *int64, stringAlias func(s string) string) *Int64 {
	return &Int64{
		v:           p,
		stringAlias: stringAlias,
	}
}

// Set for each of the types
func (f *Int64) Set(s string) error {
	v, err := parseInt64(f.stringAlias(s))
	if err != nil {
		return err
	}
	*f.v = v
	return nil
}

// TypeName 类型名称
func (f *Int64) TypeName() string { return typeNameInt64 }

// Get 返回类型值
func (f *Int64) Get() interface{} { return *f.v }

// String 获取Set设置的字符串数据？或数据转换到的？
func (f *Int64) String() string { return fmt.Sprintf("%v", *f.v) }

// Usage FlagSet使用
func (f *Int64) Usage() string { return "xconf/xflag/vars" }
