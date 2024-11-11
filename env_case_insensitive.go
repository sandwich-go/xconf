package xconf

import (
	"os"
	"strings"
)

// LookupEnvCaseInsensitive 查找环境变量，大小写不敏感
func LookupEnvCaseInsensitive(key string) (string, bool) {
	upperKey := strings.ToUpper(key)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 && strings.ToUpper(strings.TrimSpace(pair[0])) == upperKey {
			return pair[1], true
		}
	}
	return "", false
}
