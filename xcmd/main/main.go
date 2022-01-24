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
	xcmd.BindSet(cc)

	xcmd.Use(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		return next(ctx, c)
	})

	// sub命令export，继承上游命令的bind信息
	// export 派生go命令，只绑定http_address字段
	// go 派生export命令，追加绑定timeouts字段
	xcmd.SubCommand("export",
		xcmd.WithShort("export proto to golang/cs/python/lua"),
		xcmd.WithExamples(`
		只设定http: xcmd export --http_address=10.0.0.1
		只设定timeout: xcmd export --timeouts=read,6s`),
		xcmd.WithDescription(`支持默认值配置、解析
		支持多种格式，内置JSON, TOML, YAML，FLAG, ENV支持，并可注册解码器扩展格式支持
		支持多文件、多io.Reader数据加载，支持文件继承
		支持由OS ENV变量数据加载配置
		支持由命令行参数FLAGS加载数据`),
	).Use(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
		return next(ctx, c)
	}).SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
		fmt.Println("export command")
		return nil
	}).SubCommand("go", xcmd.WithShort("generate golang code")).
		Use(func(ctx context.Context, c *xcmd.Command, next xcmd.Executer) error {
			return next(ctx, c)
		}).
		BindFieldPathSet("http_address").
		SubCommand("service").
		BindFieldPathAdd("timeouts")

	// sub命令log绑定到新的配置项
	anotherBind := xcmdtest.NewLog()
	xcmd.AddCommand(xcmd.NewCommand("log",
		xcmd.WithShort("log command")).
		SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
			fmt.Println("log command")
			return nil
		}).BindSet(anotherBind))
	// sub命令同样绑定到xcmdtest.Config实例cc
	xcmd.AddCommand(xcmd.NewCommand("layout",
		xcmd.WithShort("layout command")).SetExecuter(func(ctx context.Context, c *xcmd.Command) error {
		fmt.Println("layout command")
		return nil
	}).BindSet(cc))

	// 手动绑定
	logLevel := 0
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

	xcmd.Check()
	// panicPrintErr("comamnd check with err: %v", xcmd.Check())
	panicPrintErr("comamnd Execute with err: %v", xcmd.Execute(context.Background(), os.Args[1:]...))
}

func panicPrintErr(format string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", err)
}
