package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sandwich-go/xconf/xcmd"
	"github.com/sandwich-go/xconf/xcmd/main/xcmdtest"
)

func main() {
	cc := xcmdtest.NewConfig()
	xcmd.Config().ApplyOption(xcmd.WithBind(cc))

	xcmd.Use(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string, next xcmd.Executer) error {
		return next(ctx, c, ff, args)
	})

	// sub命令export，继承上游命令的bind信息
	// export 派生go命令，只绑定http_address字段
	// go 派生export命令，追加绑定timeouts字段
	xcmd.CommandInheritBind("export",
		xcmd.WithSynopsis("export proto to golang/cs/python/lua"),
		xcmd.WithExecute(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string) error {
			fmt.Println("export command")
			return nil
		}),
	).Use(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string, next xcmd.Executer) error {
		return next(ctx, c, ff, args)
	}).
		CommandInheritBind("go",
			xcmd.WithBindFieldPath("http_address"),
			xcmd.WithSynopsis("generate golang code"),
		).Use(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string, next xcmd.Executer) error {
		return next(ctx, c, ff, args)
	}).
		CommandInheritBind("service", xcmd.WithBindFieldPathAppend("timeouts"))

	// sub命令log绑定到新的配置项
	anotherBind := xcmdtest.NewLog()
	xcmd.AddCommand(xcmd.NewCommand("log",
		xcmd.WithBind(anotherBind),
		xcmd.WithSynopsis("log command"),
		xcmd.WithExecute(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string) error {
			fmt.Println("log command")
			return nil
		})))
	// sub命令同样绑定到xcmdtest.Config实例cc
	xcmd.AddCommand(xcmd.NewCommand("layout",
		xcmd.WithBind(cc),
		xcmd.WithSynopsis("layout command"),
		xcmd.WithExecute(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string) error {
			fmt.Println("layout command")
			return nil
		})))

	// 手动绑定
	logLevel := 0
	binding := func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string, next xcmd.Executer) error {
		ff.IntVar(&logLevel, "log_level_manual", logLevel, "set log level")
		return next(ctx, c, ff, args)
	}

	manual := xcmd.NewCommand("manual",
		xcmd.WithBind(cc),
		xcmd.WithSynopsis("manual bing flag"),
		xcmd.WithExecute(func(ctx context.Context, c *xcmd.Command, ff *flag.FlagSet, args []string) error {
			fmt.Println("manual command got log_level:", logLevel)
			return nil
		}))
	xcmd.AddCommand(manual, binding)

	panicPrintErr("comamnd check with err: %v", xcmd.Check())
	panicPrintErr("comamnd Execute with err: %v", xcmd.Execute(context.Background(), os.Args[1:]...))

}

func panicPrintErr(format string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", err)
}
