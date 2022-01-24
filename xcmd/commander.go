package xcmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xutil"
)

// ErrHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.
var ErrHelp = flag.ErrHelp

// IsErrHelp 检查错误是否是ErrHelp
func IsErrHelp(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), flag.ErrHelp.Error())
}

// Command 表征一条命令或一个命令Group
// 中间件分preMiddleware与middleware
// 执行顺序为：preMiddleware -> Parser -> middleware -> Executer
// - preMiddleware执行的时候并没有进行任何的FlagSet解析，可以在此时进行一些自定义的Flag创建
// - Parser为Option传入的中间件，一般无需自行实现，Parser主要是将配置文件、FlagSet、Env等遵循XConf规则解析到Option传入的Bind对象上，如传nil，则会调用FlagSet的Parse方法
// - middleware执行的时候已经完成了参数对象的绑定解析
type Command struct {
	name     string
	cc       *config
	Output   io.Writer
	commands []*Command

	usageNamePath []string
	parent        *Command // 目前只用于确认命令是否已经有父节点

	executer              Executer
	executerMiddleware    []MiddlewareFunc // 记录设定Executer时的中间件，防止之后加入的中间件作用于Executer
	executerMiddlewarePre []MiddlewareFunc

	// 记录当前command上挂载的中间件
	middleware    []MiddlewareFunc
	middlewarePre []MiddlewareFunc

	bind          interface{} // 命令绑定的参数结构
	bindFieldPath []string    // 命令绑定的参数FieldPath,如空则全部绑定

	FlagArgs []string // 除去command名称的原始参数
	FlagSet  *flag.FlagSet
	usage    func()

	// 缓存记录由parent继承而来的flag
	flagInheritByMiddlewarePre []string
	flagLocal                  []string
}

// NewCommand 创建一条命令
func NewCommand(name string, opts ...ConfigOption) *Command {
	return newCommandWithConfig(name, NewConfig(opts...))
}

// newCommandWithConfig 创建一条命令
func newCommandWithConfig(name string, cc *config) *Command {
	c := &Command{
		name:   name,
		cc:     cc,
		Output: os.Stdout,
	}
	c.usageNamePath = []string{name}
	c.middlewarePre = append(c.middlewarePre, preMiddlewareBegin)
	c.updateUsage(nil)
	return c
}

func (c *Command) newXConf() *xconf.XConf {
	cc := xconf.NewOptions(
		xconf.WithErrorHandling(xconf.ContinueOnError),
		xconf.WithFlagSet(c.FlagSet),
		xconf.WithFlagArgs(c.FlagArgs...))
	cc.ApplyOption(c.Config().GetXConfOption()...)
	x := xconf.NewWithConf(cc)
	return x
}

func preMiddlewareBegin(ctx context.Context, cmd *Command, next Executer) error {
	cmd.FlagSet.VisitAll(func(f *flag.Flag) {
		cmd.flagInheritByMiddlewarePre = append(cmd.flagInheritByMiddlewarePre, f.Name)
	})
	return next(ctx, cmd)
}

func preMiddlewareEnd(ctx context.Context, cmd *Command, next Executer) error {
	var nowFlags []string
	cmd.FlagSet.VisitAll(func(f *flag.Flag) {
		nowFlags = append(nowFlags, f.Name)
	})
	for _, v := range nowFlags {
		if xutil.ContainString(cmd.flagInheritByMiddlewarePre, v) {
			continue
		}
		cmd.flagLocal = append(cmd.flagLocal, v)
	}
	return next(ctx, cmd)
}

// Bind 获取绑定的对象
func (c *Command) Bind() interface{} { return c.bind }

// BindSet 设定参数绑定的对象，只在解析之前生效,并重置绑定
func (c *Command) BindSet(xconfVal interface{}) *Command {
	c.bind = xconfVal
	if c.bind == nil {
		c.bindFieldPath = xconf.FieldPathList(c.bind, c.newXConf())
	}
	return c
}

