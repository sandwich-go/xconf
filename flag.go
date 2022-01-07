package xconf

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

// 使用xflag解析数据返回map[string]interface{}
// structPtr只是提供结构信息供xflag进行参数定义，不涉及数据解析
// validFieldPath合法的fieldPath信息
// args带解析的参数数据获取方法
// opts为xflag附加参数
func xflagMapstructure(
	structPtr interface{},
	validFieldPath []string,
	args func(xf *xflag.Maker) []string,
	opts ...xflag.Option) (map[string]interface{}, error) {

	emptyStructPtr := reflect.New(reflect.ValueOf(structPtr).Type().Elem()).Interface()
	xf := xflag.NewMaker(opts...)
	xf.FlagSet().Usage = func() {
		xflag.PrintDefaults(xf.FlagSet())
	} // do not print usage when error

	if err := xf.Set(emptyStructPtr); err != nil {
		return nil, fmt.Errorf("got error while xflag Set, err :%w ", err)
	}
	dataArgs := args(xf)
	err := xf.Parse(dataArgs)
	if err != nil {
		return nil, fmt.Errorf("got error while parse args, err :%w ", err)
	}
	data, err := castFlagSetToMapInterface(xf.FlagSet(), validFieldPath)
	if err != nil {
		return nil, fmt.Errorf("got error while cast flag to mapinterface,err :%w ", err)
	}
	return data, nil

}

func castFlagSetToMapInterface(fs *flag.FlagSet, keys []string) (ret map[string]interface{}, err error) {
	ret = make(map[string]interface{})
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("cast panic with %v", reason)
		}
	}()
	fs.Visit(func(f *flag.Flag) {
		arr := strings.Split(f.Name, DefaultKeyDelim)
		if !xutil.ContainString(keys, f.Name) {
			return
		}
		if len(arr) == 1 {
			if v, ok := f.Value.(interface{ Get() interface{} }); ok {
				ret[f.Name] = v.Get()
			} else {
				panic(fmt.Errorf("got error while cast flag to mapstructure, field:%s Value not interface{ Get() interface{} }", f.Name))
			}
		} else {
			lastIndex := len(arr) - 1
			var val map[string]interface{}
			valInterface, ok := ret[arr[0]]
			if ok {
				val = valInterface.(map[string]interface{})
			} else {
				val = make(map[string]interface{})
			}
			valLast := val
			for i := 1; i < len(arr); i++ {
				keyKey := arr[i]
				if lastIndex == i {
					if v, ok := f.Value.(interface{ Get() interface{} }); ok {
						valLast[keyKey] = v.Get()
					} else {
						panic(fmt.Errorf("got error while cast flag to mapstructure, field:%s Value not interface{ Get() interface{} }", f.Name))
					}
				} else {
					curr := make(map[string]interface{})
					if valInterface, ok := val[arr[i]]; ok {
						curr = valInterface.(map[string]interface{})
					}
					valLast[keyKey] = curr
					valLast = curr
				}
			}
			ret[arr[0]] = val
		}
	})
	return ret, nil
}
