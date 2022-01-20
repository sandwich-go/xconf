package xflag

import (
	"flag"
	"log"
	"strings"

	"github.com/sandwich-go/xconf/xflag/vars"
)

// KeyFormat key格式，当指定的tag无val时调用改方法格式化
type KeyFormat func(string) string

// LogFunc 日志方法
type LogFunc func(string)

// OptionsOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Name":                      "",
		"TagName":                   "xconf", // 使用的tag key,如不设定则使用
		"UsageTagName":              "usage",
		"Flatten":                   false, // 是否使用扁平模式，不使用.分割
		"FlagSet":                   (*flag.FlagSet)(flag.NewFlagSet("flagmaker", flag.ContinueOnError)),
		"FlagValueProvider":         vars.FlagValueProvider(vars.DefaultFlagValueProvider),
		"KeyFormat":                 KeyFormat(func(s string) string { return strings.ToLower(s) }),
		"FlagCreateIgnoreFiledPath": []string{},
		"LogDebug":                  LogFunc(func(s string) { log.Print("debug:" + s) }),
		"LogWarning":                LogFunc(func(s string) { log.Print("warning: " + s) }),
		"StringAlias":               func(s string) string { return s },
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
		if cc.StringAlias == nil {
			cc.StringAlias = func(s string) string { return s }
		}
	})
}
