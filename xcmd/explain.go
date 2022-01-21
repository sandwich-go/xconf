package xcmd

import (
	"fmt"
	"io"
	"sort"
	"strings"
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
		fmt.Fprintf(w, "Usage: %s <flags> <args>\n\n", strings.Join(c.usageNamePath, " "))
		return
	} else {
		fmt.Fprintf(w, "Usage: %s <subcommand> <flags> <args>\n\n", strings.Join(c.usageNamePath, " "))
	}
	sort.Sort(byGroupName(c.commands))
	fmt.Fprintf(w, "Available Commands:\n")
	sort.Sort(byGroupName(c.commands))
	var level = []bool{}
	lines := printCommand(c, level)
	// lines = xutil.TableFormatLines(lines, magic)
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

func applyPadding(filler string) string {
	var fill string
	for i := 0; i < padding-2; i++ {
		fill += filler
	}
	return fill
}

const magic = "\x00"

func printCommand(c *Command, lvl []bool) (lines []string) {
	lines = append(lines, fmt.Sprintf("%s%s(%d,%d) %s %s", getPrefix(lvl), c.name, len(c.middlewarePre), len(c.middleware), magic, c.cc.GetShort()))
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
