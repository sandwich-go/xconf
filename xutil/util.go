package xutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ReadAll read all bytes for golang 1.14 15
func ReadAll(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}

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

// StringSliceEmptyFilter StringSliceWalk的filter，过滤空字符串
var StringSliceEmptyFilter = func(s string) (string, bool) { return s, s != "" }

// StringSliceWalk 遍历vs,将f应用到每一个元素，返回更新后的数据
func StringSliceWalk(vs []string, f func(string) (string, bool)) []string {
	vsm := make([]string, 0)
	for _, v := range vs {
		ret, valid := f(v)
		if valid {
			vsm = append(vsm, ret)
		}
	}
	return vsm
}

// ToCleanStringSlice 分割字符串，trim字符
func ToCleanStringSlice(in string) []string {
	return StringSliceWalk(strings.Split(StringTrim(in), ","), func(s string) (string, bool) { return StringTrim(s), true })
}

// ContainAtLeastOneEqualFold s1是否至少含有s2中的一个元素
func ContainAtLeastOneEqualFold(s1 []string, s2 []string) bool {
	for _, v := range s2 {
		if ContainStringEqualFold(s1, v) {
			return true
		}
	}
	return false
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase 将指定的str返回SnakeCase类型
func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// StringMaxLen ss中的字符串最长值
func StringMaxLen(ss []string, lenFunc func(s string) int) (max int) {
	for _, v := range ss {
		if lenVal := lenFunc(v); lenVal > max {
			max = lenVal
		}
	}
	return max
}

// ContainString 是否含有字符串
func ContainString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainStringEqualFold 是否含有字符串不区分大小写
func ContainStringEqualFold(s []string, v string) bool {
	for _, vv := range s {
		if strings.EqualFold(vv, v) {
			return true
		}
	}
	return false
}

// PanicErr err不为nil则panic
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// WrapIfErr err不为nil则wrap
func WrapIfErr(err error, fmtStr string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(fmtStr, args...)
}

// WrapIfErrAsFisrt err不为nil则wrap，将err作为第一个fmt的参数
func WrapIfErrAsFisrt(err error, fmtStr string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	var argList []interface{}
	argList = append(argList, err)
	argList = append(argList, args...)
	return fmt.Errorf(fmtStr, argList...)
}

// PanicErrWithWrap err不为nil则panic，error类型使用fmtStr格式化
func PanicErrWithWrap(err error, fmtStr string, args ...interface{}) {
	if err != nil {
		PanicErr(fmt.Errorf(fmtStr, args...))
	}
}

func kv2map(kv ...string) (map[string]string, error) {
	ret := make(map[string]string)
	return ret, KVListApplyFunc(func(k, v string) bool {
		ret[k] = v
		return true
	}, kv...)
}

// KVListToFlagArgs 将kv转换为Flag格式字符串列表
func KVListToFlagArgs(kv ...string) ([]string, error) {
	var ret []string
	return ret, KVListApplyFunc(func(k, v string) bool {
		ret = append(ret, fmt.Sprintf("--%s=%s", k, v))
		return true
	}, kv...)
}

func kv2Environ(kv ...string) ([]string, error) {
	var ret []string
	return ret, KVListApplyFunc(func(k, v string) bool {
		ret = append(ret, fmt.Sprintf("%s=%s", k, v))
		return true
	}, kv...)
}

var _ = kv2Environ
var _ = kv2map

// KVListApplyFunc kv利用给定的f进行k，v遍历
func KVListApplyFunc(f func(k, v string) bool, kv ...string) error {
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
