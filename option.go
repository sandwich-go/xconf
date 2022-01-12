package xconf

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xflag/vars"
	"github.com/sandwich-go/xconf/xutil"
)

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

// 数据覆盖：REMOTE > ENV > FLAG > READER > FILES > DEFAULT
// OptionsOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=false --xconf=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	// annotation@NewFunc(comment="Parse时会由指定的File中加载配置")
	// annotation@Readers(comment="Parse时会由指定的Reader中加载配置")
	// annotation@FlagSet(comment="Parse使用的FlagSet，xconf会自动在flag中创建字段定义,如指定为空则不会创建")
	// annotation@FlagValueProvider(comment="FlagValueProvider，当xconf无法将字段定义到FlagSet时会回调该方法，提供一些复杂参数配置的Flag与Env支持")
	// annotation@FlagArgs(comment="FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑")
	// annotation@Environ(comment="Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.DumpInfo()获取当前支持的FlagSet与Env参数定义")
	// annotation@DecoderConfigOption(comment="xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure")
	// annotation@MapMerge(comment="map是否开启merge模式，默认情况下map是作为叶子节点覆盖的，可以通过指定noleaf标签表明key级别覆盖，但是key对应的val依然是整体覆盖，如果指定MapMerge为true，则Map及子元素都会在字段属性级别进行覆盖")
	// annotation@ErrorHandling(comment="错误处理模式")
	// annotation@FieldTagConvertor(comment="字段名转换到map key时优先使用TagName指定的名称，否则使用该函数转换")
	// annotation@ParseDefault(comment="是否解析struct标签中的default数据，解析规则参考xflag支持")
	// annotation@FieldPathDeprecated(comment="弃用的配置，解析时不会报错，但会打印warning日志")
	// annotation@TagNameDefaultValue(comment="默认值TAG名称,默认default")
	// annotation@ErrEnvBindNotExistWithoutDefault(comment="EnvBind时如果Env中不存在指定的key而且没有指定默认值时报错")
	// annotation@TagName(comment="字段TAG名称,默认xconf")
	// annotation@FieldFlagSetCreateIgnore(comment="不自动创建到FlagSet中的名称，路径")
	// annotation@Debug(comment="debug模式下输出调试信息")
	// annotation@LogDebug(comment="DEBUG日志")
	// annotation@LogWarning(comment="WARNING日志")
	// annotation@AppLabelList(comment="应用层Label，用于灰度发布场景")
	return map[string]interface{}{
		"Files":                            []string([]string{}),
		"Readers":                          []io.Reader([]io.Reader{}),
		"FlagSet":                          (*flag.FlagSet)(flag.CommandLine),
		"FlagValueProvider":                vars.FlagValueProvider(nil),
		"FlagArgs":                         []string(os.Args[1:]),
		"Environ":                          []string(os.Environ()),
		"DecoderConfigOption":              []DecoderConfigOption(nil),
		"ErrorHandling":                    (ErrorHandling)(PanicOnError),
		"MapMerge":                         false,
		"FieldTagConvertor":                FieldTagConvertor(xutil.SnakeCase),
		"TagName":                          string(DefaultTagName),
		"TagNameDefaultValue":              string(DefaultValueTagName),
		"ParseDefault":                     true,
		"FieldPathDeprecated":              []string{},
		"ErrEnvBindNotExistWithoutDefault": true,
		"FieldFlagSetCreateIgnore":         []string{},
		"Debug":                            false,
		"LogDebug":                         LogFunc(func(s string) { log.Println("[  DEBUG] " + s) }),
		"LogWarning":                       LogFunc(func(s string) { log.Println("[WARNING] " + s) }),
		"AppLabelList":                     []string([]string{}),
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
		if cc.FlagSet != nil {
			cc.FlagSet.Usage = func() {
				xflag.PrintDefaults(cc.FlagSet)
			}
		}
	})
}
