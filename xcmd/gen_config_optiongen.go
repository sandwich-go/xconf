// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package xcmd

import (
	"io"
	"os"

	"github.com/sandwich-go/xconf"
)

// config should use NewConfig to initialize it
type config struct {
	// annotation@Short(comment="少于一行的操作说明,简短")
	Short string
	// annotation@Description(comment="详细说明，可以多行，自行做格式控制")
	Description string
	// annotation@Examples(comment="例子,可以多行，自行做格式控制")
	Examples string
	// annotation@XConfOption(comment="Parser依赖的XConf配置")
	XConfOption []xconf.Option
	// annotation@Parser(comment="配置解析")
	Parser MiddlewareFunc
	// annotation@SuggestionsMinDistance(comment="推荐命令最低关联长度")
	SuggestionsMinDistance int
	// annotation@Output(comment="输出")
	Output io.Writer
	// annotation@Deprecated(comment="不推荐使用的命令说明,只有配置了该说明的命令才会显示Deprecated标签")
	Deprecated string
	// annotation@Author(comment="命令作者联系信息，只用于显示")
	Author []string
	// annotation@Alias(comment="alias command")
	Alias []string
}

// NewConfig new config
func NewConfig(opts ...ConfigOption) *config {
	cc := newDefaultConfig()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogConfig != nil {
		watchDogConfig(cc)
	}
	return cc
}

// ApplyOption apply multiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *config) ApplyOption(opts ...ConfigOption) []ConfigOption {
	var previous []ConfigOption
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// ConfigOption option func
type ConfigOption func(cc *config) ConfigOption

// WithShort 少于一行的操作说明,简短
func WithShort(v string) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Short
		cc.Short = v
		return WithShort(previous)
	}
}

// WithDescription 详细说明，可以多行，自行做格式控制
func WithDescription(v string) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Description
		cc.Description = v
		return WithDescription(previous)
	}
}

// WithExamples 例子,可以多行，自行做格式控制
func WithExamples(v string) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Examples
		cc.Examples = v
		return WithExamples(previous)
	}
}

// WithXConfOption Parser依赖的XConf配置
func WithXConfOption(v ...xconf.Option) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.XConfOption
		cc.XConfOption = v
		return WithXConfOption(previous...)
	}
}

// WithParser 配置解析
func WithParser(v MiddlewareFunc) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Parser
		cc.Parser = v
		return WithParser(previous)
	}
}

// WithSuggestionsMinDistance 推荐命令最低关联长度
func WithSuggestionsMinDistance(v int) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.SuggestionsMinDistance
		cc.SuggestionsMinDistance = v
		return WithSuggestionsMinDistance(previous)
	}
}

// WithOutput 输出
func WithOutput(v io.Writer) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Output
		cc.Output = v
		return WithOutput(previous)
	}
}

// WithDeprecated 不推荐使用的命令说明,只有配置了该说明的命令才会显示Deprecated标签
func WithDeprecated(v string) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Deprecated
		cc.Deprecated = v
		return WithDeprecated(previous)
	}
}

// WithAuthor 命令作者联系信息，只用于显示
func WithAuthor(v ...string) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Author
		cc.Author = v
		return WithAuthor(previous...)
	}
}

// WithAlias alias command
func WithAlias(v ...string) ConfigOption {
	return func(cc *config) ConfigOption {
		previous := cc.Alias
		cc.Alias = v
		return WithAlias(previous...)
	}
}

// InstallConfigWatchDog the installed func will called when NewConfig  called
func InstallConfigWatchDog(dog func(cc *config)) { watchDogConfig = dog }

// watchDogConfig global watch dog
var watchDogConfig func(cc *config)

// newDefaultConfig new default config
func newDefaultConfig() *config {
	cc := &config{}

	for _, opt := range [...]ConfigOption{
		WithShort(""),
		WithDescription(""),
		WithExamples(""),
		WithXConfOption(defaultXConfOption...),
		WithParser(ParserXConf),
		WithSuggestionsMinDistance(2),
		WithOutput(os.Stdout),
		WithDeprecated(""),
		WithAuthor(make([]string, 0)...),
		WithAlias(make([]string, 0)...),
	} {
		opt(cc)
	}

	return cc
}

// all getter func
func (cc *config) GetShort() string               { return cc.Short }
func (cc *config) GetDescription() string         { return cc.Description }
func (cc *config) GetExamples() string            { return cc.Examples }
func (cc *config) GetXConfOption() []xconf.Option { return cc.XConfOption }
func (cc *config) GetParser() MiddlewareFunc      { return cc.Parser }
func (cc *config) GetSuggestionsMinDistance() int { return cc.SuggestionsMinDistance }
func (cc *config) GetOutput() io.Writer           { return cc.Output }
func (cc *config) GetDeprecated() string          { return cc.Deprecated }
func (cc *config) GetAuthor() []string            { return cc.Author }
func (cc *config) GetAlias() []string             { return cc.Alias }

// ConfigVisitor visitor interface for config
type ConfigVisitor interface {
	GetShort() string
	GetDescription() string
	GetExamples() string
	GetXConfOption() []xconf.Option
	GetParser() MiddlewareFunc
	GetSuggestionsMinDistance() int
	GetOutput() io.Writer
	GetDeprecated() string
	GetAuthor() []string
	GetAlias() []string
}

// ConfigInterface visitor + ApplyOption interface for config
type ConfigInterface interface {
	ConfigVisitor
	ApplyOption(...ConfigOption) []ConfigOption
}
