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

// TableFormatLines fotmat return lines
func TableFormatLines(lineAll []string, magic string) []string {
	ret := append([]string{}, lineAll...)
	for {
		maxLen := 0
		for _, line := range ret {
			sidx := strings.Index(line, magic)
			if sidx > maxLen {
				maxLen = sidx
			}
		}
		if maxLen == 0 {
			break
		}
		maxLen += 2
		for index, line := range ret {
			sidx := strings.Index(line, magic)
			spacing := strings.Repeat(" ", maxLen-sidx)
			line = line[:sidx] + spacing + line[sidx+1:]
			ret[index] = line
		}
	}
	sort.Strings(ret[1:])
	return ret
}

// TableFormat table格式化lineAll，对齐
func TableFormat(lineAll []string, magic string, seperateLine bool, suffixLines ...string) string {
	lineAllFormatted := TableFormatLines(lineAll, magic)
	buf := new(bytes.Buffer)
	lineMaxLen := StringMaxLen(lineAllFormatted, realLength)
	if seperateLine {
		fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	}
	for i, line := range lineAllFormatted {
		fmt.Fprintln(buf, line)
		if i == 0 {
			if seperateLine {
				fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
			}
			continue
		}
	}
	if seperateLine {
		fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	}
	suffixLines = StringSliceWalk(suffixLines, StringSliceEmptyFilter)
	if len(suffixLines) != 0 {
		for _, v := range suffixLines {
			fmt.Fprintln(buf, v)
		}
		if seperateLine {
			fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
		}
	}
	return buf.String()
}
