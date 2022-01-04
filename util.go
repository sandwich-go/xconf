package xconf

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"
)

type ErrorHandling int

const (
	ContinueOnError ErrorHandling = iota
	ExitOnError
	PanicOnError
)

var ErrorNeedParsedFirst = errors.New("should parsed first")
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

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
