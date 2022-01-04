// Code generated by gotemplate. DO NOT EDIT.

package vars

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

//template type MapKTypeVType(KType,VType,ParseKeyFunc,ParseValFunc,SetProviderByFieldType,StringValueDelim)

var typeNameMapInt64String = ""

func init() {
	v := map[int64]string{}
	typeNameMapInt64String = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapInt64String, func(valPtr interface{}) flag.Getter {
		return NewMapInt64String(valPtr)
	})
}

type MapInt64String struct {
	s   string
	set bool
	val *map[int64]string
}

func NewMapInt64String(valPtr interface{}) *MapInt64String {
	return &MapInt64String{
		val: valPtr.(*map[int64]string),
	}
}

func (e *MapInt64String) Get() interface{} {
	vv := make(map[string]interface{})
	for k, v := range *e.val {
		vv[fmt.Sprintf("%v", k)] = v
	}
	return vv
}
func (e *MapInt64String) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars %s,k%sv%sk%sv", typeNameMapInt64String, StringValueDelim, StringValueDelim, StringValueDelim)
}
func (e *MapInt64String) String() string { return e.s }
func (e *MapInt64String) Set(s string) error {
	e.s = s
	kv := strings.Split(s, StringValueDelim)
	if len(kv)%2 == 1 {
		return errors.New("got the odd number of input pairs")
	}
	if !e.set {
		e.set = true
		*e.val = make(map[int64]string)
	}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s
			continue
		}
		keyVal, err := parseInt64(key)
		if err != nil {
			return err
		}
		val, err := parseString(s)
		if err != nil {
			return err
		}
		(*e.val)[keyVal] = val
	}
	return nil
}