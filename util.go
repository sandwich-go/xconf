package xconf

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"
)

// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
var DefaultTrimChars = string([]byte{
	'\t', // Tab.
	'\v', // Vertical tab.
	'\n', // New line (line feed).
	'\r', // Carriage return.
	'\f', // New page.
	' ',  // Ordinary space.
	0x00, // NUL-byte.
	0x85, // Delete.
	0xA0, // Non-breaking space.
})

// StringTrim trim字符串
func StringTrim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}

// StringMap 便于vs,将f应用到每一个元素，返回更新后的数据
func StringMap(vs []string, f func(string) (string, bool)) []string {
	vsm := make([]string, 0)
	for _, v := range vs {
		ret, valid := f(v)
		if valid {
			vsm = append(vsm, ret)
		}
	}
	return vsm
}

func toCleanStringSlice(in string) []string {
	return StringMap(strings.Split(StringTrim(in), ","), func(s string) (string, bool) { return StringTrim(s), true })
}

func containAtLeastOneEqualFold(s1 []string, s2 []string) bool {
	for _, v := range s2 {
		if containStringEqualFold(s1, v) {
			return true
		}
	}
	return false
}

// ErrorHandling 错误处理类型
type ErrorHandling int

const (
	// ContinueOnError 发生错误继续运行，Parse会返回错误
	ContinueOnError ErrorHandling = iota
	// ExitOnError 发生错误后退出
	ExitOnError
	// PanicOnError 发生错误后主动panic
	PanicOnError
)

var errorNeedParsedFirst = errors.New("should parsed first")
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase 将指定的str返回SnakeCase类型
func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func newFlagSet(name string) *flag.FlagSet { return flag.NewFlagSet(name, flag.ContinueOnError) }

func containString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func containStringEqualFold(s []string, v string) bool {
	for _, vv := range s {
		if strings.EqualFold(vv, v) {
			return true
		}
	}
	return false
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func wrapIfErr(err error, fmtStr string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(fmtStr, args...)
}

func wrapIfErrAsFisrt(err error, fmtStr string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	var argList []interface{}
	argList = append(argList, err)
	argList = append(argList, args...)
	return fmt.Errorf(fmtStr, argList...)
}

func panicErrWithWrap(err error, fmtStr string, args ...interface{}) {
	if err != nil {
		panicErr(fmt.Errorf(fmtStr, args...))
	}
}

func kv2map(kv ...string) (map[string]string, error) {
	ret := make(map[string]string)
	return ret, kvWithFunc(func(k, v string) bool {
		ret[k] = v
		return true
	}, kv...)
}

func kv2FlagArgs(kv ...string) ([]string, error) {
	var ret []string
	return ret, kvWithFunc(func(k, v string) bool {
		ret = append(ret, fmt.Sprintf("--%s=%s", k, v))
		return true
	}, kv...)
}

func kv2Environ(kv ...string) ([]string, error) {
	var ret []string
	return ret, kvWithFunc(func(k, v string) bool {
		ret = append(ret, fmt.Sprintf("%s=%s", k, v))
		return true
	}, kv...)
}

var _ = kv2Environ
var _ = kv2map

func kvWithFunc(f func(k, v string) bool, kv ...string) error {
	if len(kv)%2 == 1 {
		return errors.New("got the odd number of input pairs")
	}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s
			continue
		}
		if !f(key, s) {
			break
		}
	}
	return nil
}
