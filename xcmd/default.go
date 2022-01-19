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

func Use(middleware ...MiddlewareFunc) *Command             { return rootCmd.Use(middleware...) }
func UsePre(middleware ...MiddlewareFunc) *Command          { return rootCmd.UsePre(middleware...) }
func AddCommand(sub *Command, middleware ...MiddlewareFunc) { rootCmd.AddCommand(sub, middleware...) }
func Config() ConfigInterface                               { return rootCmd.cc }
func Execute(ctx context.Context, args ...string) error     { return rootCmd.Execute(ctx, args...) }
func Explain(w io.Writer)                                   { rootCmd.Explain(w) }
func Check() error                                          { return rootCmd.Check() }

func SubCommand(name string, opts ...ConfigOption) *Command {
	return rootCmd.SubCommand(name, opts...)
}

func SetRootCommand(root *Command) { rootCmd = root }
func RootCommand() *Command        { return rootCmd }