// BindFieldPathSet 设定绑定的参数FieldPath，只在解析之前生效
func (c *Command) BindFieldPathSet(filePath ...string) *Command {
	c.bindFieldPath = filePath
	return c
}

// BindFieldPathAdd 设定绑定的参数FieldPath，只在解析之前生效
func (c *Command) BindFieldPathAdd(filePath ...string) *Command {
	for _, v := range filePath {
		if xutil.ContainString(c.bindFieldPath, v) {
			continue
		}
		c.bindFieldPath = append(c.bindFieldPath, v)
	}
	return c
}

// BindFieldPathClean 清空绑定路径
func (c *Command) BindFieldPathClean() *Command {
	c.bindFieldPath = nil
	return c
}

// BindFieldPath 返回绑定额路径列表
func (c *Command) BindFieldPath() []string { return c.bindFieldPath }

// BindFieldPathReomove 移除部分绑定路径
func (c *Command) BindFieldPathReomove(filePath ...string) *Command {
	for i := 0; i < len(c.bindFieldPath); i++ {
		if xutil.ContainString(filePath, c.bindFieldPath[i]) {
			c.bindFieldPath = append(c.bindFieldPath[:i], c.bindFieldPath[i+1:]...)
			i--
		}
	}
	return c
}

// SetExecuter 设定新的Executer，会缓存此时的中间件，只有此时缓存的中间件会被应用到Executer，如果Executer为nil，则所有的中间件都会被应用到默认的Executer
func (c *Command) SetExecuter(executer Executer) *Command {
	c.executer = executer
	// 记录设定Executer时刻的中间件
	c.executerMiddleware = combineMiddlewareFunc(c.middleware)
	c.executerMiddlewarePre = combineMiddlewareFunc(c.middlewarePre)
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
	c.middlewarePre = append(c.middlewarePre, preMiddleware...)
	return c
}

