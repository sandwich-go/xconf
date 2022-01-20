// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package xflag

import (
	"flag"
	"log"
	"strings"

	"github.com/sandwich-go/xconf/xflag/vars"
)

// Options should use NewOptions to initialize it
type Options struct {
	Name                      string
	TagName                   string // 使用的tag key,如不设定则使用
	UsageTagName              string
	Flatten                   bool // 是否使用扁平模式，不使用.分割
	FlagSet                   *flag.FlagSet
	FlagValueProvider         vars.FlagValueProvider
	KeyFormat                 KeyFormat
	FlagCreateIgnoreFiledPath []string
	LogDebug                  LogFunc
	LogWarning                LogFunc
	StringAlias               func(s string) string
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

// WithName option func for filed Name
func WithName(v string) Option {
	return func(cc *Options) Option {
		previous := cc.Name
		cc.Name = v
		return WithName(previous)
	}
}

// WithTagName option func for filed TagName
func WithTagName(v string) Option {
	return func(cc *Options) Option {
		previous := cc.TagName
		cc.TagName = v
		return WithTagName(previous)
	}
}

// WithUsageTagName option func for filed UsageTagName
func WithUsageTagName(v string) Option {
	return func(cc *Options) Option {
		previous := cc.UsageTagName
		cc.UsageTagName = v
		return WithUsageTagName(previous)
	}
}

// WithFlatten option func for filed Flatten
func WithFlatten(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.Flatten
		cc.Flatten = v
		return WithFlatten(previous)
	}
}

// WithFlagSet option func for filed FlagSet
func WithFlagSet(v *flag.FlagSet) Option {
	return func(cc *Options) Option {
		previous := cc.FlagSet
		cc.FlagSet = v
		return WithFlagSet(previous)
	}
}

// WithFlagValueProvider option func for filed FlagValueProvider
func WithFlagValueProvider(v vars.FlagValueProvider) Option {
	return func(cc *Options) Option {
		previous := cc.FlagValueProvider
		cc.FlagValueProvider = v
		return WithFlagValueProvider(previous)
	}
}

// WithKeyFormat option func for filed KeyFormat
func WithKeyFormat(v KeyFormat) Option {
	return func(cc *Options) Option {
		previous := cc.KeyFormat
		cc.KeyFormat = v
		return WithKeyFormat(previous)
	}
}

// WithFlagCreateIgnoreFiledPath option func for filed FlagCreateIgnoreFiledPath
func WithFlagCreateIgnoreFiledPath(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FlagCreateIgnoreFiledPath
		cc.FlagCreateIgnoreFiledPath = v
		return WithFlagCreateIgnoreFiledPath(previous...)
	}
}

// WithFlagCreateIgnoreFiledPath option func for filed FlagCreateIgnoreFiledPath append
func WithFlagCreateIgnoreFiledPathAppend(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.FlagCreateIgnoreFiledPath
		cc.FlagCreateIgnoreFiledPath = append(cc.FlagCreateIgnoreFiledPath, v...)
		return WithFlagCreateIgnoreFiledPath(previous...)
	}
}

// WithLogDebug option func for filed LogDebug
func WithLogDebug(v LogFunc) Option {
	return func(cc *Options) Option {
		previous := cc.LogDebug
		cc.LogDebug = v
		return WithLogDebug(previous)
	}
}

// WithLogWarning option func for filed LogWarning
func WithLogWarning(v LogFunc) Option {
	return func(cc *Options) Option {
		previous := cc.LogWarning
		cc.LogWarning = v
		return WithLogWarning(previous)
	}
}

// WithStringAlias option func for filed StringAlias
func WithStringAlias(v func(s string) string) Option {
	return func(cc *Options) Option {
		previous := cc.StringAlias
		cc.StringAlias = v
		return WithStringAlias(previous)
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
		WithName(""),
		WithTagName("xconf"),
		WithUsageTagName("usage"),
		WithFlatten(false),
		WithFlagSet(flag.NewFlagSet("flagmaker", flag.ContinueOnError)),
		WithFlagValueProvider(vars.DefaultFlagValueProvider),
		WithKeyFormat(func(s string) string { return strings.ToLower(s) }),
		WithFlagCreateIgnoreFiledPath(make([]string, 0)...),
		WithLogDebug(func(s string) { log.Print("debug:" + s) }),
		WithLogWarning(func(s string) { log.Print("warning: " + s) }),
		WithStringAlias(func(s string) string {
			return s
		}),
	} {
		opt(cc)
	}

	return cc
}

// all getter func
func (cc *Options) GetName() string                              { return cc.Name }
func (cc *Options) GetTagName() string                           { return cc.TagName }
func (cc *Options) GetUsageTagName() string                      { return cc.UsageTagName }
func (cc *Options) GetFlatten() bool                             { return cc.Flatten }
func (cc *Options) GetFlagSet() *flag.FlagSet                    { return cc.FlagSet }
func (cc *Options) GetFlagValueProvider() vars.FlagValueProvider { return cc.FlagValueProvider }
func (cc *Options) GetKeyFormat() KeyFormat                      { return cc.KeyFormat }
func (cc *Options) GetFlagCreateIgnoreFiledPath() []string       { return cc.FlagCreateIgnoreFiledPath }
func (cc *Options) GetLogDebug() LogFunc                         { return cc.LogDebug }
func (cc *Options) GetLogWarning() LogFunc                       { return cc.LogWarning }
func (cc *Options) GetStringAlias() func(s string) string        { return cc.StringAlias }

// OptionsVisitor visitor interface for Options
type OptionsVisitor interface {
	GetName() string
	GetTagName() string
	GetUsageTagName() string
	GetFlatten() bool
	GetFlagSet() *flag.FlagSet
	GetFlagValueProvider() vars.FlagValueProvider
	GetKeyFormat() KeyFormat
	GetFlagCreateIgnoreFiledPath() []string
	GetLogDebug() LogFunc
	GetLogWarning() LogFunc
	GetStringAlias() func(s string) string
}

// OptionsInterface visitor + ApplyOption interface for Options
type OptionsInterface interface {
	OptionsVisitor
	ApplyOption(...Option) []Option
}
