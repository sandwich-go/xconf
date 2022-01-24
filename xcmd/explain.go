package xcmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

type byGroupName []*Command

func (p byGroupName) Len() int           { return len(p) }
func (p byGroupName) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p byGroupName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Explain 打印使用说明
func (c *Command) Explain(w io.Writer) { explainGroup(w, c) }

// explainGroup explains all the subcommands for a particular group.
func explainGroup(w io.Writer, c *Command) {
	if len(c.commands) == 0 {
		fmt.Fprintf(w, "USAGE: \n%s%s <flags> <args>\n\n", paddingContent, strings.Join(c.usageNamePath, " "))
	} else {
		fmt.Fprintf(w, "USAGE: \n%s%s <subcommand> <flags> <args>\n\n", paddingContent, strings.Join(c.usageNamePath, " "))
	}
	if len(c.commands) == 0 {
		return
	}
	if c.cc.Description != "" {
		fmt.Fprintf(w, "DESCRIPTION:\n")
		usageLine := xutil.StringSliceWalk(strings.Split(xutil.StringTrim(c.cc.Description), "\n"), func(s string) (string, bool) {
			return paddingContent + xutil.StringTrim(s), true
		})
		fmt.Fprintf(w, "%s\n\n", strings.Join(usageLine, "\n"))
	}
	if c.cc.Examples != "" {
		fmt.Fprintf(w, "EXAMPLES:\n")
		examplesLines := xutil.StringSliceWalk(strings.Split(xutil.StringTrim(c.cc.Examples), "\n"), func(s string) (string, bool) {
			return paddingContent + xutil.StringTrim(s), true
		})
		fmt.Fprintf(w, "%s\n\n", strings.Join(examplesLines, "\n"))
	}
	sort.Sort(byGroupName(c.commands))
	fmt.Fprintf(w, "AVAIABLE COMMANDS:\n")
	sort.Sort(byGroupName(c.commands))
	var level = []bool{}
	lines := printCommand(c, level)
	fmt.Fprintln(w, strings.Join(lines, "\n"))
	fmt.Fprintf(w, "\n")
}

func getPrefix(lvl []bool) string {
	var levelPrefix string
	var level = len(lvl)

	for i := 0; i < level; i++ {
		if level == 1 && lvl[i] {
			levelPrefix += fmt.Sprintf("└%s ", applyPadding("─"))
		} else if level == 1 && !lvl[i] {
			levelPrefix += fmt.Sprintf("├%s ", applyPadding("─"))
		} else if i+1 == level && !lvl[i] {
			levelPrefix += fmt.Sprintf("├%s ", applyPadding("─"))
		} else if i+1 == level && lvl[i] {
			levelPrefix += fmt.Sprintf("└%s ", applyPadding("─"))
		} else if lvl[i] {
			levelPrefix += fmt.Sprintf(" %s ", applyPadding(" "))
		} else {
			levelPrefix += fmt.Sprintf("│%s ", applyPadding(" "))
		}
	}

	return levelPrefix
}

const padding = 4
const paddingContent = "    "

func applyPadding(filler string) string {
	var fill string
	for i := 0; i < padding-2; i++ {
		fill += filler
	}
	return fill
}

const magic = "\x00"

func printCommand(c *Command, lvl []bool) (lines []string) {
	lines = append(lines, fmt.Sprintf("%s%s%s(%d,%d) %s %s", paddingContent, getPrefix(lvl), c.name, len(c.middlewarePre), len(c.middleware), magic, c.cc.GetShort()))
	var level = append(lvl, false)
	for i := 0; i < len(c.commands); i++ {
		if i+1 == len(c.commands) {
			level[len(level)-1] = true
		}
		subLines := printCommand(c.commands[i], level)
		lines = append(lines, subLines...)
	}
	return lines
}

