package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sandwich-go/xconf/xcmd"
	"github.com/sandwich-go/xconf/xcmd/main/xcmdtest"
)

func main() {
	cc := xcmdtest.NewConfig()
	xcmd.Config().ApplyOption(
		xcmd.WithAuthor("who_am_i@google.com", "i_am_who@google.com"),
		xcmd.WithDescription("xcmd依托xconf自动完成flag参数创建，绑定，解析等操作，同时支持自定义flag，支持中间件，支持子命令."),
	)

	xcmd.UsePre(func(ctx context.Context, cmd *xcmd.Command, next xcmd.Executer) error {
		// 只应该在这里定义flag等一些无副作用的操作，不论Executer是否为空都会执行
		return next(ctx, cmd)
	}).Use(func(ctx context.Context, cmd *xcmd.Command, next xcmd.Executer) error {
		fmt.Println("只有Executer不为空的时候才会执行到这里")
		return next(ctx, cmd)
	})

	logLevel := 0
	// 将Cofnig绑定到xcmd根命令
	xcmd.BindSet(cc)
	// 自定义一个全局参数，需要使用UsePre插入中间件，在参数解析前创建Flag
	xcmd.UsePre(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		c.FlagSet.IntVar(&logLevel, "log_level", logLevel, "set log level")
		// 记得继续派发
		return next(ctx, c)
	})
	cmdExport := xcmd.SubCommand("export", // 添加一个子命令
		xcmd.WithShort("export proto to golang/cs/python/lua"),
		xcmd.WithDescription(`详细描述一下，在执行: export --help时会显示该消息，并可以换行,内部会自动格式化`),
		xcmd.WithExamples(`- 只设定http: export --http_address=10.0.0.1
		- 只设定timeout: export --timeouts=read,6s`))

	cmdExport.UsePre(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		fmt.Println("export middlewar pre,UsePre插入的中间件，此时参数还没有解析")
		return next(ctx, c)
	}).Use(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		fmt.Println("export middlewar,Use插入的中间件，此时参数已经解析")
		return next(ctx, c)
	}).SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
		fmt.Println("export command")
		return nil
	}).UsePre(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		fmt.Println("由于这个中间件添加在SetExecuter之后，所以不会在export中执行，会在其sub command中执行")
		return next(ctx, c)
	})

	cmdExport.SubCommand("go",
		xcmd.WithShort("generate golang code"),
	).Use(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		return next(ctx, c)
	}).SetExecuter(func(ctx context.Context, cmd *xcmd.Command) error {
		fmt.Println("go command")
		return nil
	}).BindFieldPathSet("http_address"). //设定只绑定到http_address上
						SubCommand("service").       //扩展一个字命令service
						BindFieldPathAdd("timeouts") //service子命令在go命令基础上再绑定一个timeouts字段

	// sub命令log绑定到新的配置项
	anotherBind := xcmdtest.NewLog()
	xcmd.NewCommand("log",
		xcmd.WithShort("log command"),
	).SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
		fmt.Println("log command")
		return nil
	}).AddTo(xcmd.RootCommand()).BindSet(anotherBind)

	// sub命令同样绑定到xcmdtest.Config实例cc
	xcmd.NewCommand("layout",
		xcmd.WithShort("layout command"),
		// 设定为Deprecated
		xcmd.WithDeprecated("do not use this again, use export"),
	).SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
		fmt.Println("layout command")
		return nil
	}).BindSet(cc).AddTo(xcmd.RootCommand())

	// 手动绑定
	binding := func(ctx context.Context, cmd *xcmd.Command, next xcmd.Executer) error {
		cmd.FlagSet.IntVar(&logLevel, "log_level_manual", logLevel, "set log level")
		return next(ctx, cmd)
	}

	manual := xcmd.NewCommand("manual", xcmd.WithShort("manual bing flag")).BindSet(cc)
	manual.UsePre(binding).SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
		fmt.Println("manual command got log_level:", logLevel)
		return nil
	}).SubCommand("export_go").SubCommand("export_go_service")
	xcmd.AddCommand(manual)

	xcmd.Add("empty")

	panicPrintErr("comamnd Execute with err: %v", xcmd.Execute(context.Background(), os.Args[1:]...))
}

func panicPrintErr(format string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", err)
}
