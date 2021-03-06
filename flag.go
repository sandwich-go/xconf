package xconf

import (
	"flag"
	"fmt"
	"strings"

	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

// ErrHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.
var ErrHelp = flag.ErrHelp

// IsErrHelp 检查错误是否是ErrHelp
func IsErrHelp(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), flag.ErrHelp.Error())
}

// newFlagSetContinueOnError 新建flagset，设定错误类型为ContinueOnError
func newFlagSetContinueOnError(name string) *flag.FlagSet {
	f := flag.NewFlagSet(name, flag.ContinueOnError)
	f.Usage = func() { xflag.PrintDefaults(f) }
	return f
}

// 使用xflag解析数据返回map[string]interface{}
// structPtr只是提供结构信息供xflag进行参数定义，不涉及数据解析
// validFieldPath合法的conf关心的fieldPath信息
// args带解析的参数数据获取方法
// opts为xflag附加参数
func xflagMapstructure(
	structPtr interface{},
	validFieldPath []string,
	args func(xf *xflag.Maker) []string,
	opts ...xflag.Option) (map[string]interface{}, error) {

	xf := xflag.NewMaker(opts...)
	if err := xf.Set(structPtr); err != nil {
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
