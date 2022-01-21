package xcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xutil"
)

// ParserXConf xconf的Parser方法
func ParserXConf(ctx context.Context, cmd *Command, next Executer) error {
	cc := xconf.NewOptions(
		xconf.WithErrorHandling(xconf.ContinueOnError),
		xconf.WithFlagSet(cmd.FlagSet),
		xconf.WithFlagArgs(cmd.RawArgs...))
	cc.ApplyOption(cmd.Config().GetXConfOption()...)

	x := xconf.NewWithConf(cc)
	keysList := xconf.FieldPathList(cmd.Bind(), x)

	bindFieldPath := cmd.bindFieldPath
	// keysList中的元素如果没有包含在GetBindFieldPath中，则为不允许通过flag覆盖的item
	var ignorePath []string
	// 检查GetBindFieldPath中的key是否合法
	var invalidKeys []string
	if len(bindFieldPath) > 0 {
		for _, k := range keysList {
			if !xutil.ContainStringEqualFold(bindFieldPath, k) {
				ignorePath = append(ignorePath, k)
			}
		}
		for _, v := range bindFieldPath {
			if !xutil.ContainString(keysList, v) {
				invalidKeys = append(invalidKeys, v)
			}
		}
	}

	if len(invalidKeys) > 0 {
		return cmd.wrapErr(fmt.Errorf("option BindFieldPath has invalid item: %s valid: %v", strings.Join(invalidKeys, ","), keysList))
	}
	// 更新忽略调的绑定字段，重新狗仔xconf实例
	cc.ApplyOption(xconf.WithFlagCreateIgnoreFiledPath(ignorePath...))
	x = xconf.NewWithConf(cc)

	// 更新FlagSet的Usage，使用xconf内置版本
	cc.FlagSet.Usage = func() {
		cmd.Explain(cmd.Output)
		fmt.Fprintf(cmd.Output, "Flags:\n")
		x.UsageToWriter(cmd.Output, cmd.RawArgs...)
	}
	err := x.Parse(cmd.bind)
	if err != nil {
		if IsErrHelp(err) {
			err = ErrHelp
		} else {
			err = fmt.Errorf("[ParserXConf] %s", err.Error())
		}
	}
	if err != nil {
		return err
	}
	return next(ctx, cmd)
}
