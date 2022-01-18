// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package xconf

import (
	"flag"
	"io"
	"log"
	"os"
	"sync/atomic"
	"unsafe"

	"github.com/sandwich-go/xconf/xutil"
)

// Options should use NewOptions to initialize it
type Options struct {
	OptionUsage string `xconf:"option_usage"`
	// annotation@NewFunc(comment="Parse时会由指定的File中加载配置")
	Files []string `xconf:"files"`
	// annotation@Readers(comment="Parse时会由指定的Reader中加载配置")
	Readers []io.Reader `xconf:"readers" usage:"Parse时会由指定的Reader中加载配置"`
	// annotation@FlagSet(comment="Parse使用的FlagSet，xconf会自动在flag中创建字段定义,如指定为空则不会创建")
	FlagSet *flag.FlagSet `xconf:"flag_set" usage:"Parse使用的FlagSet，xconf会自动在flag中创建字段定义,如指定为空则不会创建"`
	// annotation@FlagArgs(comment="FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑")
	FlagArgs []string `xconf:"flag_args" usage:"FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑"`
	// annotation@Environ(comment="Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.Usage()获取当前支持的FlagSet与Env参数定义")
	Environ []string `xconf:"environ" usage:"Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.Usage()获取当前支持的FlagSet与Env参数定义"`
	// annotation@ErrorHandling(comment="错误处理模式")
	ErrorHandling ErrorHandling `xconf:"error_handling" usage:"错误处理模式"`
	// annotation@TagName(comment="xconf使用的字段TAG名称,默认:xconf")
	TagName string `xconf:"tag_name" usage:"xconf使用的字段TAG名称,默认:xconf"`
	// annotation@DecoderConfigOption(comment="xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure")
	DecoderConfigOption []DecoderConfigOption `xconf:"decoder_config_option" usage:"xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure"`
	// annotation@MapMerge(comment="map是否开启merge模式，详情见文档")
	MapMerge bool `xconf:"map_merge" usage:"map是否开启merge模式，详情见文档"`
	// annotation@FieldTagConvertor(comment="字段名转换到FiledPath时优先使用TagName指定的名称，否则使用该函数转换")
	FieldTagConvertor FieldTagConvertor `xconf:"field_tag_convertor" usage:"字段名转换到FiledPath时优先使用TagName指定的名称，否则使用该函数转换"`
	// annotation@FieldPathRemoved(comment="弃用的配置，目标结构中已经删除，但配置文件中可能存在，解析时不会认为是错误，会将该配置丢弃，并打印WARNING日志")
	FieldPathRemoved []string `xconf:"field_path_removed" usage:"弃用的配置，目标结构中已经删除，但配置文件中可能存在，解析时不会认为是错误，会将该配置丢弃，并打印WARNING日志"`
	// annotation@Debug(comment="debug模式下输出调试信息")
	Debug bool `xconf:"debug" usage:"debug模式下输出调试信息"`
	// annotation@LogDebug(comment="DEBUG日志")
	LogDebug LogFunc `xconf:"log_debug" usage:"DEBUG日志"`
	// annotation@LogWarning(comment="WARNING日志")
	LogWarning LogFunc `xconf:"log_warning" usage:"WARNING日志"`
	// annotation@AppLabelList(comment="应用层Label，用于灰度发布场景")
	AppLabelList []string `xconf:"app_label_list" usage:"应用层Label，用于灰度发布场景"`
	// annotation@EnvBindShouldErrorWhenFailed(comment="EnvBind时如果Env中不存在指定的key而且没有指定默认值时是否返回错误")
	EnvBindShouldErrorWhenFailed bool `xconf:"env_bind_should_error_when_failed" usage:"EnvBind时如果Env中不存在指定的key而且没有指定默认值时是否返回错误"`
	// annotation@FlagCreateIgnoreFiledPath(comment="不创建到FlagSet中的字段FieldPath")
	// todo: 可以通过tag中指定flagoff规避这个字段的支持
	FlagCreateIgnoreFiledPath []string `xconf:"flag_create_ignore_filed_path" usage:"不创建到FlagSet中的字段FieldPath"`
	// annotation@ParseDefault(comment="是否解析struct标签中的default数据，解析规则参考xflag支持")
	ParseDefault bool `xconf:"parse_default" usage:"是否解析struct标签中的default数据，解析规则参考xflag支持"`
	// annotation@TagNameForDefaultValue(comment="默认值TAG名称,默认default")
	TagNameForDefaultValue string `xconf:"tag_name_for_default_value" usage:"默认值TAG名称,默认default"`
	// annotation@ReplaceFlagSetUsage(comment="是否替换FlagSet的Usage，使用xconf内置版本")
	ReplaceFlagSetUsage bool `xconf:"replace_flag_set_usage" usage:"是否替换FlagSet的Usage，使用xconf内置版本"`
}

