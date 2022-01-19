package xcmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

type Executer = func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error

var DefaultXConfOption = []xconf.Option{
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
		// annotation@Synopsis(comment="少于一行的操作说明")
		"Synopsis": "",
		// annotation@Usage(comment="详细说明")
		"Usage": "",
		// annotation@Execute(comment="执行方法")
		"Execute": Executer(nil),
		// annotation@XConfOption(comment="Parser依赖的XConf配置")
		"XConfOption": ([]xconf.Option)(DefaultXConfOption),
		// annotation@Parser(comment="配置解析")
		"Parser": MiddlewareFunc(Parser),
		// annotation@Executer(comment="配置解析")
		"OnExecuterLost": Executer(func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
			return c.wrapErr(errors.New("no executer"))
		}),
	}
}

var _ = configOptionDeclareWithDefault

func Parser(ctx context.Context, c *Command, ff *flag.FlagSet, args []string, next Executer) error {
	// 默认 usage 无参
	ff.Usage = func() {
		c.Explain(c.Output)
		xflag.PrintDefaults(ff)
	}
	if c.cc.GetBind() == nil {
		return ff.Parse(args)
	}
	cc := xconf.NewOptions(
		xconf.WithErrorHandling(xconf.ContinueOnError),
		xconf.WithFlagSet(ff),
		xconf.WithFlagArgs(args...))
	cc.ApplyOption(c.cc.GetXConfOption()...)

	// 获取bindto结构合法的FieldPath，并过滤合法的BindToFieldPath
	_, fieldsMap := xconf.NewStruct(
		reflect.New(reflect.ValueOf(c.cc.GetBind()).Type().Elem()).Interface(),
		cc.TagName,
		cc.TagNameForDefaultValue,
		cc.FieldTagConvertor).Map()

	var ignorePath []string
	var invalidKeys []string

	if len(c.cc.GetBindFieldPath()) > 0 {
		for k := range fieldsMap {
			if !xutil.ContainStringEqualFold(c.cc.GetBindFieldPath(), k) {
				ignorePath = append(ignorePath, k)
			}
		}
	}
	for _, v := range c.cc.GetBindFieldPath() {
		if _, ok := fieldsMap[v]; !ok {
			invalidKeys = append(invalidKeys, v)
		}
	}

	if len(invalidKeys) > 0 {
		var keys []string
		for k := range fieldsMap {
			keys = append(keys, k)
		}
		return c.wrapErr(fmt.Errorf("option BindFieldPath has invalid item: %s valid: %v", strings.Join(invalidKeys, ","), keys))
	}

	cc.ApplyOption(xconf.WithFlagCreateIgnoreFiledPath(ignorePath...))
	x := xconf.NewWithConf(cc)

	// Available Commands + Flags
	cc.FlagSet.Usage = func() {
		c.Explain(c.Output)
		x.UsageToWriter(c.Output, args...)
	}
	err := x.Parse(c.cc.GetBind())
	if err != nil {
		if IsErrHelp(err) {
			err = nil
		} else {
			err = fmt.Errorf("[ApplyArgs] %s", err.Error())
		}
	}
	if err != nil {
		return err
	}
	return next(ctx, c, ff, args)
}
