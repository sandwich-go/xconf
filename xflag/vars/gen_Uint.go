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
type Uint uint

var typeNameUint = ""

func init() {
	var ss uint
	typeNameUint = reflect.TypeOf(ss).Name()
}

// NewVar new func
func NewUint(p *uint) *Uint { return (*Uint)(p) }

// Set for each of the types
func (f *Uint) Set(s string) error {
	v, err := parseUint(s)
	if err != nil {
		return err
	}
	*f = Uint(v)
	return nil
}

// TypeName 类型名称
func (f *Uint) TypeName() string { return typeNameUint }

// Get 返回类型值
func (f *Uint) Get() interface{} { return uint(*f) }

// String 获取Set设置的字符串数据？或数据转换到的？
func (f *Uint) String() string { return fmt.Sprintf("%v", *f) }

// Usage FlagSet使用
func (f *Uint) Usage() string { return "xconf/xflag/vars" }
