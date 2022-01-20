package xcmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xutil"
)

// ParserXConf xconf的Parser方法
func ParserXConf(ctx context.Context, c *Command, ff *flag.FlagSet, args []string, next Executer) error {
	cc := xconf.NewOptions(
		xconf.WithErrorHandling(xconf.ContinueOnError),
		xconf.WithFlagSet(ff),
		xconf.WithFlagArgs(args...))
	cc.ApplyOption(c.cc.GetXConfOption()...)

	x := xconf.NewWithConf(cc)
	keysList := xconf.FieldPathList(c.cc.GetBind(), x)

	// keysList中的元素如果没有包含在GetBindFieldPath中，则为不允许通过flag覆盖的item
	var ignorePath []string
	if len(c.cc.GetBindFieldPath()) > 0 {
		for _, k := range keysList {
			if !xutil.ContainStringEqualFold(c.cc.GetBindFieldPath(), k) {
				ignorePath = append(ignorePath, k)
			}
		}
	}
	// 检查GetBindFieldPath中的key是否合法
	var invalidKeys []string
	for _, v := range c.cc.GetBindFieldPath() {
		if !xutil.ContainString(keysList, v) {
			invalidKeys = append(invalidKeys, v)
		}
	}
	if len(invalidKeys) > 0 {
		return c.wrapErr(fmt.Errorf("option BindFieldPath has invalid item: %s valid: %v", strings.Join(invalidKeys, ","), keysList))
	}
	// 更新忽略调的绑定字段，重新狗仔xconf实例
	cc.ApplyOption(xconf.WithFlagCreateIgnoreFiledPath(ignorePath...))
	x = xconf.NewWithConf(cc)

	// 更新FlagSet的Usage，使用xconf内置版本
	cc.FlagSet.Usage = func() {
		c.Explain(c.Output)
		fmt.Fprintf(c.Output, "Flags:\n")
		x.UsageToWriter(c.Output, args...)
	}
	err := x.Parse(c.cc.GetBind())
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
	return next(ctx, c, ff, args)
}