// NewOptions new Options
func NewOptions(opts ...Option) *Options {
	cc := newDefaultOptions()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogOptions != nil {
		watchDogOptions(cc)
	}
	return cc
}

// ApplyOption apply mutiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *Options) ApplyOption(opts ...Option) []Option {
	var previous []Option
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// Option option func
type Option func(cc *Options) Option

// WithOptionUsage option func for filed OptionUsage
func WithOptionUsage(v string) Option {
	return func(cc *Options) Option {
		previous := cc.OptionUsage
		cc.OptionUsage = v
		return WithOptionUsage(previous)
	}
}

// WithFiles option func for filed Files
func WithFiles(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.Files
		cc.Files = v
		return WithFiles(previous...)
	}
}

// WithFiles option func for filed Files append
func WithFilesAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.Files
		cc.Files = append(cc.Files, v...)
		return WithFiles(previous...)
	}
}

// WithReaders Parse时会由指定的Reader中加载配置
func WithReaders(v ...io.Reader) Option {
	return func(cc *Options) Option {
		previous := cc.Readers
		cc.Readers = v
		return WithReaders(previous...)
	}
}

// WithReaders Parse时会由指定的Reader中加载配置 append
func WithReadersAppend(v ...io.Reader) Option {
	return func(cc *Options) Option {
		previous := cc.Readers
		cc.Readers = append(cc.Readers, v...)
		return WithReaders(previous...)
	}
}

// WithFlagSet Parse使用的FlagSet，xconf会自动在flag中创建字段定义,如指定为空则不会创建
func WithFlagSet(v *flag.FlagSet) Option {
	return func(cc *Options) Option {
		previous := cc.FlagSet
		cc.FlagSet = v
		return WithFlagSet(previous)
	}
}

// WithFlagArgs FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑
func WithFlagArgs(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FlagArgs
		cc.FlagArgs = v
		return WithFlagArgs(previous...)
	}
}

// WithFlagArgs FlagSet解析使用的Args列表，默认为os.Args[1:]，如指定为空则不会触发FlagSet的定义和解析逻辑 append
func WithFlagArgsAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FlagArgs
		cc.FlagArgs = append(cc.FlagArgs, v...)
		return WithFlagArgs(previous...)
	}
}

// WithEnviron Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.Usage()获取当前支持的FlagSet与Env参数定义
func WithEnviron(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.Environ
		cc.Environ = v
		return WithEnviron(previous...)
	}
}

// WithEnviron Parse解析的环境变量，内部将其转换为FlagSet处理，支持的类型参考FlagSet，可以通过xconf.Usage()获取当前支持的FlagSet与Env参数定义 append
func WithEnvironAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.Environ
		cc.Environ = append(cc.Environ, v...)
		return WithEnviron(previous...)
	}
}

// WithErrorHandling 错误处理模式
func WithErrorHandling(v ErrorHandling) Option {
	return func(cc *Options) Option {
		previous := cc.ErrorHandling
		cc.ErrorHandling = v
		return WithErrorHandling(previous)
	}
}

// WithTagName xconf使用的字段TAG名称,默认:xconf
func WithTagName(v string) Option {
	return func(cc *Options) Option {
		previous := cc.TagName
		cc.TagName = v
		return WithTagName(previous)
	}
}

// WithDecoderConfigOption xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure
func WithDecoderConfigOption(v ...DecoderConfigOption) Option {
	return func(cc *Options) Option {
		previous := cc.DecoderConfigOption
		cc.DecoderConfigOption = v
		return WithDecoderConfigOption(previous...)
	}
}

// WithDecoderConfigOption xconf内部依赖mapstructure，改方法用户用户层自定义mapstructure解析参数,参考：https://github.com/mitchellh/mapstructure append
func WithDecoderConfigOptionAppend(v ...DecoderConfigOption) Option {
	return func(cc *Options) Option {
		previous := cc.DecoderConfigOption
		cc.DecoderConfigOption = append(cc.DecoderConfigOption, v...)
		return WithDecoderConfigOption(previous...)
	}
}

// WithMapMerge map是否开启merge模式，详情见文档
func WithMapMerge(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.MapMerge
		cc.MapMerge = v
		return WithMapMerge(previous)
	}
}

// WithFieldTagConvertor 字段名转换到FiledPath时优先使用TagName指定的名称，否则使用该函数转换
func WithFieldTagConvertor(v FieldTagConvertor) Option {
	return func(cc *Options) Option {
		previous := cc.FieldTagConvertor
		cc.FieldTagConvertor = v
		return WithFieldTagConvertor(previous)
	}
}

