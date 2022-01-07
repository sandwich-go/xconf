package xutil

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// TableFormat table格式化lineAll，对齐
func TableFormat(lineAll []string, magic string, suffixLines ...string) string {
	for hasNextMagic := true; hasNextMagic; {
		maxLen := 0
		for _, line := range lineAll {
			sidx := strings.Index(line, magic)
			if sidx > maxLen {
				maxLen = sidx
			}
		}
		if maxLen == 0 {
			hasNextMagic = false
			break
		}
		maxLen += 3
		for index, line := range lineAll {
			sidx := strings.Index(line, magic)
			spacing := strings.Repeat(" ", maxLen-sidx)
			line = line[:sidx] + spacing + line[sidx+1:]
			lineAll[index] = line
		}
	}
	sort.Strings(lineAll[1:])
	buf := new(bytes.Buffer)
	lineMaxLen := StringMaxLen(lineAll)
	fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	for i, line := range lineAll {
		fmt.Fprintln(buf, line)
		if i == 0 {
			fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
			continue
		}
	}
	fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	if len(suffixLines) != 0 {
		for _, v := range suffixLines {
			fmt.Fprintln(buf, v)
		}
		fmt.Fprintln(buf, strings.Repeat("-", lineMaxLen))
	}
	return buf.String()
}
