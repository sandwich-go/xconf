package xconf

import (
	"io"

	"github.com/sandwich-go/xconf/kv"
)

var xx *XConf

func init() {
	xx = New()
}

func Parse(ccPtr interface{}, opts ...Option) error {
	xx.cc.ApplyOption(opts...)
	return xx.Parse(ccPtr)
}

// 返回指定配置的hash字符串
func HashStructure(v interface{}) (s string) { return xx.HashStructure(v) }

// 返回当前最新数据的hash字符串
func Hash() (s string) { return xx.Hash() }

// Latest将xconf内缓存的配置数据绑定到Parse时传入类型，逻辑层需要将返回的interface{}转换到相应的配置指针
func Latest() (interface{}, error) { return xx.Latest() }

// DumpInfo debug数据
func DumpInfo() { xx.DumpInfo() }

// NotifyUpdate 底层配置更新，返回的是全量的数据指针
func NotifyUpdate() <-chan interface{} { return xx.NotifyUpdate() }

// WatchUpdate confPath不会自动绑定env value,如果需要watch的路径与环境变量相关，先通过ParseEnvValue自行解析替换处理错误
func WatchUpdate(confPath string, loader kv.Loader) { xx.WatchUpdate(confPath, loader) }

// UpdateWithFieldPathValues 根据字段FieldPath更新数据, 支持的字段类型依赖于xflag,通过NotifyUpdate异步通知更新或通过Latest同步获取
func UpdateWithFieldPathValues(kv ...string) (err error) { return xx.UpdateWithFieldPathValues(kv...) }

// UpdateWithFlagArgs 提供FlagSet合法参数更新数据，通过NotifyUpdate异步通知更新或通过Latest同步获取
func UpdateWithFlagArgs(flagArgs ...string) (err error) { return xx.UpdateWithFlagArgs(flagArgs...) }

// UpdateWithEnviron 提供环境变量合法配置更新数据，通过NotifyUpdate异步通知更新或通过Latest同步获取
func UpdateWithEnviron(environ ...string) (err error) { return xx.UpdateWithEnviron(environ...) }

// UpdateWithFiles 通过文件更新配置，通过NotifyUpdate异步通知更新或通过Latest同步获取
func UpdateWithFiles(files ...string) (err error) { return xx.UpdateWithFiles(files...) }

// UpdateWithReader 通过reader更新配置，通过NotifyUpdate异步通知更新或通过Latest同步获取
func UpdateWithReader(readers ...io.Reader) (err error) { return xx.UpdateWithReader(readers...) }

// WatchFieldPath 关注特定的字段变化
func WatchFieldPath(fieldPath string, changed OnFieldUpdated) { xx.WatchFieldPath(fieldPath, changed) }

// SaveToFile 将内置解析的数据dump到文件，根据文件后缀选择codec
func SaveToFile(fileName string) error { return xx.SaveToFile(fileName) }

// SaveToWriter 将内置解析的数据dump到writer，类型为ct
func SaveToWriter(ct ConfigType, writer io.Writer) error { return xx.SaveToWriter(ct, writer) }

// SaveVarToFile 将外部传入的valPtr,写入到fileName中，根据文件后缀选择codec
func SaveVarToFile(valPtr interface{}, fileName string) error {
	return xx.SaveVarToFile(valPtr, fileName)
}

// SaveVarToWriter 将外部传入的valPtr,写入到writer中，类型为ct
func SaveVarToWriter(valPtr interface{}, ct ConfigType, writer io.Writer) error {
	return xx.SaveVarToWriter(valPtr, ct, writer)
}

func MustSaveToFile(f string)                                       { xx.MustSaveToFile(f) }
func MustSaveToWriter(ct ConfigType, writer io.Writer)              { xx.MustSaveToWriter(ct, writer) }
func MustSaveVarToFile(v interface{}, f string)                     { xx.MustSaveVarToFile(v, f) }
func MustSaveVarToWriter(v interface{}, ct ConfigType, w io.Writer) { xx.SaveVarToWriter(v, ct, w) }
func MustAsBytes(ct ConfigType) []byte                              { return xx.MustAsBytes(ct) }