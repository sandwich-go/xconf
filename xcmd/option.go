package xcmd

import (
	"github.com/sandwich-go/xconf"
)

// SetDefaultXConfOption 设定默认的XConf参数项
func SetDefaultXConfOption(opts ...xconf.Option) { defaultXConfOption = opts }

var defaultXConfOption = []xconf.Option{
	xconf.WithErrorHandling(xconf.ContinueOnError),
	xconf.WithReplaceFlagSetUsage(false),
}

//go:generate optiongen
func configOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Short(comment="少于一行的操作说明,简短")
		"Short": "",
		// annotation@Description(comment="详细说明，可以多行，自行做格式控制")
		"Description": "",
		// annotation@Examples(comment="例子,可以多行，自行做格式控制")
		"Examples": "",
		// annotation@XConfOption(comment="Parser依赖的XConf配置")
		"XConfOption": ([]xconf.Option)(defaultXConfOption),
		// annotation@Parser(comment="配置解析")
		"Parser": MiddlewareFunc(ParserXConf),
		// annotation@Executer(comment="当未配置Parser时触发该默认逻辑")
		"OnExecuterLost":         Executer(DefaultExecuter),
		"SuggestionsMinDistance": 2,
	}
}

var _ = configOptionDeclareWithDefault
