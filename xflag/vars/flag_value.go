package vars

import (
	"flag"
	"strings"
)

const PathPrefix = "xflag#"

func cleanPath(fieldPath string) string {
	if strings.HasPrefix(fieldPath, PathPrefix) {
		return fieldPath
	}
	return PathPrefix + fieldPath
}

var allProviders = make(map[string]func(valPtr interface{}) flag.Getter)

func SetProviderByFieldType(typeStr string, provider func(valPtr interface{}) flag.Getter) {
	allProviders[typeStr] = provider
}
func SetProviderByFieldPath(fieldPath string, provider func(valPtr interface{}) flag.Getter) {
	allProviders[cleanPath(fieldPath)] = provider
}

type FlagValue interface {
	flag.Getter
	Usage() string
}

type FlagValueProvider = func(fieldPath string, typeStr string, valPtr interface{}) (flag.Getter, bool)

// 优先通过filedPath匹配
var DefaultFlagValueProvider FlagValueProvider = func(fieldPath string, typeStr string, valPtr interface{}) (flag.Getter, bool) {
	provider, ok := allProviders[cleanPath(fieldPath)]
	if !ok {
		provider, ok = allProviders[typeStr]
		if !ok {
			return nil, false
		}
	}
	return provider(valPtr), true
}
