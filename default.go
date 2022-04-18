package xconf

import (
	"fmt"
	"reflect"

	"github.com/sandwich-go/xconf/xflag"
)

// ParseDefault 根据opts指定的TagNameDefaultValue解析字段默认值并绑定到valPtr
func ParseDefault(valPtr interface{}, opts ...Option) (err error) {
	xd := NewWithoutFlagEnv(opts...)
	applyXConfOptions(valPtr, xd)
	data, _, err := xd.parseDefault(valPtr)
	if err != nil {
		return fmt.Errorf("got error:%w while parse default", err)
	}
	if len(data) == 0 {
		return nil
	}
	if err = xd.decode(data, valPtr); err != nil {
		return fmt.Errorf("got error:%w while decode using map structure", err)
	}
	return nil
}

func (x *XConf) parseDefault(valPtr interface{}) (data map[string]interface{}, parsed bool, err error) {
	var flagVals []string
	var keys []string
	_, fieldInfo := NewStruct(
		reflect.New(reflect.ValueOf(valPtr).Type().Elem()).Interface(),
		x.cc.TagName,
		x.cc.TagNameForDefaultValue,
		x.cc.FieldTagConvertor).Map()
	for k, v := range fieldInfo {
		keys = append(keys, k)
		if v.DefaultGot {
			flagVals = append(flagVals, fmt.Sprintf("--%s=%s", k, v.DefaultString))
		}
	}
	if len(flagVals) == 0 {
		return nil, false, nil
	}
	xflagOpts := x.defaultXFlagOptions()
	data, err = xflagMapstructure(valPtr, keys,
		func(xf *xflag.Maker) []string {
			return flagVals
		},
		append(xflagOpts, xflag.WithFlagSet(newFlagSetContinueOnError("Default")))...)
	if err != nil {
		return data, true, fmt.Errorf("got error while xflag for default, err :%w", err)
	}

	return data, true, nil
}
