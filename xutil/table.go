package xutil

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
)

// stripAnsiEscapeRegexp is a regular expression to clean ANSI Control sequences
// feat https://stackoverflow.com/questions/14693701/how-can-i-remove-the-ansi-escape-sequences-from-a-string-in-python#33925425
var stripAnsiEscapeRegexp = regexp.MustCompile(`(\x9B|\x1B\[)[0-?]*[ -/]*[@-~]`)

// stripAnsiEscape returns string without ANSI escape sequences (colors etc)
func stripAnsiEscape(s string) string {
	return stripAnsiEscapeRegexp.ReplaceAllString(s, "")
}

// realWidth returns real string length (without ANSI escape sequences)
func realLength(s string) int {
	return runewidth.StringWidth(stripAnsiEscape(s))
}

// TableFormat table格式化lineAll，对齐
func TableFormat(lineAll []string, magic string, suffixLines ...string) string {
	for {
		maxLen := 0
		for _, line := range lineAll {
			sidx := strings.Index(line, magic)
			if sidx > maxLen {
				maxLen = sidx
			}
		}
		if maxLen == 0 {
			break
		}
		maxLen += 2
		for index, line := range lineAll {
			sidx := strings.Index(line, magic)
			spacing := strings.Repeat(" ", maxLen-sidx)
			line = line[:sidx] + spacing + line[sidx+1:]
			lineAll[index] = line
		}
	}
	sort.Strings(lineAll[1:])
	buf := new(bytes.Buffer)
	lineMaxLen := StringMaxLen(lineAll, realLength)
	fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	for i, line := range lineAll {
		fmt.Fprintln(buf, line)
		if i == 0 {
			fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
			continue
		}
	}
	fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	suffixLines = StringSliceWalk(suffixLines, StringSliceEmptyFilter)
	if len(suffixLines) != 0 {
		for _, v := range suffixLines {
			fmt.Fprintln(buf, v)
		}
		fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	}
	return buf.String()
}
