package xcmd

import (
	"context"
	"flag"

	"github.com/sandwich-go/xconf"
)

// Executer 命令执行方法
type Executer = func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error

// SetDefaultXConfOption 设定默认的XConf参数项
func SetDefaultXConfOption(opts ...xconf.Option) { defaultXConfOption = opts }

var defaultXConfOption = []xconf.Option{
	xconf.WithErrorHandling(xconf.ContinueOnError),
	xconf.WithReplaceFlagSetUsage(false),
}

//go:generate optiongen --new_func_return=interface
func configOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Bind(comment="命令绑定的参数结构")
		"Bind": (interface{})(nil),
		// annotation@BindFieldPath(comment="命令绑定的参数FieldPath,如空则全部绑定")
		"BindFieldPath": []string{},
		// annotation@Short(comment="少于一行的操作说明")
		"Short": "",
		// annotation@Usage(comment="详细说明")
		"Usage": "",
		// annotation@Execute(comment="执行方法")
		"Execute": Executer(nil),
		// annotation@XConfOption(comment="Parser依赖的XConf配置")
		"XConfOption": ([]xconf.Option)(defaultXConfOption),
		// annotation@Parser(comment="配置解析")
		"Parser": MiddlewareFunc(ParserXConf),
		// annotation@Executer(comment="当未配置Parser时触发该默认逻辑")
		"OnExecuterLost": Executer(func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
			c.Usage()
			return ErrHelp
		}),
	}
}

var _ = configOptionDeclareWithDefault
