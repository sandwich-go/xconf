package xcmd

import (
	"io"
	"os"

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
		// annotation@SuggestionsMinDistance(comment="推荐命令最低关联长度")
		"SuggestionsMinDistance": 2,
		// annotation@Output(comment="输出")
		"Output": io.Writer(os.Stdout),
		// annotation@Deprecated(comment="不推荐使用的命令说明,只有配置了该说明的命令才会显示Deprecated标签")
		"Deprecated": "",
		// annotation@Author(comment="命令作者联系信息，只用于显示")
		"Author": []string{},
	}
}

var _ = configOptionDeclareWithDefault