// WithFieldPathRemoved 弃用的配置，目标结构中已经删除，但配置文件中可能存在，解析时不会认为是错误，会将该配置丢弃，并打印WARNING日志
func WithFieldPathRemoved(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FieldPathRemoved
		cc.FieldPathRemoved = v
		return WithFieldPathRemoved(previous...)
	}
}

// WithFieldPathRemoved 弃用的配置，目标结构中已经删除，但配置文件中可能存在，解析时不会认为是错误，会将该配置丢弃，并打印WARNING日志 append
func WithFieldPathRemovedAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FieldPathRemoved
		cc.FieldPathRemoved = append(cc.FieldPathRemoved, v...)
		return WithFieldPathRemoved(previous...)
	}
}

// WithDebug debug模式下输出调试信息
func WithDebug(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.Debug
		cc.Debug = v
		return WithDebug(previous)
	}
}

// WithLogDebug DEBUG日志
func WithLogDebug(v LogFunc) Option {
	return func(cc *Options) Option {
		previous := cc.LogDebug
		cc.LogDebug = v
		return WithLogDebug(previous)
	}
}

// WithLogWarning WARNING日志
func WithLogWarning(v LogFunc) Option {
	return func(cc *Options) Option {
		previous := cc.LogWarning
		cc.LogWarning = v
		return WithLogWarning(previous)
	}
}

// WithAppLabelList 应用层Label，用于灰度发布场景
func WithAppLabelList(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.AppLabelList
		cc.AppLabelList = v
		return WithAppLabelList(previous...)
	}
}

// WithAppLabelList 应用层Label，用于灰度发布场景 append
func WithAppLabelListAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.AppLabelList
		cc.AppLabelList = append(cc.AppLabelList, v...)
		return WithAppLabelList(previous...)
	}
}

// WithEnvBindShouldErrorWhenFailed EnvBind时如果Env中不存在指定的key而且没有指定默认值时是否返回错误
func WithEnvBindShouldErrorWhenFailed(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.EnvBindShouldErrorWhenFailed
		cc.EnvBindShouldErrorWhenFailed = v
		return WithEnvBindShouldErrorWhenFailed(previous)
	}
}

// WithFlagCreateIgnoreFiledPath 不创建到FlagSet中的字段FieldPath
func WithFlagCreateIgnoreFiledPath(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FlagCreateIgnoreFiledPath
		cc.FlagCreateIgnoreFiledPath = v
		return WithFlagCreateIgnoreFiledPath(previous...)
	}
}

// WithFlagCreateIgnoreFiledPath 不创建到FlagSet中的字段FieldPath append
func WithFlagCreateIgnoreFiledPathAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FlagCreateIgnoreFiledPath
		cc.FlagCreateIgnoreFiledPath = append(cc.FlagCreateIgnoreFiledPath, v...)
		return WithFlagCreateIgnoreFiledPath(previous...)
	}
}

// WithParseDefault 是否解析struct标签中的default数据，解析规则参考xflag支持
func WithParseDefault(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.ParseDefault
		cc.ParseDefault = v
		return WithParseDefault(previous)
	}
}

// WithTagNameForDefaultValue 默认值TAG名称,默认default
func WithTagNameForDefaultValue(v string) Option {
	return func(cc *Options) Option {
		previous := cc.TagNameForDefaultValue
		cc.TagNameForDefaultValue = v
		return WithTagNameForDefaultValue(previous)
	}
}

// WithReplaceFlagSetUsage 是否替换FlagSet的Usage，使用xconf内置版本
func WithReplaceFlagSetUsage(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.ReplaceFlagSetUsage
		cc.ReplaceFlagSetUsage = v
		return WithReplaceFlagSetUsage(previous)
	}
}

// InstallOptionsWatchDog the installed func will called when NewOptions  called
func InstallOptionsWatchDog(dog func(cc *Options)) { watchDogOptions = dog }

// watchDogOptions global watch dog
var watchDogOptions func(cc *Options)

