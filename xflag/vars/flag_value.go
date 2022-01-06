package vars

import (
	"flag"
	"strings"
)

const pathPrefix = "xflag#"

func cleanPath(fieldPath string) string {
	if strings.HasPrefix(fieldPath, pathPrefix) {
		return fieldPath
	}
	return pathPrefix + fieldPath
}

var allProviders = make(map[string]func(valPtr interface{}) flag.Getter)

// SetProviderByFieldType 设定type名称的flag.Getter获取方法
func SetProviderByFieldType(typeStr string, provider func(valPtr interface{}) flag.Getter) {
	allProviders[typeStr] = provider
}

// SetProviderByFieldPath 设定fieldPath指向字段的flag.Getter获取方法
func SetProviderByFieldPath(fieldPath string, provider func(valPtr interface{}) flag.Getter) {
	allProviders[cleanPath(fieldPath)] = provider
}

// FlagValueProvider 由fieldPath与typeStr以及数值的指针返回对应的FlagValue
type FlagValueProvider = func(fieldPath string, typeStr string, valPtr interface{}) (flag.Getter, bool)

// DefaultFlagValueProvider 优先通过filedPath匹配
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
