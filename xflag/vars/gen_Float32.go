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
type Float32 struct {
	v           *float32
	stringAlias func(s string) string
}

var typeNameFloat32 = ""

func init() {
	var ss float32
	typeNameFloat32 = reflect.TypeOf(ss).Name()
}

// NewVar new func
func NewFloat32(p *float32, stringAlias func(s string) string) *Float32 {
	return &Float32{
		v:           p,
		stringAlias: stringAlias,
	}
}

// Set for each of the types
func (f *Float32) Set(s string) error {
	v, err := parseFloat32(f.stringAlias(s))
	if err != nil {
		return err
	}
	*f.v = v
	return nil
}

// TypeName 类型名称
func (f *Float32) TypeName() string { return typeNameFloat32 }

// Get 返回类型值
func (f *Float32) Get() interface{} { return *f.v }

// String 获取Set设置的字符串数据？或数据转换到的？
func (f *Float32) String() string {
	if f.v == nil {
		return ""
	}
	return fmt.Sprintf("%v", *f.v)
}

// Usage FlagSet使用
func (f *Float32) Usage() string { return "xconf/xflag/vars" }

// IsBoolFlag IsBoolFlag
func (f *Float32) IsBoolFlag() bool { return typeNameFloat32 == "bool" }