// newDefaultOptions new default Options
func newDefaultOptions() *Options {
	cc := &Options{}

	for _, opt := range [...]Option{
		WithOptionUsage(optionUsage),
		WithFiles([]string{}...),
		WithReaders([]io.Reader{}...),
		WithFlagSet(flag.CommandLine),
		WithFlagArgs(os.Args[1:]...),
		WithEnviron(os.Environ()...),
		WithErrorHandling(PanicOnError),
		WithTagName(DefaultTagName),
		WithDecoderConfigOption(nil...),
		WithMapMerge(false),
		WithFieldTagConvertor(xutil.SnakeCase),
		WithFieldPathRemoved(make([]string, 0)...),
		WithDebug(false),
		WithLogDebug(func(s string) { log.Println("[  DEBUG] " + s) }),
		WithLogWarning(func(s string) { log.Println("[WARNING] " + s) }),
		WithAppLabelList([]string{}...),
		WithEnvBindShouldErrorWhenFailed(true),
		WithFlagCreateIgnoreFiledPath(make([]string, 0)...),
		WithParseDefault(true),
		WithTagNameForDefaultValue(DefaultValueTagName),
		WithReplaceFlagSetUsage(true),
	} {
		opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *Options) AtomicSetFunc() func(interface{}) { return AtomicOptionsSet }

// atomicOptions global *Options holder
var atomicOptions unsafe.Pointer

// onAtomicOptionsSet global call back when  AtomicOptionsSet called by XConf.
// use OptionsInterface.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomicOptions
var onAtomicOptionsSet func(cc OptionsInterface) bool

// InstallCallbackOnAtomicOptionsSet install callback
func InstallCallbackOnAtomicOptionsSet(callback func(cc OptionsInterface) bool) {
	onAtomicOptionsSet = callback
}

// AtomicOptionsSet atomic setter for *Options
func AtomicOptionsSet(update interface{}) {
	cc := update.(*Options)
	if onAtomicOptionsSet != nil && !onAtomicOptionsSet(cc) {
		return
	}
	atomic.StorePointer(&atomicOptions, (unsafe.Pointer)(cc))
}

// AtomicOptions return atomic *OptionsVisitor
func AtomicOptions() OptionsVisitor {
	current := (*Options)(atomic.LoadPointer(&atomicOptions))
	if current == nil {
		defaultOne := newDefaultOptions()
		if watchDogOptions != nil {
			watchDogOptions(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomicOptions, nil, (unsafe.Pointer)(defaultOne))
		return (*Options)(atomic.LoadPointer(&atomicOptions))
	}
	return current
}

// all getter func
func (cc *Options) GetOptionUsage() string                        { return cc.OptionUsage }
func (cc *Options) GetFiles() []string                            { return cc.Files }
func (cc *Options) GetReaders() []io.Reader                       { return cc.Readers }
func (cc *Options) GetFlagSet() *flag.FlagSet                     { return cc.FlagSet }
func (cc *Options) GetFlagArgs() []string                         { return cc.FlagArgs }
func (cc *Options) GetEnviron() []string                          { return cc.Environ }
func (cc *Options) GetErrorHandling() ErrorHandling               { return cc.ErrorHandling }
func (cc *Options) GetTagName() string                            { return cc.TagName }
func (cc *Options) GetDecoderConfigOption() []DecoderConfigOption { return cc.DecoderConfigOption }
func (cc *Options) GetMapMerge() bool                             { return cc.MapMerge }
func (cc *Options) GetFieldTagConvertor() FieldTagConvertor       { return cc.FieldTagConvertor }
func (cc *Options) GetFieldPathRemoved() []string                 { return cc.FieldPathRemoved }
func (cc *Options) GetDebug() bool                                { return cc.Debug }
func (cc *Options) GetLogDebug() LogFunc                          { return cc.LogDebug }
func (cc *Options) GetLogWarning() LogFunc                        { return cc.LogWarning }
func (cc *Options) GetAppLabelList() []string                     { return cc.AppLabelList }
func (cc *Options) GetEnvBindShouldErrorWhenFailed() bool         { return cc.EnvBindShouldErrorWhenFailed }
func (cc *Options) GetFlagCreateIgnoreFiledPath() []string        { return cc.FlagCreateIgnoreFiledPath }
func (cc *Options) GetParseDefault() bool                         { return cc.ParseDefault }
func (cc *Options) GetTagNameForDefaultValue() string             { return cc.TagNameForDefaultValue }
func (cc *Options) GetReplaceFlagSetUsage() bool                  { return cc.ReplaceFlagSetUsage }

// OptionsVisitor visitor interface for Options
type OptionsVisitor interface {
	GetOptionUsage() string
	GetFiles() []string
	GetReaders() []io.Reader
	GetFlagSet() *flag.FlagSet
	GetFlagArgs() []string
	GetEnviron() []string
	GetErrorHandling() ErrorHandling
	GetTagName() string
	GetDecoderConfigOption() []DecoderConfigOption
	GetMapMerge() bool
	GetFieldTagConvertor() FieldTagConvertor
	GetFieldPathRemoved() []string
	GetDebug() bool
	GetLogDebug() LogFunc
	GetLogWarning() LogFunc
	GetAppLabelList() []string
	GetEnvBindShouldErrorWhenFailed() bool
	GetFlagCreateIgnoreFiledPath() []string
	GetParseDefault() bool
	GetTagNameForDefaultValue() string
	GetReplaceFlagSetUsage() bool
}

// OptionsInterface visitor + ApplyOption interface for Options
type OptionsInterface interface {
	OptionsVisitor
	ApplyOption(...Option) []Option
}
