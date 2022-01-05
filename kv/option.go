package kv

import (
	"github.com/sandwich-go/xconf/secconf"
)

type ContentChange = func(string, string, []byte)
type WatchError = func(string, string, error)

//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OnWatchError": WatchError(nil),
		"Decoder":      secconf.Codec(nil), // @MethodComment(允许每一个远端设定独立的加密方式)
	}
}
