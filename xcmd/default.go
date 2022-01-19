package xcmd

import (
	"context"
	"flag"
	"io"
	"os"
	"path"
)

var rootCmd = NewCommand(path.Base(os.Args[0]), WithExecute(func(ctx context.Context, c *Command, ff *flag.FlagSet, args []string) error {
	ff.Usage()
	return nil
}))

// Use 添加中间件，在此之后添加的子命令都会继承该中间件
// 执行顺序为：preMiddleware -> Parser -> middleware -> Executer
func Use(middleware ...MiddlewareFunc) *Command { return rootCmd.Use(middleware...) }

// UsePre 添加preMiddleware中间件，pre中间件运行在Parser之前
// 执行顺序为：preMiddleware -> Parser -> middleware -> Executer
func UsePre(middleware ...MiddlewareFunc) *Command { return rootCmd.UsePre(middleware...) }

// Add 添加一条子命令
func Add(name string, opts ...ConfigOption) *Command { return rootCmd.Add(name, opts...) }

// AddCommand 添加一条子命令，可以携带中间件信息，等同于Add(xxxxx).Use或者AddCommand(xxxx).Use
func AddCommand(sub *Command, middleware ...MiddlewareFunc) { rootCmd.AddCommand(sub, middleware...) }

// Config 获取配置，允许运行期调整，但只在Parser运行前生效
func Config() ConfigInterface { return rootCmd.Config() }

// Execute 执行参数解析驱动命令执行
func Execute(ctx context.Context, args ...string) error { return rootCmd.Execute(ctx, args...) }

// Explain 打印使用说明
func Explain(w io.Writer) { rootCmd.Explain(w) }

// Check 检查当前命令及子命令是否有路径绑定错误等信息
func Check() error { return rootCmd.Check() }

// SubCommand 由当前命令扩展子命令, 继承Bing，BindPath,XConfOption等参数
func SubCommand(name string, opts ...ConfigOption) *Command {
	return rootCmd.SubCommand(name, opts...)
}

// SetRootCommand 设定当前默认的根命令
func SetRootCommand(root *Command) { rootCmd = root }

// RootCommand 获取根命令
func RootCommand() *Command { return rootCmd }
