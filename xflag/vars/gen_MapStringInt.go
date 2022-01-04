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

var typeNameMapStringInt = ""

func init() {
	v := map[string]int{}
	typeNameMapStringInt = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapStringInt, func(valPtr interface{}) flag.Getter {
		return NewMapStringInt(valPtr)
	})
}

type MapStringInt struct {
	s   string
	set bool
	val *map[string]int
}

func NewMapStringInt(valPtr interface{}) *MapStringInt {
	return &MapStringInt{
		val: valPtr.(*map[string]int),
	}
}

func (e *MapStringInt) Get() interface{} {
	vv := make(map[string]interface{})
	for k, v := range *e.val {
		vv[fmt.Sprintf("%v", k)] = v
	}
	return vv
}
func (e *MapStringInt) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars %s,k%sv%sk%sv", typeNameMapStringInt, StringValueDelim, StringValueDelim, StringValueDelim)
}
func (e *MapStringInt) String() string { return e.s }
func (e *MapStringInt) Set(s string) error {
	e.s = s
	kv := strings.Split(s, StringValueDelim)
	if len(kv)%2 == 1 {
		return errors.New("got the odd number of input pairs")
	}
	if !e.set {
		e.set = true
		*e.val = make(map[string]int)
	}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s
			continue
		}
		keyVal, err := parseString(key)
		if err != nil {
			return err
		}
		val, err := parseInt(s)
		if err != nil {
			return err
		}
		(*e.val)[keyVal] = val
	}
	return nil
}