package xflag

import (
	"flag"
	"log"
	"strings"

	"github.com/sandwich-go/xconf/xflag/vars"
)

type KeyFormat func(string) string
type LogFunc func(string)

//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Name":              "",
		"TagName":           "cfg", // 使用的tag key,如不设定则使用
		"Flatten":           false, // 是否使用扁平模式，不使用.分割
		"FlagSet":           (*flag.FlagSet)(flag.NewFlagSet("flagmaker", flag.ContinueOnError)),
		"FlagValueProvider": vars.FlagValueProvider(vars.DefaultFlagValueProvider),
		"KeyFormat":         KeyFormat(func(s string) string { return strings.ToLower(s) }),
		"FlagSetIgnore":     []string{},
		"LogDebug": LogFunc(func(s string) {
			log.Print("debug:" + s)
		}),
		"LogWarning": LogFunc(func(s string) {
			log.Print("warning: " + s)
		}),
	}
}

func init() {
	InstallOptionsWatchDog(func(cc *Options) {
		if cc.Name == "" {
			if cc.FlagSet == flag.CommandLine {
				cc.Name = "CommandLine"
			} else {
				cc.Name = cc.FlagSet.Name()
			}
		}
	})
}
