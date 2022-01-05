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

var typeNameMapInt64Int64 = ""

func init() {
	v := map[int64]int64{}
	typeNameMapInt64Int64 = fmt.Sprintf("map[%s]%s", reflect.TypeOf(v).Key().Name(), reflect.TypeOf(v).Elem().Name())
	SetProviderByFieldType(typeNameMapInt64Int64, func(valPtr interface{}) flag.Getter {
		return NewMapInt64Int64(valPtr)
	})
}

type MapInt64Int64 struct {
	s   string
	set bool
	val *map[int64]int64
}

func NewMapInt64Int64(valPtr interface{}) *MapInt64Int64 {
	return &MapInt64Int64{
		val: valPtr.(*map[int64]int64),
	}
}

func (e *MapInt64Int64) Get() interface{} {
	vv := make(map[string]interface{})
	for k, v := range *e.val {
		vv[fmt.Sprintf("%v", k)] = v
	}
	return vv
}

func (e *MapInt64Int64) Usage() string {
	return fmt.Sprintf("xconf/xflag/vars, key and value split by %s", StringValueDelim)
}

func (e *MapInt64Int64) TypeName() string { return typeNameMapInt64Int64 }
func (e *MapInt64Int64) String() string   { return e.s }
func (e *MapInt64Int64) Set(s string) error {
	e.s = s
	kv := strings.Split(s, StringValueDelim)
	if len(kv)%2 == 1 {
		return errors.New("got the odd number of input pairs")
	}
	if !e.set {
		e.set = true
		*e.val = make(map[int64]int64)
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
		val, err := parseInt64(s)
		if err != nil {
			return err
		}
		(*e.val)[keyVal] = val
	}
	return nil
}
