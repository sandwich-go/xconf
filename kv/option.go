package kv

import (
	"github.com/sandwich-go/xconf/secconf"
)

// ContentChange kv数据发生变化时回调
type ContentChange = func(loaderName string, confPath string, content []byte) error

// WatchError kv.Loader.Watch发生错误时回调
type WatchError = func(loaderName string, confPath string, watchErr error)

// OptionsOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OnWatchError": WatchError(nil),
		"Decoder":      secconf.Codec(nil), // @MethodComment(允许每一个远端设定独立的加密方式)
	}
}
