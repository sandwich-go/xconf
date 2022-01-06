package xconf

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// envBindToFlags 将environ绑定到FlagSet格式,字段映射由mapping指定
func envBindToFlags(environ []string, mapping map[string]string) []string {
	var flags []string
	env := make(map[string]string)
	for _, s := range environ {
		i := strings.Index(s, "=")
		if i < 1 {
			continue
		}
		env[s[0:i]] = s[i+1:]
	}
	for k, v := range env {
		flagName := mapping[strings.ToUpper(k)]
		if flagName == "" {
			continue
		}
		flags = append(flags, fmt.Sprintf("-%s=%s", flagName, v))
	}
	return flags
}

var _ = falgBindEnv

func falgBindEnv(fs *flag.FlagSet, environ []string, mapping map[string]string) error {
	env := make(map[string]string)
	for _, s := range environ {
		i := strings.Index(s, "=")
		if i < 1 {
			continue
		}
		env[s[0:i]] = s[i+1:]
	}
	for k, v := range env {
		flagName := mapping[strings.ToUpper(k)]
		if flagName == "" {
			continue
		}
		flagInst := fs.Lookup(flagName)
		if flagInst == nil {
			continue
		}
		err := flagInst.Value.Set(v)
		if err != nil {
			return fmt.Errorf("got error while bind env:%s val:%s to flag,err :%w", k, v, err)
		}
	}
	return nil
}

// ValueGetter Env value provider func.
var ValueGetter = os.LookupEnv

// parse env value, allow:
// 	only key 	 - "${SHELL}"
// 	with default - "${NotExist|defValue}"
//	multi key 	 - "${GOPATH}/${APP_ENV | prod}/dir"
// Notice:
//  must add "?" - To ensure that there is no greedy match
//  var envRegex = regexp.MustCompile(`\${[\w-| ]+}`)
var envRegex = regexp.MustCompile(`\${.+?}`)

// ParseEnvValue parse ENV var value from input string, support default value.
func ParseEnvValue(val string, errEnvBindNotExistWithoutDefault bool) (newVal string, err error) {
	if !strings.Contains(val, "${") {
		return val, nil
	}
	var name, def string
	var valAndDefaultNotFound []string
	newVal = envRegex.ReplaceAllStringFunc(val, func(eVar string) string {
		// eVar like "${NotExist|defValue}", first remove "${" and "}", then split it
		ss := strings.SplitN(eVar[2:len(eVar)-1], "|", 2)
		// with default value. ${NotExist|defValue}
		hasDefault := false
		if len(ss) == 2 {
			name, def = strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1])
			hasDefault = true
		} else {
			def = eVar // use raw value
			name = strings.TrimSpace(ss[0])
		}
		// get ENV value by name
		eVal, ok := ValueGetter(name)
		if eVal == "" {
			if !ok && !hasDefault {
				// 指定的key不存在且没有显示提供默认值
				valAndDefaultNotFound = append(valAndDefaultNotFound, name)
			}
			eVal = def
		}
		return eVal
	})
	if len(valAndDefaultNotFound) == 0 {
		return newVal, nil
	}
	return newVal, fmt.Errorf("EnvBind lost env and default for key:%s", strings.Join(valAndDefaultNotFound, ","))
}