// AddCommand 添加一条子命令，可以携带中间件信息，等同于Add(xxxxx).Use或者AddCommand(xxxx).Use
func (c *Command) AddCommand(sub *Command, middleware ...MiddlewareFunc) {
	if sub.parent != nil {
		panic(fmt.Sprintf("command:%s has parent:%s", sub.name, sub.parent.name))
	}
	if c == sub {
		panic("same command")
	}
	sub.usageNamePath = append(c.usageNamePath, sub.usageNamePath...)

	sub.middlewarePre = combineMiddlewareFunc(c.middlewarePre, sub.middlewarePre...)
	sub.executerMiddlewarePre = combineMiddlewareFunc(c.middlewarePre, sub.executerMiddlewarePre...)
	sub.middleware = combineMiddlewareFunc(c.middleware, middleware...)
	sub.executerMiddleware = combineMiddlewareFunc(c.middleware, sub.executerMiddleware...)

	sub.parent = c
	// 如果该命令在添加子命令前没有父节点，则需要将父节点的中间件追加上
	for _, v := range sub.commands {
		v.middlewarePre = combineMiddlewareFunc(c.middlewarePre, v.middlewarePre...)
		v.executerMiddlewarePre = combineMiddlewareFunc(c.middlewarePre, v.executerMiddlewarePre...)
		v.middleware = combineMiddlewareFunc(c.middleware, v.middleware...)
		v.executerMiddleware = combineMiddlewareFunc(c.middleware, v.executerMiddleware...)
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

// AddTo 为parent添加一条子指令
func (c *Command) AddTo(parent *Command) *Command {
	parent.AddCommand(c)
	return c
}

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
func isFlagArg(arg string) bool {
	return ((len(arg) >= 3 && arg[1] == '-') || (len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
}

// Execute 执行参数解析驱动命令执行
func (c *Command) Execute(ctx context.Context, args ...string) error {
	// 存储原始的参数数据，主要debug使用
	c.FlagSet = flag.NewFlagSet(c.name, flag.ContinueOnError)
	c.FlagArgs = args
	var argFirst string
	if len(args) != 0 {
		// 以index=0的元素作为命令名尝试寻找subcommand
		argFirst = args[0]
		for _, cmd := range c.commands {
			if cmd.Name() != argFirst {
				continue
			}
			return cmd.Execute(ctx, args[1:]...)
		}
		// 没有找到进一步的子命令，但此时的args[0]可能是输入错的子命令，也可能是当前命令的arguments或者flag
	}
	c.FlagSet.Usage = c.usage
	var executerMiddleware []MiddlewareFunc
	if c.executer == nil {
		executerMiddleware = append(executerMiddleware, c.middlewarePre...)
		executerMiddleware = append(executerMiddleware, preMiddlewareEnd)
		executerMiddleware = append(executerMiddleware, parser)
		executerMiddleware = append(executerMiddleware, c.middleware...)
	} else {
		executerMiddleware = append(executerMiddleware, c.executerMiddlewarePre...)
		executerMiddleware = append(executerMiddleware, preMiddlewareEnd)
		executerMiddleware = append(executerMiddleware, parser)
		executerMiddleware = append(executerMiddleware, c.executerMiddleware...)
	}
	err := ChainMiddleware(executerMiddleware...)(ctx, c, exec)
	if err != nil && IsErrHelp(err) {
		// 当执行返回的是ErrHelp时，说明当前可能是一个父命令且未设定有效Executer，猜测此时的args[0]可能是输入错误的子命令
		if argFirst != "" && !isFlagArg(argFirst) {
			if suggestions := c.suggestionsFor(argFirst); len(suggestions) > 0 {
				fmt.Fprintf(c.Output, "\n%q is not a subcommand. Did you mean this?\n%s\n", argFirst, strings.Join(suggestions, "\n"))
			}
		}
		return nil
	}
	return err
}

func parser(ctx context.Context, cmd *Command, next Executer) error {
	if cmd.bind == nil {
		err := cmd.FlagSet.Parse(cmd.FlagArgs)
		if err != nil {
			return err
		}
		return next(ctx, cmd)
	}
	return cmd.cc.GetParser()(ctx, cmd, next)
}

func exec(ctx context.Context, cmd *Command) error {
	executer := cmd.executer
	if executer == nil {
		executer = cmd.cc.GetOnExecuterLost()
	}
	return executer(context.Background(), cmd)
}

// Name 获取当前命令名称
func (c *Command) Name() string { return c.name }

// NamePath 获取当前命令路径
func (c *Command) NamePath() []string { return c.usageNamePath }

// Usage
func (c *Command) Usage() { c.usage() }

// Short 获取当前命令Short Usage
func (c *Command) Short() string { return c.cc.GetShort() }

// SubCommand 由当前命令扩展子命令, 继承Bing，BindPath,XConfOption等参数
func (c *Command) SubCommand(name string, opts ...ConfigOption) *Command {
	config := NewConfig(WithXConfOption(c.cc.XConfOption...))
	config.ApplyOption(opts...)
	sub := newCommandWithConfig(name, config)
	sub.bind = c.bind
	sub.bindFieldPath = c.bindFieldPath
	return sub.AddTo(c)
}

// Check 检查当前命令及子命令是否有路径绑定错误等信息. 调试使用
func (c *Command) Check() error {
	for _, v := range c.commands {
		// 替换executer防止检查过程中的执行，输出
		executer := v.executer
		v.executer = func(ctx context.Context, cmd *Command) error { return nil }
		err := v.Execute(context.Background())
		if err != nil {
			return err
		}
		v.executer = executer
	}
	return nil
}

// suggestionsFor provides suggestions for the naem.
func (c *Command) suggestionsFor(naem string) []string {
	suggestions := []string{}
	for _, cmd := range c.commands {
		ld := xutil.LD(naem, cmd.Name(), true)
		suggestByLevenshtein := ld <= c.cc.SuggestionsMinDistance
		suggestByPrefix := strings.HasPrefix(strings.ToLower(cmd.Name()), strings.ToLower(naem))
		if suggestByLevenshtein || suggestByPrefix {
			suggestions = append(suggestions, strings.Join(cmd.usageNamePath, " "))
		}
		suggestions = append(suggestions, cmd.suggestionsFor(naem)...)
	}
	return suggestions
}
