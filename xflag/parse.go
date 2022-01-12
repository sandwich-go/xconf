package xflag

import (
	"regexp"
)

var (
	argumentRegex = regexp.MustCompile(`^\-{1,2}([\w\?\.\-]+)(=){0,1}(.*)$`)
)

// ParseArgsToMapStringString 简单的将arg按照flag的格式解析到map
func ParseArgsToMapStringString(args []string) map[string]string {
	parsedOptions := make(map[string]string)
	for i := 0; i < len(args); {
		match := argumentRegex.FindStringSubmatch(args[i])
		if len(match) > 2 {
			if match[2] == "=" {
				parsedOptions[match[1]] = match[3]
			} else if i < len(args)-1 {
				if len(args[i+1]) > 0 && args[i+1][0] == '-' {
					parsedOptions[match[1]] = match[3]
				} else {
					parsedOptions[match[1]] = args[i+1]
					i += 2
					continue
				}
			} else {
				parsedOptions[match[1]] = match[3]
			}
		}
		i++
	}
	return parsedOptions
}
