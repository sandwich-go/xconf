package xflag

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/sandwich-go/xconf/xutil"
)

// PrintDefaults 打印FlagSet，替换默认实现
func PrintDefaults(f *flag.FlagSet, suffixLines ...string) {
	lines := make([]string, 0)
	magic := "\x00"
	lines = append(lines, fmt.Sprintf("FLAG%sTYPE%sUSAGE", magic, magic))
	f.VisitAll(func(ff *flag.Flag) {
		line := ""
		line = fmt.Sprintf("--%s", ff.Name)
		varname, usage := flag.UnquoteUsage(ff)
		if varname == "" || varname == "value" {
			if t, ok := ff.Value.(interface{ TypeName() string }); ok {
				varname = t.TypeName()
			}
			if t, ok := ff.Value.(interface{ IsBoolFlag() bool }); ok && t.IsBoolFlag() {
				varname = "bool"
			}
		}
		line += magic
		if len(varname) > 0 {
			line += varname
		}
		line += magic
		line += usage
		if !isZeroValue(ff, ff.DefValue) {
			if varname == "string" {
				line += fmt.Sprintf(" (default %q)", ff.DefValue)
			} else {
				line += fmt.Sprintf(" (default %s)", ff.DefValue)
			}
		}
		lines = append(lines, line)
	})
	fmt.Fprint(f.Output(), xutil.TableFormat(lines, magic, suffixLines...), "\n")
}

func isZeroValue(f *flag.Flag, value string) bool {
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(flag.Value).String()
}
