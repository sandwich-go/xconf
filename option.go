package xconf

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/sandwich-go/xconf/xflag/vars"
)

var DefaultTagName = "xconf"
var DefaultValueTagName = "default"
var DefaultKeyDelim = "."

// AutoOptions provide flag: _auto_conf_files_
type AutoOptions struct {
	AutoConfFiles string `flag:"_auto_conf_files_"`
}

type LogFunc = func(string)
type FieldTagConvertor = func(fieldName string) string

// 数据覆盖：REMOTE > ENV > FLAG > READER > FILES > DEFAULT
//go:generate optiongen --option_with_struct_name=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Files":                            []string([]string{}),                                      // @MethodComment(Parse时会由指定的File中加载配置)
		"Readers":                          []io.Reader([]io.Reader{}),                                // @MethodComment(Parse时会由指定的Reader中加载配置)
		"FlagSet":                          (*flag.FlagSet)(flag.CommandLine),                         // @MethodComment(Parse使用的FlagSet，xconf会自动在flag中创建字段定义,如指定为空则不会创建)
		"FlagValueProvider":                vars.FlagValueProvider(nil),                               // @MethodComment(FlagValueProvider，当xconf无法将字段定义到FlagSet时会回调该方法，提供一些复杂参数配置的Flag与Env支持)
		"FlagArgs":                         []string(os.Args[1:]),                                     // @MethodComment(FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑)
		"Environ":                          []string(os.Environ()),                                    // @MethodComment((Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.DumpInfo()获取当前支持的FlagSet与Env参数定义)
		"DecoderConfigOption":              []DecoderConfigOption(nil),                                // @MethodComment(xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure)
		"ErrorHandling":                    (ErrorHandling)(PanicOnError),                             // @MethodComment(错误处理模式)
		"MapMerge":                         false,                                                     // @MethodComment(map是否开启merge模式，默认情况下map是作为叶子节点覆盖的，可以通过指定noleaf标签表明key级别覆盖，但是key对应的val依然是整体覆盖，如果指定MapMerge为true，则Map及子元素都会在字段属性级别进行覆盖)
		"Debug":                            false,                                                     // @MethodComment(debug模式下输出调试信息)
		"LogDebug":                         LogFunc(func(s string) { log.Println("[  DEBUG] " + s) }), // @MethodComment(DEBUG日志)
		"LogWarning":                       LogFunc(func(s string) { log.Println("[WARNING] " + s) }), // @MethodComment(WARNING日志)
		"FieldTagConvertor":                FieldTagConvertor(SnakeCase),                              // @MethodComment(字段名转换到map key时优先使用TagName指定的名称，否则使用该函数转换)
		"TagName":                          string(DefaultTagName),                                    // @MethodComment(字段TAG名称,默认xconf)
		"TagNameDefaultValue":              string(DefaultValueTagName),                               // @MethodComment(默认值TAG名称,默认default)
		"ParseDefault":                     true,                                                      // @MethodComment(是否解析struct标签中的default数据，解析规则参考xflag支持)
		"FieldPathDeprecated":              []string{},                                                // @MethodComment(弃用的配置，解析时不会报错，但会打印warning日志)
		"ErrEnvBindNotExistWithoutDefault": true,                                                      // @MethodComment(EnvBind时如果Env中不存在指定的key而且没有指定默认值时报错)
		"FieldFlagSetCreateIgnore":         []string{},                                                // @MethodComment(不自动创建到FlagSet中的名称，路径)
	}
}

func init() {
	InstallOptionsWatchDog(func(cc *Options) {
		if cc.MapMerge {
			cc.LogWarning("Map Merge Model Enabled.")
		}
	})
}
