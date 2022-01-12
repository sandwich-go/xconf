package xflag

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/sandwich-go/xconf/xutil"
)

// FlagInfo flag basic info
type FlagInfo struct {
	Name     string
	TypeName string
	Usage    string
	DefValue string
	Flag     *flag.Flag
}

// FlagList FlagInfo list
type FlagList []*FlagInfo

func (f FlagList) Flag(name string) *FlagInfo {
	for _, v := range f {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// GetFlagInfo get FlagList from given FlagSet
func GetFlagInfo(f *flag.FlagSet) (ret FlagList) {
	f.VisitAll(func(ff *flag.Flag) {
		varname, usage := flag.UnquoteUsage(ff)
		if varname == "" || varname == "value" {
			if t, ok := ff.Value.(interface{ TypeName() string }); ok {
				varname = t.TypeName()
			}
			if t, ok := ff.Value.(interface{ IsBoolFlag() bool }); ok && t.IsBoolFlag() {
				varname = "bool"
			}
		}
		v := &FlagInfo{
			Name:     ff.Name,
			Usage:    usage,
			TypeName: varname,
			DefValue: ff.DefValue,
			Flag:     ff,
		}
		ret = append(ret, v)
	})
	return ret
}

// PrintDefaults 打印FlagSet，替换默认实现
func PrintDefaults(f *flag.FlagSet, suffixLines ...string) {
	lines := make([]string, 0)
	magic := "\x00"
	lines = append(lines, fmt.Sprintf("FLAG%sTYPE%sUSAGE", magic, magic))

	allFlag := GetFlagInfo(f)
	for _, v := range allFlag {
		line := ""
		line = fmt.Sprintf("--%s", v.Name)
		line += magic
		line += v.TypeName
		line += magic
		line += v.Usage
		if !isZeroValue(v.Flag, v.DefValue) {
			if v.TypeName == "string" {
				line += fmt.Sprintf(" (default %q)", v.DefValue)
			} else {
				line += fmt.Sprintf(" (default %s)", v.DefValue)
			}
		}
		lines = append(lines, line)
	}
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
