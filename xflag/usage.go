package xflag

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

// PrintDefaults 打印FlagSet，替换默认实现
func PrintDefaults(f *flag.FlagSet) {
	buf := new(bytes.Buffer)
	lines := make([]string, 0)
	maxlenName := 0
	maxlenNameVarName := 0
	f.VisitAll(func(ff *flag.Flag) {
		line := ""
		line = fmt.Sprintf("      --%s", ff.Name)
		varname, usage := flag.UnquoteUsage(ff)
		if varname == "" || varname == "value" {
			if t, ok := ff.Value.(interface{ TypeName() string }); ok {
				varname = t.TypeName()
			}
			if t, ok := ff.Value.(interface{ IsBoolFlag() bool }); ok && t.IsBoolFlag() {
				varname = "bool"
			}
		}
		line += "\x00"
		if len(line) > maxlenName {
			maxlenName = len(line)
		}
		if len(varname) > 0 {
			line += "  " + varname
		}
		line += "\x01"
		if len(line) > maxlenNameVarName {
			maxlenNameVarName = len(line)
		}

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
	for _, line := range lines {
		{
			sidx := strings.Index(line, "\x00")
			spacing := strings.Repeat(" ", maxlenName-sidx)
			line = line[:sidx] + spacing + line[sidx+1:]
		}
		{
			sidx := strings.Index(line, "\x01")
			spacing := strings.Repeat(" ", maxlenNameVarName-sidx+2)
			line = line[:sidx] + spacing + line[sidx+1:]
		}
		fmt.Fprintln(buf, line)
	}
	fmt.Fprint(f.Output(), buf.String(), "\n")
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
