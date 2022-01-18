package xconf

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/sandwich-go/xconf/xutil"
)

// 数据覆盖：REMOTE > ENV > FLAG > READER > FILES > DEFAULT

// DefaultTagName 默认读取的tag名
var DefaultTagName = "xconf"

// DefaultValueTagName default value 默认读取的tag名
var DefaultValueTagName = "default"

// DefaultKeyDelim 默认的FilePath分割符
var DefaultKeyDelim = "."

// LogFunc 日志方法
type LogFunc = func(string)

// FieldTagConvertor filed名称转换方法
type FieldTagConvertor = func(fieldName string) string

var optionUsage = `
xconf : https://github.com/sandwich-go/xconf
optiongen: https://github.com/timestee/optiongen
xconf-providers: https://github.com/sandwich-go/xconf-providers
`

// OptionsOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=false --xconf=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OptionUsage": string(optionUsage),
		// annotation@NewFunc(comment="Parse时会由指定的File中加载配置")
		"Files": []string([]string{}),
		// annotation@Readers(comment="Parse时会由指定的Reader中加载配置")
		"Readers": []io.Reader([]io.Reader{}),
		// annotation@FlagSet(comment="Parse使用的FlagSet，xconf会自动在flag中创建字段定义,如指定为空则不会创建")
		"FlagSet": (*flag.FlagSet)(flag.CommandLine),
		// annotation@FlagArgs(comment="FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑")
		"FlagArgs": []string(os.Args[1:]),
		// annotation@Environ(comment="Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.Usage()获取当前支持的FlagSet与Env参数定义")
		"Environ": []string(os.Environ()),
		// annotation@ErrorHandling(comment="错误处理模式")
		"ErrorHandling": (ErrorHandling)(PanicOnError),
		// annotation@TagName(comment="xconf使用的字段TAG名称,默认:xconf")
		"TagName": string(DefaultTagName),
		// annotation@DecoderConfigOption(comment="xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure")
		"DecoderConfigOption": []DecoderConfigOption(nil),
		// annotation@MapMerge(comment="map是否开启merge模式，详情见文档")
		"MapMerge": false,
		// annotation@FieldTagConvertor(comment="字段名转换到FiledPath时优先使用TagName指定的名称，否则使用该函数转换")
		"FieldTagConvertor": FieldTagConvertor(xutil.SnakeCase),
		// annotation@FieldPathRemoved(comment="弃用的配置，目标结构中已经删除，但配置文件中可能存在，解析时不会认为是错误，会将该配置丢弃，并打印WARNING日志")
		"FieldPathRemoved": []string{},
		// annotation@Debug(comment="debug模式下输出调试信息")
		"Debug": false,
		// annotation@LogDebug(comment="DEBUG日志")
		"LogDebug": LogFunc(func(s string) { log.Println("[  DEBUG] " + s) }),
		// annotation@LogWarning(comment="WARNING日志")
		"LogWarning": LogFunc(func(s string) { log.Println("[WARNING] " + s) }),
		// annotation@AppLabelList(comment="应用层Label，用于灰度发布场景")
		"AppLabelList": []string([]string{}),
		// annotation@EnvBindShouldErrorWhenFailed(comment="EnvBind时如果Env中不存在指定的key而且没有指定默认值时是否返回错误")
		"EnvBindShouldErrorWhenFailed": true,
		// annotation@FlagCreateIgnoreFiledPath(comment="不创建到FlagSet中的字段FieldPath")
		// todo: 可以通过tag中指定flagoff规避这个字段的支持
		"FlagCreateIgnoreFiledPath": []string{},
		// annotation@ParseDefault(comment="是否解析struct标签中的default数据，解析规则参考xflag支持")
		"ParseDefault": true,
		// annotation@TagNameForDefaultValue(comment="默认值TAG名称,默认default")
		"TagNameForDefaultValue": string(DefaultValueTagName),
		// annotation@ReplaceFlagSetUsage(comment="是否替换FlagSet的Usage，使用xconf内置版本")
		"ReplaceFlagSetUsage": true,
	}
}

func init() {
	InstallOptionsWatchDog(func(cc *Options) {
		if cc.MapMerge {
			cc.LogWarning("Map Merge Model Enabled.")
		}
		if len(cc.AppLabelList) == 0 {
			hostName, _ := os.Hostname()
			cc.AppLabelList = append(cc.AppLabelList, hostName)
		}
	})
}
