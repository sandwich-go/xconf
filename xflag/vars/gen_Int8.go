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
type Int8 int8

var typeNameInt8 = ""

func init() {
	var ss int8
	typeNameInt8 = reflect.TypeOf(ss).Name()
}

// NewVar new
func NewInt8(p *int8) *Int8 { return (*Int8)(p) }

// Set for each of the types
func (f *Int8) Set(s string) error {
	v, err := parseInt8(s)
	if err != nil {
		return err
	}
	*f = Int8(v)
	return nil
}

// TypeName 类型名称
func (f *Int8) TypeName() string { return typeNameInt8 }

// Get 返回类型值
func (f *Int8) Get() interface{} { return int8(*f) }

// String 获取Set设置的字符串数据？或数据转换到的？
func (f *Int8) String() string { return fmt.Sprintf("%v", *f) }

// Usage FlagSet使用
func (f *Int8) Usage() string { return "xconf/xflag/vars" }
