package xcmd

import (
	"context"
	"flag"

	"github.com/sandwich-go/xconf"
)

type Executer = func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error

//go:generate optiongen --new_func_return=interface
func configOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Bind(comment="命令绑定的参数结构")
		"Bind": (interface{})(nil),
		// annotation@BindFieldPath(comment="命令绑定的参数FieldPath,如空则全部绑定")
		"BindFieldPath": []string{},
		// annotation@Synopsis(comment="少于一行的操作说明")
		"Synopsis": "",
		// annotation@Usage(comment="详细说明")
		"Usage": "",
		// annotation@Execute(comment="执行方法")
		"Execute": Executer(nil),
		"XConfOption": ([]xconf.Option)(
			[]xconf.Option{
				xconf.WithErrorHandling(xconf.ContinueOnError),
				xconf.WithReplaceFlagSetUsage(false),
			},
		),
	}
}

var _ = configOptionDeclareWithDefault