func (c *Command) updateUsage(x *xconf.XConf) {
	c.usage = func() {
		c.Explain(c.Output)
		var bindFieldPathParent []string
		bindFieldPath := c.bindFieldPath
		if c.parent != nil {
			bindFieldPathParent = c.parent.bindFieldPath
			if c.parent.bind != nil && len(bindFieldPathParent) == 0 {
				bindFieldPathParent = xconf.FieldPathList(c.parent.bind, c.parent.newXConf())
			}
		}

		if c.bind != nil && len(bindFieldPath) == 0 {
			bindFieldPath = xconf.FieldPathList(c.bind, x)
		}
		local := c.flagLocal
		for _, v := range bindFieldPath {
			if xutil.ContainString(bindFieldPathParent, v) {
				continue
			}
			local = append(local, v)
		}
		var nowFlags []string
		c.FlagSet.VisitAll(func(f *flag.Flag) {
			nowFlags = append(nowFlags, f.Name)
		})
		var inherit []string
		for _, v := range nowFlags {
			if xutil.ContainString(local, v) {
				continue
			}
			inherit = append(inherit, v)
		}

		allFlag := xflag.GetFlagInfo(c.FlagSet)
		fieldPathInfoMap := make(map[string]xconf.StructFieldPathInfo)
		if c.bind != nil {
			fieldPathInfoMap = x.ZeroStructKeysTagList(reflect.New(reflect.ValueOf(c.bind).Type().Elem()).Interface())
		}

		var linesGlobal []string
		var linesLocal []string
		magic := "\x00"
		for _, v := range allFlag.List {
			line := fmt.Sprintf("--%s", v.Name)
			line += magic
			tag := "-"
			if c.bind != nil {
				tag = xconf.FlagTypeStr(x, v.Name)
			}
			line += v.TypeName
			line += magic
			usage := ""
			if info, ok := fieldPathInfoMap[v.Name]; ok {
				usage = info.Tag.Get("usage")
			}
			if usage == "" {
				usage = v.Usage
			}
			line += fmt.Sprintf("|%s| %s", tag, usage)
			if xutil.ContainString(inherit, v.Name) {
				linesGlobal = append(linesGlobal, line)
			} else if xutil.ContainString(local, v.Name) {
				linesLocal = append(linesLocal, line)
			} else {
				panic("invalid flag name : " + v.Name)
			}
		}
		heaerLine := "FLAG" + "\x00" + "TYPE" + "\x00" + "USAGE"
		var allLine []string
		allLine = append(allLine, heaerLine)
		allLine = append(allLine, linesGlobal...)
		allLine = append(allLine, linesLocal...)
		lineAllFormatted := xutil.TableFormatLines(allLine, magic)
		lineMaxLen := xutil.StringMaxLenByRune(lineAllFormatted)

		if len(linesGlobal) > 0 {
			fmt.Fprintf(c.Output, "OPTIONS GLOBAL:\n")
			fmt.Fprintln(c.Output, paddingContent+strings.Repeat("-", lineMaxLen))
			fmt.Fprintln(c.Output, paddingContent+lineAllFormatted[0])
			fmt.Fprintln(c.Output, paddingContent+strings.Repeat("-", lineMaxLen))
			sorted := lineAllFormatted[1 : len(linesGlobal)+1]
			sort.Strings(sorted)
			for i := 0; i < len(linesGlobal); i++ {
				fmt.Fprintln(c.Output, paddingContent+sorted[i])
			}

			fmt.Fprintln(c.Output, paddingContent+strings.Repeat("-", lineMaxLen))
			fmt.Fprintln(c.Output)
		}

		if len(linesLocal) > 0 {
			fmt.Fprintf(c.Output, "OPTIONS LOCAL:\n")
			fmt.Fprintln(c.Output, paddingContent+strings.Repeat("-", lineMaxLen))
			fmt.Fprintln(c.Output, paddingContent+lineAllFormatted[0])
			fmt.Fprintln(c.Output, paddingContent+strings.Repeat("-", lineMaxLen))
			sorted := lineAllFormatted[1+len(linesGlobal):]
			for i := 0; i < len(linesLocal); i++ {
				fmt.Fprintln(c.Output, paddingContent+sorted[i])
			}
			fmt.Fprintln(c.Output, paddingContent+strings.Repeat("-", lineMaxLen))
		}
		fmt.Fprintln(c.Output)
		fmt.Fprintf(c.Output, "Use \"%s [command] --help\" for more information about a command.\n", path.Base(os.Args[0]))
	}
}
