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
type FlagList struct {
	List []*FlagInfo
}

// Flag return FlagInfo by name
func (f *FlagList) Flag(name string) *FlagInfo {
	for _, v := range f.List {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// FlagGetAndDel return FlagInfo by name and del it
func (f *FlagList) FlagGetAndDel(name string) *FlagInfo {
	for index, v := range f.List {
		if v.Name == name {
			f.List = append(f.List[:index], f.List[index+1:]...)
			return v
		}
	}
	return nil
}

// UnquoteUsage linke flag.UnquoteUsage
func UnquoteUsage(ff *flag.Flag) (name string, usage string) {
	varname, usage := flag.UnquoteUsage(ff)
	if t, ok := ff.Value.(interface{ TypeName() string }); ok {
		varname = t.TypeName()
	}
	if t, ok := ff.Value.(interface{ IsBoolFlag() bool }); ok && t.IsBoolFlag() {
		varname = "bool"
	}
	return varname, usage
}

// GetFlagInfo get FlagList from given FlagSet
func GetFlagInfo(f *flag.FlagSet) (ret FlagList) {
	f.VisitAll(func(ff *flag.Flag) {
		varname, usage := UnquoteUsage(ff)
		v := &FlagInfo{
			Name:     ff.Name,
			Usage:    usage,
			TypeName: varname,
			DefValue: ff.DefValue,
			Flag:     ff,
		}
		ret.List = append(ret.List, v)
	})
	return ret
}

// PrintDefaults 打印FlagSet，替换默认实现
func PrintDefaults(f *flag.FlagSet, suffixLines ...string) {
	lines := make([]string, 0)
	magic := "\x00"
	lines = append(lines, "FLAG"+"\x00"+"TYPE"+"\x00"+"USAGE")

	allFlag := GetFlagInfo(f)
	for _, v := range allFlag.List {
		line := ""
		line = fmt.Sprintf("--%s", v.Name)
		line += magic
		line += v.TypeName
		line += magic
		line += fmt.Sprintf("|%s| %s", "-", v.Usage)
		if !IsZeroValue(v.Flag, v.DefValue) {
			if v.TypeName == "string" {
				line += fmt.Sprintf(" (default %q)", v.DefValue)
			} else {
				line += fmt.Sprintf(" (default %s)", v.DefValue)
			}
		}
		lines = append(lines, line)
	}
	fmt.Fprint(f.Output(), xutil.TableFormat(lines, magic, true, suffixLines...), "\n")
}

// IsZeroValue check is zero val
func IsZeroValue(f *flag.Flag, value string) bool {
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(flag.Value).String()
}
