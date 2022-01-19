package xcmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type MiddlewareFunc = func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string, next Executer) error

type Command struct {
	name          string
	cc            ConfigInterface
	Output        io.Writer
	commands      []*Command
	middleware    []MiddlewareFunc
	preMiddleware []MiddlewareFunc
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

func (c *Command) UsePre(preMiddleware ...MiddlewareFunc) *Command {
	c.preMiddleware = append(c.preMiddleware, preMiddleware...)
	return c
}

func (c *Command) AddCommand(sub *Command, middleware ...MiddlewareFunc) {
	sub.usageNamePath = append(c.usageNamePath, sub.usageNamePath...)
	sub.middleware = combineMiddlewareFunc(c.middleware, middleware...)
	sub.preMiddleware = combineMiddlewareFunc(c.preMiddleware, sub.preMiddleware...)
	c.commands = append(c.commands, sub)
}

func combineMiddlewareFunc(middlewareNow []MiddlewareFunc, middleware ...MiddlewareFunc) []MiddlewareFunc {
	m := make([]MiddlewareFunc, 0, len(middlewareNow)+len(middleware))
	m = append(m, middlewareNow...)
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
	var allMiddlewares []MiddlewareFunc
	allMiddlewares = append(allMiddlewares, c.preMiddleware...)
	allMiddlewares = append(allMiddlewares, c.cc.GetParser())
	allMiddlewares = append(allMiddlewares, c.middleware...)
	return ChainMiddleware(allMiddlewares...)(ctx, c, ff, args, exec)
}
func exec(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
	executer := c.Config().GetExecute()
	if executer == nil {
		executer = c.cc.GetOnExecuterLost()
	}
	return executer(context.Background(), c, ff, args)
}

func (c *Command) GetName() string { return c.name }
func (c *Command) Usage() string   { return c.cc.GetSynopsis() }

func (c *Command) SubCommand(name string, opts ...ConfigOption) *Command {
	cc := NewConfig(
		WithBind(c.cc.GetBind()),
		WithBindFieldPath(c.cc.GetBindFieldPath()...),
		WithXConfOption(c.cc.GetXConfOption()...),
	)
	cc.ApplyOption(opts...)
	sub := NewCommandWithConfig(name, cc)
	c.AddCommand(sub)
	return sub
}

func (c *Command) Check() error {
	for _, v := range c.commands {
		binder := c.cc.GetParser()
		if binder == nil {
			return errors.New("need Parser")
		}
		ff := flag.NewFlagSet(strings.Join(v.usageNamePath, "/"), flag.ContinueOnError)
		err := binder(context.Background(), v, ff, nil, func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
			return nil
		})
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
