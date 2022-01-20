package xcmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sandwich-go/xconf/xflag"
)

// MiddlewareFunc 中间件方法
// cmd *Command : 为当前执行的命令对象
// ff *flag.FlagSet: 为当前命令对象解析中使用的FlagSet,如为pre中间件则未解析，否则为解析过后的FlagSet
// args []string : 为当前命令行参数，也就是FlagSet要去解析的参数列表
// next : 下一步要执行的方法，可能是下一个中间件或者目标Executer方法
type MiddlewareFunc = func(ctx context.Context, cmd *Command, ff *flag.FlagSet, args []string, next Executer) error

// Command 表征一条命令或一个命令Group
// 中间件分preMiddleware与middleware
// 执行顺序为：preMiddleware -> Parser -> middleware -> Executer
// - preMiddleware执行的时候并没有进行任何的FlagSet解析，可以在此时进行一些自定义的Flag创建
// - Parser为Option传入的中间件，一般无需自行实现，Parser主要是将配置文件、FlagSet、Env等遵循XConf规则解析到Option传入的Bind对象上，如传nil，则会调用FlagSet的Parse方法
// - middleware执行的时候已经完成了参数对象的绑定解析
type Command struct {
	name          string
	cc            ConfigInterface
	Output        io.Writer
	commands      []*Command
	middleware    []MiddlewareFunc
	preMiddleware []MiddlewareFunc
	usageNamePath []string
	parent        *Command // 目前只用于确认命令是否已经有父节点
}

// NewCommand 创建一条命令
func NewCommand(name string, opts ...ConfigOption) *Command {
	return NewCommandWithConfig(name, NewConfig(opts...))
}

// NewCommandWithConfig 创建一条命令
func NewCommandWithConfig(name string, cc ConfigInterface) *Command {
	c := &Command{
		name:   name,
		cc:     cc,
		Output: os.Stdout,
	}
	c.usageNamePath = []string{name}
	return c
}

// Use 添加中间件，在此之后添加的子命令都会继承该中间件
// 执行顺序为：preMiddleware -> Parser -> middleware -> Executer
func (c *Command) Use(middleware ...MiddlewareFunc) *Command {
	c.middleware = append(c.middleware, middleware...)
	return c
}

// UsePre 添加preMiddleware中间件，pre中间件运行在Parser之前
// 执行顺序为：preMiddleware -> Parser -> middleware -> Executer
func (c *Command) UsePre(preMiddleware ...MiddlewareFunc) *Command {
	c.preMiddleware = append(c.preMiddleware, preMiddleware...)
	return c
}

// AddCommand 添加一条子命令，可以携带中间件信息，等同于Add(xxxxx).Use或者AddCommand(xxxx).Use
func (c *Command) AddCommand(sub *Command, middleware ...MiddlewareFunc) {
	if sub.parent != nil {
		panic(fmt.Sprintf("command:%s has parent:%s", sub.name, sub.parent.name))
	}
	sub.usageNamePath = append(c.usageNamePath, sub.usageNamePath...)
	sub.middleware = combineMiddlewareFunc(c.middleware, middleware...)
	sub.preMiddleware = combineMiddlewareFunc(c.preMiddleware, sub.preMiddleware...)
	sub.parent = c
	// 如果该命令在添加子命令前没有父节点，则需要将父节点的中间件追加上
	for _, v := range sub.commands {
		v.middleware = combineMiddlewareFunc(c.middleware, v.middleware...)
		v.preMiddleware = combineMiddlewareFunc(c.preMiddleware, v.preMiddleware...)
	}
	c.commands = append(c.commands, sub)
}

func combineMiddlewareFunc(middlewareNow []MiddlewareFunc, middleware ...MiddlewareFunc) []MiddlewareFunc {
	m := make([]MiddlewareFunc, 0, len(middlewareNow)+len(middleware))
	m = append(m, middlewareNow...)
	m = append(m, middleware...)
	return m
}

// Config 获取配置，允许运行期调整，但只在Parser运行前生效
func (c *Command) Config() ConfigInterface { return c.cc }

// Add 添加一条子命令
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

// Execute 执行参数解析驱动命令执行
func (c *Command) Execute(ctx context.Context, args ...string) error {
	if len(args) != 0 {
		// 尝试在当前命令集下寻找子命令
		subCommandName := args[0]
		for _, cmd := range c.commands {
			if cmd.Name() != subCommandName {
				continue
			}
			return cmd.Execute(ctx, args[1:]...)
		}
	}
	ff := flag.NewFlagSet(strings.Join(c.usageNamePath, "/"), flag.ContinueOnError)
	// 默认 usage 无参
	ff.Usage = func() {
		c.Explain(c.Output)
		fmt.Fprintf(c.Output, "Flags:\n")
		xflag.PrintDefaults(ff)
	}

	var allMiddlewares []MiddlewareFunc
	allMiddlewares = append(allMiddlewares, c.preMiddleware...)
	allMiddlewares = append(allMiddlewares, parser)
	allMiddlewares = append(allMiddlewares, c.middleware...)
	err := ChainMiddleware(allMiddlewares...)(ctx, c, ff, args, exec)
	if err != nil && IsErrHelp(err) {
		return nil
	}
	return err
}

const HasParsed = "xcmd_has_parsed"

func parser(ctx context.Context, c *Command, ff *flag.FlagSet, args []string, next Executer) error {
	if v := ctx.Value(HasParsed); v != nil {
		return next(ctx, c, ff, args)
	}
	if c.cc.GetBind() == nil {
		err := ff.Parse(args)
		if err != nil {
			return err
		}
		return next(ctx, c, ff, args)
	}
	return c.cc.GetParser()(ctx, c, ff, args, next)
}

func exec(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
	executer := c.Config().GetExecute()
	if executer == nil {
		executer = c.cc.GetOnExecuterLost()
	}
	return executer(context.Background(), c, ff, args)
}

// Name 获取当前命令名称
func (c *Command) Name() string { return c.name }

// NamePath 获取当前命令路径
func (c *Command) NamePath() []string { return c.usageNamePath }

// Usage 获取当前命令路径
func (c *Command) Usage() string { return c.cc.GetUsage() }

// Short 获取当前命令Short Usage
func (c *Command) Short() string { return c.cc.GetShort() }

// SubCommand 由当前命令扩展子命令, 继承Bing，BindPath,XConfOption等参数
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

// Check 检查当前命令及子命令是否有路径绑定错误等信息
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
