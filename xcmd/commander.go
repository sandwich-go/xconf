package xcmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

type MiddlewareFunc = func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string, next Executer) error

type Command struct {
	name          string
	cc            ConfigInterface
	Output        io.Writer
	commands      []*Command
	middleware    []MiddlewareFunc
	usageNamePath []string
}

func NewCommand(name string, opts ...ConfigOption) *Command {
	return NewCommandWithConfig(name, NewConfig(opts...))
}

func NewCommandWithConfig(name string, cc ConfigInterface) *Command {
	c := &Command{
		name:   name,
		cc:     cc,
		Output: os.Stdout,
	}
	c.usageNamePath = []string{name}
	return c
}

func (c *Command) Config() ConfigInterface { return c.cc }

func (c *Command) Use(middleware ...MiddlewareFunc) *Command {
	c.middleware = append(c.middleware, middleware...)
	return c
}

func (c *Command) AddCommand(sub *Command, middleware ...MiddlewareFunc) {
	sub.usageNamePath = append(c.usageNamePath, sub.usageNamePath...)
	sub.middleware = c.combineMiddlewareFunc(middleware...)
	c.commands = append(c.commands, sub)
}

func (c *Command) combineMiddlewareFunc(middleware ...MiddlewareFunc) []MiddlewareFunc {
	m := make([]MiddlewareFunc, 0, len(c.middleware)+len(middleware))
	m = append(m, c.middleware...)
	m = append(m, middleware...)
	return m
}

func (c *Command) Add(name string, opts ...ConfigOption) *Command {
	sub := NewCommand(name, opts...)
	c.AddCommand(sub)
	return sub
}

func (c *Command) wrapErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("command: %s err:%s", strings.Join(c.usageNamePath, " "), err.Error())
}

func (c *Command) ApplyArgs(ff *flag.FlagSet, args ...string) error {
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
	if len(c.cc.GetBindFieldPath()) > 0 {
		for k := range fieldsMap {
			if !xutil.ContainStringEqualFold(c.cc.GetBindFieldPath(), k) {
				ignorePath = append(ignorePath, k)
			}
		}
	}
	var invalidKeys []string
	for _, v := range c.cc.GetBindFieldPath() {
		if _, ok := fieldsMap[v]; !ok {
			invalidKeys = append(invalidKeys, v)
		}
	}

	if len(invalidKeys) > 0 {
		return c.wrapErr(fmt.Errorf("option BindFieldPath has invalid item:%s", strings.Join(invalidKeys, ",")))
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
		return fmt.Errorf("got err:%s while Parse", err.Error())
	}
	return nil
}

func (c *Command) Execute(ctx context.Context, args ...string) error {
	if len(args) != 0 {
		// 尝试在当前命令集下寻找子命令
		subCommandName := args[0]
		for _, cmd := range c.commands {
			if cmd.GetName() != subCommandName {
				continue
			}
			return cmd.Execute(ctx, args[1:]...)
		}
	}
	ff := flag.NewFlagSet(strings.Join(c.usageNamePath, "/"), flag.ContinueOnError)
	return ChainMiddleware(c.middleware...)(ctx, c, ff, args, exec)
}

func exec(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
	// arg为空或子命令没有找到，说明传入的参数不是命令名称，在当前层执行
	if err := c.ApplyArgs(ff, args...); err != nil {
		if IsErrHelp(err) {
			return nil
		}
		return fmt.Errorf("[ApplyArgs] %s", err.Error())
	}
	return c.Config().GetExecute()(context.Background(), c, ff, args)
}

func (c *Command) GetName() string { return c.name }
func (c *Command) Usage() string   { return c.cc.GetSynopsis() }

func (c *Command) CommandInheritBind(name string, opts ...ConfigOption) *Command {
	cc := NewConfig(WithBind(c.cc.GetBind()), WithBindFieldPath(c.cc.GetBindFieldPath()...))
	cc.ApplyOption(opts...)
	sub := NewCommandWithConfig(name, cc)
	c.AddCommand(sub)
	return sub
}

func (c *Command) Check() error {
	for _, v := range c.commands {
		err := v.ApplyArgs(flag.NewFlagSet(strings.Join(c.usageNamePath, "/"), flag.ContinueOnError))
		if err != nil {
			return err
		}
	}
	return nil
}

// ChainMiddleware middleware chain
func ChainMiddleware(middlewares ...MiddlewareFunc) MiddlewareFunc {
	n := len(middlewares)
	return func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string, next Executer) error {
		chain := func(currMiddleware MiddlewareFunc, currDispatcher Executer) Executer {
			return func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
				return currMiddleware(ctx, c, ff, args, currDispatcher)
			}
		}
		chainHandlerFunc := next
		for i := n - 1; i >= 0; i-- {
			chainHandlerFunc = chain(middlewares[i], chainHandlerFunc)
		}
		return chainHandlerFunc(ctx, c, ff, args)
	}
}
