package xconf

import (
	"fmt"
	"io"
)

// UpdateWithFiles 提供files更新数据, 支持的字段类型依赖于xflag
func (x *XConf) UpdateWithFiles(files ...string) (err error) {
	return x.commonUpdateAndNotify(func() error {
		return x.updateDstDataWithFiles(files...)
	})
}

// UpdateWithReader 提供files更新数据, 支持的字段类型依赖于xflag
func (x *XConf) UpdateWithReader(readers ...io.Reader) (err error) {
	return x.commonUpdateAndNotify(func() error {
		return x.updateDstDataWithReaders(readers...)
	})
}

// UpdateWithFieldPathValues 根据字段FieldPath更新数据, 支持的字段类型依赖于xflag
func (x *XConf) UpdateWithFieldPathValues(kv ...string) (err error) {
	args, err := kv2FlagArgs(kv...)
	if err != nil {
		return fmt.Errorf("kv2FlagArgs with error:%v", err)
	}
	return x.UpdateWithFlagArgs(args...)
}

// UpdateWithFlagArgs 提供FlagSet合法参数更新数据，异步通知更新
func (x *XConf) UpdateWithFlagArgs(flagArgs ...string) (err error) {
	return x.commonUpdateAndNotify(func() error {
		return x.updateDstDataWithFlagSet(newFlagSet("UpdateWithFlagArgs"), flagArgs...)
	})
}

// UpdateWithEnviron 提供环境变量合法配置更新数据，异步通知更新
func (x *XConf) UpdateWithEnviron(environ ...string) (err error) {
	return x.commonUpdateAndNotify(func() error {
		return x.updateDstDataWithEnviron(environ...)
	})
}

func (x *XConf) commonUpdateAndNotify(f func() error) (err error) {
	x.dynamicUpdate.Lock()
	defer x.dynamicUpdate.Unlock()
	if !x.hasParsed {
		return errorNeedParsedFirst
	}
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("%v", reason)
		}
	}()
	err = f()
	panicErr(err)
	return wrapIfErrAsFisrt(x.notifyChanged(), "got error while notifyChanged:%v")
}
